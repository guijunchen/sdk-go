/**
 * @Author: jasonruan
 * @Date:   2020-11-27 15:14:08
 */
package chainmaker_sdk_go

import (
	"context"

	"chainmaker.org/chainmaker-go/chainmaker-sdk-go/pb"
)

type SDKInterface interface {
	Stop() error

	// =============================== BEGIN =============================
	// ======================== [(1/5)用户合约接口] ========================
	// 合约创建
	// 参数说明：
	//   - txId: 交易ID
	//           格式要求：长度为64bit，字符在a-z0-9
	//           可为空，若为空字符串，将自动生成，在pb.TxResponse.ContractResult.Result字段中返回该自动生成的txId
	//   - multiSignedPayload: 经多签后的payload数据
	ContractCreate(txId string, multiSignedPayload []byte) (*pb.TxResponse, error)

	// 合约升级
	// 参数说明：
	//   - txId: 交易ID
	//           格式要求：长度为64bit，字符在a-z0-9
	//           可为空，若为空字符串，将自动生成，在pb.TxResponse.ContractResult.Result字段中返回该自动生成的txId
	//   - multiSignedPayload: 经多签后的payload数据
	ContractUpgrade(txId string, multiSignedPayload []byte) (*pb.TxResponse, error)

	// 合约调用
	// 参数说明：
	//   - contractName: 合约名称
	//   - method: 合约方法
	//   - txId: 交易ID
	//           格式要求：长度为64bit，字符在a-z0-9
	//           可为空，若为空字符串，将自动生成，在pb.TxResponse.ContractResult.Result字段中返回该自动生成的txId
	//   - params: 合约参数
	ContractInvoke(contractName, method, txId string, params map[string]string) (*pb.TxResponse, error)

	// 合约查询接口调用
	// 参数说明：
	//   - contractName: 合约名称
	//   - method: 合约方法
	//   - params: 合约参数
	ContractQuery(contractName, method string, params map[string]string) (*pb.TxResponse, error)

	// ======================== [(1/5)用户合约接口] ========================
	// ================================ END ==============================

	// =============================== BEGIN =============================
	// ======================== [(2/5)系统合约接口] ========================
	// 根据交易Id查询交易
	// 参数说明：
	//   - txId: 交易ID
	GetTxByTxId(txId string) (*pb.TransactionInfo, error)

	// 根据区块高度查询区块
	// 参数说明：
	//   - blockHeight: 指定区块高度，若为-1，将返回最新区块
	//   - withRWSet: 是否返回读写集
	GetBlockByHeight(blockHeight int64, withRWSet bool) (*pb.BlockInfo, error)

	// 根据区块哈希查询区块
	// 参数说明：
	//   - blockHash: 指定区块Hash
	//   - withRWSet: 是否返回读写集
	GetBlockByHash(blockHash string, withRWSet bool) (*pb.BlockInfo, error)

	// 根据交易Id查询区块
	// 参数说明：
	//   - txId: 交易ID
	//   - withRWSet: 是否返回读写集
	GetBlockByTxId(txId string, withRWSet bool) (*pb.BlockInfo, error)

	// 查询最新的配置块
	// 参数说明：
	//   - withRWSet: 是否返回读写集
	GetLastConfigBlock(withRWSet bool) (*pb.BlockInfo, error)

	// 查询节点已部署的所有合约信息，包括：合约名、合约版本、运行环境、交易ID
	GetContractInfo() (*pb.ContractInfo, error)

	// 查询节点加入的链信息，返回ChainId清单
	GetNodeChainList() (*pb.ChainList, error)

	// 查询链信息，包括：当前链最新高度，链节点信息
	GetChainInfo() (*pb.ChainInfo, error)
	// ======================== [(2/5)系统合约接口] ========================
	// =============================== END ===============================

	// ============================== BEGIN ============================
	// ======================== [(3/5)链配置接口] ========================
	// 查询最新链配置
	ChainConfigGet() (*pb.ChainConfig, error)

	// 根据指定区块高度，之前最近的链配置，如果当前区块就是配置块，直接返回当前区块的链配置
	ChainConfigGetByBlockHeight(blockHeight int) (*pb.ChainConfig, error)

	// 查询最新链配置序号Sequence，用于链配置更新
	ChainConfigGetSeq() (int, error)

	// 链配置更新获取Payload签名
	ChainConfigPayloadCollectSign(payloadBytes []byte) ([]byte, error)

	// 链配置更新Payload签名收集&合并
	ChainConfigPayloadMergeSign(signedPayloadBytes [][]byte) ([]byte, error)

	// 发送链配置更新请求
	SendChainConfigUpdateRequest(mergeSignedPayloadBytes []byte) (*pb.TxResponse, error)

	// 以下ChainConfigCreateXXXXXXPayload方法，用于生成链配置待签名payload，在进行多签收集后(需机构Admin权限账号签名)，用于链配置的更新
	// 更新Core模块待签名payload生成
	//   - 若无需修改，请置为-1
	// 参数说明：
	//   - tx_scheduler_timeout：uint，交易调度器从交易池拿到交易后, 进行调度的时间，其值范围为[0, 60]
	//   - tx_scheduler_validate_timeout：uint，交易调度器从区块中拿到交易后, 进行验证的超时时间，其值范围为[0, 60]
	ChainConfigCreateCoreUpdatePayload(txSchedulerTimeout, txSchedulerValidateTimeout int) ([]byte, error)

	// 更新Core模块待签名payload生成
	//   - 若无需修改，请置为-1
	// 参数说明：
	//   - txTimestampVerify：bool，是否需要开启交易时间戳校验
	//   - txTimeout：uint，交易时间戳的过期时间(秒)，其值范围为[600, +∞)
	//   - blockTxCapacity：uint，区块中最大交易数，其值范围为(0, +∞]
	//   - blockSize：uint，区块最大限制，单位MB，其值范围为(0, +∞]
	//   - blockInterval：uint，出块间隔，单位:ms，其值范围为[10, +∞]
	ChainConfigCreateBlockUpdatePayload(txTimestampVerify bool, txTimeout, blockTxCapacity, blockSize, blockInterval int) ([]byte, error)

	// 添加信任组织根证书待签名payload生成
	// 参数说明：
	//   - trustRootOrgId：string，组织Id
	//   - trustRootCrt：string，根证书
	ChainConfigCreateTrustRootAddPayload(trustRootOrgId, trustRootCrt string) ([]byte, error)

	// 更新信任组织根证书待签名payload生成
	// 参数说明：
	//   - trustRootOrgId：string，组织Id
	//   - trustRootCrt：string，根证书
	ChainConfigCreateTrustRootUpdatePayload(trustRootOrgId, trustRootCrt string) ([]byte, error)

	// 删除信任组织根证书待签名payload生成
	// 参数说明：
	//   - trustRootOrgId：string，组织Id
	ChainConfigCreateTrustRootDeletePayload(trustRootOrgId string) ([]byte, error)

	// 添加权限配置待签名payload生成
	// 参数说明：
	//   - permissionResourceName：string，权限名
	//   - principle：*pb.Principle，权限规则
	ChainConfigCreatePermissionAddPayload(permissionResourceName string, principle *pb.Principle) ([]byte, error)

	// 更新权限配置待签名payload生成
	// 参数说明：
	//   - permissionResourceName：string，权限名
	//   - principle：*pb.Principle，权限规则
	ChainConfigCreatePermissionUpdatePayload(permissionResourceName string, principle *pb.Principle) ([]byte, error)

	// 删除权限配置待签名payload生成
	// 参数说明：
	//   - permissionResourceName：string，权限名
	ChainConfigCreatePermissionDeletePayload(permissionResourceName string) ([]byte, error)
	// ======================== [(3/5)链配置接口] ========================
	// =============================== END =============================

	// ============================= BEGIN =============================
	// ======================== [(4/5)证书管理] =========================
	// 用户证书添加
	// 参数说明：
	//   - 在pb.TxResponse.ContractResult.Result字段中返回成功添加的certHash
	CertAdd() (*pb.TxResponse, error)

	// 用户证书删除
	// 参数说明：
	//   - certHashes: 证书Hash列表，多个使用逗号分割
	CertDelete(certHashes string) (*pb.TxResponse, error)

	// 用户证书查询
	// 参数说明：
	//   - certHashes: 证书Hash列表，多个使用逗号分割
	// 返回值说明：
	//   - *pb.CertInfos: 包含证书Hash和证书内容的列表
	CertQuery(certHashes string) (*pb.CertInfos, error)
	// ======================== [(4/5)证书管理] =========================
	// ============================== END ==============================

	// ============================= BEGIN =============================
	// ====================== [(5/5)消息订阅接口] ========================
	// 区块订阅
	// 参数说明：
	//   - startBlock: 订阅起始区块高度，若为-1，表示订阅实时最新区块
	//   - endBlock: 订阅结束区块高度，若为-1，表示订阅实时最新区块
	//   - withRwSet: 是否返回读写集
	SubscribeBlock(ctx context.Context, startBlock, endBlock int64, withRwSet bool) (<-chan interface{}, error)

	// 交易订阅
	// 参数说明：
	//   - startBlock: 订阅起始区块高度，若为-1，表示订阅实时最新区块
	//   - endBlock: 订阅结束区块高度，若为-1，表示订阅实时最新区块
	//   - txType: 订阅交易类型,若为pb.TxType(-1)，表示订阅所有交易类型
	//   - txIds: 订阅txId列表，若为空，表示订阅所有txId
	SubscribeTx(ctx context.Context, startBlock, endBlock int64, txType pb.TxType, txIds []string) (<-chan interface{}, error)

	// 多合一订阅
	// 参数说明：
	//   - txType: 订阅交易类型，目前已支持：区块消息订阅(pb.TxType_SUBSCRIBE_BLOCK_INFO)、交易消息订阅(pb.TxType_SUBSCRIBE_TX_INFO)
	//   - payloadBytes: 消息订阅参数payload
	Subscribe(ctx context.Context, txType pb.TxType, payloadBytes []byte) (<-chan interface{}, error)
	// ====================== [(5/5)消息订阅接口] ========================
	// ============================== END ==============================
}
