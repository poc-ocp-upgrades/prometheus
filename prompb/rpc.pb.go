package prompb

import proto "github.com/gogo/protobuf/proto"
import fmt "fmt"
import math "math"
import _ "google.golang.org/genproto/googleapis/api/annotations"
import time "time"
import (
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
)
import github_com_gogo_protobuf_types "github.com/gogo/protobuf/types"
import io "io"

var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf
var _ = time.Kitchen

const _ = proto.GoGoProtoPackageIsVersion2

type TSDBSnapshotRequest struct {
	SkipHead		bool		`protobuf:"varint,1,opt,name=skip_head,json=skipHead,proto3" json:"skip_head,omitempty"`
	XXX_NoUnkeyedLiteral	struct{}	`json:"-"`
	XXX_unrecognized	[]byte		`json:"-"`
	XXX_sizecache		int32		`json:"-"`
}

func (m *TSDBSnapshotRequest) Reset() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	*m = TSDBSnapshotRequest{}
}
func (m *TSDBSnapshotRequest) String() string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return proto.CompactTextString(m)
}
func (*TSDBSnapshotRequest) ProtoMessage() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
}
func (*TSDBSnapshotRequest) Descriptor() ([]byte, []int) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return fileDescriptor_rpc_e0d54cfadc26b2e1, []int{0}
}
func (m *TSDBSnapshotRequest) XXX_Unmarshal(b []byte) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return m.Unmarshal(b)
}
func (m *TSDBSnapshotRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	if deterministic {
		return xxx_messageInfo_TSDBSnapshotRequest.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalTo(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (dst *TSDBSnapshotRequest) XXX_Merge(src proto.Message) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	xxx_messageInfo_TSDBSnapshotRequest.Merge(dst, src)
}
func (m *TSDBSnapshotRequest) XXX_Size() int {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return m.Size()
}
func (m *TSDBSnapshotRequest) XXX_DiscardUnknown() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	xxx_messageInfo_TSDBSnapshotRequest.DiscardUnknown(m)
}

var xxx_messageInfo_TSDBSnapshotRequest proto.InternalMessageInfo

type TSDBSnapshotResponse struct {
	Name			string		`protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	XXX_NoUnkeyedLiteral	struct{}	`json:"-"`
	XXX_unrecognized	[]byte		`json:"-"`
	XXX_sizecache		int32		`json:"-"`
}

func (m *TSDBSnapshotResponse) Reset() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	*m = TSDBSnapshotResponse{}
}
func (m *TSDBSnapshotResponse) String() string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return proto.CompactTextString(m)
}
func (*TSDBSnapshotResponse) ProtoMessage() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
}
func (*TSDBSnapshotResponse) Descriptor() ([]byte, []int) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return fileDescriptor_rpc_e0d54cfadc26b2e1, []int{1}
}
func (m *TSDBSnapshotResponse) XXX_Unmarshal(b []byte) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return m.Unmarshal(b)
}
func (m *TSDBSnapshotResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	if deterministic {
		return xxx_messageInfo_TSDBSnapshotResponse.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalTo(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (dst *TSDBSnapshotResponse) XXX_Merge(src proto.Message) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	xxx_messageInfo_TSDBSnapshotResponse.Merge(dst, src)
}
func (m *TSDBSnapshotResponse) XXX_Size() int {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return m.Size()
}
func (m *TSDBSnapshotResponse) XXX_DiscardUnknown() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	xxx_messageInfo_TSDBSnapshotResponse.DiscardUnknown(m)
}

var xxx_messageInfo_TSDBSnapshotResponse proto.InternalMessageInfo

type TSDBCleanTombstonesRequest struct {
	XXX_NoUnkeyedLiteral	struct{}	`json:"-"`
	XXX_unrecognized	[]byte		`json:"-"`
	XXX_sizecache		int32		`json:"-"`
}

func (m *TSDBCleanTombstonesRequest) Reset() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	*m = TSDBCleanTombstonesRequest{}
}
func (m *TSDBCleanTombstonesRequest) String() string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return proto.CompactTextString(m)
}
func (*TSDBCleanTombstonesRequest) ProtoMessage() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
}
func (*TSDBCleanTombstonesRequest) Descriptor() ([]byte, []int) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return fileDescriptor_rpc_e0d54cfadc26b2e1, []int{2}
}
func (m *TSDBCleanTombstonesRequest) XXX_Unmarshal(b []byte) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return m.Unmarshal(b)
}
func (m *TSDBCleanTombstonesRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	if deterministic {
		return xxx_messageInfo_TSDBCleanTombstonesRequest.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalTo(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (dst *TSDBCleanTombstonesRequest) XXX_Merge(src proto.Message) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	xxx_messageInfo_TSDBCleanTombstonesRequest.Merge(dst, src)
}
func (m *TSDBCleanTombstonesRequest) XXX_Size() int {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return m.Size()
}
func (m *TSDBCleanTombstonesRequest) XXX_DiscardUnknown() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	xxx_messageInfo_TSDBCleanTombstonesRequest.DiscardUnknown(m)
}

var xxx_messageInfo_TSDBCleanTombstonesRequest proto.InternalMessageInfo

type TSDBCleanTombstonesResponse struct {
	XXX_NoUnkeyedLiteral	struct{}	`json:"-"`
	XXX_unrecognized	[]byte		`json:"-"`
	XXX_sizecache		int32		`json:"-"`
}

func (m *TSDBCleanTombstonesResponse) Reset() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	*m = TSDBCleanTombstonesResponse{}
}
func (m *TSDBCleanTombstonesResponse) String() string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return proto.CompactTextString(m)
}
func (*TSDBCleanTombstonesResponse) ProtoMessage() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
}
func (*TSDBCleanTombstonesResponse) Descriptor() ([]byte, []int) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return fileDescriptor_rpc_e0d54cfadc26b2e1, []int{3}
}
func (m *TSDBCleanTombstonesResponse) XXX_Unmarshal(b []byte) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return m.Unmarshal(b)
}
func (m *TSDBCleanTombstonesResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	if deterministic {
		return xxx_messageInfo_TSDBCleanTombstonesResponse.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalTo(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (dst *TSDBCleanTombstonesResponse) XXX_Merge(src proto.Message) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	xxx_messageInfo_TSDBCleanTombstonesResponse.Merge(dst, src)
}
func (m *TSDBCleanTombstonesResponse) XXX_Size() int {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return m.Size()
}
func (m *TSDBCleanTombstonesResponse) XXX_DiscardUnknown() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	xxx_messageInfo_TSDBCleanTombstonesResponse.DiscardUnknown(m)
}

var xxx_messageInfo_TSDBCleanTombstonesResponse proto.InternalMessageInfo

type SeriesDeleteRequest struct {
	MinTime			*time.Time	`protobuf:"bytes,1,opt,name=min_time,json=minTime,proto3,stdtime" json:"min_time,omitempty"`
	MaxTime			*time.Time	`protobuf:"bytes,2,opt,name=max_time,json=maxTime,proto3,stdtime" json:"max_time,omitempty"`
	Matchers		[]LabelMatcher	`protobuf:"bytes,3,rep,name=matchers,proto3" json:"matchers"`
	XXX_NoUnkeyedLiteral	struct{}	`json:"-"`
	XXX_unrecognized	[]byte		`json:"-"`
	XXX_sizecache		int32		`json:"-"`
}

func (m *SeriesDeleteRequest) Reset() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	*m = SeriesDeleteRequest{}
}
func (m *SeriesDeleteRequest) String() string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return proto.CompactTextString(m)
}
func (*SeriesDeleteRequest) ProtoMessage() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
}
func (*SeriesDeleteRequest) Descriptor() ([]byte, []int) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return fileDescriptor_rpc_e0d54cfadc26b2e1, []int{4}
}
func (m *SeriesDeleteRequest) XXX_Unmarshal(b []byte) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return m.Unmarshal(b)
}
func (m *SeriesDeleteRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	if deterministic {
		return xxx_messageInfo_SeriesDeleteRequest.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalTo(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (dst *SeriesDeleteRequest) XXX_Merge(src proto.Message) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	xxx_messageInfo_SeriesDeleteRequest.Merge(dst, src)
}
func (m *SeriesDeleteRequest) XXX_Size() int {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return m.Size()
}
func (m *SeriesDeleteRequest) XXX_DiscardUnknown() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	xxx_messageInfo_SeriesDeleteRequest.DiscardUnknown(m)
}

var xxx_messageInfo_SeriesDeleteRequest proto.InternalMessageInfo

type SeriesDeleteResponse struct {
	XXX_NoUnkeyedLiteral	struct{}	`json:"-"`
	XXX_unrecognized	[]byte		`json:"-"`
	XXX_sizecache		int32		`json:"-"`
}

func (m *SeriesDeleteResponse) Reset() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	*m = SeriesDeleteResponse{}
}
func (m *SeriesDeleteResponse) String() string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return proto.CompactTextString(m)
}
func (*SeriesDeleteResponse) ProtoMessage() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
}
func (*SeriesDeleteResponse) Descriptor() ([]byte, []int) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return fileDescriptor_rpc_e0d54cfadc26b2e1, []int{5}
}
func (m *SeriesDeleteResponse) XXX_Unmarshal(b []byte) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return m.Unmarshal(b)
}
func (m *SeriesDeleteResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	if deterministic {
		return xxx_messageInfo_SeriesDeleteResponse.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalTo(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (dst *SeriesDeleteResponse) XXX_Merge(src proto.Message) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	xxx_messageInfo_SeriesDeleteResponse.Merge(dst, src)
}
func (m *SeriesDeleteResponse) XXX_Size() int {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return m.Size()
}
func (m *SeriesDeleteResponse) XXX_DiscardUnknown() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	xxx_messageInfo_SeriesDeleteResponse.DiscardUnknown(m)
}

var xxx_messageInfo_SeriesDeleteResponse proto.InternalMessageInfo

func init() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	proto.RegisterType((*TSDBSnapshotRequest)(nil), "prometheus.TSDBSnapshotRequest")
	proto.RegisterType((*TSDBSnapshotResponse)(nil), "prometheus.TSDBSnapshotResponse")
	proto.RegisterType((*TSDBCleanTombstonesRequest)(nil), "prometheus.TSDBCleanTombstonesRequest")
	proto.RegisterType((*TSDBCleanTombstonesResponse)(nil), "prometheus.TSDBCleanTombstonesResponse")
	proto.RegisterType((*SeriesDeleteRequest)(nil), "prometheus.SeriesDeleteRequest")
	proto.RegisterType((*SeriesDeleteResponse)(nil), "prometheus.SeriesDeleteResponse")
}

var _ context.Context
var _ grpc.ClientConn

const _ = grpc.SupportPackageIsVersion4

type AdminClient interface {
	TSDBSnapshot(ctx context.Context, in *TSDBSnapshotRequest, opts ...grpc.CallOption) (*TSDBSnapshotResponse, error)
	TSDBCleanTombstones(ctx context.Context, in *TSDBCleanTombstonesRequest, opts ...grpc.CallOption) (*TSDBCleanTombstonesResponse, error)
	DeleteSeries(ctx context.Context, in *SeriesDeleteRequest, opts ...grpc.CallOption) (*SeriesDeleteResponse, error)
}
type adminClient struct{ cc *grpc.ClientConn }

func NewAdminClient(cc *grpc.ClientConn) AdminClient {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return &adminClient{cc}
}
func (c *adminClient) TSDBSnapshot(ctx context.Context, in *TSDBSnapshotRequest, opts ...grpc.CallOption) (*TSDBSnapshotResponse, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	out := new(TSDBSnapshotResponse)
	err := c.cc.Invoke(ctx, "/prometheus.Admin/TSDBSnapshot", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}
func (c *adminClient) TSDBCleanTombstones(ctx context.Context, in *TSDBCleanTombstonesRequest, opts ...grpc.CallOption) (*TSDBCleanTombstonesResponse, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	out := new(TSDBCleanTombstonesResponse)
	err := c.cc.Invoke(ctx, "/prometheus.Admin/TSDBCleanTombstones", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}
func (c *adminClient) DeleteSeries(ctx context.Context, in *SeriesDeleteRequest, opts ...grpc.CallOption) (*SeriesDeleteResponse, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	out := new(SeriesDeleteResponse)
	err := c.cc.Invoke(ctx, "/prometheus.Admin/DeleteSeries", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

type AdminServer interface {
	TSDBSnapshot(context.Context, *TSDBSnapshotRequest) (*TSDBSnapshotResponse, error)
	TSDBCleanTombstones(context.Context, *TSDBCleanTombstonesRequest) (*TSDBCleanTombstonesResponse, error)
	DeleteSeries(context.Context, *SeriesDeleteRequest) (*SeriesDeleteResponse, error)
}

func RegisterAdminServer(s *grpc.Server, srv AdminServer) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	s.RegisterService(&_Admin_serviceDesc, srv)
}
func _Admin_TSDBSnapshot_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	in := new(TSDBSnapshotRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AdminServer).TSDBSnapshot(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/prometheus.Admin/TSDBSnapshot"}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AdminServer).TSDBSnapshot(ctx, req.(*TSDBSnapshotRequest))
	}
	return interceptor(ctx, in, info, handler)
}
func _Admin_TSDBCleanTombstones_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	in := new(TSDBCleanTombstonesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AdminServer).TSDBCleanTombstones(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/prometheus.Admin/TSDBCleanTombstones"}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AdminServer).TSDBCleanTombstones(ctx, req.(*TSDBCleanTombstonesRequest))
	}
	return interceptor(ctx, in, info, handler)
}
func _Admin_DeleteSeries_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	in := new(SeriesDeleteRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AdminServer).DeleteSeries(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/prometheus.Admin/DeleteSeries"}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AdminServer).DeleteSeries(ctx, req.(*SeriesDeleteRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _Admin_serviceDesc = grpc.ServiceDesc{ServiceName: "prometheus.Admin", HandlerType: (*AdminServer)(nil), Methods: []grpc.MethodDesc{{MethodName: "TSDBSnapshot", Handler: _Admin_TSDBSnapshot_Handler}, {MethodName: "TSDBCleanTombstones", Handler: _Admin_TSDBCleanTombstones_Handler}, {MethodName: "DeleteSeries", Handler: _Admin_DeleteSeries_Handler}}, Streams: []grpc.StreamDesc{}, Metadata: "rpc.proto"}

func (m *TSDBSnapshotRequest) Marshal() (dAtA []byte, err error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}
func (m *TSDBSnapshotRequest) MarshalTo(dAtA []byte) (int, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	var i int
	_ = i
	var l int
	_ = l
	if m.SkipHead {
		dAtA[i] = 0x8
		i++
		if m.SkipHead {
			dAtA[i] = 1
		} else {
			dAtA[i] = 0
		}
		i++
	}
	if m.XXX_unrecognized != nil {
		i += copy(dAtA[i:], m.XXX_unrecognized)
	}
	return i, nil
}
func (m *TSDBSnapshotResponse) Marshal() (dAtA []byte, err error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}
func (m *TSDBSnapshotResponse) MarshalTo(dAtA []byte) (int, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	var i int
	_ = i
	var l int
	_ = l
	if len(m.Name) > 0 {
		dAtA[i] = 0xa
		i++
		i = encodeVarintRpc(dAtA, i, uint64(len(m.Name)))
		i += copy(dAtA[i:], m.Name)
	}
	if m.XXX_unrecognized != nil {
		i += copy(dAtA[i:], m.XXX_unrecognized)
	}
	return i, nil
}
func (m *TSDBCleanTombstonesRequest) Marshal() (dAtA []byte, err error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}
func (m *TSDBCleanTombstonesRequest) MarshalTo(dAtA []byte) (int, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	var i int
	_ = i
	var l int
	_ = l
	if m.XXX_unrecognized != nil {
		i += copy(dAtA[i:], m.XXX_unrecognized)
	}
	return i, nil
}
func (m *TSDBCleanTombstonesResponse) Marshal() (dAtA []byte, err error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}
func (m *TSDBCleanTombstonesResponse) MarshalTo(dAtA []byte) (int, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	var i int
	_ = i
	var l int
	_ = l
	if m.XXX_unrecognized != nil {
		i += copy(dAtA[i:], m.XXX_unrecognized)
	}
	return i, nil
}
func (m *SeriesDeleteRequest) Marshal() (dAtA []byte, err error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}
func (m *SeriesDeleteRequest) MarshalTo(dAtA []byte) (int, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	var i int
	_ = i
	var l int
	_ = l
	if m.MinTime != nil {
		dAtA[i] = 0xa
		i++
		i = encodeVarintRpc(dAtA, i, uint64(github_com_gogo_protobuf_types.SizeOfStdTime(*m.MinTime)))
		n1, err := github_com_gogo_protobuf_types.StdTimeMarshalTo(*m.MinTime, dAtA[i:])
		if err != nil {
			return 0, err
		}
		i += n1
	}
	if m.MaxTime != nil {
		dAtA[i] = 0x12
		i++
		i = encodeVarintRpc(dAtA, i, uint64(github_com_gogo_protobuf_types.SizeOfStdTime(*m.MaxTime)))
		n2, err := github_com_gogo_protobuf_types.StdTimeMarshalTo(*m.MaxTime, dAtA[i:])
		if err != nil {
			return 0, err
		}
		i += n2
	}
	if len(m.Matchers) > 0 {
		for _, msg := range m.Matchers {
			dAtA[i] = 0x1a
			i++
			i = encodeVarintRpc(dAtA, i, uint64(msg.Size()))
			n, err := msg.MarshalTo(dAtA[i:])
			if err != nil {
				return 0, err
			}
			i += n
		}
	}
	if m.XXX_unrecognized != nil {
		i += copy(dAtA[i:], m.XXX_unrecognized)
	}
	return i, nil
}
func (m *SeriesDeleteResponse) Marshal() (dAtA []byte, err error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}
func (m *SeriesDeleteResponse) MarshalTo(dAtA []byte) (int, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	var i int
	_ = i
	var l int
	_ = l
	if m.XXX_unrecognized != nil {
		i += copy(dAtA[i:], m.XXX_unrecognized)
	}
	return i, nil
}
func encodeVarintRpc(dAtA []byte, offset int, v uint64) int {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return offset + 1
}
func (m *TSDBSnapshotRequest) Size() (n int) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.SkipHead {
		n += 2
	}
	if m.XXX_unrecognized != nil {
		n += len(m.XXX_unrecognized)
	}
	return n
}
func (m *TSDBSnapshotResponse) Size() (n int) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Name)
	if l > 0 {
		n += 1 + l + sovRpc(uint64(l))
	}
	if m.XXX_unrecognized != nil {
		n += len(m.XXX_unrecognized)
	}
	return n
}
func (m *TSDBCleanTombstonesRequest) Size() (n int) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.XXX_unrecognized != nil {
		n += len(m.XXX_unrecognized)
	}
	return n
}
func (m *TSDBCleanTombstonesResponse) Size() (n int) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.XXX_unrecognized != nil {
		n += len(m.XXX_unrecognized)
	}
	return n
}
func (m *SeriesDeleteRequest) Size() (n int) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.MinTime != nil {
		l = github_com_gogo_protobuf_types.SizeOfStdTime(*m.MinTime)
		n += 1 + l + sovRpc(uint64(l))
	}
	if m.MaxTime != nil {
		l = github_com_gogo_protobuf_types.SizeOfStdTime(*m.MaxTime)
		n += 1 + l + sovRpc(uint64(l))
	}
	if len(m.Matchers) > 0 {
		for _, e := range m.Matchers {
			l = e.Size()
			n += 1 + l + sovRpc(uint64(l))
		}
	}
	if m.XXX_unrecognized != nil {
		n += len(m.XXX_unrecognized)
	}
	return n
}
func (m *SeriesDeleteResponse) Size() (n int) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.XXX_unrecognized != nil {
		n += len(m.XXX_unrecognized)
	}
	return n
}
func sovRpc(x uint64) (n int) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	for {
		n++
		x >>= 7
		if x == 0 {
			break
		}
	}
	return n
}
func sozRpc(x uint64) (n int) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return sovRpc(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *TSDBSnapshotRequest) Unmarshal(dAtA []byte) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowRpc
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
			return fmt.Errorf("proto: TSDBSnapshotRequest: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: TSDBSnapshotRequest: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field SkipHead", wireType)
			}
			var v int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowRpc
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				v |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			m.SkipHead = bool(v != 0)
		default:
			iNdEx = preIndex
			skippy, err := skipRpc(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthRpc
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			m.XXX_unrecognized = append(m.XXX_unrecognized, dAtA[iNdEx:iNdEx+skippy]...)
			iNdEx += skippy
		}
	}
	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *TSDBSnapshotResponse) Unmarshal(dAtA []byte) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowRpc
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
			return fmt.Errorf("proto: TSDBSnapshotResponse: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: TSDBSnapshotResponse: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Name", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowRpc
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
				return ErrInvalidLengthRpc
			}
			postIndex := iNdEx + intStringLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Name = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipRpc(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthRpc
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			m.XXX_unrecognized = append(m.XXX_unrecognized, dAtA[iNdEx:iNdEx+skippy]...)
			iNdEx += skippy
		}
	}
	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *TSDBCleanTombstonesRequest) Unmarshal(dAtA []byte) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowRpc
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
			return fmt.Errorf("proto: TSDBCleanTombstonesRequest: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: TSDBCleanTombstonesRequest: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		default:
			iNdEx = preIndex
			skippy, err := skipRpc(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthRpc
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			m.XXX_unrecognized = append(m.XXX_unrecognized, dAtA[iNdEx:iNdEx+skippy]...)
			iNdEx += skippy
		}
	}
	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *TSDBCleanTombstonesResponse) Unmarshal(dAtA []byte) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowRpc
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
			return fmt.Errorf("proto: TSDBCleanTombstonesResponse: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: TSDBCleanTombstonesResponse: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		default:
			iNdEx = preIndex
			skippy, err := skipRpc(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthRpc
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			m.XXX_unrecognized = append(m.XXX_unrecognized, dAtA[iNdEx:iNdEx+skippy]...)
			iNdEx += skippy
		}
	}
	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *SeriesDeleteRequest) Unmarshal(dAtA []byte) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowRpc
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
			return fmt.Errorf("proto: SeriesDeleteRequest: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: SeriesDeleteRequest: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field MinTime", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowRpc
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
				return ErrInvalidLengthRpc
			}
			postIndex := iNdEx + msglen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.MinTime == nil {
				m.MinTime = new(time.Time)
			}
			if err := github_com_gogo_protobuf_types.StdTimeUnmarshal(m.MinTime, dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field MaxTime", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowRpc
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
				return ErrInvalidLengthRpc
			}
			postIndex := iNdEx + msglen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.MaxTime == nil {
				m.MaxTime = new(time.Time)
			}
			if err := github_com_gogo_protobuf_types.StdTimeUnmarshal(m.MaxTime, dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Matchers", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowRpc
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
				return ErrInvalidLengthRpc
			}
			postIndex := iNdEx + msglen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Matchers = append(m.Matchers, LabelMatcher{})
			if err := m.Matchers[len(m.Matchers)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipRpc(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthRpc
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			m.XXX_unrecognized = append(m.XXX_unrecognized, dAtA[iNdEx:iNdEx+skippy]...)
			iNdEx += skippy
		}
	}
	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *SeriesDeleteResponse) Unmarshal(dAtA []byte) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowRpc
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
			return fmt.Errorf("proto: SeriesDeleteResponse: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: SeriesDeleteResponse: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		default:
			iNdEx = preIndex
			skippy, err := skipRpc(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthRpc
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			m.XXX_unrecognized = append(m.XXX_unrecognized, dAtA[iNdEx:iNdEx+skippy]...)
			iNdEx += skippy
		}
	}
	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func skipRpc(dAtA []byte) (n int, err error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowRpc
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
					return 0, ErrIntOverflowRpc
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
					return 0, ErrIntOverflowRpc
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
				return 0, ErrInvalidLengthRpc
			}
			return iNdEx, nil
		case 3:
			for {
				var innerWire uint64
				var start int = iNdEx
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return 0, ErrIntOverflowRpc
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
				next, err := skipRpc(dAtA[start:])
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
	ErrInvalidLengthRpc	= fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowRpc	= fmt.Errorf("proto: integer overflow")
)

func init() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	proto.RegisterFile("rpc.proto", fileDescriptor_rpc_e0d54cfadc26b2e1)
}

var fileDescriptor_rpc_e0d54cfadc26b2e1 = []byte{0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x8c, 0x52, 0x3d, 0x8f, 0xd3, 0x40, 0x10, 0xbd, 0xbd, 0x84, 0x23, 0xd9, 0x5c, 0xe5, 0x8b, 0x20, 0xf8, 0x42, 0x1c, 0x5c, 0x70, 0xa7, 0x2b, 0x6c, 0xc9, 0x74, 0x47, 0x45, 0xb8, 0x82, 0x02, 0x1a, 0x27, 0x15, 0x4d, 0xb4, 0x8e, 0x87, 0xc4, 0x22, 0xfb, 0x81, 0x77, 0x83, 0x0e, 0x4a, 0x3a, 0x2a, 0x90, 0xf8, 0x53, 0x91, 0x68, 0x90, 0xe8, 0xf9, 0x88, 0xf8, 0x21, 0x68, 0x77, 0xed, 0xbb, 0xc4, 0x32, 0xe2, 0xba, 0xd9, 0xd9, 0xf7, 0xe6, 0xcd, 0xbc, 0x19, 0xdc, 0xce, 0xc5, 0x2c, 0x10, 0x39, 0x57, 0xdc, 0xc1, 0x22, 0xe7, 0x14, 0xd4, 0x02, 0x56, 0xd2, 0xed, 0xa8, 0x77, 0x02, 0xa4, 0xfd, 0x70, 0xbd, 0x39, 0xe7, 0xf3, 0x25, 0x84, 0xe6, 0x95, 0xac, 0x5e, 0x85, 0x2a, 0xa3, 0x20, 0x15, 0xa1, 0xa2, 0x00, 0xf4, 0x0b, 0x00, 0x11, 0x59, 0x48, 0x18, 0xe3, 0x8a, 0xa8, 0x8c, 0xb3, 0x92, 0xde, 0x9d, 0xf3, 0x39, 0x37, 0x61, 0xa8, 0x23, 0x9b, 0xf5, 0x23, 0x7c, 0x34, 0x19, 0x5f, 0x8c, 0xc6, 0x8c, 0x08, 0xb9, 0xe0, 0x2a, 0x86, 0x37, 0x2b, 0x90, 0xca, 0x39, 0xc6, 0x6d, 0xf9, 0x3a, 0x13, 0xd3, 0x05, 0x90, 0xb4, 0x87, 0x86, 0xe8, 0xb4, 0x15, 0xb7, 0x74, 0xe2, 0x19, 0x90, 0xd4, 0x3f, 0xc3, 0xdd, 0x5d, 0x8e, 0x14, 0x9c, 0x49, 0x70, 0x1c, 0xdc, 0x64, 0x84, 0x82, 0xc1, 0xb7, 0x63, 0x13, 0xfb, 0x7d, 0xec, 0x6a, 0xec, 0xd3, 0x25, 0x10, 0x36, 0xe1, 0x34, 0x91, 0x8a, 0x33, 0x90, 0x85, 0x8c, 0x7f, 0x1f, 0x1f, 0xd7, 0xfe, 0xda, 0x82, 0xfe, 0x57, 0x84, 0x8f, 0xc6, 0x90, 0x67, 0x20, 0x2f, 0x60, 0x09, 0x0a, 0xca, 0xee, 0x1e, 0xe3, 0x16, 0xcd, 0xd8, 0x54, 0xcf, 0x6f, 0xc4, 0x3a, 0x91, 0x1b, 0xd8, 0xd9, 0x83, 0xd2, 0x9c, 0x60, 0x52, 0x9a, 0x33, 0x6a, 0x7e, 0xfe, 0xe9, 0xa1, 0xf8, 0x36, 0xcd, 0x98, 0xce, 0x19, 0x32, 0xb9, 0xb4, 0xe4, 0xfd, 0x1b, 0x93, 0xc9, 0xa5, 0x21, 0x9f, 0x6b, 0xb2, 0x9a, 0x2d, 0x20, 0x97, 0xbd, 0xc6, 0xb0, 0x71, 0xda, 0x89, 0x7a, 0xc1, 0xf5, 0xbe, 0x82, 0xe7, 0x24, 0x81, 0xe5, 0x0b, 0x0b, 0x18, 0x35, 0xd7, 0x3f, 0xbc, 0xbd, 0xf8, 0x0a, 0xef, 0xdf, 0xc1, 0xdd, 0xdd, 0x61, 0xec, 0x94, 0xd1, 0xc7, 0x06, 0xbe, 0xf5, 0x24, 0xa5, 0x19, 0x73, 0x72, 0x7c, 0xb8, 0x6d, 0xac, 0xe3, 0x6d, 0xd7, 0xae, 0x59, 0x93, 0x3b, 0xfc, 0x37, 0xa0, 0xb0, 0xd0, 0xfb, 0xf0, 0xfd, 0xcf, 0x97, 0xfd, 0x7b, 0xfe, 0xdd, 0xf0, 0x6d, 0x14, 0x12, 0xad, 0x12, 0x2a, 0x99, 0x26, 0xa1, 0x2c, 0x35, 0x3e, 0x21, 0x7b, 0x01, 0x95, 0x1d, 0x38, 0x0f, 0xab, 0xa5, 0xeb, 0x57, 0xe8, 0x9e, 0xfc, 0x17, 0x57, 0x74, 0x72, 0x62, 0x3a, 0x79, 0xe0, 0x7b, 0x95, 0x4e, 0x66, 0x1a, 0x3f, 0x55, 0xd7, 0xca, 0xef, 0xf1, 0xa1, 0x75, 0xc8, 0xba, 0xb5, 0xeb, 0x42, 0xcd, 0x39, 0xec, 0xba, 0x50, 0x67, 0xf1, 0x95, 0x76, 0xbf, 0xa2, 0x9d, 0x1a, 0xd8, 0x54, 0x1a, 0xce, 0x39, 0x3a, 0x1b, 0xf5, 0xd6, 0xbf, 0x07, 0x7b, 0xeb, 0xcd, 0x00, 0x7d, 0xdb, 0x0c, 0xd0, 0xaf, 0xcd, 0x00, 0xbd, 0x3c, 0xd0, 0xb5, 0x45, 0x92, 0x1c, 0x98, 0xe3, 0x78, 0xf4, 0x37, 0x00, 0x00, 0xff, 0xff, 0x8b, 0x21, 0x42, 0x5d, 0xaa, 0x03, 0x00, 0x00}
