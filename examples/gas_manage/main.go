package main

import (
	"fmt"
	"log"

	"chainmaker.org/chainmaker/common/v2/crypto"
	"chainmaker.org/chainmaker/pb-go/v2/syscontract"
	sdk "chainmaker.org/chainmaker/sdk-go/v2"
	"chainmaker.org/chainmaker/sdk-go/v2/examples"
)

const (
	sdkConfigPKUser1Path = "../sdk_configs/sdk_config_pk_user1.yml"
)

func main() {
	client, err := examples.CreateChainClientWithSDKConf(sdkConfigPKUser1Path)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("====================== 设置client自己为 gas admin ======================")
	if err := setGasAdmin(client, client.GetPublicKey()); err != nil {
		log.Fatal(err)
	}
	fmt.Println("====================== 获取 gas admin ======================")
	if err := getGasAdmin(client); err != nil {
		log.Fatal(err)
	}
	fmt.Println("====================== 充值gas账户 100个gas ======================")
	gasAdminPubKeyStr, err := client.GetPublicKey().String()
	if err != nil {
		log.Fatal(err)
	}
	rechargeGasList := []*syscontract.RechargeGas{
		{
			PublicKey: []byte(gasAdminPubKeyStr),
			GasAmount: 900000000,
		},
	}
	if err := rechargeGas(client, rechargeGasList); err != nil {
		log.Fatal(err)
	}
	fmt.Println("====================== 查询gas账户余额 ======================")
	if err := getGasBalance(client, client.GetPublicKey()); err != nil {
		log.Fatal(err)
	}
	fmt.Println("====================== 退还gas账户 5个gas ======================")
	if err := refundGas(client, client.GetPublicKey(), 5); err != nil {
		log.Fatal(err)
	}
	fmt.Println("====================== 查询gas账户余额 ======================")
	if err := getGasBalance(client, client.GetPublicKey()); err != nil {
		log.Fatal(err)
	}
	fmt.Println("====================== 冻结指定gas账户 ======================")
	if err := frozenGasAccount(client, client.GetPublicKey()); err != nil {
		log.Fatal(err)
	}
	fmt.Println("====================== 查询gas账户的状态 ======================")
	if err := getGasAccountStatus(client, client.GetPublicKey()); err != nil {
		log.Fatal(err)
	}
	fmt.Println("====================== 解冻指定gas账户 ======================")
	if err := unfrozenGasAccount(client, client.GetPublicKey()); err != nil {
		log.Fatal(err)
	}
	fmt.Println("====================== 查询gas账户的状态 ======================")
	if err := getGasAccountStatus(client, client.GetPublicKey()); err != nil {
		log.Fatal(err)
	}
}

func setGasAdmin(cc *sdk.ChainClient, key crypto.PublicKey) error {
	payload, err := cc.CreateSetGasAdminPayload(key)
	if err != nil {
		return err
	}
	endorsers, err := examples.GetEndorsersWithAuthType(crypto.HashAlgoMap[cc.GetHashType()],
		cc.GetAuthType(), payload, examples.UserNameOrg1Admin1, examples.UserNameOrg2Admin1,
		examples.UserNameOrg3Admin1, examples.UserNameOrg4Admin1)
	if err != nil {
		return err
	}
	resp, err := cc.SendGasManageRequest(payload, endorsers, -1, true)
	if err != nil {
		return err
	}

	fmt.Printf("set gas admin resp: %+v\n", resp)
	return nil
}

func getGasAdmin(cc *sdk.ChainClient) error {
	adminAddr, err := cc.GetGasAdmin()
	if err != nil {
		return err
	}
	fmt.Printf("get gas admin address: %+v\n", adminAddr)
	return nil
}

func rechargeGas(cc *sdk.ChainClient, rechargeGasList []*syscontract.RechargeGas) error {
	payload, err := cc.CreateRechargeGasPayload(rechargeGasList)
	if err != nil {
		return err
	}
	resp, err := cc.SendGasManageRequest(payload, nil, -1, true)
	if err != nil {
		return err
	}

	fmt.Printf("recharge gas resp: %+v\n", resp)
	return nil
}

func getGasBalance(cc *sdk.ChainClient, pubKey crypto.PublicKey) error {
	balance, err := cc.GetGasBalance(pubKey)
	if err != nil {
		return err
	}

	fmt.Printf("get gas balance: %+v\n", balance)
	return nil
}

func refundGas(cc *sdk.ChainClient, pubKey crypto.PublicKey, amount int64) error {
	payload, err := cc.CreateRefundGasPayload(pubKey, amount)
	if err != nil {
		return err
	}
	resp, err := cc.SendGasManageRequest(payload, nil, -1, true)
	if err != nil {
		return err
	}

	fmt.Printf("refund gas resp: %+v\n", resp)
	return nil
}

func frozenGasAccount(cc *sdk.ChainClient, pubKey crypto.PublicKey) error {
	payload, err := cc.CreateFrozenGasAccountPayload(pubKey)
	if err != nil {
		return err
	}
	resp, err := cc.SendGasManageRequest(payload, nil, -1, true)
	if err != nil {
		return err
	}

	fmt.Printf("frozen gas account resp: %+v\n", resp)
	return nil
}

func unfrozenGasAccount(cc *sdk.ChainClient, pubKey crypto.PublicKey) error {
	payload, err := cc.CreateUnfrozenGasAccountPayload(pubKey)
	if err != nil {
		return err
	}
	resp, err := cc.SendGasManageRequest(payload, nil, -1, true)
	if err != nil {
		return err
	}

	fmt.Printf("unfrozen gas account resp: %+v\n", resp)
	return nil
}

func getGasAccountStatus(cc *sdk.ChainClient, pubKey crypto.PublicKey) error {
	status, err := cc.GetGasAccountStatus(pubKey)
	if err != nil {
		return err
	}

	fmt.Printf("get gas account status: %+v\n", status)
	return nil
}
