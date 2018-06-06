// Copyright (c) 2018 Uber Technologies, Inc.
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: github.com/m3db/m3db/src/coordinator/generated/proto/admin/database.proto

/*
	Package admin is a generated protocol buffer package.

	It is generated from these files:
		github.com/m3db/m3db/src/coordinator/generated/proto/admin/database.proto
		github.com/m3db/m3db/src/coordinator/generated/proto/admin/namespace.proto
		github.com/m3db/m3db/src/coordinator/generated/proto/admin/placement.proto

	It has these top-level messages:
		DatabaseCreateRequest
		DatabaseCreateResponse
		NamespaceGetResponse
		NamespaceAddRequest
		PlacementInitRequest
		PlacementGetResponse
		PlacementAddRequest
*/
package admin

import proto "github.com/gogo/protobuf/proto"
import fmt "fmt"
import math "math"

import io "io"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.GoGoProtoPackageIsVersion2 // please upgrade the proto package

type DatabaseCreateRequest struct {
	// Required fields
	NamespaceName string `protobuf:"bytes,1,opt,name=namespace_name,json=namespaceName,proto3" json:"namespace_name,omitempty"`
	Type          string `protobuf:"bytes,2,opt,name=type,proto3" json:"type,omitempty"`
	// Optional fields that may be inferred depending on database type
	NumShards         int32 `protobuf:"varint,3,opt,name=num_shards,json=numShards,proto3" json:"num_shards,omitempty"`
	ReplicationFactor int32 `protobuf:"varint,4,opt,name=replication_factor,json=replicationFactor,proto3" json:"replication_factor,omitempty"`
	// The below two options are used to default retention options
	RetentionPeriodNanos            int64 `protobuf:"varint,5,opt,name=retention_period_nanos,json=retentionPeriodNanos,proto3" json:"retention_period_nanos,omitempty"`
	ExpectedSeriesDatapointsPerHour int64 `protobuf:"varint,6,opt,name=expected_series_datapoints_per_hour,json=expectedSeriesDatapointsPerHour,proto3" json:"expected_series_datapoints_per_hour,omitempty"`
}

func (m *DatabaseCreateRequest) Reset()                    { *m = DatabaseCreateRequest{} }
func (m *DatabaseCreateRequest) String() string            { return proto.CompactTextString(m) }
func (*DatabaseCreateRequest) ProtoMessage()               {}
func (*DatabaseCreateRequest) Descriptor() ([]byte, []int) { return fileDescriptorDatabase, []int{0} }

func (m *DatabaseCreateRequest) GetNamespaceName() string {
	if m != nil {
		return m.NamespaceName
	}
	return ""
}

func (m *DatabaseCreateRequest) GetType() string {
	if m != nil {
		return m.Type
	}
	return ""
}

func (m *DatabaseCreateRequest) GetNumShards() int32 {
	if m != nil {
		return m.NumShards
	}
	return 0
}

func (m *DatabaseCreateRequest) GetReplicationFactor() int32 {
	if m != nil {
		return m.ReplicationFactor
	}
	return 0
}

func (m *DatabaseCreateRequest) GetRetentionPeriodNanos() int64 {
	if m != nil {
		return m.RetentionPeriodNanos
	}
	return 0
}

func (m *DatabaseCreateRequest) GetExpectedSeriesDatapointsPerHour() int64 {
	if m != nil {
		return m.ExpectedSeriesDatapointsPerHour
	}
	return 0
}

type DatabaseCreateResponse struct {
	Namespace *NamespaceGetResponse `protobuf:"bytes,1,opt,name=namespace" json:"namespace,omitempty"`
	Placement *PlacementGetResponse `protobuf:"bytes,2,opt,name=placement" json:"placement,omitempty"`
}

func (m *DatabaseCreateResponse) Reset()                    { *m = DatabaseCreateResponse{} }
func (m *DatabaseCreateResponse) String() string            { return proto.CompactTextString(m) }
func (*DatabaseCreateResponse) ProtoMessage()               {}
func (*DatabaseCreateResponse) Descriptor() ([]byte, []int) { return fileDescriptorDatabase, []int{1} }

func (m *DatabaseCreateResponse) GetNamespace() *NamespaceGetResponse {
	if m != nil {
		return m.Namespace
	}
	return nil
}

func (m *DatabaseCreateResponse) GetPlacement() *PlacementGetResponse {
	if m != nil {
		return m.Placement
	}
	return nil
}

func init() {
	proto.RegisterType((*DatabaseCreateRequest)(nil), "admin.DatabaseCreateRequest")
	proto.RegisterType((*DatabaseCreateResponse)(nil), "admin.DatabaseCreateResponse")
}
func (m *DatabaseCreateRequest) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *DatabaseCreateRequest) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	if len(m.NamespaceName) > 0 {
		dAtA[i] = 0xa
		i++
		i = encodeVarintDatabase(dAtA, i, uint64(len(m.NamespaceName)))
		i += copy(dAtA[i:], m.NamespaceName)
	}
	if len(m.Type) > 0 {
		dAtA[i] = 0x12
		i++
		i = encodeVarintDatabase(dAtA, i, uint64(len(m.Type)))
		i += copy(dAtA[i:], m.Type)
	}
	if m.NumShards != 0 {
		dAtA[i] = 0x18
		i++
		i = encodeVarintDatabase(dAtA, i, uint64(m.NumShards))
	}
	if m.ReplicationFactor != 0 {
		dAtA[i] = 0x20
		i++
		i = encodeVarintDatabase(dAtA, i, uint64(m.ReplicationFactor))
	}
	if m.RetentionPeriodNanos != 0 {
		dAtA[i] = 0x28
		i++
		i = encodeVarintDatabase(dAtA, i, uint64(m.RetentionPeriodNanos))
	}
	if m.ExpectedSeriesDatapointsPerHour != 0 {
		dAtA[i] = 0x30
		i++
		i = encodeVarintDatabase(dAtA, i, uint64(m.ExpectedSeriesDatapointsPerHour))
	}
	return i, nil
}

func (m *DatabaseCreateResponse) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *DatabaseCreateResponse) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	if m.Namespace != nil {
		dAtA[i] = 0xa
		i++
		i = encodeVarintDatabase(dAtA, i, uint64(m.Namespace.Size()))
		n1, err := m.Namespace.MarshalTo(dAtA[i:])
		if err != nil {
			return 0, err
		}
		i += n1
	}
	if m.Placement != nil {
		dAtA[i] = 0x12
		i++
		i = encodeVarintDatabase(dAtA, i, uint64(m.Placement.Size()))
		n2, err := m.Placement.MarshalTo(dAtA[i:])
		if err != nil {
			return 0, err
		}
		i += n2
	}
	return i, nil
}

func encodeVarintDatabase(dAtA []byte, offset int, v uint64) int {
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return offset + 1
}
func (m *DatabaseCreateRequest) Size() (n int) {
	var l int
	_ = l
	l = len(m.NamespaceName)
	if l > 0 {
		n += 1 + l + sovDatabase(uint64(l))
	}
	l = len(m.Type)
	if l > 0 {
		n += 1 + l + sovDatabase(uint64(l))
	}
	if m.NumShards != 0 {
		n += 1 + sovDatabase(uint64(m.NumShards))
	}
	if m.ReplicationFactor != 0 {
		n += 1 + sovDatabase(uint64(m.ReplicationFactor))
	}
	if m.RetentionPeriodNanos != 0 {
		n += 1 + sovDatabase(uint64(m.RetentionPeriodNanos))
	}
	if m.ExpectedSeriesDatapointsPerHour != 0 {
		n += 1 + sovDatabase(uint64(m.ExpectedSeriesDatapointsPerHour))
	}
	return n
}

func (m *DatabaseCreateResponse) Size() (n int) {
	var l int
	_ = l
	if m.Namespace != nil {
		l = m.Namespace.Size()
		n += 1 + l + sovDatabase(uint64(l))
	}
	if m.Placement != nil {
		l = m.Placement.Size()
		n += 1 + l + sovDatabase(uint64(l))
	}
	return n
}

func sovDatabase(x uint64) (n int) {
	for {
		n++
		x >>= 7
		if x == 0 {
			break
		}
	}
	return n
}
func sozDatabase(x uint64) (n int) {
	return sovDatabase(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *DatabaseCreateRequest) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowDatabase
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: DatabaseCreateRequest: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: DatabaseCreateRequest: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field NamespaceName", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowDatabase
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= (uint64(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthDatabase
			}
			postIndex := iNdEx + intStringLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.NamespaceName = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Type", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowDatabase
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= (uint64(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthDatabase
			}
			postIndex := iNdEx + intStringLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Type = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 3:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field NumShards", wireType)
			}
			m.NumShards = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowDatabase
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.NumShards |= (int32(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 4:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field ReplicationFactor", wireType)
			}
			m.ReplicationFactor = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowDatabase
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.ReplicationFactor |= (int32(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 5:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field RetentionPeriodNanos", wireType)
			}
			m.RetentionPeriodNanos = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowDatabase
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.RetentionPeriodNanos |= (int64(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 6:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field ExpectedSeriesDatapointsPerHour", wireType)
			}
			m.ExpectedSeriesDatapointsPerHour = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowDatabase
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.ExpectedSeriesDatapointsPerHour |= (int64(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		default:
			iNdEx = preIndex
			skippy, err := skipDatabase(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthDatabase
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
func (m *DatabaseCreateResponse) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowDatabase
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: DatabaseCreateResponse: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: DatabaseCreateResponse: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Namespace", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowDatabase
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthDatabase
			}
			postIndex := iNdEx + msglen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.Namespace == nil {
				m.Namespace = &NamespaceGetResponse{}
			}
			if err := m.Namespace.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Placement", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowDatabase
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthDatabase
			}
			postIndex := iNdEx + msglen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.Placement == nil {
				m.Placement = &PlacementGetResponse{}
			}
			if err := m.Placement.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipDatabase(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthDatabase
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
func skipDatabase(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowDatabase
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
					return 0, ErrIntOverflowDatabase
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				iNdEx++
				if dAtA[iNdEx-1] < 0x80 {
					break
				}
			}
			return iNdEx, nil
		case 1:
			iNdEx += 8
			return iNdEx, nil
		case 2:
			var length int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowDatabase
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
			iNdEx += length
			if length < 0 {
				return 0, ErrInvalidLengthDatabase
			}
			return iNdEx, nil
		case 3:
			for {
				var innerWire uint64
				var start int = iNdEx
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return 0, ErrIntOverflowDatabase
					}
					if iNdEx >= l {
						return 0, io.ErrUnexpectedEOF
					}
					b := dAtA[iNdEx]
					iNdEx++
					innerWire |= (uint64(b) & 0x7F) << shift
					if b < 0x80 {
						break
					}
				}
				innerWireType := int(innerWire & 0x7)
				if innerWireType == 4 {
					break
				}
				next, err := skipDatabase(dAtA[start:])
				if err != nil {
					return 0, err
				}
				iNdEx = start + next
			}
			return iNdEx, nil
		case 4:
			return iNdEx, nil
		case 5:
			iNdEx += 4
			return iNdEx, nil
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
	}
	panic("unreachable")
}

var (
	ErrInvalidLengthDatabase = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowDatabase   = fmt.Errorf("proto: integer overflow")
)

func init() {
	proto.RegisterFile("github.com/m3db/m3db/src/coordinator/generated/proto/admin/database.proto", fileDescriptorDatabase)
}

var fileDescriptorDatabase = []byte{
	// 388 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xa4, 0x91, 0xdf, 0x6a, 0xdb, 0x30,
	0x14, 0xc6, 0xa7, 0xfc, 0x83, 0x68, 0x6c, 0x6c, 0x62, 0x0b, 0x66, 0x63, 0x5e, 0xc8, 0x18, 0xe4,
	0x66, 0x36, 0x2c, 0xbb, 0xd9, 0xed, 0x16, 0xb6, 0xb6, 0x94, 0x10, 0x9c, 0x07, 0x30, 0xb2, 0x74,
	0x9a, 0x08, 0xa2, 0x3f, 0x95, 0x64, 0x68, 0x5f, 0xa2, 0xf4, 0x21, 0xfa, 0x30, 0xbd, 0xec, 0x23,
	0x94, 0xf4, 0x45, 0x8a, 0x95, 0xd8, 0x69, 0x7b, 0x9b, 0x1b, 0x73, 0xf8, 0xbe, 0xdf, 0xf9, 0x2c,
	0xbe, 0x83, 0x8f, 0x97, 0xc2, 0xaf, 0xca, 0x22, 0x61, 0x5a, 0xa6, 0x72, 0xc2, 0x8b, 0xed, 0xc7,
	0x59, 0x96, 0x32, 0xad, 0x2d, 0x17, 0x8a, 0x7a, 0x6d, 0xd3, 0x25, 0x28, 0xb0, 0xd4, 0x03, 0x4f,
	0x8d, 0xd5, 0x5e, 0xa7, 0x94, 0x4b, 0xa1, 0x52, 0x4e, 0x3d, 0x2d, 0xa8, 0x83, 0x24, 0x88, 0xa4,
	0x1b, 0xd4, 0x4f, 0x27, 0x07, 0x24, 0x2a, 0x2a, 0xc1, 0x19, 0xca, 0x76, 0x91, 0x07, 0x65, 0x99,
	0x35, 0x65, 0x20, 0x41, 0xf9, 0x6d, 0xd6, 0xe8, 0xa6, 0x85, 0x3f, 0x4e, 0x77, 0x2f, 0xfe, 0x6b,
	0x81, 0x7a, 0xc8, 0xe0, 0xbc, 0x04, 0xe7, 0xc9, 0x77, 0xfc, 0xb6, 0xf9, 0x71, 0x5e, 0x4d, 0x11,
	0x1a, 0xa2, 0x71, 0x3f, 0x7b, 0xd3, 0xa8, 0x33, 0x2a, 0x81, 0x10, 0xdc, 0xf1, 0x97, 0x06, 0xa2,
	0x56, 0x30, 0xc3, 0x4c, 0xbe, 0x60, 0xac, 0x4a, 0x99, 0xbb, 0x15, 0xb5, 0xdc, 0x45, 0xed, 0x21,
	0x1a, 0x77, 0xb3, 0xbe, 0x2a, 0xe5, 0x22, 0x08, 0xe4, 0x07, 0x26, 0x16, 0xcc, 0x5a, 0x30, 0xea,
	0x85, 0x56, 0xf9, 0x19, 0x65, 0x5e, 0xdb, 0xa8, 0x13, 0xb0, 0xf7, 0x4f, 0x9c, 0x7f, 0xc1, 0x20,
	0xbf, 0xf0, 0xc0, 0x82, 0x07, 0x15, 0x60, 0x03, 0x56, 0x68, 0x9e, 0x2b, 0xaa, 0xb4, 0x8b, 0xba,
	0x43, 0x34, 0x6e, 0x67, 0x1f, 0x1a, 0x77, 0x1e, 0xcc, 0x59, 0xe5, 0x91, 0x53, 0xfc, 0x0d, 0x2e,
	0x0c, 0x30, 0x0f, 0x3c, 0x77, 0x60, 0x05, 0xb8, 0xbc, 0xba, 0x8c, 0xd1, 0x42, 0x79, 0x57, 0xc5,
	0xe4, 0x2b, 0x5d, 0xda, 0xa8, 0x17, 0x22, 0xbe, 0xd6, 0xe8, 0x22, 0x90, 0xd3, 0x06, 0x9c, 0x83,
	0x3d, 0xd2, 0xa5, 0x1d, 0x5d, 0x21, 0x3c, 0x78, 0x59, 0x93, 0x33, 0x5a, 0x39, 0x20, 0xbf, 0x71,
	0xbf, 0x69, 0x24, 0x54, 0xf4, 0xfa, 0xe7, 0xe7, 0x24, 0x94, 0x9d, 0xcc, 0x6a, 0xfd, 0x3f, 0xf8,
	0x9a, 0xcf, 0xf6, 0x74, 0xb5, 0xda, 0xdc, 0x23, 0x14, 0xb8, 0x5f, 0x9d, 0xd7, 0xfa, 0xb3, 0xd5,
	0x86, 0xfe, 0xf3, 0xee, 0x76, 0x13, 0xa3, 0xbb, 0x4d, 0x8c, 0xee, 0x37, 0x31, 0xba, 0x7e, 0x88,
	0x5f, 0x15, 0xbd, 0x70, 0xd0, 0xc9, 0x63, 0x00, 0x00, 0x00, 0xff, 0xff, 0x99, 0x16, 0x86, 0x4e,
	0xbc, 0x02, 0x00, 0x00,
}
