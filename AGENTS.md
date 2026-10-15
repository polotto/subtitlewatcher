# subtitlewatcher

App em Go que baixa legendas do OpenSubtitles (API REST) para um vídeo ou para vídeos novos numa pasta. Sem argumentos abre uma UI web local; com `get`/`watch` é CLI.

## Comandos
- testes: `go vet ./... && go test ./...` (um teste: `go test -run TestFetch`)
- build: `CGO_ENABLED=0 go build -ldflags="-s -w" .`
- pacotes de todas as plataformas: `scripts/dist.sh <versão>` (gera `dist/`)
- release: tag `vX.Y.Z` em `main`; o `ci.yml` roda o `dist.sh` e publica. Passo a passo em `docs/RELEASING.md`.

## Arquitetura
- `main.go`: flags, config salva (`os.UserConfigDir()`), comandos `get` e `watch` (polling com `filepath.WalkDir`).
- `gui.go`: servidor em 127.0.0.1 com token por sessão e checagem de Host, diálogos nativos via `osascript`/PowerShell/`zenity`.
- `ui.html`: página única embutida com `go:embed`, textos em `strings.en` e `strings.pt`.
- `opensubtitles.go`: cliente da API, escolha da legenda e hash do arquivo.

## Regras do código
- **Só stdlib e sem CGO.** O objetivo do projeto é ser mínimo; dependência nova precisa de motivo forte.
- A chave da API do app entra só no build de tag, via `-X main.defaultAPIKey` e o secret `OPENSUBTITLES_API_KEY`. Nunca versionar chave.
- Erros que o usuário precisa entender viram código em `errorCode` (`gui.go`) e texto em `ui.html`.
- Mensagens da CLI em inglês. Texto novo na UI entra em `en` e `pt` no `ui.html`; erros conhecidos vão como código (`noKey`, `notFound`) e a página traduz.
- Toda rota `/api/*` passa pelo `guard`. Não abrir CORS nem escutar fora de 127.0.0.1.
- Docs para o usuário: `README.md` (inglês) e `README.pt-BR.md`, sempre os dois juntos. Screenshots em `docs/images/`.
