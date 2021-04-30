package chainmaker_sdk_go

import (
	"chainmaker.org/chainmaker-go/common/crypto"
	"chainmaker.org/chainmaker-sdk-go/pb/protogo/common"
	"fmt"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

// test contract functionName
const (
	// save Hibe Message
	saveHibeMsg = "save_hibe_msg"

	// save params
	saveHibeParams = "save_hibe_params"

	// find params by ogrId
	findParamsByOrgId = "find_params_by_org_id"
)

// test data
const (
	hibeContractByteCodePath = "./testdata/counter-go-demo/contract-hibe-1.0.0.wasm"

	//
	hibeContractName = "contract-hibe-1"

	// 本地 hibe params 文件路径
	localHibeParamsFilePath = "./testdata/hibe-data/wx-org1.chainmaker.org/wx-org1.chainmaker.org.params"

	// 测试源消息
	msg = "这是一条HIBE测试存证 ✔✔✔"

	// hibe_msg 的消息 Id
	bizId2 = "1234567890123452"

	// Id 和 对应私钥文件路径 这里测试3组
	localTopLevelId                 = "wx-topL"
	localTopLevelHibePrvKeyFilePath = "./testdata/hibe-data/wx-org1.chainmaker.org/privateKeys/wx-topL.privateKey"

	localSecondLevelId                 = "wx-topL/secondL"
	localSecondLevelHibePrvKeyFilePath = "./testdata/hibe-data/wx-org1.chainmaker.org/privateKeys/wx-topL_secondL.privateKey"

	localThirdLevelId                 = "wx-topL/secondL/thirdL"
	localThirdLevelHibePrvKeyFilePath = "./testdata/hibe-data/wx-org1.chainmaker.org/privateKeys/wx-topL.privateKey"
)

var txId = ""

func TestHibeContractCounterGo(t *testing.T) {

	txId = GetRandTxId()
	client, err := createClientWithConfig()
	require.Nil(t, err)

	admin1, err := createAdmin(orgId1)
	require.Nil(t, err)
	admin2, err := createAdmin(orgId2)
	require.Nil(t, err)
	admin3, err := createAdmin(orgId3)
	require.Nil(t, err)
	admin4, err := createAdmin(orgId4)
	require.Nil(t, err)

	fmt.Println("====================== 创建合约（异步）======================")
	testUserHibeContractCounterGoCreate(t, client, admin1, admin2, admin3, admin4, false)
	time.Sleep(5 * time.Second)

	fmt.Println("====================== 调用合约 params 上链 （异步）======================")
	testUserHibeContractParamsGoInvoke(t, client, saveHibeParams, false)
	time.Sleep(5 * time.Second)

	fmt.Println("====================== 执行合约 params 查询接口 ======================")
	testUserHibeContractParamsGoQuery(t, client, findParamsByOrgId, nil)

	fmt.Println("====================== 调用合约 加密数据上链（异步）======================")
	testUserHibeContractMsgGoInvoke(t, client, saveHibeMsg, false)
	time.Sleep(5 * time.Second)

	fmt.Println("====================== 执行合约 加密数据查询接口 ======================")
	testUserHibeContractMsgGoQuery(t, client)
}

// 创建Hibe合约
func testUserHibeContractCounterGoCreate(t *testing.T, client *ChainClient,
	admin1, admin2, admin3, admin4 *ChainClient, withSyncResult bool) {

	resp, err := createUserHibeContract(client, admin1, admin2, admin3, admin4,
		hibeContractName, version, hibeContractByteCodePath, common.RuntimeType_GASM, []*common.KeyValuePair{}, withSyncResult)
	require.Nil(t, err)

	t.Logf("CREATE contract-hibe-1 contract resp: %+v\n", resp)
}

// 调用Hibe合约
// params 上链
func testUserHibeContractParamsGoInvoke(t *testing.T, client *ChainClient,
	method string, withSyncResult bool) {
	err := invokeUserHibeContractParams(t, client, hibeContractName, method, "", withSyncResult)
	require.Nil(t, err)
}

// params 查询
func testUserHibeContractParamsGoQuery(t *testing.T, client *ChainClient,
	method string, params map[string]string) {
	hibeParams, err := client.QueryHibeParamsWithOrgId(hibeContractName, findParamsByOrgId, orgId1, -1)
	require.Nil(t, err)
	t.Logf("QUERY %s contract resp -> hibeParams:%s\n", hibeContractName, hibeParams)
}

// 加密数据上链
func testUserHibeContractMsgGoInvoke(t *testing.T, client *ChainClient,
	method string, withSyncResult bool) {
	err := invokeUserHibeContractMsg(t, client, hibeContractName, method, txId, withSyncResult)
	require.Nil(t, err)
}

// 获取加密数据
func testUserHibeContractMsgGoQuery(t *testing.T, client *ChainClient) {
	//keyType := crypto.AES
	keyType := crypto.SM4

	localParams, err := ReadHibeParamsWithFilePath(localHibeParamsFilePath)
	require.Nil(t, err)

	topHibePrvKey, err := ReadHibePrvKeysWithFilePath(localTopLevelHibePrvKeyFilePath)
	require.Nil(t, err)

	msgBytes1, err := client.DecryptHibeTxByTxId(localTopLevelId, localParams, topHibePrvKey, txId, keyType)
	require.Nil(t, err)
	t.Logf("QUERY hibe-contract-go-1 contract resp DecryptHibeTxByBizId [Decrypt Msg By TopLevel privateKey] message: %s\n", string(msgBytes1))

	msgBytes2, err := client.DecryptHibeTxByTxId(localTopLevelId, localParams, topHibePrvKey, txId, keyType)
	require.Nil(t, err)
	t.Logf("QUERY hibe-contract-go-1 contract resp DecryptHibeTxByBizId [Decrypt Msg By SecondLevel privateKey] message: %s\n", string(msgBytes2))

	msgBytes3, err := client.DecryptHibeTxByTxId(localTopLevelId, localParams, topHibePrvKey, txId, keyType)
	require.Nil(t, err)
	t.Logf("QUERY hibe-contract-go-1 contract resp DecryptHibeTxByBizId [Decrypt Msg By ThirdLevel privateKey] message: %s\n", string(msgBytes3))

}

func createUserHibeContract(client *ChainClient, admin1, admin2, admin3, admin4 *ChainClient,
	contractName, version, byteCodePath string, runtime common.RuntimeType, kvs []*common.KeyValuePair, withSyncResult bool) (*common.TxResponse, error) {

	payloadBytes, err := client.CreateContractCreatePayload(contractName, version, byteCodePath, runtime, kvs)
	if err != nil {
		return nil, err
	}

	// 各组织Admin权限用户签名
	signedPayloadBytes1, err := admin1.SignContractManagePayload(payloadBytes)
	if err != nil {
		return nil, err
	}

	signedPayloadBytes2, err := admin2.SignContractManagePayload(payloadBytes)
	if err != nil {
		return nil, err
	}

	signedPayloadBytes3, err := admin3.SignContractManagePayload(payloadBytes)
	if err != nil {
		return nil, err
	}

	signedPayloadBytes4, err := admin4.SignContractManagePayload(payloadBytes)
	if err != nil {
		return nil, err
	}

	// 收集并合并签名
	mergeSignedPayloadBytes, err := client.MergeContractManageSignedPayload([][]byte{signedPayloadBytes1,
		signedPayloadBytes2, signedPayloadBytes3, signedPayloadBytes4})
	if err != nil {
		return nil, err
	}

	// 发送创建合约请求
	resp, err := client.SendContractManageRequest(mergeSignedPayloadBytes, createContractTimeout, withSyncResult)
	if err != nil {
		return nil, err
	}

	err = checkProposalRequestResp(resp, true)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func invokeUserHibeContractParams(t *testing.T, client *ChainClient, contractName, method, txId string, withSyncResult bool) error {
	localParams, err := ReadHibeParamsWithFilePath(localHibeParamsFilePath)
	if err != nil {
		return err
	}
	payloadParams, err := client.CreateHibeInitParamsTxPayloadParams(orgId1, localParams)
	resp, err := client.InvokeContract(contractName, method, txId, payloadParams, -1, withSyncResult)
	if err != nil {
		return err
	}

	if resp.Code != common.TxStatusCode_SUCCESS {
		return fmt.Errorf("invoke contract failed, [code:%d]/[msg:%s]\n", resp.Code, resp.Message)
	}

	if !withSyncResult {
		t.Logf("invoke contract success, resp: [code:%d]/[msg:%s]/[txId:%s]\n", resp.Code, resp.Message, resp.ContractResult.Result)
	} else {
		t.Logf("invoke contract success, resp: [code:%d]/[msg:%s]/[contractResult:%s]\n", resp.Code, resp.Message, resp.ContractResult)
	}

	return nil
}

func invokeUserHibeContractMsg(t *testing.T, client *ChainClient, contractName, method, txId string, withSyncResult bool) error {
	receiverId := make([]string, 3)
	receiverId[0] = localSecondLevelId
	receiverId[1] = localThirdLevelId
	receiverId[2] = localTopLevelId

	// fetch orgId []string from receiverId []string
	org := make([]string, len(receiverId))
	org[0] = "wx-org1.chainmaker.org"
	org[1] = "wx-org1.chainmaker.org"
	org[2] = "wx-org1.chainmaker.org"

	// query params
	paramsBytesList := make([][]byte, 0)
	for _, id := range org {
		hibeParamsBytes, err := client.QueryHibeParamsWithOrgId(hibeContractName, findParamsByOrgId, id, -1)
		if err != nil {
			//t.Logf("QUERY hibe-contract-go-1 contract resp: %+v\n", hibeParams)
			return fmt.Errorf("client.QueryHibeParamsWithOrgId(hibeContractName, id, -1) failed, err: %v\n", err)
		}

		if len(hibeParamsBytes) == 0 {
			return fmt.Errorf("no souch params of %s's org, please check it", id)
		}

		paramsBytesList = append(paramsBytesList, hibeParamsBytes)
	}

	//keyType := crypto.AES
	keyType := crypto.SM4
	params, err := client.CreateHibeTxPayloadParamsWithHibeParams([]byte(msg), receiverId, paramsBytesList, txId, keyType)
	if err != nil {
		return err
	}

	resp, err := client.InvokeContract(contractName, method, txId, params, -1, withSyncResult)
	if err != nil {
		return err
	}

	if resp.Code != common.TxStatusCode_SUCCESS {
		return fmt.Errorf("invoke contract failed, [code:%d]/[msg:%s]\n", resp.Code, resp.Message)
	}

	if !withSyncResult {
		t.Logf("invoke contract success, resp: [code:%d]/[msg:%s]/[txId:%s]\n", resp.Code, resp.Message, resp.ContractResult.Result)
	} else {
		t.Logf("invoke contract success, resp: [code:%d]/[msg:%s]/[contractResult:%s]\n", resp.Code, resp.Message, resp.ContractResult)
	}

	return nil
}
