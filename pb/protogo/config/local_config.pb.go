// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: config/local_config.proto

package config

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

// rquest for debug configuration
type DebugConfigRequest struct {
	Pairs []*common.KeyValuePair `protobuf:"bytes,1,rep,name=pairs,proto3" json:"pairs,omitempty"`
}

func (m *DebugConfigRequest) Reset()         { *m = DebugConfigRequest{} }
func (m *DebugConfigRequest) String() string { return proto.CompactTextString(m) }
func (*DebugConfigRequest) ProtoMessage()    {}
func (*DebugConfigRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_cd062d6c1b6d65a6, []int{0}
}
func (m *DebugConfigRequest) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *DebugConfigRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_DebugConfigRequest.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *DebugConfigRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_DebugConfigRequest.Merge(m, src)
}
func (m *DebugConfigRequest) XXX_Size() int {
	return m.Size()
}
func (m *DebugConfigRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_DebugConfigRequest.DiscardUnknown(m)
}

var xxx_messageInfo_DebugConfigRequest proto.InternalMessageInfo

func (m *DebugConfigRequest) GetPairs() []*common.KeyValuePair {
	if m != nil {
		return m.Pairs
	}
	return nil
}

// Rrsponse for debug configuration
type DebugConfigResponse struct {
	// 0 success
	// 1 fail
	Code int32 `protobuf:"varint,1,opt,name=code,proto3" json:"code,omitempty"`
	// failure message
	Message string `protobuf:"bytes,2,opt,name=message,proto3" json:"message,omitempty"`
}

func (m *DebugConfigResponse) Reset()         { *m = DebugConfigResponse{} }
func (m *DebugConfigResponse) String() string { return proto.CompactTextString(m) }
func (*DebugConfigResponse) ProtoMessage()    {}
func (*DebugConfigResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_cd062d6c1b6d65a6, []int{1}
}
func (m *DebugConfigResponse) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *DebugConfigResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_DebugConfigResponse.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *DebugConfigResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_DebugConfigResponse.Merge(m, src)
}
func (m *DebugConfigResponse) XXX_Size() int {
	return m.Size()
}
func (m *DebugConfigResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_DebugConfigResponse.DiscardUnknown(m)
}

var xxx_messageInfo_DebugConfigResponse proto.InternalMessageInfo

func (m *DebugConfigResponse) GetCode() int32 {
	if m != nil {
		return m.Code
	}
	return 0
}

func (m *DebugConfigResponse) GetMessage() string {
	if m != nil {
		return m.Message
	}
	return ""
}

// request for check new block configuration
type CheckNewBlockChainConfigRequest struct {
}

func (m *CheckNewBlockChainConfigRequest) Reset()         { *m = CheckNewBlockChainConfigRequest{} }
func (m *CheckNewBlockChainConfigRequest) String() string { return proto.CompactTextString(m) }
func (*CheckNewBlockChainConfigRequest) ProtoMessage()    {}
func (*CheckNewBlockChainConfigRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_cd062d6c1b6d65a6, []int{2}
}
func (m *CheckNewBlockChainConfigRequest) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *CheckNewBlockChainConfigRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_CheckNewBlockChainConfigRequest.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *CheckNewBlockChainConfigRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_CheckNewBlockChainConfigRequest.Merge(m, src)
}
func (m *CheckNewBlockChainConfigRequest) XXX_Size() int {
	return m.Size()
}
func (m *CheckNewBlockChainConfigRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_CheckNewBlockChainConfigRequest.DiscardUnknown(m)
}

var xxx_messageInfo_CheckNewBlockChainConfigRequest proto.InternalMessageInfo

// response for check new block configuration
type CheckNewBlockChainConfigResponse struct {
	// 0 success
	// 1 fail
	Code int32 `protobuf:"varint,1,opt,name=code,proto3" json:"code,omitempty"`
	// failure message
	Message string `protobuf:"bytes,2,opt,name=message,proto3" json:"message,omitempty"`
}

func (m *CheckNewBlockChainConfigResponse) Reset()         { *m = CheckNewBlockChainConfigResponse{} }
func (m *CheckNewBlockChainConfigResponse) String() string { return proto.CompactTextString(m) }
func (*CheckNewBlockChainConfigResponse) ProtoMessage()    {}
func (*CheckNewBlockChainConfigResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_cd062d6c1b6d65a6, []int{3}
}
func (m *CheckNewBlockChainConfigResponse) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *CheckNewBlockChainConfigResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_CheckNewBlockChainConfigResponse.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *CheckNewBlockChainConfigResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_CheckNewBlockChainConfigResponse.Merge(m, src)
}
func (m *CheckNewBlockChainConfigResponse) XXX_Size() int {
	return m.Size()
}
func (m *CheckNewBlockChainConfigResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_CheckNewBlockChainConfigResponse.DiscardUnknown(m)
}

var xxx_messageInfo_CheckNewBlockChainConfigResponse proto.InternalMessageInfo

func (m *CheckNewBlockChainConfigResponse) GetCode() int32 {
	if m != nil {
		return m.Code
	}
	return 0
}

func (m *CheckNewBlockChainConfigResponse) GetMessage() string {
	if m != nil {
		return m.Message
	}
	return ""
}

func init() {
	proto.RegisterType((*DebugConfigRequest)(nil), "config.DebugConfigRequest")
	proto.RegisterType((*DebugConfigResponse)(nil), "config.DebugConfigResponse")
	proto.RegisterType((*CheckNewBlockChainConfigRequest)(nil), "config.CheckNewBlockChainConfigRequest")
	proto.RegisterType((*CheckNewBlockChainConfigResponse)(nil), "config.CheckNewBlockChainConfigResponse")
}

func init() { proto.RegisterFile("config/local_config.proto", fileDescriptor_cd062d6c1b6d65a6) }

var fileDescriptor_cd062d6c1b6d65a6 = []byte{
	// 278 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x9c, 0x90, 0xbf, 0x4e, 0xc3, 0x30,
	0x18, 0xc4, 0x6b, 0xa0, 0x45, 0x98, 0xcd, 0x74, 0x08, 0x0c, 0x26, 0x64, 0xaa, 0x90, 0x6a, 0x4b,
	0xe5, 0x05, 0x50, 0xc3, 0x86, 0x84, 0xaa, 0x0c, 0x0c, 0x2c, 0xc8, 0x71, 0xbf, 0xba, 0x51, 0xfe,
	0x7c, 0xc1, 0x6e, 0x84, 0x78, 0x0b, 0x1e, 0x8b, 0xb1, 0x23, 0x23, 0x4a, 0x5e, 0x04, 0x11, 0x0b,
	0x41, 0x07, 0x16, 0x36, 0xff, 0x7c, 0x77, 0x3e, 0xeb, 0xe8, 0xa9, 0xc6, 0x6a, 0x95, 0x19, 0x59,
	0xa0, 0x56, 0xc5, 0xa3, 0x07, 0x51, 0x5b, 0xdc, 0x20, 0x1b, 0x79, 0x3a, 0x1b, 0x6b, 0x2c, 0x4b,
	0xac, 0xa4, 0x85, 0xa7, 0x06, 0xdc, 0xc6, 0xab, 0xd1, 0x35, 0x65, 0x37, 0x90, 0x36, 0x26, 0xee,
	0x4d, 0x89, 0xd7, 0xd8, 0x25, 0x1d, 0xd6, 0x2a, 0xb3, 0x2e, 0x20, 0xe1, 0xfe, 0xe4, 0x78, 0x36,
	0x16, 0x3e, 0x2b, 0x6e, 0xe1, 0xe5, 0x5e, 0x15, 0x0d, 0x2c, 0x54, 0x66, 0x13, 0x6f, 0x89, 0x62,
	0x7a, 0xb2, 0xf3, 0x82, 0xab, 0xb1, 0x72, 0xc0, 0x18, 0x3d, 0xd0, 0xb8, 0x84, 0x80, 0x84, 0x64,
	0x32, 0x4c, 0xfa, 0x33, 0x0b, 0xe8, 0x61, 0x09, 0xce, 0x29, 0x03, 0xc1, 0x5e, 0x48, 0x26, 0x47,
	0xc9, 0x37, 0x46, 0x17, 0xf4, 0x3c, 0x5e, 0x83, 0xce, 0xef, 0xe0, 0x79, 0x5e, 0xa0, 0xce, 0xe3,
	0xb5, 0xca, 0xaa, 0x9d, 0x3f, 0x45, 0x0b, 0x1a, 0xfe, 0x6d, 0xf9, 0x4f, 0xe9, 0x7c, 0xf5, 0xd6,
	0x72, 0xb2, 0x6d, 0x39, 0xf9, 0x68, 0x39, 0x79, 0xed, 0xf8, 0x60, 0xdb, 0xf1, 0xc1, 0x7b, 0xc7,
	0x07, 0x34, 0x40, 0x6b, 0x84, 0xfe, 0x2a, 0x28, 0x55, 0x0e, 0x56, 0xd4, 0xa9, 0xf0, 0x2b, 0x3e,
	0xcc, 0x7e, 0xdd, 0xa2, 0x35, 0xf2, 0x07, 0xa7, 0x6e, 0x99, 0x4f, 0x0d, 0xca, 0x3a, 0x95, 0xfd,
	0xb4, 0x06, 0xa5, 0xcf, 0xa4, 0xa3, 0x9e, 0xaf, 0x3e, 0x03, 0x00, 0x00, 0xff, 0xff, 0xee, 0x13,
	0x2c, 0x79, 0xa5, 0x01, 0x00, 0x00,
}

func (m *DebugConfigRequest) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *DebugConfigRequest) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *DebugConfigRequest) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.Pairs) > 0 {
		for iNdEx := len(m.Pairs) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.Pairs[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintLocalConfig(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0xa
		}
	}
	return len(dAtA) - i, nil
}

func (m *DebugConfigResponse) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *DebugConfigResponse) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *DebugConfigResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.Message) > 0 {
		i -= len(m.Message)
		copy(dAtA[i:], m.Message)
		i = encodeVarintLocalConfig(dAtA, i, uint64(len(m.Message)))
		i--
		dAtA[i] = 0x12
	}
	if m.Code != 0 {
		i = encodeVarintLocalConfig(dAtA, i, uint64(m.Code))
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func (m *CheckNewBlockChainConfigRequest) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *CheckNewBlockChainConfigRequest) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *CheckNewBlockChainConfigRequest) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	return len(dAtA) - i, nil
}

func (m *CheckNewBlockChainConfigResponse) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *CheckNewBlockChainConfigResponse) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *CheckNewBlockChainConfigResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.Message) > 0 {
		i -= len(m.Message)
		copy(dAtA[i:], m.Message)
		i = encodeVarintLocalConfig(dAtA, i, uint64(len(m.Message)))
		i--
		dAtA[i] = 0x12
	}
	if m.Code != 0 {
		i = encodeVarintLocalConfig(dAtA, i, uint64(m.Code))
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func encodeVarintLocalConfig(dAtA []byte, offset int, v uint64) int {
	offset -= sovLocalConfig(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *DebugConfigRequest) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if len(m.Pairs) > 0 {
		for _, e := range m.Pairs {
			l = e.Size()
			n += 1 + l + sovLocalConfig(uint64(l))
		}
	}
	return n
}

func (m *DebugConfigResponse) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.Code != 0 {
		n += 1 + sovLocalConfig(uint64(m.Code))
	}
	l = len(m.Message)
	if l > 0 {
		n += 1 + l + sovLocalConfig(uint64(l))
	}
	return n
}

func (m *CheckNewBlockChainConfigRequest) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	return n
}

func (m *CheckNewBlockChainConfigResponse) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.Code != 0 {
		n += 1 + sovLocalConfig(uint64(m.Code))
	}
	l = len(m.Message)
	if l > 0 {
		n += 1 + l + sovLocalConfig(uint64(l))
	}
	return n
}

func sovLocalConfig(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozLocalConfig(x uint64) (n int) {
	return sovLocalConfig(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *DebugConfigRequest) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowLocalConfig
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
			return fmt.Errorf("proto: DebugConfigRequest: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: DebugConfigRequest: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Pairs", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowLocalConfig
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
				return ErrInvalidLengthLocalConfig
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthLocalConfig
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Pairs = append(m.Pairs, &common.KeyValuePair{})
			if err := m.Pairs[len(m.Pairs)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipLocalConfig(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthLocalConfig
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
func (m *DebugConfigResponse) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowLocalConfig
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
			return fmt.Errorf("proto: DebugConfigResponse: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: DebugConfigResponse: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Code", wireType)
			}
			m.Code = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowLocalConfig
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Code |= int32(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Message", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowLocalConfig
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
				return ErrInvalidLengthLocalConfig
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthLocalConfig
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Message = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipLocalConfig(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthLocalConfig
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
func (m *CheckNewBlockChainConfigRequest) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowLocalConfig
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
			return fmt.Errorf("proto: CheckNewBlockChainConfigRequest: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: CheckNewBlockChainConfigRequest: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		default:
			iNdEx = preIndex
			skippy, err := skipLocalConfig(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthLocalConfig
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
func (m *CheckNewBlockChainConfigResponse) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowLocalConfig
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
			return fmt.Errorf("proto: CheckNewBlockChainConfigResponse: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: CheckNewBlockChainConfigResponse: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Code", wireType)
			}
			m.Code = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowLocalConfig
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Code |= int32(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Message", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowLocalConfig
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
				return ErrInvalidLengthLocalConfig
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthLocalConfig
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Message = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipLocalConfig(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthLocalConfig
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
func skipLocalConfig(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowLocalConfig
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
					return 0, ErrIntOverflowLocalConfig
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
					return 0, ErrIntOverflowLocalConfig
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
				return 0, ErrInvalidLengthLocalConfig
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupLocalConfig
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthLocalConfig
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthLocalConfig        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowLocalConfig          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupLocalConfig = fmt.Errorf("proto: unexpected end of group")
)
