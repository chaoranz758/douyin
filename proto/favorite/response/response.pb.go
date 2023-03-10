// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v3.20.1
// source: favorite/response/response.proto

package response

import (
	response "douyin/proto/video/response"
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

type DouyinFavoriteActionResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Code int32 `protobuf:"varint,2,opt,name=code,proto3" json:"code,omitempty"`
}

func (x *DouyinFavoriteActionResponse) Reset() {
	*x = DouyinFavoriteActionResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_favorite_response_response_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DouyinFavoriteActionResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DouyinFavoriteActionResponse) ProtoMessage() {}

func (x *DouyinFavoriteActionResponse) ProtoReflect() protoreflect.Message {
	mi := &file_favorite_response_response_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DouyinFavoriteActionResponse.ProtoReflect.Descriptor instead.
func (*DouyinFavoriteActionResponse) Descriptor() ([]byte, []int) {
	return file_favorite_response_response_proto_rawDescGZIP(), []int{0}
}

func (x *DouyinFavoriteActionResponse) GetCode() int32 {
	if x != nil {
		return x.Code
	}
	return 0
}

type DouyinFavoriteListResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	VideoList []*response.Video `protobuf:"bytes,1,rep,name=video_list,json=videoList,proto3" json:"video_list,omitempty"`
	Code      int32             `protobuf:"varint,2,opt,name=code,proto3" json:"code,omitempty"`
}

func (x *DouyinFavoriteListResponse) Reset() {
	*x = DouyinFavoriteListResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_favorite_response_response_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DouyinFavoriteListResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DouyinFavoriteListResponse) ProtoMessage() {}

func (x *DouyinFavoriteListResponse) ProtoReflect() protoreflect.Message {
	mi := &file_favorite_response_response_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DouyinFavoriteListResponse.ProtoReflect.Descriptor instead.
func (*DouyinFavoriteListResponse) Descriptor() ([]byte, []int) {
	return file_favorite_response_response_proto_rawDescGZIP(), []int{1}
}

func (x *DouyinFavoriteListResponse) GetVideoList() []*response.Video {
	if x != nil {
		return x.VideoList
	}
	return nil
}

func (x *DouyinFavoriteListResponse) GetCode() int32 {
	if x != nil {
		return x.Code
	}
	return 0
}

type DouyinFavoriteCountResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	FavoriteCount []int64 `protobuf:"varint,1,rep,packed,name=favorite_count,json=favoriteCount,proto3" json:"favorite_count,omitempty"`
	IsFavorite    []bool  `protobuf:"varint,3,rep,packed,name=is_favorite,json=isFavorite,proto3" json:"is_favorite,omitempty"`
	Code          int32   `protobuf:"varint,2,opt,name=code,proto3" json:"code,omitempty"`
}

func (x *DouyinFavoriteCountResponse) Reset() {
	*x = DouyinFavoriteCountResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_favorite_response_response_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DouyinFavoriteCountResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DouyinFavoriteCountResponse) ProtoMessage() {}

func (x *DouyinFavoriteCountResponse) ProtoReflect() protoreflect.Message {
	mi := &file_favorite_response_response_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DouyinFavoriteCountResponse.ProtoReflect.Descriptor instead.
func (*DouyinFavoriteCountResponse) Descriptor() ([]byte, []int) {
	return file_favorite_response_response_proto_rawDescGZIP(), []int{2}
}

func (x *DouyinFavoriteCountResponse) GetFavoriteCount() []int64 {
	if x != nil {
		return x.FavoriteCount
	}
	return nil
}

func (x *DouyinFavoriteCountResponse) GetIsFavorite() []bool {
	if x != nil {
		return x.IsFavorite
	}
	return nil
}

func (x *DouyinFavoriteCountResponse) GetCode() int32 {
	if x != nil {
		return x.Code
	}
	return 0
}

type DouyinFavoriteIdListResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Code    int32   `protobuf:"varint,2,opt,name=code,proto3" json:"code,omitempty"`
	VideoId []int64 `protobuf:"varint,1,rep,packed,name=video_id,json=videoId,proto3" json:"video_id,omitempty"`
}

func (x *DouyinFavoriteIdListResponse) Reset() {
	*x = DouyinFavoriteIdListResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_favorite_response_response_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DouyinFavoriteIdListResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DouyinFavoriteIdListResponse) ProtoMessage() {}

func (x *DouyinFavoriteIdListResponse) ProtoReflect() protoreflect.Message {
	mi := &file_favorite_response_response_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DouyinFavoriteIdListResponse.ProtoReflect.Descriptor instead.
func (*DouyinFavoriteIdListResponse) Descriptor() ([]byte, []int) {
	return file_favorite_response_response_proto_rawDescGZIP(), []int{3}
}

func (x *DouyinFavoriteIdListResponse) GetCode() int32 {
	if x != nil {
		return x.Code
	}
	return 0
}

func (x *DouyinFavoriteIdListResponse) GetVideoId() []int64 {
	if x != nil {
		return x.VideoId
	}
	return nil
}

type DouyinFavoriteListIdListResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Code                int32                  `protobuf:"varint,2,opt,name=code,proto3" json:"code,omitempty"`
	FavoriteVideoIdList []*FavoriteVideoIdList `protobuf:"bytes,1,rep,name=favorite_video_id_list,json=favoriteVideoIdList,proto3" json:"favorite_video_id_list,omitempty"`
}

func (x *DouyinFavoriteListIdListResponse) Reset() {
	*x = DouyinFavoriteListIdListResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_favorite_response_response_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DouyinFavoriteListIdListResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DouyinFavoriteListIdListResponse) ProtoMessage() {}

func (x *DouyinFavoriteListIdListResponse) ProtoReflect() protoreflect.Message {
	mi := &file_favorite_response_response_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DouyinFavoriteListIdListResponse.ProtoReflect.Descriptor instead.
func (*DouyinFavoriteListIdListResponse) Descriptor() ([]byte, []int) {
	return file_favorite_response_response_proto_rawDescGZIP(), []int{4}
}

func (x *DouyinFavoriteListIdListResponse) GetCode() int32 {
	if x != nil {
		return x.Code
	}
	return 0
}

func (x *DouyinFavoriteListIdListResponse) GetFavoriteVideoIdList() []*FavoriteVideoIdList {
	if x != nil {
		return x.FavoriteVideoIdList
	}
	return nil
}

type FavoriteVideoIdList struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	VideoId []int64 `protobuf:"varint,1,rep,packed,name=video_id,json=videoId,proto3" json:"video_id,omitempty"`
}

func (x *FavoriteVideoIdList) Reset() {
	*x = FavoriteVideoIdList{}
	if protoimpl.UnsafeEnabled {
		mi := &file_favorite_response_response_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *FavoriteVideoIdList) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*FavoriteVideoIdList) ProtoMessage() {}

func (x *FavoriteVideoIdList) ProtoReflect() protoreflect.Message {
	mi := &file_favorite_response_response_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use FavoriteVideoIdList.ProtoReflect.Descriptor instead.
func (*FavoriteVideoIdList) Descriptor() ([]byte, []int) {
	return file_favorite_response_response_proto_rawDescGZIP(), []int{5}
}

func (x *FavoriteVideoIdList) GetVideoId() []int64 {
	if x != nil {
		return x.VideoId
	}
	return nil
}

type DouyinGetUserFavoritedCountResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	FavoritedCount []int64 `protobuf:"varint,1,rep,packed,name=favorited_count,json=favoritedCount,proto3" json:"favorited_count,omitempty"`
	FavoriteCount  []int64 `protobuf:"varint,2,rep,packed,name=favorite_count,json=favoriteCount,proto3" json:"favorite_count,omitempty"`
	Code           int32   `protobuf:"varint,3,opt,name=code,proto3" json:"code,omitempty"`
}

func (x *DouyinGetUserFavoritedCountResponse) Reset() {
	*x = DouyinGetUserFavoritedCountResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_favorite_response_response_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DouyinGetUserFavoritedCountResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DouyinGetUserFavoritedCountResponse) ProtoMessage() {}

func (x *DouyinGetUserFavoritedCountResponse) ProtoReflect() protoreflect.Message {
	mi := &file_favorite_response_response_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DouyinGetUserFavoritedCountResponse.ProtoReflect.Descriptor instead.
func (*DouyinGetUserFavoritedCountResponse) Descriptor() ([]byte, []int) {
	return file_favorite_response_response_proto_rawDescGZIP(), []int{6}
}

func (x *DouyinGetUserFavoritedCountResponse) GetFavoritedCount() []int64 {
	if x != nil {
		return x.FavoritedCount
	}
	return nil
}

func (x *DouyinGetUserFavoritedCountResponse) GetFavoriteCount() []int64 {
	if x != nil {
		return x.FavoriteCount
	}
	return nil
}

func (x *DouyinGetUserFavoritedCountResponse) GetCode() int32 {
	if x != nil {
		return x.Code
	}
	return 0
}

var File_favorite_response_response_proto protoreflect.FileDescriptor

var file_favorite_response_response_proto_rawDesc = []byte{
	0x0a, 0x20, 0x66, 0x61, 0x76, 0x6f, 0x72, 0x69, 0x74, 0x65, 0x2f, 0x72, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x2f, 0x72, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x12, 0x08, 0x72, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x1a, 0x1d, 0x76, 0x69,
	0x64, 0x65, 0x6f, 0x2f, 0x72, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x2f, 0x72, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x35, 0x0a, 0x1f, 0x64,
	0x6f, 0x75, 0x79, 0x69, 0x6e, 0x5f, 0x66, 0x61, 0x76, 0x6f, 0x72, 0x69, 0x74, 0x65, 0x5f, 0x61,
	0x63, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x72, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x12,
	0x0a, 0x04, 0x63, 0x6f, 0x64, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x05, 0x52, 0x04, 0x63, 0x6f,
	0x64, 0x65, 0x22, 0x63, 0x0a, 0x1d, 0x64, 0x6f, 0x75, 0x79, 0x69, 0x6e, 0x5f, 0x66, 0x61, 0x76,
	0x6f, 0x72, 0x69, 0x74, 0x65, 0x5f, 0x6c, 0x69, 0x73, 0x74, 0x5f, 0x72, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x12, 0x2e, 0x0a, 0x0a, 0x76, 0x69, 0x64, 0x65, 0x6f, 0x5f, 0x6c, 0x69, 0x73,
	0x74, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x0f, 0x2e, 0x72, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x73, 0x65, 0x2e, 0x56, 0x69, 0x64, 0x65, 0x6f, 0x52, 0x09, 0x76, 0x69, 0x64, 0x65, 0x6f, 0x4c,
	0x69, 0x73, 0x74, 0x12, 0x12, 0x0a, 0x04, 0x63, 0x6f, 0x64, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x05, 0x52, 0x04, 0x63, 0x6f, 0x64, 0x65, 0x22, 0x7c, 0x0a, 0x1e, 0x64, 0x6f, 0x75, 0x79, 0x69,
	0x6e, 0x5f, 0x66, 0x61, 0x76, 0x6f, 0x72, 0x69, 0x74, 0x65, 0x5f, 0x63, 0x6f, 0x75, 0x6e, 0x74,
	0x5f, 0x72, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x25, 0x0a, 0x0e, 0x66, 0x61, 0x76,
	0x6f, 0x72, 0x69, 0x74, 0x65, 0x5f, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x18, 0x01, 0x20, 0x03, 0x28,
	0x03, 0x52, 0x0d, 0x66, 0x61, 0x76, 0x6f, 0x72, 0x69, 0x74, 0x65, 0x43, 0x6f, 0x75, 0x6e, 0x74,
	0x12, 0x1f, 0x0a, 0x0b, 0x69, 0x73, 0x5f, 0x66, 0x61, 0x76, 0x6f, 0x72, 0x69, 0x74, 0x65, 0x18,
	0x03, 0x20, 0x03, 0x28, 0x08, 0x52, 0x0a, 0x69, 0x73, 0x46, 0x61, 0x76, 0x6f, 0x72, 0x69, 0x74,
	0x65, 0x12, 0x12, 0x0a, 0x04, 0x63, 0x6f, 0x64, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x05, 0x52,
	0x04, 0x63, 0x6f, 0x64, 0x65, 0x22, 0x51, 0x0a, 0x20, 0x64, 0x6f, 0x75, 0x79, 0x69, 0x6e, 0x5f,
	0x66, 0x61, 0x76, 0x6f, 0x72, 0x69, 0x74, 0x65, 0x5f, 0x69, 0x64, 0x5f, 0x6c, 0x69, 0x73, 0x74,
	0x5f, 0x72, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x63, 0x6f, 0x64,
	0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x05, 0x52, 0x04, 0x63, 0x6f, 0x64, 0x65, 0x12, 0x19, 0x0a,
	0x08, 0x76, 0x69, 0x64, 0x65, 0x6f, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x03, 0x28, 0x03, 0x52,
	0x07, 0x76, 0x69, 0x64, 0x65, 0x6f, 0x49, 0x64, 0x22, 0x92, 0x01, 0x0a, 0x25, 0x64, 0x6f, 0x75,
	0x79, 0x69, 0x6e, 0x5f, 0x66, 0x61, 0x76, 0x6f, 0x72, 0x69, 0x74, 0x65, 0x5f, 0x6c, 0x69, 0x73,
	0x74, 0x5f, 0x69, 0x64, 0x5f, 0x6c, 0x69, 0x73, 0x74, 0x5f, 0x72, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x73, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x63, 0x6f, 0x64, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x05,
	0x52, 0x04, 0x63, 0x6f, 0x64, 0x65, 0x12, 0x55, 0x0a, 0x16, 0x66, 0x61, 0x76, 0x6f, 0x72, 0x69,
	0x74, 0x65, 0x5f, 0x76, 0x69, 0x64, 0x65, 0x6f, 0x5f, 0x69, 0x64, 0x5f, 0x6c, 0x69, 0x73, 0x74,
	0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x20, 0x2e, 0x72, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x2e, 0x66, 0x61, 0x76, 0x6f, 0x72, 0x69, 0x74, 0x65, 0x5f, 0x76, 0x69, 0x64, 0x65, 0x6f,
	0x5f, 0x69, 0x64, 0x5f, 0x6c, 0x69, 0x73, 0x74, 0x52, 0x13, 0x66, 0x61, 0x76, 0x6f, 0x72, 0x69,
	0x74, 0x65, 0x56, 0x69, 0x64, 0x65, 0x6f, 0x49, 0x64, 0x4c, 0x69, 0x73, 0x74, 0x22, 0x33, 0x0a,
	0x16, 0x66, 0x61, 0x76, 0x6f, 0x72, 0x69, 0x74, 0x65, 0x5f, 0x76, 0x69, 0x64, 0x65, 0x6f, 0x5f,
	0x69, 0x64, 0x5f, 0x6c, 0x69, 0x73, 0x74, 0x12, 0x19, 0x0a, 0x08, 0x76, 0x69, 0x64, 0x65, 0x6f,
	0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x03, 0x28, 0x03, 0x52, 0x07, 0x76, 0x69, 0x64, 0x65, 0x6f,
	0x49, 0x64, 0x22, 0x8e, 0x01, 0x0a, 0x28, 0x64, 0x6f, 0x75, 0x79, 0x69, 0x6e, 0x5f, 0x67, 0x65,
	0x74, 0x5f, 0x75, 0x73, 0x65, 0x72, 0x5f, 0x66, 0x61, 0x76, 0x6f, 0x72, 0x69, 0x74, 0x65, 0x64,
	0x5f, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x5f, 0x72, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12,
	0x27, 0x0a, 0x0f, 0x66, 0x61, 0x76, 0x6f, 0x72, 0x69, 0x74, 0x65, 0x64, 0x5f, 0x63, 0x6f, 0x75,
	0x6e, 0x74, 0x18, 0x01, 0x20, 0x03, 0x28, 0x03, 0x52, 0x0e, 0x66, 0x61, 0x76, 0x6f, 0x72, 0x69,
	0x74, 0x65, 0x64, 0x43, 0x6f, 0x75, 0x6e, 0x74, 0x12, 0x25, 0x0a, 0x0e, 0x66, 0x61, 0x76, 0x6f,
	0x72, 0x69, 0x74, 0x65, 0x5f, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x18, 0x02, 0x20, 0x03, 0x28, 0x03,
	0x52, 0x0d, 0x66, 0x61, 0x76, 0x6f, 0x72, 0x69, 0x74, 0x65, 0x43, 0x6f, 0x75, 0x6e, 0x74, 0x12,
	0x12, 0x0a, 0x04, 0x63, 0x6f, 0x64, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x05, 0x52, 0x04, 0x63,
	0x6f, 0x64, 0x65, 0x42, 0x20, 0x5a, 0x1e, 0x64, 0x6f, 0x75, 0x79, 0x69, 0x6e, 0x2f, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x2f, 0x66, 0x61, 0x76, 0x6f, 0x72, 0x69, 0x74, 0x65, 0x2f, 0x72, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_favorite_response_response_proto_rawDescOnce sync.Once
	file_favorite_response_response_proto_rawDescData = file_favorite_response_response_proto_rawDesc
)

func file_favorite_response_response_proto_rawDescGZIP() []byte {
	file_favorite_response_response_proto_rawDescOnce.Do(func() {
		file_favorite_response_response_proto_rawDescData = protoimpl.X.CompressGZIP(file_favorite_response_response_proto_rawDescData)
	})
	return file_favorite_response_response_proto_rawDescData
}

var file_favorite_response_response_proto_msgTypes = make([]protoimpl.MessageInfo, 7)
var file_favorite_response_response_proto_goTypes = []interface{}{
	(*DouyinFavoriteActionResponse)(nil),        // 0: response.douyin_favorite_action_response
	(*DouyinFavoriteListResponse)(nil),          // 1: response.douyin_favorite_list_response
	(*DouyinFavoriteCountResponse)(nil),         // 2: response.douyin_favorite_count_response
	(*DouyinFavoriteIdListResponse)(nil),        // 3: response.douyin_favorite_id_list_response
	(*DouyinFavoriteListIdListResponse)(nil),    // 4: response.douyin_favorite_list_id_list_response
	(*FavoriteVideoIdList)(nil),                 // 5: response.favorite_video_id_list
	(*DouyinGetUserFavoritedCountResponse)(nil), // 6: response.douyin_get_user_favorited_count_response
	(*response.Video)(nil),                      // 7: response.Video
}
var file_favorite_response_response_proto_depIdxs = []int32{
	7, // 0: response.douyin_favorite_list_response.video_list:type_name -> response.Video
	5, // 1: response.douyin_favorite_list_id_list_response.favorite_video_id_list:type_name -> response.favorite_video_id_list
	2, // [2:2] is the sub-list for method output_type
	2, // [2:2] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_favorite_response_response_proto_init() }
func file_favorite_response_response_proto_init() {
	if File_favorite_response_response_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_favorite_response_response_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DouyinFavoriteActionResponse); i {
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
		file_favorite_response_response_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DouyinFavoriteListResponse); i {
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
		file_favorite_response_response_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DouyinFavoriteCountResponse); i {
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
		file_favorite_response_response_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DouyinFavoriteIdListResponse); i {
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
		file_favorite_response_response_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DouyinFavoriteListIdListResponse); i {
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
		file_favorite_response_response_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*FavoriteVideoIdList); i {
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
		file_favorite_response_response_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DouyinGetUserFavoritedCountResponse); i {
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
			RawDescriptor: file_favorite_response_response_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   7,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_favorite_response_response_proto_goTypes,
		DependencyIndexes: file_favorite_response_response_proto_depIdxs,
		MessageInfos:      file_favorite_response_response_proto_msgTypes,
	}.Build()
	File_favorite_response_response_proto = out.File
	file_favorite_response_response_proto_rawDesc = nil
	file_favorite_response_response_proto_goTypes = nil
	file_favorite_response_response_proto_depIdxs = nil
}
