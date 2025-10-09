# Tutorial BIP38CLI CLI - Guia Completo em Português

Este tutorial ensina como usar o BIP38CLI CLI para criptografar chaves privadas Bitcoin com segurança máxima.

##  Índice

1. [Instalação](#instalação)
2. [Conceitos Básicos](#conceitos-básicos)
3. [Uso Básico](#uso-básico)
4. [Cenários Avançados](#cenários-avançados)
5. [Automação](#automação)
6. [Segurança e Boas Práticas](#segurança-e-boas-práticas)
7. [Solução de Problemas](#solução-de-problemas)

##  Instalação

### Opção 1: Download Direto (Recomendado)
```bash
# Baixe a versão mais recente
wget https://github.com/mannkind/bip38cli-cli/releases/latest/download/bip38cli-linux-amd64.tar.gz

# Extraia o arquivo
tar -xzf bip38cli-linux-amd64.tar.gz

# Torne executável
chmod +x bip38cli

# Mova para PATH (opcional)
sudo mv bip38cli /usr/local/bin/
```

### Opção 2: Compilar do Código
```bash
# Clone o repositório
git clone https://github.com/mannkind/bip38cli-cli.git
cd bip38cli-cli

# Compile
make build

# O binário estará em ./bin/bip38cli
```

### Verificação da Instalação
```bash
bip38cli --version
# Deve mostrar: bip38cli version x.x.x (built: ...)

bip38cli --help
# Deve mostrar a ajuda completa
```

##  Conceitos Básicos

### O que é BIP38?
BIP38 é um padrão para criptografar chaves privadas Bitcoin com uma senha. Permite:
- **Armazenamento seguro** de chaves privadas
- **Backup protegido por senha**
- **Herança digital** (família pode acessar com senha)
- **Segurança de dois fatores**

### Tipos de Chaves
```
Chave Privada (WIF):     5HueCGU8rMjxEXxiPuD5BDku4MkFqeZyd4dZ1jvhTVqvbTLvyTJ
Chave Criptografada:     6PRVWUbkzzsbcVac2qwfssoUJAN1Xhrg6bNk8J7Nzm5H7kxEbn2Nh2ZoGg
Código Intermediário:    passphraseabc123def456ghi789jkl012mno345pqr...
```

##  Uso Básico

### 1. Criptografar uma Chave Privada

**Modo Interativo (Mais Seguro):**
```bash
bip38cli encrypt
# Digite a chave WIF quando solicitado
# Digite a senha (não será exibida)
# Confirme a senha
# Resultado: chave criptografada 6P...
```

**Modo Direto:**
```bash
bip38cli encrypt 5HueCGU8rMjxEXxiPuD5BDku4MkFqeZyd4dZ1jvhTVqvbTLvyTJ
# Digite apenas a senha
# Resultado: 6PRVWUbkzzsbcVac2qwfssoUJAN1Xhrg6bNk8J7Nzm5H7kxEbn2Nh2ZoGg
```

**Forçar Compressão:**
```bash
bip38cli encrypt --compressed 5HueCGU8rMjxEXxiPuD5BDku4MkFqeZyd4dZ1jvhTVqvbTLvyTJ
```

### 2. Descriptografar uma Chave

**Modo Interativo:**
```bash
bip38cli decrypt
# Digite a chave criptografada 6P...
# Digite a senha
# Resultado: chave WIF original
```

**Com Endereço Bitcoin:**
```bash
bip38cli decrypt --show-address 6PRVWUbkzzsbcVac2qwfssoUJAN1Xhrg6bNk8J7Nzm5H7kxEbn2Nh2ZoGg
# Mostra a chave WIF + endereço Bitcoin correspondente
```

### 3. Códigos Intermediários (Segurança Avançada)

**Gerar Código:**
```bash
bip38cli intermediate generate
# Digite a senha
# Resultado: código intermediário passphrase...
```

**Validar Código:**
```bash
bip38cli intermediate validate passphraseabc123...
# Verifica se o código é válido
```

##  Cenários Avançados

### Cenário 1: Criando seu "Cofre Digital" Pessoal

**Situação:** Você acumulou uma quantidade significativa de Bitcoin em uma carteira de software (como Electrum, Sparrow ou no seu celular) e começa a se preocupar. E se o computador ou celular quebrar? Um simples backup do arquivo da carteira não é seguro contra ladrões. Vamos criar um backup robusto e protegido por senha.

**O Fluxo:**

#### Passo 1: Obter a Chave Privada (O Material Bruto)

O objetivo é extrair sua chave privada em formato WIF (Wallet Import Format).

1.  **Localize a opção de exportação:** Em sua carteira, procure por opções como "Private Keys", "Export" ou "Sweep".
2.  **Aviso de Segurança:** Este é um momento delicado. Expor sua chave privada é como abrir seu cofre. Faça isso em um ambiente offline, se possível, e certifique-se de que ninguém está vendo sua tela.

Para nosso exemplo, digamos que você exportou a chave: `5KYZdUEo39z3FPrtuX2QbbwGnNP5zTd7yyr2SC1j299sBCnWjss`

---

#### Passo 2: O Ritual de Criptografia

Agora, vamos trancar essa chave em um cofre digital usando uma senha forte.

```bash
# Inicie o comando de forma interativa para máxima segurança
$ bip38cli encrypt --compressed

# O programa vai pedir a chave. Cole a chave que você exportou:
Enter WIF private key: 5KYZdUEo39z3FPrtuX2QbbwGnNP5zTd7yyr2SC1j299sBCnWjss

# Agora, crie e digite uma senha mestra.
# Dica: Use o método "diceware" ou uma frase com 5+ palavras aleatórias.
Enter passphrase: 
# Ex: "correto cavalo bateria grampo janela trovão"
Confirm passphrase: 

# Resultado (Sua chave criptografada! Guarde-a.)
6PRW5VbB4Qj2a2m2b2GzQvH3JPTgS3t2fA1Zytx
```
*Use a flag `--compressed` porque a maioria das carteiras modernas usa endereços comprimidos.*

---

#### Passo 3: Verificação (Confie, mas Verifique)

**Este é o passo mais importante.** Nunca confie em um backup cego. Você precisa ter 100% de certeza de que a chave criptografada e sua senha funcionam perfeitamente.

```bash
# Use o comando de descriptografia com a flag --show-address
$ bip38cli decrypt --show-address 6PRW5VbB4Qj2a2m2b2GzQvH3JPTgS3t2fA1Zytx

Enter passphrase: 
# Digite a mesma senha mestra

# Resultado Esperado:
WIF: 5KYZdUEo39z3FPrtuX2QbbwGnNP5zTd7yyr2SC1j299sBCnWjss
Address: 1CC3X2gu58d6wXUWMffpuzN9JAfTUw4K9A 
```
**Confirme que o `Address` (endereço) exibido é exatamente o mesmo** da sua carteira original de onde você exportou a chave. Se for, seu cofre digital está validado.

---

#### Passo 4: O Princípio da Separação

Agora, guarde os dois "pedaços" do seu cofre em locais completamente separados. A lógica é que um ladrão precisaria encontrar ambos para roubar seus fundos.

*   **A Chave Criptografada (`6PR...`):**
    *   Salve em um arquivo de texto no seu Google Drive, Dropbox ou em um pendrive.
    *   Imprima em um papel e guarde em um local físico.
*   **Sua Senha Mestra:**
    *   Guarde em um gerenciador de senhas confiável (Bitwarden, 1Password).
    *   Escreva em um papel e guarde em um cofre físico, em um local diferente do papel da chave criptografada.

Com isso, seu backup está seguro e resiliente.

### Cenário 2: Planejamento de Herança Digital Simples

**Situação:** Você quer garantir que seus entes queridos possam acessar seus Bitcoins caso algo lhe aconteça. Deixar uma chave privada anotada é perigoso e um arquivo de carteira pode ser confuso para leigos. Uma chave BIP38 com um plano de recuperação bem pensado é uma solução equilibrada.

**O Fluxo:**

#### Passo 1: Escolher a Chave e a "Senha de Família"

Criptografe a chave privada principal do seu patrimônio. A parte crucial é a senha. Ela não deve ser trivial, mas deve ser recuperável por sua família através de um método que só eles conheçam.

*   **Ideia Ruim:** `aniversario-do-cachorro` (muito fácil de adivinhar).
*   **Ideia Boa:** Uma frase de um livro que todos na casa conhecem, seguida do número da página. Ex: `frase-do-livro-pagina-42`. O método para chegar na senha é o segredo.

Vamos usar a senha `tesouro-da-familia-reunida-2015` para este exemplo.

---

#### Passo 2: Criptografar e Criar o "Pacote de Herança"

Criptografe a chave e prepare um kit para sua família.

```bash
# Criptografe a chave com a "Senha de Família"
$ bip38cli encrypt 5KYZdUEo39z3FPrtuX2QbbwGnNP5zTd7yyr2SC1j299sBCnWjss --compressed
# Digite a senha...
# Resultado: 6Pz... (anote esta chave)
```

Agora, crie uma pasta ou um envelope físico chamado **"Pacote de Herança Bitcoin"** contendo:

1.  **A Chave Criptografada:** Um arquivo de texto (`chave_heranca.txt`) contendo apenas a chave `6Pz...`.
2.  **O Software Necessário:** Inclua o executável `bip38cli` ou um link direto para a página de releases do GitHub (`https://github.com/carlosrabelo/bip38cli/releases`).
3.  **Instruções Claras:** Um arquivo `INSTRUCOES.txt`.

**Template para `INSTRUCOES.txt`:**

```
=================================================
INSTRUÇÕES PARA ACESSO AO LEGADO BITCOIN
=================================================

Para acessar os fundos, siga estes passos:

1.  Abra o programa "bip38cli" que está nesta pasta (ou baixe do link fornecido).

2.  Execute o seguinte comando no terminal (prompt de comando):
    
    ./bip38cli decrypt [a chave que está no arquivo chave_heranca.txt]

3.  O programa pedirá uma senha. A pista para a senha é:
    [ESCREVA SUA PISTA AQUI. Ex: "O nome do nosso primeiro barco e o ano em que o compramos, tudo junto com hífens."]

4.  O programa irá gerar uma chave privada longa (começando com 'K' ou 'L').

5.  Importe essa chave em uma carteira de Bitcoin moderna (como BlueWallet ou Electrum) para ter acesso aos fundos.

Com amor.
```

---

#### Passo 3: Armazenar o Pacote e as Pistas

A segurança deste método depende da separação.

*   **O "Pacote de Herança"** (com a chave criptografada e as instruções) deve ser guardado em um local seguro e conhecido pela sua família ou advogado (ex: cofre, caixa de segurança).
*   **As Pistas para a Senha** (ou a senha em si) devem ser guardadas em um local **diferente**. Pode ser em um envelope lacrado com seu advogado, em um gerenciador de senhas compartilhado com seu cônjuge, ou outra solução criativa.

Este método cria uma barreira robusta contra roubo, mas um caminho claro para sua família.

### Cenário 3: Delegação Segura de Criação de Chaves (com Código Intermediário)

**Situação:** Imagine que **Alice** quer que **Bob** gere novas carteiras de papel (paper wallets) para ela. Alice tem uma senha mestra super segura que ela não quer compartilhar com ninguém, nem mesmo com Bob. Como Bob pode criar chaves já criptografadas com a senha da Alice, sem nunca saber qual é a senha?

É aqui que o **código intermediário** entra. Ele funciona como uma "autorização" que permite criptografar, mas não descriptografar.

**O Fluxo:**

#### Passo 1: Alice Gera o Código Intermediário

Alice, em seu computador seguro, usa sua senha mestra para gerar um código intermediário. Este código não é a senha dela, mas é derivado dela.

```bash
# No computador da ALICE
$ bip38cli intermediate generate

Enter passphrase: 
# Alice digita sua senha secreta (ex: "do-not-share-this-secret-phrase")
Confirm passphrase: 
# Alice confirma a senha

# Resultado (este é o código que Alice vai compartilhar)
passphrasezctFpQWj9H252m2b2GzQvH3JPTgS3t2fA
```
**Importante:** Alice guarda sua senha secreta e **envia apenas o código `passphrasez...` para o Bob**.

---

#### Passo 2: Bob Gera uma Nova Chave Privada

Bob, em seu sistema (que pode ser online, não precisa ser super seguro para esta parte), gera um novo par de chaves Bitcoin. Ele obtém uma chave privada em formato WIF.

Para este exemplo, vamos supor que Bob gerou a seguinte chave (não-comprimida):

`5Jag5pY5aWJtL3A2YdDB5s2b2GzQvH3JPTgS3t2fA`

---

#### Passo 3: Bob Criptografa a Chave para Alice

Agora, Bob usa o **código intermediário** que recebeu da Alice para criptografar a nova chave privada. Ele **não usa a senha da Alice** (ele nem a conhece).

```bash
# No computador do BOB
$ bip38cli encrypt 5Jag5pY5aWJtL3A2YdDB5s2b2GzQvH3JPTgS3t2fA --intermediate passphrasezctFpQWj9H252m2b2GzQvH3JPTgS3t2fA

# O comando NÃO pede senha! Ele usa o código intermediário.

# Resultado (chave criptografada)
6PRW5VbB4Qj2a2m2b2GzQvH3JPTgS3t2fA1Zytx
```
Bob agora pode entregar esta chave criptografada (`6PRW...`) para a Alice. Ele pode imprimi-la em uma paper wallet ou enviá-la digitalmente. Bob não consegue descriptografar esta chave.

---

#### Passo 4: Alice Verifica e Usa a Chave

Alice recebe a chave criptografada de Bob. Para ter certeza de que funciona, ela a descriptografa usando sua **senha mestra original**.

```bash
# De volta ao computador da ALICE
$ bip38cli decrypt 6PRW5VbB4Qj2a2m2b2GzQvH3JPTgS3t2fA1Zytx

Enter passphrase: 
# Alice digita sua senha secreta original ("do-not-share-this-secret-phrase")

# Resultado (a chave privada que Bob gerou)
5Jag5pY5aWJtL3A2YdDB5s2b2GzQvH3JPTgS3t2fA
```

**Conclusão e Vantagem de Segurança:**

O processo foi um sucesso!

*   **Bob** conseguiu criar uma chave BIP38 para Alice.
*   **Alice** é a única que consegue descriptografá-la.
*   A **senha mestra de Alice nunca saiu** de seu computador seguro.

Este método é ideal para serviços que geram paper wallets, para que os usuários possam ter certeza de que nem o serviço que gerou a carteira tem acesso aos seus fundos.

##  Automação

### Script de Backup em Lote

**backup_carteiras.sh:**
```bash
#!/bin/bash

# Lista de chaves privadas
CHAVES=(
    "5HueCGU8rMjxEXxiPuD5BDku4MkFqeZyd4dZ1jvhTVqvbTLvyTJ"
    "5J3mBbAH58CpQ3Y5RNJpUKPE62SQ5tfcvU2JpbnkeyhfsYB1Jcn"
    # ... mais chaves
)

# Senha (NUNCA hardcode em produção!)
read -s -p "Digite a senha para criptografia: " SENHA
echo

# Arquivo de saída
OUTPUT="backup_criptografado_$(date +%Y%m%d).txt"

echo "Iniciando backup de ${#CHAVES[@]} chaves..."

for i in "${!CHAVES[@]}"; do
    echo "Processando chave $((i+1))..."

    # Criptografa a chave (simula entrada da senha)
    RESULTADO=$(echo -e "${CHAVES[i]}\n$SENHA\n$SENHA" | bip38cli encrypt)

    if [[ $RESULTADO == 6P* ]]; then
        echo "Chave $((i+1)): $RESULTADO" >> "$OUTPUT"
        echo "✓ Chave $((i+1)) criptografada com sucesso"
    else
        echo "✗ Erro na chave $((i+1))"
    fi
done

echo "Backup salvo em: $OUTPUT"
echo "IMPORTANTE: Guarde a senha em local seguro e separado!"
```

### Script de Verificação

**verificar_backup.sh:**
```bash
#!/bin/bash

ARQUIVO_BACKUP="$1"
if [[ ! -f "$ARQUIVO_BACKUP" ]]; then
    echo "Uso: $0 arquivo_backup.txt"
    exit 1
fi

read -s -p "Digite a senha do backup: " SENHA
echo

echo "Verificando integridade do backup..."

while IFS= read -r linha; do
    if [[ $linha == *":"* ]]; then
        NUMERO=$(echo "$linha" | cut -d: -f1)
        CHAVE_CRIPTO=$(echo "$linha" | cut -d: -f2 | tr -d ' ')

        # Tenta descriptografar
        RESULTADO=$(echo -e "$CHAVE_CRIPTO\n$SENHA" | bip38cli decrypt 2>/dev/null)

        if [[ $RESULTADO == 5* ]] || [[ $RESULTADO == K* ]] || [[ $RESULTADO == L* ]]; then
            echo "✓ $NUMERO: OK"
        else
            echo "✗ $NUMERO: ERRO"
        fi
    fi
done < "$ARQUIVO_BACKUP"

echo "Verificação concluída!"
```

##  Segurança e Boas Práticas

###  FAÇA

- **Use senhas muito fortes** (mínimo 15 caracteres, melhor ainda: frase de 6+ palavras)
- **Sempre teste a descriptografia** antes de guardar o backup
- **Use modo offline** (desconecte internet durante uso)
- **Guarde senha e chave em locais DIFERENTES**
- **Faça múltiplas cópias** do backup criptografado
- **Use `--compressed`** para chaves de carteiras modernas
- **Verifique endereços** com `--show-address`

###  NÃO FAÇA

- **NUNCA** use senhas simples (123456, password, etc.)
- **NUNCA** guarde senha junto com chave criptografada
- **NUNCA** confie apenas em uma cópia
- **NUNCA** use em computador infectado/público
- **NUNCA** compartilhe chaves não criptografadas
- **NUNCA** hardcode senhas em scripts
- **NUNCA** use internet durante operações críticas

###  Níveis de Segurança

**Nível 1 - Básico:**
- Senha forte (15+ caracteres)
- Backup criptografado em 2 locais
- Teste de descriptografia

**Nível 2 - Avançado:**
- Frase-senha de 6+ palavras
- Backup em 3+ locais diferentes
- Computador dedicado/live USB
- Verificação periódica

**Nível 3 - Paranóico:**
- Senha gerada por dados físicos
- Air-gapped computer
- Multiple encryption layers
- Dead man's switch para família

##  Solução de Problemas

### Erro: "invalid WIF private key"
```bash
# Problema: Formato de chave inválido
# Solução: Verifique se a chave começa com 5, K ou L
bip38cli encrypt KwYgW8gcxj1JWJXhPSu4Fqwzfhp5Yfi42mdYmMa4XqK7NJxXUSK7
```

### Erro: "incorrect passphrase"
```bash
# Problema: Senha errada na descriptografia
# Solução: Verifique caps lock, teclado, caracteres especiais
# Tente variações da senha se necessário
```

### Erro: "command not found"
```bash
# Problema: bip38cli não está no PATH
# Solução: Use caminho completo
./bin/bip38cli --help

# Ou adicione ao PATH:
export PATH=$PATH:$(pwd)/bin
```

### Performance Lenta
```bash
# Problema: Criptografia demora muito
# Causa: BIP38 usa scrypt (intencionalmente lento)
# Normal: 2-5 segundos por operação
# Se > 30 segundos: computador muito lento ou problema
```

### Chave Não Funciona na Carteira
```bash
# Problema: Chave descriptografada não funciona
# Possíveis causas:
# 1. Formato compressed/uncompressed errado
bip38cli decrypt --show-address sua_chave_6P
# Compare endereço com carteira original

# 2. Rede errada (testnet vs mainnet)
# 3. Erro na senha original
```

##  Exemplos Práticos Completos

### Exemplo 1: Primeira Vez
```bash
# Instalar
wget https://github.com/mannkind/bip38cli-cli/releases/latest/download/bip38cli-linux-amd64.tar.gz
tar -xzf bip38cli-linux-amd64.tar.gz
chmod +x bip38cli

# Testar com chave de exemplo
echo "Testando com chave de exemplo..."
./bip38cli encrypt 5HueCGU8rMjxEXxiPuD5BDku4MkFqeZyd4dZ1jvhTVqvbTLvyTJ
# Digite senha: minhasenhateste123
# Resultado: 6PRVWUbkzz...

# Verificar
./bip38cli decrypt --show-address
# Digite a chave 6P resultante
# Digite senha: minhasenhateste123
# Deve retornar chave original + endereço
```

### Exemplo 2: Backup Real
```bash
# 1. Preparar ambiente seguro
sudo systemctl stop NetworkManager  # Desconectar internet
cd /tmp  # Usar pasta temporária

# 2. Criptografar chave real
./bip38cli encrypt --compressed
# Digite sua chave real da carteira
# Digite senha super forte
# ANOTE o resultado 6P...

# 3. Testar imediatamente
./bip38cli decrypt --show-address
# Digite a chave 6P
# Digite a senha
# CONFIRME que endereço confere

# 4. Limpar e reconectar
cd /
rm -rf /tmp/*
sudo systemctl start NetworkManager
```

---

##  AVISO IMPORTANTE

**Este software lida com chaves privadas Bitcoin. Erros podem resultar em PERDA PERMANENTE de fundos.**

- **SEMPRE teste** com pequenas quantias primeiro
- **SEMPRE faça backup** das chaves originais
- **SEMPRE verifique** endereços após descriptografia
- **NUNCA use** em ambientes não confiáveis

**O autor não se responsabiliza por perdas. Use por sua conta e risco.**

---

##  Suporte

- **Documentação:** README.md
- **Problemas:** https://github.com/mannkind/bip38cli-cli/issues
- **Código:** https://github.com/mannkind/bip38cli-cli

**Lembre-se: Segurança primeiro, sempre!** 