# BIP38CLI - Ferramenta de Criptografia de Chaves Bitcoin

Uma aplicação de linha de comando que implementa o padrão [BIP38](https://github.com/bitcoin/bips/blob/master/bip-0038.mediawiki) para criptografar e descriptografar chaves privadas Bitcoin com proteção por senha.

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go Version](https://img.shields.io/badge/Go-1.24%2B-blue.svg)](https://go.dev/)
[![Release](https://img.shields.io/github/release/carlosrabelo/bip38cli.svg)](https://github.com/carlosrabelo/bip38cli/releases)

## Destaques

- Criptografa e descriptografa chaves em formato WIF com rotinas compatíveis com o BIP38
- Gera e valida códigos intermediários para fluxos de criação de chaves em dois fatores
- Zera buffers de senha assim que possível para reduzir exposição na memória
- Entrada oculta de senha no terminal, com alternância de compressão e modo verboso
- Descoberta inteligente de configuração: `~/.bip38cli.yaml` → `./bip38cli.yaml` → `/etc/bip38cli/config.yaml`
- Gera autocompletes para bash, zsh, fish e PowerShell

## Estrutura do Projeto

```
core/
  cmd/bip38cli/             # Ponto de entrada do binário CLI em Go
  internal/app/cli/         # Comandos Cobra e fluxos de interação
  internal/domain/bip38/    # Lógica de domínio BIP38 e testes
  pkg/                      # Reservado para pacotes públicos futuros
  Makefile                  # Auxiliares de build específicos de Go
bin/
  .gitkeep                  # Placeholder para saída de binários
README-PT.md              # Documentação em português
docs/
  TUTORIAL-EN.md            # Tutorial em inglês
  TUTORIAL-PT.md            # Tutorial em português
scripts/
  install.sh                # Auxiliar de instalação do binário
  uninstall.sh              # Auxiliar de remoção do binário
Makefile                    # Makefile raiz
```

## Início Rápido

### Compilar a partir do código

```bash
git clone https://github.com/carlosrabelo/bip38cli.git
cd bip38cli
make build
./bin/bip38cli --version
```

### Instalar o binário

Instale em `$HOME/.local/bin` (recomendado para usuários sem sudo):

```bash
./scripts/install.sh --user
```

Instale em `/usr/local/bin` (requer permissões apropriadas):

```bash
sudo ./scripts/install.sh
```

Remova o binário mais tarde com o script correspondente:

```bash
./scripts/uninstall.sh --user
# ou
sudo ./scripts/uninstall.sh
```

### Executar via Docker

```bash
# Exibir ajuda dentro do container
./scripts/bip38cli-docker.sh --help

# Executar comandos sem instalar Go localmente
./scripts/bip38cli-docker.sh encrypt --verbose
```

O script mantém os artefatos em `docker/` e constrói a imagem local sob demanda. Veja `docker/README.md` para cenários avançados.

## Uso

### Criptografar uma chave WIF

```bash
bip38cli encrypt KwYgW8gcxj1JWJXhPSu4Fqwzfhp5Yfi42mdYmMa4XqK7NJxXUSK7
# O terminal solicita a senha duas vezes
# Resultado: 6PRV...
```

### Descriptografar uma chave BIP38

```bash
bip38cli decrypt 6PRVWUbkzzsbcVac2qwfssoUJAN1Xhrg6bNk8J7Nzm5H7kxEbn2Nh2ZoGg
# O terminal solicita a senha
# Opcional: --show-address para exibir o endereço derivado
```

### Trabalhar com códigos intermediários

```bash
# Gerar código com metadados de lot/sequence
bip38cli intermediate generate --lot 123 --sequence 456 --use-lot-sequence

# Validar um código existente
bip38cli intermediate validate passphraseabc123...
```

Gerar autocompletes para o seu shell:

```bash
bip38cli completion bash       > /usr/share/bash-completion/completions/bip38cli
bip38cli completion zsh        > /usr/share/zsh/site-functions/_bip38cli
bip38cli completion fish       > ~/.config/fish/completions/bip38cli.fish
bip38cli completion powershell | Out-String | Invoke-Expression
```

## Configuração

O BIP38CLI lê configurações com Viper seguindo esta ordem:

1. `--config /caminho/para/arquivo.yaml`
2. `~/.bip38cli.yaml`
3. `./bip38cli.yaml`
4. `/etc/bip38cli/config.yaml`

Valores padrão embutidos no binário:

```yaml
defaults:
  compressed: true
output:
  format: text
  colors: true
```

Defina `verbose: true` para exibir o caminho de configuração em uso e mensagens adicionais.

## Documentação

- [Tutorial em inglês](docs/TUTORIAL-EN.md)
- [Tutorial em português](docs/TUTORIAL-PT.md)

## Desenvolvimento

Comandos úteis na raiz do repositório:

```bash
make build      # Compila o binário em bin/bip38cli
make fmt        # Formata os arquivos Go com gofmt
make test       # Executa go test ./...
make lint       # Roda golangci-lint (se disponível)
make clean      # Remove artefatos de build
```

O módulo Go reside em `github.com/carlosrabelo/bip38cli/core`. Os testes do domínio BIP38 executam derivação scrypt real e podem levar alguns segundos.

## Notas de Segurança

- Use senhas fortes (15+ caracteres ou frases longas)
- Teste a descriptografia imediatamente após criptografar, antes de armazenar backups
- Mantenha senhas separadas das chaves criptografadas e evite ferramentas de cópia em rede
- Prefira máquinas isoladas (air-gapped) para grandes volumes ou carteiras de alto valor
- Trate códigos intermediários com o mesmo cuidado das chaves criptografadas

## Doações

Se o BIP38CLI é útil para você, considere apoiar o desenvolvimento:

**BTC**: `bc1qw2raw7urfuu2032uyyx9k5pryan5gu6gmz6exm`  
**ETH**: `0xdb4d2517C81bE4FE110E223376dD9B23ca3C762E`  
**SOL**: `A3tNpXSb8rHw2PJYALQeZzwvR4pRWk72YwJdeXGKmS1q`  
**TRX**: `TTznF3FeDCqLmL5gx8GingeahUyLsJJ68A`

## Licença

Distribuído sob a [Licença MIT](LICENSE).
