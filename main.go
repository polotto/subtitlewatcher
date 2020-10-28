package main

import (
	"fyne.io/fyne"
	"fyne.io/fyne/app"
	"fyne.io/fyne/dialog"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/storage"
	"fyne.io/fyne/widget"
	"strings"
	maincontroller "subtitlewatcher/controllers"
	"subtitlewatcher/messenger"
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

	var msgs = messenger.ReadMessages()

	watchStr["disabled"] = msgs["watchEnabledButton"]
	watchStr["enabled"] = msgs["watchDisabledButton"]
	var watchStarted = false

	var appMain = app.New()
	var w = appMain.NewWindow(msgs["appTitle"])

	w.Resize(fyne.NewSize(600, 500))

	openFileBtn := widget.NewButton(msgs["downloadButton"], func() {
		openFileDialog(w, maincontroller.FileFormats, func(filePath string) {
			maincontroller.DownloadSubtitle(msgs["subtitleNotFoundError"], filePath, func() {
				dialog.ShowInformation(msgs["doneDownloadTitle"], msgs["doneDownloadMsg"]+filePath, w)
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
				progress := dialog.NewProgressInfinite(msgs["loadingWatcherTitle"], msgs["loadingWatcherMsg"], w)
				progress.Show()
				maincontroller.SubtitleWatcherStart(msgs["subtitleNotFoundError"], folderPath, func(folderPath string) {
					watchFolderBtn.Text = watchStr["enabled"]
					watchFolderBtn.Importance = widget.HighImportance
					watchFolderBtn.Refresh()
					progress.Hide()
					dialog.ShowInformation(msgs["doneWatcherTitle"], msgs["doneWatcherMsg"]+folderPath, w)
				}, func(err error) {
					progress.Hide()
					dialog.ShowError(err, w)
				})
			})
		} else {
			progress := dialog.NewProgressInfinite(msgs["loadingStopWatcherTitle"], msgs["loadingStopWatcherMsg"], w)
			progress.Show()
			maincontroller.SubtitleWatcherStop(func() {
				watchFolderBtn.Text = watchStr["disabled"]
				watchFolderBtn.Importance = widget.MediumImportance
				watchFolderBtn.Refresh()
				progress.Hide()
				dialog.ShowInformation(msgs["stopWatcherTitle"], msgs["stopWatcherMsg"], w)
			})
		}
	})

	hContainer := fyne.NewContainerWithLayout(layout.NewHBoxLayout(), layout.NewSpacer(), openFileBtn, watchFolderBtn, layout.NewSpacer())
	vContainer := fyne.NewContainerWithLayout(layout.NewVBoxLayout(), layout.NewSpacer(), hContainer, layout.NewSpacer())
	w.SetContent(vContainer)
	w.ShowAndRun()

	defer func() {
		if watchStarted {
			maincontroller.SubtitleWatcherStop(nil)
		}
	}()
}
