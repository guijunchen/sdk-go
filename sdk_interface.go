/**
 * @Author: jasonruan
 * @Date:   2020-11-27 15:14:08
 */
package chainmaker_sdk_go

import "chainmaker.org/chainmaker-go/chainmaker-sdk-go/pb"

type SDKInterface interface {
	Stop() error

	// 合约调用
	ContractInvoke(contractName, method, txId string, params map[string]string) (*pb.TxResponse, error)

	SendTransaction(tx *pb.Transaction) (*pb.TxResponse, error)

	GetChainConfigBeforeBlockHeight(blockHeight int) (*pb.ChainConfig, error)

	GetTxByTxId(txId string) (*pb.TransactionInfo, error)

}
