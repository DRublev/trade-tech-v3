// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.32.0
// 	protoc        v4.25.2
// source: shares.proto

package grpcGW_shares

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type Quatation struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Units int32 `protobuf:"varint,1,opt,name=units,proto3" json:"units,omitempty"`
	Nano  int32 `protobuf:"varint,2,opt,name=nano,proto3" json:"nano,omitempty"`
}

func (x *Quatation) Reset() {
	*x = Quatation{}
	if protoimpl.UnsafeEnabled {
		mi := &file_shares_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Quatation) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Quatation) ProtoMessage() {}

func (x *Quatation) ProtoReflect() protoreflect.Message {
	mi := &file_shares_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Quatation.ProtoReflect.Descriptor instead.
func (*Quatation) Descriptor() ([]byte, []int) {
	return file_shares_proto_rawDescGZIP(), []int{0}
}

func (x *Quatation) GetUnits() int32 {
	if x != nil {
		return x.Units
	}
	return 0
}

func (x *Quatation) GetNano() int32 {
	if x != nil {
		return x.Nano
	}
	return 0
}

type Share struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Figi                string                 `protobuf:"bytes,1,opt,name=figi,proto3" json:"figi,omitempty"`
	Name                string                 `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	Exchange            string                 `protobuf:"bytes,3,opt,name=exchange,proto3" json:"exchange,omitempty"`
	Ticker              string                 `protobuf:"bytes,4,opt,name=ticker,proto3" json:"ticker,omitempty"`
	Lot                 int32                  `protobuf:"varint,5,opt,name=lot,proto3" json:"lot,omitempty"`
	IpoDate             *timestamppb.Timestamp `protobuf:"bytes,6,opt,name=ipoDate,proto3" json:"ipoDate,omitempty"`
	TradingStatus       int32                  `protobuf:"varint,7,opt,name=tradingStatus,proto3" json:"tradingStatus,omitempty"`
	MinPriceIncrement   *Quatation             `protobuf:"bytes,8,opt,name=minPriceIncrement,proto3" json:"minPriceIncrement,omitempty"`
	Uid                 string                 `protobuf:"bytes,9,opt,name=uid,proto3" json:"uid,omitempty"`
	First1MinCandleDate *timestamppb.Timestamp `protobuf:"bytes,10,opt,name=first1minCandleDate,proto3" json:"first1minCandleDate,omitempty"`
	First1DayCandleDate *timestamppb.Timestamp `protobuf:"bytes,11,opt,name=first1dayCandleDate,proto3" json:"first1dayCandleDate,omitempty"`
}

func (x *Share) Reset() {
	*x = Share{}
	if protoimpl.UnsafeEnabled {
		mi := &file_shares_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Share) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Share) ProtoMessage() {}

func (x *Share) ProtoReflect() protoreflect.Message {
	mi := &file_shares_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Share.ProtoReflect.Descriptor instead.
func (*Share) Descriptor() ([]byte, []int) {
	return file_shares_proto_rawDescGZIP(), []int{1}
}

func (x *Share) GetFigi() string {
	if x != nil {
		return x.Figi
	}
	return ""
}

func (x *Share) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *Share) GetExchange() string {
	if x != nil {
		return x.Exchange
	}
	return ""
}

func (x *Share) GetTicker() string {
	if x != nil {
		return x.Ticker
	}
	return ""
}

func (x *Share) GetLot() int32 {
	if x != nil {
		return x.Lot
	}
	return 0
}

func (x *Share) GetIpoDate() *timestamppb.Timestamp {
	if x != nil {
		return x.IpoDate
	}
	return nil
}

func (x *Share) GetTradingStatus() int32 {
	if x != nil {
		return x.TradingStatus
	}
	return 0
}

func (x *Share) GetMinPriceIncrement() *Quatation {
	if x != nil {
		return x.MinPriceIncrement
	}
	return nil
}

func (x *Share) GetUid() string {
	if x != nil {
		return x.Uid
	}
	return ""
}

func (x *Share) GetFirst1MinCandleDate() *timestamppb.Timestamp {
	if x != nil {
		return x.First1MinCandleDate
	}
	return nil
}

func (x *Share) GetFirst1DayCandleDate() *timestamppb.Timestamp {
	if x != nil {
		return x.First1DayCandleDate
	}
	return nil
}

type GetInstrumentsRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	InstrumentStatus int32 `protobuf:"varint,2,opt,name=instrumentStatus,proto3" json:"instrumentStatus,omitempty"`
}

func (x *GetInstrumentsRequest) Reset() {
	*x = GetInstrumentsRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_shares_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetInstrumentsRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetInstrumentsRequest) ProtoMessage() {}

func (x *GetInstrumentsRequest) ProtoReflect() protoreflect.Message {
	mi := &file_shares_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetInstrumentsRequest.ProtoReflect.Descriptor instead.
func (*GetInstrumentsRequest) Descriptor() ([]byte, []int) {
	return file_shares_proto_rawDescGZIP(), []int{2}
}

func (x *GetInstrumentsRequest) GetInstrumentStatus() int32 {
	if x != nil {
		return x.InstrumentStatus
	}
	return 0
}

type GetSharesResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Instruments []*Share `protobuf:"bytes,1,rep,name=instruments,proto3" json:"instruments,omitempty"`
}

func (x *GetSharesResponse) Reset() {
	*x = GetSharesResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_shares_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetSharesResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetSharesResponse) ProtoMessage() {}

func (x *GetSharesResponse) ProtoReflect() protoreflect.Message {
	mi := &file_shares_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetSharesResponse.ProtoReflect.Descriptor instead.
func (*GetSharesResponse) Descriptor() ([]byte, []int) {
	return file_shares_proto_rawDescGZIP(), []int{3}
}

func (x *GetSharesResponse) GetInstruments() []*Share {
	if x != nil {
		return x.Instruments
	}
	return nil
}

var File_shares_proto protoreflect.FileDescriptor

var file_shares_proto_rawDesc = []byte{
	0x0a, 0x0c, 0x73, 0x68, 0x61, 0x72, 0x65, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x06,
	0x73, 0x68, 0x61, 0x72, 0x65, 0x73, 0x1a, 0x1f, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d,
	0x70, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x35, 0x0a, 0x09, 0x51, 0x75, 0x61, 0x74, 0x61,
	0x74, 0x69, 0x6f, 0x6e, 0x12, 0x14, 0x0a, 0x05, 0x75, 0x6e, 0x69, 0x74, 0x73, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x05, 0x52, 0x05, 0x75, 0x6e, 0x69, 0x74, 0x73, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61,
	0x6e, 0x6f, 0x18, 0x02, 0x20, 0x01, 0x28, 0x05, 0x52, 0x04, 0x6e, 0x61, 0x6e, 0x6f, 0x22, 0xc0,
	0x03, 0x0a, 0x05, 0x53, 0x68, 0x61, 0x72, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x66, 0x69, 0x67, 0x69,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x66, 0x69, 0x67, 0x69, 0x12, 0x12, 0x0a, 0x04,
	0x6e, 0x61, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65,
	0x12, 0x1a, 0x0a, 0x08, 0x65, 0x78, 0x63, 0x68, 0x61, 0x6e, 0x67, 0x65, 0x18, 0x03, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x08, 0x65, 0x78, 0x63, 0x68, 0x61, 0x6e, 0x67, 0x65, 0x12, 0x16, 0x0a, 0x06,
	0x74, 0x69, 0x63, 0x6b, 0x65, 0x72, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x74, 0x69,
	0x63, 0x6b, 0x65, 0x72, 0x12, 0x10, 0x0a, 0x03, 0x6c, 0x6f, 0x74, 0x18, 0x05, 0x20, 0x01, 0x28,
	0x05, 0x52, 0x03, 0x6c, 0x6f, 0x74, 0x12, 0x34, 0x0a, 0x07, 0x69, 0x70, 0x6f, 0x44, 0x61, 0x74,
	0x65, 0x18, 0x06, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74,
	0x61, 0x6d, 0x70, 0x52, 0x07, 0x69, 0x70, 0x6f, 0x44, 0x61, 0x74, 0x65, 0x12, 0x24, 0x0a, 0x0d,
	0x74, 0x72, 0x61, 0x64, 0x69, 0x6e, 0x67, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x18, 0x07, 0x20,
	0x01, 0x28, 0x05, 0x52, 0x0d, 0x74, 0x72, 0x61, 0x64, 0x69, 0x6e, 0x67, 0x53, 0x74, 0x61, 0x74,
	0x75, 0x73, 0x12, 0x3f, 0x0a, 0x11, 0x6d, 0x69, 0x6e, 0x50, 0x72, 0x69, 0x63, 0x65, 0x49, 0x6e,
	0x63, 0x72, 0x65, 0x6d, 0x65, 0x6e, 0x74, 0x18, 0x08, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x11, 0x2e,
	0x73, 0x68, 0x61, 0x72, 0x65, 0x73, 0x2e, 0x51, 0x75, 0x61, 0x74, 0x61, 0x74, 0x69, 0x6f, 0x6e,
	0x52, 0x11, 0x6d, 0x69, 0x6e, 0x50, 0x72, 0x69, 0x63, 0x65, 0x49, 0x6e, 0x63, 0x72, 0x65, 0x6d,
	0x65, 0x6e, 0x74, 0x12, 0x10, 0x0a, 0x03, 0x75, 0x69, 0x64, 0x18, 0x09, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x03, 0x75, 0x69, 0x64, 0x12, 0x4c, 0x0a, 0x13, 0x66, 0x69, 0x72, 0x73, 0x74, 0x31, 0x6d,
	0x69, 0x6e, 0x43, 0x61, 0x6e, 0x64, 0x6c, 0x65, 0x44, 0x61, 0x74, 0x65, 0x18, 0x0a, 0x20, 0x01,
	0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x13,
	0x66, 0x69, 0x72, 0x73, 0x74, 0x31, 0x6d, 0x69, 0x6e, 0x43, 0x61, 0x6e, 0x64, 0x6c, 0x65, 0x44,
	0x61, 0x74, 0x65, 0x12, 0x4c, 0x0a, 0x13, 0x66, 0x69, 0x72, 0x73, 0x74, 0x31, 0x64, 0x61, 0x79,
	0x43, 0x61, 0x6e, 0x64, 0x6c, 0x65, 0x44, 0x61, 0x74, 0x65, 0x18, 0x0b, 0x20, 0x01, 0x28, 0x0b,
	0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62,
	0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x13, 0x66, 0x69,
	0x72, 0x73, 0x74, 0x31, 0x64, 0x61, 0x79, 0x43, 0x61, 0x6e, 0x64, 0x6c, 0x65, 0x44, 0x61, 0x74,
	0x65, 0x22, 0x43, 0x0a, 0x15, 0x47, 0x65, 0x74, 0x49, 0x6e, 0x73, 0x74, 0x72, 0x75, 0x6d, 0x65,
	0x6e, 0x74, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x2a, 0x0a, 0x10, 0x69, 0x6e,
	0x73, 0x74, 0x72, 0x75, 0x6d, 0x65, 0x6e, 0x74, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x05, 0x52, 0x10, 0x69, 0x6e, 0x73, 0x74, 0x72, 0x75, 0x6d, 0x65, 0x6e, 0x74,
	0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x22, 0x44, 0x0a, 0x11, 0x47, 0x65, 0x74, 0x53, 0x68, 0x61,
	0x72, 0x65, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x2f, 0x0a, 0x0b, 0x69,
	0x6e, 0x73, 0x74, 0x72, 0x75, 0x6d, 0x65, 0x6e, 0x74, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b,
	0x32, 0x0d, 0x2e, 0x73, 0x68, 0x61, 0x72, 0x65, 0x73, 0x2e, 0x53, 0x68, 0x61, 0x72, 0x65, 0x52,
	0x0b, 0x69, 0x6e, 0x73, 0x74, 0x72, 0x75, 0x6d, 0x65, 0x6e, 0x74, 0x73, 0x32, 0x4f, 0x0a, 0x06,
	0x53, 0x68, 0x61, 0x72, 0x65, 0x73, 0x12, 0x45, 0x0a, 0x09, 0x47, 0x65, 0x74, 0x53, 0x68, 0x61,
	0x72, 0x65, 0x73, 0x12, 0x1d, 0x2e, 0x73, 0x68, 0x61, 0x72, 0x65, 0x73, 0x2e, 0x47, 0x65, 0x74,
	0x49, 0x6e, 0x73, 0x74, 0x72, 0x75, 0x6d, 0x65, 0x6e, 0x74, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x1a, 0x19, 0x2e, 0x73, 0x68, 0x61, 0x72, 0x65, 0x73, 0x2e, 0x47, 0x65, 0x74, 0x53,
	0x68, 0x61, 0x72, 0x65, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x42, 0x0f, 0x5a,
	0x0d, 0x67, 0x72, 0x70, 0x63, 0x47, 0x57, 0x2e, 0x73, 0x68, 0x61, 0x72, 0x65, 0x73, 0x62, 0x06,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_shares_proto_rawDescOnce sync.Once
	file_shares_proto_rawDescData = file_shares_proto_rawDesc
)

func file_shares_proto_rawDescGZIP() []byte {
	file_shares_proto_rawDescOnce.Do(func() {
		file_shares_proto_rawDescData = protoimpl.X.CompressGZIP(file_shares_proto_rawDescData)
	})
	return file_shares_proto_rawDescData
}

var file_shares_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_shares_proto_goTypes = []interface{}{
	(*Quatation)(nil),             // 0: shares.Quatation
	(*Share)(nil),                 // 1: shares.Share
	(*GetInstrumentsRequest)(nil), // 2: shares.GetInstrumentsRequest
	(*GetSharesResponse)(nil),     // 3: shares.GetSharesResponse
	(*timestamppb.Timestamp)(nil), // 4: google.protobuf.Timestamp
}
var file_shares_proto_depIdxs = []int32{
	4, // 0: shares.Share.ipoDate:type_name -> google.protobuf.Timestamp
	0, // 1: shares.Share.minPriceIncrement:type_name -> shares.Quatation
	4, // 2: shares.Share.first1minCandleDate:type_name -> google.protobuf.Timestamp
	4, // 3: shares.Share.first1dayCandleDate:type_name -> google.protobuf.Timestamp
	1, // 4: shares.GetSharesResponse.instruments:type_name -> shares.Share
	2, // 5: shares.Shares.GetShares:input_type -> shares.GetInstrumentsRequest
	3, // 6: shares.Shares.GetShares:output_type -> shares.GetSharesResponse
	6, // [6:7] is the sub-list for method output_type
	5, // [5:6] is the sub-list for method input_type
	5, // [5:5] is the sub-list for extension type_name
	5, // [5:5] is the sub-list for extension extendee
	0, // [0:5] is the sub-list for field type_name
}

func init() { file_shares_proto_init() }
func file_shares_proto_init() {
	if File_shares_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_shares_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Quatation); i {
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
		file_shares_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Share); i {
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
		file_shares_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetInstrumentsRequest); i {
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
		file_shares_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetSharesResponse); i {
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
			RawDescriptor: file_shares_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_shares_proto_goTypes,
		DependencyIndexes: file_shares_proto_depIdxs,
		MessageInfos:      file_shares_proto_msgTypes,
	}.Build()
	File_shares_proto = out.File
	file_shares_proto_rawDesc = nil
	file_shares_proto_goTypes = nil
	file_shares_proto_depIdxs = nil
}
