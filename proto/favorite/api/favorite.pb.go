// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v3.20.1
// source: favorite/api/favorite.proto

package api

import (
	request "douyin/proto/favorite/request"
	response "douyin/proto/favorite/response"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

var File_favorite_api_favorite_proto protoreflect.FileDescriptor

var file_favorite_api_favorite_proto_rawDesc = []byte{
	0x0a, 0x1b, 0x66, 0x61, 0x76, 0x6f, 0x72, 0x69, 0x74, 0x65, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x66,
	0x61, 0x76, 0x6f, 0x72, 0x69, 0x74, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x03, 0x61,
	0x70, 0x69, 0x1a, 0x1e, 0x66, 0x61, 0x76, 0x6f, 0x72, 0x69, 0x74, 0x65, 0x2f, 0x72, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x2f, 0x72, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x1a, 0x20, 0x66, 0x61, 0x76, 0x6f, 0x72, 0x69, 0x74, 0x65, 0x2f, 0x72, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x2f, 0x72, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x32, 0xb3, 0x05, 0x0a, 0x08, 0x46, 0x61, 0x76, 0x6f, 0x72, 0x69, 0x74,
	0x65, 0x12, 0x63, 0x0a, 0x0d, 0x46, 0x61, 0x76, 0x6f, 0x72, 0x69, 0x74, 0x65, 0x56, 0x69, 0x64,
	0x65, 0x6f, 0x12, 0x27, 0x2e, 0x72, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x2e, 0x64, 0x6f, 0x75,
	0x79, 0x69, 0x6e, 0x5f, 0x66, 0x61, 0x76, 0x6f, 0x72, 0x69, 0x74, 0x65, 0x5f, 0x61, 0x63, 0x74,
	0x69, 0x6f, 0x6e, 0x5f, 0x72, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x29, 0x2e, 0x72, 0x65,
	0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x2e, 0x64, 0x6f, 0x75, 0x79, 0x69, 0x6e, 0x5f, 0x66, 0x61,
	0x76, 0x6f, 0x72, 0x69, 0x74, 0x65, 0x5f, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x72, 0x65,
	0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x66, 0x0a, 0x14, 0x47, 0x65, 0x74, 0x46, 0x61, 0x76,
	0x6f, 0x72, 0x69, 0x74, 0x65, 0x56, 0x69, 0x64, 0x65, 0x6f, 0x4c, 0x69, 0x73, 0x74, 0x12, 0x25,
	0x2e, 0x72, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x2e, 0x64, 0x6f, 0x75, 0x79, 0x69, 0x6e, 0x5f,
	0x66, 0x61, 0x76, 0x6f, 0x72, 0x69, 0x74, 0x65, 0x5f, 0x6c, 0x69, 0x73, 0x74, 0x5f, 0x72, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x27, 0x2e, 0x72, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65,
	0x2e, 0x64, 0x6f, 0x75, 0x79, 0x69, 0x6e, 0x5f, 0x66, 0x61, 0x76, 0x6f, 0x72, 0x69, 0x74, 0x65,
	0x5f, 0x6c, 0x69, 0x73, 0x74, 0x5f, 0x72, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x72,
	0x0a, 0x1a, 0x47, 0x65, 0x74, 0x55, 0x73, 0x65, 0x72, 0x46, 0x61, 0x76, 0x6f, 0x72, 0x69, 0x74,
	0x65, 0x56, 0x69, 0x64, 0x65, 0x6f, 0x49, 0x64, 0x4c, 0x69, 0x73, 0x74, 0x12, 0x28, 0x2e, 0x72,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x2e, 0x64, 0x6f, 0x75, 0x79, 0x69, 0x6e, 0x5f, 0x66, 0x61,
	0x76, 0x6f, 0x72, 0x69, 0x74, 0x65, 0x5f, 0x69, 0x64, 0x5f, 0x6c, 0x69, 0x73, 0x74, 0x5f, 0x72,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x2a, 0x2e, 0x72, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x2e, 0x64, 0x6f, 0x75, 0x79, 0x69, 0x6e, 0x5f, 0x66, 0x61, 0x76, 0x6f, 0x72, 0x69, 0x74,
	0x65, 0x5f, 0x69, 0x64, 0x5f, 0x6c, 0x69, 0x73, 0x74, 0x5f, 0x72, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x73, 0x65, 0x12, 0x80, 0x01, 0x0a, 0x1e, 0x47, 0x65, 0x74, 0x55, 0x73, 0x65, 0x72, 0x4c, 0x69,
	0x73, 0x74, 0x46, 0x61, 0x76, 0x6f, 0x72, 0x69, 0x74, 0x65, 0x56, 0x69, 0x64, 0x65, 0x6f, 0x49,
	0x64, 0x4c, 0x69, 0x73, 0x74, 0x12, 0x2d, 0x2e, 0x72, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x2e,
	0x64, 0x6f, 0x75, 0x79, 0x69, 0x6e, 0x5f, 0x66, 0x61, 0x76, 0x6f, 0x72, 0x69, 0x74, 0x65, 0x5f,
	0x6c, 0x69, 0x73, 0x74, 0x5f, 0x69, 0x64, 0x5f, 0x6c, 0x69, 0x73, 0x74, 0x5f, 0x72, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x1a, 0x2f, 0x2e, 0x72, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x2e,
	0x64, 0x6f, 0x75, 0x79, 0x69, 0x6e, 0x5f, 0x66, 0x61, 0x76, 0x6f, 0x72, 0x69, 0x74, 0x65, 0x5f,
	0x6c, 0x69, 0x73, 0x74, 0x5f, 0x69, 0x64, 0x5f, 0x6c, 0x69, 0x73, 0x74, 0x5f, 0x72, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x7d, 0x0a, 0x15, 0x47, 0x65, 0x74, 0x55, 0x73, 0x65, 0x72,
	0x46, 0x61, 0x76, 0x6f, 0x72, 0x69, 0x74, 0x65, 0x64, 0x43, 0x6f, 0x75, 0x6e, 0x74, 0x12, 0x30,
	0x2e, 0x72, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x2e, 0x64, 0x6f, 0x75, 0x79, 0x69, 0x6e, 0x5f,
	0x67, 0x65, 0x74, 0x5f, 0x75, 0x73, 0x65, 0x72, 0x5f, 0x66, 0x61, 0x76, 0x6f, 0x72, 0x69, 0x74,
	0x65, 0x64, 0x5f, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x5f, 0x72, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x1a, 0x32, 0x2e, 0x72, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x2e, 0x64, 0x6f, 0x75, 0x79,
	0x69, 0x6e, 0x5f, 0x67, 0x65, 0x74, 0x5f, 0x75, 0x73, 0x65, 0x72, 0x5f, 0x66, 0x61, 0x76, 0x6f,
	0x72, 0x69, 0x74, 0x65, 0x64, 0x5f, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x5f, 0x72, 0x65, 0x73, 0x70,
	0x6f, 0x6e, 0x73, 0x65, 0x12, 0x64, 0x0a, 0x10, 0x47, 0x65, 0x74, 0x46, 0x61, 0x76, 0x6f, 0x72,
	0x69, 0x74, 0x65, 0x43, 0x6f, 0x75, 0x6e, 0x74, 0x12, 0x26, 0x2e, 0x72, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x2e, 0x64, 0x6f, 0x75, 0x79, 0x69, 0x6e, 0x5f, 0x66, 0x61, 0x76, 0x6f, 0x72, 0x69,
	0x74, 0x65, 0x5f, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x5f, 0x72, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x1a, 0x28, 0x2e, 0x72, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x2e, 0x64, 0x6f, 0x75, 0x79,
	0x69, 0x6e, 0x5f, 0x66, 0x61, 0x76, 0x6f, 0x72, 0x69, 0x74, 0x65, 0x5f, 0x63, 0x6f, 0x75, 0x6e,
	0x74, 0x5f, 0x72, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x42, 0x1b, 0x5a, 0x19, 0x64, 0x6f,
	0x75, 0x79, 0x69, 0x6e, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x66, 0x61, 0x76, 0x6f, 0x72,
	0x69, 0x74, 0x65, 0x2f, 0x61, 0x70, 0x69, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var file_favorite_api_favorite_proto_goTypes = []interface{}{
	(*request.DouyinFavoriteActionRequest)(nil),          // 0: request.douyin_favorite_action_request
	(*request.DouyinFavoriteListRequest)(nil),            // 1: request.douyin_favorite_list_request
	(*request.DouyinFavoriteIdListRequest)(nil),          // 2: request.douyin_favorite_id_list_request
	(*request.DouyinFavoriteListIdListRequest)(nil),      // 3: request.douyin_favorite_list_id_list_request
	(*request.DouyinGetUserFavoritedCountRequest)(nil),   // 4: request.douyin_get_user_favorited_count_request
	(*request.DouyinFavoriteCountRequest)(nil),           // 5: request.douyin_favorite_count_request
	(*response.DouyinFavoriteActionResponse)(nil),        // 6: response.douyin_favorite_action_response
	(*response.DouyinFavoriteListResponse)(nil),          // 7: response.douyin_favorite_list_response
	(*response.DouyinFavoriteIdListResponse)(nil),        // 8: response.douyin_favorite_id_list_response
	(*response.DouyinFavoriteListIdListResponse)(nil),    // 9: response.douyin_favorite_list_id_list_response
	(*response.DouyinGetUserFavoritedCountResponse)(nil), // 10: response.douyin_get_user_favorited_count_response
	(*response.DouyinFavoriteCountResponse)(nil),         // 11: response.douyin_favorite_count_response
}
var file_favorite_api_favorite_proto_depIdxs = []int32{
	0,  // 0: api.Favorite.FavoriteVideo:input_type -> request.douyin_favorite_action_request
	1,  // 1: api.Favorite.GetFavoriteVideoList:input_type -> request.douyin_favorite_list_request
	2,  // 2: api.Favorite.GetUserFavoriteVideoIdList:input_type -> request.douyin_favorite_id_list_request
	3,  // 3: api.Favorite.GetUserListFavoriteVideoIdList:input_type -> request.douyin_favorite_list_id_list_request
	4,  // 4: api.Favorite.GetUserFavoritedCount:input_type -> request.douyin_get_user_favorited_count_request
	5,  // 5: api.Favorite.GetFavoriteCount:input_type -> request.douyin_favorite_count_request
	6,  // 6: api.Favorite.FavoriteVideo:output_type -> response.douyin_favorite_action_response
	7,  // 7: api.Favorite.GetFavoriteVideoList:output_type -> response.douyin_favorite_list_response
	8,  // 8: api.Favorite.GetUserFavoriteVideoIdList:output_type -> response.douyin_favorite_id_list_response
	9,  // 9: api.Favorite.GetUserListFavoriteVideoIdList:output_type -> response.douyin_favorite_list_id_list_response
	10, // 10: api.Favorite.GetUserFavoritedCount:output_type -> response.douyin_get_user_favorited_count_response
	11, // 11: api.Favorite.GetFavoriteCount:output_type -> response.douyin_favorite_count_response
	6,  // [6:12] is the sub-list for method output_type
	0,  // [0:6] is the sub-list for method input_type
	0,  // [0:0] is the sub-list for extension type_name
	0,  // [0:0] is the sub-list for extension extendee
	0,  // [0:0] is the sub-list for field type_name
}

func init() { file_favorite_api_favorite_proto_init() }
func file_favorite_api_favorite_proto_init() {
	if File_favorite_api_favorite_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_favorite_api_favorite_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   0,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_favorite_api_favorite_proto_goTypes,
		DependencyIndexes: file_favorite_api_favorite_proto_depIdxs,
	}.Build()
	File_favorite_api_favorite_proto = out.File
	file_favorite_api_favorite_proto_rawDesc = nil
	file_favorite_api_favorite_proto_goTypes = nil
	file_favorite_api_favorite_proto_depIdxs = nil
}
