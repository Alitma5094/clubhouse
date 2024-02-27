package main

import (
	"errors"
	"testing"
)

func TestValidateMessage(t *testing.T) {
	var tests = []struct {
		name string
		text string
		want error
	}{
		{"Valid Message", "Hello, World", nil},
		{"Empty Text", "", ErrMessageEmptyText},
		{"Too Long", "Knowledge nay estimable questions repulsive daughters boy. Solicitude gay way unaffected expression for. His mistress ladyship required off horrible disposed rejoiced. Unpleasing pianoforte unreserved as oh he unpleasant no inquietude insipidity. Advantages can discretion possession add favourable cultivated admiration far. Why rather assure how esteem end hunted nearer and before. By an truth after heard going early given he. Charmed to it excited females whether at examine. Him abilities suffering may are yet dependent.", ErrMessageTooLong},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ans := ValidateMessageInput(tt.text)
			if !errors.Is(ans, tt.want) {
				t.Errorf("got %s, want %s", ans, tt.want)
			}
		})
	}
}
