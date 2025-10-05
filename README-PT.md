# BIP38CLI - Ferramenta de Criptografia de Chaves Bitcoin

Uma ferramenta moderna de linha de comando para criptografia e descriptografia de chaves privadas BIP38 (Bitcoin Improvement Proposal 38), escrita em Go.

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go Version](https://img.shields.io/badge/Go-1.24%2B-blue.svg)](https://golang.org)
[![Release](https://img.shields.io/github/release/carlosrabelo/bip38cli.svg)](https://github.com/carlosrabelo/bip38cli/releases)

## Características

- **Criptografar/Descriptografar chaves privadas Bitcoin** usando padrão BIP38
- **Gerar códigos de senha intermediários** para criptografia de dois fatores
- **Suporte para chaves comprimidas e descomprimidas**
- **Manuseio seguro de senhas** com entrada oculta
- **Rápido e eficiente** - construído com Go e BTCSuite
- **Suporte multiplataforma** (Linux, macOS, Windows)
- **Autocompletar** para bash, zsh, fish e PowerShell

## Início Rápido

### Instalação

**Baixar binário pré-compilado:**
```bash
# Linux/macOS
curl -LO https://github.com/carlosrabelo/bip38cli/releases/latest/download/bip38cli-linux-amd64.tar.gz
tar -xzf bip38cli-*.tar.gz
sudo mv bip38cli /usr/local/bin/

# Verificar instalação
bip38cli --version
```

**Compilar do código fonte:**
```bash
git clone https://github.com/carlosrabelo/bip38cli.git
cd bip38cli
make build
./bin/bip38cli --version
```

### Uso Básico

**Criptografar uma chave privada:**
```bash
bip38cli encrypt 5HueCGU8rMjxEXxiPuD5BDku4MkFqeZyd4dZ1jvhTVqvbTLvyTJ
# Digite a senha quando solicitado
# Saída: 6PRVWUbkzzsbcVac2qwfssoUJAN1Xhrg6bNk8J7Nzm5H7kxEbn2Nh2ZoGg
```

**Descriptografar uma chave privada:**
```bash
bip38cli decrypt 6PRVWUbkzzsbcVac2qwfssoUJAN1Xhrg6bNk8J7Nzm5H7kxEbn2Nh2ZoGg
# Digite a senha quando solicitado
# Saída: 5HueCGU8rMjxEXxiPuD5BDku4MkFqeZyd4dZ1jvhTVqvbTLvyTJ
```

## Comandos

### `encrypt` - Criptografar Chaves Privadas
```bash
# Modo interativo (mais seguro)
bip38cli encrypt

# Modo direto
bip38cli encrypt [CHAVE_PRIVADA_WIF]

# Forçar formato comprimido
bip38cli encrypt --compressed [CHAVE_PRIVADA_WIF]
```

### `decrypt` - Descriptografar Chaves Privadas
```bash
# Modo interativo
bip38cli decrypt

# Modo direto
bip38cli decrypt [CHAVE_CRIPTOGRAFADA]

# Mostrar endereço Bitcoin
bip38cli decrypt --show-address [CHAVE_CRIPTOGRAFADA]
```

### `intermediate` - Criptografia de Dois Fatores
```bash
# Gerar código intermediário
bip38cli intermediate generate

# Gerar com números de lote/sequência
bip38cli intermediate generate --lot 123 --sequence 456

# Validar código intermediário
bip38cli intermediate validate [CODIGO_INTERMEDIARIO]
```

## Autocompletar

O BIP38CLI suporta autocompletar para bash, zsh, fish e PowerShell.

### Bash
```bash
# Adicionar ao ~/.bashrc
echo 'source <(bip38cli completion bash)' >> ~/.bashrc
source ~/.bashrc
```

### Zsh
```bash
# Adicionar ao ~/.zshrc
echo 'source <(bip38cli completion zsh)' >> ~/.zshrc
source ~/.zshrc
```

### Fish
```bash
# Gerar arquivo de autocompletar
bip38cli completion fish > ~/.config/fish/completions/bip38cli.fish
```

### PowerShell
```powershell
# Adicionar ao perfil do PowerShell
bip38cli completion powershell | Out-String | Invoke-Expression
```

## Configuração

Crie um arquivo de configuração em `~/.bip38cli.yaml`:

```yaml
# Configurações de comportamento padrão
verbose: false
compressed: true

# Preferências de formato de saída
output:
  format: "text"  # text, json
  colors: true
```

## O que é BIP38?

BIP38 é um Bitcoin Improvement Proposal que define um método para criptografar chaves privadas Bitcoin com uma senha. Isso permite:

1. **Armazenamento protegido por senha** - Chaves podem ser armazenadas ou transmitidas com segurança
2. **Segurança de dois fatores** - Usando códigos intermediários para geração segura de chaves
3. **Formato padronizado** - Compatível com outras implementações BIP38
4. **Backups seguros** - Chaves criptografadas podem ser armazenadas em múltiplos locais

## Melhores Práticas de Segurança

- **Use senhas fortes** (15+ caracteres ou frases de 6+ palavras)
- **Sempre teste a descriptografia** antes de armazenar chaves criptografadas
- **Use ambientes offline** para operações críticas
- **Armazene senhas separadamente** das chaves criptografadas
- **Faça múltiplos backups** de chaves criptografadas
- **Verifique endereços** após descriptografia

- **Nunca use senhas fracas** (123456, password, etc.)
- **Nunca armazene senhas com chaves criptografadas**
- **Nunca use em computadores infectados**
- **Nunca compartilhe chaves privadas não criptografadas**

## Exemplos

### Backup Seguro de Carteira
```bash
# 1. Exporte a chave privada da sua carteira
# 2. Criptografe com senha forte
bip38cli encrypt --compressed
# 3. Teste a descriptografia imediatamente
bip38cli decrypt --show-address
# 4. Armazene chave criptografada e senha separadamente
```

### Processamento em Lote
```bash
# Criptografar múltiplas chaves de arquivo
while read -r key; do
    echo "Processando: $key"
    echo "$key" | bip38cli encrypt
done < chaves_privadas.txt
```

### Herança Digital
```bash
# Criar backup criptografado para família
bip38cli encrypt 5HueCGU8rMjxEXxiPuD5BDku4MkFqeZyd4dZ1jvhTVqvbTLvyTJ
# Use senha memorável que a família possa deduzir
# Deixe instruções claras para descriptografia
```

## Desenvolvimento

### Requisitos
- Go 1.24.0 ou posterior
- Make

### Compilação
```bash
# Instalar dependências
make deps

# Executar testes
make test

# Compilar binário
make build

# Executar todas as verificações
make all
```

### Testes
```bash
# Executar testes com cobertura
make test-coverage

# Executar linting
make lint

# Benchmark de performance
go test -bench=. ./internal/bip38/
```

## Referência da API

### Códigos de Saída
- `0` - Sucesso
- `1` - Erro geral
- `2` - Argumentos inválidos
- `3` - Falha na criptografia/descriptografia
- `4` - Formato de chave inválido

### Variáveis de Ambiente
- `BIP38CLI_CONFIG` - Caminho do arquivo de configuração
- `BIP38CLI_VERBOSE` - Habilitar saída detalhada
- `NO_COLOR` - Desabilitar saída colorida

## Contribuindo

1. Fork o repositório
2. Crie uma branch de funcionalidade (`git checkout -b feature/funcionalidade-incrivel`)
3. Faça suas alterações
4. Adicione testes para nova funcionalidade
5. Execute `make all` para verificar se tudo funciona
6. Commit suas alterações (`git commit -m 'Adicionar funcionalidade incrível'`)
7. Push para a branch (`git push origin feature/funcionalidade-incrivel`)
8. Abra um Pull Request

## Licença

Este projeto está licenciado sob a Licença MIT - veja o arquivo [LICENSE](LICENSE) para detalhes.

## Agradecimentos

- Construído com [BTCSuite](https://github.com/btcsuite) - a biblioteca Bitcoin padrão da indústria para Go
- Usa [Cobra](https://github.com/spf13/cobra) para framework CLI
- Implementa a [especificação BIP38](https://github.com/bitcoin/bips/blob/master/bip-0038.mediawiki)

## Suporte

- **Documentação**: [Tutoriais completos](TUTORIAL-PT.md)
- **Problemas**: [GitHub Issues](https://github.com/carlosrabelo/bip38cli/issues)
- **Discussões**: [GitHub Discussions](https://github.com/carlosrabelo/bip38cli/discussions)

---

**⚠️ Importante**: Este software lida com chaves privadas Bitcoin. Sempre teste com pequenas quantias primeiro e certifique-se de entender os riscos. Os autores não são responsáveis por qualquer perda de fundos.