// Package credentials содержит компоненты и логику для создания и редактирования учетных данных в TUI приложении.
package credentials

import (
	"beliaev-aa/GophKeeper/internal/client/storage"
	"beliaev-aa/GophKeeper/internal/client/tui"
	"beliaev-aa/GophKeeper/internal/client/tui/components"
	"beliaev-aa/GophKeeper/internal/client/tui/screens"
	"beliaev-aa/GophKeeper/internal/client/tui/styles"
	"beliaev-aa/GophKeeper/pkg/models"
	"context"
	"errors"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbletea"
	"time"
)

const (
	credTitle = iota
	credMetadata
	credLogin
	credPassword
)

// CredentialEditScreen структура для экрана редактирования учетных данных.
type CredentialEditScreen struct {
	inputGroup components.InputGroup
	secret     *models.Secret
	storage    storage.Storage
}

type inputOpts struct {
	placeholder string
	charLimit   int
	focus       bool
}

// Make создаёт новый экран CredentialEditScreen на основе переданных параметров.
func (s *CredentialEditScreen) Make(msg tui.NavigationMsg, _, _ int) (tui.TeaLike, error) {
	return NewCredentialEditScreen(msg.Secret, msg.Storage), nil
}

// NewCredentialEditScreen создаёт и инициализирует новый экземпляр CredentialEditScreen.
func NewCredentialEditScreen(secret *models.Secret, store storage.Storage) *CredentialEditScreen {
	m := CredentialEditScreen{
		secret:  secret,
		storage: store,
	}

	inputs := make([]textinput.Model, 4)
	inputs[credTitle] = newInput(inputOpts{placeholder: "Title", charLimit: 64})
	inputs[credMetadata] = newInput(inputOpts{placeholder: "Metadata", charLimit: 64})
	inputs[credLogin] = newInput(inputOpts{placeholder: "Login", charLimit: 64})
	inputs[credPassword] = newInput(inputOpts{placeholder: "Password", charLimit: 64})

	var buttons []components.Button
	buttons = append(buttons, components.Button{Title: "[ Submit ]", Cmd: func() tea.Cmd {
		if err := m.Submit(); err != nil {
			return tui.ReportError(err)
		}
		return tui.SetBodyPane(tui.StorageBrowseScreen, tui.WithStorage(m.storage))
	}})

	buttons = append(buttons, components.Button{Title: "[ Back ]", Cmd: func() tea.Cmd {
		return tui.SetBodyPane(tui.StorageBrowseScreen, tui.WithStorage(m.storage))
	}})

	if secret.ID > 0 {
		inputs[credTitle].SetValue(secret.Title)
		inputs[credMetadata].SetValue(secret.Metadata)
		inputs[credLogin].SetValue(secret.Creds.Login)
		inputs[credPassword].SetValue(secret.Creds.Password)
	}

	m.inputGroup = components.NewInputGroup(inputs, buttons)

	return &m
}

// Init инициализирует компоненты экрана.
func (s *CredentialEditScreen) Init() tea.Cmd {
	return s.inputGroup.Init()
}

// Update обрабатывает пользовательский ввод и обновляет состояние экрана.
func (s *CredentialEditScreen) Update(msg tea.Msg) tea.Cmd {
	var (
		cmd      tea.Cmd
		commands []tea.Cmd
	)

	ig, cmd := s.inputGroup.Update(msg)
	s.inputGroup = ig.(components.InputGroup)

	commands = append(commands, cmd)

	return tea.Batch(commands...)
}

// Submit обрабатывает отправку данных учетной записи в хранилище.
func (s *CredentialEditScreen) Submit() error {
	var (
		err error
	)

	title := s.inputGroup.Inputs[credTitle].Value()
	metadata := s.inputGroup.Inputs[credMetadata].Value()
	login := s.inputGroup.Inputs[credLogin].Value()
	password := s.inputGroup.Inputs[credPassword].Value()

	if len(metadata) == 0 {
		return errors.New("please enter metadata")
	}

	if len(title) == 0 {
		return errors.New("please enter title")
	}

	if len(login) == 0 {
		return errors.New("please enter login")
	}

	if len(password) == 0 {
		return errors.New("please enter password")
	}

	s.secret.Title = title
	s.secret.Metadata = metadata
	s.secret.Creds = &models.Credentials{Login: login, Password: password}
	s.secret.UpdatedAt = time.Now()

	if s.secret.ID == 0 {
		s.secret.CreatedAt = time.Now()
		err = s.storage.Create(context.Background(), s.secret)
	} else {
		err = s.storage.Update(context.Background(), s.secret)
	}

	return err
}

// View отображает текущее состояние экрана в виде строки.
func (s *CredentialEditScreen) View() string {
	return screens.RenderContent("Fill in credential details:", s.inputGroup.View())
}

// newInput создаёт новую модель ввода текста с заданными параметрами.
func newInput(opts inputOpts) textinput.Model {
	t := textinput.New()
	t.CharLimit = opts.charLimit
	t.Placeholder = opts.placeholder

	if opts.focus {
		t.Focus()
		t.PromptStyle = styles.Focused
		t.TextStyle = styles.Focused
	}

	return t
}
