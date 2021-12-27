package utils

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestInitConfig(t *testing.T) {
	tests := []struct {
		name     string
		confPath string
	}{
		{
			"good",
			"../testdata/sdk_config.yml",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := InitConfig(tt.confPath)
			require.Nil(t, err)
		})
	}
}
