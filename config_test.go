/*
Copyright (C) BABEC. All rights reserved.
Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.

SPDX-License-Identifier: Apache-2.0
*/

package chainmaker_sdk_go

import (
	"fmt"
	"testing"

	"github.com/hokaccha/go-prettyjson"
	"github.com/stretchr/testify/require"
)

func TestInitConfig(t *testing.T) {
	err := InitConfig("./testdata/sdk_config_for_ut.yml")
	require.Nil(t, err)

	json, err := prettyjson.Marshal(Config)
	require.Nil(t, err)

	fmt.Println(string(json))
}
