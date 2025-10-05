# BIP38CLI - Bitcoin Private Key Encryption Tool

A modern command-line tool for BIP38 (Bitcoin Improvement Proposal 38) private key encryption and decryption, written in Go.

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go Version](https://img.shields.io/badge/Go-1.24%2B-blue.svg)](https://golang.org)
[![Release](https://img.shields.io/github/release/yourusername/bip38cli.svg)](https://github.com/yourusername/bip38cli/releases)

## Features

- **Encrypt/Decrypt Bitcoin private keys** using BIP38 standard
- **Generate intermediate passphrase codes** for two-factor encryption
- **Support for compressed and uncompressed keys**
- **Secure passphrase handling** with hidden input
- **Fast and efficient** - built with Go and BTCSuite
- **Cross-platform support** (Linux, macOS, Windows)
- **Shell completion** for bash, zsh, fish, and PowerShell

## Quick Start

### Installation

**Download pre-built binary:**
```bash
# Linux/macOS
curl -LO https://github.com/yourusername/bip38cli/releases/latest/download/bip38cli-linux-amd64.tar.gz
tar -xzf bip38cli-*.tar.gz
sudo mv bip38cli /usr/local/bin/

# Verify installation
bip38cli --version
```

**Build from source:**
```bash
git clone https://github.com/yourusername/bip38cli.git
cd bip38cli
make build
./bin/bip38cli --version
```

### Basic Usage

**Encrypt a private key:**
```bash
bip38cli encrypt 5HueCGU8rMjxEXxiPuD5BDku4MkFqeZyd4dZ1jvhTVqvbTLvyTJ
# Enter passphrase when prompted
# Output: 6PRVWUbkzzsbcVac2qwfssoUJAN1Xhrg6bNk8J7Nzm5H7kxEbn2Nh2ZoGg
```

**Decrypt a private key:**
```bash
bip38cli decrypt 6PRVWUbkzzsbcVac2qwfssoUJAN1Xhrg6bNk8J7Nzm5H7kxEbn2Nh2ZoGg
# Enter passphrase when prompted
# Output: 5HueCGU8rMjxEXxiPuD5BDku4MkFqeZyd4dZ1jvhTVqvbTLvyTJ
```

## Commands

### `encrypt` - Encrypt Private Keys
```bash
# Interactive mode (most secure)
bip38cli encrypt

# Direct mode
bip38cli encrypt [WIF_PRIVATE_KEY]

# Force compressed format
bip38cli encrypt --compressed [WIF_PRIVATE_KEY]
```

### `decrypt` - Decrypt Private Keys
```bash
# Interactive mode
bip38cli decrypt

# Direct mode
bip38cli decrypt [ENCRYPTED_KEY]

# Show Bitcoin address
bip38cli decrypt --show-address [ENCRYPTED_KEY]
```

### `intermediate` - Two-Factor Encryption
```bash
# Generate intermediate code
bip38cli intermediate generate

# Generate with lot/sequence numbers
bip38cli intermediate generate --lot 123 --sequence 456

# Validate intermediate code
bip38cli intermediate validate [INTERMEDIATE_CODE]
```

## Shell Completion

BIP38CLI supports auto-completion for bash, zsh, fish, and PowerShell.

### Bash
```bash
# Add to ~/.bashrc
echo 'source <(bip38cli completion bash)' >> ~/.bashrc
source ~/.bashrc
```

### Zsh
```bash
# Add to ~/.zshrc
echo 'source <(bip38cli completion zsh)' >> ~/.zshrc
source ~/.zshrc
```

### Fish
```bash
# Generate completion file
bip38cli completion fish > ~/.config/fish/completions/bip38cli.fish
```

### PowerShell
```powershell
# Add to PowerShell profile
bip38cli completion powershell | Out-String | Invoke-Expression
```

## Configuration

Create a config file at `~/.bip38cli.yaml`:

```yaml
# Default behavior settings
verbose: false
compressed: true

# Output format preferences
output:
  format: "text"  # text, json
  colors: true
```

## What is BIP38?

BIP38 is a Bitcoin Improvement Proposal that defines a method for encrypting Bitcoin private keys with a passphrase. This enables:

1. **Password-protected storage** - Keys can be safely stored or transmitted
2. **Two-factor security** - Using intermediate codes for secure key generation
3. **Standardized format** - Compatible with other BIP38 implementations
4. **Safe backups** - Encrypted keys can be stored in multiple locations

## Security Best Practices

- **Use strong passphrases** (15+ characters or 6+ word phrases)
- **Always test decryption** before storing encrypted keys
- **Use offline environments** for critical operations
- **Store passphrases separately** from encrypted keys
- **Make multiple backups** of encrypted keys
- **Verify addresses** after decryption

- **Never use weak passwords** (123456, password, etc.)
- **Never store passphrases with encrypted keys**
- **Never use on infected computers**
- **Never share unencrypted private keys**

## Examples

### Secure Wallet Backup
```bash
# 1. Export private key from your wallet
# 2. Encrypt with strong passphrase
bip38cli encrypt --compressed
# 3. Test decryption immediately
bip38cli decrypt --show-address
# 4. Store encrypted key and passphrase separately
```

### Batch Processing
```bash
# Encrypt multiple keys from file
while read -r key; do
    echo "Processing: $key"
    echo "$key" | bip38cli encrypt
done < private_keys.txt
```

### Digital Inheritance
```bash
# Create encrypted backup for family
bip38cli encrypt 5HueCGU8rMjxEXxiPuD5BDku4MkFqeZyd4dZ1jvhTVqvbTLvyTJ
# Use memorable passphrase family can deduce
# Leave clear instructions for decryption
```

## Development

### Requirements
- Go 1.24.0 or later
- Make

### Building
```bash
# Install dependencies
make deps

# Run tests
make test

# Build binary
make build

# Run all checks
make all
```

### Testing
```bash
# Run tests with coverage
make test-coverage

# Run linting
make lint

# Benchmark performance
go test -bench=. ./internal/bip38/
```

## API Reference

### Exit Codes
- `0` - Success
- `1` - General error
- `2` - Invalid arguments
- `3` - Encryption/decryption failed
- `4` - Invalid key format

### Environment Variables
- `BIP38CLI_CONFIG` - Config file path
- `BIP38CLI_VERBOSE` - Enable verbose output
- `NO_COLOR` - Disable colored output

## Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Make your changes
4. Add tests for new functionality
5. Run `make all` to verify everything works
6. Commit your changes (`git commit -m 'Add amazing feature'`)
7. Push to the branch (`git push origin feature/amazing-feature`)
8. Open a Pull Request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Acknowledgments

- Built with [BTCSuite](https://github.com/btcsuite) - the industry standard Bitcoin library for Go
- Uses [Cobra](https://github.com/spf13/cobra) for CLI framework
- Implements the [BIP38 specification](https://github.com/bitcoin/bips/blob/master/bip-0038.mediawiki)

## Support

- **Documentation**: [Complete tutorials](TUTORIAL-EN.md)
- **Issues**: [GitHub Issues](https://github.com/yourusername/bip38cli/issues)
- **Discussions**: [GitHub Discussions](https://github.com/yourusername/bip38cli/discussions)

---

**⚠️ Important**: This software handles Bitcoin private keys. Always test with small amounts first and ensure you understand the risks. The authors are not responsible for any loss of funds.