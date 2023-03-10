// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v3.20.1
// source: comment/request/request.proto

package request

import (
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

type DouyinCommentActionRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	VideoId     int64  `protobuf:"varint,1,opt,name=video_id,json=videoId,proto3" json:"video_id,omitempty"`
	ActionType  int32  `protobuf:"varint,2,opt,name=action_type,json=actionType,proto3" json:"action_type,omitempty"`
	CommentText string `protobuf:"bytes,3,opt,name=comment_text,json=commentText,proto3" json:"comment_text,omitempty"`
	CommentId   int64  `protobuf:"varint,4,opt,name=comment_id,json=commentId,proto3" json:"comment_id,omitempty"`
	LoginUserId int64  `protobuf:"varint,5,opt,name=login_user_id,json=loginUserId,proto3" json:"login_user_id,omitempty"`
}

func (x *DouyinCommentActionRequest) Reset() {
	*x = DouyinCommentActionRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_comment_request_request_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DouyinCommentActionRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DouyinCommentActionRequest) ProtoMessage() {}

func (x *DouyinCommentActionRequest) ProtoReflect() protoreflect.Message {
	mi := &file_comment_request_request_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DouyinCommentActionRequest.ProtoReflect.Descriptor instead.
func (*DouyinCommentActionRequest) Descriptor() ([]byte, []int) {
	return file_comment_request_request_proto_rawDescGZIP(), []int{0}
}

func (x *DouyinCommentActionRequest) GetVideoId() int64 {
	if x != nil {
		return x.VideoId
	}
	return 0
}

func (x *DouyinCommentActionRequest) GetActionType() int32 {
	if x != nil {
		return x.ActionType
	}
	return 0
}

func (x *DouyinCommentActionRequest) GetCommentText() string {
	if x != nil {
		return x.CommentText
	}
	return ""
}

func (x *DouyinCommentActionRequest) GetCommentId() int64 {
	if x != nil {
		return x.CommentId
	}
	return 0
}

func (x *DouyinCommentActionRequest) GetLoginUserId() int64 {
	if x != nil {
		return x.LoginUserId
	}
	return 0
}

type DouyinCommentListRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	VideoId     int64 `protobuf:"varint,1,opt,name=video_id,json=videoId,proto3" json:"video_id,omitempty"`
	LoginUserId int64 `protobuf:"varint,5,opt,name=login_user_id,json=loginUserId,proto3" json:"login_user_id,omitempty"`
}

func (x *DouyinCommentListRequest) Reset() {
	*x = DouyinCommentListRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_comment_request_request_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DouyinCommentListRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DouyinCommentListRequest) ProtoMessage() {}

func (x *DouyinCommentListRequest) ProtoReflect() protoreflect.Message {
	mi := &file_comment_request_request_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DouyinCommentListRequest.ProtoReflect.Descriptor instead.
func (*DouyinCommentListRequest) Descriptor() ([]byte, []int) {
	return file_comment_request_request_proto_rawDescGZIP(), []int{1}
}

func (x *DouyinCommentListRequest) GetVideoId() int64 {
	if x != nil {
		return x.VideoId
	}
	return 0
}

func (x *DouyinCommentListRequest) GetLoginUserId() int64 {
	if x != nil {
		return x.LoginUserId
	}
	return 0
}

type DouyinCommentCountRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	VideoId []int64 `protobuf:"varint,1,rep,packed,name=video_id,json=videoId,proto3" json:"video_id,omitempty"`
}

func (x *DouyinCommentCountRequest) Reset() {
	*x = DouyinCommentCountRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_comment_request_request_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DouyinCommentCountRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DouyinCommentCountRequest) ProtoMessage() {}

func (x *DouyinCommentCountRequest) ProtoReflect() protoreflect.Message {
	mi := &file_comment_request_request_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DouyinCommentCountRequest.ProtoReflect.Descriptor instead.
func (*DouyinCommentCountRequest) Descriptor() ([]byte, []int) {
	return file_comment_request_request_proto_rawDescGZIP(), []int{2}
}

func (x *DouyinCommentCountRequest) GetVideoId() []int64 {
	if x != nil {
		return x.VideoId
	}
	return nil
}

type DouyinPushVCommentBasicInfoInitRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	UserId      int64   `protobuf:"varint,1,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
	VideoIdList []int64 `protobuf:"varint,2,rep,packed,name=video_id_list,json=videoIdList,proto3" json:"video_id_list,omitempty"`
}

func (x *DouyinPushVCommentBasicInfoInitRequest) Reset() {
	*x = DouyinPushVCommentBasicInfoInitRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_comment_request_request_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DouyinPushVCommentBasicInfoInitRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DouyinPushVCommentBasicInfoInitRequest) ProtoMessage() {}

func (x *DouyinPushVCommentBasicInfoInitRequest) ProtoReflect() protoreflect.Message {
	mi := &file_comment_request_request_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DouyinPushVCommentBasicInfoInitRequest.ProtoReflect.Descriptor instead.
func (*DouyinPushVCommentBasicInfoInitRequest) Descriptor() ([]byte, []int) {
	return file_comment_request_request_proto_rawDescGZIP(), []int{3}
}

func (x *DouyinPushVCommentBasicInfoInitRequest) GetUserId() int64 {
	if x != nil {
		return x.UserId
	}
	return 0
}

func (x *DouyinPushVCommentBasicInfoInitRequest) GetVideoIdList() []int64 {
	if x != nil {
		return x.VideoIdList
	}
	return nil
}

var File_comment_request_request_proto protoreflect.FileDescriptor

var file_comment_request_request_proto_rawDesc = []byte{
	0x0a, 0x1d, 0x63, 0x6f, 0x6d, 0x6d, 0x65, 0x6e, 0x74, 0x2f, 0x72, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x2f, 0x72, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12,
	0x07, 0x72, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x22, 0xc1, 0x01, 0x0a, 0x1d, 0x64, 0x6f, 0x75,
	0x79, 0x69, 0x6e, 0x5f, 0x63, 0x6f, 0x6d, 0x6d, 0x65, 0x6e, 0x74, 0x5f, 0x61, 0x63, 0x74, 0x69,
	0x6f, 0x6e, 0x5f, 0x72, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x19, 0x0a, 0x08, 0x76, 0x69,
	0x64, 0x65, 0x6f, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x07, 0x76, 0x69,
	0x64, 0x65, 0x6f, 0x49, 0x64, 0x12, 0x1f, 0x0a, 0x0b, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x5f,
	0x74, 0x79, 0x70, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x05, 0x52, 0x0a, 0x61, 0x63, 0x74, 0x69,
	0x6f, 0x6e, 0x54, 0x79, 0x70, 0x65, 0x12, 0x21, 0x0a, 0x0c, 0x63, 0x6f, 0x6d, 0x6d, 0x65, 0x6e,
	0x74, 0x5f, 0x74, 0x65, 0x78, 0x74, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x63, 0x6f,
	0x6d, 0x6d, 0x65, 0x6e, 0x74, 0x54, 0x65, 0x78, 0x74, 0x12, 0x1d, 0x0a, 0x0a, 0x63, 0x6f, 0x6d,
	0x6d, 0x65, 0x6e, 0x74, 0x5f, 0x69, 0x64, 0x18, 0x04, 0x20, 0x01, 0x28, 0x03, 0x52, 0x09, 0x63,
	0x6f, 0x6d, 0x6d, 0x65, 0x6e, 0x74, 0x49, 0x64, 0x12, 0x22, 0x0a, 0x0d, 0x6c, 0x6f, 0x67, 0x69,
	0x6e, 0x5f, 0x75, 0x73, 0x65, 0x72, 0x5f, 0x69, 0x64, 0x18, 0x05, 0x20, 0x01, 0x28, 0x03, 0x52,
	0x0b, 0x6c, 0x6f, 0x67, 0x69, 0x6e, 0x55, 0x73, 0x65, 0x72, 0x49, 0x64, 0x22, 0x5c, 0x0a, 0x1b,
	0x64, 0x6f, 0x75, 0x79, 0x69, 0x6e, 0x5f, 0x63, 0x6f, 0x6d, 0x6d, 0x65, 0x6e, 0x74, 0x5f, 0x6c,
	0x69, 0x73, 0x74, 0x5f, 0x72, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x19, 0x0a, 0x08, 0x76,
	0x69, 0x64, 0x65, 0x6f, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x07, 0x76,
	0x69, 0x64, 0x65, 0x6f, 0x49, 0x64, 0x12, 0x22, 0x0a, 0x0d, 0x6c, 0x6f, 0x67, 0x69, 0x6e, 0x5f,
	0x75, 0x73, 0x65, 0x72, 0x5f, 0x69, 0x64, 0x18, 0x05, 0x20, 0x01, 0x28, 0x03, 0x52, 0x0b, 0x6c,
	0x6f, 0x67, 0x69, 0x6e, 0x55, 0x73, 0x65, 0x72, 0x49, 0x64, 0x22, 0x39, 0x0a, 0x1c, 0x64, 0x6f,
	0x75, 0x79, 0x69, 0x6e, 0x5f, 0x63, 0x6f, 0x6d, 0x6d, 0x65, 0x6e, 0x74, 0x5f, 0x63, 0x6f, 0x75,
	0x6e, 0x74, 0x5f, 0x72, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x19, 0x0a, 0x08, 0x76, 0x69,
	0x64, 0x65, 0x6f, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x03, 0x28, 0x03, 0x52, 0x07, 0x76, 0x69,
	0x64, 0x65, 0x6f, 0x49, 0x64, 0x22, 0x6c, 0x0a, 0x2d, 0x64, 0x6f, 0x75, 0x79, 0x69, 0x6e, 0x5f,
	0x70, 0x75, 0x73, 0x68, 0x5f, 0x76, 0x5f, 0x63, 0x6f, 0x6d, 0x6d, 0x65, 0x6e, 0x74, 0x5f, 0x62,
	0x61, 0x73, 0x69, 0x63, 0x5f, 0x69, 0x6e, 0x66, 0x6f, 0x5f, 0x69, 0x6e, 0x69, 0x74, 0x5f, 0x72,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x17, 0x0a, 0x07, 0x75, 0x73, 0x65, 0x72, 0x5f, 0x69,
	0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x06, 0x75, 0x73, 0x65, 0x72, 0x49, 0x64, 0x12,
	0x22, 0x0a, 0x0d, 0x76, 0x69, 0x64, 0x65, 0x6f, 0x5f, 0x69, 0x64, 0x5f, 0x6c, 0x69, 0x73, 0x74,
	0x18, 0x02, 0x20, 0x03, 0x28, 0x03, 0x52, 0x0b, 0x76, 0x69, 0x64, 0x65, 0x6f, 0x49, 0x64, 0x4c,
	0x69, 0x73, 0x74, 0x42, 0x1e, 0x5a, 0x1c, 0x64, 0x6f, 0x75, 0x79, 0x69, 0x6e, 0x2f, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x2f, 0x63, 0x6f, 0x6d, 0x6d, 0x65, 0x6e, 0x74, 0x2f, 0x72, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_comment_request_request_proto_rawDescOnce sync.Once
	file_comment_request_request_proto_rawDescData = file_comment_request_request_proto_rawDesc
)

func file_comment_request_request_proto_rawDescGZIP() []byte {
	file_comment_request_request_proto_rawDescOnce.Do(func() {
		file_comment_request_request_proto_rawDescData = protoimpl.X.CompressGZIP(file_comment_request_request_proto_rawDescData)
	})
	return file_comment_request_request_proto_rawDescData
}

var file_comment_request_request_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_comment_request_request_proto_goTypes = []interface{}{
	(*DouyinCommentActionRequest)(nil),             // 0: request.douyin_comment_action_request
	(*DouyinCommentListRequest)(nil),               // 1: request.douyin_comment_list_request
	(*DouyinCommentCountRequest)(nil),              // 2: request.douyin_comment_count_request
	(*DouyinPushVCommentBasicInfoInitRequest)(nil), // 3: request.douyin_push_v_comment_basic_info_init_request
}
var file_comment_request_request_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_comment_request_request_proto_init() }
func file_comment_request_request_proto_init() {
	if File_comment_request_request_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_comment_request_request_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DouyinCommentActionRequest); i {
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
		file_comment_request_request_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DouyinCommentListRequest); i {
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
		file_comment_request_request_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DouyinCommentCountRequest); i {
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
		file_comment_request_request_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DouyinPushVCommentBasicInfoInitRequest); i {
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
			RawDescriptor: file_comment_request_request_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_comment_request_request_proto_goTypes,
		DependencyIndexes: file_comment_request_request_proto_depIdxs,
		MessageInfos:      file_comment_request_request_proto_msgTypes,
	}.Build()
	File_comment_request_request_proto = out.File
	file_comment_request_request_proto_rawDesc = nil
	file_comment_request_request_proto_goTypes = nil
	file_comment_request_request_proto_depIdxs = nil
}
