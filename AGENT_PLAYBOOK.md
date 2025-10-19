# PLAYBOOK DO AGENTE

## 1. Preparação Inicial
- Garantir que a cópia local está alinhada com `main` (`git fetch --all` + `git status` limpo).
- Conferir `STATE.md` e `CONTEXT.md` para entender estado atual e convenções vigentes.
- Operar sempre em português e sem emojis, conforme política do projeto.
- Antes de alterações, validar se o ambiente possui Go 1.24+, `golangci-lint` (quando for executar lint) e os scripts marcados como executáveis.

## 2. Rotina de Execução por Tarefa
- Ler a solicitação do usuário e confirmar escopo; registrar dúvidas imediatamente.
- Mapear diretórios relevantes (`core/internal/domain/bip38`, `core/internal/app/cli`, etc.) antes de editar.
- Para mudanças em código Go:
  - Escrever testes novos quando necessário para cobrir comportamento.
  - Executar `make fmt` (ou `go fmt` localizado) para manter estilo.
  - Rodar `make test`; se lint for requerido ou tocado por mudança, executar `make lint`.
- Para ajustes em scripts ou documentos:
  - Checar reflexos nas duas línguas quando houver manual em EN/PT.
  - Atualizar `CONTEXT.md` ou `STATE.md` quando decisão de arquitetura/processo mudar.

## 3. Padrões de Qualidade
- Manter suporte a saídas `text` e `json` em qualquer funcionalidade exposta.
- Garantir cobertura para fluxos sensíveis (ex.: WIF inválido, senha incorreta, uso de `--verbose`).
- Explicar em comentários apenas trechos complexos; evitar comentários redundantes.
- Nenhum segredo ou chave deve aparecer em logs, fixtures ou mensagens.

## 4. Segurança e Compliance
- Reutilizar utilitários de `secureZero` ao lidar com buffers sensíveis.
- Confirmar que leitura de senhas continua sendo feita via `term.ReadPassword`.
- Verificar que scripts não expõem variáveis sensíveis nem exigem sudo desnecessário.
- Respeitar sandbox: não executar comandos destrutivos (`git reset --hard`, `rm -rf`) sem pedido explícito.
- Obedecer aos hooks de commit configurados; ajustar o fluxo quando um hook bloquear a alteração.

## 5. Documentação e Sincronia
- `README.md` e `README-PT.md` devem permanecer equivalentes; qualquer adição precisa de tradução correspondente.
- `docs/TUTORIAL-EN.md` e `docs/TUTORIAL-PT.md` seguem mesma regra; garantir que exemplos de CLI estejam atualizados.
- Atualizar `STATE.md` quando fluxos de processo, ferramentas obrigatórias ou políticas mudarem.
- Atualizar `CONTEXT.md` quando arquitetura, dependências chave ou convenções estruturais forem alteradas.

## 6. Checklist de Entrega
- Código formatado (`make fmt`).
- Testes e lint relevantes executados, com resultado positivo.
- Documentação perturbada revisada nas duas línguas.
- Mensagem final ao usuário menciona arquivos tocados, resultados de testes e próximos passos naturais (se existirem).
- Git working tree limpo ou com mudanças intencionais únicas ao escopo da tarefa.
