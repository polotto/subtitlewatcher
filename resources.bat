@echo off

if not exist ".\resources\tmp\locales" mkdir .\resources\tmp\locales
fyne bundle -package locales -prefix Res .\resources\locales > .\resources\tmp\locales\locales.go
if not exist ".\resources\tmp\images" mkdir .\resources\tmp\images
fyne bundle -package images -prefix Res .\resources\images > .\resources\tmp\images\images.go