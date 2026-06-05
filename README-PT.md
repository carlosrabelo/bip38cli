# BIP38CLI - Ferramenta de Criptografia de Chaves Bitcoin

Uma aplicação de linha de comando que implementa o padrão [BIP38](https://github.com/bitcoin/bips/blob/master/bip-0038.mediawiki) para criptografar e descriptografar chaves privadas Bitcoin com proteção por senha.

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go Version](https://img.shields.io/badge/Go-1.26%2B-blue.svg)](https://go.dev/)
[![Release](https://img.shields.io/github/release/carlosrabelo/bip38cli.svg)](https://github.com/carlosrabelo/bip38cli/releases)

## Destaques

- Criptografa e descriptografa chaves em formato WIF com rotinas compatíveis com o BIP38
- Gera novas chaves WIF para qualquer rede Bitcoin com opção de criptografia BIP38
- Exibe endereços SegWit nativos (BIP84) em formato bech32 para chaves comprimidas, mantendo fallback legado para WIFs não comprimidas
- Gera e valida códigos intermediários para fluxos de criação de chaves em dois fatores
- Zera buffers de senha assim que possível para reduzir exposição na memória
- Entrada oculta de senha no terminal, com alternância de compressão e detalhes adicionais no modo verboso
- Flags de linha de comando para todas as opções de configuração
- Gera autocompletes para bash, zsh, fish e PowerShell

## Visão Geral

`bip38cli` é uma ferramenta CLI para trabalhar com chaves privadas Bitcoin criptografadas com BIP38. O foco é oferecer fluxos interativos e automatizáveis, mantendo segredos fora do histórico do shell ao solicitar senhas de forma segura.

Se você quiser um passo a passo, veja:

- [Tutorial em inglês](docs/TUTORIAL-EN.md)
- [Tutorial em português](docs/TUTORIAL-PT.md)

## Pré-requisitos

- **Go 1.26+** — necessário para compilar a partir do código; [download](https://go.dev/dl/)
- **golangci-lint** — necessário para `make lint`; instale com `go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest`

## Instalação

### Compilar a partir do código

```bash
git clone https://github.com/carlosrabelo/bip38cli.git
cd bip38cli
make build
./bin/bip38cli --version
```

Instale em `$HOME/.local/bin` (recomendado para usuários sem sudo):

```bash
make install
```

Instale em `/usr/local/bin` (requer permissões apropriadas):

```bash
sudo make install
```

Remova o binário depois:

```bash
make uninstall
# ou
sudo make uninstall
```

## Início Rápido

```bash
make build
./bin/bip38cli --help
```


## Uso

### Gerar uma nova carteira (WIF)

```bash
# Gerar uma chave WIF comprimida na mainnet (padrão)
bip38cli wallet generate

# Escolher outra rede (ex.: testnet) e exibir o endereço derivado
bip38cli wallet generate --network testnet --show-address

# Criptografar a chave gerada com BIP38 (senha solicitada no terminal)
bip38cli wallet generate --encrypt

# Saída em JSON (inclui WIF, rede, compressão e, opcionalmente, endereço/chave BIP38)
bip38cli wallet generate --output-format json --show-address

# Gerar endereço legado P2PKH (BIP44) em vez de bech32 (padrão BIP84)
bip38cli wallet generate --address-type bip44 --show-address

# Gerar uma chave não comprimida (endereços legados implícitos)
bip38cli wallet generate --uncompressed
```

> Chaves comprimidas seguem o padrão BIP84 (bech32). Caso você gere explicitamente uma chave não comprimida, o CLI retorna o endereço legado P2PKH.

Use `--address-type bip44` sempre que precisar do endereço legado P2PKH para compatibilidade com carteiras antigas.

### Inspecionar uma WIF existente

```bash
# Exibir rede, compressão e endereço
bip38cli wallet inspect 5HueCGU8rMjxEXxiPuD5BDku4MkFqeZyd4dZ1jvhTVqvbTLvyTJ

# Saída em JSON (para automações)
bip38cli wallet inspect --output-format json KwYgW8gcxj1JWJXhPSu4Fqwzfhp5Yfi42mdYmMa4XqK7NJxXUSK7

# Forçar saída legada P2PKH
bip38cli wallet inspect --address-type bip44 KwYgW8gcxj1JWJXhPSu4Fqwzfhp5Yfi42mdYmMa4XqK7NJxXUSK7
```

### Criptografar uma chave WIF

```bash
# Criptografia básica
bip38cli encrypt KwYgW8gcxj1JWJXhPSu4Fqwzfhp5Yfi42mdYmMa4XqK7NJxXUSK7
# O terminal solicita a senha duas vezes
# Resultado: 6PRV...

# Forçar formato comprimido (igual ao padrão)
bip38cli encrypt --compressed KwYgW8gcxj1JWJXhPSu4Fqwzfhp5Yfi42mdYmMa4XqK7NJxXUSK7

# Forçar formato não comprimido
bip38cli encrypt --uncompressed KwYgW8gcxj1JWJXhPSu4Fqwzfhp5Yfi42mdYmMa4XqK7NJxXUSK7

# Usar flag global de compressão (padrão: comprimido)
bip38cli encrypt --compressed KwYgW8gcxj1JWJXhPSu4Fqwzfhp5Yfi42mdYmMa4XqK7NJxXUSK7

# Usar formato não comprimido como padrão
bip38cli encrypt --compressed=false KwYgW8gcxj1JWJXhPSu4Fqwzfhp5Yfi42mdYmMa4XqK7NJxXUSK7

# Saída em JSON
bip38cli encrypt --output-format json KwYgW8gcxj1JWJXhPSu4Fqwzfhp5Yfi42mdYmMa4XqK7NJxXUSK7
```

### Descriptografar uma chave BIP38

```bash
# Descriptografia básica
bip38cli decrypt 6PRVWUbkzzsbcVac2qwfssoUJAN1Xhrg6bNk8J7Nzm5H7kxEbn2Nh2ZoGg
# O terminal solicita a senha

# Exibir o endereço derivado
bip38cli decrypt --show-address 6PRVWUbkzzsbcVac2qwfssoUJAN1Xhrg6bNk8J7Nzm5H7kxEbn2Nh2ZoGg

# Saída em JSON com endereço
bip38cli decrypt --show-address --output-format json 6PRVWUbkzzsbcVac2qwfssoUJAN1Xhrg6bNk8J7Nzm5H7kxEbn2Nh2ZoGg
```

### Trabalhar com códigos intermediários

```bash
# Gerar código intermediário básico
bip38cli intermediate generate

# Gerar código com metadados de lot/sequence
bip38cli intermediate generate --lot 123 --sequence 456 --use-lot-sequence

# Gerar com saída JSON
bip38cli intermediate generate --output-format json --lot 123 --sequence 456 --use-lot-sequence

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

O BIP38CLI utiliza apenas flags de linha de comando para ajustes de comportamento; não há suporte a arquivos de configuração externos.

Flags globais:
- `--verbose, -v`: ativa saída detalhada com logs adicionais.
- `--output-format`: controla o formato (`text`|`json`, padrão: `text`).
- `--compressed, -c`: define o uso padrão de chaves comprimidas.
- `--uncompressed`: força o formato não comprimido (sobrescreve `--compressed`).

Flags específicas por comando:
- `encrypt --compressed`: gera chave criptografada em formato comprimido.
- `encrypt --uncompressed`: gera chave criptografada em formato não comprimido.
- `decrypt --show-address`: exibe o endereço Bitcoin derivado da chave descriptografada.
- `decrypt --address-type <bip84|bip44>`: controla o formato do endereço ao usar `--show-address` (padrão: `bip84`).
- `intermediate generate --lot <número>`: informa o número de lote (0-1048575).
- `intermediate generate --sequence <número>`: informa o número de sequência (0-4095).
- `intermediate generate --use-lot-sequence`: inclui lote e sequência no código intermediário.
- `wallet generate --address-type <bip84|bip44>`: escolhe entre bech32 (bip84) ou legado P2PKH (bip44).
- `wallet generate --uncompressed`: produz uma chave não comprimida (endereços legados).
- `wallet inspect --address-type <bip84|bip44>`: inspeciona WIFs usando o tipo de endereço desejado.
- `wallet generate --network <nome>`: escolhe a rede (`mainnet`, `testnet`, `regtest`, `simnet`, `signet`).
- `wallet generate --encrypt`: envolve a chave recém-gerada com BIP38 (senha interativa).
- `wallet generate --show-address`: apresenta o endereço Bitcoin derivado da nova chave.

### Exemplos com saída JSON

```bash
# Criptografar com saída JSON
bip38cli encrypt --output-format json KwYgW8gcxj1JWJXhPSu4Fqwzfhp5Yfi42mdYmMa4XqK7NJxXUSK7
# Saída: {"encrypted_key": "6PRV...", "compressed": true}

# Descriptografar com saída JSON e endereço
bip38cli decrypt --show-address --output-format json 6PRVWUbkzzsbcVac2qwfssoUJAN1Xhrg6bNk8J7Nzm5H7kxEbn2Nh2ZoGg
# Saída: {"private_key": "KwYg...", "compressed": true, "address": "bc1qklnjad76qxxxy833ggfjsjyjc29vdrgnpnju5d"}
```

## Estrutura do Projeto

```
bip38cli/cmd/bip38cli/    # Ponto de entrada do binário CLI em Go
bip38cli/internal/bip38/  # Lógica de domínio BIP38 e testes
bip38cli/internal/cli/    # Comandos Cobra e fluxos de interação
bip38cli/internal/errors/ # Tipos de erro
bip38cli/internal/logger/ # Utilitários de log
bip38cli/internal/metrics/# Utilitários de métricas
bin/
  .gitkeep                  # Placeholder para saída de binários
.make/                     # Scripts de build/test/install
README-PT.md              # Documentação em português
docs/
  TUTORIAL-EN.md            # Tutorial em inglês
  TUTORIAL-PT.md            # Tutorial em português
Makefile                    # Makefile raiz
```

## Desenvolvimento

Comandos úteis na raiz do repositório:

```bash
make build      # Compila o binário em bin/bip38cli
make fmt        # Formata os arquivos Go
make test       # Executa go test ./...
make quality    # Roda todas as checagens de qualidade
make lint       # Roda golangci-lint (se disponível)
make install    # Instala em $HOME/.local/bin
make uninstall  # Remove de $HOME/.local/bin
make version    # Mostra informações de versão
make clean      # Remove artefatos de build
```

O módulo Go reside em `github.com/carlosrabelo/bip38cli`. Os testes do domínio BIP38 executam derivação scrypt real e podem levar alguns segundos.

## Licença

Distribuído sob a [Licença MIT](LICENSE).
