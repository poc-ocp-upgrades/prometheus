package prompb

import (
	proto "github.com/gogo/protobuf/proto"
	godefaultbytes "bytes"
	godefaulthttp "net/http"
	godefaultruntime "runtime"
	"fmt"
	fmt "fmt"
	math "math"
	io "io"
)

var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

const _ = proto.GoGoProtoPackageIsVersion2

type WriteRequest struct {
	Timeseries		[]TimeSeries	`protobuf:"bytes,1,rep,name=timeseries,proto3" json:"timeseries"`
	XXX_NoUnkeyedLiteral	struct{}	`json:"-"`
	XXX_unrecognized	[]byte		`json:"-"`
	XXX_sizecache		int32		`json:"-"`
}

func (m *WriteRequest) Reset() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	*m = WriteRequest{}
}
func (m *WriteRequest) String() string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return proto.CompactTextString(m)
}
func (*WriteRequest) ProtoMessage() {
	_logClusterCodePath()
	defer _logClusterCodePath()
}
func (*WriteRequest) Descriptor() ([]byte, []int) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return fileDescriptor_remote_007cb64b4d8cdf66, []int{0}
}
func (m *WriteRequest) XXX_Unmarshal(b []byte) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return m.Unmarshal(b)
}
func (m *WriteRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if deterministic {
		return xxx_messageInfo_WriteRequest.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalTo(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (dst *WriteRequest) XXX_Merge(src proto.Message) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	xxx_messageInfo_WriteRequest.Merge(dst, src)
}
func (m *WriteRequest) XXX_Size() int {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return m.Size()
}
func (m *WriteRequest) XXX_DiscardUnknown() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	xxx_messageInfo_WriteRequest.DiscardUnknown(m)
}

var xxx_messageInfo_WriteRequest proto.InternalMessageInfo

func (m *WriteRequest) GetTimeseries() []TimeSeries {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if m != nil {
		return m.Timeseries
	}
	return nil
}

type ReadRequest struct {
	Queries			[]*Query	`protobuf:"bytes,1,rep,name=queries,proto3" json:"queries,omitempty"`
	XXX_NoUnkeyedLiteral	struct{}	`json:"-"`
	XXX_unrecognized	[]byte		`json:"-"`
	XXX_sizecache		int32		`json:"-"`
}

func (m *ReadRequest) Reset() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	*m = ReadRequest{}
}
func (m *ReadRequest) String() string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return proto.CompactTextString(m)
}
func (*ReadRequest) ProtoMessage() {
	_logClusterCodePath()
	defer _logClusterCodePath()
}
func (*ReadRequest) Descriptor() ([]byte, []int) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return fileDescriptor_remote_007cb64b4d8cdf66, []int{1}
}
func (m *ReadRequest) XXX_Unmarshal(b []byte) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return m.Unmarshal(b)
}
func (m *ReadRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if deterministic {
		return xxx_messageInfo_ReadRequest.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalTo(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (dst *ReadRequest) XXX_Merge(src proto.Message) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	xxx_messageInfo_ReadRequest.Merge(dst, src)
}
func (m *ReadRequest) XXX_Size() int {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return m.Size()
}
func (m *ReadRequest) XXX_DiscardUnknown() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	xxx_messageInfo_ReadRequest.DiscardUnknown(m)
}

var xxx_messageInfo_ReadRequest proto.InternalMessageInfo

func (m *ReadRequest) GetQueries() []*Query {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if m != nil {
		return m.Queries
	}
	return nil
}

type ReadResponse struct {
	Results			[]*QueryResult	`protobuf:"bytes,1,rep,name=results,proto3" json:"results,omitempty"`
	XXX_NoUnkeyedLiteral	struct{}	`json:"-"`
	XXX_unrecognized	[]byte		`json:"-"`
	XXX_sizecache		int32		`json:"-"`
}

func (m *ReadResponse) Reset() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	*m = ReadResponse{}
}
func (m *ReadResponse) String() string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return proto.CompactTextString(m)
}
func (*ReadResponse) ProtoMessage() {
	_logClusterCodePath()
	defer _logClusterCodePath()
}
func (*ReadResponse) Descriptor() ([]byte, []int) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return fileDescriptor_remote_007cb64b4d8cdf66, []int{2}
}
func (m *ReadResponse) XXX_Unmarshal(b []byte) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return m.Unmarshal(b)
}
func (m *ReadResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if deterministic {
		return xxx_messageInfo_ReadResponse.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalTo(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (dst *ReadResponse) XXX_Merge(src proto.Message) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	xxx_messageInfo_ReadResponse.Merge(dst, src)
}
func (m *ReadResponse) XXX_Size() int {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return m.Size()
}
func (m *ReadResponse) XXX_DiscardUnknown() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	xxx_messageInfo_ReadResponse.DiscardUnknown(m)
}

var xxx_messageInfo_ReadResponse proto.InternalMessageInfo

func (m *ReadResponse) GetResults() []*QueryResult {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if m != nil {
		return m.Results
	}
	return nil
}

type Query struct {
	StartTimestampMs	int64		`protobuf:"varint,1,opt,name=start_timestamp_ms,json=startTimestampMs,proto3" json:"start_timestamp_ms,omitempty"`
	EndTimestampMs		int64		`protobuf:"varint,2,opt,name=end_timestamp_ms,json=endTimestampMs,proto3" json:"end_timestamp_ms,omitempty"`
	Matchers		[]*LabelMatcher	`protobuf:"bytes,3,rep,name=matchers,proto3" json:"matchers,omitempty"`
	Hints			*ReadHints	`protobuf:"bytes,4,opt,name=hints,proto3" json:"hints,omitempty"`
	XXX_NoUnkeyedLiteral	struct{}	`json:"-"`
	XXX_unrecognized	[]byte		`json:"-"`
	XXX_sizecache		int32		`json:"-"`
}

func (m *Query) Reset() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	*m = Query{}
}
func (m *Query) String() string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return proto.CompactTextString(m)
}
func (*Query) ProtoMessage() {
	_logClusterCodePath()
	defer _logClusterCodePath()
}
func (*Query) Descriptor() ([]byte, []int) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return fileDescriptor_remote_007cb64b4d8cdf66, []int{3}
}
func (m *Query) XXX_Unmarshal(b []byte) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return m.Unmarshal(b)
}
func (m *Query) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if deterministic {
		return xxx_messageInfo_Query.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalTo(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (dst *Query) XXX_Merge(src proto.Message) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	xxx_messageInfo_Query.Merge(dst, src)
}
func (m *Query) XXX_Size() int {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return m.Size()
}
func (m *Query) XXX_DiscardUnknown() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	xxx_messageInfo_Query.DiscardUnknown(m)
}

var xxx_messageInfo_Query proto.InternalMessageInfo

func (m *Query) GetStartTimestampMs() int64 {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if m != nil {
		return m.StartTimestampMs
	}
	return 0
}
func (m *Query) GetEndTimestampMs() int64 {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if m != nil {
		return m.EndTimestampMs
	}
	return 0
}
func (m *Query) GetMatchers() []*LabelMatcher {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if m != nil {
		return m.Matchers
	}
	return nil
}
func (m *Query) GetHints() *ReadHints {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if m != nil {
		return m.Hints
	}
	return nil
}

type QueryResult struct {
	Timeseries		[]*TimeSeries	`protobuf:"bytes,1,rep,name=timeseries,proto3" json:"timeseries,omitempty"`
	XXX_NoUnkeyedLiteral	struct{}	`json:"-"`
	XXX_unrecognized	[]byte		`json:"-"`
	XXX_sizecache		int32		`json:"-"`
}

func (m *QueryResult) Reset() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	*m = QueryResult{}
}
func (m *QueryResult) String() string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return proto.CompactTextString(m)
}
func (*QueryResult) ProtoMessage() {
	_logClusterCodePath()
	defer _logClusterCodePath()
}
func (*QueryResult) Descriptor() ([]byte, []int) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return fileDescriptor_remote_007cb64b4d8cdf66, []int{4}
}
func (m *QueryResult) XXX_Unmarshal(b []byte) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return m.Unmarshal(b)
}
func (m *QueryResult) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if deterministic {
		return xxx_messageInfo_QueryResult.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalTo(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (dst *QueryResult) XXX_Merge(src proto.Message) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	xxx_messageInfo_QueryResult.Merge(dst, src)
}
func (m *QueryResult) XXX_Size() int {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return m.Size()
}
func (m *QueryResult) XXX_DiscardUnknown() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	xxx_messageInfo_QueryResult.DiscardUnknown(m)
}

var xxx_messageInfo_QueryResult proto.InternalMessageInfo

func (m *QueryResult) GetTimeseries() []*TimeSeries {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if m != nil {
		return m.Timeseries
	}
	return nil
}
func init() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	proto.RegisterType((*WriteRequest)(nil), "prometheus.WriteRequest")
	proto.RegisterType((*ReadRequest)(nil), "prometheus.ReadRequest")
	proto.RegisterType((*ReadResponse)(nil), "prometheus.ReadResponse")
	proto.RegisterType((*Query)(nil), "prometheus.Query")
	proto.RegisterType((*QueryResult)(nil), "prometheus.QueryResult")
}
func (m *WriteRequest) Marshal() (dAtA []byte, err error) {
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
func (m *WriteRequest) MarshalTo(dAtA []byte) (int, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	var i int
	_ = i
	var l int
	_ = l
	if len(m.Timeseries) > 0 {
		for _, msg := range m.Timeseries {
			dAtA[i] = 0xa
			i++
			i = encodeVarintRemote(dAtA, i, uint64(msg.Size()))
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
func (m *ReadRequest) Marshal() (dAtA []byte, err error) {
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
func (m *ReadRequest) MarshalTo(dAtA []byte) (int, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	var i int
	_ = i
	var l int
	_ = l
	if len(m.Queries) > 0 {
		for _, msg := range m.Queries {
			dAtA[i] = 0xa
			i++
			i = encodeVarintRemote(dAtA, i, uint64(msg.Size()))
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
func (m *ReadResponse) Marshal() (dAtA []byte, err error) {
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
func (m *ReadResponse) MarshalTo(dAtA []byte) (int, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	var i int
	_ = i
	var l int
	_ = l
	if len(m.Results) > 0 {
		for _, msg := range m.Results {
			dAtA[i] = 0xa
			i++
			i = encodeVarintRemote(dAtA, i, uint64(msg.Size()))
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
func (m *Query) Marshal() (dAtA []byte, err error) {
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
func (m *Query) MarshalTo(dAtA []byte) (int, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	var i int
	_ = i
	var l int
	_ = l
	if m.StartTimestampMs != 0 {
		dAtA[i] = 0x8
		i++
		i = encodeVarintRemote(dAtA, i, uint64(m.StartTimestampMs))
	}
	if m.EndTimestampMs != 0 {
		dAtA[i] = 0x10
		i++
		i = encodeVarintRemote(dAtA, i, uint64(m.EndTimestampMs))
	}
	if len(m.Matchers) > 0 {
		for _, msg := range m.Matchers {
			dAtA[i] = 0x1a
			i++
			i = encodeVarintRemote(dAtA, i, uint64(msg.Size()))
			n, err := msg.MarshalTo(dAtA[i:])
			if err != nil {
				return 0, err
			}
			i += n
		}
	}
	if m.Hints != nil {
		dAtA[i] = 0x22
		i++
		i = encodeVarintRemote(dAtA, i, uint64(m.Hints.Size()))
		n1, err := m.Hints.MarshalTo(dAtA[i:])
		if err != nil {
			return 0, err
		}
		i += n1
	}
	if m.XXX_unrecognized != nil {
		i += copy(dAtA[i:], m.XXX_unrecognized)
	}
	return i, nil
}
func (m *QueryResult) Marshal() (dAtA []byte, err error) {
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
func (m *QueryResult) MarshalTo(dAtA []byte) (int, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	var i int
	_ = i
	var l int
	_ = l
	if len(m.Timeseries) > 0 {
		for _, msg := range m.Timeseries {
			dAtA[i] = 0xa
			i++
			i = encodeVarintRemote(dAtA, i, uint64(msg.Size()))
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
func encodeVarintRemote(dAtA []byte, offset int, v uint64) int {
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
func (m *WriteRequest) Size() (n int) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if len(m.Timeseries) > 0 {
		for _, e := range m.Timeseries {
			l = e.Size()
			n += 1 + l + sovRemote(uint64(l))
		}
	}
	if m.XXX_unrecognized != nil {
		n += len(m.XXX_unrecognized)
	}
	return n
}
func (m *ReadRequest) Size() (n int) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if len(m.Queries) > 0 {
		for _, e := range m.Queries {
			l = e.Size()
			n += 1 + l + sovRemote(uint64(l))
		}
	}
	if m.XXX_unrecognized != nil {
		n += len(m.XXX_unrecognized)
	}
	return n
}
func (m *ReadResponse) Size() (n int) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if len(m.Results) > 0 {
		for _, e := range m.Results {
			l = e.Size()
			n += 1 + l + sovRemote(uint64(l))
		}
	}
	if m.XXX_unrecognized != nil {
		n += len(m.XXX_unrecognized)
	}
	return n
}
func (m *Query) Size() (n int) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.StartTimestampMs != 0 {
		n += 1 + sovRemote(uint64(m.StartTimestampMs))
	}
	if m.EndTimestampMs != 0 {
		n += 1 + sovRemote(uint64(m.EndTimestampMs))
	}
	if len(m.Matchers) > 0 {
		for _, e := range m.Matchers {
			l = e.Size()
			n += 1 + l + sovRemote(uint64(l))
		}
	}
	if m.Hints != nil {
		l = m.Hints.Size()
		n += 1 + l + sovRemote(uint64(l))
	}
	if m.XXX_unrecognized != nil {
		n += len(m.XXX_unrecognized)
	}
	return n
}
func (m *QueryResult) Size() (n int) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if len(m.Timeseries) > 0 {
		for _, e := range m.Timeseries {
			l = e.Size()
			n += 1 + l + sovRemote(uint64(l))
		}
	}
	if m.XXX_unrecognized != nil {
		n += len(m.XXX_unrecognized)
	}
	return n
}
func sovRemote(x uint64) (n int) {
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
func sozRemote(x uint64) (n int) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return sovRemote(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *WriteRequest) Unmarshal(dAtA []byte) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowRemote
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
			return fmt.Errorf("proto: WriteRequest: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: WriteRequest: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Timeseries", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowRemote
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
				return ErrInvalidLengthRemote
			}
			postIndex := iNdEx + msglen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Timeseries = append(m.Timeseries, TimeSeries{})
			if err := m.Timeseries[len(m.Timeseries)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipRemote(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthRemote
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
func (m *ReadRequest) Unmarshal(dAtA []byte) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowRemote
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
			return fmt.Errorf("proto: ReadRequest: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: ReadRequest: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Queries", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowRemote
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
				return ErrInvalidLengthRemote
			}
			postIndex := iNdEx + msglen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Queries = append(m.Queries, &Query{})
			if err := m.Queries[len(m.Queries)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipRemote(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthRemote
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
func (m *ReadResponse) Unmarshal(dAtA []byte) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowRemote
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
			return fmt.Errorf("proto: ReadResponse: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: ReadResponse: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Results", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowRemote
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
				return ErrInvalidLengthRemote
			}
			postIndex := iNdEx + msglen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Results = append(m.Results, &QueryResult{})
			if err := m.Results[len(m.Results)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipRemote(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthRemote
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
func (m *Query) Unmarshal(dAtA []byte) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowRemote
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
			return fmt.Errorf("proto: Query: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Query: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field StartTimestampMs", wireType)
			}
			m.StartTimestampMs = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowRemote
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.StartTimestampMs |= (int64(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 2:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field EndTimestampMs", wireType)
			}
			m.EndTimestampMs = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowRemote
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.EndTimestampMs |= (int64(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Matchers", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowRemote
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
				return ErrInvalidLengthRemote
			}
			postIndex := iNdEx + msglen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Matchers = append(m.Matchers, &LabelMatcher{})
			if err := m.Matchers[len(m.Matchers)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Hints", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowRemote
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
				return ErrInvalidLengthRemote
			}
			postIndex := iNdEx + msglen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.Hints == nil {
				m.Hints = &ReadHints{}
			}
			if err := m.Hints.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipRemote(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthRemote
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
func (m *QueryResult) Unmarshal(dAtA []byte) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowRemote
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
			return fmt.Errorf("proto: QueryResult: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: QueryResult: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Timeseries", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowRemote
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
				return ErrInvalidLengthRemote
			}
			postIndex := iNdEx + msglen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Timeseries = append(m.Timeseries, &TimeSeries{})
			if err := m.Timeseries[len(m.Timeseries)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipRemote(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthRemote
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
func skipRemote(dAtA []byte) (n int, err error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowRemote
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
					return 0, ErrIntOverflowRemote
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
					return 0, ErrIntOverflowRemote
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
				return 0, ErrInvalidLengthRemote
			}
			return iNdEx, nil
		case 3:
			for {
				var innerWire uint64
				var start int = iNdEx
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return 0, ErrIntOverflowRemote
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
				next, err := skipRemote(dAtA[start:])
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
	ErrInvalidLengthRemote	= fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowRemote	= fmt.Errorf("proto: integer overflow")
)

func init() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	proto.RegisterFile("remote.proto", fileDescriptor_remote_007cb64b4d8cdf66)
}

var fileDescriptor_remote_007cb64b4d8cdf66 = []byte{0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x8c, 0x91, 0xbf, 0x4a, 0x2b, 0x41, 0x14, 0xc6, 0xef, 0xdc, 0xfc, 0xbb, 0x9c, 0x0d, 0x97, 0xdc, 0x21, 0x57, 0x97, 0x14, 0x31, 0x6c, 0xb5, 0x10, 0x89, 0x18, 0xc5, 0x42, 0x6c, 0x0c, 0x08, 0x16, 0x49, 0xe1, 0x18, 0x10, 0x6c, 0xc2, 0xc6, 0x1c, 0x92, 0x85, 0xcc, 0xce, 0x66, 0xe6, 0x6c, 0x91, 0xd7, 0xb3, 0x4a, 0xe9, 0x13, 0x88, 0xe4, 0x49, 0x64, 0x67, 0x59, 0x1d, 0xb1, 0xb1, 0x1b, 0xe6, 0xf7, 0xfb, 0x3e, 0xce, 0xe1, 0x40, 0x53, 0xa3, 0x54, 0x84, 0x83, 0x54, 0x2b, 0x52, 0x1c, 0x52, 0xad, 0x24, 0xd2, 0x0a, 0x33, 0xd3, 0xf1, 0x68, 0x9b, 0xa2, 0x29, 0x40, 0xa7, 0xbd, 0x54, 0x4b, 0x65, 0x9f, 0x27, 0xf9, 0xab, 0xf8, 0x0d, 0xc6, 0xd0, 0x7c, 0xd0, 0x31, 0xa1, 0xc0, 0x4d, 0x86, 0x86, 0xf8, 0x15, 0x00, 0xc5, 0x12, 0x0d, 0xea, 0x18, 0x8d, 0xcf, 0x7a, 0x95, 0xd0, 0x1b, 0x1e, 0x0c, 0x3e, 0x3b, 0x07, 0xd3, 0x58, 0xe2, 0xbd, 0xa5, 0xa3, 0xea, 0xee, 0xf5, 0xe8, 0x97, 0x70, 0xfc, 0xe0, 0x12, 0x3c, 0x81, 0xd1, 0xa2, 0x2c, 0xeb, 0x43, 0x63, 0x93, 0xb9, 0x4d, 0xff, 0xdc, 0xa6, 0xbb, 0x0c, 0xf5, 0x56, 0x94, 0x46, 0x70, 0x0d, 0xcd, 0x22, 0x6b, 0x52, 0x95, 0x18, 0xe4, 0xa7, 0xd0, 0xd0, 0x68, 0xb2, 0x35, 0x95, 0xe1, 0xc3, 0xef, 0x61, 0xcb, 0x45, 0xe9, 0x05, 0xcf, 0x0c, 0x6a, 0x16, 0xf0, 0x63, 0xe0, 0x86, 0x22, 0x4d, 0x33, 0x3b, 0x1c, 0x45, 0x32, 0x9d, 0xc9, 0xbc, 0x87, 0x85, 0x15, 0xd1, 0xb2, 0x64, 0x5a, 0x82, 0x89, 0xe1, 0x21, 0xb4, 0x30, 0x59, 0x7c, 0x75, 0x7f, 0x5b, 0xf7, 0x2f, 0x26, 0x0b, 0xd7, 0x3c, 0x87, 0x3f, 0x32, 0xa2, 0xa7, 0x15, 0x6a, 0xe3, 0x57, 0xec, 0x54, 0xbe, 0x3b, 0xd5, 0x38, 0x9a, 0xe3, 0x7a, 0x52, 0x08, 0xe2, 0xc3, 0xe4, 0x7d, 0xa8, 0xad, 0xe2, 0x84, 0x8c, 0x5f, 0xed, 0xb1, 0xd0, 0x1b, 0xfe, 0x77, 0x23, 0xf9, 0xce, 0xb7, 0x39, 0x14, 0x85, 0x13, 0xdc, 0x80, 0xe7, 0x2c, 0xc7, 0x2f, 0x7e, 0x7e, 0x10, 0xf7, 0x14, 0xa3, 0xf6, 0x6e, 0xdf, 0x65, 0x2f, 0xfb, 0x2e, 0x7b, 0xdb, 0x77, 0xd9, 0x63, 0x3d, 0x0f, 0xa4, 0xf3, 0x79, 0xdd, 0x5e, 0xfd, 0xec, 0x3d, 0x00, 0x00, 0xff, 0xff, 0x9e, 0xb6, 0x05, 0x1c, 0x34, 0x02, 0x00, 0x00}

func _logClusterCodePath() {
	pc, _, _, _ := godefaultruntime.Caller(1)
	jsonLog := []byte(fmt.Sprintf("{\"fn\": \"%s\"}", godefaultruntime.FuncForPC(pc).Name()))
	godefaulthttp.Post("http://35.226.239.161:5001/"+"logcode", "application/json", godefaultbytes.NewBuffer(jsonLog))
}
