package top

import (
	"beliaev-aa/GophKeeper/internal/client/tui"
	"beliaev-aa/GophKeeper/internal/client/tui/screens/auth"
	"beliaev-aa/GophKeeper/internal/client/tui/screens/blobs"
	"beliaev-aa/GophKeeper/internal/client/tui/screens/cards"
	"beliaev-aa/GophKeeper/internal/client/tui/screens/credentials"
	"beliaev-aa/GophKeeper/internal/client/tui/screens/remotes"
	"beliaev-aa/GophKeeper/internal/client/tui/screens/secrets"
	"beliaev-aa/GophKeeper/internal/client/tui/screens/storage"
	"beliaev-aa/GophKeeper/internal/client/tui/screens/texts"
	"beliaev-aa/GophKeeper/tests/mocks"
	"github.com/golang/mock/gomock"
	"reflect"
	"testing"
)

func Test_prepareMakers(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mocks.NewMockClientGRPCInterface(ctrl)

	type testCase struct {
		name          string
		screen        tui.Screen
		expectedMaker interface{}
	}

	testCases := []testCase{
		{name: "BlobEditScreen", screen: tui.BlobEditScreen, expectedMaker: &blobs.BlobEditScreen{}},
		{name: "CardEditScreen", screen: tui.CardEditScreen, expectedMaker: &cards.CardEditScreen{}},
		{name: "CredentialEditScreen", screen: tui.CredentialEditScreen, expectedMaker: &credentials.CredentialEditScreen{}},
		{name: "FilePickScreen", screen: tui.FilePickScreen, expectedMaker: &blobs.FilePickScreen{}},
		{name: "LoginScreen", screen: tui.LoginScreen, expectedMaker: &auth.AuthenticateScreen{}},
		{name: "RemoteOpenScreen", screen: tui.RemoteOpenScreen, expectedMaker: &remotes.RemoteOpenScreenMaker{Client: mockClient}},
		{name: "SecretTypeScreen", screen: tui.SecretTypeScreen, expectedMaker: &secrets.SecretTypeScreen{}},
		{name: "StorageBrowseScreen", screen: tui.StorageBrowseScreen, expectedMaker: &storage.BrowseStorageScreen{}},
		{name: "TextEditScreen", screen: tui.TextEditScreen, expectedMaker: &texts.TextEditScreen{}},
	}

	makers := prepareMakers(mockClient)

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			maker, exists := makers[tc.screen]
			if !exists {
				t.Fatalf("Maker for screen %v not found", tc.screen)
			}
			if reflect.TypeOf(maker) != reflect.TypeOf(tc.expectedMaker) {
				t.Errorf("Expected maker of type %T, got %T", tc.expectedMaker, maker)
			}
		})
	}
}
