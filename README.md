# BIP38CLI - Bitcoin Private Key Encryption Tool

A focused command-line application that implements the [BIP38](https://github.com/bitcoin/bips/blob/master/bip-0038.mediawiki) standard to encrypt and decrypt Bitcoin private keys with passphrase protection.

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go Version](https://img.shields.io/badge/Go-1.24%2B-blue.svg)](https://go.dev/)
[![Release](https://img.shields.io/github/release/carlosrabelo/bip38cli.svg)](https://github.com/carlosrabelo/bip38cli/releases)

## Highlights

- Encrypt and decrypt Wallet Import Format (WIF) keys using spec-compliant BIP38 routines
- Generate and validate intermediate codes for two-factor key creation flows
- Zero out passphrase buffers as soon as possible to reduce memory exposure
- Hidden terminal input for passphrases with compression toggles and verbose insights
- Smart configuration discovery: `~/.bip38cli.yaml` → `./bip38cli.yaml` → `/etc/bip38cli/config.yaml`
- Shell completion generation for bash, zsh, fish, and PowerShell

## Project Layout

```
core/
  cmd/bip38cli/             # Go entry point for the CLI binary
  internal/app/cli/         # Cobra commands and user interaction flows
  internal/domain/bip38/    # BIP38 domain logic and tests
  pkg/                      # Reserved for future public packages
  Makefile                  # Go-specific build helpers
bin/
  .gitkeep                  # Binary output directory placeholder
README-PT.md              # Portuguese documentation
docs/
  TUTORIAL-EN.md            # English tutorial and walkthroughs
  TUTORIAL-PT.md            # Portuguese tutorial and walkthroughs
scripts/
  install.sh                # Binary installer helper
  uninstall.sh              # Binary removal helper
Makefile                    # Root orchestration makefile
```

## Quick Start

### Build from Source

```bash
git clone https://github.com/carlosrabelo/bip38cli.git
cd bip38cli
make build
./bin/bip38cli --version
```

### Install the Binary

Install to `$HOME/.local/bin` (recommended for non-root users):

```bash
./scripts/install.sh --user
```

Install to `/usr/local/bin` (requires appropriate permissions):

```bash
sudo ./scripts/install.sh
```

Remove the binary later with the matching uninstall script:

```bash
./scripts/uninstall.sh --user
# or
sudo ./scripts/uninstall.sh
```

### Run via Docker

```bash
# Show CLI help inside the container
./scripts/bip38cli-docker.sh --help

# Execute commands without installing Go locally
./scripts/bip38cli-docker.sh encrypt --verbose
```

The helper script keeps Docker artefacts under `docker/` and will build a local image on demand. For advanced scenarios see `docker/README.md`.

## Usage

### Encrypt a WIF Key

```bash
bip38cli encrypt KwYgW8gcxj1JWJXhPSu4Fqwzfhp5Yfi42mdYmMa4XqK7NJxXUSK7
# Hidden prompts ask for passphrase twice
# Result: 6PRV...
```

### Decrypt an Encrypted Key

```bash
bip38cli decrypt 6PRVWUbkzzsbcVac2qwfssoUJAN1Xhrg6bNk8J7Nzm5H7kxEbn2Nh2ZoGg
# Hidden prompt asks for passphrase
# Optional: --show-address to print derived address
```

### Work with Intermediate Codes

```bash
# Generate code with lot/sequence metadata
bip38cli intermediate generate --lot 123 --sequence 456 --use-lot-sequence

# Validate a provided code
bip38cli intermediate validate passphraseabc123...
```

Generate shell completions for your environment:

```bash
bip38cli completion bash       > /usr/share/bash-completion/completions/bip38cli
bip38cli completion zsh        > /usr/share/zsh/site-functions/_bip38cli
bip38cli completion fish       > ~/.config/fish/completions/bip38cli.fish
bip38cli completion powershell | Out-String | Invoke-Expression
```

## Configuration

BIP38CLI reads configuration using Viper with the following precedence:

1. `--config /path/to/file.yaml`
2. `~/.bip38cli.yaml`
3. `./bip38cli.yaml`
4. `/etc/bip38cli/config.yaml`

Default values baked into the binary:

```yaml
defaults:
  compressed: true
output:
  format: text
  colors: true
```

Set `verbose: true` to display the config path in use and additional diagnostic lines.

## Documentation

- [English tutorial](docs/TUTORIAL-EN.md)
- [Portuguese tutorial](docs/TUTORIAL-PT.md)

## Development

Run the following commands from the repository root:

```bash
make build      # Compile the CLI into bin/bip38cli
make fmt        # Format Go sources with gofmt
make test       # Execute go test ./...
make lint       # Run golangci-lint if available
make clean      # Remove build artifacts
```

The Go module lives at `github.com/carlosrabelo/bip38cli/core`. Tests for the BIP38 domain perform real scrypt key derivations and can take a few seconds.

## Security Notes

- Use strong passphrases (15+ characters or multi-word phrases)
- Test decryption immediately after encryption before storing backups
- Keep passphrases separate from encrypted keys and avoid networked copy-paste tools
- Prefer air-gapped machines for large batches or high-value wallets
- Treat intermediate codes with the same care as encrypted keys

## Donations

If BIP38CLI is useful to you, consider supporting development:

**BTC**: `bc1qw2raw7urfuu2032uyyx9k5pryan5gu6gmz6exm`  
**ETH**: `0xdb4d2517C81bE4FE110E223376dD9B23ca3C762E`  
**SOL**: `A3tNpXSb8rHw2PJYALQeZzwvR4pRWk72YwJdeXGKmS1q`  
**TRX**: `TTznF3FeDCqLmL5gx8GingeahUyLsJJ68A`

## License

Distributed under the [MIT License](LICENSE).
