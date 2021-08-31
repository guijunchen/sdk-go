/*
Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.

SPDX-License-Identifier: Apache-2.0
*/

package chainmaker_sdk_go

import (
	"testing"

	"github.com/gogo/protobuf/proto"
	"github.com/stretchr/testify/require"

	"chainmaker.org/chainmaker/common/v2/crypto"
	bcx509 "chainmaker.org/chainmaker/common/v2/crypto/x509"
	"chainmaker.org/chainmaker/pb-go/v2/common"
	"chainmaker.org/chainmaker/sdk-go/v2/utils"
)

func TestSignPayloadWithPath(t *testing.T) {
	tests := []struct {
		name            string
		unsignedPayload *common.Payload
		wantErr         bool
	}{
		{
			"good",
			&common.Payload{
				ChainId: "chain1",
			},
			false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cli, err := newMockChainClient(nil, nil, WithConfPath(sdkConfigPathForUT))
			require.Nil(t, err)
			defer cli.Stop()

			e, err := SignPayloadWithPath(utils.Config.ChainClientConfig.UserSignKeyFilePath,
				utils.Config.ChainClientConfig.UserSignCrtFilePath, tt.unsignedPayload)
			require.Equal(t, tt.wantErr, err != nil)

			payloadBz, err := proto.Marshal(tt.unsignedPayload)
			require.Nil(t, err)

			var opts crypto.SignOpts
			hashalgo, err := bcx509.GetHashFromSignatureAlgorithm(cli.userCrt.SignatureAlgorithm)
			require.Nil(t, err)

			opts.Hash = hashalgo
			opts.UID = crypto.CRYPTO_DEFAULT_UID

			verified, err := cli.userCrt.PublicKey.VerifyWithOpts(payloadBz, e.Signature, &opts)
			require.Nil(t, err)
			require.True(t, verified)
		})
	}
}
