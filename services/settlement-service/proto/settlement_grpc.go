package proto

import (
	"context"

	"google.golang.org/grpc"
)

type SettlementServiceServer interface {
	ComputeSettlements(context.Context, *ComputeSettlementsRequest) (*ComputeSettlementsResponse, error)
	RecordSettlement(context.Context, *RecordSettlementRequest) (*RecordSettlementResponse, error)
	MarkSettlementCompleted(context.Context, *MarkSettlementCompletedRequest) (*MarkSettlementCompletedResponse, error)
	ListSettlements(context.Context, *ListSettlementsRequest) (*ListSettlementsResponse, error)
	mustEmbedUnimplementedSettlementServiceServer()
}

type UnimplementedSettlementServiceServer struct{}

func (UnimplementedSettlementServiceServer) ComputeSettlements(context.Context, *ComputeSettlementsRequest) (*ComputeSettlementsResponse, error) {
	return nil, nil
}
func (UnimplementedSettlementServiceServer) RecordSettlement(context.Context, *RecordSettlementRequest) (*RecordSettlementResponse, error) {
	return nil, nil
}
func (UnimplementedSettlementServiceServer) MarkSettlementCompleted(context.Context, *MarkSettlementCompletedRequest) (*MarkSettlementCompletedResponse, error) {
	return nil, nil
}
func (UnimplementedSettlementServiceServer) ListSettlements(context.Context, *ListSettlementsRequest) (*ListSettlementsResponse, error) {
	return nil, nil
}
func (UnimplementedSettlementServiceServer) mustEmbedUnimplementedSettlementServiceServer() {}

type SettlementServiceClient interface {
	ComputeSettlements(ctx context.Context, in *ComputeSettlementsRequest, opts ...grpc.CallOption) (*ComputeSettlementsResponse, error)
	RecordSettlement(ctx context.Context, in *RecordSettlementRequest, opts ...grpc.CallOption) (*RecordSettlementResponse, error)
	MarkSettlementCompleted(ctx context.Context, in *MarkSettlementCompletedRequest, opts ...grpc.CallOption) (*MarkSettlementCompletedResponse, error)
	ListSettlements(ctx context.Context, in *ListSettlementsRequest, opts ...grpc.CallOption) (*ListSettlementsResponse, error)
}

type settlementServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewSettlementServiceClient(cc grpc.ClientConnInterface) SettlementServiceClient {
	return &settlementServiceClient{cc}
}

func (c *settlementServiceClient) ComputeSettlements(ctx context.Context, in *ComputeSettlementsRequest, opts ...grpc.CallOption) (*ComputeSettlementsResponse, error) {
	out := new(ComputeSettlementsResponse)
	err := c.cc.Invoke(ctx, "/settlement.SettlementService/ComputeSettlements", in, out, opts...)
	return out, err
}

func (c *settlementServiceClient) RecordSettlement(ctx context.Context, in *RecordSettlementRequest, opts ...grpc.CallOption) (*RecordSettlementResponse, error) {
	out := new(RecordSettlementResponse)
	err := c.cc.Invoke(ctx, "/settlement.SettlementService/RecordSettlement", in, out, opts...)
	return out, err
}

func (c *settlementServiceClient) MarkSettlementCompleted(ctx context.Context, in *MarkSettlementCompletedRequest, opts ...grpc.CallOption) (*MarkSettlementCompletedResponse, error) {
	out := new(MarkSettlementCompletedResponse)
	err := c.cc.Invoke(ctx, "/settlement.SettlementService/MarkSettlementCompleted", in, out, opts...)
	return out, err
}

func (c *settlementServiceClient) ListSettlements(ctx context.Context, in *ListSettlementsRequest, opts ...grpc.CallOption) (*ListSettlementsResponse, error) {
	out := new(ListSettlementsResponse)
	err := c.cc.Invoke(ctx, "/settlement.SettlementService/ListSettlements", in, out, opts...)
	return out, err
}

func RegisterSettlementServiceServer(s grpc.ServiceRegistrar, srv SettlementServiceServer) {
	s.RegisterService(&SettlementService_ServiceDesc, srv)
}

var SettlementService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "settlement.SettlementService",
	HandlerType: (*SettlementServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{MethodName: "ComputeSettlements", Handler: _SettlementService_ComputeSettlements_Handler},
		{MethodName: "RecordSettlement", Handler: _SettlementService_RecordSettlement_Handler},
		{MethodName: "MarkSettlementCompleted", Handler: _SettlementService_MarkSettlementCompleted_Handler},
		{MethodName: "ListSettlements", Handler: _SettlementService_ListSettlements_Handler},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "proto/settlement/settlement.proto",
}

func _SettlementService_ComputeSettlements_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ComputeSettlementsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SettlementServiceServer).ComputeSettlements(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/settlement.SettlementService/ComputeSettlements"}
	return interceptor(ctx, in, info, func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SettlementServiceServer).ComputeSettlements(ctx, req.(*ComputeSettlementsRequest))
	})
}

func _SettlementService_RecordSettlement_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RecordSettlementRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SettlementServiceServer).RecordSettlement(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/settlement.SettlementService/RecordSettlement"}
	return interceptor(ctx, in, info, func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SettlementServiceServer).RecordSettlement(ctx, req.(*RecordSettlementRequest))
	})
}

func _SettlementService_MarkSettlementCompleted_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MarkSettlementCompletedRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SettlementServiceServer).MarkSettlementCompleted(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/settlement.SettlementService/MarkSettlementCompleted"}
	return interceptor(ctx, in, info, func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SettlementServiceServer).MarkSettlementCompleted(ctx, req.(*MarkSettlementCompletedRequest))
	})
}

func _SettlementService_ListSettlements_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListSettlementsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SettlementServiceServer).ListSettlements(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/settlement.SettlementService/ListSettlements"}
	return interceptor(ctx, in, info, func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SettlementServiceServer).ListSettlements(ctx, req.(*ListSettlementsRequest))
	})
}
