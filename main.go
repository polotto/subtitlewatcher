// https://medium.com/@adiach3nko/package-management-with-go-modules-the-pragmatic-guide-c831b4eaaf31
// https://github.com/fyne-io/fyne
// https://github.com/matcornic/subify
// https://github.com/radovskyb/watcher

// https://github.com/fyne-io/fyne/issues/941
// https://github.com/fyne-io/fyne/blob/86d26ebe4d97a525aa5cf1b6720186fc76d3b669/cmd/fyne_demo/screens/window.go

// PR
// https://github.com/fyne-io/fyne/pull/1222
// https://github.com/fyne-io/fyne/pull/1449
// https://github.com/fyne-io/fyne/issues/21
// https://github.com/fyne-io/fyne/blob/master/cmd/fyne_demo/screens/widget.go#L100
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

func uriToPath(uri string) string {
	return strings.Replace(uri, "file://", "", 1)
}

func openFileDialog(w fyne.Window, fileFormats []string, chosen func(filePath string)) {
	fd := dialog.NewFileOpen(func(reader fyne.URIReadCloser, err error) {
		if err == nil && reader == nil {
			return
		}
		if err != nil {
			dialog.ShowError(err, w)
			return
		}
		filePath := uriToPath(reader.URI().String())
		chosen(filePath)
	}, w)
	fd.SetFilter(storage.NewExtensionFileFilter(fileFormats))
	fd.Show()
}

func openFolderDialog(w fyne.Window, chosen func(folderPath string)) {
	dialog.ShowFolderOpen(func(list fyne.ListableURI, err error) {
		if err != nil {
			dialog.ShowError(err, w)
			return
		}
		if list == nil {
			return
		}
		folderPath := uriToPath(list.String())
		chosen(folderPath)
	}, w)
}

func main() {
	var watchStr = map[string]string{}
	watchStr["disabled"] = "Start watch folder and download subtitle automatically"
	watchStr["enabled"] = "Stop watch folder and download subtitle automatically"
	var watchStarted = false

	var appMain = app.New()
	var w = appMain.NewWindow("Subtitle watcher")

	w.Resize(fyne.NewSize(600, 500))

	openFileBtn := widget.NewButton("Download subtitle for a file", func() {
		openFileDialog(w, FileFormats, func(filePath string) {
			downloadSubtitle(filePath, func() {
				dialog.ShowInformation("Done!", "Subtitle downloaded for the file: "+filePath, w)
			}, func(err error) {
				dialog.ShowError(err, w)
			})
		})
	})

	var watchFolderBtn *widget.Button
	watchFolderBtn = widget.NewButton(watchStr["disabled"], func() {
		watchStarted = !watchStarted
		if watchStarted {
			openFolderDialog(w, func(folderPath string) {
				progress := dialog.NewProgressInfinite("Folder watcher", "Starting folder watcher, please wait...", w)
				progress.Show()
				subtitleWatcherStart(folderPath, func(folderPath string) {
					watchFolderBtn.Text = watchStr["enabled"]
					watchFolderBtn.Importance = widget.HighImportance
					watchFolderBtn.Refresh()
					progress.Hide()
					dialog.ShowInformation("Folder watcher", "Watching files in the dir '"+folderPath, w)
				}, func(err error) {
					progress.Hide()
					dialog.ShowError(err, w)
				})
			})
		} else {
			progress := dialog.NewProgressInfinite("Folder watcher", "Starting folder watcher, please wait...", w)
			progress.Show()
			subtitleWatcherStop(func() {
				watchFolderBtn.Text = watchStr["disabled"]
				watchFolderBtn.Importance = widget.MediumImportance
				watchFolderBtn.Refresh()
				progress.Hide()
				dialog.ShowInformation("Folder watcher", "Folder watcher stopped", w)
			})
		}
	})

	//grid := fyne.NewContainerWithLayout(layout.NewGridWrapLayout(fyne.NewSize(250, 250)), openFileBtn, watchFolderBtn)
	hContainer := fyne.NewContainerWithLayout(layout.NewHBoxLayout(), layout.NewSpacer(), openFileBtn, watchFolderBtn, layout.NewSpacer())
	vContainer := fyne.NewContainerWithLayout(layout.NewVBoxLayout(), layout.NewSpacer(), hContainer, layout.NewSpacer())
	w.SetContent(vContainer)
	w.ShowAndRun()
}
