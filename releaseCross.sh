#!/bin/bash

# GO111MODULE=on go get github.com/lucor/fyne-cross/v2/cmd/fyne-cross

NAME_BIN="subtitlewatcher"
APP_ID="subtitlewatcher/polotto/com.github"

./resources.sh

fyne-cross windows -arch=amd64 -icon "icon.png"  -output "$NAME_BIN"_windows.exe &&
fyne-cross linux -arch=amd64  -icon "icon.png" -output "$NAME_BIN"_linux &&
fyne-cross darwin -arch=amd64  -icon "icon.png" -output "$NAME_BIN"_darwin -app-id "$APP_ID"

MAC_DIST="./fyne-cross/dist/darwin-amd64/"

zip -r "$MAC_DIST""$NAME_BIN"_darwin.zip "$MAC_DIST""$NAME_BIN"_darwin.app