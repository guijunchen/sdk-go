/*
Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.

SPDX-License-Identifier: Apache-2.0
*/

package chainmaker_sdk_go

//import (
//	"testing"
//
//	"github.com/stretchr/testify/require"
//)
//
//func TestGetChainMakerServerVersion(t *testing.T) {
//	tests := []struct {
//		name    string
//		version string
//		err     error
//	}{
//		{
//			"valid request",
//			"1.0.0",
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
//			version, err := cli.GetChainMakerServerVersion()
//			require.Nil(t, err)
//			require.Equal(t, tt.version, version)
//		})
//	}
//}
