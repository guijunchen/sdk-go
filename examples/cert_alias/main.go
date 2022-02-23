package main

import (
	"fmt"
	"log"

	"chainmaker.org/chainmaker/common/v2/crypto"
	sdk "chainmaker.org/chainmaker/sdk-go/v2"
	"chainmaker.org/chainmaker/sdk-go/v2/examples"
)

const (
	sdkConfigPath     = "../sdk_configs/sdk_config_org1_client1.yml"
	certAlias         = "mycertalias"
	certAliasUpdateTo = "mycertalias_updated"
)

func main() {
	cc, err := sdk.NewChainClient(
		sdk.WithConfPath(sdkConfigPath),
		sdk.WithChainClientAlias(certAlias),
	)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("====================== add alias ======================")
	resp, err := cc.AddAlias()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("resp: %+v\n", resp)

	fmt.Println("====================== query alias ======================")
	aliasInfos, err := cc.QueryCertsAlias([]string{certAlias})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("aliasInfos: %+v\n", aliasInfos)

	fmt.Println("====================== update alias ======================")
	updateAliasPayload := cc.CreateUpdateAliasPayload(certAliasUpdateTo, string(cc.GetCertPEM()))

	endorsers, err := examples.GetEndorsersWithAuthType(crypto.HashAlgoMap[cc.GetHashType()],
		cc.GetAuthType(), updateAliasPayload, examples.UserNameOrg1Admin1, examples.UserNameOrg2Admin1,
		examples.UserNameOrg3Admin1, examples.UserNameOrg4Admin1)
	if err != nil {
		log.Fatal(err)
	}
	resp2, err := cc.UpdateAlias(updateAliasPayload, endorsers, -1, true)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("resp: %+v\n", resp2)

	fmt.Println("====================== query alias ======================")
	aliasInfos2, err := cc.QueryCertsAlias([]string{certAlias})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("aliasInfos: %+v\n", aliasInfos2)

	fmt.Println("====================== delete alias ======================")
	deleteAliasPayload := cc.CreateDeleteCertsAliasPayload([]string{certAliasUpdateTo})

	endorsers2, err := examples.GetEndorsersWithAuthType(crypto.HashAlgoMap[cc.GetHashType()],
		cc.GetAuthType(), deleteAliasPayload, examples.UserNameOrg1Admin1, examples.UserNameOrg2Admin1,
		examples.UserNameOrg3Admin1, examples.UserNameOrg4Admin1)
	if err != nil {
		log.Fatal(err)
	}
	resp3, err := cc.UpdateAlias(deleteAliasPayload, endorsers2, -1, true)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("resp: %+v\n", resp3)
}
