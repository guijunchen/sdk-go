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

	// 合约升级
	ContractUpgrade(txId string, multiSignPayload []byte) (*pb.TxResponse, error)

	// 合约调用
	ContractInvoke(contractName, method, txId string, params map[string]string) (*pb.TxResponse, error)

	// 合约查询
	ContractQuery(contractName, method, txId string, params map[string]string) (*pb.TxResponse, error)

	// [系统合约接口]
	// 根据交易Id查询交易
	GetTxByTxId(txId string) (*pb.TransactionInfo, error)

	// 根据区块高度查询区块
	GetBlockByHeight(blockHeight int64, withRWSet bool) (*pb.BlockInfo, error)

	// 根据区块哈希查询区块
	GetBlockByHash(blockHash string, withRWSet bool) (*pb.BlockInfo, error)

	// 根据交易Id查询区块
	GetBlockByTxId(txId string, withRWSet bool) (*pb.BlockInfo, error)

	// 查询最新的配置块
	GetLastConfigBlock(withRWSet bool) (*pb.BlockInfo, error)

	// 查询合约信息
	GetContractInfo() (*pb.ContractInfo, error)

	// 查询节点加入的链信息
	GetNodeChainList() (*pb.ChainList, error)

	// 查询链信息
	GetChainInfo() (*pb.ChainInfo, error)

	// TODO: [待分类接口]
	SendTransaction(tx *pb.Transaction) (*pb.TxResponse, error)

	GetChainConfigBeforeBlockHeight(blockHeight int) (*pb.ChainConfig, error)
}
