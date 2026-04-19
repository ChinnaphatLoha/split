package handler

import (
	"context"
	"log"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/split/services/group-service/internal/model"
	"github.com/split/services/group-service/internal/service"
	pb "github.com/split/services/group-service/proto"
)

type GroupHandler struct {
	pb.UnimplementedGroupServiceServer
	svc *service.GroupService
}

func NewGroupHandler(svc *service.GroupService) *GroupHandler {
	return &GroupHandler{svc: svc}
}

func (h *GroupHandler) CreateGroup(ctx context.Context, req *pb.CreateGroupRequest) (*pb.CreateGroupResponse, error) {
	if req.Name == "" || req.OwnerId == "" {
		return nil, status.Error(codes.InvalidArgument, "name and owner_id are required")
	}

	group, err := h.svc.CreateGroup(ctx, req.Name, req.Description, req.Currency, req.OwnerId)
	if err != nil {
		log.Printf("create group error: %v", err)
		return nil, status.Error(codes.Internal, "failed to create group")
	}

	return &pb.CreateGroupResponse{Group: groupToProto(group)}, nil
}

func (h *GroupHandler) GetGroup(ctx context.Context, req *pb.GetGroupRequest) (*pb.GetGroupResponse, error) {
	if req.Id == "" {
		return nil, status.Error(codes.InvalidArgument, "group id is required")
	}

	group, err := h.svc.GetGroup(ctx, req.Id)
	if err != nil {
		if err == service.ErrGroupNotFound {
			return nil, status.Error(codes.NotFound, "group not found")
		}
		log.Printf("get group error: %v", err)
		return nil, status.Error(codes.Internal, "failed to get group")
	}

	return &pb.GetGroupResponse{Group: groupToProto(group)}, nil
}

func (h *GroupHandler) ListUserGroups(ctx context.Context, req *pb.ListUserGroupsRequest) (*pb.ListUserGroupsResponse, error) {
	if req.UserId == "" {
		return nil, status.Error(codes.InvalidArgument, "user_id is required")
	}

	groups, err := h.svc.ListUserGroups(ctx, req.UserId)
	if err != nil {
		log.Printf("list user groups error: %v", err)
		return nil, status.Error(codes.Internal, "failed to list groups")
	}

	var pbGroups []*pb.Group
	for _, g := range groups {
		pbGroups = append(pbGroups, groupToProto(g))
	}

	return &pb.ListUserGroupsResponse{Groups: pbGroups}, nil
}

func (h *GroupHandler) UpdateGroup(ctx context.Context, req *pb.UpdateGroupRequest) (*pb.UpdateGroupResponse, error) {
	if req.Id == "" {
		return nil, status.Error(codes.InvalidArgument, "group id is required")
	}

	group, err := h.svc.UpdateGroup(ctx, req.Id, req.Name, req.Description, req.Currency)
	if err != nil {
		if err == service.ErrGroupNotFound {
			return nil, status.Error(codes.NotFound, "group not found")
		}
		log.Printf("update group error: %v", err)
		return nil, status.Error(codes.Internal, "failed to update group")
	}

	return &pb.UpdateGroupResponse{Group: groupToProto(group)}, nil
}

func (h *GroupHandler) DeleteGroup(ctx context.Context, req *pb.DeleteGroupRequest) (*pb.DeleteGroupResponse, error) {
	if req.Id == "" || req.UserId == "" {
		return nil, status.Error(codes.InvalidArgument, "group id and user_id are required")
	}

	err := h.svc.DeleteGroup(ctx, req.Id, req.UserId)
	if err != nil {
		if err == service.ErrGroupNotFound {
			return nil, status.Error(codes.NotFound, "group not found")
		}
		if err == service.ErrNotOwner {
			return nil, status.Error(codes.PermissionDenied, "only the owner can delete the group")
		}
		log.Printf("delete group error: %v", err)
		return nil, status.Error(codes.Internal, "failed to delete group")
	}

	return &pb.DeleteGroupResponse{Success: true}, nil
}

func (h *GroupHandler) AddMember(ctx context.Context, req *pb.AddMemberRequest) (*pb.AddMemberResponse, error) {
	if req.GroupId == "" || req.UserId == "" {
		return nil, status.Error(codes.InvalidArgument, "group_id and user_id are required")
	}

	group, err := h.svc.AddMember(ctx, req.GroupId, req.UserId, req.Role)
	if err != nil {
		log.Printf("add member error: %v", err)
		return nil, status.Error(codes.Internal, "failed to add member")
	}

	return &pb.AddMemberResponse{Group: groupToProto(group)}, nil
}

func (h *GroupHandler) RemoveMember(ctx context.Context, req *pb.RemoveMemberRequest) (*pb.RemoveMemberResponse, error) {
	if req.GroupId == "" || req.UserId == "" {
		return nil, status.Error(codes.InvalidArgument, "group_id and user_id are required")
	}

	group, err := h.svc.RemoveMember(ctx, req.GroupId, req.UserId)
	if err != nil {
		if err == service.ErrCannotRemoveOwner {
			return nil, status.Error(codes.PermissionDenied, "cannot remove the group owner")
		}
		log.Printf("remove member error: %v", err)
		return nil, status.Error(codes.Internal, "failed to remove member")
	}

	return &pb.RemoveMemberResponse{Group: groupToProto(group)}, nil
}

func (h *GroupHandler) GenerateInviteCode(ctx context.Context, req *pb.GenerateInviteCodeRequest) (*pb.GenerateInviteCodeResponse, error) {
	if req.GroupId == "" || req.UserId == "" {
		return nil, status.Error(codes.InvalidArgument, "group_id and user_id are required")
	}

	code, err := h.svc.GenerateInviteCode(ctx, req.GroupId, req.UserId)
	if err != nil {
		if err == service.ErrNotOwner {
			return nil, status.Error(codes.PermissionDenied, "only the owner can generate invite codes")
		}
		log.Printf("generate invite code error: %v", err)
		return nil, status.Error(codes.Internal, "failed to generate invite code")
	}

	return &pb.GenerateInviteCodeResponse{InviteCode: code}, nil
}

func (h *GroupHandler) JoinByInviteCode(ctx context.Context, req *pb.JoinByInviteCodeRequest) (*pb.JoinByInviteCodeResponse, error) {
	if req.InviteCode == "" || req.UserId == "" {
		return nil, status.Error(codes.InvalidArgument, "invite_code and user_id are required")
	}

	group, err := h.svc.JoinByInviteCode(ctx, req.InviteCode, req.UserId)
	if err != nil {
		if err == service.ErrInvalidInvite {
			return nil, status.Error(codes.NotFound, "invalid invite code")
		}
		log.Printf("join by invite code error: %v", err)
		return nil, status.Error(codes.Internal, "failed to join group")
	}

	return &pb.JoinByInviteCodeResponse{Group: groupToProto(group)}, nil
}

func groupToProto(g *model.Group) *pb.Group {
	pbGroup := &pb.Group{
		Id:          g.ID,
		Name:        g.Name,
		Description: g.Description,
		Currency:    g.Currency,
		InviteCode:  g.InviteCode,
		OwnerId:     g.OwnerID,
		CreatedAt:   timestamppb.New(g.CreatedAt),
		UpdatedAt:   timestamppb.New(g.UpdatedAt),
	}

	for _, m := range g.Members {
		pbGroup.Members = append(pbGroup.Members, &pb.GroupMember{
			UserId:   m.UserID,
			Role:     m.Role,
			JoinedAt: timestamppb.New(m.JoinedAt),
		})
	}

	return pbGroup
}
