# BIP38CLI - Complete Tutorial in English

This tutorial teaches you how to use BIP38CLI to securely encrypt Bitcoin private keys with maximum security.

##  Table of Contents

1. [Installation](#installation)
2. [Basic Concepts](#basic-concepts)
3. [Basic Usage](#basic-usage)
4. [Advanced Scenarios](#advanced-scenarios)
5. [Automation](#automation)
6. [Security and Best Practices](#security-and-best-practices)
7. [Troubleshooting](#troubleshooting)

##  Installation

### Option 1: Direct Download (Recommended)
```bash
# Download the latest release
wget https://github.com/carlosrabelo/bip38cli/releases/latest/download/bip38cli-linux-amd64.tar.gz

# Extract the file
tar -xzf bip38cli-linux-amd64.tar.gz

# Make executable
chmod +x bip38cli

# Move to PATH (optional)
sudo mv bip38cli /usr/local/bin/
```

### Option 2: Build from Source
```bash
# Clone the repository
git clone https://github.com/carlosrabelo/bip38cli.git
cd bip38cli

# Build
make build

# Binary will be at ./bin/bip38cli
```

### Installation Verification
```bash
bip38cli --version
# Should show: bip38cli version x.x.x (built: ...)

bip38cli --help
# Should show complete help
```

##  Basic Concepts

### What is BIP38?
BIP38 is a standard for encrypting Bitcoin private keys with a passphrase. It enables:
- **Secure storage** of private keys
- **Password-protected backups**
- **Digital inheritance** (family can access with password)
- **Two-factor security**

### Key Types
```
Private Key (WIF):       5HueCGU8rMjxEXxiPuD5BDku4MkFqeZyd4dZ1jvhTVqvbTLvyTJ
Encrypted Key:           6PRVWUbkzzsbcVac2qwfssoUJAN1Xhrg6bNk8J7Nzm5H7kxEbn2Nh2ZoGg
Intermediate Code:       passphraseabc123def456ghi789jkl012mno345pqr...
```

##  Basic Usage

### Generate a New Wallet (WIF)

```bash
# Create a compressed mainnet key (default)
bip38cli wallet generate

# Target another network (e.g., testnet) and reveal the address
bip38cli wallet generate --network testnet --show-address

# Immediately wrap the key with BIP38 encryption (interactive passphrase)
bip38cli wallet generate --encrypt

# Force legacy P2PKH output (BIP44) or uncompressed keys
bip38cli wallet generate --address-type bip44 --show-address
bip38cli wallet generate --uncompressed
```

Compressed keys produce BIP84 (bech32) addresses. If you deliberately use an
uncompressed key, the CLI falls back to a legacy P2PKH address.

Use `--address-type bip44` with either `wallet generate` or `wallet inspect`
when you explicitly need legacy P2PKH output from a compressed key.

### Inspect an Existing WIF

```bash
# Show network, compression flag and derived address
bip38cli wallet inspect 5HueCGU8rMjxEXxiPuD5BDku4MkFqeZyd4dZ1jvhTVqvbTLvyTJ

# JSON output for scripts
bip38cli wallet inspect --output-format json KwYgW8gcxj1JWJXhPSu4Fqwzfhp5Yfi42mdYmMa4XqK7NJxXUSK7
```

### 1. Encrypt a Private Key

**Interactive Mode (More Secure):**
```bash
bip38cli encrypt
# Enter WIF key when prompted
# Enter password (won't be displayed)
# Confirm password
# Result: encrypted key 6P...
```

**Direct Mode:**
```bash
bip38cli encrypt 5HueCGU8rMjxEXxiPuD5BDku4MkFqeZyd4dZ1jvhTVqvbTLvyTJ
# Enter password only
# Result: 6PRVWUbkzzsbcVac2qwfssoUJAN1Xhrg6bNk8J7Nzm5H7kxEbn2Nh2ZoGg
```

**Force Compression:**
```bash
bip38cli encrypt --compressed 5HueCGU8rMjxEXxiPuD5BDku4MkFqeZyd4dZ1jvhTVqvbTLvyTJ
```

### 2. Decrypt a Key

**Interactive Mode:**
```bash
bip38cli decrypt
# Enter encrypted key 6P...
# Enter password
# Result: original WIF key
```

**With Bitcoin Address:**
```bash
bip38cli decrypt --show-address 6PRVWUbkzzsbcVac2qwfssoUJAN1Xhrg6bNk8J7Nzm5H7kxEbn2Nh2ZoGg
# Shows WIF key + corresponding Bitcoin address
```

### 3. Intermediate Codes (Advanced Security)

**Generate Code:**
```bash
bip38cli intermediate generate
# Enter password
# Result: intermediate code passphrase...
```

**Validate Code:**
```bash
bip38cli intermediate validate passphraseabc123...
# Verifies if the code is valid
```

##  Advanced Scenarios

### Scenario 1: Secure Wallet Backup

**Situation:** You have a Bitcoin wallet and want ultra-secure backup.

```bash
# 1. Export private key from your wallet
# (in Electrum: Wallet > Private Keys > Export)

# 2. Encrypt with strong password
bip38cli encrypt --compressed
# Enter your private key
# Enter VERY strong password (e.g., 6-word phrase)
# WRITE DOWN the resulting encrypted key

# 3. TEST decryption
bip38cli decrypt --show-address
# Enter encrypted key
# Enter password
# CONFIRM address matches your wallet

# 4. Store securely:
# - Encrypted key in cloud/paper
# - Password in DIFFERENT location
```

### Scenario 2: Digital Inheritance

**Situation:** Leave Bitcoin for family to access in emergency.

```bash
# 1. Encrypt your keys
bip38cli encrypt 5HueCGU8rMjxEXxiPuD5BDku4MkFqeZyd4dZ1jvhTVqvbTLvyTJ
# Use password family can figure out (e.g., birthdate + pet name)

# 2. Create simple instructions:
echo "In emergency:
1. Download bip38cli from github.com/carlosrabelo/bip38cli
2. Run: bip38cli decrypt 6PRVWUbkzzsbcVac2qwfssoUJAN1Xhrg6bNk8J7Nzm5H7kxEbn2Nh2ZoGg
3. Password is: [your hint here]
4. Use resulting WIF key in any Bitcoin wallet" > inheritance_instructions.txt
```

### Scenario 3: Two People, Maximum Security

**Situation:** One person knows password, another generates keys.

```bash
# Person A (knows password):
bip38cli intermediate generate
# Enter secret password
# Send intermediate code to Person B

# Person B (generates keys, doesn't know password):
# [would use 'generate' command - not implemented yet]
# Sends encrypted keys to Person A

# Person A can decrypt when needed
```

##  Automation

### Batch Backup Script

**backup_wallets.sh:**
```bash
#!/bin/bash

# List of private keys
KEYS=(
    "5HueCGU8rMjxEXxiPuD5BDku4MkFqeZyd4dZ1jvhTVqvbTLvyTJ"
    "5J3mBbAH58CpQ3Y5RNJpUKPE62SQ5tfcvU2JpbnkeyhfsYB1Jcn"
    # ... more keys
)

# Password (NEVER hardcode in production!)
read -s -p "Enter password for encryption: " PASSWORD
echo

# Output file
OUTPUT="encrypted_backup_$(date +%Y%m%d).txt"

echo "Starting backup of ${#KEYS[@]} keys..."

for i in "${!KEYS[@]}"; do
    echo "Processing key $((i+1))..."

    # Encrypt key (simulates password input)
    RESULT=$(echo -e "${KEYS[i]}\n$PASSWORD\n$PASSWORD" | bip38cli encrypt)

    if [[ $RESULT == 6P* ]]; then
        echo "Key $((i+1)): $RESULT" >> "$OUTPUT"
        echo "✓ Key $((i+1)) encrypted successfully"
    else
        echo "✗ Error with key $((i+1))"
    fi
done

echo "Backup saved to: $OUTPUT"
echo "IMPORTANT: Store password in secure and separate location!"
```

### Verification Script

**verify_backup.sh:**
```bash
#!/bin/bash

BACKUP_FILE="$1"
if [[ ! -f "$BACKUP_FILE" ]]; then
    echo "Usage: $0 backup_file.txt"
    exit 1
fi

read -s -p "Enter backup password: " PASSWORD
echo

echo "Verifying backup integrity..."

while IFS= read -r line; do
    if [[ $line == *":"* ]]; then
        NUMBER=$(echo "$line" | cut -d: -f1)
        ENCRYPTED_KEY=$(echo "$line" | cut -d: -f2 | tr -d ' ')

        # Try to decrypt
        RESULT=$(echo -e "$ENCRYPTED_KEY\n$PASSWORD" | bip38cli decrypt 2>/dev/null)

        if [[ $RESULT == 5* ]] || [[ $RESULT == K* ]] || [[ $RESULT == L* ]]; then
            echo "✓ $NUMBER: OK"
        else
            echo "✗ $NUMBER: ERROR"
        fi
    fi
done < "$BACKUP_FILE"

echo "Verification complete!"
```

##  Security and Best Practices

###  DO

- **Use very strong passwords** (minimum 15 characters, better: 6+ word phrases)
- **Always test decryption** before storing backup
- **Use offline mode** (disconnect internet during use)
- **Store password and key in DIFFERENT locations**
- **Make multiple copies** of encrypted backup
- **Use `--compressed`** for modern wallet keys
- **Verify addresses** with `--show-address`

###  DON'T

- **NEVER** use simple passwords (123456, password, etc.)
- **NEVER** store password with encrypted key
- **NEVER** trust only one copy
- **NEVER** use on infected/public computers
- **NEVER** share unencrypted keys
- **NEVER** hardcode passwords in scripts
- **NEVER** use internet during critical operations

###  Security Levels

**Level 1 - Basic:**
- Strong password (15+ characters)
- Encrypted backup in 2 locations
- Decryption test

**Level 2 - Advanced:**
- 6+ word passphrase
- Backup in 3+ different locations
- Dedicated computer/live USB
- Periodic verification

**Level 3 - Paranoid:**
- Password generated from physical entropy
- Air-gapped computer
- Multiple encryption layers
- Dead man's switch for family

##  Troubleshooting

### Error: "invalid WIF private key"
```bash
# Problem: Invalid key format
# Solution: Check if key starts with 5, K or L
bip38cli encrypt KwYgW8gcxj1JWJXhPSu4Fqwzfhp5Yfi42mdYmMa4XqK7NJxXUSK7
```

### Error: "incorrect passphrase"
```bash
# Problem: Wrong password during decryption
# Solution: Check caps lock, keyboard, special characters
# Try password variations if necessary
```

### Error: "command not found"
```bash
# Problem: bip38cli not in PATH
# Solution: Use full path
./bin/bip38cli --help

# Or add to PATH:
export PATH=$PATH:$(pwd)/bin
```

### Slow Performance
```bash
# Problem: Encryption takes too long
# Cause: BIP38 uses scrypt (intentionally slow)
# Normal: 2-5 seconds per operation
# If > 30 seconds: very slow computer or problem
```

### Key Doesn't Work in Wallet
```bash
# Problem: Decrypted key doesn't work
# Possible causes:
# 1. Wrong compressed/uncompressed format
bip38cli decrypt --show-address your_6P_key
# Compare address with original wallet

# 2. Wrong network (testnet vs mainnet)
# 3. Error in original password
```

##  Complete Practical Examples

### Example 1: First Time
```bash
# Install
wget https://github.com/carlosrabelo/bip38cli/releases/latest/download/bip38cli-linux-amd64.tar.gz
tar -xzf bip38cli-linux-amd64.tar.gz
chmod +x bip38cli

# Test with example key
echo "Testing with example key..."
./bip38cli encrypt 5HueCGU8rMjxEXxiPuD5BDku4MkFqeZyd4dZ1jvhTVqvbTLvyTJ
# Enter password: mytestpassword123
# Result: 6PRVWUbkzz...

# Verify
./bip38cli decrypt --show-address
# Enter resulting 6P key
# Enter password: mytestpassword123
# Should return original key + address
```

### Example 2: Real Backup
```bash
# 1. Prepare secure environment
sudo systemctl stop NetworkManager  # Disconnect internet
cd /tmp  # Use temporary folder

# 2. Encrypt real key
./bip38cli encrypt --compressed
# Enter your real wallet key
# Enter super strong password
# WRITE DOWN the 6P result

# 3. Test immediately
./bip38cli decrypt --show-address
# Enter the 6P key
# Enter the password
# CONFIRM address matches

# 4. Clean up and reconnect
cd /
rm -rf /tmp/*
sudo systemctl start NetworkManager
```

---

##  IMPORTANT WARNING

**This software handles Bitcoin private keys. Mistakes can result in PERMANENT LOSS of funds.**

- **ALWAYS test** with small amounts first
- **ALWAYS backup** original keys
- **ALWAYS verify** addresses after decryption
- **NEVER use** in untrusted environments

**The author is not responsible for losses. Use at your own risk.**

---

##  Support

- **Documentation:** README.md
- **Issues:** https://github.com/carlosrabelo/bip38cli/issues
- **Source:** https://github.com/carlosrabelo/bip38cli

**Remember: Security first, always!** 
