package models

import (
	"testing"
)

func TestSecret_ToClipboard(t *testing.T) {
	tests := []struct {
		name     string
		secret   *Secret
		expected string
	}{
		{
			name: "Credential_Secret",
			secret: func() *Secret {
				s := NewSecret(CredSecret)
				s.Creds = &Credentials{
					Login:    "user1",
					Password: "pass123",
				}
				return s
			}(),
			expected: "login: user1\npassword: pass123",
		},
		{
			name: "Card_Secret",
			secret: func() *Secret {
				s := NewSecret(CardSecret)
				s.Card = &Card{
					Number:   "1234567890123456",
					ExpMonth: 12,
					ExpYear:  2030,
					CVV:      123,
				}
				return s
			}(),
			expected: "Card Number: 1234567890123456\nExp: 12/2030CVV: 123",
		},
		{
			name: "Text_Secret",
			secret: func() *Secret {
				s := NewSecret(TextSecret)
				s.Text = &Text{
					Content: "Hello World",
				}
				return s
			}(),
			expected: "Text: Hello World\n",
		},
		{
			name: "Blob_Secret",
			secret: func() *Secret {
				s := NewSecret(BlobSecret)
				return s
			}(),
			expected: "File data cannot be moved to clipboard\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.secret.ToClipboard(); got != tt.expected {
				t.Errorf("Secret.ToClipboard() = %v, want %v", got, tt.expected)
			}
		})
	}
}
