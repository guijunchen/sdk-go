/*
Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.

SPDX-License-Identifier: Apache-2.0
*/
package utils

import (
	"fmt"
	"io/ioutil"

	"chainmaker.org/chainmaker/common/v2/crypto"
	"chainmaker.org/chainmaker/common/v2/crypto/asym"
	bcx509 "chainmaker.org/chainmaker/common/v2/crypto/x509"
	"chainmaker.org/chainmaker/pb-go/v2/accesscontrol"
	"chainmaker.org/chainmaker/pb-go/v2/common"
	"github.com/gogo/protobuf/proto"
)

func SignPayload(privateKey crypto.PrivateKey, cert *bcx509.Certificate, payload *common.Payload) ([]byte, error) {
	payloadBytes, err := proto.Marshal(payload)
	if err != nil {
		return nil, err
	}

	return SignPayloadBytes(privateKey, cert, payloadBytes)
}

func SignPayloadBytes(privateKey crypto.PrivateKey, cert *bcx509.Certificate, payloadBytes []byte) ([]byte, error) {
	var opts crypto.SignOpts
	hashalgo, err := bcx509.GetHashFromSignatureAlgorithm(cert.SignatureAlgorithm)
	if err != nil {
		return nil, fmt.Errorf("invalid algorithm: %v", err)
	}

	opts.Hash = hashalgo
	opts.UID = crypto.CRYPTO_DEFAULT_UID

	return privateKey.SignWithOpts(payloadBytes, &opts)
}

func SignPayloadV2(privateKey crypto.PrivateKey, hashType crypto.HashType, payload *common.Payload) ([]byte, error) {
	payloadBytes, err := proto.Marshal(payload)
	if err != nil {
		return nil, err
	}

	return SignPayloadBytesV2(privateKey, hashType, payloadBytes)
}

func SignPayloadBytesV2(privateKey crypto.PrivateKey, hashType crypto.HashType, payloadBytes []byte) ([]byte, error) {

	var opts crypto.SignOpts
	opts.Hash = hashType
	opts.UID = crypto.CRYPTO_DEFAULT_UID

	return privateKey.SignWithOpts(payloadBytes, &opts)
}

func SignPayloadWithPath(keyFilePath, crtFilePath string, payload *common.Payload) ([]byte, error) {
	// 读取私钥
	keyPem, err := ioutil.ReadFile(keyFilePath)
	if err != nil {
		return nil, fmt.Errorf("read key file failed, %s", err)
	}

	// 读取证书
	certPem, err := ioutil.ReadFile(crtFilePath)
	if err != nil {
		return nil, fmt.Errorf("read cert file failed, %s", err)
	}

	key, err := asym.PrivateKeyFromPEM(keyPem, nil)
	if err != nil {
		return nil, err
	}

	cert, err := ParseCert(certPem)
	if err != nil {
		return nil, err
	}

	return SignPayload(key, cert, payload)
}

func NewEndorser(orgId string, certPem []byte, sig []byte) *common.EndorsementEntry {
	return &common.EndorsementEntry{
		Signer: &accesscontrol.Member{
			OrgId:      orgId,
			MemberInfo: certPem,
			MemberType: accesscontrol.MemberType_CERT,
		},
		Signature: sig,
	}
}

func MakeEndorserWithPem(keyPem, certPem []byte, payload *common.Payload) (*common.EndorsementEntry, error) {
	key, err := asym.PrivateKeyFromPEM(keyPem, nil)
	if err != nil {
		return nil, err
	}

	cert, err := ParseCert(certPem)
	if err != nil {
		return nil, err
	}

	signature, err := SignPayload(key, cert, payload)
	if err != nil {
		return nil, err
	}

	var orgId string
	if len(cert.Subject.Organization) != 0 {
		orgId = cert.Subject.Organization[0]
	}

	e := NewEndorser(orgId, certPem, signature)
	return e, nil
}

func MakeEndorserWithPath(keyFilePath, crtFilePath string, payload *common.Payload) (*common.EndorsementEntry, error) {
	// 读取私钥
	keyPem, err := ioutil.ReadFile(keyFilePath)
	if err != nil {
		return nil, fmt.Errorf("read key file failed, %s", err)
	}

	// 读取证书
	certPem, err := ioutil.ReadFile(crtFilePath)
	if err != nil {
		return nil, fmt.Errorf("read cert file failed, %s", err)
	}

	return MakeEndorserWithPem(keyPem, certPem, payload)
}
