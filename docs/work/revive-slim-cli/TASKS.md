# Reviver o subtitlewatcher: CLI + UI enxuta

- [x] Remover GUI Fyne, recursos gerados, registro do Windows e scripts antigos
- [x] Cliente da API REST do OpenSubtitles (busca por hash e nome, download, login opcional)
- [x] Comandos `get` e `watch` só com stdlib
- [x] UI web local (sem argumentos): baixar arquivo, monitorar pasta, atividade, configurações, pt-BR/EN
- [x] Config salva em arquivo e usada também pela CLI
- [x] Login: chave embutida no build, conta opcional, teste de login com cota, erros traduzidos
- [x] `scripts/dist.sh` e CI: testes em 3 sistemas, pacotes para macOS universal, Windows, Linux, FreeBSD
- [x] README em inglês e pt-BR com screenshots, `docs/RELEASING.md`
- [ ] Criar o consumer no OpenSubtitles e o secret `OPENSUBTITLES_API_KEY` (Angelo)
- [ ] Teste real com chave válida: download, login certo e cota
- [ ] Testar o diálogo nativo "Escolher…" no macOS, Windows e Linux
- [ ] PR para `main` e tag `v1.0.0` (só com pedido)
- [ ] Apagar a branch remota `develop` e atualizar a descrição do repo no GitHub (decisão do Angelo)
