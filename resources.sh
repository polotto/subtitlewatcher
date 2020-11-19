#!/bin/bash

mkdir -p ./resources/tmp/locales
"$(go env GOPATH)"/bin/fyne bundle -package "locales" -prefix "Res" ./resources/locales > ./resources/tmp/locales/locales.go
mkdir -p ./resources/tmp/images
"$(go env GOPATH)"/bin/fyne bundle -package "images" -prefix "Res" ./resources/images > ./resources/tmp/images/images.go