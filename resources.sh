#!/bin/bash

mkdir -p ./resources/tmp/locales
fyne bundle -package "locales" -prefix "Res" ./resources/locales > ./resources/tmp/locales/locales.go
mkdir -p ./resources/tmp/images
fyne bundle -package "images" -prefix "Res" ./resources/images > ./resources/tmp/images/images.go