package main

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"math/big"
	"testing"
	"time"

	"chainmaker.org/chainmaker/common/crypto/paillier"

	cmlog "chainmaker.org/chainmaker/common/log"
	sdk "chainmaker.org/chainmaker/sdk-go"

	"chainmaker.org/chainmaker/pb-go/common"
	"chainmaker.org/chainmaker/sdk-go/examples"
	"github.com/stretchr/testify/require"
)

const (
	sdkConfigOrg1Client1Path = "../sdk_configs/sdk_config_org1_client1.yml"
	chainId                  = "chain1"
	orgId1                   = "wx-org1.chainmaker.org"
	orgId2                   = "wx-org2.chainmaker.org"
	orgId3                   = "wx-org3.chainmaker.org"
	orgId4                   = "wx-org4.chainmaker.org"
	orgId5                   = "wx-org5.chainmaker.org"
	orgId6                   = "wx-org6.chainmaker.org"
	//contractName   = "counter-go-1"
	certPathPrefix = "../../testdata"
	tlsHostName    = "chainmaker.org"
	version        = "1.0.0"
	upgradeVersion = "2.0.0"

	nodeAddr1 = "127.0.0.1:12301"
	connCnt1  = 5

	nodeAddr2 = "127.0.0.1:12301"
	connCnt2  = 5

	certPathFormat = "/crypto-config/%s/ca"

	createContractTimeout = 5
)

var (
	caPaths = []string{
		certPathPrefix + fmt.Sprintf(certPathFormat, orgId1),
	}

	adminKeyPath = certPathPrefix + "/crypto-config/%s/user/admin1/admin1.tls.key"
	adminCrtPath = certPathPrefix + "/crypto-config/%s/user/admin1/admin1.tls.crt"
)

const (

	// go 合约
	paillierContractName = "pailliergo10001"
	paillierByteCodePath = "../../testdata/paillier-wasm-demo/contract-paillier.wasm"
	runtime              = common.RuntimeType_GASM

	// rust 合约
	//paillierContractName = "paillier-rust-10001"
	//paillierByteCodePath = "./testdata/counter-go-demo/chainmaker_contract.wasm"
	//runtime              = common.RuntimeType_WASMER

	paillierPubKeyFilePath = "../../testdata/paillier-key/test1.pubKey"
	paillierPrvKeyFilePath = "../../testdata/paillier-key/test1.prvKey"
)

func main() {
	TestPaillierContractCounterGo()
}

func TestPaillierContractCounterGo() {
	t := new(testing.T)
	client, err := examples.CreateChainClientWithSDKConf(sdkConfigOrg1Client1Path)
	require.Nil(t, err)

	admin1, err := createAdmin(orgId1)
	require.Nil(t, err)
	admin2, err := createAdmin(orgId2)
	require.Nil(t, err)
	admin3, err := createAdmin(orgId3)
	require.Nil(t, err)
	admin4, err := createAdmin(orgId4)
	require.Nil(t, err)

	fmt.Println("======================================= 创建合约（异步）=======================================")
	testPaillierCreate(client, admin1, admin2, admin3, admin4, false)
	time.Sleep(5 * time.Second)

	fmt.Println("======================================= 调用合约运算（异步）=======================================")
	testPaillierOperation(t, client, "paillier_test_set", false)
	time.Sleep(5 * time.Second)

	fmt.Println("======================================= 查询结果并解密（异步）=======================================")
	testPaillierQueryResult(t, client, "paillier_test_get")
}

func createAdmin(orgId string) (*sdk.ChainClient, error) {
	if node1 == nil {
		node1 = createNode(nodeAddr1, connCnt1)
	}

	if node2 == nil {
		node2 = createNode(nodeAddr2, connCnt2)
	}

	config := cmlog.LogConfig{
		Module:       "[SDK]",
		LogPath:      "./sdk.log",
		LogLevel:     cmlog.LEVEL_DEBUG,
		MaxAge:       30,
		JsonFormat:   false,
		ShowLine:     true,
		LogInConsole: false,
	}

	logger, _ := cmlog.InitSugarLogger(&config)

	adminClient, err := sdk.NewChainClient(
		sdk.WithChainClientOrgId(orgId),
		sdk.WithChainClientChainId(chainId),
		sdk.WithChainClientLogger(logger),
		sdk.WithUserKeyFilePath(fmt.Sprintf(adminKeyPath, orgId)),
		sdk.WithUserCrtFilePath(fmt.Sprintf(adminCrtPath, orgId)),
		sdk.AddChainClientNodeConfig(node1),
		sdk.AddChainClientNodeConfig(node2),
	)
	if err != nil {
		return nil, err
	}

	//启用证书压缩（开启证书压缩可以减小交易包大小，提升处理性能）
	err = adminClient.EnableCertHash()
	if err != nil {
		return nil, err
	}

	return adminClient, nil
}

var (
	node1 *sdk.NodeConfig
	node2 *sdk.NodeConfig
)

// 创建节点
func createNode(nodeAddr string, connCnt int) *sdk.NodeConfig {
	node := sdk.NewNodeConfig(
		// 节点地址，格式：127.0.0.1:12301
		sdk.WithNodeAddr(nodeAddr),
		// 节点连接数
		sdk.WithNodeConnCnt(connCnt),
		// 节点是否启用TLS认证
		sdk.WithNodeUseTLS(true),
		// 根证书路径，支持多个
		sdk.WithNodeCAPaths(caPaths),
		// TLS Hostname
		sdk.WithNodeTLSHostName(tlsHostName),
	)

	return node
}

func testPaillierCreate(client *sdk.ChainClient, admin1, admin2, admin3,
	admin4 *sdk.ChainClient, withSyncResult bool) {

	resp, err := createUserContract(client, admin1, admin2, admin3, admin4,
		paillierContractName, examples.Version, paillierByteCodePath, runtime, []*common.KeyValuePair{}, withSyncResult)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Printf("CREATE contract-hibe-1 contract resp: %+v\n", resp)
}

func createUserContract(client *sdk.ChainClient, admin1, admin2, admin3, admin4 *sdk.ChainClient,
	contractName, version, byteCodePath string, runtime common.RuntimeType, kvs []*common.KeyValuePair, withSyncResult bool) (*common.TxResponse, error) {

	payload, err := client.CreateContractCreatePayload(contractName, version, byteCodePath, runtime, kvs)
	if err != nil {
		return nil, err
	}

	endorsers, err := examples.GetEndorsers(payload, admin1, admin2, admin3, admin4)
	if err != nil {
		return nil, err
	}

	resp, err := client.SendContractManageRequest(payload, endorsers, createContractTimeout, withSyncResult)
	if err != nil {
		return nil, err
	}

	// TODO: ??
	//err = examples.CheckProposalRequestResp(resp, true)
	//if err != nil {
	//	return nil, err
	//}

	return resp, nil
}

// 调用合约进行同态运算
func testPaillierOperation(t *testing.T, client *sdk.ChainClient, s string, b bool) {
	pubKeyBytes, err := ioutil.ReadFile(paillierPubKeyFilePath)
	//require.Nil(t, err)
	if err != nil {
		log.Fatalln(err)
	}
	payloadParams, err := CreatePaillierTransactionPayloadParams(pubKeyBytes, 1, 1000000)
	resp, err := client.InvokeContract(paillierContractName, s, "", payloadParams, -1, b)
	//require.Nil(t, err)
	if err != nil {
		log.Fatalln(err)
	}

	if resp.Code != common.TxStatusCode_SUCCESS {
		fmt.Printf("invoke contract failed, [code:%d]/[msg:%s]\n", resp.Code, resp.Message)
	}

	//if !b {
	//	fmt.Printf("invoke contract success, resp: [code:%d]/[msg:%s]/[txId:%s]\n", resp.Code, resp.Message, resp.ContractResult.Result)
	//} else {
	//	fmt.Printf("invoke contract success, resp: [code:%d]/[msg:%s]/[contractResult:%s]\n", resp.Code, resp.Message, resp.ContractResult)
	//}
}

func CreatePaillierTransactionPayloadParams(pubKeyBytes []byte, plaintext1, plaintext2 int64) ([]*common.KeyValuePair, error) {
	pubKey := new(paillier.PubKey)
	err := pubKey.Unmarshal(pubKeyBytes)
	if err != nil {

	}

	pt1 := new(big.Int).SetInt64(plaintext1)
	ciphertext1, err := pubKey.Encrypt(pt1)
	if err != nil {
		return nil, err
	}
	ct1Bytes, err := ciphertext1.Marshal()
	if err != nil {
		return nil, err
	}

	prv2, _ := paillier.GenKey()

	pub2, _ := prv2.GetPubKey()
	_, _ = pub2.Marshal()

	pt2 := new(big.Int).SetInt64(plaintext2)
	ciphertext2, err := pubKey.Encrypt(pt2)
	if err != nil {
		return nil, err
	}
	ct2Bytes, err := ciphertext2.Marshal()
	if err != nil {
		return nil, err
	}

	ct1Str := base64.StdEncoding.EncodeToString(ct1Bytes)
	ct2Str := base64.StdEncoding.EncodeToString(ct2Bytes)

	payloadParams := []*common.KeyValuePair{
		{
			Key:   "handletype",
			Value: []byte("SubCiphertext"),
		},
		{
			Key:   "para1",
			Value: []byte(ct1Str),
		},
		{
			Key:   "para2",
			Value: []byte(ct2Str),
		},
		{
			Key:   "pubkey",
			Value: pubKeyBytes,
		},
	}
	/*
		old
		payloadParams := make(map[string]string)
		//payloadParams["handletype"] = "AddCiphertext"
		//payloadParams["handletype"] = "AddPlaintext"
		payloadParams["handletype"] = "SubCiphertext"
		//payloadParams["handletype"] = "SubCiphertextStr"
		//payloadParams["handletype"] = "SubPlaintext"
		//payloadParams["handletype"] = "NumMul"

		payloadParams["para1"] = ct1Str
		payloadParams["para2"] = ct2Str
		payloadParams["pubkey"] = string(pubKeyBytes)
	*/

	return payloadParams, nil
}

// 查询同态执行结果并解密
func testPaillierQueryResult(t *testing.T, c *sdk.ChainClient, s string) {
	//paillierMethod := "AddCiphertext"
	//paillierMethod := "AddPlaintext"
	paillierMethod := "SubCiphertext"
	//paillierMethod := "SubCiphertextStr"
	//paillierMethod := "SubPlaintext"
	params1, err := QueryPaillierResult(c, paillierContractName, s, paillierMethod, -1, paillierPrvKeyFilePath)
	require.Nil(t, err)
	fmt.Printf("QUERY %s contract resp -> encrypt(cipher 10): %d\n", paillierContractName, params1)
}

func QueryPaillierResult(c *sdk.ChainClient, contractName, method, paillierDataItemId string, timeout int64, paillierPrvKeyPath string) (int64, error) {

	resultStr, err := QueryPaillierResultById(c, contractName, method, paillierDataItemId, timeout)
	if err != nil {
		return 0, err
	}

	ct := new(paillier.Ciphertext)
	resultBytes, err := base64.StdEncoding.DecodeString(string(resultStr))
	if err != nil {
		return 0, err
	}
	err = ct.Unmarshal(resultBytes)
	if err != nil {
		return 0, err
	}

	prvKey := new(paillier.PrvKey)

	prvKeyBytes, err := ioutil.ReadFile(paillierPrvKeyPath)
	if err != nil {
		return 0, fmt.Errorf("open paillierKey file failed, [err:%s]", err)
	}

	err = prvKey.Unmarshal(prvKeyBytes)
	if err != nil {
		return 0, err
	}

	decrypt, err := prvKey.Decrypt(ct)
	if err != nil {
		return 0, err
	}

	return decrypt.Int64(), nil
}

func QueryPaillierResultById(c *sdk.ChainClient, contractName, method, paillierMethod string, timeout int64) ([]byte, error) {
	pairs := []*common.KeyValuePair{
		{Key: "handletype", Value: []byte(paillierMethod)},
	}

	/*
		old
		pairsMap := make(map[string]string)
		pairsMap["handletype"] = paillierMethod
	*/

	resp, err := c.QueryContract(contractName, method, pairs, timeout)
	if err != nil {
		return nil, err
	}

	result := resp.ContractResult.Result

	return result, nil
}
