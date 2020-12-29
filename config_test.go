/**
 * @Author: jasonruan
 * @Date:   2020-12-29 11:26:01
 **/
package chainmaker_sdk_go

import (
	"fmt"
	"github.com/hokaccha/go-prettyjson"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestInitConfig(t *testing.T) {
	err := InitConfig("/home/jason/Work/ChainMaker/chainmaker-sdk-go/config_demo.yml")
	require.Nil(t, err)

	json, err := prettyjson.Marshal(Config)
	require.Nil(t, err)

	fmt.Println(string(json))
}
