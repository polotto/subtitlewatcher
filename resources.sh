#!/bin/bash

mkdir -p ./resources/tmp/locales
"$GOROOT"/bin/fyne bundle -package "locales" -prefix "Res" ./resources/locales > ./resources/tmp/locales/locales.go
mkdir -p ./resources/tmp/images
"$GOROOT"/bin/fyne bundle -package "images" -prefix "Res" ./resources/images > ./resources/tmp/images/images.go