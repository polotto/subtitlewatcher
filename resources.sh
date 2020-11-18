#!/bin/bash

mkdir -p ./resources/tmp/locales
"$GOPATH"/bin/fyne bundle -package "locales" -prefix "Res" ./resources/locales > ./resources/tmp/locales/locales.go
mkdir -p ./resources/tmp/images
"$GOPATH"/bin/fyne bundle -package "images" -prefix "Res" ./resources/images > ./resources/tmp/images/images.go