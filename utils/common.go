/*
Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.

SPDX-License-Identifier: Apache-2.0
*/
package utils

import (
	"encoding/pem"
	"errors"
	"fmt"
	"os"

	"chainmaker.org/chainmaker/common/v2/crypto/hash"
	bcx509 "chainmaker.org/chainmaker/common/v2/crypto/x509"
	"chainmaker.org/chainmaker/common/v2/random/uuid"
	"chainmaker.org/chainmaker/pb-go/v2/common"
)

const (
	// SUCCESS ContractResult success code
	SUCCESS uint32 = 0
)

func GetRandTxId() string {
	return uuid.GetUUID() + uuid.GetUUID()
}

func CheckProposalRequestResp(resp *common.TxResponse, needContractResult bool) error {
	if resp.Code != common.TxStatusCode_SUCCESS {
		return errors.New(resp.Message)
	}

	if needContractResult && resp.ContractResult == nil {
		return fmt.Errorf("contract result is nil")
	}

	if resp.ContractResult != nil && resp.ContractResult.Code != SUCCESS {
		return errors.New(resp.ContractResult.Message)
	}

	return nil
}

func GetCertificateId(certPEM []byte, hashType string) ([]byte, error) {
	if certPEM == nil {
		return nil, fmt.Errorf("get cert certPEM == nil")
	}
	certDer, _ := pem.Decode(certPEM)
	if certDer == nil {
		return nil, fmt.Errorf("invalid certificate")
	}
	return GetCertificateIdFromDER(certDer.Bytes, hashType)
}

func GetCertificateIdFromDER(certDER []byte, hashType string) ([]byte, error) {
	if certDER == nil {
		return nil, fmt.Errorf("get cert from der certDER == nil")
	}
	id, err := hash.GetByStrType(hashType, certDER)
	if err != nil {
		return nil, err
	}
	return id, nil
}

func ParseCert(crtPEM []byte) (*bcx509.Certificate, error) {
	certBlock, _ := pem.Decode(crtPEM)
	if certBlock == nil {
		return nil, fmt.Errorf("decode pem failed, invalid certificate")
	}

	cert, err := bcx509.ParseCertificate(certBlock.Bytes)
	if err != nil {
		return nil, fmt.Errorf("x509 parse cert failed, %s", err)
	}

	return cert, nil
}

func Exists(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		return os.IsExist(err)
	}
	return true
}

func IsArchived(txStatusCode common.TxStatusCode) bool {
	return txStatusCode == common.TxStatusCode_ARCHIVED_BLOCK || txStatusCode == common.TxStatusCode_ARCHIVED_TX
}

func IsArchivedString(txStatusCode string) bool {
	return txStatusCode == common.TxStatusCode_ARCHIVED_BLOCK.String() ||
		txStatusCode == common.TxStatusCode_ARCHIVED_TX.String()
}
