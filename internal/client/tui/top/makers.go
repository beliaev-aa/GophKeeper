// Package top содержит функцию prepareMakers, которая инициализирует карту создателей экранов для TUI-приложения.
package top

import (
	"beliaev-aa/GophKeeper/internal/client/grpc"
	"beliaev-aa/GophKeeper/internal/client/tui"
	"beliaev-aa/GophKeeper/internal/client/tui/screens/auth"
	"beliaev-aa/GophKeeper/internal/client/tui/screens/blobs"
	"beliaev-aa/GophKeeper/internal/client/tui/screens/cards"
	"beliaev-aa/GophKeeper/internal/client/tui/screens/credentials"
	"beliaev-aa/GophKeeper/internal/client/tui/screens/remotes"
	"beliaev-aa/GophKeeper/internal/client/tui/screens/secrets"
	"beliaev-aa/GophKeeper/internal/client/tui/screens/storage"
	"beliaev-aa/GophKeeper/internal/client/tui/screens/texts"
)

func prepareMakers(client *grpc.ClientGRPC) map[tui.Screen]tui.ScreenMaker {

	return map[tui.Screen]tui.ScreenMaker{
		tui.BlobEditScreen:       &blobs.BlobEditScreen{},
		tui.CardEditScreen:       &cards.CardEditScreen{},
		tui.CredentialEditScreen: &credentials.CredentialEditScreen{},
		tui.FilePickScreen:       &blobs.FilePickScreen{},
		tui.LoginScreen:          &auth.AuthenticateScreen{},
		tui.RemoteOpenScreen:     &remotes.RemoteOpenScreenMaker{Client: client},
		tui.SecretTypeScreen:     &secrets.SecretTypeScreen{},
		tui.StorageBrowseScreen:  &storage.BrowseStorageScreen{},
		tui.TextEditScreen:       &texts.TextEditScreen{},
	}
}
