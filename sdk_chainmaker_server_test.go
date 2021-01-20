package chainmaker_sdk_go

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestChainClient_GetChainMakerServerVersion(t *testing.T) {
	client, err := createClient()
	require.Nil(t, err)
	version, err := client.GetChainMakerServerVersion()
	require.Nil(t, err)
	fmt.Println("get chainmaker server version:", version)
}
