NAME_BIN="subtitlewatcher"
APP_ID="subtitlewatcher/polotto/com.github"

DIST="./fyne-cross"
DIST_WIN="$DIST/windows/"
DIST_LIN="$DIST/linux/"
DIST_DAR="$DIST/darwin/"

mkdir -p $DIST_WIN
mkdir -p $DIST_LIN
mkdir -p $DIST_DAR

NAME_WIN="$DIST_WIN""$NAME_BIN""_windows"
NAME_LIN="$DIST_LIN""$NAME_BIN""_linux"
NAME_DAR="$DIST_LIN""$NAME_BIN""_darwin"

fyne package -os windows -icon "icon.png" -release -name "$NAME_WIN".exe &&
fyne package -os linux -icon "icon.png" -release -name $NAME_LIN &&
fyne package -os darwin -icon "icon.png" -release -name $NAME_DAR -appID "$APP_ID"

zip -r "$NAME_WIN".zip "$NAME_WIN".exe
zip -r "$NAME_DAR".zip "$NAME_DAR".app