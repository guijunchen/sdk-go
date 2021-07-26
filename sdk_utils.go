package chainmaker_sdk_go

import (
	"encoding/pem"
	"errors"
	"fmt"
	"io/ioutil"

	"chainmaker.org/chainmaker/common/crypto/asym"
	bcx509 "chainmaker.org/chainmaker/common/crypto/x509"
	"chainmaker.org/chainmaker/pb-go/accesscontrol"
	"chainmaker.org/chainmaker/pb-go/common"
	"chainmaker.org/chainmaker/sdk-go/utils"
)

func SignPayload(keyBytes, crtBytes []byte, payload *common.Payload) (*common.EndorsementEntry, error) {
	key, err := asym.PrivateKeyFromPEM(keyBytes, nil)
	if err != nil {
		return nil, fmt.Errorf("asym.PrivateKeyFromPEM failed, %s", err)
	}

	blockCrt, _ := pem.Decode(crtBytes)
	crt, err := bcx509.ParseCertificate(blockCrt.Bytes)
	if err != nil {
		return nil, fmt.Errorf("bcx509.ParseCertificate failed, %s", err)
	}

	if len(crt.Subject.Organization) != 1 {
		return nil, errors.New("invalid certificate, certificate must contain one Organization")
	}

	signature, err := utils.SignPayload(key, crt, payload)
	if err != nil {
		return nil, fmt.Errorf("SignPayload failed, %s", err)
	}

	sender := &accesscontrol.Member{
		OrgId:      crt.Subject.Organization[0],
		MemberInfo: crtBytes,
		MemberType: accesscontrol.MemberType_CERT,
	}

	entry := &common.EndorsementEntry{
		Signer:    sender,
		Signature: signature,
	}

	return entry, nil
}

func SignPayloadWithPath(keyFilePath, crtFilePath string, payload *common.Payload) (*common.EndorsementEntry, error) {
	// 读取私钥
	keyBytes, err := ioutil.ReadFile(keyFilePath)
	if err != nil {
		return nil, fmt.Errorf("read key file failed, %s", err)
	}

	// 读取证书
	crtBytes, err := ioutil.ReadFile(crtFilePath)
	if err != nil {
		return nil, fmt.Errorf("read crt file failed, %s", err)
	}

	return SignPayload(keyBytes, crtBytes, payload)
}
