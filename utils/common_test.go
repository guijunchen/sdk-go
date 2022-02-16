package utils

import (
	"testing"

	"chainmaker.org/chainmaker/pb-go/v2/common"
	"github.com/stretchr/testify/require"
)

func TestGetRandTxId(t *testing.T) {
	txId := GetRandTxId()
	require.Len(t, txId, 64)
}

func TestCheckProposalRequestResp(t *testing.T) {
	tests := []struct {
		name               string
		serverTxResp       *common.TxResponse
		needContractResult bool
		wantErr            bool
	}{
		{
			"good",
			&common.TxResponse{Code: common.TxStatusCode_SUCCESS, ContractResult: &common.ContractResult{
				Code: SUCCESS,
			}},
			true,
			false,
		},
		{
			"bad",
			&common.TxResponse{Code: common.TxStatusCode_CONTRACT_FAIL},
			false,
			true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := CheckProposalRequestResp(tt.serverTxResp, tt.needContractResult)
			require.Equal(t, err != nil, tt.wantErr)
		})
	}
}
