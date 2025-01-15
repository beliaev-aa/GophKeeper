// Package remotes содержит компоненты для управления удаленным доступом к данным в TUI.
package remotes

import (
	"beliaev-aa/GophKeeper/internal/client/grpc"
	"beliaev-aa/GophKeeper/internal/client/storage"
	"beliaev-aa/GophKeeper/internal/client/tui"
	"github.com/charmbracelet/bubbletea"
)

// RemoteOpenScreen представляет экран для открытия удаленного хранилища.
type RemoteOpenScreen struct {
	client *grpc.ClientGRPC
}

// RemoteOpenScreenMaker структура для создания экрана RemoteOpenScreen.
type RemoteOpenScreenMaker struct {
	Client *grpc.ClientGRPC
}

// Make создает новый экран RemoteOpenScreen.
func (m RemoteOpenScreenMaker) Make(_ tui.NavigationMsg, _, _ int) (tui.TeaLike, error) {
	return NewRemoteOpenScreen(m.Client), nil
}

// Make создает новый экран RemoteOpenScreen, используя переданное сообщение навигации.
func (s *RemoteOpenScreen) Make(msg tui.NavigationMsg, _, _ int) (tui.TeaLike, error) {
	return NewRemoteOpenScreen(msg.Client), nil
}

// NewRemoteOpenScreen создает и инициализирует новый экран удаленного открытия.
func NewRemoteOpenScreen(client *grpc.ClientGRPC) *RemoteOpenScreen {
	return &RemoteOpenScreen{
		client: client,
	}
}

// Init инициализирует экран и определяет начальное поведение в зависимости от того,
// аутентифицирован ли пользователь.
func (s *RemoteOpenScreen) Init() tea.Cmd {
	var commands []tea.Cmd

	if len(s.client.GetToken()) > 0 {
		store, err := storage.NewRemoteStorage(s.client)
		if err != nil {
			commands = append(commands, tui.ReportError(err))
		} else {
			commands = append(commands, tui.SetBodyPane(tui.StorageBrowseScreen, tui.WithStorage(store)))
		}
	} else {
		commands = append(commands, tui.SetBodyPane(tui.LoginScreen, tui.WithClient(s.client)))
	}

	return tea.Batch(commands...)
}

// Update обрабатывает сообщения и обновляет состояние экрана.
func (s *RemoteOpenScreen) Update(_ tea.Msg) tea.Cmd {
	return nil
}

// View возвращает визуальное представление экрана. Поскольку экран не имеет визуального контента,
// функция возвращает пустую строку.
func (s *RemoteOpenScreen) View() string {
	return ""
}
