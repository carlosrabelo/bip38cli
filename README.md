# BIP38CLI - Bitcoin Private Key Encryption Tool

A focused command-line application that implements the [BIP38](https://github.com/bitcoin/bips/blob/master/bip-0038.mediawiki) standard to encrypt and decrypt Bitcoin private keys with passphrase protection.

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go Version](https://img.shields.io/badge/Go-1.25%2B-blue.svg)](https://go.dev/)
[![Release](https://img.shields.io/github/release/carlosrabelo/bip38cli.svg)](https://github.com/carlosrabelo/bip38cli/releases)

## Highlights

- Encrypt and decrypt Wallet Import Format (WIF) keys using spec-compliant BIP38 routines
- Generate fresh WIF keys for any Bitcoin network with optional BIP38 encryption
- Display native SegWit (BIP84) bech32 addresses for compressed keys, with legacy fallback for uncompressed WIFs
- Generate and validate intermediate codes for two-factor key creation flows
- Zero out passphrase buffers as soon as possible to reduce memory exposure
- Hidden terminal input for passphrases with compression toggles and verbose insights
- Command-line flags for all configuration options
- Shell completion generation for bash, zsh, fish, and PowerShell

## Project Layout

```
core/
  cmd/bip38cli/             # Go entry point for the CLI binary
  internal/cli/             # Cobra commands and user interaction flows
  internal/bip38/           # BIP38 domain logic and tests
  pkg/                      # Internal packages and utilities
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


## Usage

### Generate a New Wallet (WIF)

```bash
# Generate a compressed mainnet WIF (default settings)
bip38cli wallet generate

# Target another network (e.g. testnet) and show the derived address
bip38cli wallet generate --network testnet --show-address

# Encrypt the generated key with BIP38 (interactive passphrase prompt)
bip38cli wallet generate --encrypt

# JSON output (includes WIF, network, compression, and optional address/encrypted key)
bip38cli wallet generate --output-format json --show-address

# Legacy P2PKH (BIP44) address instead of bech32 (BIP84 default)
bip38cli wallet generate --address-type bip44 --show-address

# Generate an uncompressed key (forces legacy P2PKH output)
bip38cli wallet generate --uncompressed
```

> Addresses for compressed keys follow BIP84 (bech32). If you explicitly generate an uncompressed key, the CLI falls back to legacy P2PKH output.

Use `--address-type bip44` whenever you need to force a legacy P2PKH address for compatibility with older wallets.

### Inspect an Existing WIF

```bash
# Show network, compression and address
bip38cli wallet inspect 5HueCGU8rMjxEXxiPuD5BDku4MkFqeZyd4dZ1jvhTVqvbTLvyTJ

# JSON output (machine friendly)
bip38cli wallet inspect --output-format json KwYgW8gcxj1JWJXhPSu4Fqwzfhp5Yfi42mdYmMa4XqK7NJxXUSK7

# Force legacy P2PKH output
bip38cli wallet inspect --address-type bip44 KwYgW8gcxj1JWJXhPSu4Fqwzfhp5Yfi42mdYmMa4XqK7NJxXUSK7
```

### Encrypt a WIF Key

```bash
# Basic encryption
bip38cli encrypt KwYgW8gcxj1JWJXhPSu4Fqwzfhp5Yfi42mdYmMa4XqK7NJxXUSK7
# Hidden prompts ask for passphrase twice
# Result: 6PRV...

# Force compressed format (same as default)
bip38cli encrypt --compressed KwYgW8gcxj1JWJXhPSu4Fqwzfhp5Yfi42mdYmMa4XqK7NJxXUSK7

# Force uncompressed format  
bip38cli encrypt --uncompressed KwYgW8gcxj1JWJXhPSu4Fqwzfhp5Yfi42mdYmMa4XqK7NJxXUSK7

# Use global compressed flag (default: compressed)
bip38cli encrypt --compressed KwYgW8gcxj1JWJXhPSu4Fqwzfhp5Yfi42mdYmMa4XqK7NJxXUSK7

# Use uncompressed by default
bip38cli encrypt --compressed=false KwYgW8gcxj1JWJXhPSu4Fqwzfhp5Yfi42mdYmMa4XqK7NJxXUSK7

# JSON output
bip38cli encrypt --output-format json KwYgW8gcxj1JWJXhPSu4Fqwzfhp5Yfi42mdYmMa4XqK7NJxXUSK7
```

### Decrypt an Encrypted Key

```bash
# Basic decryption
bip38cli decrypt 6PRVWUbkzzsbcVac2qwfssoUJAN1Xhrg6bNk8J7Nzm5H7kxEbn2Nh2ZoGg
# Hidden prompt asks for passphrase

# Show derived address
bip38cli decrypt --show-address 6PRVWUbkzzsbcVac2qwfssoUJAN1Xhrg6bNk8J7Nzm5H7kxEbn2Nh2ZoGg

# JSON output with address
bip38cli decrypt --show-address --output-format json 6PRVWUbkzzsbcVac2qwfssoUJAN1Xhrg6bNk8J7Nzm5H7kxEbn2Nh2ZoGg
```

### Work with Intermediate Codes

```bash
# Generate basic intermediate code
bip38cli intermediate generate

# Generate code with lot/sequence metadata
bip38cli intermediate generate --lot 123 --sequence 456 --use-lot-sequence

# Generate with JSON output
bip38cli intermediate generate --output-format json --lot 123 --sequence 456 --use-lot-sequence

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

BIP38CLI uses command-line flags for all configuration options. No configuration files are required.

Global flags:
- `--verbose, -v`: Enable verbose output for additional diagnostic information
- `--output-format`: Output format (text|json, default: text)
- `--compressed, -c`: Use compressed public key format (default: true)
- `--uncompressed`: Use uncompressed public key format (overrides --compressed)

Command-specific flags:
- `encrypt --compressed`: Force compressed public key format
- `encrypt --uncompressed`: Force uncompressed public key format  
- `decrypt --show-address`: Show the Bitcoin address for the decrypted key
- `decrypt --address-type <bip84|bip44>`: Control address encoding when `--show-address` is used (default: bip84)
- `intermediate generate --lot <number>`: Specify lot number (0-1048575)
- `intermediate generate --sequence <number>`: Specify sequence number (0-4095)
- `intermediate generate --use-lot-sequence`: Use lot and sequence numbers
- `wallet generate --address-type <bip84|bip44>`: Choose bech32 (bip84) or legacy P2PKH (bip44) output
- `wallet generate --uncompressed`: Produce an uncompressed key (implicitly legacy address)
- `wallet inspect --address-type <bip84|bip44>`: Inspect WIFs using the desired address encoding
- `wallet generate --network <name>`: Choose network (`mainnet`, `testnet`, `regtest`, `simnet`, `signet`)
- `wallet generate --encrypt`: Encrypt the generated key with BIP38 (interactive passphrase)
- `wallet generate --show-address`: Display the derived Bitcoin address for the new key

### Examples with JSON Output

```bash
# Encrypt with JSON output
bip38cli encrypt --output-format json KwYgW8gcxj1JWJXhPSu4Fqwzfhp5Yfi42mdYmMa4XqK7NJxXUSK7
# Output: {"encrypted_key": "6PRV...", "compressed": true}

# Decrypt with JSON output and address
bip38cli decrypt --show-address --output-format json 6PRVWUbkzzsbcVac2qwfssoUJAN1Xhrg6bNk8J7Nzm5H7kxEbn2Nh2ZoGg
# Output: {"private_key": "KwYg...", "compressed": true, "address": "bc1qklnjad76qxxxy833ggfjsjyjc29vdrgnpnju5d"}
```

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
