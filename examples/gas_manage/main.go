package main

import (
	"fmt"
	"log"

	"chainmaker.org/chainmaker/common/v2/crypto"
	"chainmaker.org/chainmaker/common/v2/crypto/asym"
	"chainmaker.org/chainmaker/pb-go/v2/syscontract"
	sdk "chainmaker.org/chainmaker/sdk-go/v2"
	"chainmaker.org/chainmaker/sdk-go/v2/examples"
)

const (
	sdkConfigOrg1Client1Path = "../sdk_configs/sdk_config_pk_user1.yml"

	gasAdminPubKeyPem = "-----BEGIN PUBLIC KEY-----\nMIIBCgKCAQEAx+IzSqBDeZMPaEhbBg9i4vgXbbaWUZ5ISWvQqgajt910xEAaNW1x\n9XldcJn8G3HynPgyhBEruKDNEmMH3KszCGUXEbY0VssfXD/OaeFJqXBdfYq4lmKd\nnypO+CCrJh6Cu0QuUi2DgI0ZsnM/VJ2JKby8JSAFhBQOPN9QdyFVQRY4fAqQ4p9T\nxv4x2KMliJKLHDqonW7Puk9UUYA2AIpehBapGDR4Zwj2S4ExPCD38uR/y7cB/0KN\ntcNamFqO5GRuZfO5KC9CUZzbFi0iKq/N8lShRoAWFAmR5FzTlUwZ5R2Wn+ckZHZE\nu7LJiX56XcS8d19c/7k7LAM6cl7EcigNeQIDAQAB\n-----END PUBLIC KEY-----"
)

func main() {
	client, err := examples.CreateChainClientWithSDKConf(sdkConfigOrg1Client1Path)
	if err != nil {
		log.Fatalln(err)
	}

	adminPubKey, err := asym.PublicKeyFromPEM([]byte(gasAdminPubKeyPem))
	if err != nil {
		log.Fatal(err)
	}
	adminPubKeyBytes, err := adminPubKey.Bytes()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("====================== 设置 gas admin ======================")
	if err := setGasAdmin(client, adminPubKey); err != nil {
		log.Fatal(err)
	}
	fmt.Println("====================== 获取 gas admin ======================")
	if err := getGasAdmin(client); err != nil {
		log.Fatal(err)
	}
	fmt.Println("====================== 充值gas账户 100个gas ======================")
	rechargeGasList := []*syscontract.RechargeGas{
		{
			PublicKey: adminPubKeyBytes,
			GasAmount: 100,
		},
	}
	if err := rechargeGas(client, rechargeGasList); err != nil {
		log.Fatal(err)
	}
	fmt.Println("====================== 扣除gas账户 10个gas ======================")
	if err := chargeGas(client, adminPubKey, 10); err != nil {
		log.Fatal(err)
	}
	fmt.Println("====================== 查询gas账户余额 ======================")
	if err := getGasBalance(client, adminPubKey); err != nil {
		log.Fatal(err)
	}
	fmt.Println("====================== 退还gas账户 5个gas ======================")
	if err := refundGas(client, adminPubKey, 5); err != nil {
		log.Fatal(err)
	}
	fmt.Println("====================== 查询gas账户余额 ======================")
	if err := getGasBalance(client, adminPubKey); err != nil {
		log.Fatal(err)
	}
	fmt.Println("====================== 冻结指定gas账户 ======================")
	if err := frozenGasAccount(client, adminPubKey); err != nil {
		log.Fatal(err)
	}
	fmt.Println("====================== 查询gas账户的状态 ======================")
	if err := getGasAccountStatus(client, adminPubKey); err != nil {
		log.Fatal(err)
	}
	fmt.Println("====================== 解冻指定gas账户 ======================")
	if err := unfrozenGasAccount(client, adminPubKey); err != nil {
		log.Fatal(err)
	}
	fmt.Println("====================== 查询gas账户的状态 ======================")
	if err := getGasAccountStatus(client, adminPubKey); err != nil {
		log.Fatal(err)
	}
}

func setGasAdmin(cc *sdk.ChainClient, key crypto.PublicKey) error {
	payload, err := cc.CreateSetGasAdminPayload(key)
	if err != nil {
		return err
	}
	resp, err := cc.SendGasManageRequest(payload, nil, -1, true)
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

func chargeGas(cc *sdk.ChainClient, pubKey crypto.PublicKey, amount int64) error {
	payload, err := cc.CreateChargeGasPayload(pubKey, amount)
	if err != nil {
		return err
	}
	resp, err := cc.SendGasManageRequest(payload, nil, -1, true)
	if err != nil {
		return err
	}

	fmt.Printf("charge gas resp: %+v\n", resp)
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
