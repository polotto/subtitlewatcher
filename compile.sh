#!/bin/bash
# https://stackoverflow.com/questions/36915134/go-golang-cross-compile-from-mac-to-windows-fatal-error-windows-h-file-not-f
# env GOOS="windows" GOARCH="386"   CGO_ENABLED="1" CC="i686-w64-mingw32-gcc"
env GOOS="windows" GOARCH="amd64" CGO_ENABLED="1" CC="x86_64-w64-mingw32-gcc" fyne package -os windows -icon icon.png
