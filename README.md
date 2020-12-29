# ChainMaker Go SDK README

[TOC]

## 1 基本概念

- **Node（节点）**：代表一个链节点的基本信息，包括：节点地址、连接数、是否启用`TLS`认证等信息
- **ChainClient（链客户端）**：所有客户端对链节点的操作接口都来自`ChainClient`
- **压缩证书**：可以为`ChainClient`开启证书压缩功能，开启后可以减小交易包大小，提升处理性能

## 2 下载安装

```bash
$ git clone --recursive git@git.code.tencent.com:ChainMaker/chainmaker-sdk-go.git
```

## 3 使用示例

### 3.1 创建节点

```go
// 创建节点
func createNode(nodeAddr string, connCnt int) *NodeConfig {
	node := NewNodeConfig(
		// 节点地址，格式：127.0.0.1:12301
		WithNodeAddr(nodeAddr),
		// 节点连接数
		WithNodeConnCnt(connCnt),
		// 节点是否启用TLS认证
		WithNodeUseTLS(true),
		// 根证书路径，支持多个
		WithNodeCAPaths(caPaths),
		// TLS Hostname
		WithNodeTLSHostName(tlsHostName),
	)

	return node
}
```

### 3.2 创建ChainClient

```go
// 创建ChainClient
func createClient() (*ChainClient, error) {
	if node1 == nil {
		// 创建节点1
		node1 = createNode(nodeAddr1, connCnt1)
	}

	if node2 == nil {
		// 创建节点2
		node2 = createNode(nodeAddr2, connCnt2)
	}

	chainClient, err := NewChainClient(
		// 设置归属组织
		WithChainClientOrgId(chainOrgId),
		// 设置链ID
		WithChainClientChainId(chainId),
		// 设置logger句柄，若不设置，将采用默认日志文件输出日志
		WithChainClientLogger(getDefaultLogger()),
		// 设置客户端用户私钥路径
		WithUserKeyFilePath(userKeyPath),
		// 设置客户端用户证书
		WithUserCrtFilePath(userCrtPath),
		// 添加节点1
		AddChainClientNodeConfig(node1),
		// 添加节点2
		AddChainClientNodeConfig(node2),
		)

	if err != nil {
		return nil, err
	}

	//启用证书压缩（开启证书压缩可以减小交易包大小，提升处理性能）
	err = chainClient.EnableCertHash()
	if err != nil {
		log.Fatal(err)
	}

	return chainClient, nil
}
```

### 3.3 接口调用

> 具体接口调用示例，请参看单元测试用例中的用法。

| 功能     | 单测代码                    |
| -------- | --------------------------- |
| 用户合约 | sdk_user_contract_test.go   |
| 系统合约 | sdk_system_contract_test.go |
| 链配置   | sdk_chain_config_test.go    |
| 证书管理 | sdk_cert_manage_test.go     |
| 消息订阅 | sdk_subscribe_test.go       |

## 4. 接口说明

### 4.1 用户合约接口
#### 4.1.1 创建合约待签名payload生成
**参数说明**
  - contractName: 合约名
  - version: 版本号
  - byteCodePath: 合约路径
  - runtime: 合约运行环境
  - kvs: 合约初始化参数
```go
CreateContractCreatePayload(contractName, version, byteCodePath string, runtime pb.RuntimeType, kvs []*pb.KeyValuePair) ([]byte, error)
```

#### 4.1.2 升级合约待签名payload生成
**参数说明**
  - contractName: 合约名
  - version: 版本号
  - byteCodePath: 合约路径
  - runtime: 合约运行环境
  - kvs: 合约升级参数
```go
CreateContractUpgradePayload(contractName, version, byteCodePath string, runtime pb.RuntimeType, kvs []*pb.KeyValuePair) ([]byte, error)
```

#### 4.1.3 合约管理获取Payload签名
```go
SignContractManagePayload(payloadBytes []byte) ([]byte, error)
```

#### 4.1.4 合约管理Payload签名收集&合并
```go
MergeContractManageSignedPayload(signedPayloadBytes [][]byte) ([]byte, error)
```

#### 4.1.5 发送创建合约请求
**参数说明**
  - multiSignedPayload: 多签结果
  - timeout: 超时时间，单位：s，若传入-1，将使用默认超时时间：10s
  - withSyncResult: 是否同步获取交易执行结果
           当为true时，若成功调用，pb.TxResponse.ContractResult.Result为pb.TransactionInfo
           当为false时，若成功调用，pb.TxResponse.ContractResult.Result为txId
```go
SendContractCreateRequest(mergeSignedPayloadBytes []byte, timeout int64, withSyncResult bool) (*pb.TxResponse, error)
```

#### 4.1.6 发送升级合约请求
**参数说明**
  - multiSignedPayload: 多签结果
  - timeout: 超时时间，单位：s，若传入-1，将使用默认超时时间：10s
  - withSyncResult: 是否同步获取交易执行结果
           当为true时，若成功调用，pb.TxResponse.ContractResult.Result为pb.TransactionInfo
           当为false时，若成功调用，pb.TxResponse.ContractResult.Result为txId
```go
SendContractUpgradeRequest(mergeSignedPayloadBytes []byte, timeout int64, withSyncResult bool) (*pb.TxResponse, error)
```

#### 4.1.7 合约调用
**参数说明**
  - contractName: 合约名称
  - method: 合约方法
  - txId: 交易ID
          格式要求：长度为64bit，字符在a-z0-9
          可为空，若为空字符串，将自动生成txId
  - params: 合约参数
  - timeout: 超时时间，单位：s，若传入-1，将使用默认超时时间：10s
  - withSyncResult: 是否同步获取交易执行结果
           当为true时，若成功调用，pb.TxResponse.ContractResult.Result为pb.TransactionInfo
           当为false时，若成功调用，pb.TxResponse.ContractResult.Result为txId
```go
InvokeContract(contractName, method, txId string, params map[string]string, timeout int64, withSyncResult bool) (*pb.TxResponse, error)
```

#### 4.1.8 合约查询接口调用
**参数说明**
  - contractName: 合约名称
  - method: 合约方法
  - params: 合约参数
  - timeout: 超时时间，单位：s，若传入-1，将使用默认超时时间：10s
```go
QueryContract(contractName, method string, params map[string]string, timeout int64) (*pb.TxResponse, error)
```

### 4.2 系统合约接口
#### 4.2.1 根据交易Id查询交易
**参数说明**
  - txId: 交易ID
```go
GetTxByTxId(txId string) (*pb.TransactionInfo, error)
```

#### 4.2.2 根据区块高度查询区块
**参数说明**
  - blockHeight: 指定区块高度，若为-1，将返回最新区块
  - withRWSet: 是否返回读写集
```go
GetBlockByHeight(blockHeight int64, withRWSet bool) (*pb.BlockInfo, error)
```

#### 4.2.3 根据区块哈希查询区块
**参数说明**
  - blockHash: 指定区块Hash
  - withRWSet: 是否返回读写集
```go
GetBlockByHash(blockHash string, withRWSet bool) (*pb.BlockInfo, error)
```

#### 4.2.4 根据交易Id查询区块
**参数说明**
  - txId: 交易ID
  - withRWSet: 是否返回读写集
```go
GetBlockByTxId(txId string, withRWSet bool) (*pb.BlockInfo, error)
```

#### 4.2.5 查询最新的配置块
**参数说明**
  - withRWSet: 是否返回读写集
```go
GetLastConfigBlock(withRWSet bool) (*pb.BlockInfo, error)
```

#### 4.2.6 查询节点已部署的所有合约信息
   - 包括：合约名、合约版本、运行环境、交易ID
```go
GetContractInfo() (*pb.ContractInfo, error)
```

#### 4.2.7 查询节点加入的链信息
   - 返回ChainId清单
```go
GetNodeChainList() (*pb.ChainList, error)
```

#### 4.2.8 查询链信息
  - 包括：当前链最新高度，链节点信息
```go
GetChainInfo() (*pb.ChainInfo, error)
```

### 4.3 链配置接口
#### 4.3.1 查询最新链配置
```go
GetChainConfig() (*pb.ChainConfig, error)
```

#### 4.3.2 根据指定区块高度查询最近链配置
  - 如果当前区块就是配置块，直接返回当前区块的链配置
```go
GetChainConfigByBlockHeight(blockHeight int) (*pb.ChainConfig, error)
```

#### 4.3.3 查询最新链配置序号Sequence
  - 用于链配置更新
```go
GetChainConfigSequence() (int, error)
```

#### 4.3.4 链配置更新获取Payload签名
```go
SignChainConfigPayload(payloadBytes []byte) ([]byte, error)
```

#### 4.3.5 链配置更新Payload签名收集&合并
```go
MergeChainConfigSignedPayload(signedPayloadBytes [][]byte) ([]byte, error)
```

#### 4.3.6 发送链配置更新请求
```go
SendChainConfigUpdateRequest(mergeSignedPayloadBytes []byte) (*pb.TxResponse, error)
```

> 以下CreateChainConfigXXXXXXPayload方法，用于生成链配置待签名payload，在进行多签收集后(需机构Admin权限账号签名)，用于链配置的更新

#### 4.3.7 更新Core模块待签名payload生成
**参数说明**
  - txSchedulerTimeout: 交易调度器从交易池拿到交易后, 进行调度的时间，其值范围为[0, 60]，若无需修改，请置为-1
  - txSchedulerValidateTimeout: 交易调度器从区块中拿到交易后, 进行验证的超时时间，其值范围为[0, 60]，若无需修改，请置为-1
```go
CreateChainConfigCoreUpdatePayload(txSchedulerTimeout, txSchedulerValidateTimeout int) ([]byte, error)
```

#### 4.3.8 更新Core模块待签名payload生成
**参数说明**
  - txTimestampVerify: 是否需要开启交易时间戳校验
  - (以下参数，若无需修改，请置为-1)
  - txTimeout: 交易时间戳的过期时间(秒)，其值范围为[600, +∞)
  - blockTxCapacity: 区块中最大交易数，其值范围为(0, +∞]
  - blockSize: 区块最大限制，单位MB，其值范围为(0, +∞]
  - blockInterval: 出块间隔，单位:ms，其值范围为[10, +∞]
```go
CreateChainConfigBlockUpdatePayload(txTimestampVerify bool, txTimeout, blockTxCapacity, blockSize, blockInterval int) ([]byte, error)
```

#### 4.3.9 添加信任组织根证书待签名payload生成
**参数说明**
  - trustRootOrgId: 组织Id
  - trustRootCrt: 根证书
```go
CreateChainConfigTrustRootAddPayload(trustRootOrgId, trustRootCrt string) ([]byte, error)
```

#### 4.3.10 更新信任组织根证书待签名payload生成
**参数说明**
  - trustRootOrgId: 组织Id
  - trustRootCrt: 根证书
```go
CreateChainConfigTrustRootUpdatePayload(trustRootOrgId, trustRootCrt string) ([]byte, error)
```

#### 4.3.11 删除信任组织根证书待签名payload生成
**参数说明**
  - trustRootOrgId: 组织Id
```go
CreateChainConfigTrustRootDeletePayload(trustRootOrgId string) ([]byte, error)
```

#### 4.3.12 添加权限配置待签名payload生成
**参数说明**
  - permissionResourceName: 权限名
  - principle: 权限规则
```go
CreateChainConfigPermissionAddPayload(permissionResourceName string, principle *pb.Principle) ([]byte, error)
```

#### 4.3.13 更新权限配置待签名payload生成
**参数说明**
  - permissionResourceName: 权限名
  - principle: 权限规则
```go
CreateChainConfigPermissionUpdatePayload(permissionResourceName string, principle *pb.Principle) ([]byte, error)
```

#### 4.3.14 删除权限配置待签名payload生成
**参数说明**
  - permissionResourceName: 权限名
```go
CreateChainConfigPermissionDeletePayload(permissionResourceName string) ([]byte, error)
```

#### 4.3.15 添加共识节点地址待签名payload生成
**参数说明**
  - nodeOrgId: 节点组织Id
  - nodeAddresses: 节点地址
```go
CreateChainConfigConsensusNodeAddrAddPayload(nodeOrgId string, nodeAddresses []string) ([]byte, error)
```

#### 4.3.16 更新共识节点地址待签名payload生成
**参数说明**
  - nodeOrgId: 节点组织Id
  - nodeOldAddress: 节点原地址
  - nodeNewAddress: 节点新地址
```go
CreateChainConfigConsensusNodeAddrUpdatePayload(nodeOrgId, nodeOldAddress, nodeNewAddress string) ([]byte, error)
```

#### 4.3.17 删除共识节点地址待签名payload生成
**参数说明**
  - nodeOrgId: 节点组织Id
  - nodeAddress: 节点地址
```go
CreateChainConfigConsensusNodeAddrDeletePayload(nodeOrgId, nodeAddress string) ([]byte, error)
```

#### 4.3.18 添加共识节点待签名payload生成
**参数说明**
  - nodeOrgId: 节点组织Id
  - nodeAddresses: 节点地址
```go
CreateChainConfigConsensusNodeOrgAddPayload(nodeOrgId string, nodeAddresses []string) ([]byte, error)
```

#### 4.3.19 更新共识节点待签名payload生成
**参数说明**
  - nodeOrgId: 节点组织Id
  - nodeAddresses: 节点地址
```go
CreateChainConfigConsensusNodeOrgUpdatePayload(nodeOrgId string, nodeAddresses []string) ([]byte, error)
```

#### 4.3.20 删除共识节点待签名payload生成
**参数说明**
  - nodeOrgId: 节点组织Id
```go
CreateChainConfigConsensusNodeOrgDeletePayload(nodeOrgId string) ([]byte, error)
```

#### 4.3.21 添加共识扩展字段待签名payload生成
**参数说明**
  - kvs: 字段key、value对
```go
CreateChainConfigConsensusExtAddPayload(kvs []*pb.KeyValuePair) ([]byte, error)
```

#### 4.3.22 添加共识扩展字段待签名payload生成
**参数说明**
  - kvs: 字段key、value对
```go
CreateChainConfigConsensusExtUpdatePayload(kvs []*pb.KeyValuePair) ([]byte, error)
```

#### 4.3.23 添加共识扩展字段待签名payload生成
**参数说明**
  - keys: 待删除字段
```go
CreateChainConfigConsensusExtDeletePayload(keys []string) ([]byte, error)
```

### 4.4 证书管理接口
#### 4.4.1 用户证书添加
**参数说明**
  - 在pb.TxResponse.ContractResult.Result字段中返回成功添加的certHash
```go
AddCert() (*pb.TxResponse, error)
```

#### 4.4.2 用户证书删除
**参数说明**
  - certHashes: 证书Hash列表
```go
DeleteCert(certHashes []string) (*pb.TxResponse, error)
```

#### 4.4.3 用户证书查询
**参数说明**
  - certHashes: 证书Hash列表
返回值说明：
  - *pb.CertInfos: 包含证书Hash和证书内容的列表
```go
QueryCert(certHashes []string) (*pb.CertInfos, error)
```

#### 4.4.4 获取用户证书哈希
```go
GetCertHash() ([]byte, error)
```

### 4.5 消息订阅接口
#### 4.5.1 区块订阅
**参数说明**
  - startBlock: 订阅起始区块高度，若为-1，表示订阅实时最新区块
  - endBlock: 订阅结束区块高度，若为-1，表示订阅实时最新区块
  - withRwSet: 是否返回读写集
```go
SubscribeBlock(ctx context.Context, startBlock, endBlock int64, withRwSet bool) (<-chan interface{}, error)
```

#### 4.5.2 交易订阅
**参数说明**
  - startBlock: 订阅起始区块高度，若为-1，表示订阅实时最新区块
  - endBlock: 订阅结束区块高度，若为-1，表示订阅实时最新区块
  - txType: 订阅交易类型,若为pb.TxType(-1)，表示订阅所有交易类型
  - txIds: 订阅txId列表，若为空，表示订阅所有txId
```go
SubscribeTx(ctx context.Context, startBlock, endBlock int64, txType pb.TxType, txIds []string) (<-chan interface{}, error)
```

#### 4.5.3 多合一订阅
**参数说明**
  - txType: 订阅交易类型，目前已支持：区块消息订阅(pb.TxType_SUBSCRIBE_BLOCK_INFO)、交易消息订阅(pb.TxType_SUBSCRIBE_TX_INFO)
  - payloadBytes: 消息订阅参数payload
```go
Subscribe(ctx context.Context, txType pb.TxType, payloadBytes []byte) (<-chan interface{}, error)
```

### 4.6 证书压缩

*开启证书压缩可以减小交易包大小，提升处理性能*

#### 4.6.1 启用压缩证书功能

```go
EnableCertHash() error
```

#### 4.6.2 停用压缩证书功能

```go
DisableCertHash() error
```

### 4.7 管理类接口

#### 4.7.1 SDK停止接口

*关闭连接池连接，释放资源*

```go
Stop() error
```


