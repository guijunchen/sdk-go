/**
 * @Author: zghh
 * @Date:   2020-12-03 10:16:38
 **/
package chainmaker_sdk_go

import "chainmaker.org/chainmaker-go/chainmaker-sdk-go/pb"

func (cc ChainClient) CertAdd() (*pb.TxResponse, error) {
	panic("implement me")
}

func (cc ChainClient) CertDelete(certHashes string) (*pb.TxResponse, error) {
	panic("implement me")
}

func (cc ChainClient) CertQuery(certHashes string) (*pb.CertInfos, error) {
	panic("implement me")
}
