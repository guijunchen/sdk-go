/**
 * @Author: zghh
 * @Date:   2020-12-03 10:16:49
 **/
package chainmaker_sdk_go

import (
	"encoding/hex"
	"fmt"
	"testing"
	"time"

	"chainmaker.org/chainmaker-go/chainmaker-sdk-go/pb"
	"github.com/stretchr/testify/require"
)

func TestCertHash(t *testing.T) {
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

func TestCertManage(t *testing.T) {
	var (
		err error
		client1, admin1 *ChainClient

		// org2 client证书
		certs = []string{
			"-----BEGIN CERTIFICATE-----\nMIICiDCCAi6gAwIBAgIDCuSTMAoGCCqGSM49BAMCMIGKMQswCQYDVQQGEwJDTjEQ\nMA4GA1UECBMHQmVpamluZzEQMA4GA1UEBxMHQmVpamluZzEfMB0GA1UEChMWd3gt\nb3JnMi5jaGFpbm1ha2VyLm9yZzESMBAGA1UECxMJcm9vdC1jZXJ0MSIwIAYDVQQD\nExljYS53eC1vcmcyLmNoYWlubWFrZXIub3JnMB4XDTIwMTExNjA2NDYwNFoXDTI1\nMTExNTA2NDYwNFowgZAxCzAJBgNVBAYTAkNOMRAwDgYDVQQIEwdCZWlqaW5nMRAw\nDgYDVQQHEwdCZWlqaW5nMR8wHQYDVQQKExZ3eC1vcmcyLmNoYWlubWFrZXIub3Jn\nMQ8wDQYDVQQLEwZjbGllbnQxKzApBgNVBAMTImNsaWVudDEudGxzLnd4LW9yZzIu\nY2hhaW5tYWtlci5vcmcwWTATBgcqhkjOPQIBBggqhkjOPQMBBwNCAAQjsmPDqPjx\nikMpRPkmWH8RFgUXwpzwaoMF9OQY6sAty2U8Q6TPlafMbm/xBls//UPZpi5uhwTv\neunkar0HqfvRo3sweTAOBgNVHQ8BAf8EBAMCAaYwDwYDVR0lBAgwBgYEVR0lADAp\nBgNVHQ4EIgQgjqe9Y2WHp+WC/GfKlvwummg3xvKPi9hbDja0QVFKa/EwKwYDVR0j\nBCQwIoAgmZcrtYWpTzN56LDZdqiHah3fG5w0kPaLoEBtyC8GfaEwCgYIKoZIzj0E\nAwIDSAAwRQIgbz8Du0bvtlWVJfBFzUamyfY2OodQDGBbKnr/eFXNeIECIQDnnJs5\nAX2NCT42Be3et+jhwxshehNsYm3WOOdTq/y+yg==\n-----END CERTIFICATE-----\n",
		}
	)

	client1, err = createClientWithOrgId(orgId1)
	require.Nil(t, err)

	admin1, err = createAdmin(orgId1)
	require.Nil(t, err)

	fmt.Println("====================== 证书冻结 ======================")
	testCertManageFrozen(t, client1, admin1, certs)
	fmt.Println("====================== 已冻结，不可用 ======================")
	testCertIsAvailable(t, false)

	fmt.Println("====================== 证书解冻 ======================")
	testCertManageUnfrozen(t, client1, admin1, certs)
	fmt.Println("====================== 已解冻，可用 ======================")
	testCertIsAvailable(t, true)
}

func testCertIsAvailable(t *testing.T, isAvailable bool) {
	_, err := createClientWithOrgId(orgId2)
	if isAvailable {
		require.Nil(t, err)
	} else {
		require.NotNil(t, err)
	}
}

func testCertManageFrozen(t *testing.T, client1, admin1 *ChainClient, certs []string) {
	var (
		err error
		payload []byte
		signedPayloadBytes []byte
		resp *pb.TxResponse
		result string
	)

	payload, err = client1.CreateCertManageFrozenPayload(certs)
	require.Nil(t, err)

	signedPayloadBytes, err = admin1.SignCertManagePayload(payload)
	require.Nil(t, err)

	resp, err = client1.SendCertManageRequest(signedPayloadBytes, -1, true)
	require.Nil(t, err)

	result = string(resp.ContractResult.Result)

	fmt.Printf("resp: %+v, result:%+s\n", resp, result)
}

func testCertManageUnfrozen(t *testing.T, client1, admin1 *ChainClient, certs []string) {
	var (
		err error
		payload []byte
		signedPayloadBytes []byte
		resp *pb.TxResponse
		result string
	)

	payload, err = client1.CreateCertManageUnfrozenPayload(certs)
	require.Nil(t, err)

	signedPayloadBytes, err = admin1.SignCertManagePayload(payload)
	require.Nil(t, err)

	resp, err = client1.SendCertManageRequest(signedPayloadBytes, -1, true)
	require.Nil(t, err)

	result = string(resp.ContractResult.Result)

	fmt.Printf("resp: %+v, result:%+s\n", resp, result)
}

func testCertManageRevoke(t *testing.T) {

}

func testCertAdd(t *testing.T, client *ChainClient) string {
	resp, err := client.AddCert()
	require.Nil(t, err)
	return hex.EncodeToString(resp.ContractResult.Result)
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
