/*
Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.

SPDX-License-Identifier: Apache-2.0
*/

package chainmaker_sdk_go

import (
	"testing"

	"github.com/stretchr/testify/require"

	"chainmaker.org/chainmaker/pb-go/common"
)

func TestCheckNewBlockChainConfig(t *testing.T) {
	tests := []struct {
		name string
		res  *common.TxResponse
		err  error
	}{
		{
			"valid request",
			&common.TxResponse{Code: common.TxStatusCode_SUCCESS},
			nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cli, err := newMockChainClient(WithConfPath(sdkConfigForUtPath))
			require.Nil(t, err)
			defer cli.Stop()

			err = cli.CheckNewBlockChainConfig()
			require.Nil(t, err)
		})
	}
}
