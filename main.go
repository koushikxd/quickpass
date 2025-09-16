package main

import (
	"crypto/rand"
	"flag"
	"fmt"
	"log"
	"math/big"
	"os/exec"
	"runtime"
	"strings"
)

const defaultLength = 10
const defaultCharset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%^&*()-_=+[]{}|;:,.<>?/"

func main() {
	length := flag.Int("n", defaultLength, "Password length")
	hide := flag.Bool("h", false, "Hide password output in terminal")
	flag.Parse()

	charset := defaultCharset

	password, err := generatePassword(*length, charset)
	if err != nil {
		log.Fatalf("Error generating password: %v", err)
	}

	err = copyToClipboard(password)
	if err != nil {
		log.Printf("Warning: Could not copy to clipboard: %v", err)
	} 

	if !*hide {
		fmt.Println(password)
	}
}

func generatePassword(length int, charset string) (string, error) {
	if length <= 0 {
		return "", fmt.Errorf("password length must be positive")
	}

	password := make([]byte, length)
	charsetLen := big.NewInt(int64(len(charset)))

	for i := 0; i < length; i++ {
		randomIndex, err := rand.Int(rand.Reader, charsetLen)
		if err != nil {
			return "", fmt.Errorf("failed to generate random number: %v", err)
		}
		password[i] = charset[randomIndex.Int64()]
	}

	return string(password), nil
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
