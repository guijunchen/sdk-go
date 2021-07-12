/*
Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.

SPDX-License-Identifier: Apache-2.0
*/

package chainmaker_sdk_go

//func TestSendArchiveBlockRequest(t *testing.T) {
//	tests := []struct {
//		name   string
//		height int64
//		res    *common.TxResponse
//		err    error
//	}{
//		{
//			"valid request",
//			100,
//			&common.TxResponse{Code: common.TxStatusCode_SUCCESS},
//			nil,
//		},
//	}
//
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			cli, err := newMockChainClient(WithConfPath(sdkConfigPathForUT))
//			require.Nil(t, err)
//			defer cli.Stop()
//
//			var (
//				payload            []byte
//				signedPayloadBytes []byte
//				resp               *common.TxResponse
//			)
//
//			payload, err = cli.CreateArchiveBlockPayload(tt.height)
//			require.Nil(t, err)
//
//			signedPayloadBytes = payload
//
//			resp, err = cli.SendArchiveBlockRequest(signedPayloadBytes, -1)
//			require.Nil(t, err)
//
//			if resp.Code != tt.res.Code {
//				t.Error("error: expected", tt.res, "received", resp)
//			}
//		})
//	}
//}
//
//func TestSendRestoreBlockRequest(t *testing.T) {
//	tests := []struct {
//		name      string
//		fullblock []byte
//		res       *common.TxResponse
//		err       error
//	}{
//		{
//			"valid request",
//			[]byte("fullblock"),
//			&common.TxResponse{Code: common.TxStatusCode_SUCCESS},
//			nil,
//		},
//	}
//
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			cli, err := newMockChainClient(WithConfPath(sdkConfigPathForUT))
//			require.Nil(t, err)
//			defer cli.Stop()
//
//			var (
//				payload            []byte
//				signedPayloadBytes []byte
//				resp               *common.TxResponse
//			)
//
//			payload, err = cli.CreateRestoreBlockPayload(tt.fullblock)
//			require.Nil(t, err)
//
//			signedPayloadBytes = payload
//
//			resp, err = cli.SendRestoreBlockRequest(signedPayloadBytes, -1)
//			require.Nil(t, err)
//
//			if resp.Code != tt.res.Code {
//				t.Error("error: expected", tt.res, "received", resp)
//			}
//		})
//	}
//}
