package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"github.com/sashabaranov/go-openai"
)

const defaultLength = 24

func main() {
	length := flag.Int("n", defaultLength, "Password length")
	hide := flag.Bool("h", false, "Hide password output in terminal")
	flag.Parse()

	err := godotenv.Load(".env.local")
	if err != nil {
		log.Printf("Warning: Could not load .env.local file: %v", err)
	}

	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		log.Fatalf("OPENAI_API_KEY environment variable is required")
	}

	fmt.Println("Asking AI for help...")
	
	password, err := generateAIPassword(*length, apiKey)
	if err != nil {
		log.Fatalf("Error generating AI-powered password: %v", err)
	}

	err = copyToClipboard(password)
	if err != nil {
		log.Printf("Warning: Could not copy to clipboard: %v", err)
	}

	if !*hide {
		fmt.Println(password)
		fmt.Println("Don't show this to anyone else")
	} else {
		fmt.Println("Password copied to clipboard")
	}
}

func generateAIPassword(length int, apiKey string) (string, error) {
	if length <= 0 {
		return "", fmt.Errorf("password length must be positive")
	}

	client := openai.NewClient(apiKey)

	prompt := fmt.Sprintf(`You are an expert password generator. Please generate exactly ONE secure password that is EXACTLY %d characters long. 

Requirements:
- Must be exactly %d characters
- Include uppercase letters, lowercase letters, numbers, and special characters
- Be cryptographically secure and random
- Do NOT include any explanations, quotes, or additional text
- Respond ONLY with the password itself

Generate the password now:`, length, length)

	done := make(chan bool)
	var resp openai.ChatCompletionResponse
	var err error

	go func() {
		resp, err = client.CreateChatCompletion(
			context.Background(),
			openai.ChatCompletionRequest{
				Model: "gpt-4o-mini",
				Messages: []openai.ChatCompletionMessage{
					{
						Role:    openai.ChatMessageRoleUser,
						Content: prompt,
					},
				},
				MaxTokens:   100,
				Temperature: 1.0,
			},
		)
		done <- true
	}()

	dots := ""
	for {
		select {
		case <-done:
			fmt.Print("\n")
			goto finished
		case <-time.After(500 * time.Millisecond):
			dots += "."
			if len(dots) > 3 {
				dots = ""
			} else {
				fmt.Print(".")
			}
		}
	}

finished:
	if err != nil {
		return "", fmt.Errorf("failed to generate AI password: %v", err)
	}

	if len(resp.Choices) == 0 {
		return "", fmt.Errorf("no response from AI")
	}

	password := strings.TrimSpace(resp.Choices[0].Message.Content)
	
	if len(password) > length {
		password = password[:length]
	} else if len(password) < length {
		for len(password) < length {
			password += "!"
		}
	}

	return password, nil
}

func copyToClipboard(text string) error {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("pbcopy")
	case "linux":
		if _, err := exec.LookPath("xclip"); err == nil {
			cmd = exec.Command("xclip", "-selection", "clipboard")
		} else if _, err := exec.LookPath("xsel"); err == nil {
			cmd = exec.Command("xsel", "--clipboard", "--input")
		} else {
			return fmt.Errorf("no clipboard utility found (install xclip or xsel)")
		}
	case "windows":
		cmd = exec.Command("clip")
	default:
		return fmt.Errorf("unsupported operating system: %s", runtime.GOOS)
	}

	if cmd == nil {
		return fmt.Errorf("could not create clipboard command")
	}

	cmd.Stdin = strings.NewReader(text)
	return cmd.Run()
}
