package proto

import (
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
)

type Group struct {
	Id          string                 `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Name        string                 `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	Description string                 `protobuf:"bytes,3,opt,name=description,proto3" json:"description,omitempty"`
	Currency    string                 `protobuf:"bytes,4,opt,name=currency,proto3" json:"currency,omitempty"`
	InviteCode  string                 `protobuf:"bytes,5,opt,name=invite_code,json=inviteCode,proto3" json:"invite_code,omitempty"`
	OwnerId     string                 `protobuf:"bytes,6,opt,name=owner_id,json=ownerId,proto3" json:"owner_id,omitempty"`
	CreatedAt   *timestamppb.Timestamp `protobuf:"bytes,7,opt,name=created_at,json=createdAt,proto3" json:"created_at,omitempty"`
	UpdatedAt   *timestamppb.Timestamp `protobuf:"bytes,8,opt,name=updated_at,json=updatedAt,proto3" json:"updated_at,omitempty"`
	Members     []*GroupMember         `protobuf:"bytes,9,rep,name=members,proto3" json:"members,omitempty"`
}

type GroupMember struct {
	UserId   string                 `protobuf:"bytes,1,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
	Role     string                 `protobuf:"bytes,2,opt,name=role,proto3" json:"role,omitempty"`
	JoinedAt *timestamppb.Timestamp `protobuf:"bytes,3,opt,name=joined_at,json=joinedAt,proto3" json:"joined_at,omitempty"`
}

type CreateGroupRequest struct {
	Name        string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Description string `protobuf:"bytes,2,opt,name=description,proto3" json:"description,omitempty"`
	Currency    string `protobuf:"bytes,3,opt,name=currency,proto3" json:"currency,omitempty"`
	OwnerId     string `protobuf:"bytes,4,opt,name=owner_id,json=ownerId,proto3" json:"owner_id,omitempty"`
}

type CreateGroupResponse struct {
	Group *Group `protobuf:"bytes,1,opt,name=group,proto3" json:"group,omitempty"`
}

type GetGroupRequest struct {
	Id string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
}

type GetGroupResponse struct {
	Group *Group `protobuf:"bytes,1,opt,name=group,proto3" json:"group,omitempty"`
}

type ListUserGroupsRequest struct {
	UserId string `protobuf:"bytes,1,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
}

type ListUserGroupsResponse struct {
	Groups []*Group `protobuf:"bytes,1,rep,name=groups,proto3" json:"groups,omitempty"`
}

type UpdateGroupRequest struct {
	Id          string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Name        string `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	Description string `protobuf:"bytes,3,opt,name=description,proto3" json:"description,omitempty"`
	Currency    string `protobuf:"bytes,4,opt,name=currency,proto3" json:"currency,omitempty"`
}

type UpdateGroupResponse struct {
	Group *Group `protobuf:"bytes,1,opt,name=group,proto3" json:"group,omitempty"`
}

type DeleteGroupRequest struct {
	Id     string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	UserId string `protobuf:"bytes,2,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
}

type DeleteGroupResponse struct {
	Success bool `protobuf:"varint,1,opt,name=success,proto3" json:"success,omitempty"`
}

type AddMemberRequest struct {
	GroupId string `protobuf:"bytes,1,opt,name=group_id,json=groupId,proto3" json:"group_id,omitempty"`
	UserId  string `protobuf:"bytes,2,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
	Role    string `protobuf:"bytes,3,opt,name=role,proto3" json:"role,omitempty"`
}

type AddMemberResponse struct {
	Group *Group `protobuf:"bytes,1,opt,name=group,proto3" json:"group,omitempty"`
}

type RemoveMemberRequest struct {
	GroupId string `protobuf:"bytes,1,opt,name=group_id,json=groupId,proto3" json:"group_id,omitempty"`
	UserId  string `protobuf:"bytes,2,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
}

type RemoveMemberResponse struct {
	Group *Group `protobuf:"bytes,1,opt,name=group,proto3" json:"group,omitempty"`
}

type GenerateInviteCodeRequest struct {
	GroupId string `protobuf:"bytes,1,opt,name=group_id,json=groupId,proto3" json:"group_id,omitempty"`
	UserId  string `protobuf:"bytes,2,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
}

type GenerateInviteCodeResponse struct {
	InviteCode string `protobuf:"bytes,1,opt,name=invite_code,json=inviteCode,proto3" json:"invite_code,omitempty"`
}

type JoinByInviteCodeRequest struct {
	InviteCode string `protobuf:"bytes,1,opt,name=invite_code,json=inviteCode,proto3" json:"invite_code,omitempty"`
	UserId     string `protobuf:"bytes,2,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
}

type JoinByInviteCodeResponse struct {
	Group *Group `protobuf:"bytes,1,opt,name=group,proto3" json:"group,omitempty"`
}
