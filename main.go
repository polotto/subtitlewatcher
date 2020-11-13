package main

import (
	"fyne.io/fyne"
	"fyne.io/fyne/app"
	"fyne.io/fyne/dialog"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/storage"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"
	"strings"
	maincontroller "subtitlewatcher/controllers"
	"subtitlewatcher/messenger"
	"subtitlewatcher/resources/tmp/images"
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
	appMain.Settings().SetTheme(theme.DarkTheme())

	res := fyne.NewStaticResource("icon", images.ResIconPng.StaticContent)
	appMain.SetIcon(res)

	var w = appMain.NewWindow(msgs["appTitleWindow"])

	if err := maincontroller.LoadSettings(); err != nil {
		dialog.ShowError(err, w)
	}

	w.Resize(fyne.NewSize(600, 500))

	lang1Label := widget.NewLabel(msgs["lang1Label"])
	lang1Select := widget.NewSelect(maincontroller.Languages(), func(s string) {
		maincontroller.Select(s, 0)
	})
	lang1Select.SetSelected(maincontroller.FindSelection(0))

	lang2Label := widget.NewLabel(msgs["lang2Label"])
	lang2Select := widget.NewSelect(maincontroller.Languages(), func(s string) {
		maincontroller.Select(s, 1)
	})
	lang2Select.SetSelected(maincontroller.FindSelection(1))

	actionsLabel := widget.NewLabel(msgs["langActions"])

	openFileBtn := widget.NewButton(msgs["downloadButton"], func() {
		openFileDialog(w, maincontroller.FileFormats, func(filePath string) {
			progress := dialog.NewProgressInfinite(msgs["downloadSubtitleTitleDialogProgress"], msgs["downloadSubtitleMsgDialogProgress"], w)
			progress.Show()
			maincontroller.DownloadSubtitle(msgs["downloadSubtitleNotFoundDialogError"], filePath, func() {
				progress.Hide()
				dialog.ShowInformation(msgs["downloadDoneTitleDialogInfo"], msgs["downloadDoneMsgDialogInfo"]+filePath, w)
			}, func(err error) {
				progress.Hide()
				dialog.ShowError(err, w)
			})
		})
	})

	var watchFolderBtn *widget.Button
	watchFolderBtn = widget.NewButton(watchStr["disabled"], func() {
		watchStarted = !watchStarted
		if watchStarted {
			openFolderDialog(w, func(folderPath string) {
				progress := dialog.NewProgressInfinite(msgs["watchLoadingTitleDialogProgress"], msgs["watchLoadingMsgDialogProgress"], w)
				progress.Show()
				maincontroller.SubtitleWatcherStart(msgs["downloadSubtitleNotFoundDialogError"], folderPath, func(folderPath string) {
					watchFolderBtn.Text = watchStr["enabled"]
					watchFolderBtn.Importance = widget.HighImportance
					watchFolderBtn.Refresh()
					progress.Hide()
					dialog.ShowInformation(msgs["watchDoneTitleDialogInfo"], msgs["watchDoneMsgDialogInfo"]+folderPath, w)
				}, func(err error) {
					progress.Hide()
					dialog.ShowError(err, w)
				})
			})
		} else {
			progress := dialog.NewProgressInfinite(msgs["watchLoadingStopTitleDialogProgress"], msgs["watchLoadingStopMsgDialogProgress"], w)
			progress.Show()
			maincontroller.SubtitleWatcherStop(func() {
				watchFolderBtn.Text = watchStr["disabled"]
				watchFolderBtn.Importance = widget.MediumImportance
				watchFolderBtn.Refresh()
				progress.Hide()
				dialog.ShowInformation(msgs["watchStopTitleDialogInfo"], msgs["watchStopMsgDialogInfo"], w)
			})
		}
	})

	vContainer := fyne.NewContainerWithLayout(layout.NewVBoxLayout(), lang1Label, lang1Select,
		lang2Label, lang2Select,
		layout.NewSpacer(), actionsLabel, openFileBtn, watchFolderBtn)
	w.SetContent(vContainer)
	w.ShowAndRun()

	defer func() {
		if watchStarted {
			maincontroller.SubtitleWatcherStop(nil)
		}
		if err := maincontroller.SaveSettings(); err != nil {
			dialog.ShowError(err, w)
		}
	}()
}
