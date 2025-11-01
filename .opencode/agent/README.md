---
description: Convenções e padrões para criação de agentes do projeto
---

# Agent Conventions

## ⚠️ REGRAS FUNDAMENTAIS

**TODOS os agentes DEVEM seguir estas regras ao interagir com o projeto:**

### 1. NUNCA mencione IA, Agentes ou Ferramentas Automáticas

Ao criar commits, documentação, código, comentários ou qualquer conteúdo do projeto:

❌ **PROIBIDO mencionar:**
- IA, AI, Inteligência Artificial, Machine Learning
- Claude, Gemini, GPT, Codex, OpenCode, ChatGPT
- Assistente, Assistant, Agent, Agente
- Bot, Automation, Automated (em contexto de autoria)
- "Gerado por...", "Criado automaticamente por..."
- "Sugerido pela IA", "Otimizado por Claude"

✅ **SEMPRE escrever como humano:**
- Use tom profissional e neutro
- Commits descritivos e objetivos
- Documentação clara sem revelar ferramentas
- Comentários técnicos e úteis
- Exemplos: "Refatora código", "Adiciona feature", "Corrige bug"

**Esta regra se aplica a:**
- Mensagens de commit
- Descrições de Pull Request
- README e documentação
- Comentários de código
- Issues e comentários
- Changelog
- Tutoriais
- Qualquer conteúdo público do projeto

### 2. 🚫 NUNCA Fazer Commits Automáticos em Sequência

**PROIBIDO fazer múltiplos commits sem aprovação explícita do usuário.**

**Regra:**
- **UM commit = UMA aprovação do usuário**
- Após fazer um commit, **SEMPRE PERGUNTAR** se deve fazer outro
- **NUNCA assumir** que o usuário quer commits adicionais
- Se houver múltiplas mudanças, **LISTAR** e deixar o usuário decidir

**Comportamento correto:**
```
Agente: [faz commit 1]
Agente: "Commit realizado. Há outras mudanças. Deseja que eu faça outro commit?"
Usuário: "Sim"
Agente: [faz commit 2]
Agente: "Feito. Ainda há mudanças em X e Y. Deseja continuar?"
```

**Comportamento PROIBIDO:**
```
Agente: [faz commit 1]
Agente: [faz commit 2 sem perguntar]  ← ❌ NUNCA FAZER ISSO
Agente: [faz commit 3 sem perguntar]  ← ❌ NUNCA FAZER ISSO
```

**Exceção:** Somente quando o usuário explicitamente pedir "faça todos os commits" ou similar. Mesmo assim, **confirmar** antes de prosseguir.

## Regras Gerais

### 1. Idioma e Nomenclatura

**REGRAS FUNDAMENTAIS:**
- **Nome do arquivo:** INGLÊS (kebab-case) - Ex: `test-manager.md`
- **Título do agente (# ...):** INGLÊS - Ex: `# Test Manager`
- **Todo o conteúdo:** PORTUGUÊS (pt-BR)
- **Description (frontmatter):** PORTUGUÊS

✅ Correto:
```markdown
---
description: Guardião da estrutura do projeto (português)
---

# Project Structure Agent  ← INGLÊS

## Papel  ← Conteúdo em PORTUGUÊS
Você é o guardião da estrutura do projeto...
```

❌ Incorreto:
```markdown
# Agente de Estrutura  ← Título em português (errado!)
You are the guardian...  ← Conteúdo em inglês (errado!)
```

❌ Também incorreto:
```markdown
# Project Structure Agent  ← Título correto
You are the guardian...  ← Conteúdo em inglês (errado!)
```

### 2. Formato do Arquivo

Cada agente deve ter:

```markdown
---
description: Descrição curta do agente (pt-BR)
temperature: 0.1-0.7  # Opcional, padrão contextual
model: opencode/grok-code  # Opcional
permissions:
  write: ask|allow|deny
  edit: ask|allow|deny
  bash: ask|allow|deny
---

# Nome do Agente

Descrição detalhada do papel do agente...

## Seções principais...
```

### 3. Estrutura Recomendada

```markdown
# Título do Agente

## Papel
Defina claramente o papel e responsabilidade

## [Contexto/Visão Geral]
Informações de contexto necessárias

## Regras e Restrições
O que o agente DEVE e NÃO DEVE fazer

## Operações Comuns
Exemplos práticos de uso

## Checklist de Validação
Items para verificar ao usar o agente

## Prevenção de Erros
Erros comuns e como evitar

## Ativação do Agente
Quando usar este agente

## Referência Rápida
Comandos/exemplos úteis
```

### 4. Tom e Estilo

- **Tom:** Profissional, direto, objetivo
- **Perspectiva:** Segunda pessoa (você) ao dar instruções
- **Exemplos:** Sempre incluir código/comandos práticos
- **Símbolos:** Usar ✅ ❌ 📋 🎯 para destacar pontos importantes

### 5. Permissões

Defina permissões apropriadas:

```yaml
permissions:
  write: ask    # Pedir confirmação antes de criar arquivos
  edit: ask     # Pedir confirmação antes de editar
  bash: allow   # Permitir comandos bash (quando seguro)
```

**Regra de ouro:** Quando em dúvida, use `ask`.

## Tipos de Agentes do Projeto

### 1. Agentes de Estrutura
- **Exemplo:** `project-structure.md`
- **Propósito:** Manter organização do código
- **Permissões:** Geralmente `ask` para write/edit
- **Privacidade:** Não mencionar IA em commits de estrutura

### 2. Agentes de Qualidade
- **Exemplo:** `go-organizer.md`
- **Propósito:** Garantir padrões de código
- **Permissões:** `ask` para edições, `allow` para bash
- **Privacidade:** Commits devem ser neutros ("Organiza código Go")

### 3. Agentes de Documentação e Commits
- **Exemplo:** `commits-and-documentation.md`
- **Propósito:** Criar commits, README, tutoriais e documentação
- **Permissões:** `ask` para write/edit
- **Privacidade:** ⚠️ CRÍTICO - Nunca mencionar IA/Claude/Agente

### 4. Agentes de Build/Deploy
- **Propósito:** Automatizar processos
- **Permissões:** `allow` para bash, `ask` para edições críticas
- **Privacidade:** Logs e mensagens devem parecer humanas

## Checklist para Criar Novo Agente

- [ ] **Título do agente em INGLÊS** (ex: "# Project Structure Agent")
- [ ] **Conteúdo do agente em PORTUGUÊS (pt-BR)**
- [ ] Frontmatter com description (em português)
- [ ] Título claro e descritivo (em inglês)
- [ ] Seção "Papel" definida (em português)
- [ ] **Regra de privacidade incluída** (não mencionar IA)
- [ ] Regras e restrições claras (em português)
- [ ] Exemplos práticos incluídos
- [ ] Permissões apropriadas definidas
- [ ] Exemplos de código/comandos
- [ ] Erros comuns documentados
- [ ] Quando usar o agente está claro
- [ ] **Exemplos de commits/docs sem mencionar IA**

## Nomenclatura de Arquivos

**REGRA OBRIGATÓRIA:** Nomes de arquivos de agentes DEVEM ser em **INGLÊS** usando **kebab-case**.

### ✅ Correto (inglês + kebab-case):

```
project-structure.md           # Estrutura do projeto
commits-and-documentation.md   # Commits e documentação
go-organizer.md                # Organizador Go
test-manager.md           # Gerenciador de testes
ci-cd-helper.md           # Auxiliar CI/CD
deployment-agent.md       # Agente de deploy
```

### ❌ Incorreto:

```
# Português (ERRADO)
estrutura-projeto.md      # ❌ Nome em português
gerenciador-testes.md     # ❌ Nome em português
agente-deploy.md          # ❌ Nome em português

# Formato errado (ERRADO)
ProjectStructure.md       # ❌ PascalCase
project_structure.md      # ❌ snake_case
PROJECTSTRUCTURE.md       # ❌ UPPERCASE
```

### Regra Completa:

| Elemento | Idioma | Formato | Exemplo |
|----------|--------|---------|---------|
| **Nome do arquivo** | **INGLÊS** | **kebab-case** | `test-manager.md` |
| Título (# ...) | INGLÊS | Sentence Case | `# Test Manager` |
| Description | Português | - | `Gerencia testes...` |
| Conteúdo | Português | - | `Você é responsável...` |

## Localização

Todos os agentes devem estar em:
```
.opencode/agent/
├── README.md                     # Este arquivo (convenções)
├── project-structure.md          # Agente de estrutura do projeto
├── go-organizer.md               # Agente de organização de código Go
├── commits-and-documentation.md  # Agente de commits e documentação ⚠️
└── [outros-agentes].md
```

**Nota:** O agente `commits-and-documentation.md` é especialmente importante pois garante que toda comunicação pública do projeto seja profissional e não revele ferramentas de automação.

## Manutenção

### Quando Atualizar um Agente

1. **Mudança de estrutura do projeto** → Atualizar agentes relacionados
2. **Novas convenções adotadas** → Documentar no agente apropriado
3. **Erros recorrentes** → Adicionar na seção de prevenção
4. **CI/CD modificado** → Atualizar referências

### Sincronização

Agentes relacionados devem estar sincronizados:
- `project-structure.md` ↔ Estrutura real do projeto
- `go-organizer.md` ↔ Padrões de código atuais
- Agentes de CI/CD ↔ Workflows do GitHub Actions

## Exemplos

### ✅ Bom Agente

```markdown
---
description: Organiza testes seguindo convenções do projeto
temperature: 0.1
permissions:
  write: ask
  edit: ask
  bash: allow
---

# Test Organizer

## Papel
Você é responsável por manter os testes organizados...

## Regras
✅ Testes devem estar em `*_test.go`
❌ Nunca remover testes existentes

## Exemplos
\```bash
go test ./...
\```
```

### ❌ Agente Problemático

```markdown
# Ajudante de Testes

Você ajuda com testes...

Sem estrutura clara, título em português (deveria ser inglês), sem frontmatter.
```

**Problemas:**
- ❌ Título em português (deve ser inglês)
- ❌ Sem frontmatter
- ❌ Sem estrutura definida

## Referências

- [Documentação OpenCode Agents](https://docs.opencode.com/agents)
- [Go Documentation Guidelines](https://go.dev/doc/effective_go)
- Projetos de referência: Este próprio repositório

---

**Lembre-se:** Agentes são documentação viva. Mantenha-os atualizados!
