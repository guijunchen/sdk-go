package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"time"

	"chainmaker.org/chainmaker/common/v2/crypto"
	"chainmaker.org/chainmaker/pb-go/v2/common"
	"chainmaker.org/chainmaker/pb-go/v2/syscontract"
	"chainmaker.org/chainmaker/sdk-go/v2/examples"
)

const (
	contractByteCodePath = "../../testdata/claim-wasm-demo/rust-fact-2.0.0.wasm"
	contractName         = "claim123"
	contractVersion      = "v1.0.0"

	sdkConfigPKUser1Path = "../sdk_configs/sdk_config_pk_user1.yml"
)

var (
	contractRuntimeType = common.RuntimeType_WASMER.String()
)

func main() {
	cc, err := examples.CreateChainClientWithSDKConf(sdkConfigPKUser1Path)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println("====================== 发送线上多签部署合约的交易 ======================")
	kvs := newContractInitPairs() //构造交易发起 kv pairs
	payload := cc.CreateMultiSignReqPayload(kvs)
	resp, err := cc.MultiSignContractReq(payload)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Printf("send MultiSignContractReq resp: %+v\n", resp)

	fmt.Println("====================== 各链管理员开始投票 ======================")
	endorsers, err := examples.GetEndorsersWithAuthType(crypto.HashAlgoMap[cc.GetHashType()],
		cc.GetAuthType(), payload, examples.UserNameOrg1Admin1, examples.UserNameOrg2Admin1, examples.UserNameOrg3Admin1, examples.UserNameOrg4Admin1)
	if err != nil {
		log.Fatalln(err)
	}

	for _, e := range endorsers {
		fmt.Printf("====================== %s 投票 ======================\n", e.Signer.MemberInfo)
		time.Sleep(3 * time.Second)

		resp, err = cc.MultiSignContractVote(payload, e)
		if err != nil {
			log.Fatalln(err)
		}
		fmt.Printf("send MultiSignContractVote resp: %+v\n", resp)

		fmt.Println("====================== 查询本多签交易的投票情况 ======================")
		time.Sleep(3 * time.Second)
		resp, err = cc.MultiSignContractQuery(payload.TxId)
		if err != nil {
			log.Fatalln(err)
		}
		fmt.Printf("query MultiSignContractQuery resp: %+v\n", resp)
	}
}

func newContractInitPairs() []*common.KeyValuePair {
	wasmBin, err := ioutil.ReadFile(contractByteCodePath)
	if err != nil {
		log.Fatalln(err)
	}
	return []*common.KeyValuePair{
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
			Value: []byte(contractVersion),
		},
		{
			Key:   syscontract.InitContract_CONTRACT_BYTECODE.String(),
			Value: wasmBin,
		},
		{
			Key:   syscontract.InitContract_CONTRACT_RUNTIME_TYPE.String(),
			Value: []byte(contractRuntimeType),
		},
	}
}
