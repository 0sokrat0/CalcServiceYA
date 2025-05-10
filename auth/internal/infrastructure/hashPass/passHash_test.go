package hashpass

import (
	"reflect"
	"testing"
)

func TestNewPassHasher(t *testing.T) {
	tests := []struct {
		name string
		want PasswordHasher
	}{
		{
			name: "should return Hasher instance",
			want: &Hasher{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewPassHasher(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewPassHasher() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHasher_Hash(t *testing.T) {
	h := &Hasher{}

	tests := []struct {
		name    string
		h       *Hasher
		args    string
		wantErr bool
	}{
		{
			name:    "hash valid password",
			h:       h,
			args:    "validpassword",
			wantErr: false,
		},
		{
			name:    "hash empty password",
			h:       h,
			args:    "",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.h.Hash(tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("Hasher.Hash() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if got == "" {
					t.Error("Hasher.Hash() returned empty hash")
				}
				if !h.Compare(got, tt.args) {
					t.Error("Hashed password does not match original")
				}
			}
		})
	}
}

func TestHasher_Compare(t *testing.T) {
	h := &Hasher{}
	correctPassword := "correctpassword"
	correctHash, err := h.Hash(correctPassword)
	if err != nil {
		t.Fatalf("Could not hash password: %v", err)
	}

	otherHash, err := h.Hash("otherpassword")
	if err != nil {
		t.Fatalf("Could not hash other password: %v", err)
	}

	tests := []struct {
		name     string
		hash     string
		password string
		want     bool
	}{
		{
			name:     "correct password",
			hash:     correctHash,
			password: correctPassword,
			want:     true,
		},
		{
			name:     "incorrect password",
			hash:     correctHash,
			password: "incorrect",
			want:     false,
		},
		{
			name:     "invalid hash format",
			hash:     "$invalidhash",
			password: correctPassword,
			want:     false,
		},
		{
			name:     "empty password",
			hash:     correctHash,
			password: "",
			want:     false,
		},
		{
			name:     "different hash for password",
			hash:     otherHash,
			password: correctPassword,
			want:     false,
		},
		{
			name:     "empty hash",
			hash:     "",
			password: "any",
			want:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := h.Compare(tt.hash, tt.password); got != tt.want {
				t.Errorf("Hasher.Compare() = %v, want %v", got, tt.want)
			}
		})
	}
}
