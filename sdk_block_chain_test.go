package chainmaker_sdk_go

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestChainClient_CheckNewBlockChainConfig(t *testing.T) {
	client, err := createClient()
	require.Nil(t, err)
	err = client.CheckNewBlockChainConfig()
	require.Nil(t, err)
	fmt.Println("check new block chain config: ok")
}
