/**
 * @Author: zghh
 * @Date:   2020-12-03 10:16:49
 **/
package chainmaker_sdk_go

import (
	"testing"
	"time"

	"chainmaker.org/chainmaker-go/chainmaker-sdk-go/pb"
	"github.com/stretchr/testify/require"
)

func TestCertManageGo(t *testing.T) {
	client, err := createClient()
	require.Nil(t, err)

	certHash := testCertAdd(t, client)
	time.Sleep(3 * time.Second)

	certInfos := testQueryCert(t, client, []string{certHash})
	require.Equal(t, 1, len(certInfos.CertInfos))

	testDeleteCert(t, client, []string{certHash})
	time.Sleep(3 * time.Second)

	var bytesNil []byte
	bytesNil = nil
	certInfos = testQueryCert(t, client, []string{certHash})
	require.Equal(t, 1, len(certInfos.CertInfos))
	require.Equal(t, bytesNil, certInfos.CertInfos[0].Cert)
}

func testCertAdd(t *testing.T, client *ChainClient) string {
	resp, err := client.AddCert()
	require.Nil(t, err)
	return string(resp.ContractResult.Result)
}

func testQueryCert(t *testing.T, client *ChainClient, certHashes []string) *pb.CertInfos {
	certInfos, err := client.QueryCert(certHashes)
	require.Nil(t, err)
	return certInfos
}

func testDeleteCert(t *testing.T, client *ChainClient, certHashes []string) {
	_, err := client.DeleteCert(certHashes)
	require.Nil(t, err)
}
