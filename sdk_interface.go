/*
Copyright (C) BABEC. All rights reserved.
Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.

SPDX-License-Identifier: Apache-2.0
*/

package chainmaker_sdk_go

import (
	"chainmaker.org/chainmaker-sdk-go/pb/protogo/accesscontrol"
	"chainmaker.org/chainmaker-sdk-go/pb/protogo/common"
	"chainmaker.org/chainmaker-sdk-go/pb/protogo/config"
	"chainmaker.org/chainmaker-sdk-go/pb/protogo/discovery"
	"context"
)

// # ChainMaker Go SDK 接口说明
type SDKInterface interface {
	// ## 1 用户合约接口
	// ### 1.1 创建合约待签名payload生成
	// **参数说明**
	//   - contractName: 合约名
	//   - version: 版本号
	//   - byteCode: 支持传入合约二进制文件路径或Base64编码的二进制内容
	//   - runtime: 合约运行环境
	//   - kvs: 合约初始化参数
	// ```go
	CreateContractCreatePayload(contractName, version, byteCode string, runtime common.RuntimeType, kvs []*common.KeyValuePair) ([]byte, error)
	// ```

	// ### 1.2 升级合约待签名payload生成
	// **参数说明**
	//   - contractName: 合约名
	//   - version: 版本号
	//   - byteCode: 支持传入合约二进制文件路径或Base64编码的二进制内容
	//   - runtime: 合约运行环境
	//   - kvs: 合约升级参数
	// ```go
	CreateContractUpgradePayload(contractName, version, byteCode string, runtime common.RuntimeType, kvs []*common.KeyValuePair) ([]byte, error)
	// ```

	// ### 1.3 冻结合约payload生成
	// **参数说明**
	//   - contractName: 合约名
	// ```go
	CreateContractFreezePayload(contractName string) ([]byte, error)
	// ```

	// ### 1.4 解冻合约payload生成
	// **参数说明**
	//   - contractName: 合约名
	// ```go
	CreateContractUnfreezePayload(contractName string) ([]byte, error)
	// ```

	// ### 1.5 吊销合约payload生成
	// **参数说明**
	//   - contractName: 合约名
	// ```go
	CreateContractRevokePayload(contractName string) ([]byte, error)
	// ```

	// ### 1.6 合约管理获取Payload签名
	// **参数说明**
	//   - payloadBytes: 待签名payload
	// ```go
	SignContractManagePayload(payloadBytes []byte) ([]byte, error)
	// ```

	// ### 1.7 合约管理Payload签名收集&合并
	// **参数说明**
	//   - signedPayloadBytes: 已签名payload列表
	// ```go
	MergeContractManageSignedPayload(signedPayloadBytes [][]byte) ([]byte, error)
	// ```

	// ### 1.8 发送合约管理请求（创建、更新、冻结、解冻、吊销）
	// **参数说明**
	//   - multiSignedPayload: 多签结果
	//   - timeout: 超时时间，单位：s，若传入-1，将使用默认超时时间：10s
	//   - withSyncResult: 是否同步获取交易执行结果
	//            当为true时，若成功调用，common.TxResponse.ContractResult.Result为common.TransactionInfo
	//            当为false时，若成功调用，common.TxResponse.ContractResult.Result为txId
	// ```go
	SendContractManageRequest(mergeSignedPayloadBytes []byte, timeout int64, withSyncResult bool) (*common.TxResponse, error)
	// ```

	// ### 1.9 合约调用
	// **参数说明**
	//   - contractName: 合约名称
	//   - method: 合约方法
	//   - txId: 交易ID
	//           格式要求：长度为64bit，字符在a-z0-9
	//           可为空，若为空字符串，将自动生成txId
	//   - params: 合约参数
	//   - timeout: 超时时间，单位：s，若传入-1，将使用默认超时时间：10s
	//   - withSyncResult: 是否同步获取交易执行结果
	//            当为true时，若成功调用，common.TxResponse.ContractResult.Result为common.TransactionInfo
	//            当为false时，若成功调用，common.TxResponse.ContractResult.Result为txId
	// ```go
	InvokeContract(contractName, method, txId string, params map[string]string, timeout int64, withSyncResult bool) (*common.TxResponse, error)
	// ```

	// ### 1.10 合约查询接口调用
	// **参数说明**
	//   - contractName: 合约名称
	//   - method: 合约方法
	//   - params: 合约参数
	//   - timeout: 超时时间，单位：s，若传入-1，将使用默认超时时间：10s
	// ```go
	QueryContract(contractName, method string, params map[string]string, timeout int64) (*common.TxResponse, error)
	// ```

	// ## 2 系统合约接口
	// ### 2.1 根据交易Id查询交易
	// **参数说明**
	//   - txId: 交易ID
	// ```go
	GetTxByTxId(txId string) (*common.TransactionInfo, error)
	// ```

	// ### 2.2 根据区块高度查询区块
	// **参数说明**
	//   - blockHeight: 指定区块高度，若为-1，将返回最新区块
	//   - withRWSet: 是否返回读写集
	// ```go
	GetBlockByHeight(blockHeight int64, withRWSet bool) (*common.BlockInfo, error)
	// ```

	// ### 2.3 根据区块哈希查询区块
	// **参数说明**
	//   - blockHash: 指定区块Hash
	//   - withRWSet: 是否返回读写集
	// ```go
	GetBlockByHash(blockHash string, withRWSet bool) (*common.BlockInfo, error)
	// ```

	// ### 2.4 根据交易Id查询区块
	// **参数说明**
	//   - txId: 交易ID
	//   - withRWSet: 是否返回读写集
	// ```go
	GetBlockByTxId(txId string, withRWSet bool) (*common.BlockInfo, error)
	// ```

	// ### 2.5 查询最新的配置块
	// **参数说明**
	//   - withRWSet: 是否返回读写集
	// ```go
	GetLastConfigBlock(withRWSet bool) (*common.BlockInfo, error)
	// ```

	// ### 2.6 查询节点加入的链信息
	//    - 返回ChainId清单
	// ```go
	GetNodeChainList() (*discovery.ChainList, error)
	// ```

	// ### 2.7 查询链信息
	//   - 包括：当前链最新高度，链节点信息
	// ```go
	GetChainInfo() (*discovery.ChainInfo, error)
	// ```

	// ## 3 链配置接口
	// ### 3.1 查询最新链配置
	// ```go
	GetChainConfig() (*config.ChainConfig, error)
	// ```

	// ### 3.2 根据指定区块高度查询最近链配置
	//   - 如果当前区块就是配置块，直接返回当前区块的链配置
	// ```go
	GetChainConfigByBlockHeight(blockHeight int) (*config.ChainConfig, error)
	// ```

	// ### 3.3 查询最新链配置序号Sequence
	//   - 用于链配置更新
	// ```go
	GetChainConfigSequence() (int, error)
	// ```

	// ### 3.4 链配置更新获取Payload签名
	// ```go
	SignChainConfigPayload(payloadBytes []byte) ([]byte, error)
	// ```

	// ### 3.5 链配置更新Payload签名收集&合并
	// ```go
	MergeChainConfigSignedPayload(signedPayloadBytes [][]byte) ([]byte, error)
	// ```

	// ### 3.6 发送链配置更新请求
	// ```go
	SendChainConfigUpdateRequest(mergeSignedPayloadBytes []byte) (*common.TxResponse, error)
	// ```

	// > 以下CreateChainConfigXXXXXXPayload方法，用于生成链配置待签名payload，在进行多签收集后(需机构Admin权限账号签名)，用于链配置的更新

	// ### 3.7 更新Core模块待签名payload生成
	// **参数说明**
	//   - txSchedulerTimeout: 交易调度器从交易池拿到交易后, 进行调度的时间，其值范围为[0, 60]，若无需修改，请置为-1
	//   - txSchedulerValidateTimeout: 交易调度器从区块中拿到交易后, 进行验证的超时时间，其值范围为[0, 60]，若无需修改，请置为-1
	// ```go
	CreateChainConfigCoreUpdatePayload(txSchedulerTimeout, txSchedulerValidateTimeout int) ([]byte, error)
	// ```

	// ### 3.8 更新Core模块待签名payload生成
	// **参数说明**
	//   - txTimestampVerify: 是否需要开启交易时间戳校验
	//   - (以下参数，若无需修改，请置为-1)
	//   - txTimeout: 交易时间戳的过期时间(秒)，其值范围为[600, +∞)
	//   - blockTxCapacity: 区块中最大交易数，其值范围为(0, +∞]
	//   - blockSize: 区块最大限制，单位MB，其值范围为(0, +∞]
	//   - blockInterval: 出块间隔，单位:ms，其值范围为[10, +∞]
	// ```go
	CreateChainConfigBlockUpdatePayload(txTimestampVerify bool, txTimeout, blockTxCapacity, blockSize, blockInterval int) ([]byte, error)
	// ```

	// ### 3.9 添加信任组织根证书待签名payload生成
	// **参数说明**
	//   - trustRootOrgId: 组织Id
	//   - trustRootCrt: 根证书
	// ```go
	CreateChainConfigTrustRootAddPayload(trustRootOrgId, trustRootCrt string) ([]byte, error)
	// ```

	// ### 3.10 更新信任组织根证书待签名payload生成
	// **参数说明**
	//   - trustRootOrgId: 组织Id
	//   - trustRootCrt: 根证书
	// ```go
	CreateChainConfigTrustRootUpdatePayload(trustRootOrgId, trustRootCrt string) ([]byte, error)
	// ```

	// ### 3.11 删除信任组织根证书待签名payload生成
	// **参数说明**
	//   - trustRootOrgId: 组织Id
	// ```go
	CreateChainConfigTrustRootDeletePayload(trustRootOrgId string) ([]byte, error)
	// ```

	// ### 3.12 添加权限配置待签名payload生成
	// **参数说明**
	//   - permissionResourceName: 权限名
	//   - policy: 权限规则
	// ```go
	CreateChainConfigPermissionAddPayload(permissionResourceName string, policy *accesscontrol.Policy) ([]byte, error)
	// ```

	// ### 3.13 更新权限配置待签名payload生成
	// **参数说明**
	//   - permissionResourceName: 权限名
	//   - policy: 权限规则
	// ```go
	CreateChainConfigPermissionUpdatePayload(permissionResourceName string, policy *accesscontrol.Policy) ([]byte, error)
	// ```

	// ### 3.14 删除权限配置待签名payload生成
	// **参数说明**
	//   - permissionResourceName: 权限名
	// ```go
	CreateChainConfigPermissionDeletePayload(permissionResourceName string) ([]byte, error)
	// ```

	// ### 3.15 添加共识节点地址待签名payload生成
	// **参数说明**
	//   - nodeOrgId: 节点组织Id
	//   - nodeAddresses: 节点地址
	// ```go
	CreateChainConfigConsensusNodeAddrAddPayload(nodeOrgId string, nodeAddresses []string) ([]byte, error)
	// ```

	// ### 3.16 更新共识节点地址待签名payload生成
	// **参数说明**
	//   - nodeOrgId: 节点组织Id
	//   - nodeOldAddress: 节点原地址
	//   - nodeNewAddress: 节点新地址
	// ```go
	CreateChainConfigConsensusNodeAddrUpdatePayload(nodeOrgId, nodeOldAddress, nodeNewAddress string) ([]byte, error)
	// ```

	// ### 3.17 删除共识节点地址待签名payload生成
	// **参数说明**
	//   - nodeOrgId: 节点组织Id
	//   - nodeAddress: 节点地址
	// ```go
	CreateChainConfigConsensusNodeAddrDeletePayload(nodeOrgId, nodeAddress string) ([]byte, error)
	// ```

	// ### 3.18 添加共识节点待签名payload生成
	// **参数说明**
	//   - nodeOrgId: 节点组织Id
	//   - nodeAddresses: 节点地址
	// ```go
	CreateChainConfigConsensusNodeOrgAddPayload(nodeOrgId string, nodeAddresses []string) ([]byte, error)
	// ```

	// ### 3.19 更新共识节点待签名payload生成
	// **参数说明**
	//   - nodeOrgId: 节点组织Id
	//   - nodeAddresses: 节点地址
	// ```go
	CreateChainConfigConsensusNodeOrgUpdatePayload(nodeOrgId string, nodeAddresses []string) ([]byte, error)
	// ```

	// ### 3.20 删除共识节点待签名payload生成
	// **参数说明**
	//   - nodeOrgId: 节点组织Id
	// ```go
	CreateChainConfigConsensusNodeOrgDeletePayload(nodeOrgId string) ([]byte, error)
	// ```

	// ### 3.21 添加共识扩展字段待签名payload生成
	// **参数说明**
	//   - kvs: 字段key、value对
	// ```go
	CreateChainConfigConsensusExtAddPayload(kvs []*common.KeyValuePair) ([]byte, error)
	// ```

	// ### 3.22 添加共识扩展字段待签名payload生成
	// **参数说明**
	//   - kvs: 字段key、value对
	// ```go
	CreateChainConfigConsensusExtUpdatePayload(kvs []*common.KeyValuePair) ([]byte, error)
	// ```

	// ### 3.23 添加共识扩展字段待签名payload生成
	// **参数说明**
	//   - keys: 待删除字段
	// ```go
	CreateChainConfigConsensusExtDeletePayload(keys []string) ([]byte, error)
	// ```

	// ## 4 证书管理接口
	// ### 4.1 用户证书添加
	// **参数说明**
	//   - 在common.TxResponse.ContractResult.Result字段中返回成功添加的certHash
	// ```go
	AddCert() (*common.TxResponse, error)
	// ```

	// ### 4.2 用户证书删除
	// **参数说明**
	//   - certHashes: 证书Hash列表
	// ```go
	DeleteCert(certHashes []string) (*common.TxResponse, error)
	// ```

	// ### 4.3 用户证书查询
	// **参数说明**
	//   - certHashes: 证书Hash列表
	// 返回值说明：
	//   - *common.CertInfos: 包含证书Hash和证书内容的列表
	// ```go
	QueryCert(certHashes []string) (*common.CertInfos, error)
	// ```

	// ### 4.4 获取用户证书哈希
	// ```go
	GetCertHash() ([]byte, error)
	// ```

	// ### 4.5 生成证书管理操作Payload（三合一接口）
	// **参数说明**
	//   - method: CERTS_FROZEN(证书冻结)/CERTS_UNFROZEN(证书解冻)/CERTS_REVOCATION(证书吊销)
	//   - kvs: 证书管理操作参数
	// ```go
	CreateCertManagePayload(method string, kvs []*common.KeyValuePair) ([]byte, error)
	// ```

	// ### 4.6 生成证书冻结操作Payload
	// **参数说明**
	//   - certs: X509证书列表
	// ```go
	CreateCertManageFrozenPayload(certs []string) ([]byte, error)
	// ```

	// ### 4.7 生成证书解冻操作Payload
	// **参数说明**
	//   - certs: X509证书列表
	// ```go
	CreateCertManageUnfrozenPayload(certs []string) ([]byte, error)
	// ```

	// ### 4.8 生成证书吊销操作Payload
	// **参数说明**
	//   - certs: X509证书列表
	// ```go
	CreateCertManageRevocationPayload(certCrl string) ([]byte, error)
	// ```

	// ### 4.9 待签payload签名
	//  *一般需要使用具有管理员权限账号进行签名*
	// **参数说明**
	//   - payloadBytes: 待签名payload
	// ```go
	SignCertManagePayload(payloadBytes []byte) ([]byte, error)
	// ```

	// ### 4.10 证书管理Payload签名收集&合并
	// **参数说明**
	//   - signedPayloadBytes: 已签名payload列表
	// ```go
	MergeCertManageSignedPayload(signedPayloadBytes [][]byte) ([]byte, error)
	// ```

	// ### 4.11 发送证书管理请求（证书冻结、解冻、吊销）
	// **参数说明**
	//   - multiSignedPayload: 多签结果
	//   - timeout: 超时时间，单位：s，若传入-1，将使用默认超时时间：10s
	//   - withSyncResult: 是否同步获取交易执行结果
	//            当为true时，若成功调用，common.TxResponse.ContractResult.Result为common.TransactionInfo
	//            当为false时，若成功调用，common.TxResponse.ContractResult.Result为txId
	// ```go
	SendCertManageRequest(mergeSignedPayloadBytes []byte, timeout int64, withSyncResult bool) (*common.TxResponse, error)
	// ```

	// ## 5 在线多签接口
	// ### 5.1 待签payload签名
	//  *一般需要使用具有管理员权限账号进行签名*
	// **参数说明**
	//   - payloadBytes: 待签名payload
	// ```go
	SignMultiSignPayload(payloadBytes []byte) (*common.EndorsementEntry, error)
	// ```

	// ### 5.2 多签请求
	// **参数说明**
	//   - txType: 多签payload交易类型
	//   - payloadBytes: 待签名payload
	//   - endorsementEntry: 签名收集信息
	//   - deadlineBlockHeight: 过期的区块高度，若设置为0，表示永不过期
	//   - timeout: 超时时间，单位：s，若传入-1，将使用默认超时时间：10s
	// **返回值**
	//   若成功调用，common.TxResponse.ContractResult.Result为txId
	// ```go
	SendMultiSignReq(txType common.TxType, payloadBytes []byte, endorsementEntry *common.EndorsementEntry, deadlineBlockHeight int,
		timeout int64) (*common.TxResponse, error)
	// ```

	// ### 5.3 多签投票
	// **参数说明**
	//   - voteStatus: 投票状态（赞成、反对）
	//   - multiSignReqTxId: 多签请求交易ID(txId或payloadHash至少填其一，txId优先)
	//   - payloadHash: 待多签payload hash(txId或payloadHash至少填其一，txId优先)
	//   - payloadBytes: 待签名payload
	//   - endorsementEntry: 签名收集信息
	//   - timeout: 超时时间，单位：s，若传入-1，将使用默认超时时间：10s
	// **返回值**
	//   若成功调用，common.TxResponse.ContractResult.Result为txId
	// ```go
	SendMultiSignVote(voteStatus common.VoteStatus, multiSignReqTxId, payloadHash string,
		endorsementEntry *common.EndorsementEntry, timeout int64) (*common.TxResponse, error)
	// ```

	// ### 5.4 投票查询
	// **参数说明**
	//   - multiSignReqTxId: 多签请求交易ID(txId或payloadHash至少填其一，txId优先)
	//   - payloadHash: 待多签payload hash(txId或payloadHash至少填其一，txId优先)
	// ```go
	QueryMultiSignResult(multiSignReqTxId, payloadHash string) (*common.TxResponse, error)
	// ```

	// ## 6 消息订阅接口
	// ### 6.1 区块订阅
	// **参数说明**
	//   - startBlock: 订阅起始区块高度，若为-1，表示订阅实时最新区块
	//   - endBlock: 订阅结束区块高度，若为-1，表示订阅实时最新区块
	//   - withRwSet: 是否返回读写集
	// ```go
	SubscribeBlock(ctx context.Context, startBlock, endBlock int64, withRwSet bool) (<-chan interface{}, error)
	// ```

	// ### 6.2 交易订阅
	// **参数说明**
	//   - startBlock: 订阅起始区块高度，若为-1，表示订阅实时最新区块
	//   - endBlock: 订阅结束区块高度，若为-1，表示订阅实时最新区块
	//   - txType: 订阅交易类型,若为common.TxType(-1)，表示订阅所有交易类型
	//   - txIds: 订阅txId列表，若为空，表示订阅所有txId
	// ```go
	SubscribeTx(ctx context.Context, startBlock, endBlock int64, txType common.TxType, txIds []string) (<-chan interface{}, error)
	// ```

	// ### 6.3 多合一订阅
	// **参数说明**
	//   - txType: 订阅交易类型，目前已支持：区块消息订阅(common.TxType_SUBSCRIBE_BLOCK_INFO)、交易消息订阅(common.TxType_SUBSCRIBE_TX_INFO)
	//   - payloadBytes: 消息订阅参数payload
	// ```go
	Subscribe(ctx context.Context, txType common.TxType, payloadBytes []byte) (<-chan interface{}, error)
	// ```

	// ## 7 证书压缩
	// *开启证书压缩可以减小交易包大小，提升处理性能*
	// ### 7.1 启用压缩证书功能
	// ```go
	EnableCertHash() error
	// ```

	// ### 7.2 停用压缩证书功能
	// ```go
	DisableCertHash() error
	// ```

	// ## 9 编解码类
	// ### 9.1 将EasyCodec编码解码成map
	// ```go
	EasyCodecBytesToParamsMap(data []byte) map[string]string
	// ```

	// ## 10 系统类接口
	// ### 10.1 SDK停止接口
	// *关闭连接池连接，释放资源*
	// ```go
	Stop() error
	// ```

	// ### 10.2 获取链版本
	// ```go
	GetChainMakerServerVersion() (string, error)
	// ```

	// ## 11 系统隐私合约类接口
	// ### 11.1 证书上链及验证
	// ```go
	SaveCert(enclaveCert, enclaveId, txId string, withSyncResult bool, timeout int64) (*common.TxResponse, error)
	//```

	// ### 11.2 隐私目录上链
	// ```go
	SaveDir(orderId, txId string, privateDir *common.StrSlice, withSyncResult bool, timeout int64) (*common.TxResponse, error)
	//```

	// ### 11.3 隐私合约代码查询
	// ```go
	GetContract(contractName, codeHash string) (*common.PrivateGetContract, error)
	//```

	// ### 11.4 隐私计算结果上链
	// ```go
	SaveData(contractName string, contractVersion string, codeHash []byte, reportHash []byte, result *common.ContractResult, txId string, rwSet *common.TxRWSet, reportSign []byte, events *common.StrSlice, userCert []byte, clientSign []byte, orgId string, withSyncResult bool, timeout int64) (*common.TxResponse, error)
	//```

	// ### 11.5 隐私计算结果查询
	// ```go
	GetData(contractName, key string) ([]byte, error)
	//```

	// ### 11.6 隐私合约代码上链
	// ```go
	SaveContract(codeBytes []byte, codeHash, contractName, version, txId string, withSyncResult bool, timeout int64) (*common.TxResponse, error)
	//```

	// ### 11.7 enclave通过⽹关调⽤
	// ```go
	SaveQuote(enclaveId, quoteId, quote, sign, txId string, withSyncResult bool, timeout int64) (*common.TxResponse, error)
	//```

	// ### 11.8 隐私计算证书查询
	// ```go
	GetCert(enclaveId string) ([]byte, error)
	// ```

	// ### 11.9 隐私计算隐私目录查询
	// ```go
	GetDir(orderId string) ([]byte, error)
	// ```

	// ### 11.10 隐私计算隐私目录查询
	// ```go
	GetQuote(quoteId string) ([]byte, error)
	// ```

	// ###  11.11 隐私计算调用者权限验证
	// ```go
	CheckCallerCertAuth(userCert, clientSign, payload, orgId string) (*common.TxResponse, error)
	// ```
}
