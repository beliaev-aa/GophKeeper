package blobs

import (
	"beliaev-aa/GophKeeper/internal/client/storage"
	"beliaev-aa/GophKeeper/internal/client/tui"
	"beliaev-aa/GophKeeper/internal/client/tui/screens"
	"beliaev-aa/GophKeeper/internal/client/tui/styles"
	"beliaev-aa/GophKeeper/pkg/models"
	"fmt"
	"github.com/charmbracelet/bubbles/filepicker"
	"github.com/charmbracelet/bubbletea"
	"os"
	"path/filepath"
	"strings"
)

type FilePickScreen struct {
	secret     *models.Secret
	storage    storage.Storage
	filePicker filepicker.Model
	callback   tui.NavigationCallback
}

func (s *FilePickScreen) Make(msg tui.NavigationMsg, _, _ int) (tui.TeaLike, error) {
	return NewFilePickScreen(msg.Secret, msg.Storage, msg.Callback), nil
}

func NewFilePickScreen(secret *models.Secret, store storage.Storage, callback tui.NavigationCallback) *FilePickScreen {
	defaultPath, err := os.UserHomeDir()
	if err != nil {
		panic("Error getting working directory: %v\n")
	}

	fp := filepicker.New()
	fp.CurrentDirectory = filepath.Join(defaultPath)
	fp.AutoHeight = false
	fp.Height = 10

	m := FilePickScreen{
		filePicker: fp,
		secret:     secret,
		storage:    store,
		callback:   callback,
	}

	return &m
}

func (s *FilePickScreen) Init() tea.Cmd {
	return s.filePicker.Init()
}

func (s *FilePickScreen) Update(msg tea.Msg) tea.Cmd {
	var (
		cmd      tea.Cmd
		commands []tea.Cmd
	)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "b":
			commands = append(commands, tui.SetBodyPane(tui.BlobEditScreen, tui.WithStorage(s.storage), tui.WithSecret(s.secret)))
		}
	case tea.WindowSizeMsg:
		s.filePicker.Height = msg.Height - styles.FilePickerBotPadding
	}

	s.filePicker, cmd = s.filePicker.Update(msg)
	commands = append(commands, cmd)

	if selected, path := s.filePicker.DidSelectFile(msg); selected {
		commands = append(commands, tui.ReportInfo("selected: %v", path))
		commands = append(commands, s.callback(path))
	}

	return tea.Batch(commands...)
}

func (s *FilePickScreen) View() string {

	var b strings.Builder
	b.WriteString(fmt.Sprintf("%20s%s:\n", "", s.filePicker.CurrentDirectory))
	b.WriteString(s.filePicker.View())

	return screens.RenderContent("Select file to store. Use ←, ↑, →, ↓ to navigate. Press B to go back", b.String())
}
