# qpass - Secure Password Generator

A simple, fast CLI tool for generating cryptographically secure passwords.

## Features

- Generates cryptographically secure passwords using `crypto/rand`
- Automatically copies passwords to clipboard
- Cross-platform support (macOS, Linux, Windows)
- Customizable password length
- Option to hide password output

## Installation

### From Source

```bash
go install github.com/koushikyemula/qpass@latest
```

### Build Locally

```bash
git clone https://github.com/koushikyemula/qpass
cd qpass
go build -o qpass
```

## Usage

### Basic Usage

Generate an 10-character password:

```bash
qpass
```

### Options

- `-n <length>`: Set password length (default: 10)
- `-h`: Hide password output in terminal (still copies to clipboard)

### Examples

```bash
# Generate a 12-character password
qpass -n 12

# Generate password without showing it in terminal
qpass -h

# Combine options
qpass -n 16 -h
```

## Requirements

### Clipboard Utilities

- **macOS**: Uses `pbcopy` (built-in)
- **Linux**: Requires `xclip` or `xsel`

  ```bash
  # Ubuntu/Debian
  sudo apt-get install xclip
  # or
  sudo apt-get install xsel
  ```

- **Windows**: Uses `clip` (built-in)

## Character Set

The password generator uses the following character set by default:

- Lowercase letters: `a-z`
- Uppercase letters: `A-Z`
- Numbers: `0-9`
- Special characters: `!@#$%^&*()-_=+[]{}|;:,.<>?/`

## Security

- Uses Go's `crypto/rand` package for cryptographically secure random number generation
- No password storage or logging
- Passwords are generated fresh each time

## License

MIT License
