package models

import (
	"bytes"
	"encoding/json"
	"testing"
)

var testKey = bytes.Repeat([]byte{0xAA}, 32)

func TestEncryptContent(t *testing.T) {
	type testCase struct {
		name    string
		stored  *StoredData
		key     []byte
		wantErr bool
	}

	testCases := []testCase{
		{
			name: "login_password_valid",
			stored: &StoredData{
				ID:          "lp1",
				Type:        LoginPasswordType,
				Description: "test login/password",
				MetaInfo:    []MetaData{},
				Content: LoginPassword{
					Login:    "myLogin",
					Password: "myPassword",
				},
			},
			key:     testKey,
			wantErr: false,
		},
		{
			name: "text_data_valid",
			stored: &StoredData{
				ID:          "txt1",
				Type:        TextDataType,
				Description: "test text",
				Content:     TextData{Text: "Hello, world!"},
			},
			key:     testKey,
			wantErr: false,
		},
		{
			name: "binary_data_valid",
			stored: &StoredData{
				ID:          "bin1",
				Type:        BinaryDataType,
				Description: "test binary",
				Content:     BinaryData{Data: []byte{0xDE, 0xAD, 0xBE, 0xEF}},
			},
			key:     testKey,
			wantErr: false,
		},
		{
			name: "card_data_valid",
			stored: &StoredData{
				ID:          "card1",
				Type:        CardDataType,
				Description: "test card",
				Content: CardData{
					CardNumber:     "1111 2222 3333 4444",
					ExpiryDate:     "10/28",
					CVV:            "123",
					CardHolderName: "John Doe",
				},
			},
			key:     testKey,
			wantErr: false,
		},
		{
			name: "otp_valid",
			stored: &StoredData{
				ID:          "otp1",
				Type:        OTPType,
				Description: "test otp",
				Content: OneTimePassword{
					OTPCode: "123456",
					Expiry:  "2030-12-31",
				},
			},
			key:     testKey,
			wantErr: false,
		},
		{
			name: "nil_content",
			stored: &StoredData{
				ID:      "nil1",
				Type:    LoginPasswordType,
				Content: nil,
			},
			key:     testKey,
			wantErr: true,
		},
		{
			name: "unsupported_data_type",
			stored: &StoredData{
				ID:          "unknown_type",
				Type:        DataType("some_random_type"),
				Description: "unsupported",
				Content:     "whatever",
			},
			key:     testKey,
			wantErr: true,
		},
		{
			name: "json_marshal_error",
			stored: &StoredData{
				ID:          "json_err",
				Type:        LoginPasswordType,
				Description: "marshal error test",
				Content: struct {
					Ch chan int
				}{
					Ch: make(chan int),
				},
			},
			key:     testKey,
			wantErr: true,
		},
		{
			name: "invalid_key_length",
			stored: &StoredData{
				ID:          "key_err",
				Type:        TextDataType,
				Description: "invalid key test",
				Content:     TextData{Text: "test"},
			},
			key:     bytes.Repeat([]byte{0xAA}, 10),
			wantErr: true,
		},
		{
			name: "another_unsupported_data_type",
			stored: &StoredData{
				ID:          "unknown_type_2",
				Type:        DataType("some_other_random_type"),
				Description: "unsupported_2",
				Content:     "whatever_2",
			},
			key:     testKey,
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.stored.EncryptContent(tc.key)
			if (err != nil) != tc.wantErr {
				t.Fatalf("EncryptContent() error = %v, wantErr = %v", err, tc.wantErr)
			}

			if !tc.wantErr {
				encrypted, ok := tc.stored.Content.(string)
				if !ok {
					t.Errorf("EncryptContent() = %v, want string in stored.Content", tc.stored.Content)
				} else if encrypted == "" {
					t.Error("EncryptContent() returned empty ciphertext, expected non-empty")
				}
			}
		})
	}
}

func TestDecryptContent(t *testing.T) {
	type testCase struct {
		name           string
		stored         *StoredData
		key            []byte
		wantErr        bool
		wantContentVal interface{}
	}

	lp := &StoredData{
		ID:          "lp_ok",
		Type:        LoginPasswordType,
		Description: "test decrypt lp",
		Content: LoginPassword{
			Login:    "myLogin",
			Password: "myPassword",
		},
	}
	if err := lp.EncryptContent(testKey); err != nil {
		panic("failed to prepare login/password encryption in test: " + err.Error())
	}

	txt := &StoredData{
		ID:          "txt_ok",
		Type:        TextDataType,
		Description: "test decrypt text",
		Content:     TextData{Text: "Hello, world!"},
	}
	if err := txt.EncryptContent(testKey); err != nil {
		panic("failed to prepare text encryption in test: " + err.Error())
	}

	binOk := &StoredData{
		ID:          "bin_ok",
		Type:        BinaryDataType,
		Description: "test decrypt binary",
		Content:     BinaryData{Data: []byte{0xAB, 0xCD, 0xEF}},
	}
	if err := binOk.EncryptContent(testKey); err != nil {
		panic("failed to prepare binary encryption in test: " + err.Error())
	}

	cardOk := &StoredData{
		ID:          "card_ok",
		Type:        CardDataType,
		Description: "test decrypt card",
		Content: CardData{
			CardNumber:     "5555 6666 7777 8888",
			ExpiryDate:     "05/30",
			CVV:            "999",
			CardHolderName: "Jane Doe",
		},
	}
	if err := cardOk.EncryptContent(testKey); err != nil {
		panic("failed to prepare card encryption in test: " + err.Error())
	}

	otpOk := &StoredData{
		ID:          "otp_ok",
		Type:        OTPType,
		Description: "test decrypt otp",
		Content: OneTimePassword{
			OTPCode: "654321",
			Expiry:  "2040-01-01",
		},
	}
	if err := otpOk.EncryptContent(testKey); err != nil {
		panic("failed to prepare OTP encryption in test: " + err.Error())
	}

	testCases := []testCase{
		{
			name:    "login_password_ok",
			stored:  lp,
			key:     testKey,
			wantErr: false,
			wantContentVal: LoginPassword{
				Login:    "myLogin",
				Password: "myPassword",
			},
		},
		{
			name:    "text_data_ok",
			stored:  txt,
			key:     testKey,
			wantErr: false,
			wantContentVal: TextData{
				Text: "Hello, world!",
			},
		},
		{
			name:    "binary_data_ok",
			stored:  binOk,
			key:     testKey,
			wantErr: false,
			wantContentVal: BinaryData{
				Data: []byte{0xAB, 0xCD, 0xEF},
			},
		},
		{
			name:    "card_data_ok",
			stored:  cardOk,
			key:     testKey,
			wantErr: false,
			wantContentVal: CardData{
				CardNumber:     "5555 6666 7777 8888",
				ExpiryDate:     "05/30",
				CVV:            "999",
				CardHolderName: "Jane Doe",
			},
		},
		{
			name:    "otp_data_ok",
			stored:  otpOk,
			key:     testKey,
			wantErr: false,
			wantContentVal: OneTimePassword{
				OTPCode: "654321",
				Expiry:  "2040-01-01",
			},
		},
		{
			name: "content_is_not_string",
			stored: &StoredData{
				ID:          "not_string",
				Type:        TextDataType,
				Description: "should fail",
				Content:     12345,
			},
			key:     testKey,
			wantErr: true,
		},
		{
			name: "unsupported_data_type_during_decrypt",
			stored: &StoredData{
				ID:      "unsupported_type",
				Type:    DataType("weird"),
				Content: "some_encrypted_string",
			},
			key:     testKey,
			wantErr: true,
		},
		{
			name: "invalid_hex_string",
			stored: &StoredData{
				ID:      "invalid_hex",
				Type:    TextDataType,
				Content: "XYZ_not_hex_data_123",
			},
			key:     testKey,
			wantErr: true,
		},
		{
			name: "invalid_ciphertext",
			stored: &StoredData{
				ID:      "invalid_cipher",
				Type:    OTPType,
				Content: "ABCD1234",
			},
			key:     testKey,
			wantErr: true,
		},
		{
			name: "decrypt_with_wrong_key",
			stored: &StoredData{
				ID:      "wrong_key",
				Type:    TextDataType,
				Content: txt.Content,
			},
			key:     bytes.Repeat([]byte{0xBB}, 32),
			wantErr: true,
		},
		{
			name: "binary_data_invalid_json",
			stored: &StoredData{
				ID:      "bin_json_fail",
				Type:    BinaryDataType,
				Content: "7b2244617461223a",
			},
			key:     testKey,
			wantErr: true,
		},
		{
			name: "card_data_invalid_json",
			stored: &StoredData{
				ID:      "card_json_fail",
				Type:    CardDataType,
				Content: "7b22436172644e756d626572223a5b",
			},
			key:     testKey,
			wantErr: true,
		},
		{
			name: "otp_data_invalid_json",
			stored: &StoredData{
				ID:      "otp_json_fail",
				Type:    OTPType,
				Content: "7b224f5450436f6465223a313233",
			},
			key:     testKey,
			wantErr: true,
		},
		{
			name: "another_unsupported_type_during_decrypt",
			stored: &StoredData{
				ID:      "unsupported_type_2",
				Type:    DataType("completely_weird_type"),
				Content: "some_encrypted_string",
			},
			key:     testKey,
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.stored.DecryptContent(tc.key)
			if (err != nil) != tc.wantErr {
				t.Fatalf("DecryptContent() error = %v, wantErr = %v", err, tc.wantErr)
			}
			if !tc.wantErr {

				gotBytes, err := json.Marshal(tc.stored.Content)
				if err != nil {
					t.Fatalf("failed to marshal decrypted content: %v", err)
				}
				wantBytes, err := json.Marshal(tc.wantContentVal)
				if err != nil {
					t.Fatalf("failed to marshal wantContentVal: %v", err)
				}
				if !bytes.Equal(gotBytes, wantBytes) {
					t.Errorf("DecryptContent() got = %s, want = %s", gotBytes, wantBytes)
				}
			}
		})
	}
}
