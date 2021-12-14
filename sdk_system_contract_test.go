/*
Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.

SPDX-License-Identifier: Apache-2.0
*/

package chainmaker_sdk_go

import (
	"testing"

	"github.com/gogo/protobuf/proto"
	"github.com/stretchr/testify/require"

	"chainmaker.org/chainmaker/pb-go/v2/common"
)

func TestGetTxByTxId(t *testing.T) {
	var txID = "b374f23e4e6747e4b5fcb3ca975ef1655ad56555adfd4534ae8676cd9f1eb145"

	goodTxInfoBz, err := proto.Marshal(&common.TransactionInfo{Transaction: &common.Transaction{Payload: &common.Payload{TxId: txID}}})
	require.Nil(t, err)

	tests := []struct {
		name         string
		serverTxResp *common.TxResponse
		serverTxErr  error
		wantErr      bool
	}{
		{
			"bad",
			&common.TxResponse{
				Code: common.TxStatusCode_SUCCESS,
				ContractResult: &common.ContractResult{
					Result: []byte("this is a bad *common.TransactionInfo bytes"),
				},
			},
			nil,
			true,
		},
		{
			"good",
			&common.TxResponse{
				Code: common.TxStatusCode_SUCCESS,
				ContractResult: &common.ContractResult{
					Result: goodTxInfoBz,
				},
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

			txInfo, err := cli.GetTxByTxId(txID)
			require.Equal(t, tt.wantErr, err != nil)

			if txInfo != nil {
				bz, err := proto.Marshal(txInfo)
				require.Nil(t, err)
				require.Equal(t, tt.serverTxResp.ContractResult.Result, bz)
			}
		})
	}
}
