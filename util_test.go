package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHexToInt(t *testing.T) {
	tests := map[string]struct {
		hexStr  string
		wantNum int
		wantErr bool
	}{
		"convert hex to int successfully": {
			hexStr:  "0x0f1",
			wantNum: 241,
			wantErr: false,
		},
		"string contain special character": {
			hexStr:  "0\n123",
			wantNum: 0,
			wantErr: true,
		},
		"string without 0x prefix": {
			hexStr:  "0f1",
			wantNum: 241,
			wantErr: false,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			num, err := HexToInt(tt.hexStr)
			assert.Equal(t, tt.wantNum, num)
			assert.Equal(t, tt.wantErr, err != nil)
		})
	}
}
