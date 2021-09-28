package main

import (
	sdkutils "chainmaker.org/chainmaker/sdk-go/v2/utils"
	"fmt"
	"io/ioutil"
	"log"
	"time"

	"chainmaker.org/chainmaker/pb-go/v2/common"
	"chainmaker.org/chainmaker/pb-go/v2/syscontract"
	sdk "chainmaker.org/chainmaker/sdk-go/v2"
	"chainmaker.org/chainmaker/sdk-go/v2/examples"
)

var (
	WasmPath        = ""
	WasmUpgradePath = ""
	runtimeType     common.RuntimeType
	contractName    = ""
	payload         *common.Payload
	pairs           []*common.KeyValuePair
	signKeyPath     = "../../testdata/crypto-config/wx-org3.chainmaker.org/user/admin1/admin1.sign.key"
	signCrtPath     = "../../testdata/crypto-config/wx-org3.chainmaker.org/user/admin1/admin1.sign.crt"
)

const (
	sdkConfigOrg1Client1Path = "../sdk_configs/sdk_config_org1_client1.yml"
)

func init() {
	WasmPath = "../../testdata/claim-wasm-demo/rust-fact-2.0.0.wasm"
	WasmUpgradePath = WasmPath
	contractName = "contract101"
	runtimeType = common.RuntimeType_WASMER
}

func main() {
	client, err := examples.CreateChainClientWithSDKConf(sdkConfigOrg1Client1Path)
	if err != nil {
		log.Fatalln(err)
	}
	pairs = initContractInitPairs() //构造交易发起pairs
	payload = testMultiSignReq(client)

	time.Sleep(2 * time.Second)
	testMultiSignVote(client, payload)

	time.Sleep(2 * time.Second)
	testMultiSignQuery(client, payload.TxId)
}

func testMultiSignReq(client *sdk.ChainClient) *common.Payload {
	payload = client.CreateMultiSignReqPayload(pairs)
	resp, err := client.MultiSignContractReq(payload)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Printf("testMultiSignReq resp: %+v\n", resp)
	return payload
}

func testMultiSignVote(client *sdk.ChainClient, multiSignReqPayload *common.Payload) {

	endorser, err := sdkutils.MakeEndorserWithPath(signKeyPath, signCrtPath, multiSignReqPayload)
	if err != nil {
		log.Fatalln(err)
	}
	resp, err := client.MultiSignContractVote(multiSignReqPayload, endorser)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Printf("testMultiSignVote resp: %+v\n", resp)
}

func testMultiSignQuery(client *sdk.ChainClient, txId string) {
	resp, err := client.MultiSignContractQuery(txId)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Printf("testMultiSignQuery resp: %+v\n", resp)
}

func initContractInitPairs() []*common.KeyValuePair {
	wasmBin, err := ioutil.ReadFile(WasmPath)
	if err != nil {
		panic(err)
	}
	pairs := []*common.KeyValuePair{
		{
			Key:   syscontract.MultiReq_SYS_CONTRACT_NAME.String(),
			Value: []byte(syscontract.SystemContract_CONTRACT_MANAGE.String()),
		},
		{
			Key:   syscontract.MultiReq_SYS_METHOD.String(),
			Value: []byte(syscontract.ContractManageFunction_INIT_CONTRACT.String()),
		},
		{
			Key:   syscontract.InitContract_CONTRACT_NAME.String(),
			Value: []byte(contractName),
		},
		{
			Key:   syscontract.InitContract_CONTRACT_VERSION.String(),
			Value: []byte("1.0"),
		},
		{
			Key:   syscontract.InitContract_CONTRACT_BYTECODE.String(),
			Value: wasmBin,
		},
		{
			Key:   syscontract.InitContract_CONTRACT_RUNTIME_TYPE.String(),
			Value: []byte(runtimeType.String()),
		},
	}

	return pairs
}
