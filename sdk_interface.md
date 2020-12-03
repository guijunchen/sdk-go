# ChainMaker Go SDK 接口说明
## 1 用户合约接口
### 1.1 合约创建
**参数说明**
  - txId: 交易ID
          格式要求：长度为64bit，字符在a-z0-9
          可为空，若为空字符串，将自动生成，在pb.TxResponse.ContractResult.Result字段中返回该自动生成的txId
  - multiSignedPayload: 经多签后的payload数据
```go
ContractCreate(txId string, multiSignedPayload []byte) (*pb.TxResponse, error)
```

### 1.2 合约升级
**参数说明**
  - txId: 交易ID
          格式要求：长度为64bit，字符在a-z0-9
          可为空，若为空字符串，将自动生成，在pb.TxResponse.ContractResult.Result字段中返回该自动生成的txId
  - multiSignedPayload: 经多签后的payload数据
```go
ContractUpgrade(txId string, multiSignedPayload []byte) (*pb.TxResponse, error)
```

### 1.3 合约调用
**参数说明**
  - contractName: 合约名称
  - method: 合约方法
  - txId: 交易ID
          格式要求：长度为64bit，字符在a-z0-9
          可为空，若为空字符串，将自动生成，在pb.TxResponse.ContractResult.Result字段中返回该自动生成的txId
  - params: 合约参数
```go
ContractInvoke(contractName, method, txId string, params map[string]string) (*pb.TxResponse, error)
```

### 1.4 合约查询接口调用
**参数说明**
  - contractName: 合约名称
  - method: 合约方法
  - params: 合约参数
```go
ContractQuery(contractName, method string, params map[string]string) (*pb.TxResponse, error)
```

## 2 系统合约接口
### 2.1 根据交易Id查询交易
**参数说明**
  - txId: 交易ID
```go
GetTxByTxId(txId string) (*pb.TransactionInfo, error)
```

### 2.2 根据区块高度查询区块
**参数说明**
  - blockHeight: 指定区块高度，若为-1，将返回最新区块
  - withRWSet: 是否返回读写集
```go
GetBlockByHeight(blockHeight int64, withRWSet bool) (*pb.BlockInfo, error)
```

### 2.3 根据区块哈希查询区块
**参数说明**
  - blockHash: 指定区块Hash
  - withRWSet: 是否返回读写集
```go
GetBlockByHash(blockHash string, withRWSet bool) (*pb.BlockInfo, error)
```

### 2.4 根据交易Id查询区块
**参数说明**
  - txId: 交易ID
  - withRWSet: 是否返回读写集
```go
GetBlockByTxId(txId string, withRWSet bool) (*pb.BlockInfo, error)
```

### 2.5 查询最新的配置块
**参数说明**
  - withRWSet: 是否返回读写集
```go
GetLastConfigBlock(withRWSet bool) (*pb.BlockInfo, error)
```

### 2.6 查询节点已部署的所有合约信息
   - 包括：合约名、合约版本、运行环境、交易ID
```go
GetContractInfo() (*pb.ContractInfo, error)
```

### 2.7 查询节点加入的链信息
   - 返回ChainId清单
```go
GetNodeChainList() (*pb.ChainList, error)
```

### 2.8 查询链信息
  - 包括：当前链最新高度，链节点信息
```go
GetChainInfo() (*pb.ChainInfo, error)
```

## 3 链配置接口
### 3.1 查询最新链配置
```go
ChainConfigGet() (*pb.ChainConfig, error)
```

### 3.2 根据指定区块高度查询最近链配置
  - 如果当前区块就是配置块，直接返回当前区块的链配置
```go
ChainConfigGetByBlockHeight(blockHeight int) (*pb.ChainConfig, error)
```

### 3.3 查询最新链配置序号Sequence
  - 用于链配置更新
```go
ChainConfigGetSeq() (int, error)
```

### 3.4 链配置更新获取Payload签名
```go
ChainConfigPayloadCollectSign(payloadBytes []byte) ([]byte, error)
```

### 3.5 链配置更新Payload签名收集&合并
```go
ChainConfigPayloadMergeSign(signedPayloadBytes [][]byte) ([]byte, error)
```

### 3.6 发送链配置更新请求
```go
SendChainConfigUpdateRequest(mergeSignedPayloadBytes []byte) (*pb.TxResponse, error)
```

> 以下ChainConfigCreateXXXXXXPayload方法，用于生成链配置待签名payload，在进行多签收集后(需机构Admin权限账号签名)，用于链配置的更新
### 3.7 更新Core模块待签名payload生成
**参数说明**
  - txSchedulerTimeout: 交易调度器从交易池拿到交易后, 进行调度的时间，其值范围为[0, 60]，若无需修改，请置为-1
  - txSchedulerValidateTimeout: 交易调度器从区块中拿到交易后, 进行验证的超时时间，其值范围为[0, 60]，若无需修改，请置为-1
```go
ChainConfigCreateCoreUpdatePayload(txSchedulerTimeout, txSchedulerValidateTimeout int) ([]byte, error)
```

### 3.8 更新Core模块待签名payload生成
**参数说明**
  - txTimestampVerify: 是否需要开启交易时间戳校验
  - (以下参数，若无需修改，请置为-1)
  - txTimeout: 交易时间戳的过期时间(秒)，其值范围为[600, +∞)
  - blockTxCapacity: 区块中最大交易数，其值范围为(0, +∞]
  - blockSize: 区块最大限制，单位MB，其值范围为(0, +∞]
  - blockInterval: 出块间隔，单位:ms，其值范围为[10, +∞]
```go
ChainConfigCreateBlockUpdatePayload(txTimestampVerify bool, txTimeout, blockTxCapacity, blockSize, blockInterval int) ([]byte, error)
```

### 3.9 添加信任组织根证书待签名payload生成
**参数说明**
  - trustRootOrgId: 组织Id
  - trustRootCrt: 根证书
```go
ChainConfigCreateTrustRootAddPayload(trustRootOrgId, trustRootCrt string) ([]byte, error)
```

### 3.10 更新信任组织根证书待签名payload生成
**参数说明**
  - trustRootOrgId: 组织Id
  - trustRootCrt: 根证书
```go
ChainConfigCreateTrustRootUpdatePayload(trustRootOrgId, trustRootCrt string) ([]byte, error)
```

### 3.11 删除信任组织根证书待签名payload生成
**参数说明**
  - trustRootOrgId: 组织Id
```go
ChainConfigCreateTrustRootDeletePayload(trustRootOrgId string) ([]byte, error)
```

### 3.12 添加权限配置待签名payload生成
**参数说明**
  - permissionResourceName: 权限名
  - principle: 权限规则
```go
ChainConfigCreatePermissionAddPayload(permissionResourceName string, principle *pb.Principle) ([]byte, error)
```

### 3.13 更新权限配置待签名payload生成
**参数说明**
  - permissionResourceName: 权限名
  - principle: 权限规则
```go
ChainConfigCreatePermissionUpdatePayload(permissionResourceName string, principle *pb.Principle) ([]byte, error)
```

### 3.14 删除权限配置待签名payload生成
**参数说明**
  - permissionResourceName: 权限名
```go
ChainConfigCreatePermissionDeletePayload(permissionResourceName string) ([]byte, error)
```

### 3.15 添加共识节点地址待签名payload生成
**参数说明**
  - nodeOrgId: 节点组织Id
  - nodeAddresses: 节点地址
```go
ChainConfigCreateConsensusNodeAddrAddPayload(nodeOrgId string, nodeAddresses []string) ([]byte, error)
```

### 3.16 更新共识节点地址待签名payload生成
**参数说明**
  - nodeOrgId: 节点组织Id
  - nodeOldAddress: 节点原地址
  - nodeNewAddress: 节点新地址
```go
ChainConfigCreateConsensusNodeAddrUpdatePayload(nodeOrgId, nodeOldAddress, nodeNewAddress string) ([]byte, error)
```

### 3.17 删除共识节点地址待签名payload生成
**参数说明**
  - nodeOrgId: 节点组织Id
  - nodeAddress: 节点地址
```go
ChainConfigCreateConsensusNodeAddrDeletePayload(nodeOrgId, nodeAddress string) ([]byte, error)
```

### 3.18 添加共识节点待签名payload生成
**参数说明**
  - nodeOrgId: 节点组织Id
  - nodeAddresses: 节点地址
```go
ChainConfigCreateConsensusNodeOrgAddPayload(nodeOrgId string, nodeAddresses []string) ([]byte, error)
```

### 3.19 更新共识节点待签名payload生成
**参数说明**
  - nodeOrgId: 节点组织Id
  - nodeAddresses: 节点地址
```go
ChainConfigCreateConsensusNodeOrgUpdatePayload(nodeOrgId string, nodeAddresses []string) ([]byte, error)
```

### 3.20 删除共识节点待签名payload生成
**参数说明**
  - nodeOrgId: 节点组织Id
```go
ChainConfigCreateConsensusNodeOrgDeletePayload(nodeOrgId string) ([]byte, error)
```

### 3.21 添加共识扩展字段待签名payload生成
**参数说明**
  - kvs: 字段key、value对
```go
ChainConfigCreateConsensusExtAddPayload(kvs []*pb.KeyValuePair) ([]byte, error)
```

### 3.22 添加共识扩展字段待签名payload生成
**参数说明**
  - kvs: 字段key、value对
```go
ChainConfigCreateConsensusExtUpdatePayload(kvs []*pb.KeyValuePair) ([]byte, error)
```

### 3.23 添加共识扩展字段待签名payload生成
**参数说明**
  - keys: 待删除字段
```go
ChainConfigCreateConsensusExtDeletePayload(keys []string) ([]byte, error)
```

## 4 证书管理接口
### 4.1 用户证书添加
**参数说明**
  - 在pb.TxResponse.ContractResult.Result字段中返回成功添加的certHash
```go
CertAdd() (*pb.TxResponse, error)
```

### 4.2 用户证书删除
**参数说明**
  - certHashes: 证书Hash列表，多个使用逗号分割
```go
CertDelete(certHashes string) (*pb.TxResponse, error)
```

### 4.3 用户证书查询
**参数说明**
  - certHashes: 证书Hash列表，多个使用逗号分割
返回值说明：
  - *pb.CertInfos: 包含证书Hash和证书内容的列表
```go
CertQuery(certHashes string) (*pb.CertInfos, error)
```

## 5 消息订阅接口
### 5.1 区块订阅
**参数说明**
  - startBlock: 订阅起始区块高度，若为-1，表示订阅实时最新区块
  - endBlock: 订阅结束区块高度，若为-1，表示订阅实时最新区块
  - withRwSet: 是否返回读写集
```go
SubscribeBlock(ctx context.Context, startBlock, endBlock int64, withRwSet bool) (<-chan interface{}, error)
```

### 5.2 交易订阅
**参数说明**
  - startBlock: 订阅起始区块高度，若为-1，表示订阅实时最新区块
  - endBlock: 订阅结束区块高度，若为-1，表示订阅实时最新区块
  - txType: 订阅交易类型,若为pb.TxType(-1)，表示订阅所有交易类型
  - txIds: 订阅txId列表，若为空，表示订阅所有txId
```go
SubscribeTx(ctx context.Context, startBlock, endBlock int64, txType pb.TxType, txIds []string) (<-chan interface{}, error)
```

### 5.3 多合一订阅
**参数说明**
  - txType: 订阅交易类型，目前已支持：区块消息订阅(pb.TxType_SUBSCRIBE_BLOCK_INFO)、交易消息订阅(pb.TxType_SUBSCRIBE_TX_INFO)
  - payloadBytes: 消息订阅参数payload
```go
Subscribe(ctx context.Context, txType pb.TxType, payloadBytes []byte) (<-chan interface{}, error)
```

## 6 管理类接口
### 6.1 SDK停止接口：关闭连接池连接，释放资源
```go
Stop() error
```
