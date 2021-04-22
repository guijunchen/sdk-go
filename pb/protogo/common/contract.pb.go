// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: common/contract.proto

package common

import (
	fmt "fmt"
	proto "github.com/gogo/protobuf/proto"
	io "io"
	math "math"
	math_bits "math/bits"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.GoGoProtoPackageIsVersion3 // please upgrade the proto package

type ContractName int32

const (
	// system chain configuration contract
	// used to add, delete and change the chain configuration
	ContractName_SYSTEM_CONTRACT_CHAIN_CONFIG ContractName = 0
	// system chain query contract
	// used to query the configuration on the chain
	ContractName_SYSTEM_CONTRACT_QUERY ContractName = 1
	// system certificate storage contract
	// used to manage certificates
	ContractName_SYSTEM_CONTRACT_CERT_MANAGE ContractName = 2
	// governance contract
	ContractName_SYSTEM_CONTRACT_GOVERNANCE ContractName = 3
	// multi signature contract on chain
	ContractName_SYSTEM_CONTRACT_MULTI_SIGN ContractName = 4
)

var ContractName_name = map[int32]string{
	0: "SYSTEM_CONTRACT_CHAIN_CONFIG",
	1: "SYSTEM_CONTRACT_QUERY",
	2: "SYSTEM_CONTRACT_CERT_MANAGE",
	3: "SYSTEM_CONTRACT_GOVERNANCE",
	4: "SYSTEM_CONTRACT_MULTI_SIGN",
}

var ContractName_value = map[string]int32{
	"SYSTEM_CONTRACT_CHAIN_CONFIG": 0,
	"SYSTEM_CONTRACT_QUERY":        1,
	"SYSTEM_CONTRACT_CERT_MANAGE":  2,
	"SYSTEM_CONTRACT_GOVERNANCE":   3,
	"SYSTEM_CONTRACT_MULTI_SIGN":   4,
}

func (x ContractName) String() string {
	return proto.EnumName(ContractName_name, int32(x))
}

func (ContractName) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_a1095c55e7168440, []int{0}
}

type QueryFunction int32

const (
	// get block by transactionId
	QueryFunction_GET_BLOCK_BY_TX_ID QueryFunction = 0
	// get contract information
	QueryFunction_GET_CONTRACT_INFO QueryFunction = 1
	// get transaction by transactionId
	QueryFunction_GET_TX_BY_TX_ID QueryFunction = 2
	// get block by block height
	QueryFunction_GET_BLOCK_BY_HEIGHT QueryFunction = 3
	// get chain information
	QueryFunction_GET_CHAIN_INFO QueryFunction = 4
	// get the last configuration block
	QueryFunction_GET_LAST_CONFIG_BLOCK QueryFunction = 5
	// get block by block hash
	QueryFunction_GET_BLOCK_BY_HASH QueryFunction = 6
	// get the list of node
	QueryFunction_GET_NODE_CHAIN_LIST QueryFunction = 7
	// get governance information
	QueryFunction_GET_GOVERNANCE_CONTRACT QueryFunction = 8
	// get read/write set information by eight
	QueryFunction_GET_BLOCK_WITH_TXRWSETS_BY_HEIGHT QueryFunction = 9
	// get read/write set information by hash
	QueryFunction_GET_BLOCK_WITH_TXRWSETS_BY_HASH QueryFunction = 10
	// get the last block
	QueryFunction_GET_LAST_BLOCK QueryFunction = 11
)

var QueryFunction_name = map[int32]string{
	0:  "GET_BLOCK_BY_TX_ID",
	1:  "GET_CONTRACT_INFO",
	2:  "GET_TX_BY_TX_ID",
	3:  "GET_BLOCK_BY_HEIGHT",
	4:  "GET_CHAIN_INFO",
	5:  "GET_LAST_CONFIG_BLOCK",
	6:  "GET_BLOCK_BY_HASH",
	7:  "GET_NODE_CHAIN_LIST",
	8:  "GET_GOVERNANCE_CONTRACT",
	9:  "GET_BLOCK_WITH_TXRWSETS_BY_HEIGHT",
	10: "GET_BLOCK_WITH_TXRWSETS_BY_HASH",
	11: "GET_LAST_BLOCK",
}

var QueryFunction_value = map[string]int32{
	"GET_BLOCK_BY_TX_ID":                0,
	"GET_CONTRACT_INFO":                 1,
	"GET_TX_BY_TX_ID":                   2,
	"GET_BLOCK_BY_HEIGHT":               3,
	"GET_CHAIN_INFO":                    4,
	"GET_LAST_CONFIG_BLOCK":             5,
	"GET_BLOCK_BY_HASH":                 6,
	"GET_NODE_CHAIN_LIST":               7,
	"GET_GOVERNANCE_CONTRACT":           8,
	"GET_BLOCK_WITH_TXRWSETS_BY_HEIGHT": 9,
	"GET_BLOCK_WITH_TXRWSETS_BY_HASH":   10,
	"GET_LAST_BLOCK":                    11,
}

func (x QueryFunction) String() string {
	return proto.EnumName(QueryFunction_name, int32(x))
}

func (QueryFunction) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_a1095c55e7168440, []int{1}
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

var RuntimeType_name = map[int32]string{
	0: "INVALID",
	1: "NATIVE",
	2: "WASMER",
	3: "WXVM",
	4: "GASM",
	5: "EVM",
	6: "DOCKER_GO",
	7: "DOCKER_JAVA",
}

var RuntimeType_value = map[string]int32{
	"INVALID":     0,
	"NATIVE":      1,
	"WASMER":      2,
	"WXVM":        3,
	"GASM":        4,
	"EVM":         5,
	"DOCKER_GO":   6,
	"DOCKER_JAVA": 7,
}

func (x RuntimeType) String() string {
	return proto.EnumName(RuntimeType_name, int32(x))
}

func (RuntimeType) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_a1095c55e7168440, []int{2}
}

type ConfigFunction int32

const (
	// get chain configuration
	ConfigFunction_GET_CHAIN_CONFIG ConfigFunction = 0
	// get the latest configuration block
	// the incoming blockheight must exist in the database
	// 如果传入的blockHeight是配置块，直接返回当前的配置信息
	ConfigFunction_GET_CHAIN_CONFIG_AT ConfigFunction = 1
	// update core
	ConfigFunction_CORE_UPDATE ConfigFunction = 2
	// update block
	ConfigFunction_BLOCK_UPDATE ConfigFunction = 3
	// add trusted certficate (org_id and root)
	ConfigFunction_TRUST_ROOT_ADD ConfigFunction = 4
	// [self] modify an individual's own trusted root certificate [org_id must exist in the original trust_roots,
	// and the new root certificate must be different from other certificates]
	ConfigFunction_TRUST_ROOT_UPDATE ConfigFunction = 5
	// delete trusted root certificate [org_ ID should be in trust_ The nodes in nodes need to be deleted]
	ConfigFunction_TRUST_ROOT_DELETE ConfigFunction = 6
	// organization add node address
	// org_id must already exist in nodes，you can add addresses in batches
	// the parameter is addresses. Single addresses are separated by ","
	// ip+port and peerid cannot be repeated
	ConfigFunction_NODE_ADDR_ADD ConfigFunction = 7
	// [self]the organization updates an address
	//[org_id and address must already exist in nodes, new_address is the new address. ip+port and peerId cannot be duplicated]
	ConfigFunction_NODE_ADDR_UPDATE ConfigFunction = 8
	// organization delete node address [org_id and address must already exist in nodes]
	ConfigFunction_NODE_ADDR_DELETE ConfigFunction = 9
	// organization add node address in batches [org_id在nodes不存在，批量添加地址，参数为addresses，单地址用逗号","隔开。ip+port和peerId不能重复]
	ConfigFunction_NODE_ORG_ADD ConfigFunction = 10
	// organization update
	// org_id must already exist in nodes，the parameter is addresses，Single addresses are separated by ","
	// ip+port and peerid cannot be repeated
	ConfigFunction_NODE_ORG_UPDATE ConfigFunction = 11
	// organization delete, org_id must already exist in nodes
	ConfigFunction_NODE_ORG_DELETE ConfigFunction = 12
	// add consensus parameters, key is not exit in ext_config
	ConfigFunction_CONSENSUS_EXT_ADD ConfigFunction = 13
	// update onsensus parameters, key exit in ext_config
	ConfigFunction_CONSENSUS_EXT_UPDATE ConfigFunction = 14
	// delete onsensus parameters, key exit in ext_config
	ConfigFunction_CONSENSUS_EXT_DELETE ConfigFunction = 15
	// add permission
	ConfigFunction_PERMISSION_ADD ConfigFunction = 16
	// update permission
	ConfigFunction_PERMISSION_UPDATE ConfigFunction = 17
	// delete permission
	ConfigFunction_PERMISSION_DELETE ConfigFunction = 18
)

var ConfigFunction_name = map[int32]string{
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

var ConfigFunction_value = map[string]int32{
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

func (x ConfigFunction) String() string {
	return proto.EnumName(ConfigFunction_name, int32(x))
}

func (ConfigFunction) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_a1095c55e7168440, []int{3}
}

// Methods of certificate management
type CertManageFunction int32

const (
	// add certficate
	CertManageFunction_CERT_ADD CertManageFunction = 0
	// delete certficate
	CertManageFunction_CERTS_DELETE CertManageFunction = 1
	// query certficate
	CertManageFunction_CERTS_QUERY CertManageFunction = 2
	// freeze certificate
	CertManageFunction_CERTS_FREEZE CertManageFunction = 3
	// unfreezing certificate
	CertManageFunction_CERTS_UNFREEZE CertManageFunction = 4
	// Revocation of certificate
	CertManageFunction_CERTS_REVOKE CertManageFunction = 5
)

var CertManageFunction_name = map[int32]string{
	0: "CERT_ADD",
	1: "CERTS_DELETE",
	2: "CERTS_QUERY",
	3: "CERTS_FREEZE",
	4: "CERTS_UNFREEZE",
	5: "CERTS_REVOKE",
}

var CertManageFunction_value = map[string]int32{
	"CERT_ADD":       0,
	"CERTS_DELETE":   1,
	"CERTS_QUERY":    2,
	"CERTS_FREEZE":   3,
	"CERTS_UNFREEZE": 4,
	"CERTS_REVOKE":   5,
}

func (x CertManageFunction) String() string {
	return proto.EnumName(CertManageFunction_name, int32(x))
}

func (CertManageFunction) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_a1095c55e7168440, []int{4}
}

// methods of managing multi signature
type MultiSignFunction int32

const (
	// multi signature request
	MultiSignFunction_REQ MultiSignFunction = 0
	// multi signature voting
	MultiSignFunction_VOTE MultiSignFunction = 1
	// multi signature query
	MultiSignFunction_QUERY MultiSignFunction = 2
)

var MultiSignFunction_name = map[int32]string{
	0: "REQ",
	1: "VOTE",
	2: "QUERY",
}

var MultiSignFunction_value = map[string]int32{
	"REQ":   0,
	"VOTE":  1,
	"QUERY": 2,
}

func (x MultiSignFunction) String() string {
	return proto.EnumName(MultiSignFunction_name, int32(x))
}

func (MultiSignFunction) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_a1095c55e7168440, []int{5}
}

// methods of user management contract
type ManageUserContractFunction int32

const (
	// init contract
	ManageUserContractFunction_INIT_CONTRACT ManageUserContractFunction = 0
	// upgrade contract
	ManageUserContractFunction_UPGRADE_CONTRACT ManageUserContractFunction = 1
	// freeze  contract
	ManageUserContractFunction_FREEZE_CONTRACT ManageUserContractFunction = 2
	// unfreezing contract
	ManageUserContractFunction_UNFREEZE_CONTRACT ManageUserContractFunction = 3
	// Revocation of contract
	ManageUserContractFunction_REVOKE_CONTRACT ManageUserContractFunction = 4
)

var ManageUserContractFunction_name = map[int32]string{
	0: "INIT_CONTRACT",
	1: "UPGRADE_CONTRACT",
	2: "FREEZE_CONTRACT",
	3: "UNFREEZE_CONTRACT",
	4: "REVOKE_CONTRACT",
}

var ManageUserContractFunction_value = map[string]int32{
	"INIT_CONTRACT":     0,
	"UPGRADE_CONTRACT":  1,
	"FREEZE_CONTRACT":   2,
	"UNFREEZE_CONTRACT": 3,
	"REVOKE_CONTRACT":   4,
}

func (x ManageUserContractFunction) String() string {
	return proto.EnumName(ManageUserContractFunction_name, int32(x))
}

func (ManageUserContractFunction) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_a1095c55e7168440, []int{6}
}

// the unique identifier of a smart contract
type ContractId struct {
	// smart contract name, set by contract creator, can have multiple versions
	ContractName string `protobuf:"bytes,1,opt,name=contract_name,json=contractName,proto3" json:"contract_name,omitempty"`
	// smart contract version, set by contract creator, name + version should be unique
	ContractVersion string `protobuf:"bytes,2,opt,name=contract_version,json=contractVersion,proto3" json:"contract_version,omitempty"`
	// smart contract runtime type, set by contract creator
	RuntimeType RuntimeType `protobuf:"varint,3,opt,name=runtime_type,json=runtimeType,proto3,enum=common.RuntimeType" json:"runtime_type,omitempty"`
}

func (m *ContractId) Reset()         { *m = ContractId{} }
func (m *ContractId) String() string { return proto.CompactTextString(m) }
func (*ContractId) ProtoMessage()    {}
func (*ContractId) Descriptor() ([]byte, []int) {
	return fileDescriptor_a1095c55e7168440, []int{0}
}
func (m *ContractId) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *ContractId) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_ContractId.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *ContractId) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ContractId.Merge(m, src)
}
func (m *ContractId) XXX_Size() int {
	return m.Size()
}
func (m *ContractId) XXX_DiscardUnknown() {
	xxx_messageInfo_ContractId.DiscardUnknown(m)
}

var xxx_messageInfo_ContractId proto.InternalMessageInfo

func (m *ContractId) GetContractName() string {
	if m != nil {
		return m.ContractName
	}
	return ""
}

func (m *ContractId) GetContractVersion() string {
	if m != nil {
		return m.ContractVersion
	}
	return ""
}

func (m *ContractId) GetRuntimeType() RuntimeType {
	if m != nil {
		return m.RuntimeType
	}
	return RuntimeType_INVALID
}

type ContractInfo struct {
	ContractTransaction []*ContractTransaction `protobuf:"bytes,1,rep,name=contract_transaction,json=contractTransaction,proto3" json:"contract_transaction,omitempty"`
}

func (m *ContractInfo) Reset()         { *m = ContractInfo{} }
func (m *ContractInfo) String() string { return proto.CompactTextString(m) }
func (*ContractInfo) ProtoMessage()    {}
func (*ContractInfo) Descriptor() ([]byte, []int) {
	return fileDescriptor_a1095c55e7168440, []int{1}
}
func (m *ContractInfo) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *ContractInfo) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_ContractInfo.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *ContractInfo) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ContractInfo.Merge(m, src)
}
func (m *ContractInfo) XXX_Size() int {
	return m.Size()
}
func (m *ContractInfo) XXX_DiscardUnknown() {
	xxx_messageInfo_ContractInfo.DiscardUnknown(m)
}

var xxx_messageInfo_ContractInfo proto.InternalMessageInfo

func (m *ContractInfo) GetContractTransaction() []*ContractTransaction {
	if m != nil {
		return m.ContractTransaction
	}
	return nil
}

type ContractTransaction struct {
	ContractId *ContractId `protobuf:"bytes,1,opt,name=contract_id,json=contractId,proto3" json:"contract_id,omitempty"`
	TxId       string      `protobuf:"bytes,2,opt,name=tx_id,json=txId,proto3" json:"tx_id,omitempty"`
}

func (m *ContractTransaction) Reset()         { *m = ContractTransaction{} }
func (m *ContractTransaction) String() string { return proto.CompactTextString(m) }
func (*ContractTransaction) ProtoMessage()    {}
func (*ContractTransaction) Descriptor() ([]byte, []int) {
	return fileDescriptor_a1095c55e7168440, []int{2}
}
func (m *ContractTransaction) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *ContractTransaction) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_ContractTransaction.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *ContractTransaction) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ContractTransaction.Merge(m, src)
}
func (m *ContractTransaction) XXX_Size() int {
	return m.Size()
}
func (m *ContractTransaction) XXX_DiscardUnknown() {
	xxx_messageInfo_ContractTransaction.DiscardUnknown(m)
}

var xxx_messageInfo_ContractTransaction proto.InternalMessageInfo

func (m *ContractTransaction) GetContractId() *ContractId {
	if m != nil {
		return m.ContractId
	}
	return nil
}

func (m *ContractTransaction) GetTxId() string {
	if m != nil {
		return m.TxId
	}
	return ""
}

func init() {
	proto.RegisterEnum("common.ContractName", ContractName_name, ContractName_value)
	proto.RegisterEnum("common.QueryFunction", QueryFunction_name, QueryFunction_value)
	proto.RegisterEnum("common.RuntimeType", RuntimeType_name, RuntimeType_value)
	proto.RegisterEnum("common.ConfigFunction", ConfigFunction_name, ConfigFunction_value)
	proto.RegisterEnum("common.CertManageFunction", CertManageFunction_name, CertManageFunction_value)
	proto.RegisterEnum("common.MultiSignFunction", MultiSignFunction_name, MultiSignFunction_value)
	proto.RegisterEnum("common.ManageUserContractFunction", ManageUserContractFunction_name, ManageUserContractFunction_value)
	proto.RegisterType((*ContractId)(nil), "common.ContractId")
	proto.RegisterType((*ContractInfo)(nil), "common.ContractInfo")
	proto.RegisterType((*ContractTransaction)(nil), "common.ContractTransaction")
}

func init() { proto.RegisterFile("common/contract.proto", fileDescriptor_a1095c55e7168440) }

var fileDescriptor_a1095c55e7168440 = []byte{
	// 937 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x7c, 0x95, 0xcf, 0x6f, 0xe2, 0x46,
	0x14, 0xc7, 0x31, 0x3f, 0xc3, 0x03, 0x92, 0x61, 0x48, 0xba, 0x74, 0x53, 0xb1, 0x69, 0x56, 0x95,
	0x52, 0xa4, 0x82, 0x9a, 0x95, 0x7a, 0x77, 0xec, 0x89, 0x71, 0x83, 0xed, 0x64, 0x3c, 0x10, 0xb2,
	0x87, 0x5a, 0x84, 0x78, 0x29, 0x6a, 0xb1, 0x91, 0xe3, 0x54, 0x9b, 0x43, 0xaf, 0x3d, 0xf7, 0xd0,
	0x3f, 0x63, 0xff, 0x90, 0x3d, 0xee, 0xb1, 0xc7, 0x2a, 0xf9, 0x47, 0xaa, 0x19, 0xff, 0x22, 0x64,
	0xb5, 0x37, 0xe6, 0x33, 0xdf, 0xf7, 0x7d, 0x6f, 0xde, 0x1b, 0x0f, 0xb0, 0x37, 0xf3, 0x97, 0x4b,
	0xdf, 0xeb, 0xcf, 0x7c, 0x2f, 0x0c, 0xa6, 0xb3, 0xb0, 0xb7, 0x0a, 0xfc, 0xd0, 0xc7, 0xe5, 0x08,
	0x1f, 0xfe, 0x23, 0x01, 0x28, 0xf1, 0x96, 0x7e, 0x83, 0x5f, 0x43, 0x23, 0x11, 0x3a, 0xde, 0x74,
	0xe9, 0xb6, 0xa5, 0x03, 0xe9, 0xa8, 0x4a, 0xeb, 0x09, 0x34, 0xa7, 0x4b, 0x17, 0x7f, 0x0f, 0x28,
	0x15, 0xfd, 0xe1, 0x06, 0xb7, 0x0b, 0xdf, 0x6b, 0xe7, 0x85, 0x6e, 0x27, 0xe1, 0xe3, 0x08, 0xe3,
	0x9f, 0xa0, 0x1e, 0xdc, 0x79, 0xe1, 0x62, 0xe9, 0x3a, 0xe1, 0xfd, 0xca, 0x6d, 0x17, 0x0e, 0xa4,
	0xa3, 0xed, 0xe3, 0x56, 0x2f, 0xca, 0xde, 0xa3, 0xd1, 0x1e, 0xbb, 0x5f, 0xb9, 0xb4, 0x16, 0x64,
	0x8b, 0xc3, 0x5f, 0xa0, 0x9e, 0x56, 0xe5, 0xbd, 0xf3, 0xb1, 0x09, 0xbb, 0x69, 0xca, 0x30, 0x98,
	0x7a, 0xb7, 0xd3, 0x59, 0xc8, 0xd3, 0x4a, 0x07, 0x85, 0xa3, 0xda, 0xf1, 0x7e, 0xe2, 0x97, 0xc4,
	0xb0, 0x4c, 0x42, 0x5b, 0xb3, 0xe7, 0xf0, 0xd0, 0x81, 0xd6, 0x67, 0xb4, 0xf8, 0x0d, 0xd4, 0xd2,
	0x34, 0x8b, 0x1b, 0x71, 0xf8, 0xda, 0x31, 0xde, 0x74, 0xd7, 0x6f, 0x28, 0xcc, 0xb2, 0x9e, 0xb5,
	0xa0, 0x14, 0xbe, 0xe7, 0xf2, 0xa8, 0x07, 0xc5, 0xf0, 0xbd, 0x7e, 0xd3, 0xfd, 0x20, 0x65, 0x27,
	0x10, 0x4d, 0x3b, 0x80, 0x6f, 0xec, 0x2b, 0x9b, 0x11, 0xc3, 0x51, 0x2c, 0x93, 0x51, 0x59, 0x61,
	0x8e, 0x32, 0x90, 0x75, 0x93, 0x2f, 0x4f, 0x75, 0x0d, 0xe5, 0xf0, 0xd7, 0xb0, 0xb7, 0xa9, 0xb8,
	0x18, 0x11, 0x7a, 0x85, 0x24, 0xfc, 0x0a, 0xf6, 0x9f, 0x05, 0x13, 0xca, 0x1c, 0x43, 0x36, 0x65,
	0x8d, 0xa0, 0x3c, 0xee, 0xc0, 0xcb, 0x4d, 0x81, 0x66, 0x8d, 0x09, 0x35, 0x65, 0x53, 0x21, 0xa8,
	0xf0, 0xb9, 0x7d, 0x63, 0x34, 0x64, 0xba, 0x63, 0xeb, 0x9a, 0x89, 0x8a, 0xdd, 0x8f, 0x79, 0x68,
	0x5c, 0xdc, 0xb9, 0xc1, 0xfd, 0xe9, 0x9d, 0x17, 0xb5, 0xe2, 0x2b, 0xc0, 0x1a, 0x61, 0xce, 0xc9,
	0xd0, 0x52, 0xce, 0x9c, 0x93, 0x2b, 0x87, 0x4d, 0x1c, 0x5d, 0x45, 0x39, 0xbc, 0x07, 0x4d, 0xce,
	0x53, 0x1b, 0xdd, 0x3c, 0xb5, 0x90, 0x84, 0x5b, 0xb0, 0xc3, 0x31, 0x9b, 0x64, 0xda, 0x3c, 0x7e,
	0x01, 0xad, 0x27, 0x1e, 0x03, 0xa2, 0x6b, 0x03, 0x86, 0x0a, 0x18, 0xc3, 0xb6, 0x30, 0x11, 0x0d,
	0x10, 0x0e, 0x45, 0x7e, 0x7c, 0xce, 0x86, 0xb2, 0xcd, 0xe2, 0x9e, 0x44, 0x81, 0xa8, 0x94, 0xe4,
	0xcc, 0x7c, 0x64, 0x7b, 0x80, 0xca, 0x89, 0xbd, 0x69, 0xa9, 0x24, 0xb6, 0x1a, 0xea, 0x36, 0x43,
	0x15, 0xbc, 0x0f, 0x2f, 0xf8, 0x46, 0xd6, 0x81, 0xb4, 0x5c, 0xb4, 0x85, 0xbf, 0x83, 0x6f, 0x33,
	0xb3, 0x4b, 0x9d, 0x0d, 0x1c, 0x36, 0xa1, 0x97, 0x36, 0x61, 0xf6, 0x5a, 0x89, 0x55, 0xfc, 0x1a,
	0x5e, 0x7d, 0x49, 0xc6, 0x2b, 0x80, 0xe4, 0x1c, 0xa2, 0xe6, 0xa8, 0xd8, 0x5a, 0xd7, 0x87, 0xda,
	0xda, 0xb5, 0xc6, 0x35, 0xa8, 0xe8, 0xe6, 0x58, 0x1e, 0x8a, 0xe6, 0x01, 0x94, 0x4d, 0x99, 0xe9,
	0x63, 0x82, 0x24, 0xfe, 0xfb, 0x52, 0xb6, 0x0d, 0x42, 0x51, 0x1e, 0x6f, 0x41, 0xf1, 0x72, 0x32,
	0x36, 0x50, 0x81, 0xff, 0xd2, 0x64, 0xdb, 0x40, 0x45, 0x5c, 0x81, 0x02, 0x19, 0x1b, 0xa8, 0x84,
	0x1b, 0x50, 0x55, 0x2d, 0xe5, 0x8c, 0x50, 0x47, 0xb3, 0x50, 0x19, 0xef, 0x40, 0x2d, 0x5e, 0xfe,
	0x2c, 0x8f, 0x65, 0x54, 0xe9, 0x7e, 0x28, 0xc0, 0xb6, 0xe2, 0x7b, 0xef, 0x16, 0xf3, 0x74, 0x78,
	0xbb, 0x80, 0xb2, 0xfe, 0xa6, 0x17, 0x2c, 0xee, 0xd7, 0x3a, 0x75, 0x64, 0x86, 0x24, 0x6e, 0xa9,
	0x58, 0x94, 0x38, 0xa3, 0x73, 0x55, 0x66, 0xfc, 0x3a, 0x21, 0xa8, 0x47, 0x07, 0x8f, 0x89, 0x98,
	0x18, 0xa3, 0x23, 0x9b, 0x39, 0xd4, 0xb2, 0x98, 0x23, 0xab, 0x2a, 0x2a, 0xf2, 0xb1, 0xac, 0xb1,
	0x58, 0x5a, 0xda, 0xc0, 0x2a, 0x19, 0x12, 0x46, 0x50, 0x19, 0x37, 0xa1, 0x21, 0x26, 0x25, 0xab,
	0x2a, 0x15, 0x06, 0x15, 0x5e, 0x66, 0x86, 0xe2, 0xf8, 0xad, 0xa7, 0x34, 0x0e, 0xaf, 0xf2, 0x92,
	0x04, 0xb5, 0xa8, 0x26, 0xa2, 0x81, 0x5f, 0xb9, 0x94, 0xc4, 0xc1, 0xb5, 0x27, 0x30, 0x8e, 0xad,
	0xf3, 0x8a, 0x14, 0xcb, 0xb4, 0x89, 0x69, 0x8f, 0x6c, 0x87, 0x4c, 0xa2, 0xfa, 0x1b, 0xb8, 0x0d,
	0xbb, 0x4f, 0x71, 0xec, 0xb2, 0xfd, 0x7c, 0x27, 0xb6, 0xda, 0xe1, 0x7d, 0x38, 0x27, 0xd4, 0xd0,
	0x6d, 0x5b, 0xb7, 0x4c, 0xe1, 0x83, 0xb8, 0xfd, 0x1a, 0x8b, 0x4d, 0x9a, 0x1b, 0x38, 0x76, 0xc0,
	0xdd, 0x3f, 0x01, 0x2b, 0x6e, 0x10, 0x1a, 0x53, 0x6f, 0x3a, 0x77, 0xd3, 0x89, 0xd5, 0x61, 0x4b,
	0x7c, 0xd1, 0xdc, 0x31, 0xc7, 0x0f, 0xcb, 0x57, 0x76, 0x12, 0x15, 0x8d, 0x48, 0x90, 0xe8, 0x49,
	0xc8, 0x67, 0x92, 0x53, 0x4a, 0xc8, 0xdb, 0x78, 0x44, 0x11, 0x19, 0x99, 0x31, 0x2b, 0x66, 0x2a,
	0x4a, 0xc6, 0xd6, 0x19, 0x41, 0xa5, 0xee, 0x8f, 0xd0, 0x34, 0xee, 0x7e, 0x0f, 0x17, 0xf6, 0x62,
	0xee, 0xa5, 0xd9, 0x2b, 0x50, 0xa0, 0xe4, 0x02, 0xe5, 0xf8, 0xf5, 0x1b, 0x5b, 0x22, 0x61, 0x15,
	0x4a, 0x71, 0xaa, 0xee, 0x5f, 0x12, 0xbc, 0x8c, 0xca, 0x1d, 0xdd, 0xba, 0x41, 0xf2, 0xaa, 0xa5,
	0xc1, 0x4d, 0x68, 0xe8, 0xa6, 0x9e, 0x3d, 0x09, 0x28, 0xc7, 0x47, 0x38, 0x3a, 0xd7, 0xa8, 0xac,
	0xae, 0x7d, 0x79, 0xe2, 0x8d, 0x88, 0x0a, 0xcb, 0x60, 0x9e, 0x77, 0x29, 0xa9, 0x37, 0xc3, 0x05,
	0xae, 0x8d, 0x4a, 0xce, 0x60, 0xf1, 0xe4, 0xfa, 0xe3, 0x43, 0x47, 0xfa, 0xf4, 0xd0, 0x91, 0xfe,
	0x7b, 0xe8, 0x48, 0x7f, 0x3f, 0x76, 0x72, 0x9f, 0x1e, 0x3b, 0xb9, 0x7f, 0x1f, 0x3b, 0x39, 0x68,
	0xfb, 0xc1, 0xbc, 0x37, 0xfb, 0x75, 0xba, 0xf0, 0x96, 0xd3, 0xdf, 0xdc, 0xa0, 0xb7, 0xba, 0x8e,
	0x1f, 0xed, 0xb7, 0xeb, 0xd4, 0x0f, 0xe6, 0xfd, 0x6c, 0xf9, 0xc3, 0xdc, 0xef, 0xaf, 0xae, 0xfb,
	0xe2, 0xff, 0x70, 0xee, 0xf7, 0x23, 0xfd, 0x75, 0x59, 0xac, 0xdf, 0xfc, 0x1f, 0x00, 0x00, 0xff,
	0xff, 0xa4, 0x1b, 0x4e, 0x71, 0x38, 0x07, 0x00, 0x00,
}

func (m *ContractId) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *ContractId) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *ContractId) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.RuntimeType != 0 {
		i = encodeVarintContract(dAtA, i, uint64(m.RuntimeType))
		i--
		dAtA[i] = 0x18
	}
	if len(m.ContractVersion) > 0 {
		i -= len(m.ContractVersion)
		copy(dAtA[i:], m.ContractVersion)
		i = encodeVarintContract(dAtA, i, uint64(len(m.ContractVersion)))
		i--
		dAtA[i] = 0x12
	}
	if len(m.ContractName) > 0 {
		i -= len(m.ContractName)
		copy(dAtA[i:], m.ContractName)
		i = encodeVarintContract(dAtA, i, uint64(len(m.ContractName)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *ContractInfo) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *ContractInfo) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *ContractInfo) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.ContractTransaction) > 0 {
		for iNdEx := len(m.ContractTransaction) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.ContractTransaction[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintContract(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0xa
		}
	}
	return len(dAtA) - i, nil
}

func (m *ContractTransaction) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *ContractTransaction) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *ContractTransaction) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.TxId) > 0 {
		i -= len(m.TxId)
		copy(dAtA[i:], m.TxId)
		i = encodeVarintContract(dAtA, i, uint64(len(m.TxId)))
		i--
		dAtA[i] = 0x12
	}
	if m.ContractId != nil {
		{
			size, err := m.ContractId.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintContract(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func encodeVarintContract(dAtA []byte, offset int, v uint64) int {
	offset -= sovContract(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *ContractId) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.ContractName)
	if l > 0 {
		n += 1 + l + sovContract(uint64(l))
	}
	l = len(m.ContractVersion)
	if l > 0 {
		n += 1 + l + sovContract(uint64(l))
	}
	if m.RuntimeType != 0 {
		n += 1 + sovContract(uint64(m.RuntimeType))
	}
	return n
}

func (m *ContractInfo) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if len(m.ContractTransaction) > 0 {
		for _, e := range m.ContractTransaction {
			l = e.Size()
			n += 1 + l + sovContract(uint64(l))
		}
	}
	return n
}

func (m *ContractTransaction) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.ContractId != nil {
		l = m.ContractId.Size()
		n += 1 + l + sovContract(uint64(l))
	}
	l = len(m.TxId)
	if l > 0 {
		n += 1 + l + sovContract(uint64(l))
	}
	return n
}

func sovContract(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozContract(x uint64) (n int) {
	return sovContract(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *ContractId) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowContract
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: ContractId: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: ContractId: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field ContractName", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowContract
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthContract
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthContract
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.ContractName = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field ContractVersion", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowContract
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthContract
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthContract
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.ContractVersion = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 3:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field RuntimeType", wireType)
			}
			m.RuntimeType = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowContract
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.RuntimeType |= RuntimeType(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		default:
			iNdEx = preIndex
			skippy, err := skipContract(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthContract
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *ContractInfo) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowContract
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: ContractInfo: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: ContractInfo: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field ContractTransaction", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowContract
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthContract
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthContract
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.ContractTransaction = append(m.ContractTransaction, &ContractTransaction{})
			if err := m.ContractTransaction[len(m.ContractTransaction)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipContract(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthContract
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *ContractTransaction) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowContract
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: ContractTransaction: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: ContractTransaction: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field ContractId", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowContract
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthContract
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthContract
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.ContractId == nil {
				m.ContractId = &ContractId{}
			}
			if err := m.ContractId.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field TxId", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowContract
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthContract
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthContract
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.TxId = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipContract(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthContract
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func skipContract(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowContract
			}
			if iNdEx >= l {
				return 0, io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		wireType := int(wire & 0x7)
		switch wireType {
		case 0:
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowContract
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				iNdEx++
				if dAtA[iNdEx-1] < 0x80 {
					break
				}
			}
		case 1:
			iNdEx += 8
		case 2:
			var length int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowContract
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				length |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if length < 0 {
				return 0, ErrInvalidLengthContract
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupContract
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthContract
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthContract        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowContract          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupContract = fmt.Errorf("proto: unexpected end of group")
)
