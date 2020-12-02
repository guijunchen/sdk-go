/**
 * @Author: jasonruan
 * @Date:   2020-11-27 15:14:08
 */
package chainmaker_sdk_go

import "chainmaker.org/chainmaker-go/chainmaker-sdk-go/pb"

type SDKInterface interface {
	Stop() error

	// [用户合约接口]
	// 合约创建
	ContractCreate(txId string, multiSignPayload []byte) (*pb.TxResponse, error)

	// 合约调用
	ContractInvoke(contractName, method, txId string, params map[string]string) (*pb.TxResponse, error)

	// 合约查询
	ContractQuery(contractName, method, txId string, params map[string]string) (*pb.TxResponse, error)

	// [系统合约接口]
	GetTxByTxId(txId string) (*pb.TransactionInfo, error)


	// TODO: [待分类接口]
	SendTransaction(tx *pb.Transaction) (*pb.TxResponse, error)

	GetChainConfigBeforeBlockHeight(blockHeight int) (*pb.ChainConfig, error)
}
