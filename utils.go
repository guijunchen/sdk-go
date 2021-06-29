/*
Copyright (C) BABEC. All rights reserved.
Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.

SPDX-License-Identifier: Apache-2.0
*/

package chainmaker_sdk_go

import (
	"bytes"
	"chainmaker.org/chainmaker-go/common/crypto"
	"chainmaker.org/chainmaker-go/common/crypto/hash"
	localhibe "chainmaker.org/chainmaker-go/common/crypto/hibe"
	bcx509 "chainmaker.org/chainmaker-go/common/crypto/x509"
	"chainmaker.org/chainmaker-go/common/random/uuid"
	"chainmaker.org/chainmaker-sdk-go/pb/protogo/accesscontrol"
	"chainmaker.org/chainmaker-sdk-go/pb/protogo/common"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"github.com/gogo/protobuf/proto"
	"github.com/samkumar/hibe"
	"io/ioutil"
	"os"
	"time"
)

const (
	defaultSequence = 0
)

func GetRandTxId() string {
	return uuid.GetUUID() + uuid.GetUUID()
}

// calculate unsigned transaction bytes [header bytes || request payload bytes]
func CalcUnsignedTxBytes(t *common.Transaction) ([]byte, error) {
	if t == nil {
		return nil, errors.New("calc unsigned tx bytes error, tx == nil")
	}

	headerBytes, err := proto.Marshal(t.Header)
	if err != nil {
		return nil, err
	}

	rawTxBytes := bytes.Join([][]byte{headerBytes, t.RequestPayload}, []byte{})

	return rawTxBytes, nil
}

// calculate unsigned transaction request bytes
func CalcUnsignedTxRequestBytes(txReq *common.TxRequest) ([]byte, error) {
	if txReq == nil {
		return nil, errors.New("calc unsigned tx request bytes error, tx == nil")
	}

	return CalcUnsignedTxBytes(&common.Transaction{
		Header:         txReq.Header,
		RequestPayload: txReq.Payload,
	})
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

func SignTx(privateKey crypto.PrivateKey, cert *bcx509.Certificate, msg []byte) ([]byte, error) {
	var opts crypto.SignOpts
	hashalgo, err := bcx509.GetHashFromSignatureAlgorithm(cert.SignatureAlgorithm)
	if err != nil {
		return nil, fmt.Errorf("invalid algorithm: %v", err)
	}

	opts.Hash = hashalgo
	opts.UID = crypto.CRYPTO_DEFAULT_UID

	return privateKey.SignWithOpts(msg, &opts)
}

func signPayload(privateKey crypto.PrivateKey, cert *bcx509.Certificate, msg []byte) ([]byte, error) {
	return SignTx(privateKey, cert, msg)
}

func paramsMap2KVPairs(params map[string]string) (kvPairs []*common.KeyValuePair) {
	for key, val := range params {
		kvPair := &common.KeyValuePair{
			Key:   key,
			Value: val,
		}

		kvPairs = append(kvPairs, kvPair)
	}

	return
}

func constructQueryPayload(contractName, method string, pairs []*common.KeyValuePair) ([]byte, error) {
	payload := &common.QueryPayload{
		ContractName: contractName,
		Method:       method,
		Parameters:   pairs,
	}

	payloadBytes, err := proto.Marshal(payload)
	if err != nil {
		return nil, err
	}

	return payloadBytes, nil
}

func constructTransactPayload(contractName, method string, pairs []*common.KeyValuePair) ([]byte, error) {
	payload := &common.TransactPayload{
		ContractName: contractName,
		Method:       method,
		Parameters:   pairs,
	}

	payloadBytes, err := proto.Marshal(payload)
	if err != nil {
		return nil, err
	}

	return payloadBytes, nil
}

func constructSystemContractPayload(chainId, contractName, method string, pairs []*common.KeyValuePair, sequence uint64) ([]byte, error) {
	payload := &common.SystemContractPayload{
		ChainId:      chainId,
		ContractName: contractName,
		Method:       method,
		Parameters:   pairs,
		Sequence:     sequence,
	}

	payloadBytes, err := proto.Marshal(payload)
	if err != nil {
		return nil, err
	}

	return payloadBytes, nil
}

func constructConfigUpdatePayload(chainId, contractName, method string, pairs []*common.KeyValuePair, sequence int) ([]byte, error) {
	payload := &common.SystemContractPayload{
		ChainId:      chainId,
		ContractName: contractName,
		Method:       method,
		Parameters:   pairs,
		Sequence:     uint64(sequence),
	}

	payloadBytes, err := proto.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("ConfigUpdatePayload marshal failed, %s", err)
	}

	return payloadBytes, nil
}

func constructSubscribeBlockPayload(startBlock, endBlock int64, withRwSet bool) ([]byte, error) {
	payload := &common.SubscribeBlockPayload{
		StartBlock: startBlock,
		EndBlock:   endBlock,
		WithRwSet:  withRwSet,
	}

	payloadBytes, err := proto.Marshal(payload)
	if err != nil {
		return nil, err
	}

	return payloadBytes, nil
}

func constructSubscribeContractEventPayload(topic, contractName string) ([]byte, error) {
	payload := &common.SubscribeContractEventPayload{
		Topic:        topic,
		ContractName: contractName,
	}

	payloadBytes, err := proto.Marshal(payload)
	if err != nil {
		return nil, err
	}

	return payloadBytes, nil
}

func constructSubscribeTxPayload(startBlock, endBlock int64, txType common.TxType, txIds []string) ([]byte, error) {
	payload := &common.SubscribeTxPayload{
		StartBlock: startBlock,
		EndBlock:   endBlock,
		TxType:     txType,
		TxIds:      txIds,
	}

	payloadBytes, err := proto.Marshal(payload)
	if err != nil {
		return nil, err
	}

	return payloadBytes, nil
}

func IsArchived(txStatusCode common.TxStatusCode) bool {
	if txStatusCode == common.TxStatusCode_ARCHIVED_BLOCK || txStatusCode == common.TxStatusCode_ARCHIVED_TX {
		return true
	}

	return false
}

func IsArchivedString(txStatusCode string) bool {
	if txStatusCode == common.TxStatusCode_ARCHIVED_BLOCK.String() ||
		txStatusCode == common.TxStatusCode_ARCHIVED_TX.String() {

		return true
	}

	return false
}

func checkProposalRequestResp(resp *common.TxResponse, needContractResult bool) error {
	if resp.Code != common.TxStatusCode_SUCCESS {
		return errors.New(resp.Message)
	}

	if needContractResult && resp.ContractResult == nil {
		return fmt.Errorf("contract result is nil")
	}

	if resp.ContractResult != nil && resp.ContractResult.Code != common.ContractResultCode_OK {
		return errors.New(resp.ContractResult.Message)
	}

	return nil
}

func (cc *ChainClient) signSystemContractPayload(payloadBytes []byte) ([]byte, error) {
	payload := &common.SystemContractPayload{}
	if err := proto.Unmarshal(payloadBytes, payload); err != nil {
		return nil, fmt.Errorf("unmarshal system contract payload failed, %s", err)
	}

	signBytes, err := signPayload(cc.privateKey, cc.userCrt, payloadBytes)
	if err != nil {
		return nil, fmt.Errorf("SignPayload failed, %s", err)
	}

	sender := &accesscontrol.SerializedMember{
		OrgId:      cc.orgId,
		MemberInfo: cc.userCrtBytes,
		IsFullCert: true,
	}

	entry := &common.EndorsementEntry{
		Signer:    sender,
		Signature: signBytes,
	}

	payload.Endorsement = []*common.EndorsementEntry{
		entry,
	}

	signedPayloadBytes, err := proto.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("marshal system contract sigend payload failed, %s", err)
	}

	return signedPayloadBytes, nil
}


func mergeSystemContractSignedPayload(signedPayloadBytes [][]byte) ([]byte, error) {
	if len(signedPayloadBytes) == 0 {
		return nil, fmt.Errorf("input params is empty")
	}

	allPayload := &common.SystemContractPayload{}
	if err := proto.Unmarshal(signedPayloadBytes[0], allPayload); err != nil {
		return nil, fmt.Errorf("unmarshal No.0 signed payload failed, %s", err)
	}

	if len(allPayload.Endorsement) != 1 || allPayload.Endorsement[0] == nil {
		return nil, fmt.Errorf("No.0 signed payload endorsement is empty")
	}

	allPayloadCopy := proto.Clone(allPayload)
	allPayloadCopy.(*common.SystemContractPayload).Endorsement = nil

	for i := 1; i < len(signedPayloadBytes); i++ {

		payload := &common.SystemContractPayload{}
		if err := proto.Unmarshal(signedPayloadBytes[i], payload); err != nil {
			return nil, fmt.Errorf("unmarshal No.%d signed payload failed, %s", i, err)
		}

		if len(payload.Endorsement) != 1 || payload.Endorsement[0] == nil {
			return nil, fmt.Errorf("No.%d signed payload endorsement is empty", i)
		}

		payloadCopy := proto.Clone(payload)
		payloadCopy.(*common.SystemContractPayload).Endorsement = nil

		if !checkPayloads(allPayloadCopy, payloadCopy) {
			return nil, fmt.Errorf("No.%d signed payload not all the same", i)
		}

		allPayload.Endorsement = append(allPayload.Endorsement, payload.Endorsement[0])
	}

	mergeSignedPayloadBytes, err := proto.Marshal(allPayload)
	if err != nil {
		return nil, fmt.Errorf("marshal merge signed payload failed, %s", err)
	}

	return mergeSignedPayloadBytes, nil
}

func mergeContractManageSignedPayload(signedPayloadBytes [][]byte) ([]byte, error) {
	if len(signedPayloadBytes) == 0 {
		return nil, fmt.Errorf("input params is empty")
	}

	allPayload := &common.ContractMgmtPayload{}
	if err := proto.Unmarshal(signedPayloadBytes[0], allPayload); err != nil {
		return nil, fmt.Errorf("unmarshal No.0 signed payload failed, %s", err)
	}

	if len(allPayload.Endorsement) != 1 || allPayload.Endorsement[0] == nil {
		return nil, fmt.Errorf("No.0 signed payload endorsement is empty")
	}

	allPayloadCopy := proto.Clone(allPayload)
	allPayloadCopy.(*common.ContractMgmtPayload).Endorsement = nil

	for i := 1; i < len(signedPayloadBytes); i++ {

		payload := &common.ContractMgmtPayload{}
		if err := proto.Unmarshal(signedPayloadBytes[i], payload); err != nil {
			return nil, fmt.Errorf("unmarshal No.%d signed payload failed, %s", i, err)
		}

		if len(payload.Endorsement) != 1 || payload.Endorsement[0] == nil {
			return nil, fmt.Errorf("No.%d signed payload endorsement is empty", i)
		}

		payloadCopy := proto.Clone(payload)
		payloadCopy.(*common.ContractMgmtPayload).Endorsement = nil

		if !checkPayloads(allPayloadCopy, payloadCopy) {
			return nil, fmt.Errorf("No.%d signed payload not all the same", i)
		}

		allPayload.Endorsement = append(allPayload.Endorsement, payload.Endorsement[0])
	}

	mergeSignedPayloadBytes, err := proto.Marshal(allPayload)
	if err != nil {
		return nil, fmt.Errorf("marshal merge signed payload failed, %s", err)
	}

	return mergeSignedPayloadBytes, nil
}

func checkPayloads(a, b proto.Message) bool {
	aBytes, err := proto.Marshal(a)
	if err != nil {
		return false
	}

	bBytes, err := proto.Marshal(b)
	if err != nil {
		return false
	}

	return bytes.Equal(aBytes, bBytes)
}

// on input a certificate in PEM format, a hash algorithm (should be the one in chain configuration), output the identity of the certificate in the form of a string (under hexadecimal encoding)
func getCertificateIdHex(certPEM []byte, hashType string) (string, error) {
	id, err := getCertificateId(certPEM, hashType)
	if err != nil {
		return "", err
	}
	idHex := hex.EncodeToString(id)
	return idHex, nil
}

func getCertificateId(certPEM []byte, hashType string) ([]byte, error) {
	if certPEM == nil {
		return nil, fmt.Errorf("get cert certPEM == nil")
	}
	certDer, _ := pem.Decode(certPEM)
	if certDer == nil {
		return nil, fmt.Errorf("invalid certificate")
	}
	return getCertificateIdFromDER(certDer.Bytes, hashType)
}

func getCertificateIdFromDER(certDER []byte, hashType string) ([]byte, error) {
	if certDER == nil {
		return nil, fmt.Errorf("get cert from der certDER == nil")
	}
	id, err := hash.GetByStrType(hashType, certDER)
	if err != nil {
		return nil, err
	}
	return id, nil
}

// return current unix timestamp in seconds
func CurrentTimeSeconds() int64 {
	return time.Now().Unix()
}

func CurrentTimeMillisSeconds() int64 {
	return time.Now().UnixNano() / 1e6
}

func BytesToInt(bys []byte, order binary.ByteOrder) (int32, error) {
	bytebuff := bytes.NewBuffer(bys)
	var data int32
	if err := binary.Read(bytebuff, order, &data); err != nil {
		return -1, err
	}

	return data, nil
}

func Exists(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}

// Returns the serialized byte array of hibeParams
func ReadHibeParamsWithFilePath(hibeParamsFilePath string) ([]byte, error) {
	paramsBytes, err := ioutil.ReadFile(hibeParamsFilePath)
	if err != nil {
		return nil, fmt.Errorf("open hibe params file failed, [err:%s]", err)
	}

	return paramsBytes, nil
}

// Returns the serialized byte array of hibePrvKey
func ReadHibePrvKeysWithFilePath(hibePrvKeyFilePath string) ([]byte, error) {
	prvKeyBytes, err := ioutil.ReadFile(hibePrvKeyFilePath)
	if err != nil {
		return nil, fmt.Errorf("open hibe privateKey file failed, [err:%s]", err)
	}

	return prvKeyBytes, nil
}

func DecryptHibeTx(localId string, hibeParams []byte, hibePrvKey []byte, tx *common.Transaction, keyType crypto.KeyType) ([]byte, error) {
	localParams, ok := new(hibe.Params).Unmarshal(hibeParams)
	if !ok {
		return nil, errors.New("hibe.Params.Unmarshal failed, please check your file")
	}

	prvKey, ok := new(hibe.PrivateKey).Unmarshal(hibePrvKey)
	if !ok {
		return nil, errors.New("hibe.PrivateKey.Unmarshal failed, please check your file")
	}

	// get hibe_msg from tx
	requestPayload := &common.QueryPayload{}
	err := proto.Unmarshal(tx.RequestPayload, requestPayload)
	if err != nil {
		return nil, err
	}
	hibeMsgMap := make(map[string]string)
	for _, item := range requestPayload.Parameters {
		if item.Key == HibeMsgKey {
			err = json.Unmarshal([]byte(item.Value), &hibeMsgMap)
			if err != nil {
				return nil, err
			}
		}
	}

	if hibeMsgMap == nil {
		return nil, errors.New("no such message, please check transaction")
	}
	return localhibe.DecryptHibeMsg(localId, localParams, prvKey, hibeMsgMap, keyType)
}
