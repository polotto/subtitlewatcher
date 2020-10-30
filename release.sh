#!/bin/bash

# GO111MODULE=on go get github.com/lucor/fyne-cross/v2/cmd/fyne-cross

TAG="v0.1.0"
NAME_BIN="subtitlewatcher"

./resources.sh

fyne-cross windows -arch=amd64 -icon "icon.png"  -output ""$NAME_BIN"_windows_"$TAG".exe" &&
fyne-cross linux -arch=amd64  -icon "icon.png" -output ""$NAME_BIN"_linux_"$TAG"" &&
fyne-cross darwin -arch=amd64  -icon "icon.png" -output ""$NAME_BIN"_darwin_"$TAG""  -app-id "io.subtitlewatcher"

MAC_DIST="./fyne-cross/dist/darwin-amd64/"

zip -r ""$MAC_DIST""$NAME_BIN"_darwin_"$TAG".zip" ""$MAC_DIST""$NAME_BIN"_darwin_"$TAG".app"