package proto

import (
	"context"

	"google.golang.org/grpc"
)

type GroupServiceServer interface {
	CreateGroup(context.Context, *CreateGroupRequest) (*CreateGroupResponse, error)
	GetGroup(context.Context, *GetGroupRequest) (*GetGroupResponse, error)
	ListUserGroups(context.Context, *ListUserGroupsRequest) (*ListUserGroupsResponse, error)
	UpdateGroup(context.Context, *UpdateGroupRequest) (*UpdateGroupResponse, error)
	DeleteGroup(context.Context, *DeleteGroupRequest) (*DeleteGroupResponse, error)
	AddMember(context.Context, *AddMemberRequest) (*AddMemberResponse, error)
	RemoveMember(context.Context, *RemoveMemberRequest) (*RemoveMemberResponse, error)
	GenerateInviteCode(context.Context, *GenerateInviteCodeRequest) (*GenerateInviteCodeResponse, error)
	JoinByInviteCode(context.Context, *JoinByInviteCodeRequest) (*JoinByInviteCodeResponse, error)
	mustEmbedUnimplementedGroupServiceServer()
}

type UnimplementedGroupServiceServer struct{}

func (UnimplementedGroupServiceServer) CreateGroup(context.Context, *CreateGroupRequest) (*CreateGroupResponse, error) {
	return nil, nil
}
func (UnimplementedGroupServiceServer) GetGroup(context.Context, *GetGroupRequest) (*GetGroupResponse, error) {
	return nil, nil
}
func (UnimplementedGroupServiceServer) ListUserGroups(context.Context, *ListUserGroupsRequest) (*ListUserGroupsResponse, error) {
	return nil, nil
}
func (UnimplementedGroupServiceServer) UpdateGroup(context.Context, *UpdateGroupRequest) (*UpdateGroupResponse, error) {
	return nil, nil
}
func (UnimplementedGroupServiceServer) DeleteGroup(context.Context, *DeleteGroupRequest) (*DeleteGroupResponse, error) {
	return nil, nil
}
func (UnimplementedGroupServiceServer) AddMember(context.Context, *AddMemberRequest) (*AddMemberResponse, error) {
	return nil, nil
}
func (UnimplementedGroupServiceServer) RemoveMember(context.Context, *RemoveMemberRequest) (*RemoveMemberResponse, error) {
	return nil, nil
}
func (UnimplementedGroupServiceServer) GenerateInviteCode(context.Context, *GenerateInviteCodeRequest) (*GenerateInviteCodeResponse, error) {
	return nil, nil
}
func (UnimplementedGroupServiceServer) JoinByInviteCode(context.Context, *JoinByInviteCodeRequest) (*JoinByInviteCodeResponse, error) {
	return nil, nil
}
func (UnimplementedGroupServiceServer) mustEmbedUnimplementedGroupServiceServer() {}

type UnsafeGroupServiceServer interface {
	mustEmbedUnimplementedGroupServiceServer()
}

type GroupServiceClient interface {
	CreateGroup(ctx context.Context, in *CreateGroupRequest, opts ...grpc.CallOption) (*CreateGroupResponse, error)
	GetGroup(ctx context.Context, in *GetGroupRequest, opts ...grpc.CallOption) (*GetGroupResponse, error)
	ListUserGroups(ctx context.Context, in *ListUserGroupsRequest, opts ...grpc.CallOption) (*ListUserGroupsResponse, error)
	UpdateGroup(ctx context.Context, in *UpdateGroupRequest, opts ...grpc.CallOption) (*UpdateGroupResponse, error)
	DeleteGroup(ctx context.Context, in *DeleteGroupRequest, opts ...grpc.CallOption) (*DeleteGroupResponse, error)
	AddMember(ctx context.Context, in *AddMemberRequest, opts ...grpc.CallOption) (*AddMemberResponse, error)
	RemoveMember(ctx context.Context, in *RemoveMemberRequest, opts ...grpc.CallOption) (*RemoveMemberResponse, error)
	GenerateInviteCode(ctx context.Context, in *GenerateInviteCodeRequest, opts ...grpc.CallOption) (*GenerateInviteCodeResponse, error)
	JoinByInviteCode(ctx context.Context, in *JoinByInviteCodeRequest, opts ...grpc.CallOption) (*JoinByInviteCodeResponse, error)
}

type groupServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewGroupServiceClient(cc grpc.ClientConnInterface) GroupServiceClient {
	return &groupServiceClient{cc}
}

func (c *groupServiceClient) CreateGroup(ctx context.Context, in *CreateGroupRequest, opts ...grpc.CallOption) (*CreateGroupResponse, error) {
	out := new(CreateGroupResponse)
	err := c.cc.Invoke(ctx, "/group.GroupService/CreateGroup", in, out, opts...)
	return out, err
}

func (c *groupServiceClient) GetGroup(ctx context.Context, in *GetGroupRequest, opts ...grpc.CallOption) (*GetGroupResponse, error) {
	out := new(GetGroupResponse)
	err := c.cc.Invoke(ctx, "/group.GroupService/GetGroup", in, out, opts...)
	return out, err
}

func (c *groupServiceClient) ListUserGroups(ctx context.Context, in *ListUserGroupsRequest, opts ...grpc.CallOption) (*ListUserGroupsResponse, error) {
	out := new(ListUserGroupsResponse)
	err := c.cc.Invoke(ctx, "/group.GroupService/ListUserGroups", in, out, opts...)
	return out, err
}

func (c *groupServiceClient) UpdateGroup(ctx context.Context, in *UpdateGroupRequest, opts ...grpc.CallOption) (*UpdateGroupResponse, error) {
	out := new(UpdateGroupResponse)
	err := c.cc.Invoke(ctx, "/group.GroupService/UpdateGroup", in, out, opts...)
	return out, err
}

func (c *groupServiceClient) DeleteGroup(ctx context.Context, in *DeleteGroupRequest, opts ...grpc.CallOption) (*DeleteGroupResponse, error) {
	out := new(DeleteGroupResponse)
	err := c.cc.Invoke(ctx, "/group.GroupService/DeleteGroup", in, out, opts...)
	return out, err
}

func (c *groupServiceClient) AddMember(ctx context.Context, in *AddMemberRequest, opts ...grpc.CallOption) (*AddMemberResponse, error) {
	out := new(AddMemberResponse)
	err := c.cc.Invoke(ctx, "/group.GroupService/AddMember", in, out, opts...)
	return out, err
}

func (c *groupServiceClient) RemoveMember(ctx context.Context, in *RemoveMemberRequest, opts ...grpc.CallOption) (*RemoveMemberResponse, error) {
	out := new(RemoveMemberResponse)
	err := c.cc.Invoke(ctx, "/group.GroupService/RemoveMember", in, out, opts...)
	return out, err
}

func (c *groupServiceClient) GenerateInviteCode(ctx context.Context, in *GenerateInviteCodeRequest, opts ...grpc.CallOption) (*GenerateInviteCodeResponse, error) {
	out := new(GenerateInviteCodeResponse)
	err := c.cc.Invoke(ctx, "/group.GroupService/GenerateInviteCode", in, out, opts...)
	return out, err
}

func (c *groupServiceClient) JoinByInviteCode(ctx context.Context, in *JoinByInviteCodeRequest, opts ...grpc.CallOption) (*JoinByInviteCodeResponse, error) {
	out := new(JoinByInviteCodeResponse)
	err := c.cc.Invoke(ctx, "/group.GroupService/JoinByInviteCode", in, out, opts...)
	return out, err
}

func RegisterGroupServiceServer(s grpc.ServiceRegistrar, srv GroupServiceServer) {
	s.RegisterService(&GroupService_ServiceDesc, srv)
}

var GroupService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "group.GroupService",
	HandlerType: (*GroupServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{MethodName: "CreateGroup", Handler: _GroupService_CreateGroup_Handler},
		{MethodName: "GetGroup", Handler: _GroupService_GetGroup_Handler},
		{MethodName: "ListUserGroups", Handler: _GroupService_ListUserGroups_Handler},
		{MethodName: "UpdateGroup", Handler: _GroupService_UpdateGroup_Handler},
		{MethodName: "DeleteGroup", Handler: _GroupService_DeleteGroup_Handler},
		{MethodName: "AddMember", Handler: _GroupService_AddMember_Handler},
		{MethodName: "RemoveMember", Handler: _GroupService_RemoveMember_Handler},
		{MethodName: "GenerateInviteCode", Handler: _GroupService_GenerateInviteCode_Handler},
		{MethodName: "JoinByInviteCode", Handler: _GroupService_JoinByInviteCode_Handler},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "proto/group/group.proto",
}

func _GroupService_CreateGroup_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateGroupRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GroupServiceServer).CreateGroup(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/group.GroupService/CreateGroup"}
	return interceptor(ctx, in, info, func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GroupServiceServer).CreateGroup(ctx, req.(*CreateGroupRequest))
	})
}

func _GroupService_GetGroup_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetGroupRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GroupServiceServer).GetGroup(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/group.GroupService/GetGroup"}
	return interceptor(ctx, in, info, func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GroupServiceServer).GetGroup(ctx, req.(*GetGroupRequest))
	})
}

func _GroupService_ListUserGroups_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListUserGroupsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GroupServiceServer).ListUserGroups(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/group.GroupService/ListUserGroups"}
	return interceptor(ctx, in, info, func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GroupServiceServer).ListUserGroups(ctx, req.(*ListUserGroupsRequest))
	})
}

func _GroupService_UpdateGroup_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateGroupRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GroupServiceServer).UpdateGroup(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/group.GroupService/UpdateGroup"}
	return interceptor(ctx, in, info, func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GroupServiceServer).UpdateGroup(ctx, req.(*UpdateGroupRequest))
	})
}

func _GroupService_DeleteGroup_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteGroupRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GroupServiceServer).DeleteGroup(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/group.GroupService/DeleteGroup"}
	return interceptor(ctx, in, info, func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GroupServiceServer).DeleteGroup(ctx, req.(*DeleteGroupRequest))
	})
}

func _GroupService_AddMember_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AddMemberRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GroupServiceServer).AddMember(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/group.GroupService/AddMember"}
	return interceptor(ctx, in, info, func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GroupServiceServer).AddMember(ctx, req.(*AddMemberRequest))
	})
}

func _GroupService_RemoveMember_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RemoveMemberRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GroupServiceServer).RemoveMember(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/group.GroupService/RemoveMember"}
	return interceptor(ctx, in, info, func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GroupServiceServer).RemoveMember(ctx, req.(*RemoveMemberRequest))
	})
}

func _GroupService_GenerateInviteCode_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GenerateInviteCodeRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GroupServiceServer).GenerateInviteCode(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/group.GroupService/GenerateInviteCode"}
	return interceptor(ctx, in, info, func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GroupServiceServer).GenerateInviteCode(ctx, req.(*GenerateInviteCodeRequest))
	})
}

func _GroupService_JoinByInviteCode_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(JoinByInviteCodeRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GroupServiceServer).JoinByInviteCode(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/group.GroupService/JoinByInviteCode"}
	return interceptor(ctx, in, info, func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GroupServiceServer).JoinByInviteCode(ctx, req.(*JoinByInviteCodeRequest))
	})
}
