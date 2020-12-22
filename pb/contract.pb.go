// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.25.0
// 	protoc        v3.13.0
// source: contract.proto

package pb

import (
	proto "github.com/golang/protobuf/proto"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// This is a compile-time assertion that a sufficiently up-to-date version
// of the legacy proto package is being used.
const _ = proto.ProtoPackageIsVersion4

type ContractName int32

const (
	// 系统链配置合约
	ContractName_SYSTEM_CONTRACT_CHAIN_CONFIG ContractName = 0
	// 系统查询合约
	ContractName_SYSTEM_CONTRACT_QUERY ContractName = 1
	// 系统证书存储合约
	ContractName_SYSTEM_CONTRACT_CERT_MANAGE ContractName = 2
)

// Enum value maps for ContractName.
var (
	ContractName_name = map[int32]string{
		0: "SYSTEM_CONTRACT_CHAIN_CONFIG",
		1: "SYSTEM_CONTRACT_QUERY",
		2: "SYSTEM_CONTRACT_CERT_MANAGE",
	}
	ContractName_value = map[string]int32{
		"SYSTEM_CONTRACT_CHAIN_CONFIG": 0,
		"SYSTEM_CONTRACT_QUERY":        1,
		"SYSTEM_CONTRACT_CERT_MANAGE":  2,
	}
)

func (x ContractName) Enum() *ContractName {
	p := new(ContractName)
	*p = x
	return p
}

func (x ContractName) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (ContractName) Descriptor() protoreflect.EnumDescriptor {
	return file_contract_proto_enumTypes[0].Descriptor()
}

func (ContractName) Type() protoreflect.EnumType {
	return &file_contract_proto_enumTypes[0]
}

func (x ContractName) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use ContractName.Descriptor instead.
func (ContractName) EnumDescriptor() ([]byte, []int) {
	return file_contract_proto_rawDescGZIP(), []int{0}
}

type QueryFunction int32

const (
	// 通过交易ID取区块
	QueryFunction_GET_BLOCK_BY_TX_ID QueryFunction = 0
	// 获取合约信息
	QueryFunction_GET_CONTRACT_INFO QueryFunction = 1
	// 通过交易ID取交易
	QueryFunction_GET_TX_BY_TX_ID QueryFunction = 2
	// 通过区块高度取区块
	QueryFunction_GET_BLOCK_BY_HEIGHT QueryFunction = 3
	// 获取链的基本信息
	QueryFunction_GET_CHAIN_INFO QueryFunction = 4
	// 获取最后一个配置块
	QueryFunction_GET_LAST_CONFIG_BLOCK QueryFunction = 5
	// 通过区块哈希取区块
	QueryFunction_GET_BLOCK_BY_HASH QueryFunction = 6
	// 节点加入的链列表
	QueryFunction_GET_NODE_CHAIN_LIST QueryFunction = 7
)

// Enum value maps for QueryFunction.
var (
	QueryFunction_name = map[int32]string{
		0: "GET_BLOCK_BY_TX_ID",
		1: "GET_CONTRACT_INFO",
		2: "GET_TX_BY_TX_ID",
		3: "GET_BLOCK_BY_HEIGHT",
		4: "GET_CHAIN_INFO",
		5: "GET_LAST_CONFIG_BLOCK",
		6: "GET_BLOCK_BY_HASH",
		7: "GET_NODE_CHAIN_LIST",
	}
	QueryFunction_value = map[string]int32{
		"GET_BLOCK_BY_TX_ID":    0,
		"GET_CONTRACT_INFO":     1,
		"GET_TX_BY_TX_ID":       2,
		"GET_BLOCK_BY_HEIGHT":   3,
		"GET_CHAIN_INFO":        4,
		"GET_LAST_CONFIG_BLOCK": 5,
		"GET_BLOCK_BY_HASH":     6,
		"GET_NODE_CHAIN_LIST":   7,
	}
)

func (x QueryFunction) Enum() *QueryFunction {
	p := new(QueryFunction)
	*p = x
	return p
}

func (x QueryFunction) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (QueryFunction) Descriptor() protoreflect.EnumDescriptor {
	return file_contract_proto_enumTypes[1].Descriptor()
}

func (QueryFunction) Type() protoreflect.EnumType {
	return &file_contract_proto_enumTypes[1]
}

func (x QueryFunction) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use QueryFunction.Descriptor instead.
func (QueryFunction) EnumDescriptor() ([]byte, []int) {
	return file_contract_proto_rawDescGZIP(), []int{1}
}

// smart contract runtime, contains vm type and language type
type RuntimeType int32

const (
	RuntimeType_INVALID RuntimeType = 0
	// native implement in chainmaker-go
	RuntimeType_NATIVE RuntimeType = 1
	// vm-wasmer, language-c++
	RuntimeType_WASMER RuntimeType = 2
	// vm-wxvm, language-cpp
	RuntimeType_WXVM RuntimeType = 3
	// wasm interpreter in go
	RuntimeType_GASM RuntimeType = 4
	// vm-evm
	RuntimeType_EVM RuntimeType = 5
	// vm-docker, language-golang
	RuntimeType_DOCKER_GO RuntimeType = 6
	// vm-docker, language-java
	RuntimeType_DOCKER_JAVA RuntimeType = 7
)

// Enum value maps for RuntimeType.
var (
	RuntimeType_name = map[int32]string{
		0: "INVALID",
		1: "NATIVE",
		2: "WASMER",
		3: "WXVM",
		4: "GASM",
		5: "EVM",
		6: "DOCKER_GO",
		7: "DOCKER_JAVA",
	}
	RuntimeType_value = map[string]int32{
		"INVALID":     0,
		"NATIVE":      1,
		"WASMER":      2,
		"WXVM":        3,
		"GASM":        4,
		"EVM":         5,
		"DOCKER_GO":   6,
		"DOCKER_JAVA": 7,
	}
)

func (x RuntimeType) Enum() *RuntimeType {
	p := new(RuntimeType)
	*p = x
	return p
}

func (x RuntimeType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (RuntimeType) Descriptor() protoreflect.EnumDescriptor {
	return file_contract_proto_enumTypes[2].Descriptor()
}

func (RuntimeType) Type() protoreflect.EnumType {
	return &file_contract_proto_enumTypes[2]
}

func (x RuntimeType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use RuntimeType.Descriptor instead.
func (RuntimeType) EnumDescriptor() ([]byte, []int) {
	return file_contract_proto_rawDescGZIP(), []int{2}
}

type ConfigFunction int32

const (
	// 查询链配置
	ConfigFunction_GET_CHAIN_CONFIG ConfigFunction = 0
	// 获取最近的配置区块信息
	// 传入的blockHeight必须在数据库中存在
	// 如果传入的blockHeight是配置块，直接返回当前的配置信息
	ConfigFunction_GET_CHAIN_CONFIG_AT ConfigFunction = 1
	// 更新core
	ConfigFunction_CORE_UPDATE ConfigFunction = 2
	// 更新Block
	ConfigFunction_BLOCK_UPDATE ConfigFunction = 3
	// 添加可信根证书【org_id和root和原来的需不同】
	ConfigFunction_TRUST_ROOT_ADD ConfigFunction = 4
	// [self]修改个人自己的可信任根证书【org_id必须存在原trust_roots中，并且新的root和其他证书需不同】
	ConfigFunction_TRUST_ROOT_UPDATE ConfigFunction = 5
	// 删除可信任根证书【org_id需在trust_roots中，并且需要将nodes中相关的节点删除】
	ConfigFunction_TRUST_ROOT_DELETE ConfigFunction = 6
	// 机构添加节点地址【org_id需在nodes已存在，可批量添加地址，参数为addresses，单地址用逗号","隔开。ip+port和peerId不能重复】
	ConfigFunction_NODE_ADDR_ADD ConfigFunction = 7
	// [self]机构更新某一地址【org_id和address需在nodes已存在，new_address为新地址。ip+port和peerId不能重复】
	ConfigFunction_NODE_ADDR_UPDATE ConfigFunction = 8
	// 机构删除某一地址【org_id和address需在nodes已存在】
	ConfigFunction_NODE_ADDR_DELETE ConfigFunction = 9
	// 机构添加【org_id在nodes不存在，批量添加地址，参数为addresses，单地址用逗号","隔开。ip+port和peerId不能重复】
	ConfigFunction_NODE_ORG_ADD ConfigFunction = 10
	// 机构更新【org_id需在nodes已存在，参数为addresses，单地址用逗号","隔开。ip+port和peerId不能重复】
	ConfigFunction_NODE_ORG_UPDATE ConfigFunction = 11
	// 机构删除【org_id需在nodes已存在】
	ConfigFunction_NODE_ORG_DELETE ConfigFunction = 12
	// 共识参数添加【key在ext_config中不存在】
	ConfigFunction_CONSENSUS_EXT_ADD ConfigFunction = 13
	// 共识参数更新【key在ext_config中存在】
	ConfigFunction_CONSENSUS_EXT_UPDATE ConfigFunction = 14
	// 共识参数删除【key在ext_config中存在】
	ConfigFunction_CONSENSUS_EXT_DELETE ConfigFunction = 15
	// 权限添加
	ConfigFunction_PERMISSION_ADD ConfigFunction = 16
	// 权限添加
	ConfigFunction_PERMISSION_UPDATE ConfigFunction = 17
	// 权限删除
	ConfigFunction_PERMISSION_DELETE ConfigFunction = 18
)

// Enum value maps for ConfigFunction.
var (
	ConfigFunction_name = map[int32]string{
		0:  "GET_CHAIN_CONFIG",
		1:  "GET_CHAIN_CONFIG_AT",
		2:  "CORE_UPDATE",
		3:  "BLOCK_UPDATE",
		4:  "TRUST_ROOT_ADD",
		5:  "TRUST_ROOT_UPDATE",
		6:  "TRUST_ROOT_DELETE",
		7:  "NODE_ADDR_ADD",
		8:  "NODE_ADDR_UPDATE",
		9:  "NODE_ADDR_DELETE",
		10: "NODE_ORG_ADD",
		11: "NODE_ORG_UPDATE",
		12: "NODE_ORG_DELETE",
		13: "CONSENSUS_EXT_ADD",
		14: "CONSENSUS_EXT_UPDATE",
		15: "CONSENSUS_EXT_DELETE",
		16: "PERMISSION_ADD",
		17: "PERMISSION_UPDATE",
		18: "PERMISSION_DELETE",
	}
	ConfigFunction_value = map[string]int32{
		"GET_CHAIN_CONFIG":     0,
		"GET_CHAIN_CONFIG_AT":  1,
		"CORE_UPDATE":          2,
		"BLOCK_UPDATE":         3,
		"TRUST_ROOT_ADD":       4,
		"TRUST_ROOT_UPDATE":    5,
		"TRUST_ROOT_DELETE":    6,
		"NODE_ADDR_ADD":        7,
		"NODE_ADDR_UPDATE":     8,
		"NODE_ADDR_DELETE":     9,
		"NODE_ORG_ADD":         10,
		"NODE_ORG_UPDATE":      11,
		"NODE_ORG_DELETE":      12,
		"CONSENSUS_EXT_ADD":    13,
		"CONSENSUS_EXT_UPDATE": 14,
		"CONSENSUS_EXT_DELETE": 15,
		"PERMISSION_ADD":       16,
		"PERMISSION_UPDATE":    17,
		"PERMISSION_DELETE":    18,
	}
)

func (x ConfigFunction) Enum() *ConfigFunction {
	p := new(ConfigFunction)
	*p = x
	return p
}

func (x ConfigFunction) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (ConfigFunction) Descriptor() protoreflect.EnumDescriptor {
	return file_contract_proto_enumTypes[3].Descriptor()
}

func (ConfigFunction) Type() protoreflect.EnumType {
	return &file_contract_proto_enumTypes[3]
}

func (x ConfigFunction) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use ConfigFunction.Descriptor instead.
func (ConfigFunction) EnumDescriptor() ([]byte, []int) {
	return file_contract_proto_rawDescGZIP(), []int{3}
}

// 证书管理的方法集
type CertManageFunction int32

const (
	// 证书添加
	CertManageFunction_CERT_ADD CertManageFunction = 0
	// 证书集删除
	CertManageFunction_CERTS_DELETE CertManageFunction = 1
	// 证书集查询
	CertManageFunction_CERTS_QUERY CertManageFunction = 2
)

// Enum value maps for CertManageFunction.
var (
	CertManageFunction_name = map[int32]string{
		0: "CERT_ADD",
		1: "CERTS_DELETE",
		2: "CERTS_QUERY",
	}
	CertManageFunction_value = map[string]int32{
		"CERT_ADD":     0,
		"CERTS_DELETE": 1,
		"CERTS_QUERY":  2,
	}
)

func (x CertManageFunction) Enum() *CertManageFunction {
	p := new(CertManageFunction)
	*p = x
	return p
}

func (x CertManageFunction) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (CertManageFunction) Descriptor() protoreflect.EnumDescriptor {
	return file_contract_proto_enumTypes[4].Descriptor()
}

func (CertManageFunction) Type() protoreflect.EnumType {
	return &file_contract_proto_enumTypes[4]
}

func (x CertManageFunction) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use CertManageFunction.Descriptor instead.
func (CertManageFunction) EnumDescriptor() ([]byte, []int) {
	return file_contract_proto_rawDescGZIP(), []int{4}
}

// the unique identifier of a smart contract
type ContractId struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// smart contract name, set by contract creator, can have multiple versions
	ContractName string `protobuf:"bytes,1,opt,name=contract_name,json=contractName,proto3" json:"contract_name,omitempty"`
	// smart contract version, set by contract creator, name + version should be unique
	ContractVersion string `protobuf:"bytes,2,opt,name=contract_version,json=contractVersion,proto3" json:"contract_version,omitempty"`
	// smart contract runtime type, set by contract creator
	RuntimeType RuntimeType `protobuf:"varint,3,opt,name=runtime_type,json=runtimeType,proto3,enum=pb.RuntimeType" json:"runtime_type,omitempty"`
}

func (x *ContractId) Reset() {
	*x = ContractId{}
	if protoimpl.UnsafeEnabled {
		mi := &file_contract_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ContractId) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ContractId) ProtoMessage() {}

func (x *ContractId) ProtoReflect() protoreflect.Message {
	mi := &file_contract_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ContractId.ProtoReflect.Descriptor instead.
func (*ContractId) Descriptor() ([]byte, []int) {
	return file_contract_proto_rawDescGZIP(), []int{0}
}

func (x *ContractId) GetContractName() string {
	if x != nil {
		return x.ContractName
	}
	return ""
}

func (x *ContractId) GetContractVersion() string {
	if x != nil {
		return x.ContractVersion
	}
	return ""
}

func (x *ContractId) GetRuntimeType() RuntimeType {
	if x != nil {
		return x.RuntimeType
	}
	return RuntimeType_INVALID
}

type ContractInfo struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ContractTransaction []*ContractTransaction `protobuf:"bytes,1,rep,name=contract_transaction,json=contractTransaction,proto3" json:"contract_transaction,omitempty"`
}

func (x *ContractInfo) Reset() {
	*x = ContractInfo{}
	if protoimpl.UnsafeEnabled {
		mi := &file_contract_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ContractInfo) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ContractInfo) ProtoMessage() {}

func (x *ContractInfo) ProtoReflect() protoreflect.Message {
	mi := &file_contract_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ContractInfo.ProtoReflect.Descriptor instead.
func (*ContractInfo) Descriptor() ([]byte, []int) {
	return file_contract_proto_rawDescGZIP(), []int{1}
}

func (x *ContractInfo) GetContractTransaction() []*ContractTransaction {
	if x != nil {
		return x.ContractTransaction
	}
	return nil
}

type ContractTransaction struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ContractId *ContractId `protobuf:"bytes,1,opt,name=contract_id,json=contractId,proto3" json:"contract_id,omitempty"`
	TxId       string      `protobuf:"bytes,2,opt,name=tx_id,json=txId,proto3" json:"tx_id,omitempty"`
}

func (x *ContractTransaction) Reset() {
	*x = ContractTransaction{}
	if protoimpl.UnsafeEnabled {
		mi := &file_contract_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ContractTransaction) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ContractTransaction) ProtoMessage() {}

func (x *ContractTransaction) ProtoReflect() protoreflect.Message {
	mi := &file_contract_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ContractTransaction.ProtoReflect.Descriptor instead.
func (*ContractTransaction) Descriptor() ([]byte, []int) {
	return file_contract_proto_rawDescGZIP(), []int{2}
}

func (x *ContractTransaction) GetContractId() *ContractId {
	if x != nil {
		return x.ContractId
	}
	return nil
}

func (x *ContractTransaction) GetTxId() string {
	if x != nil {
		return x.TxId
	}
	return ""
}

var File_contract_proto protoreflect.FileDescriptor

var file_contract_proto_rawDesc = []byte{
	0x0a, 0x0e, 0x63, 0x6f, 0x6e, 0x74, 0x72, 0x61, 0x63, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x12, 0x02, 0x70, 0x62, 0x22, 0x90, 0x01, 0x0a, 0x0a, 0x43, 0x6f, 0x6e, 0x74, 0x72, 0x61, 0x63,
	0x74, 0x49, 0x64, 0x12, 0x23, 0x0a, 0x0d, 0x63, 0x6f, 0x6e, 0x74, 0x72, 0x61, 0x63, 0x74, 0x5f,
	0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0c, 0x63, 0x6f, 0x6e, 0x74,
	0x72, 0x61, 0x63, 0x74, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x29, 0x0a, 0x10, 0x63, 0x6f, 0x6e, 0x74,
	0x72, 0x61, 0x63, 0x74, 0x5f, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x0f, 0x63, 0x6f, 0x6e, 0x74, 0x72, 0x61, 0x63, 0x74, 0x56, 0x65, 0x72, 0x73,
	0x69, 0x6f, 0x6e, 0x12, 0x32, 0x0a, 0x0c, 0x72, 0x75, 0x6e, 0x74, 0x69, 0x6d, 0x65, 0x5f, 0x74,
	0x79, 0x70, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x0f, 0x2e, 0x70, 0x62, 0x2e, 0x52,
	0x75, 0x6e, 0x74, 0x69, 0x6d, 0x65, 0x54, 0x79, 0x70, 0x65, 0x52, 0x0b, 0x72, 0x75, 0x6e, 0x74,
	0x69, 0x6d, 0x65, 0x54, 0x79, 0x70, 0x65, 0x22, 0x5a, 0x0a, 0x0c, 0x43, 0x6f, 0x6e, 0x74, 0x72,
	0x61, 0x63, 0x74, 0x49, 0x6e, 0x66, 0x6f, 0x12, 0x4a, 0x0a, 0x14, 0x63, 0x6f, 0x6e, 0x74, 0x72,
	0x61, 0x63, 0x74, 0x5f, 0x74, 0x72, 0x61, 0x6e, 0x73, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x18,
	0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x17, 0x2e, 0x70, 0x62, 0x2e, 0x43, 0x6f, 0x6e, 0x74, 0x72,
	0x61, 0x63, 0x74, 0x54, 0x72, 0x61, 0x6e, 0x73, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x13,
	0x63, 0x6f, 0x6e, 0x74, 0x72, 0x61, 0x63, 0x74, 0x54, 0x72, 0x61, 0x6e, 0x73, 0x61, 0x63, 0x74,
	0x69, 0x6f, 0x6e, 0x22, 0x5b, 0x0a, 0x13, 0x43, 0x6f, 0x6e, 0x74, 0x72, 0x61, 0x63, 0x74, 0x54,
	0x72, 0x61, 0x6e, 0x73, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x2f, 0x0a, 0x0b, 0x63, 0x6f,
	0x6e, 0x74, 0x72, 0x61, 0x63, 0x74, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32,
	0x0e, 0x2e, 0x70, 0x62, 0x2e, 0x43, 0x6f, 0x6e, 0x74, 0x72, 0x61, 0x63, 0x74, 0x49, 0x64, 0x52,
	0x0a, 0x63, 0x6f, 0x6e, 0x74, 0x72, 0x61, 0x63, 0x74, 0x49, 0x64, 0x12, 0x13, 0x0a, 0x05, 0x74,
	0x78, 0x5f, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x74, 0x78, 0x49, 0x64,
	0x2a, 0x6c, 0x0a, 0x0c, 0x43, 0x6f, 0x6e, 0x74, 0x72, 0x61, 0x63, 0x74, 0x4e, 0x61, 0x6d, 0x65,
	0x12, 0x20, 0x0a, 0x1c, 0x53, 0x59, 0x53, 0x54, 0x45, 0x4d, 0x5f, 0x43, 0x4f, 0x4e, 0x54, 0x52,
	0x41, 0x43, 0x54, 0x5f, 0x43, 0x48, 0x41, 0x49, 0x4e, 0x5f, 0x43, 0x4f, 0x4e, 0x46, 0x49, 0x47,
	0x10, 0x00, 0x12, 0x19, 0x0a, 0x15, 0x53, 0x59, 0x53, 0x54, 0x45, 0x4d, 0x5f, 0x43, 0x4f, 0x4e,
	0x54, 0x52, 0x41, 0x43, 0x54, 0x5f, 0x51, 0x55, 0x45, 0x52, 0x59, 0x10, 0x01, 0x12, 0x1f, 0x0a,
	0x1b, 0x53, 0x59, 0x53, 0x54, 0x45, 0x4d, 0x5f, 0x43, 0x4f, 0x4e, 0x54, 0x52, 0x41, 0x43, 0x54,
	0x5f, 0x43, 0x45, 0x52, 0x54, 0x5f, 0x4d, 0x41, 0x4e, 0x41, 0x47, 0x45, 0x10, 0x02, 0x2a, 0xcb,
	0x01, 0x0a, 0x0d, 0x51, 0x75, 0x65, 0x72, 0x79, 0x46, 0x75, 0x6e, 0x63, 0x74, 0x69, 0x6f, 0x6e,
	0x12, 0x16, 0x0a, 0x12, 0x47, 0x45, 0x54, 0x5f, 0x42, 0x4c, 0x4f, 0x43, 0x4b, 0x5f, 0x42, 0x59,
	0x5f, 0x54, 0x58, 0x5f, 0x49, 0x44, 0x10, 0x00, 0x12, 0x15, 0x0a, 0x11, 0x47, 0x45, 0x54, 0x5f,
	0x43, 0x4f, 0x4e, 0x54, 0x52, 0x41, 0x43, 0x54, 0x5f, 0x49, 0x4e, 0x46, 0x4f, 0x10, 0x01, 0x12,
	0x13, 0x0a, 0x0f, 0x47, 0x45, 0x54, 0x5f, 0x54, 0x58, 0x5f, 0x42, 0x59, 0x5f, 0x54, 0x58, 0x5f,
	0x49, 0x44, 0x10, 0x02, 0x12, 0x17, 0x0a, 0x13, 0x47, 0x45, 0x54, 0x5f, 0x42, 0x4c, 0x4f, 0x43,
	0x4b, 0x5f, 0x42, 0x59, 0x5f, 0x48, 0x45, 0x49, 0x47, 0x48, 0x54, 0x10, 0x03, 0x12, 0x12, 0x0a,
	0x0e, 0x47, 0x45, 0x54, 0x5f, 0x43, 0x48, 0x41, 0x49, 0x4e, 0x5f, 0x49, 0x4e, 0x46, 0x4f, 0x10,
	0x04, 0x12, 0x19, 0x0a, 0x15, 0x47, 0x45, 0x54, 0x5f, 0x4c, 0x41, 0x53, 0x54, 0x5f, 0x43, 0x4f,
	0x4e, 0x46, 0x49, 0x47, 0x5f, 0x42, 0x4c, 0x4f, 0x43, 0x4b, 0x10, 0x05, 0x12, 0x15, 0x0a, 0x11,
	0x47, 0x45, 0x54, 0x5f, 0x42, 0x4c, 0x4f, 0x43, 0x4b, 0x5f, 0x42, 0x59, 0x5f, 0x48, 0x41, 0x53,
	0x48, 0x10, 0x06, 0x12, 0x17, 0x0a, 0x13, 0x47, 0x45, 0x54, 0x5f, 0x4e, 0x4f, 0x44, 0x45, 0x5f,
	0x43, 0x48, 0x41, 0x49, 0x4e, 0x5f, 0x4c, 0x49, 0x53, 0x54, 0x10, 0x07, 0x2a, 0x6f, 0x0a, 0x0b,
	0x52, 0x75, 0x6e, 0x74, 0x69, 0x6d, 0x65, 0x54, 0x79, 0x70, 0x65, 0x12, 0x0b, 0x0a, 0x07, 0x49,
	0x4e, 0x56, 0x41, 0x4c, 0x49, 0x44, 0x10, 0x00, 0x12, 0x0a, 0x0a, 0x06, 0x4e, 0x41, 0x54, 0x49,
	0x56, 0x45, 0x10, 0x01, 0x12, 0x0a, 0x0a, 0x06, 0x57, 0x41, 0x53, 0x4d, 0x45, 0x52, 0x10, 0x02,
	0x12, 0x08, 0x0a, 0x04, 0x57, 0x58, 0x56, 0x4d, 0x10, 0x03, 0x12, 0x08, 0x0a, 0x04, 0x47, 0x41,
	0x53, 0x4d, 0x10, 0x04, 0x12, 0x07, 0x0a, 0x03, 0x45, 0x56, 0x4d, 0x10, 0x05, 0x12, 0x0d, 0x0a,
	0x09, 0x44, 0x4f, 0x43, 0x4b, 0x45, 0x52, 0x5f, 0x47, 0x4f, 0x10, 0x06, 0x12, 0x0f, 0x0a, 0x0b,
	0x44, 0x4f, 0x43, 0x4b, 0x45, 0x52, 0x5f, 0x4a, 0x41, 0x56, 0x41, 0x10, 0x07, 0x2a, 0xac, 0x03,
	0x0a, 0x0e, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x46, 0x75, 0x6e, 0x63, 0x74, 0x69, 0x6f, 0x6e,
	0x12, 0x14, 0x0a, 0x10, 0x47, 0x45, 0x54, 0x5f, 0x43, 0x48, 0x41, 0x49, 0x4e, 0x5f, 0x43, 0x4f,
	0x4e, 0x46, 0x49, 0x47, 0x10, 0x00, 0x12, 0x17, 0x0a, 0x13, 0x47, 0x45, 0x54, 0x5f, 0x43, 0x48,
	0x41, 0x49, 0x4e, 0x5f, 0x43, 0x4f, 0x4e, 0x46, 0x49, 0x47, 0x5f, 0x41, 0x54, 0x10, 0x01, 0x12,
	0x0f, 0x0a, 0x0b, 0x43, 0x4f, 0x52, 0x45, 0x5f, 0x55, 0x50, 0x44, 0x41, 0x54, 0x45, 0x10, 0x02,
	0x12, 0x10, 0x0a, 0x0c, 0x42, 0x4c, 0x4f, 0x43, 0x4b, 0x5f, 0x55, 0x50, 0x44, 0x41, 0x54, 0x45,
	0x10, 0x03, 0x12, 0x12, 0x0a, 0x0e, 0x54, 0x52, 0x55, 0x53, 0x54, 0x5f, 0x52, 0x4f, 0x4f, 0x54,
	0x5f, 0x41, 0x44, 0x44, 0x10, 0x04, 0x12, 0x15, 0x0a, 0x11, 0x54, 0x52, 0x55, 0x53, 0x54, 0x5f,
	0x52, 0x4f, 0x4f, 0x54, 0x5f, 0x55, 0x50, 0x44, 0x41, 0x54, 0x45, 0x10, 0x05, 0x12, 0x15, 0x0a,
	0x11, 0x54, 0x52, 0x55, 0x53, 0x54, 0x5f, 0x52, 0x4f, 0x4f, 0x54, 0x5f, 0x44, 0x45, 0x4c, 0x45,
	0x54, 0x45, 0x10, 0x06, 0x12, 0x11, 0x0a, 0x0d, 0x4e, 0x4f, 0x44, 0x45, 0x5f, 0x41, 0x44, 0x44,
	0x52, 0x5f, 0x41, 0x44, 0x44, 0x10, 0x07, 0x12, 0x14, 0x0a, 0x10, 0x4e, 0x4f, 0x44, 0x45, 0x5f,
	0x41, 0x44, 0x44, 0x52, 0x5f, 0x55, 0x50, 0x44, 0x41, 0x54, 0x45, 0x10, 0x08, 0x12, 0x14, 0x0a,
	0x10, 0x4e, 0x4f, 0x44, 0x45, 0x5f, 0x41, 0x44, 0x44, 0x52, 0x5f, 0x44, 0x45, 0x4c, 0x45, 0x54,
	0x45, 0x10, 0x09, 0x12, 0x10, 0x0a, 0x0c, 0x4e, 0x4f, 0x44, 0x45, 0x5f, 0x4f, 0x52, 0x47, 0x5f,
	0x41, 0x44, 0x44, 0x10, 0x0a, 0x12, 0x13, 0x0a, 0x0f, 0x4e, 0x4f, 0x44, 0x45, 0x5f, 0x4f, 0x52,
	0x47, 0x5f, 0x55, 0x50, 0x44, 0x41, 0x54, 0x45, 0x10, 0x0b, 0x12, 0x13, 0x0a, 0x0f, 0x4e, 0x4f,
	0x44, 0x45, 0x5f, 0x4f, 0x52, 0x47, 0x5f, 0x44, 0x45, 0x4c, 0x45, 0x54, 0x45, 0x10, 0x0c, 0x12,
	0x15, 0x0a, 0x11, 0x43, 0x4f, 0x4e, 0x53, 0x45, 0x4e, 0x53, 0x55, 0x53, 0x5f, 0x45, 0x58, 0x54,
	0x5f, 0x41, 0x44, 0x44, 0x10, 0x0d, 0x12, 0x18, 0x0a, 0x14, 0x43, 0x4f, 0x4e, 0x53, 0x45, 0x4e,
	0x53, 0x55, 0x53, 0x5f, 0x45, 0x58, 0x54, 0x5f, 0x55, 0x50, 0x44, 0x41, 0x54, 0x45, 0x10, 0x0e,
	0x12, 0x18, 0x0a, 0x14, 0x43, 0x4f, 0x4e, 0x53, 0x45, 0x4e, 0x53, 0x55, 0x53, 0x5f, 0x45, 0x58,
	0x54, 0x5f, 0x44, 0x45, 0x4c, 0x45, 0x54, 0x45, 0x10, 0x0f, 0x12, 0x12, 0x0a, 0x0e, 0x50, 0x45,
	0x52, 0x4d, 0x49, 0x53, 0x53, 0x49, 0x4f, 0x4e, 0x5f, 0x41, 0x44, 0x44, 0x10, 0x10, 0x12, 0x15,
	0x0a, 0x11, 0x50, 0x45, 0x52, 0x4d, 0x49, 0x53, 0x53, 0x49, 0x4f, 0x4e, 0x5f, 0x55, 0x50, 0x44,
	0x41, 0x54, 0x45, 0x10, 0x11, 0x12, 0x15, 0x0a, 0x11, 0x50, 0x45, 0x52, 0x4d, 0x49, 0x53, 0x53,
	0x49, 0x4f, 0x4e, 0x5f, 0x44, 0x45, 0x4c, 0x45, 0x54, 0x45, 0x10, 0x12, 0x2a, 0x45, 0x0a, 0x12,
	0x43, 0x65, 0x72, 0x74, 0x4d, 0x61, 0x6e, 0x61, 0x67, 0x65, 0x46, 0x75, 0x6e, 0x63, 0x74, 0x69,
	0x6f, 0x6e, 0x12, 0x0c, 0x0a, 0x08, 0x43, 0x45, 0x52, 0x54, 0x5f, 0x41, 0x44, 0x44, 0x10, 0x00,
	0x12, 0x10, 0x0a, 0x0c, 0x43, 0x45, 0x52, 0x54, 0x53, 0x5f, 0x44, 0x45, 0x4c, 0x45, 0x54, 0x45,
	0x10, 0x01, 0x12, 0x0f, 0x0a, 0x0b, 0x43, 0x45, 0x52, 0x54, 0x53, 0x5f, 0x51, 0x55, 0x45, 0x52,
	0x59, 0x10, 0x02, 0x42, 0x3b, 0x0a, 0x18, 0x6f, 0x72, 0x67, 0x2e, 0x63, 0x68, 0x61, 0x69, 0x6e,
	0x6d, 0x61, 0x6b, 0x65, 0x72, 0x2e, 0x70, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x5a,
	0x1f, 0x63, 0x68, 0x61, 0x69, 0x6e, 0x6d, 0x61, 0x6b, 0x65, 0x72, 0x2e, 0x6f, 0x72, 0x67, 0x2f,
	0x63, 0x68, 0x61, 0x69, 0x6e, 0x6d, 0x61, 0x6b, 0x65, 0x72, 0x2d, 0x67, 0x6f, 0x2f, 0x70, 0x62,
	0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_contract_proto_rawDescOnce sync.Once
	file_contract_proto_rawDescData = file_contract_proto_rawDesc
)

func file_contract_proto_rawDescGZIP() []byte {
	file_contract_proto_rawDescOnce.Do(func() {
		file_contract_proto_rawDescData = protoimpl.X.CompressGZIP(file_contract_proto_rawDescData)
	})
	return file_contract_proto_rawDescData
}

var file_contract_proto_enumTypes = make([]protoimpl.EnumInfo, 5)
var file_contract_proto_msgTypes = make([]protoimpl.MessageInfo, 3)
var file_contract_proto_goTypes = []interface{}{
	(ContractName)(0),           // 0: pb.ContractName
	(QueryFunction)(0),          // 1: pb.QueryFunction
	(RuntimeType)(0),            // 2: pb.RuntimeType
	(ConfigFunction)(0),         // 3: pb.ConfigFunction
	(CertManageFunction)(0),     // 4: pb.CertManageFunction
	(*ContractId)(nil),          // 5: pb.ContractId
	(*ContractInfo)(nil),        // 6: pb.ContractInfo
	(*ContractTransaction)(nil), // 7: pb.ContractTransaction
}
var file_contract_proto_depIdxs = []int32{
	2, // 0: pb.ContractId.runtime_type:type_name -> pb.RuntimeType
	7, // 1: pb.ContractInfo.contract_transaction:type_name -> pb.ContractTransaction
	5, // 2: pb.ContractTransaction.contract_id:type_name -> pb.ContractId
	3, // [3:3] is the sub-list for method output_type
	3, // [3:3] is the sub-list for method input_type
	3, // [3:3] is the sub-list for extension type_name
	3, // [3:3] is the sub-list for extension extendee
	0, // [0:3] is the sub-list for field type_name
}

func init() { file_contract_proto_init() }
func file_contract_proto_init() {
	if File_contract_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_contract_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ContractId); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_contract_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ContractInfo); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_contract_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ContractTransaction); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_contract_proto_rawDesc,
			NumEnums:      5,
			NumMessages:   3,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_contract_proto_goTypes,
		DependencyIndexes: file_contract_proto_depIdxs,
		EnumInfos:         file_contract_proto_enumTypes,
		MessageInfos:      file_contract_proto_msgTypes,
	}.Build()
	File_contract_proto = out.File
	file_contract_proto_rawDesc = nil
	file_contract_proto_goTypes = nil
	file_contract_proto_depIdxs = nil
}
