/**
 * @Author: jasonruan
 * @Date:   2020-12-01 10:12:25
 **/
package chainmaker_sdk_go

import (
	"bytes"
	"encoding/hex"
	"encoding/pem"
	"errors"
	"fmt"

	"chainmaker.org/chainmaker-go/chainmaker-sdk-go/pb"
	"chainmaker.org/chainmaker-go/common/crypto"
	"chainmaker.org/chainmaker-go/common/crypto/hash"
	bcx509 "chainmaker.org/chainmaker-go/common/crypto/x509"
	"chainmaker.org/chainmaker-go/common/random/uuid"
	"github.com/golang/protobuf/proto"
)

const (
	defaultSequence = 0
)

func GetRandTxId() string {
	return uuid.GetUUID() + uuid.GetUUID()
}

// calculate unsigned transaction bytes [header bytes || request payload bytes]
func CalcUnsignedTxBytes(t *pb.Transaction) ([]byte, error) {
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
func CalcUnsignedTxRequestBytes(txReq *pb.TxRequest) ([]byte, error) {
	if txReq == nil {
		return nil, errors.New("calc unsigned tx request bytes error, tx == nil")
	}

	return CalcUnsignedTxBytes(&pb.Transaction{
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

func paramsMap2KVPairs(params map[string]string) (kvPairs []*pb.KeyValuePair) {
	for key, val := range params {
		kvPair := &pb.KeyValuePair{
			Key:   key,
			Value: val,
		}

		kvPairs = append(kvPairs, kvPair)
	}

	return
}

func constructQueryPayload(contractName, method string, pairs []*pb.KeyValuePair) ([]byte, error) {
	payload := &pb.QueryPayload{
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

func constructTransactPayload(contractName, method string, pairs []*pb.KeyValuePair) ([]byte, error) {
	payload := &pb.TransactPayload{
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

func constructSystemContractPayload(chainId, contractName, method string, pairs []*pb.KeyValuePair, sequence uint64) ([]byte, error) {
	payload := &pb.SystemContractPayload{
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

func constructConfigUpdatePayload(chainId, contractName, method string, pairs []*pb.KeyValuePair, sequence int) ([]byte, error) {
	payload := &pb.SystemContractPayload{
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
	payload := &pb.SubscribeBlockPayload{
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

func constructSubscribeTxPayload(startBlock, endBlock int64, txType pb.TxType, txIds []string) ([]byte, error) {
	payload := &pb.SubscribeTxPayload{
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

func checkProposalRequestResp(resp *pb.TxResponse, needContractResult bool) error {
	if resp.Code != pb.TxStatusCode_SUCCESS {
		return errors.New(resp.Message)
	}

	if needContractResult && resp.ContractResult == nil {
		return fmt.Errorf("contract result is nil")
	}

	if resp.ContractResult != nil {
		if resp.ContractResult.Code != pb.ContractResultCode_OK {
			return errors.New(resp.ContractResult.Message)
		}
	}

	return nil
}

func mergeConfigUpdateSignedPayload(signedPayloadBytes [][]byte) ([]byte, error) {
	if len(signedPayloadBytes) == 0 {
		return nil, fmt.Errorf("input params is empty")
	}

	allPayload := &pb.SystemContractPayload{}
	if err := proto.Unmarshal(signedPayloadBytes[0], allPayload); err != nil {
		return nil, fmt.Errorf("unmarshal No.0 signed payload failed, %s", err)
	}

	if len(allPayload.Endorsement) != 1 || allPayload.Endorsement[0] == nil {
		return nil, fmt.Errorf("No.0 signed payload endorsement is empty")
	}

	allPayloadCopy := proto.Clone(allPayload)
	allPayloadCopy.(*pb.SystemContractPayload).Endorsement = nil

	for i := 1; i < len(signedPayloadBytes); i++ {

		payload := &pb.SystemContractPayload{}
		if err := proto.Unmarshal(signedPayloadBytes[i], payload); err != nil {
			return nil, fmt.Errorf("unmarshal No.%d signed payload failed, %s", i, err)
		}

		if len(payload.Endorsement) != 1 || payload.Endorsement[0] == nil {
			return nil, fmt.Errorf("No.%d signed payload endorsement is empty", i)
		}

		payloadCopy := proto.Clone(payload)
		payloadCopy.(*pb.SystemContractPayload).Endorsement = nil

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

	allPayload := &pb.ContractMgmtPayload{}
	if err := proto.Unmarshal(signedPayloadBytes[0], allPayload); err != nil {
		return nil, fmt.Errorf("unmarshal No.0 signed payload failed, %s", err)
	}

	if len(allPayload.Endorsement) != 1 || allPayload.Endorsement[0] == nil {
		return nil, fmt.Errorf("No.0 signed payload endorsement is empty")
	}

	allPayloadCopy := proto.Clone(allPayload)
	allPayloadCopy.(*pb.ContractMgmtPayload).Endorsement = nil

	for i := 1; i < len(signedPayloadBytes); i++ {

		payload := &pb.ContractMgmtPayload{}
		if err := proto.Unmarshal(signedPayloadBytes[i], payload); err != nil {
			return nil, fmt.Errorf("unmarshal No.%d signed payload failed, %s", i, err)
		}

		if len(payload.Endorsement) != 1 || payload.Endorsement[0] == nil {
			return nil, fmt.Errorf("No.%d signed payload endorsement is empty", i)
		}

		payloadCopy := proto.Clone(payload)
		payloadCopy.(*pb.ContractMgmtPayload).Endorsement = nil

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
