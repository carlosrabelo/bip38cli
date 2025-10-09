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
wget https://github.com/mannkind/bip38cli-cli/releases/latest/download/bip38cli-linux-amd64.tar.gz

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
git clone https://github.com/mannkind/bip38cli-cli.git
cd bip38cli-cli

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

### Scenario 1: Creating Your Personal "Digital Vault"

**Situation:** You've accumulated a significant amount of Bitcoin in a software wallet (like Electrum, Sparrow, or on your phone) and you're starting to worry. What if your computer or phone breaks? A simple wallet file backup isn't safe from thieves. Let's create a robust, password-protected backup.

**The Flow:**

#### Step 1: Get the Private Key (The Raw Material)

The goal is to extract your private key in WIF (Wallet Import Format).

1.  **Find the export option:** In your wallet, look for options like "Private Keys," "Export," or "Sweep."
2.  **Security Warning:** This is a sensitive moment. Exposing your private key is like opening your vault. Do this in an offline environment if possible, and make sure no one is watching your screen.

For our example, let's say you exported the key: `5KYZdUEo39z3FPrtuX2QbbwGnNP5zTd7yyr2SC1j299sBCnWjss`

---

#### Step 2: The Encryption Ritual

Now, let's lock this key in a digital vault using a strong passphrase.

```bash
# Start the command interactively for maximum security
$ bip38cli encrypt --compressed

# The program will ask for the key. Paste the key you exported:
Enter WIF private key: 5KYZdUEo39z3FPrtuX2QbbwGnNP5zTd7yyr2SC1j299sBCnWjss

# Now, create and enter a master passphrase.
# Tip: Use the "diceware" method or a phrase with 5+ random words.
Enter passphrase: 
# Ex: "correct horse battery staple window thunder"
Confirm passphrase: 

# Result (Your encrypted key! Save it.)
6PRW5VbB4Qj2a2m2b2GzQvH3JPTgS3t2fA1Zytx
```
*Use the `--compressed` flag because most modern wallets use compressed addresses.*

---

#### Step 3: Verification (Trust, but Verify)

**This is the most important step.** Never trust a blind backup. You need to be 100% certain that the encrypted key and your passphrase work perfectly.

```bash
# Use the decrypt command with the --show-address flag
$ bip38cli decrypt --show-address 6PRW5VbB4Qj2a2m2b2GzQvH3JPTgS3t2fA1Zytx

Enter passphrase: 
# Enter the same master passphrase

# Expected Result:
WIF: 5KYZdUEo39z3FPrtuX2QbbwGnNP5zTd7yyr2SC1j299sBCnWjss
Address: 1CC3X2gu58d6wXUWMffpuzN9JAfTUw4K9A 
```
**Confirm that the `Address` displayed is exactly the same** as the one in your original wallet from where you exported the key. If it matches, your digital vault is validated.

---

#### Step 4: The Separation Principle

Now, store the two "pieces" of your vault in completely separate locations. The logic is that a thief would need to find both to steal your funds.

*   **The Encrypted Key (`6PR...`):**
    *   Save it in a text file on your Google Drive, Dropbox, or a USB drive.
    *   Print it on paper and keep it in a physical location.
*   **Your Master Passphrase:**
    *   Store it in a trusted password manager (Bitwarden, 1Password).
    *   Write it on a piece of paper and keep it in a physical safe, in a different location from the encrypted key paper.

With this, your backup is secure and resilient.

### Scenario 2: Simple Digital Inheritance Planning

**Situation:** You want to ensure your loved ones can access your Bitcoin if something happens to you. Leaving a written-down private key is dangerous, and a wallet file can be confusing for non-technical people. A BIP38 key with a well-thought-out recovery plan is a balanced solution.

**The Flow:**

#### Step 1: Choose the Key and a "Family Passphrase"

Encrypt the main private key of your holdings. The crucial part is the passphrase. It shouldn't be trivial, but it should be recoverable by your family through a method only they would know.

*   **Bad Idea:** `dogs-birthday` (too easy to guess).
*   **Good Idea:** A phrase from a book you all share, followed by the page number. Ex: `phrase-from-book-page-42`. The method to arrive at the passphrase is the secret.

Let's use the passphrase `family-treasure-reunion-2015` for this example.

---

#### Step 2: Encrypt and Create the "Inheritance Package"

Encrypt the key and prepare a kit for your family.

```bash
# Encrypt the key with the "Family Passphrase"
$ bip38cli encrypt 5KYZdUEo39z3FPrtuX2QbbwGnNP5zTd7yyr2SC1j299sBCnWjss --compressed
# Enter the passphrase...
# Result: 6Pz... (write this key down)
```

Now, create a folder or a physical envelope named **"Bitcoin Inheritance Package"** containing:

1.  **The Encrypted Key:** A text file (`inheritance_key.txt`) containing only the `6Pz...` key.
2.  **The Necessary Software:** Include the `bip38cli` executable or a direct link to the GitHub releases page (`https://github.com/carlosrabelo/bip38cli/releases`).
3.  **Clear Instructions:** A file named `INSTRUCTIONS.txt`.

**Template for `INSTRUCTIONS.txt`:**

```
=================================================
INSTRUCTIONS FOR ACCESSING BITCOIN LEGACY
=================================================

To access the funds, follow these steps:

1.  Open the "bip38cli" program located in this folder (or download it from the provided link).

2.  Run the following command in the terminal (command prompt):
    
    ./bip38cli decrypt [the key from the inheritance_key.txt file]

3.  The program will ask for a passphrase. The hint for the passphrase is:
    [WRITE YOUR HINT HERE. Ex: "The name of our first boat and the year we bought it, all joined with hyphens."]

4.  The program will output a long private key (starting with 'K' or 'L').

5.  Import this key into a modern Bitcoin wallet (like BlueWallet or Electrum) to access the funds.

With love.
```

---

#### Step 3: Store the Package and the Clues

The security of this method depends on separation.

*   **The "Inheritance Package"** (with the encrypted key and instructions) should be kept in a secure location known to your family or lawyer (e.g., a safe, a safety deposit box).
*   **The Clues for the Passphrase** (or the passphrase itself) must be stored in a **different** location. This could be in a sealed envelope with your lawyer, in a shared password manager vault with your spouse, or another creative solution.

This method creates a robust barrier against theft while providing a clear path for your family.

### Scenario 3: Secure Delegation of Key Creation (with Intermediate Code)

**Situation:** Imagine **Alice** wants **Bob** to generate new paper wallets for her. Alice has a super-secure master passphrase that she doesn't want to share with anyone, not even Bob. How can Bob create keys that are already encrypted with Alice's passphrase, without ever knowing what it is?

This is where the **intermediate code** comes in. It acts as an "authorization" that allows encryption, but not decryption.

**The Flow:**

#### Step 1: Alice Generates the Intermediate Code

Alice, on her secure computer, uses her master passphrase to generate an intermediate code. This code is not her passphrase, but is derived from it.

```bash
# On ALICE's computer
$ bip38cli intermediate generate

Enter passphrase: 
# Alice types her secret passphrase (e.g., "do-not-share-this-secret-phrase")
Confirm passphrase: 

# Result (this is the code Alice will share)
passphrasezctFpQWj9H252m2b2GzQvH3JPTgS3t2fA
```
**Important:** Alice keeps her secret passphrase safe and **sends only the `passphrasez...` code to Bob**.

---

#### Step 2: Bob Generates a New Private Key

Bob, on his system (which can be online, it doesn't need to be super-secure for this part), generates a new Bitcoin key pair. He gets a private key in WIF format.

For this example, let's assume Bob generated the following (uncompressed) key:

`5Jag5pY5aWJtL3A2YdDB5s2b2GzQvH3JPTgS3t2fA`

---

#### Step 3: Bob Encrypts the Key for Alice

Now, Bob uses the **intermediate code** he received from Alice to encrypt the new private key. He **does not use Alice's passphrase** (he doesn't even know it).

```bash
# On BOB's computer
$ bip38cli encrypt 5Jag5pY5aWJtL3A2YdDB5s2b2GzQvH3JPTgS3t2fA --intermediate passphrasezctFpQWj9H252m2b2GzQvH3JPTgS3t2fA

# The command does NOT ask for a passphrase! It uses the intermediate code.

# Result (encrypted key)
6PRW5VbB4Qj2a2m2b2GzQvH3JPTgS3t2fA1Zytx
```
Bob can now give this encrypted key (`6PRW...`) to Alice. He can print it on a paper wallet or send it digitally. Bob cannot decrypt this key.

---

#### Step 4: Alice Verifies and Uses the Key

Alice receives the encrypted key from Bob. To make sure it works, she decrypts it using her **original master passphrase**.

```bash
# Back on ALICE's computer
$ bip38cli decrypt 6PRW5VbB4Qj2a2m2b2GzQvH3JPTgS3t2fA1Zytx

Enter passphrase: 
# Alice types her original secret passphrase ("do-not-share-this-secret-phrase")

# Result (the private key Bob generated)
5Jag5pY5aWJtL3A2YdDB5s2b2GzQvH3JPTgS3t2fA
```

**Conclusion and Security Advantage:**

The process was a success!

*   **Bob** was able to create a BIP38 encrypted key for Alice.
*   **Alice** is the only one who can decrypt it.
*   **Alice's master passphrase never left** her secure computer.

This method is ideal for services that generate paper wallets, so users can be sure that not even the service that created the wallet has access to their funds.

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
wget https://github.com/mannkind/bip38cli-cli/releases/latest/download/bip38cli-linux-amd64.tar.gz
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
- **Issues:** https://github.com/mannkind/bip38cli-cli/issues
- **Source:** https://github.com/mannkind/bip38cli-cli

**Remember: Security first, always!** 