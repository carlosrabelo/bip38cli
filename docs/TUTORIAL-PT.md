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
wget https://github.com/carlosrabelo/bip38cli/releases/latest/download/bip38cli-linux-amd64.tar.gz

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
git clone https://github.com/carlosrabelo/bip38cli.git
cd bip38cli

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

### Cenário 1: Backup Seguro de Carteira

**Situação:** Você tem uma carteira Bitcoin e quer fazer backup ultra-seguro.

```bash
# 1. Exporte a chave privada da sua carteira
# (no Electrum: Wallet > Private Keys > Export)

# 2. Criptografe com senha forte
bip38cli encrypt --compressed
# Digite sua chave privada
# Digite senha MUITO forte (ex: frase de 6 palavras)
# ANOTE a chave criptografada resultante

# 3. TESTE a descriptografia
bip38cli decrypt --show-address
# Digite a chave criptografada
# Digite a senha
# CONFIRME que o endereço confere com sua carteira

# 4. Guarde com segurança:
# - Chave criptografada em nuvem/papel
# - Senha em local DIFERENTE
```

### Cenário 2: Herança Digital

**Situação:** Deixar Bitcoin para família acessar em emergência.

```bash
# 1. Criptografe suas chaves
bip38cli encrypt 5HueCGU8rMjxEXxiPuD5BDku4MkFqeZyd4dZ1jvhTVqvbTLvyTJ
# Use senha que família pode descobrir (ex: data nascimento + nome pet)

# 2. Crie instruções simples:
echo "Em emergência:
1. Baixe bip38cli de github.com/carlosrabelo/bip38cli
2. Execute: bip38cli decrypt 6PRVWUbkzzsbcVac2qwfssoUJAN1Xhrg6bNk8J7Nzm5H7kxEbn2Nh2ZoGg
3. Senha é: [sua pista aqui]
4. Use a chave WIF resultante em qualquer carteira Bitcoin" > instrucoes_heranca.txt
```

### Cenário 3: Duas Pessoas, Máxima Segurança

**Situação:** Uma pessoa conhece a senha, outra gera as chaves.

```bash
# Pessoa A (conhece senha):
bip38cli intermediate generate
# Digite senha secreta
# Envie código intermediário para Pessoa B

# Pessoa B (gera chaves, mas não conhece senha):
# [usaria comando 'generate' - não implementado ainda]
# Envia chaves criptografadas para Pessoa A

# Pessoa A pode descriptografar quando necessário
```

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
wget https://github.com/carlosrabelo/bip38cli/releases/latest/download/bip38cli-linux-amd64.tar.gz
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
- **Problemas:** https://github.com/carlosrabelo/bip38cli/issues
- **Código:** https://github.com/carlosrabelo/bip38cli

**Lembre-se: Segurança primeiro, sempre!** 