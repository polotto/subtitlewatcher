# Publicar uma versão

## Uma vez só: chave da API do app

Os binários oficiais trazem uma chave da API do OpenSubtitles embutida, para o usuário não precisar criar uma.

1. Entrar em [opensubtitles.com/consumers](https://www.opensubtitles.com/consumers) com a conta do mantenedor e criar um consumer chamado `subtitlewatcher`.
2. Gravar a chave como secret do repositório:

   ```sh
   gh secret set OPENSUBTITLES_API_KEY --repo polotto/subtitlewatcher
   ```

A chave entra só nos builds de tag (`-X main.defaultAPIKey`). Ela não fica no código nem nos artefatos de PR, mas quem baixa o binário consegue extraí-la; é assim em qualquer app open source que usa essa API. Se ela for abusada, gerar outra no mesmo painel, trocar o secret e publicar uma versão nova.

Sem o secret, o CI publica normalmente e o app pede a chave em Configurações, Avançado.

## Cada versão

1. Mergear o PR em `main` com o CI verde.
2. Criar e enviar a tag:

   ```sh
   git switch main && git pull
   git tag v1.2.3
   git push origin v1.2.3
   ```

3. O workflow `CI` testa em Linux, macOS e Windows, roda `scripts/dist.sh 1.2.3` num runner macOS e cria a release com os pacotes e o `checksums.txt`.
4. Conferir a release: baixar o pacote do seu sistema, rodar `subtitlewatcher version`, abrir a janela e baixar uma legenda.

## Pacotes gerados

| Pacote | Conteúdo |
| --- | --- |
| `macos.zip` | binário universal (Intel e Apple Silicon), via `lipo` |
| `windows-x64.zip`, `windows-arm64.zip` | `subtitlewatcher.exe` |
| `linux-x64`, `linux-arm64`, `linux-armv7` `.tar.gz` | binário Linux; `armv7` serve para Raspberry Pi |
| `freebsd-x64.tar.gz` | binário FreeBSD |

Todos levam `README.md` e `LICENSE`. Para gerar localmente: `OPENSUBTITLES_API_KEY=... scripts/dist.sh 1.2.3`. Fora do macOS, sem `lipo`, saem `macos-arm64` e `macos-x64` separados.

## Limites conhecidos

- Os binários não são assinados nem notarizados. O README explica como liberar no macOS e no Windows. Assinar exige conta paga da Apple e certificado de code signing para Windows.
- O `User-Agent` enviado à API é `subtitlewatcher v<versão>`, como o OpenSubtitles pede.
