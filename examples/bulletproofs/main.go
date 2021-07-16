package main

import (
	"encoding/base64"
	"fmt"
	"log"
	"strconv"
	"testing"
	"time"

	"chainmaker.org/chainmaker/common/crypto/bulletproofs"

	cmlog "chainmaker.org/chainmaker/common/log"
	sdk "chainmaker.org/chainmaker/sdk-go"
	"chainmaker.org/chainmaker/sdk-go/examples"

	"chainmaker.org/chainmaker/pb-go/common"
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
	// rust
	//bulletproofsContractName = "bulletproofs-rust-1001"
	//bulletproofsByteCodePath = "../../testdata/counter-go-demo/chainmaker_contract.wasm"
	//bulletproofsRuntime      = common.RuntimeType_WASMER

	// go
	bulletproofsContractName = "bulletproofsgo1001"
	bulletproofsByteCodePath = "../../testdata/bulletproofs-wasm-demo/contract-bulletproofs.wasm"
	bulletproofsRuntime      = common.RuntimeType_GASM

	// 链上合约SDK接口
	BulletProofsOpTypePedersenAddNum        = "PedersenAddNum"
	BulletProofsOpTypePedersenAddCommitment = "PedersenAddCommitment"
	BulletProofsOpTypePedersenSubNum        = "PedersenSubNum"
	BulletProofsOpTypePedersenSubCommitment = "PedersenSubCommitment"
	BulletProofsOpTypePedersenMulNum        = "PedersenMulNum"
	BulletProofsVerify                      = "BulletproofsVerify"

	// 测试数据
)

var (
	// 测试数据
	A            uint64 = 100
	X            uint64 = 20
	commitmentA1 []byte
	commitmentA2 []byte
	proofA1      []byte
	proofA2      []byte
	openingA1    []byte
)

func main() {
	TestBulletproofsContractCounterGo()
}

func TestBulletproofsContractCounterGo() {
	t := new(testing.T)
	client, err := examples.CreateChainClientWithSDKConf(sdkConfigOrg1Client1Path)
	require.Nil(t, err)

	//admin1, err := createAdmin(orgId1)
	//require.Nil(t, err)
	//admin2, err := createAdmin(orgId2)
	//require.Nil(t, err)
	//admin3, err := createAdmin(orgId3)
	//require.Nil(t, err)
	//admin4, err := createAdmin(orgId4)
	//require.Nil(t, err)
	//
	//fmt.Println("======================================= 创建合约（异步）=======================================")
	//testUserBulletproofsContractCounterGoCreate(client, admin1, admin2, admin3, admin4, false)
	//time.Sleep(5 * time.Second)

	funcName := BulletProofsOpTypePedersenAddNum
	//funcName := BulletProofsOpTypePedersenAddCommitment
	//funcName := BulletProofsOpTypePedersenSubNum
	//funcName := BulletProofsOpTypePedersenSubCommitment
	//funcName := BulletProofsOpTypePedersenMulNum

	fmt.Printf("============================= 调用合约 链上计算并存储 =============================\n")
	testBulletproofsSet(client, "bulletproofs_test_set", funcName, true)
	time.Sleep(5 * time.Second)

	fmt.Printf("============================= 查询计算结果 =============================\n")
	testBulletProofsGetOpResult(t, client, "bulletproofs_test_get", funcName, false)
	time.Sleep(5 * time.Second)

	fmt.Printf("============================= 调用合约验证 proof 和 查询的 commitment =============================\n")
	testBulletproofsVerify(client, "bulletproofs_test_set", BulletProofsVerify, true)
	time.Sleep(5 * time.Second)

	fmt.Printf("============================= 查询验证结果 =============================\n")
	testBulletProofsGetVerifyResult(t, client, "bulletproofs_test_get", BulletProofsVerify, false)
	time.Sleep(5 * time.Second)
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

func testUserBulletproofsContractCounterGoCreate(client *sdk.ChainClient, admin1, admin2, admin3,
	admin4 *sdk.ChainClient, withSyncResult bool) {

	resp, err := createUserContract(client, admin1, admin2, admin3, admin4,
		bulletproofsContractName, examples.Version, bulletproofsByteCodePath, bulletproofsRuntime, []*common.KeyValuePair{}, withSyncResult)
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

// 调用合约 链上计算并存储计算结果
func testBulletproofsSet(client *sdk.ChainClient, method string, opType string, b bool) {
	// 构造payloadParams
	payloadParams, err := constructBulletproofsSetData(opType)
	if err != nil {
		return
	}
	resp, err := client.InvokeContract(bulletproofsContractName, method, "", payloadParams, -1, b)
	if err != nil {
		fmt.Println(err)
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

func constructBulletproofsSetData(opType string) ([]*common.KeyValuePair, error) {
	// 1. 对原始数据生成承诺和证明
	var err error
	proofA1, commitmentA1, openingA1, err = bulletproofs.ProveRandomOpening(A)
	if err != nil {
		return nil, err
	}

	//_, commitmentX, openingX, err := bulletproofs.ProveRandomOpening(X)
	//if err != nil {
	// return nil, err
	//}

	// 2. 计算并生成证明
	proofA2, _, err = bulletproofs.ProveAfterAddNum(A, X, openingA1, commitmentA1)
	//proofA2, _, err = bulletproofs.ProveAfterSubNum(A, X, openingA1, commitmentA1)
	//proofA2, _, _, err = bulletproofs.ProveAfterAddCommitment(A, X, openingA1, openingX, commitmentA1, commitmentX)
	//proofA2, _, _, err = bulletproofs.ProveAfterSubCommitment(A, X, openingA1, openingX, commitmentA1, commitmentX)
	//proofA2, _, _, err = bulletproofs.ProveAfterMulNum(A, 10, openingA1, commitmentA1)
	//proofA2, _, _, err = bulletproofs.ProveAfterMulNum(A, X, openingA1, commitmentA1)
	if err != nil {
		return nil, err
	}

	// 3. 原始 commitment-proof 对儿 和 新生成的 proof 上链
	// 3.1. 构造上链 payloadParams
	base64CommitmentA1Str := base64.StdEncoding.EncodeToString(commitmentA1)
	XStr := strconv.FormatInt(int64(X), 10)
	//base64X := base64.StdEncoding.EncodeToString([]byte(XStr))
	//base64X := base64.StdEncoding.EncodeToString(commitmentX)

	//fmt.Printf("commitment: %s\n", base64CommitmentA1Str)
	//fmt.Printf("X: %s\n", base64X)

	payloadParams := []*common.KeyValuePair{
		{
			Key:   "handletype",
			Value: []byte(opType),
		},
		{
			Key:   "para1",
			Value: []byte(base64CommitmentA1Str),
		},
		{
			Key:   "para2",
			Value: []byte(XStr),
		},
	}

	/*
		old
		payloadParams := make(map[string]string, 4)
		payloadParams["para1"] = base64CommitmentA1Str
		//payloadParams["para2"] = base64X
		payloadParams["para2"] = XStr
		payloadParams["handletype"] = opType
	*/

	return payloadParams, nil
}

// 查询计算结果
func testBulletProofsGetOpResult(t *testing.T, client *sdk.ChainClient, method string, opType string, b bool) {
	var err error
	commitmentA2, err = queryBulletproofsCommitment(client, bulletproofsContractName, method, opType, -1)
	require.Nil(t, err)
	fmt.Printf("QUERY %s contract resp -> : %s\n", bulletproofsContractName, commitmentA2)
}

func queryBulletproofsCommitment(client *sdk.ChainClient, contractName, method, bpMethod string, timeout int64) ([]byte, error) {

	resultBytes, err := queryBulletProofsCommitmentByHandleType(client, contractName, method, bpMethod, timeout)
	if err != nil {
		return nil, err
	}

	if bpMethod != BulletProofsVerify {
		resultBytes, err = base64.StdEncoding.DecodeString(string(resultBytes))
		if err != nil {
			return nil, err
		}
	}

	return resultBytes, nil
}

func queryBulletProofsCommitmentByHandleType(client *sdk.ChainClient, contractName, method, bpMethod string, timeout int64) ([]byte, error) {
	pair := []*common.KeyValuePair{
		{Key: "handletype", Value: []byte(bpMethod)},
	}

	/*
		old
		pairsMap := make(map[string]string)
		pairsMap["handletype"] = bpMethod
	*/

	resp, err := client.QueryContract(contractName, method, pair, timeout)
	if err != nil {
		return nil, err
	}

	result := resp.ContractResult.Result

	return result, nil
}

// 调用合约验证 proof 和 查询的 commitment
func testBulletproofsVerify(client *sdk.ChainClient, method string, opType string, b bool) {
	// 构造payloadParams
	payloadParams, err := constructBulletproofsVerifyData(opType)
	if err != nil {
		return
	}
	resp, err := client.InvokeContract(bulletproofsContractName, method, "", payloadParams, -1, b)
	if err != nil {
		fmt.Println(err)
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

func constructBulletproofsVerifyData(opType string) ([]*common.KeyValuePair, error) {
	// 1. 对原始数据生成承诺和证明
	//var err error
	//proofA1, commitmentA1, openingA1, err = bulletproofs.ProveRandomOpening(A)
	//if err != nil {
	//	return nil, err
	//}
	//
	//// 2. 计算并生成证明
	//proofA2, _, err := bulletproofs.ProveAfterAddNum(A, X, openingA1, commitmentA1)
	//if err != nil {
	//	return nil, err
	//}

	// 3. 原始 commitment-proof 对儿 和 新生成的 proof 上链
	// 3.1. 构造上链 payloadParams
	base64CommitmentA2Str := base64.StdEncoding.EncodeToString(commitmentA2)
	base64ProofA2Str := base64.StdEncoding.EncodeToString(proofA2)
	//base64ProofA2Str := base64.StdEncoding.EncodeToString(proofA1)

	payloadParams := []*common.KeyValuePair{
		{
			Key:   "handletype",
			Value: []byte(opType),
		},
		{
			Key:   "para1",
			Value: []byte(base64ProofA2Str),
		},
		{
			Key:   "para2",
			Value: []byte(base64CommitmentA2Str),
		},
	}

	/*
		old
		payloadParams := make(map[string]string, 4)
		payloadParams["para1"] = base64ProofA2Str
		payloadParams["para2"] = base64CommitmentA2Str
		payloadParams["handletype"] = opType
	*/

	return payloadParams, nil
}

// 查询验证结果
func testBulletProofsGetVerifyResult(t *testing.T, client *sdk.ChainClient, method string, opType string, b bool) {
	result, err := queryBulletproofsCommitment(client, bulletproofsContractName, method, opType, -1)
	require.Nil(t, err)
	fmt.Printf("QUERY %s contract resp -> : %s\n", bulletproofsContractName, result)
}
