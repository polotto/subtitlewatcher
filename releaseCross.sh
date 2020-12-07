#!/bin/bash

NAME_BIN="subtitlewatcher"
APP_ID="subtitlewatcher/polotto/com.github"

./resources.sh

DIST="./fyne-cross/dist"

NAME_WIN="$NAME_BIN""_windows"
NAME_LIN="$NAME_BIN""_linux"
NAME_DAR="$NAME_BIN""_darwin"

fyne-cross windows -arch=amd64 -icon "icon.png"  -output "$NAME_WIN".exe &&
fyne-cross linux -arch=amd64  -icon "icon.png" -output $NAME_LIN &&
fyne-cross darwin -arch=amd64  -icon "icon.png" -output $NAME_DAR -app-id "$APP_ID"

zip -r "$DIST/$NAME_DAR".zip "$DIST/darwin-amd64/$NAME_DAR".app
cp "$DIST/linux-amd64/subtitlewatcher_linux.tar.gz" "$DIST/subtitlewatcher_linux.tar.gz"
cp "$DIST/windows-amd64/subtitlewatcher_windows.exe.zip" "$DIST/subtitlewatcher_windows.zip"