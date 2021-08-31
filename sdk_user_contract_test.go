/*
Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.

SPDX-License-Identifier: Apache-2.0
*/

package chainmaker_sdk_go

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"

	"chainmaker.org/chainmaker/pb-go/v2/common"
)

func TestSendTxRequest(t *testing.T) {
	var txID = "b374f23e4e6747e4b5fcb3ca975ef1655ad56555adfd4534ae8676cd9f1eb145"

	tests := []struct {
		name         string
		cliTxReq     *common.TxRequest
		serverTxResp *common.TxResponse
		serverTxErr  error
		wantErr      bool
	}{
		{
			"bad",
			&common.TxRequest{Payload: &common.Payload{TxId: txID}},
			&common.TxResponse{
				Code: common.TxStatusCode_CONTRACT_FAIL,
			},
			errors.New("rpc server throw an error"),
			true,
		},
		{
			"good",
			&common.TxRequest{Payload: &common.Payload{TxId: txID}},
			&common.TxResponse{
				Code: common.TxStatusCode_SUCCESS,
			},
			nil,
			false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cli, err := newMockChainClient(tt.serverTxResp, tt.serverTxErr, WithConfPath(sdkConfigPathForUT))
			require.Nil(t, err)
			defer cli.Stop()

			txResp, err := cli.SendTxRequest(tt.cliTxReq, -1, false)
			require.Equal(t, tt.wantErr, err != nil)
			if err != nil {
				require.Contains(t, txResp.Message, tt.serverTxErr.Error())
			} else {
				require.Equal(t, tt.serverTxResp.Code, txResp.Code)
			}
		})
	}
}
