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

// Deprecated: This function will be deleted when appropriate. Please use SignPayloadV2
func SignPayload(privateKey crypto.PrivateKey, cert *bcx509.Certificate, payload *common.Payload) ([]byte, error) {
	payloadBytes, err := proto.Marshal(payload)
	if err != nil {
		return nil, err
	}

	return SignPayloadBytes(privateKey, cert, payloadBytes)
}

// Deprecated: This function will be deleted when appropriate. Please use SignPayloadBytesV2
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

	hashAlgo, err := bcx509.GetHashFromSignatureAlgorithm(cert.SignatureAlgorithm)
	if err != nil {
		return nil, err
	}

	return SignPayloadV2(key, hashAlgo, payload)
}

func SignPayloadWithPkPath(keyFilePath, hashType string, payload *common.Payload) ([]byte, error) {
	keyPem, err := ioutil.ReadFile(keyFilePath)
	if err != nil {
		return nil, fmt.Errorf("read key file failed, %s", err)
	}

	key, err := asym.PrivateKeyFromPEM(keyPem, nil)
	if err != nil {
		return nil, err
	}

	return SignPayloadV2(key, crypto.HashAlgoMap[hashType], payload)
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

func NewPkEndorser(orgId string, pk []byte, sig []byte) *common.EndorsementEntry {
	return &common.EndorsementEntry{
		Signer: &accesscontrol.Member{
			OrgId:      orgId,
			MemberInfo: pk,
			MemberType: accesscontrol.MemberType_PUBLIC_KEY,
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

	hashAlgo, err := bcx509.GetHashFromSignatureAlgorithm(cert.SignatureAlgorithm)
	if err != nil {
		return nil, err
	}

	signature, err := SignPayloadV2(key, hashAlgo, payload)
	if err != nil {
		return nil, err
	}

	var orgId string
	if len(cert.Subject.Organization) != 0 {
		orgId = cert.Subject.Organization[0]
	}

	return NewEndorser(orgId, certPem, signature), nil
}

func MakePkEndorserWithPem(keyPem []byte, hashType, orgId string, paylaod *common.Payload) (*common.EndorsementEntry, error) {
	key, err := asym.PrivateKeyFromPEM(keyPem, nil)
	if err != nil {
		return nil, err
	}

	signature, err := SignPayloadV2(key, crypto.HashAlgoMap[hashType], paylaod)
	if err != nil {
		return nil, err
	}

	return NewPkEndorser(orgId, keyPem, signature), nil
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

func MakePkEndorserWithPath(keyFilePath, hashType, orgId string, payload *common.Payload) (*common.EndorsementEntry, error) {
	keyPem, err := ioutil.ReadFile(keyFilePath)
	if err != nil {
		return nil, fmt.Errorf("read key file failed, %s", err)
	}

	return MakePkEndorserWithPem(keyPem, hashType, orgId, payload)
}
