// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: txpool/transaction_pool.proto

package txpool

import (
	common "chainmaker.org/chainmaker-sdk-go/pb/protogo/common"
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

// SignalType is a transaction event type
type SignalType int32

const (
	// no transaction
	SignalType_NO_EVENT SignalType = 0
	// new transaction
	SignalType_TRANSACTION_INCOME SignalType = 1
	// packing block
	SignalType_BLOCK_PROPOSE SignalType = 2
)

var SignalType_name = map[int32]string{
	0: "NO_EVENT",
	1: "TRANSACTION_INCOME",
	2: "BLOCK_PROPOSE",
}

var SignalType_value = map[string]int32{
	"NO_EVENT":           0,
	"TRANSACTION_INCOME": 1,
	"BLOCK_PROPOSE":      2,
}

func (x SignalType) String() string {
	return proto.EnumName(SignalType_name, int32(x))
}

func (SignalType) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_0bc7127f197678cd, []int{0}
}

// TxPoolSignal is used by tx pool to send signal to block proposer
type TxPoolSignal struct {
	// transaction event type
	SignalType SignalType `protobuf:"varint,2,opt,name=signalType,proto3,enum=txpool.SignalType" json:"signalType,omitempty"`
	// chainId
	ChainId string `protobuf:"bytes,3,opt,name=chainId,proto3" json:"chainId,omitempty"`
}

func (m *TxPoolSignal) Reset()         { *m = TxPoolSignal{} }
func (m *TxPoolSignal) String() string { return proto.CompactTextString(m) }
func (*TxPoolSignal) ProtoMessage()    {}
func (*TxPoolSignal) Descriptor() ([]byte, []int) {
	return fileDescriptor_0bc7127f197678cd, []int{0}
}
func (m *TxPoolSignal) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *TxPoolSignal) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_TxPoolSignal.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *TxPoolSignal) XXX_Merge(src proto.Message) {
	xxx_messageInfo_TxPoolSignal.Merge(m, src)
}
func (m *TxPoolSignal) XXX_Size() int {
	return m.Size()
}
func (m *TxPoolSignal) XXX_DiscardUnknown() {
	xxx_messageInfo_TxPoolSignal.DiscardUnknown(m)
}

var xxx_messageInfo_TxPoolSignal proto.InternalMessageInfo

func (m *TxPoolSignal) GetSignalType() SignalType {
	if m != nil {
		return m.SignalType
	}
	return SignalType_NO_EVENT
}

func (m *TxPoolSignal) GetChainId() string {
	if m != nil {
		return m.ChainId
	}
	return ""
}

// transaction batch
type TxBatch struct {
	// batch ID
	BatchId int32 `protobuf:"varint,1,opt,name=batchId,proto3" json:"batchId,omitempty"`
	// node ID
	NodeId string `protobuf:"bytes,2,opt,name=nodeId,proto3" json:"nodeId,omitempty"`
	// batch size
	Size_ int32 `protobuf:"varint,3,opt,name=size,proto3" json:"size,omitempty"`
	// transaction list
	Txs []*common.Transaction `protobuf:"bytes,4,rep,name=txs,proto3" json:"txs,omitempty"`
	// Map: transaction ID mapping record( key: transaction ID, value: transaction index in txs)
	TxIdsMap map[string]int32 `protobuf:"bytes,5,rep,name=txIdsMap,proto3" json:"txIdsMap,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"varint,2,opt,name=value,proto3"`
}

func (m *TxBatch) Reset()         { *m = TxBatch{} }
func (m *TxBatch) String() string { return proto.CompactTextString(m) }
func (*TxBatch) ProtoMessage()    {}
func (*TxBatch) Descriptor() ([]byte, []int) {
	return fileDescriptor_0bc7127f197678cd, []int{1}
}
func (m *TxBatch) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *TxBatch) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_TxBatch.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *TxBatch) XXX_Merge(src proto.Message) {
	xxx_messageInfo_TxBatch.Merge(m, src)
}
func (m *TxBatch) XXX_Size() int {
	return m.Size()
}
func (m *TxBatch) XXX_DiscardUnknown() {
	xxx_messageInfo_TxBatch.DiscardUnknown(m)
}

var xxx_messageInfo_TxBatch proto.InternalMessageInfo

func (m *TxBatch) GetBatchId() int32 {
	if m != nil {
		return m.BatchId
	}
	return 0
}

func (m *TxBatch) GetNodeId() string {
	if m != nil {
		return m.NodeId
	}
	return ""
}

func (m *TxBatch) GetSize_() int32 {
	if m != nil {
		return m.Size_
	}
	return 0
}

func (m *TxBatch) GetTxs() []*common.Transaction {
	if m != nil {
		return m.Txs
	}
	return nil
}

func (m *TxBatch) GetTxIdsMap() map[string]int32 {
	if m != nil {
		return m.TxIdsMap
	}
	return nil
}

func init() {
	proto.RegisterEnum("txpool.SignalType", SignalType_name, SignalType_value)
	proto.RegisterType((*TxPoolSignal)(nil), "txpool.TxPoolSignal")
	proto.RegisterType((*TxBatch)(nil), "txpool.TxBatch")
	proto.RegisterMapType((map[string]int32)(nil), "txpool.TxBatch.TxIdsMapEntry")
}

func init() { proto.RegisterFile("txpool/transaction_pool.proto", fileDescriptor_0bc7127f197678cd) }

var fileDescriptor_0bc7127f197678cd = []byte{
	// 415 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x54, 0x51, 0x4d, 0x6b, 0xdb, 0x40,
	0x14, 0xf4, 0xda, 0x91, 0x13, 0xbf, 0x26, 0xc5, 0xdd, 0x96, 0xb0, 0x04, 0x22, 0x4c, 0xa0, 0x60,
	0x0a, 0x59, 0x81, 0x7b, 0xe9, 0xc7, 0x29, 0x36, 0x3a, 0x88, 0x36, 0x92, 0x59, 0x2f, 0x3d, 0x94,
	0x82, 0x59, 0x4b, 0xaa, 0x22, 0x2c, 0x6b, 0x85, 0xb4, 0x2d, 0x72, 0x7f, 0x45, 0x7f, 0x56, 0x8f,
	0x39, 0xf6, 0x58, 0xec, 0x7b, 0x7f, 0x43, 0xd1, 0x4a, 0xfe, 0xe8, 0x6d, 0xe6, 0xcd, 0xe8, 0xe9,
	0xcd, 0x2c, 0x5c, 0xab, 0x32, 0x93, 0x32, 0xb1, 0x54, 0x2e, 0xd2, 0x42, 0xf8, 0x2a, 0x96, 0xe9,
	0xbc, 0x1a, 0xd0, 0x2c, 0x97, 0x4a, 0xe2, 0x6e, 0x2d, 0x5f, 0x11, 0x5f, 0xae, 0x56, 0x32, 0x3d,
	0xb6, 0xd5, 0x8e, 0x9b, 0x2f, 0x70, 0xce, 0xcb, 0xa9, 0x94, 0xc9, 0x2c, 0x8e, 0x52, 0x91, 0xe0,
	0x11, 0x40, 0xa1, 0x11, 0x5f, 0x67, 0x21, 0x69, 0x0f, 0xd0, 0xf0, 0xe9, 0x08, 0xd3, 0x7a, 0x0d,
	0x9d, 0xed, 0x15, 0x76, 0xe4, 0xc2, 0x04, 0x4e, 0xfd, 0x07, 0x11, 0xa7, 0x4e, 0x40, 0x3a, 0x03,
	0x34, 0xec, 0xb1, 0x1d, 0xbd, 0xf9, 0x8b, 0xe0, 0x94, 0x97, 0x63, 0xa1, 0xfc, 0x87, 0xca, 0xb5,
	0xa8, 0x80, 0x13, 0x10, 0x34, 0x40, 0x43, 0x83, 0xed, 0x28, 0xbe, 0x84, 0x6e, 0x2a, 0x83, 0xd0,
	0x09, 0xf4, 0xff, 0x7a, 0xac, 0x61, 0x18, 0xc3, 0x49, 0x11, 0xff, 0x08, 0xf5, 0x52, 0x83, 0x69,
	0x8c, 0x5f, 0x42, 0x47, 0x95, 0x05, 0x39, 0x19, 0x74, 0x86, 0x4f, 0x46, 0xcf, 0x69, 0x9d, 0x8b,
	0xf2, 0x43, 0x2e, 0x56, 0xe9, 0xf8, 0x2d, 0x9c, 0xa9, 0xd2, 0x09, 0x8a, 0x7b, 0x91, 0x11, 0x43,
	0x7b, 0xaf, 0x77, 0x21, 0x9a, 0x7b, 0x28, 0x6f, 0x74, 0x3b, 0x55, 0xf9, 0x9a, 0xed, 0xed, 0x57,
	0xef, 0xe1, 0xe2, 0x3f, 0x09, 0xf7, 0xa1, 0xb3, 0x0c, 0xd7, 0xfa, 0xe8, 0x1e, 0xab, 0x20, 0x7e,
	0x01, 0xc6, 0x77, 0x91, 0x7c, 0xab, 0xfb, 0x31, 0x58, 0x4d, 0xde, 0xb5, 0xdf, 0xa0, 0x57, 0x36,
	0xc0, 0xa1, 0x24, 0x7c, 0x0e, 0x67, 0xae, 0x37, 0xb7, 0x3f, 0xd9, 0x2e, 0xef, 0xb7, 0xf0, 0x25,
	0x60, 0xce, 0xee, 0xdc, 0xd9, 0xdd, 0x84, 0x3b, 0x9e, 0x3b, 0x77, 0xdc, 0x89, 0x77, 0x6f, 0xf7,
	0x11, 0x7e, 0x06, 0x17, 0xe3, 0x8f, 0xde, 0xe4, 0xc3, 0x7c, 0xca, 0xbc, 0xa9, 0x37, 0xb3, 0xfb,
	0xed, 0xf1, 0xd7, 0x5f, 0x1b, 0x13, 0x3d, 0x6e, 0x4c, 0xf4, 0x67, 0x63, 0xa2, 0x9f, 0x5b, 0xb3,
	0xf5, 0xb8, 0x35, 0x5b, 0xbf, 0xb7, 0x66, 0x0b, 0x88, 0xcc, 0x23, 0xaa, 0xeb, 0x5d, 0x89, 0x65,
	0x98, 0xd3, 0x6c, 0xd1, 0xe4, 0xfa, 0x3c, 0x3a, 0x9a, 0xca, 0x3c, 0xb2, 0x0e, 0xf4, 0xb6, 0x08,
	0x96, 0xb7, 0x91, 0xb4, 0xb2, 0x85, 0xa5, 0x1f, 0x3d, 0x92, 0x56, 0xfd, 0xcd, 0xa2, 0xab, 0xf9,
	0xeb, 0x7f, 0x01, 0x00, 0x00, 0xff, 0xff, 0xe5, 0xd0, 0x0a, 0xd3, 0x47, 0x02, 0x00, 0x00,
}

func (m *TxPoolSignal) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *TxPoolSignal) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *TxPoolSignal) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.ChainId) > 0 {
		i -= len(m.ChainId)
		copy(dAtA[i:], m.ChainId)
		i = encodeVarintTransactionPool(dAtA, i, uint64(len(m.ChainId)))
		i--
		dAtA[i] = 0x1a
	}
	if m.SignalType != 0 {
		i = encodeVarintTransactionPool(dAtA, i, uint64(m.SignalType))
		i--
		dAtA[i] = 0x10
	}
	return len(dAtA) - i, nil
}

func (m *TxBatch) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *TxBatch) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *TxBatch) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.TxIdsMap) > 0 {
		for k := range m.TxIdsMap {
			v := m.TxIdsMap[k]
			baseI := i
			i = encodeVarintTransactionPool(dAtA, i, uint64(v))
			i--
			dAtA[i] = 0x10
			i -= len(k)
			copy(dAtA[i:], k)
			i = encodeVarintTransactionPool(dAtA, i, uint64(len(k)))
			i--
			dAtA[i] = 0xa
			i = encodeVarintTransactionPool(dAtA, i, uint64(baseI-i))
			i--
			dAtA[i] = 0x2a
		}
	}
	if len(m.Txs) > 0 {
		for iNdEx := len(m.Txs) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.Txs[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintTransactionPool(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0x22
		}
	}
	if m.Size_ != 0 {
		i = encodeVarintTransactionPool(dAtA, i, uint64(m.Size_))
		i--
		dAtA[i] = 0x18
	}
	if len(m.NodeId) > 0 {
		i -= len(m.NodeId)
		copy(dAtA[i:], m.NodeId)
		i = encodeVarintTransactionPool(dAtA, i, uint64(len(m.NodeId)))
		i--
		dAtA[i] = 0x12
	}
	if m.BatchId != 0 {
		i = encodeVarintTransactionPool(dAtA, i, uint64(m.BatchId))
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func encodeVarintTransactionPool(dAtA []byte, offset int, v uint64) int {
	offset -= sovTransactionPool(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *TxPoolSignal) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.SignalType != 0 {
		n += 1 + sovTransactionPool(uint64(m.SignalType))
	}
	l = len(m.ChainId)
	if l > 0 {
		n += 1 + l + sovTransactionPool(uint64(l))
	}
	return n
}

func (m *TxBatch) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.BatchId != 0 {
		n += 1 + sovTransactionPool(uint64(m.BatchId))
	}
	l = len(m.NodeId)
	if l > 0 {
		n += 1 + l + sovTransactionPool(uint64(l))
	}
	if m.Size_ != 0 {
		n += 1 + sovTransactionPool(uint64(m.Size_))
	}
	if len(m.Txs) > 0 {
		for _, e := range m.Txs {
			l = e.Size()
			n += 1 + l + sovTransactionPool(uint64(l))
		}
	}
	if len(m.TxIdsMap) > 0 {
		for k, v := range m.TxIdsMap {
			_ = k
			_ = v
			mapEntrySize := 1 + len(k) + sovTransactionPool(uint64(len(k))) + 1 + sovTransactionPool(uint64(v))
			n += mapEntrySize + 1 + sovTransactionPool(uint64(mapEntrySize))
		}
	}
	return n
}

func sovTransactionPool(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozTransactionPool(x uint64) (n int) {
	return sovTransactionPool(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *TxPoolSignal) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowTransactionPool
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
			return fmt.Errorf("proto: TxPoolSignal: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: TxPoolSignal: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 2:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field SignalType", wireType)
			}
			m.SignalType = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTransactionPool
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.SignalType |= SignalType(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field ChainId", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTransactionPool
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
				return ErrInvalidLengthTransactionPool
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthTransactionPool
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.ChainId = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipTransactionPool(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthTransactionPool
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
func (m *TxBatch) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowTransactionPool
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
			return fmt.Errorf("proto: TxBatch: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: TxBatch: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field BatchId", wireType)
			}
			m.BatchId = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTransactionPool
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.BatchId |= int32(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field NodeId", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTransactionPool
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
				return ErrInvalidLengthTransactionPool
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthTransactionPool
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.NodeId = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 3:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Size_", wireType)
			}
			m.Size_ = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTransactionPool
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Size_ |= int32(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Txs", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTransactionPool
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
				return ErrInvalidLengthTransactionPool
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthTransactionPool
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Txs = append(m.Txs, &common.Transaction{})
			if err := m.Txs[len(m.Txs)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 5:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field TxIdsMap", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTransactionPool
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
				return ErrInvalidLengthTransactionPool
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthTransactionPool
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.TxIdsMap == nil {
				m.TxIdsMap = make(map[string]int32)
			}
			var mapkey string
			var mapvalue int32
			for iNdEx < postIndex {
				entryPreIndex := iNdEx
				var wire uint64
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return ErrIntOverflowTransactionPool
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
				if fieldNum == 1 {
					var stringLenmapkey uint64
					for shift := uint(0); ; shift += 7 {
						if shift >= 64 {
							return ErrIntOverflowTransactionPool
						}
						if iNdEx >= l {
							return io.ErrUnexpectedEOF
						}
						b := dAtA[iNdEx]
						iNdEx++
						stringLenmapkey |= uint64(b&0x7F) << shift
						if b < 0x80 {
							break
						}
					}
					intStringLenmapkey := int(stringLenmapkey)
					if intStringLenmapkey < 0 {
						return ErrInvalidLengthTransactionPool
					}
					postStringIndexmapkey := iNdEx + intStringLenmapkey
					if postStringIndexmapkey < 0 {
						return ErrInvalidLengthTransactionPool
					}
					if postStringIndexmapkey > l {
						return io.ErrUnexpectedEOF
					}
					mapkey = string(dAtA[iNdEx:postStringIndexmapkey])
					iNdEx = postStringIndexmapkey
				} else if fieldNum == 2 {
					for shift := uint(0); ; shift += 7 {
						if shift >= 64 {
							return ErrIntOverflowTransactionPool
						}
						if iNdEx >= l {
							return io.ErrUnexpectedEOF
						}
						b := dAtA[iNdEx]
						iNdEx++
						mapvalue |= int32(b&0x7F) << shift
						if b < 0x80 {
							break
						}
					}
				} else {
					iNdEx = entryPreIndex
					skippy, err := skipTransactionPool(dAtA[iNdEx:])
					if err != nil {
						return err
					}
					if (skippy < 0) || (iNdEx+skippy) < 0 {
						return ErrInvalidLengthTransactionPool
					}
					if (iNdEx + skippy) > postIndex {
						return io.ErrUnexpectedEOF
					}
					iNdEx += skippy
				}
			}
			m.TxIdsMap[mapkey] = mapvalue
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipTransactionPool(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthTransactionPool
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
func skipTransactionPool(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowTransactionPool
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
					return 0, ErrIntOverflowTransactionPool
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
					return 0, ErrIntOverflowTransactionPool
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
				return 0, ErrInvalidLengthTransactionPool
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupTransactionPool
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthTransactionPool
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthTransactionPool        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowTransactionPool          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupTransactionPool = fmt.Errorf("proto: unexpected end of group")
)
