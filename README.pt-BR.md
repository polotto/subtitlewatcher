# subtitlewatcher

[English](README.md)

Baixa legendas do [OpenSubtitles](https://www.opensubtitles.com) para os seus vídeos. Escolha um vídeo e pegue a legenda, ou aponte para uma pasta e todo vídeo novo que aparecer nela ganha legenda sozinho.

![Tela principal](docs/images/main.jpg)

- Janela simples que abre no navegador, em português ou inglês.
- Procura a legenda primeiro pela impressão digital do arquivo, então ela bate com a sua versão do vídeo.
- Dois idiomas em ordem de preferência, por exemplo português do Brasil e depois inglês.
- Linha de comando para scripts e servidores.
- Um arquivo pequeno, sem instalação, para macOS, Windows, Linux e FreeBSD.

## Download

Baixe a última versão na [página de releases](https://github.com/polotto/subtitlewatcher/releases/latest):

| Sistema | Arquivo |
| --- | --- |
| macOS (Intel e Apple Silicon) | `subtitlewatcher-<versão>-macos.zip` |
| Windows | `subtitlewatcher-<versão>-windows-x64.zip` (ou `-arm64` para notebooks ARM) |
| Linux | `subtitlewatcher-<versão>-linux-x64.tar.gz` (também `arm64` e `armv7` para Raspberry Pi) |
| FreeBSD | `subtitlewatcher-<versão>-freebsd-x64.tar.gz` |

Cada release tem um `checksums.txt` com o SHA-256 de cada arquivo.

## Primeira vez

**macOS.** Descompacte e dê dois cliques em `subtitlewatcher`. O app não é assinado pela Apple, então na primeira vez o macOS bloqueia. Abra **Ajustes do Sistema > Privacidade e Segurança**, clique em **Abrir Mesmo Assim** ao lado da mensagem do subtitlewatcher e confirme. Uma janela do Terminal abre junto com o app: deixe ela aberta enquanto usa.

**Windows.** Descompacte e dê dois cliques em `subtitlewatcher.exe`. Se o SmartScreen mostrar "O Windows protegeu o computador", clique em **Mais informações** e depois em **Executar assim mesmo**. Uma janela de console abre junto com o app: deixe ela aberta enquanto usa.

**Linux e FreeBSD.** Extraia e rode `./subtitlewatcher` no terminal.

O app abre uma página no navegador padrão. Se não abrir, copie o endereço que aparece na janela do terminal.

## Como usar

- **Legenda para um vídeo**: clique em **Escolher…**, selecione o vídeo e clique em **Baixar legenda**. A legenda é salva ao lado do vídeo com o mesmo nome, e os players carregam sozinhos.
- **Monitorar uma pasta**: escolha a pasta (a de Downloads, por exemplo) e clique em **Começar a monitorar**. Todo vídeo que aparecer nela ou nas subpastas ganha legenda assim que termina de copiar. Vídeos que já têm `.srt` são ignorados. O monitoramento só funciona com o app aberto.
- **Atividade** mostra o que aconteceu, inclusive quantos downloads ainda restam hoje.
- **Sair** fecha o app. Fechar a janela do terminal faz o mesmo.

## Conta do OpenSubtitles (opcional)

O app funciona sem conta, com limite de **5 legendas por dia**. Para baixar mais:

1. Crie uma conta gratuita em [opensubtitles.com](https://www.opensubtitles.com/en/users/sign_up).
2. No app, abra **Configurações**, preencha **Usuário** e **Senha** em "Conta do OpenSubtitles" e clique em **Salvar**.
3. O app entra na hora e mostra seu nível e quantos downloads restam hoje. Se a senha estiver errada, ele avisa e não tenta de novo até você salvar outros dados.

![Configurações](docs/images/settings.jpg)

A seção **Avançado** tem um campo de chave da API. Com os downloads oficiais você não precisa dele, porque eles já trazem uma chave. Preencha só para usar sua própria chave de [opensubtitles.com/consumers](https://www.opensubtitles.com/consumers) ou se você compilou o app.

## Linha de comando

```sh
subtitlewatcher get filme.mkv                   # um ou mais vídeos
subtitlewatcher -lang en,es get *.mp4           # outros idiomas
subtitlewatcher watch ~/Downloads               # monitora até Ctrl+C
subtitlewatcher -every 30s watch /media/videos  # varre com menos frequência
subtitlewatcher version
```

A linha de comando usa as mesmas configurações salvas pela janela. Estas variáveis de ambiente têm prioridade:

| Variável | Uso |
| --- | --- |
| `OPENSUBTITLES_USERNAME`, `OPENSUBTITLES_PASSWORD` | Login da conta |
| `OPENSUBTITLES_API_KEY` | Sua própria chave da API |

## Problemas comuns

| Mensagem | O que fazer |
| --- | --- |
| Nenhuma legenda encontrada nos idiomas escolhidos | Tente um segundo idioma nas Configurações. Alguns vídeos ainda não têm legenda no OpenSubtitles. |
| Limite diário de downloads atingido | Espere até amanhã ou entre com uma conta. |
| Usuário ou senha errados | Corrija nas Configurações e salve. |
| O OpenSubtitles recusou a chave da API | Apague a chave em Configurações, Avançado, para usar a que vem no app, ou crie uma nova. |
| Este sistema não tem janela de escolha | No Linux instale `zenity` ou `kdialog`, ou cole o caminho no campo. |

## Privacidade

- As configurações ficam em `config.json` dentro de `~/Library/Application Support/subtitlewatcher` (macOS), `%AppData%\subtitlewatcher` (Windows) ou `~/.config/subtitlewatcher` (Linux). A senha fica lá em texto puro, legível só pelo seu usuário.
- O app só conversa com `api.opensubtitles.com`. A página roda no seu próprio computador (`127.0.0.1`) e outros sites não conseguem controlá-la.

## Compilar

```sh
go build .                                   # precisa do Go 1.24
OPENSUBTITLES_API_KEY=... go run .           # sem chave embutida, passe a sua
scripts/dist.sh 1.2.3                        # pacotes de todas as plataformas em dist/
```

Para publicar uma versão: [docs/RELEASING.md](docs/RELEASING.md).

## Licença

[MIT](LICENSE)
