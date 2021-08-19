package chainmaker_sdk_go

import (
	"encoding/hex"
	"encoding/pem"
	"errors"
	"fmt"
	"io/ioutil"

	"chainmaker.org/chainmaker/common/crypto/asym"
	bcx509 "chainmaker.org/chainmaker/common/crypto/x509"
	"chainmaker.org/chainmaker/common/evmutils"
	"chainmaker.org/chainmaker/common/serialize"
	"chainmaker.org/chainmaker/pb-go/accesscontrol"
	"chainmaker.org/chainmaker/pb-go/common"
	"chainmaker.org/chainmaker/sdk-go/utils"
)

func SignPayload(keyBytes, crtBytes []byte, payload *common.Payload) (*common.EndorsementEntry, error) {
	key, err := asym.PrivateKeyFromPEM(keyBytes, nil)
	if err != nil {
		return nil, fmt.Errorf("asym.PrivateKeyFromPEM failed, %s", err)
	}

	blockCrt, rest := pem.Decode(crtBytes)
	if len(rest) != 0 {
		return nil, errors.New("pem.Decode failed, invalid cert")
	}
	crt, err := bcx509.ParseCertificate(blockCrt.Bytes)
	if err != nil {
		return nil, fmt.Errorf("bcx509.ParseCertificate failed, %s", err)
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

func GetEVMAddressFromCertPath(certFilePath string) (string, error) {
	certBytes, err := ioutil.ReadFile(certFilePath)
	if err != nil {
		return "", fmt.Errorf("read cert file [%s] failed, %s", certFilePath, err)
	}

	return GetEVMAddressFromCertBytes(certBytes)
}

func GetEVMAddressFromCertBytes(certBytes []byte) (string, error) {
	block, _ := pem.Decode(certBytes)
	cert, err := bcx509.ParseCertificate(block.Bytes)
	if err != nil {
		return "", fmt.Errorf("ParseCertificate cert failed, %s", err)
	}

	ski := hex.EncodeToString(cert.SubjectKeyId)
	addrInt, err := evmutils.MakeAddressFromHex(ski)
	if err != nil {
		return "", fmt.Errorf("make address from cert SKI failed, %s", err)
	}

	return addrInt.String(), nil
}

func (cc *ChainClient) EasyCodecItemToParamsMap(items []*serialize.EasyCodecItem) map[string][]byte {
	return serialize.EasyCodecItemToParamsMap(items)
}
