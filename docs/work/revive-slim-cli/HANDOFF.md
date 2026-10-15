# Handoff: revive-slim-cli

- Objetivo: app funcionando de novo, o menor possível, com UI para quem não usa terminal e CLI para quem usa.
- Branch: `feat/revive-slim-cli`, sem commit ainda.
- Estado: `main.go` (CLI, config, watch), `opensubtitles.go` (API, conta, erros com status), `gui.go` + `ui.html` (UI web local), `scripts/dist.sh` + `.github/workflows/ci.yml` (pacotes e release), README em en e pt-BR, `docs/RELEASING.md`.
- Login: a chave da API vem embutida no build de tag (secret `OPENSUBTITLES_API_KEY`). Sem conta, 5 downloads por dia; com conta, o usuário preenche usuário e senha em Configurações e o app mostra a cota.
- Verificado: `verify` (vet, testes, build), `actionlint`, `scripts/dist.sh` local gerando os 7 pacotes com macOS universal. UI testada no Chrome com chave falsa embutida: chave recusada traduzida, sem conta, rodapé com versão.
- Não testado: download e login reais (precisa da chave), diálogo nativo, o workflow rodando no GitHub.
- Próxima ação: criar o consumer e o secret, testar com chave real, abrir o PR.
