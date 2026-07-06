package gui

import (
	"fmt"
	"image/color"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"

	"github.com/Serge-Nook/topor/internal/core"
	"github.com/Serge-Nook/topor/internal/handlers"
	"github.com/Serge-Nook/topor/internal/utils"
)

const (
	appVersion = "1.0.0"
	appAuthor  = "Горшков Сергей Владимирович"
	appWebsite = "https://nookbat.ru/"
	donateURL  = "https://nookbat.ru/donate"
)

// ToporApp holds the main application state.
type ToporApp struct {
	app    fyne.App
	window fyne.Window
	files  []string

	// UI elements
	fileTable      *widget.Table
	fileMeta       []core.FileMetadata
	logText        *widget.Entry
	progressBar    *widget.ProgressBar
	statusLabel    *widget.Label
	btnApply       *widget.Button
	btnCancel      *widget.Button
	chkRecursive   *widget.Check
	chkBackup      *widget.Check
	maskEntry      *widget.Entry

	// Time settings
	creationMode   *widget.RadioGroup
	creationDate   *widget.Entry
	creationDays   *widget.Entry
	creationHours  *widget.Entry
	creationMins   *widget.Entry
	modifyMode     *widget.RadioGroup
	modifyDate     *widget.Entry
	modifyDays     *widget.Entry
	modifyHours    *widget.Entry
	modifyMins     *widget.Entry

	// Author settings
	authorAction   *widget.RadioGroup
	authorEntry    *widget.Entry
	lastModByEntry *widget.Entry

	// Processing state
	cancelMu   sync.Mutex
	cancelled  bool
	processing bool
}

// NewToporApp creates and runs the application.
func NewToporApp(app fyne.App) {
	ta := &ToporApp{app: app}
	ta.window = app.NewWindow(fmt.Sprintf("Топор — Массовый редактор метаданных v%s", appVersion))
	ta.window.Resize(fyne.NewSize(1100, 750))
	ta.window.SetContent(ta.buildUI())
	ta.updateStatus()
	ta.window.ShowAndRun()
}

func (ta *ToporApp) buildUI() fyne.CanvasObject {
	// Donate button (top right)
	donateBtn := widget.NewButton("₽", func() {
		utils.OpenURL(donateURL)
	})
	donateBtn.Importance = widget.HighImportance
	donateRow := container.NewHBox(layout.NewSpacer(), donateBtn)

	// File selection buttons
	btnAddFiles := widget.NewButton("Добавить файлы...", ta.onAddFiles)
	btnAddFolder := widget.NewButton("Добавить папку...", ta.onAddFolder)
	btnClear := widget.NewButton("Очистить список", ta.onClear)

	fileButtons := container.NewHBox(btnAddFiles, btnAddFolder, layout.NewSpacer(), btnClear)

	// Options
	ta.chkRecursive = widget.NewCheck("Включая подкаталоги", nil)
	ta.maskEntry = widget.NewEntry()
	ta.maskEntry.SetPlaceHolder("Например: *.docx;*.xlsx")
	ta.chkBackup = widget.NewCheck("Создавать .bak копии", nil)

	maskLabel := widget.NewLabel("Фильтр по маске:")
	optionsRow := container.NewHBox(ta.chkRecursive, maskLabel, ta.maskEntry, ta.chkBackup)

	fileSection := container.NewVBox(fileButtons, optionsRow)
	fileGroup := widget.NewCard("Входные данные", "", fileSection)

	// File table
	ta.fileMeta = nil
	ta.fileTable = ta.buildFileTable()
	tableContainer := container.NewStack(ta.fileTable)

	// Tabs
	timeTab := ta.buildTimePanel()
	authorTab := ta.buildAuthorPanel()
	logTab := ta.buildLogPanel()
	aboutTab := ta.buildAboutPanel()

	tabs := container.NewAppTabs(
		container.NewTabItem("Временные штампы", timeTab),
		container.NewTabItem("Метаданные автора", authorTab),
		container.NewTabItem("Журнал операций", logTab),
		container.NewTabItem("Об авторе", aboutTab),
	)

	// Progress and action buttons
	ta.progressBar = widget.NewProgressBar()
	ta.btnApply = widget.NewButton("Применить", ta.onApply)
	ta.btnApply.Importance = widget.HighImportance
	ta.btnCancel = widget.NewButton("Отмена", ta.onCancel)
	ta.btnCancel.Importance = widget.DangerImportance
	ta.btnCancel.Disable()

	actionRow := container.NewHBox(ta.progressBar, ta.btnApply, ta.btnCancel)

	// Status bar
	ta.statusLabel = widget.NewLabel("")
	statusBar := container.NewHBox(ta.statusLabel)

	// Use a border layout with the table taking most space
	mainSplit := container.NewVSplit(
		container.NewVBox(donateRow, fileGroup, tableContainer),
		container.NewVBox(tabs, actionRow, statusBar),
	)
	mainSplit.SetOffset(0.45)

	return mainSplit
}

func (ta *ToporApp) buildFileTable() *widget.Table {
	headers := []string{"Файл", "Путь", "Размер", "Формат", "Дата создания", "Дата изменения", "Автор", "Последний редактор"}

	table := widget.NewTable(
		func() (int, int) {
			return len(ta.fileMeta) + 1, len(headers)
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("placeholder text here")
		},
		func(id widget.TableCellID, cell fyne.CanvasObject) {
			label := cell.(*widget.Label)
			label.Truncation = fyne.TextTruncateEllipsis
			if id.Row == 0 {
				label.TextStyle = fyne.TextStyle{Bold: true}
				label.SetText(headers[id.Col])
				return
			}
			label.TextStyle = fyne.TextStyle{}
			idx := id.Row - 1
			if idx >= len(ta.fileMeta) {
				label.SetText("")
				return
			}
			m := ta.fileMeta[idx]
			switch id.Col {
			case 0:
				label.SetText(m.Name)
			case 1:
				label.SetText(m.Path)
			case 2:
				label.SetText(formatSize(m.Size))
			case 3:
				label.SetText(m.Format)
			case 4:
				if !m.CreationTime.IsZero() {
					label.SetText(m.CreationTime.Format("02.01.2006 15:04:05"))
				} else {
					label.SetText("")
				}
			case 5:
				if !m.ModifyTime.IsZero() {
					label.SetText(m.ModifyTime.Format("02.01.2006 15:04:05"))
				} else {
					label.SetText("")
				}
			case 6:
				label.SetText(m.Author)
			case 7:
				label.SetText(m.LastModBy)
			}
		},
	)

	table.SetColumnWidth(0, 150)
	table.SetColumnWidth(1, 200)
	table.SetColumnWidth(2, 80)
	table.SetColumnWidth(3, 60)
	table.SetColumnWidth(4, 140)
	table.SetColumnWidth(5, 140)
	table.SetColumnWidth(6, 120)
	table.SetColumnWidth(7, 120)

	return table
}

func (ta *ToporApp) buildTimePanel() fyne.CanvasObject {
	// Creation time settings
	ta.creationMode = widget.NewRadioGroup([]string{
		"Не изменять",
		"Установить абсолютную дату",
		"Сдвиг относительно текущей",
		"Приравнять к дате изменения",
	}, nil)
	ta.creationMode.SetSelected("Не изменять")

	ta.creationDate = widget.NewEntry()
	ta.creationDate.SetPlaceHolder("ДД.ММ.ГГГГ ЧЧ:ММ:СС")
	ta.creationDays = widget.NewEntry()
	ta.creationDays.SetPlaceHolder("0")
	ta.creationHours = widget.NewEntry()
	ta.creationHours.SetPlaceHolder("0")
	ta.creationMins = widget.NewEntry()
	ta.creationMins.SetPlaceHolder("0")

	creationAbsRow := container.NewHBox(widget.NewLabel("Дата:"), ta.creationDate)
	creationOffRow := container.NewHBox(
		widget.NewLabel("Дни:"), ta.creationDays,
		widget.NewLabel("Часы:"), ta.creationHours,
		widget.NewLabel("Мин:"), ta.creationMins,
	)

	creationGroup := widget.NewCard("Дата создания (Creation Time)", "", container.NewVBox(
		ta.creationMode, creationAbsRow, creationOffRow,
	))

	// Modify time settings
	ta.modifyMode = widget.NewRadioGroup([]string{
		"Не изменять",
		"Установить абсолютную дату",
		"Сдвиг относительно текущей",
		"Приравнять к дате создания",
	}, nil)
	ta.modifyMode.SetSelected("Не изменять")

	ta.modifyDate = widget.NewEntry()
	ta.modifyDate.SetPlaceHolder("ДД.ММ.ГГГГ ЧЧ:ММ:СС")
	ta.modifyDays = widget.NewEntry()
	ta.modifyDays.SetPlaceHolder("0")
	ta.modifyHours = widget.NewEntry()
	ta.modifyHours.SetPlaceHolder("0")
	ta.modifyMins = widget.NewEntry()
	ta.modifyMins.SetPlaceHolder("0")

	modifyAbsRow := container.NewHBox(widget.NewLabel("Дата:"), ta.modifyDate)
	modifyOffRow := container.NewHBox(
		widget.NewLabel("Дни:"), ta.modifyDays,
		widget.NewLabel("Часы:"), ta.modifyHours,
		widget.NewLabel("Мин:"), ta.modifyMins,
	)

	modifyGroup := widget.NewCard("Дата изменения (Modification Time)", "", container.NewVBox(
		ta.modifyMode, modifyAbsRow, modifyOffRow,
	))

	return container.NewVBox(creationGroup, modifyGroup)
}

func (ta *ToporApp) buildAuthorPanel() fyne.CanvasObject {
	ta.authorAction = widget.NewRadioGroup([]string{
		"Не изменять",
		"Удалить",
		"Заменить на:",
	}, nil)
	ta.authorAction.SetSelected("Не изменять")

	ta.authorEntry = widget.NewEntry()
	ta.authorEntry.SetPlaceHolder("Новый автор")
	ta.lastModByEntry = widget.NewEntry()
	ta.lastModByEntry.SetPlaceHolder("Последний редактор")

	return widget.NewCard("Метаданные автора", "", container.NewVBox(
		ta.authorAction,
		container.NewHBox(widget.NewLabel("Автор:"), ta.authorEntry),
		container.NewHBox(widget.NewLabel("Последний редактор:"), ta.lastModByEntry),
	))
}

func (ta *ToporApp) buildLogPanel() fyne.CanvasObject {
	ta.logText = widget.NewMultiLineEntry()
	ta.logText.Disable()
	ta.logText.SetMinRowsVisible(8)
	return ta.logText
}

func (ta *ToporApp) buildAboutPanel() fyne.CanvasObject {
	titleText := canvas.NewText(fmt.Sprintf("Топор v%s", appVersion), color.NRGBA{R: 160, G: 255, B: 240, A: 255})
	titleText.TextSize = 22
	titleText.TextStyle = fyne.TextStyle{Bold: true}

	descText := canvas.NewText("Массовый редактор метаданных и временных штампов файлов", color.NRGBA{R: 128, G: 229, B: 208, A: 255})
	descText.TextSize = 15

	authorText := canvas.NewText(fmt.Sprintf("Автор: %s", appAuthor), color.NRGBA{R: 128, G: 229, B: 208, A: 255})
	authorText.TextSize = 14

	siteBtn := widget.NewButton("Сайт автора: nookbat.ru", func() {
		utils.OpenURL(appWebsite)
	})
	siteBtn.Importance = widget.LowImportance

	donateBtn := widget.NewButton("Сделать пожертвование (nookbat.ru/donate)", func() {
		utils.OpenURL(donateURL)
	})
	donateBtn.Importance = widget.LowImportance

	licenseText := canvas.NewText("Лицензия: MIT", color.NRGBA{R: 92, G: 184, B: 165, A: 255})
	licenseText.TextSize = 13

	return container.NewVBox(
		titleText, descText,
		widget.NewSeparator(),
		authorText, siteBtn, donateBtn,
		widget.NewSeparator(),
		licenseText,
		layout.NewSpacer(),
	)
}

// ─── File Selection ─────────────────────────────────────

func (ta *ToporApp) onAddFiles() {
	fd := dialog.NewFileOpen(func(reader fyne.URIReadCloser, err error) {
		if err != nil || reader == nil {
			return
		}
		path := reader.URI().Path()
		reader.Close()
		ta.addPaths([]string{path})
	}, ta.window)
	fd.Show()
}

func (ta *ToporApp) onAddFolder() {
	fd := dialog.NewFolderOpen(func(uri fyne.ListableURI, err error) {
		if err != nil || uri == nil {
			return
		}
		ta.addPaths([]string{uri.Path()})
	}, ta.window)
	fd.Show()
}

func (ta *ToporApp) addPaths(paths []string) {
	mask := ta.maskEntry.Text
	filterExts := core.ParseMask(mask)

	newFiles := core.DiscoverFiles(paths, ta.chkRecursive.Checked, filterExts)

	existing := make(map[string]bool)
	for _, f := range ta.files {
		existing[f] = true
	}

	added := 0
	for _, f := range newFiles {
		if !existing[f] {
			ta.files = append(ta.files, f)
			existing[f] = true
			added++
		}
	}

	ta.logAppend(fmt.Sprintf("Добавлено файлов: %d (всего: %d)", added, len(ta.files)))
	ta.refreshTable()
}

func (ta *ToporApp) onClear() {
	ta.files = nil
	ta.fileMeta = nil
	ta.fileTable.Refresh()
	ta.progressBar.SetValue(0)
	ta.updateStatus()
	ta.logAppend("Список файлов очищен.")
}

func (ta *ToporApp) refreshTable() {
	ta.fileMeta = nil
	for _, path := range ta.files {
		meta := readFileMetadata(path)
		ta.fileMeta = append(ta.fileMeta, meta)
	}
	ta.fileTable.Refresh()
	ta.updateStatus()
}

// ─── Processing ─────────────────────────────────────────

func (ta *ToporApp) onApply() {
	if len(ta.files) == 0 {
		dialog.ShowInformation("Нет файлов", "Добавьте файлы для обработки.", ta.window)
		return
	}

	if ta.processing {
		return
	}

	settings := ta.getSettings()
	ta.processing = true
	ta.cancelMu.Lock()
	ta.cancelled = false
	ta.cancelMu.Unlock()

	ta.btnApply.Disable()
	ta.btnCancel.Enable()
	ta.progressBar.SetValue(0)
	ta.logAppend(fmt.Sprintf("Начало обработки: %d файлов...", len(ta.files)))

	go ta.processFiles(settings)
}

func (ta *ToporApp) onCancel() {
	ta.cancelMu.Lock()
	ta.cancelled = true
	ta.cancelMu.Unlock()
	ta.logAppend("Отмена обработки...")
}

func (ta *ToporApp) processFiles(settings core.ProcessingSettings) {
	total := len(ta.files)
	for i, path := range ta.files {
		ta.cancelMu.Lock()
		if ta.cancelled {
			ta.cancelMu.Unlock()
			ta.logAppend("Обработка отменена.")
			break
		}
		ta.cancelMu.Unlock()

		result := ta.processFile(path, settings)
		progress := float64(i+1) / float64(total)
		ta.progressBar.SetValue(progress)

		switch result.Status {
		case core.StatusSuccess:
			ta.logAppend(fmt.Sprintf("✓ %s", filepath.Base(path)))
		case core.StatusError:
			ta.logAppend(fmt.Sprintf("✗ %s: %s", filepath.Base(path), result.Message))
		case core.StatusSkipped:
			ta.logAppend(fmt.Sprintf("⊘ %s: %s", filepath.Base(path), result.Message))
		}
	}

	ta.processing = false
	ta.btnApply.Enable()
	ta.btnCancel.Disable()
	ta.logAppend("Обработка завершена.")
	ta.refreshTable()
}

func (ta *ToporApp) processFile(path string, settings core.ProcessingSettings) core.ProcessingResult {
	result := core.ProcessingResult{Path: path}

	// Check access
	if _, err := os.Stat(path); err != nil {
		result.Status = core.StatusSkipped
		result.Message = "файл недоступен"
		return result
	}

	// Create backup if needed
	if settings.CreateBackup {
		if err := utils.CreateBackup(path); err != nil {
			result.Status = core.StatusError
			result.Message = fmt.Sprintf("ошибка бэкапа: %v", err)
			return result
		}
	}

	handler := handlers.GetHandler(path)
	if handler == nil {
		result.Status = core.StatusSkipped
		result.Message = "формат не поддерживается"
		return result
	}

	// Apply metadata changes
	var authorSettings *core.AuthorSettings
	if settings.Author.Action != core.AuthorActionNone {
		authorSettings = &settings.Author
	}
	var timeSettings *core.TimeSettings
	if settings.Time.CreationMode != core.DateModeNone || settings.Time.ModifyMode != core.DateModeNone {
		timeSettings = &settings.Time
	}

	if authorSettings != nil || timeSettings != nil {
		if err := handler.WriteMetadata(path, authorSettings, timeSettings); err != nil {
			result.Status = core.StatusError
			result.Message = err.Error()
			return result
		}
	}

	// Apply filesystem timestamps
	if timeSettings != nil {
		if err := core.ApplyTimestamps(path, settings.Time); err != nil {
			result.Status = core.StatusError
			result.Message = fmt.Sprintf("ошибка установки времени: %v", err)
			return result
		}
	}

	result.Status = core.StatusSuccess
	return result
}

func (ta *ToporApp) getSettings() core.ProcessingSettings {
	settings := core.ProcessingSettings{
		CreateBackup: ta.chkBackup.Checked,
	}

	// Time settings
	settings.Time.CreationMode = parseTimeMode(ta.creationMode.Selected, false)
	settings.Time.ModifyMode = parseTimeMode(ta.modifyMode.Selected, true)

	if settings.Time.CreationMode == core.DateModeAbsolute {
		settings.Time.CreationAbsolute = parseDateTime(ta.creationDate.Text)
	}
	if settings.Time.CreationMode == core.DateModeOffset {
		settings.Time.CreationOffsetD = parseInt(ta.creationDays.Text)
		settings.Time.CreationOffsetH = parseInt(ta.creationHours.Text)
		settings.Time.CreationOffsetM = parseInt(ta.creationMins.Text)
	}

	if settings.Time.ModifyMode == core.DateModeAbsolute {
		settings.Time.ModifyAbsolute = parseDateTime(ta.modifyDate.Text)
	}
	if settings.Time.ModifyMode == core.DateModeOffset {
		settings.Time.ModifyOffsetD = parseInt(ta.modifyDays.Text)
		settings.Time.ModifyOffsetH = parseInt(ta.modifyHours.Text)
		settings.Time.ModifyOffsetM = parseInt(ta.modifyMins.Text)
	}

	// Author settings
	switch ta.authorAction.Selected {
	case "Удалить":
		settings.Author.Action = core.AuthorActionDelete
	case "Заменить на:":
		settings.Author.Action = core.AuthorActionReplace
		settings.Author.NewAuthor = ta.authorEntry.Text
		settings.Author.NewLastModBy = ta.lastModByEntry.Text
	}

	return settings
}

// ─── Helpers ────────────────────────────────────────────

func (ta *ToporApp) logAppend(msg string) {
	ts := time.Now().Format("15:04:05")
	line := fmt.Sprintf("[%s] %s\n", ts, msg)
	ta.logText.Enable()
	ta.logText.SetText(ta.logText.Text + line)
	ta.logText.Disable()
}

func (ta *ToporApp) updateStatus() {
	ta.statusLabel.SetText(fmt.Sprintf(
		"Файлов: %d  |  Автор: Горшков С.В.  |  %s",
		len(ta.files), appWebsite,
	))
}

func readFileMetadata(path string) core.FileMetadata {
	fi, err := os.Stat(path)
	meta := core.FileMetadata{
		Path:   path,
		Name:   filepath.Base(path),
		Format: strings.ToUpper(strings.TrimPrefix(filepath.Ext(path), ".")),
	}
	if err != nil {
		return meta
	}
	meta.Size = fi.Size()
	meta.ModifyTime = fi.ModTime()
	meta.CreationTime = fi.ModTime()

	handler := handlers.GetHandler(path)
	if handler != nil {
		info, err := handler.ReadMetadata(path)
		if err == nil && info != nil {
			if info.Author != "" {
				meta.Author = info.Author
			}
			if info.LastModBy != "" {
				meta.LastModBy = info.LastModBy
			}
			if !info.CreationTime.IsZero() {
				meta.CreationTime = info.CreationTime
			}
			if !info.ModifyTime.IsZero() {
				meta.ModifyTime = info.ModifyTime
			}
		}
	}
	return meta
}

func formatSize(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d Б", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	suffix := []string{"КБ", "МБ", "ГБ", "ТБ"}
	return fmt.Sprintf("%.1f %s", float64(bytes)/float64(div), suffix[exp])
}

func parseTimeMode(selected string, isModify bool) core.DateMode {
	switch selected {
	case "Установить абсолютную дату":
		return core.DateModeAbsolute
	case "Сдвиг относительно текущей":
		return core.DateModeOffset
	case "Приравнять к дате изменения":
		return core.DateModeEqualCreateToModify
	case "Приравнять к дате создания":
		return core.DateModeEqualModifyToCreate
	default:
		return core.DateModeNone
	}
}

func parseDateTime(s string) time.Time {
	s = strings.TrimSpace(s)
	formats := []string{
		"02.01.2006 15:04:05",
		"02.01.2006 15:04",
		"2006-01-02 15:04:05",
		"2006-01-02T15:04:05",
	}
	for _, f := range formats {
		if t, err := time.Parse(f, s); err == nil {
			return t
		}
	}
	return time.Now()
}

func parseInt(s string) int {
	s = strings.TrimSpace(s)
	if s == "" {
		return 0
	}
	var v int
	fmt.Sscanf(s, "%d", &v)
	return v
}
