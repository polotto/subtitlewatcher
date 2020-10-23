// https://medium.com/@adiach3nko/package-management-with-go-modules-the-pragmatic-guide-c831b4eaaf31
// https://github.com/fyne-io/fyne
// https://github.com/matcornic/subify
// https://github.com/radovskyb/watcher

// https://github.com/fyne-io/fyne/issues/941
// https://github.com/fyne-io/fyne/blob/86d26ebe4d97a525aa5cf1b6720186fc76d3b669/cmd/fyne_demo/screens/window.go
package main

import (
	"fyne.io/fyne"
	"fyne.io/fyne/app"
	"fyne.io/fyne/dialog"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/storage"
	"fyne.io/fyne/widget"
	"strings"
)

func openFileDialog(w fyne.Window, fileFormats []string, chosen func(filePath string)) {
	fd := dialog.NewFileOpen(func(reader fyne.URIReadCloser, err error) {
		if err == nil && reader == nil {
			return
		}
		if err != nil {
			dialog.ShowError(err, w)
			return
		}
		filePath := strings.Replace(reader.URI().String(), "file://", "", 1)
		chosen(filePath)
	}, w)
	fd.SetFilter(storage.NewExtensionFileFilter(fileFormats))
	fd.Show()
}

func main() {
	var watchStr = map[string]string{}
	watchStr["disabled"] = "Start watch folder and download subtitle automatically"
	watchStr["enabled"] = "Stop watch folder and download subtitle automatically"

	var appMain = app.New()
	var w = appMain.NewWindow("Subtitle watcher")

	w.Resize(fyne.NewSize(600, 500))

	openFileBtn := widget.NewButton("Download subtitle for a file", func() {
		openFileDialog(w, FileFormats, func(filePath string) {
			downloadSubtitle(filePath, func() {
				dialog.ShowInformation("Done!", "Subtitle downloaded for the file: " + filePath, w)
			}, func(err error) {
				dialog.ShowError(err, w)
			})
		})
	})

	var watchFolderBtn *widget.Button
	watchFolderBtn = widget.NewButton(watchStr["disabled"], func() {
		progress := dialog.NewProgressInfinite("Folder watcher", "Starting folder watcher, please wait...", w)
		progress.Show()
		subtitleWatcher(func(home string, name string) {
			watchFolderBtn.Text = watchStr["enabled"]
			watchFolderBtn.Style = widget.PrimaryButton
			progress.Hide()
			dialog.ShowInformation("Folder watcher", "Watching files in the dir '"+home+"' of the user: "+name, w)
		}, func() {
			watchFolderBtn.Text = watchStr["disabled"]
			watchFolderBtn.Style = widget.DefaultButton
			progress.Hide()
			dialog.ShowInformation("Folder watcher", "Folder watcher stopped", w)
		}, func(err error) {
			progress.Hide()
			dialog.ShowError(err, w)
		})
	})

	//grid := fyne.NewContainerWithLayout(layout.NewGridWrapLayout(fyne.NewSize(250, 250)), openFileBtn, watchFolderBtn)
	hContainer := fyne.NewContainerWithLayout(layout.NewHBoxLayout(), layout.NewSpacer(), openFileBtn, watchFolderBtn, layout.NewSpacer())
	vContainer := fyne.NewContainerWithLayout(layout.NewVBoxLayout(), layout.NewSpacer(), hContainer, layout.NewSpacer())
	w.SetContent(vContainer)
	w.ShowAndRun()
}
