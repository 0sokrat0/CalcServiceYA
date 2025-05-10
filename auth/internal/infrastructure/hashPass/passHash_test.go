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
		// TODO: Add test cases.
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
	type args struct {
		password string
	}
	tests := []struct {
		name    string
		h       *Hasher
		args    args
		want    string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.h.Hash(tt.args.password)
			if (err != nil) != tt.wantErr {
				t.Errorf("Hasher.Hash() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Hasher.Hash() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHasher_Compare(t *testing.T) {
	type args struct {
		hash     string
		password string
	}
	tests := []struct {
		name string
		h    *Hasher
		args args
		want bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.h.Compare(tt.args.hash, tt.args.password); got != tt.want {
				t.Errorf("Hasher.Compare() = %v, want %v", got, tt.want)
			}
		})
	}
}
