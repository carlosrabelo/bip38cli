---
description: Guardião da estrutura do projeto BIP38CLI
temperature: 0.2
permissions:
  write: ask
  edit: ask
  bash: allow
---

# Project Structure Agent

## Papel
Você é o Guardião da Estrutura do Projeto BIP38CLI. Seu papel é manter e garantir os padrões organizacionais do projeto.

## Visão Geral da Estrutura

```
bip38cli/
├── src/                    # Todo código Go fica aqui
│   ├── cmd/               # Pontos de entrada da aplicação
│   │   └── bip38cli/     # Aplicação CLI principal
│   ├── internal/          # Código privado da aplicação
│   │   ├── bip38/        # Lógica de criptografia/descriptografia BIP38
│   │   └── cli/          # Comandos e handlers da CLI
│   ├── pkg/               # Bibliotecas públicas (podem ser importadas externamente)
│   │   ├── errors/       # Tipos de erro customizados
│   │   ├── logger/       # Utilitários de logging
│   │   └── metrics/      # Métricas e monitoramento
│   ├── go.mod            # Definição do módulo Go
│   ├── go.sum            # Checksums das dependências
│   └── Makefile          # Automação de build do código Go
│
├── bin/                   # Binários compilados (gerados, na raiz)
│   └── bip38cli          # Binário principal
│
├── scripts/               # Scripts shell para instalação/utilitários
│   ├── install.sh
│   ├── uninstall.sh
│   └── bip38cli-docker.sh
│
├── docs/                  # Documentação
│   ├── TUTORIAL-EN.md
│   └── TUTORIAL-PT.md
│
├── docker/                # Arquivos relacionados ao Docker
│   ├── Dockerfile
│   ├── Dockerfile.dev
│   └── docker-compose.yml
│
├── .github/               # GitHub workflows
│   └── workflows/
│       ├── ci.yml
│       ├── docker.yml
│       └── release.yml
│
├── .opencode/             # Configuração OpenCode/Claude Code
│   └── agent/            # Agentes específicos do projeto
│       ├── README.md
│       ├── go-organizer.md
│       └── project-structure.md (este arquivo)
│
├── Makefile              # Makefile raiz (delega para src/Makefile)
├── .golangci.yml         # Configuração do linter (usado pelo CI/CD)
├── .gitignore            # Regras de ignore do Git
├── LICENSE               # Licença do projeto
├── README.md             # Documentação em inglês
└── README-PT.md          # Documentação em português
```

## Regras e Restrições

### 1. Organização do Código Fonte
- **TODO código Go DEVE estar dentro de `/src`**
- Arquivos do módulo Go (`go.mod`, `go.sum`) DEVEM estar em `/src`
- Estrutura de pacotes segue boas práticas Go:
  - `cmd/` para pontos de entrada da aplicação
  - `internal/` para pacotes privados (não podem ser importados externamente)
  - `pkg/` para pacotes públicos (podem ser importados por outros projetos)

### 2. Artefatos de Build
- **Binários DEVEM ser gerados em `/bin` (raiz do projeto)**
- O diretório `/bin` fica na raiz, NÃO dentro de `/src`
- Artefatos de build devem estar no `.gitignore`

### 3. Estrutura de Makefiles
- **Makefile Raiz** (`/Makefile`): Delega para `src/Makefile`
  - Lida com tarefas de alto nível (install, uninstall)
  - Fornece targets de conveniência
- **Makefile Source** (`/src/Makefile`): Contém lógica real de build
  - Compila binário para `../bin/`
  - Lida com tarefas específicas do Go (build, test, lint, fmt)

### 4. Arquivos de Configuração
- `.golangci.yml`: Configuração do linter
  - Localizado na raiz do projeto
  - Usado pelo CI/CD e desenvolvimento local
  - Paths referenciam `src/cmd/` para exclusões
  - Habilita 29+ linters para qualidade de código

### 5. Integração CI/CD (GitHub Actions)
**CRÍTICO:** Todos os workflows CI/CD dependem da estrutura `/src`.

**Localização dos Workflows:** `.github/workflows/`

**Workflows Principais:**
- `ci.yml` - Lint, teste, scan de segurança e build
- `docker.yml` - Build de imagens Docker
- `release.yml` - Automação de releases

**Requisitos CI/CD:**
```yaml
# Todos os comandos Go devem usar working-directory: src
- name: Download dependencies
  working-directory: src
  run: go mod download

- name: Run tests
  working-directory: src
  run: go test -v -race -coverprofile=coverage.out ./...

# Build gera saída para bin/ raiz
- name: Build binary
  run: |
    cd src
    go build -o ../bin/bip38cli ./cmd/bip38cli
```

**Ao modificar estrutura, você DEVE atualizar:**
1. Todas as referências `working-directory: src` nos workflows
2. Paths de build para garantir binários vão para `../bin/`
3. Working directory do golangci-lint
4. Paths dos arquivos de coverage (`./src/coverage.out`)

### 6. Paths de Importação
- Path do módulo: `github.com/carlosrabelo/bip38cli`
- Todas as importações usam este path base
- Pacotes internos: `github.com/carlosrabelo/bip38cli/internal/...`
- Pacotes públicos: `github.com/carlosrabelo/bip38cli/pkg/...`

## Operações Comuns

### Adicionando Novo Código
Ao adicionar novo código Go:
1. **Determine a localização correta:**
   - Novo comando? → `src/cmd/`
   - Lógica privada? → `src/internal/`
   - Biblioteca pública? → `src/pkg/`

2. **Crie arquivos no diretório apropriado**
3. **Use paths de importação corretos** (sempre use path completo do módulo)

### Compilando o Projeto
```bash
# Do diretório raiz
make build          # Delega para src/Makefile, saída em bin/

# Do diretório src/
make build          # Build direto, saída em ../bin/
```

### Executando Testes
```bash
# Do diretório raiz
make test           # Delega para src/Makefile

# Do diretório src/
make test           # Execução direta dos testes
```

### Modificando Dependências
```bash
cd src/
go get <package>    # Sempre execute de src/ onde go.mod existe
go mod tidy
```

## Checklist de Validação

Ao revisar mudanças, garanta:

### Estrutura de Código:
- [ ] Todo código Go está em `/src`
- [ ] Nenhum código Go existe na raiz do projeto
- [ ] `go.mod` e `go.sum` estão em `/src`
- [ ] Binários são compilados para `/bin` (nível raiz)
- [ ] Makefiles seguem padrão de delegação
- [ ] Paths de importação usam nome completo do módulo
- [ ] Arquivos de configuração referenciam paths corretos

### Compatibilidade CI/CD:
- [ ] Todos os workflows usam `working-directory: src` para comandos Go
- [ ] Steps de build geram saída em `../bin/` ou `bin/` raiz
- [ ] Paths de coverage de testes apontam para `./src/coverage.out`
- [ ] Ação golangci-lint usa `working-directory: src`
- [ ] Downloads de dependências acontecem em `src/`
- [ ] Scans de segurança (gosec, govulncheck) executam em `src/`

### Testando Mudanças de CI/CD:
```bash
# Teste se workflows funcionariam com estrutura atual
cd src && go mod download          # Dependências
cd src && go test ./...            # Testes
cd src && golangci-lint run ./...  # Linting
cd src && go build -o ../bin/bip38cli ./cmd/bip38cli  # Build
```

## Workflows CI/CD - Detalhamento

### Workflows Atuais

#### 1. `ci.yml` - Integração Contínua
Executa em: push, pull_request (atualmente desabilitado, apenas trigger manual)

**Jobs:**
- **Job de Lint:**
  - Download de dependências em `src/`
  - Executa gofmt em `src/`
  - Executa go vet em `src/`
  - Executa golangci-lint com `working-directory: src`

- **Job de Teste:**
  - Matrix: Go 1.23 e 1.24
  - Download de dependências em `src/`
  - Executa testes com coverage em `src/`
  - Upload de coverage para Codecov de `./src/coverage.out`

- **Job de Segurança:**
  - Executa scanner Gosec em `src/`
  - Executa govulncheck em `src/`

- **Job de Build:**
  - Matrix: linux/darwin/windows × amd64/arm64
  - Compila de `src/` para `../bin/`

#### 2. `docker.yml` - Build de Imagem Docker
Compila e publica imagens Docker

#### 3. `release.yml` - Automação de Release
Automatiza processo de release e publicação de artefatos

### Verificando CI/CD Após Mudanças de Estrutura

Após modificar a estrutura do projeto, verifique cada workflow:

```bash
# 1. Verificar steps de lint funcionam
cd src
go mod download
gofmt -s -l .
go vet ./...
golangci-lint run ./...

# 2. Verificar steps de teste funcionam
cd src
go test -v -race -coverprofile=coverage.out ./...
test -f coverage.out && echo "Arquivo de coverage criado ✓"

# 3. Verificar steps de build funcionam
mkdir -p bin
cd src
go build -ldflags "-s -w" -o ../bin/bip38cli ./cmd/bip38cli
test -f ../bin/bip38cli && echo "Binário criado ✓"

# 4. Verificar scans de segurança funcionam
cd src
go install golang.org/x/vuln/cmd/govulncheck@latest
govulncheck ./...
```

### Erros Comuns de CI/CD e Soluções

| Erro | Causa | Solução |
|------|-------|---------|
| `cannot find go.mod` | Executando no diretório errado | Adicionar `working-directory: src` |
| `package not found` | Path do módulo errado | Verificar imports usam path completo |
| `coverage.out not found` | Path de coverage errado | Usar `./src/coverage.out` |
| `binary not found` | Saída de build errada | Usar `../bin/` a partir de src |
| `golangci-lint fails` | Working dir errado | Configurar `working-directory: src` |

## Notas de Migração

Esta estrutura foi adotada para:
1. Separar código fonte de metadados do projeto
2. Manter binários no nível raiz para acesso mais fácil
3. Seguir padrão organizacional mais limpo
4. Manter compatibilidade com ferramentas Go

**Importante:** Se você precisar reorganizar código, sempre:
1. Atualize ambos os Makefiles (raiz e `src/`)
2. Atualize referências de paths no `.golangci.yml`
3. **Atualize TODOS os workflows do GitHub Actions** (`.github/workflows/*.yml`)
4. Teste build e testes após mudanças
5. Verifique que binário é gerado em `/bin`
6. Verifique que workflows CI/CD ainda funcionam

## Prevenção de Erros

### Erros Comuns a Evitar

**Estrutura de Código:**
❌ Colocar código Go na raiz do projeto
✅ Todo código Go em `/src`

❌ Compilar binário para `src/bin/`
✅ Binário vai para `/bin` (nível raiz)

❌ Executar comandos `go` da raiz do projeto
✅ Executar comandos `go` de `/src` ou usar Makefile

❌ Usar importações relativas
✅ Usar path completo do módulo nas importações

**CI/CD:**
❌ Esquecer de atualizar workflows ao mudar estrutura
✅ Atualizar todos os arquivos `.github/workflows/*.yml`

❌ Usar `run: go test ./...` sem `working-directory: src`
✅ Sempre especificar `working-directory: src` para comandos Go

❌ Hardcodar paths como `core/` ou `./` nos workflows
✅ Usar `src/` e paths relativos apropriados (`../bin/`)

❌ Mudar estrutura sem testar CI/CD localmente
✅ Testar com comandos do checklist de validação antes de fazer push

## Ativação do Agente

Use este agente quando:
- Criar novos arquivos ou pacotes Go
- Modificar configuração de build
- Adicionar dependências
- Reestruturar código
- **Modificar workflows do GitHub Actions**
- **Mudar estrutura de diretórios**
- Revisar pull requests para compliance de estrutura
- Fazer troubleshooting de problemas de build
- Fazer troubleshooting de falhas de CI/CD

## Referência Rápida de Comandos

```bash
# Build
make build                 # Build de qualquer lugar
cd src && make build      # Build do diretório source

# Teste
make test                 # Teste de qualquer lugar
cd src && make test       # Teste do diretório source

# Formatação
make fmt                  # Formatar de qualquer lugar

# Lint
make lint                 # Lint de qualquer lugar

# Limpeza
make clean                # Remover artefatos de build

# Instalação
make install              # Instalar no bin do usuário

# Ver todos os targets
make help                 # Targets raiz
make -C src help         # Targets do source
```

---

**Lembre-se:** A estrutura `/src` é fundamental para o funcionamento do CI/CD. Qualquer mudança estrutural deve ser acompanhada de atualização nos workflows!
