/**
 * @Author: jasonruan
 * @Date:   2020-12-01 10:12:25
 **/
package chainmaker_sdk_go

import (
	"bytes"
	"chainmaker.org/chainmaker-go/chainmaker-sdk-go/pb"
	"chainmaker.org/chainmaker-go/common/crypto"
	bcx509 "chainmaker.org/chainmaker-go/common/crypto/x509"
	"chainmaker.org/chainmaker-go/common/random/uuid"
	"encoding/pem"
	"errors"
	"fmt"
	"github.com/golang/protobuf/proto"
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

func paramsMap2KVPairs(params map[string]string) (kvPairs []*pb.KeyValuePair) {
	for key, val := range params {
		kvPair := &pb.KeyValuePair{
			Key: key,
			Value: val,
		}

		kvPairs = append(kvPairs, kvPair)
	}

	return
}
