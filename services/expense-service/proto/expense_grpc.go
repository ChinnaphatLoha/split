package proto

import (
	"context"

	"google.golang.org/grpc"
)

type ExpenseServiceServer interface {
	CreateExpense(context.Context, *CreateExpenseRequest) (*CreateExpenseResponse, error)
	GetExpense(context.Context, *GetExpenseRequest) (*GetExpenseResponse, error)
	ListGroupExpenses(context.Context, *ListGroupExpensesRequest) (*ListGroupExpensesResponse, error)
	UpdateExpense(context.Context, *UpdateExpenseRequest) (*UpdateExpenseResponse, error)
	DeleteExpense(context.Context, *DeleteExpenseRequest) (*DeleteExpenseResponse, error)
	GetGroupBalances(context.Context, *GetGroupBalancesRequest) (*GetGroupBalancesResponse, error)
	mustEmbedUnimplementedExpenseServiceServer()
}

type UnimplementedExpenseServiceServer struct{}

func (UnimplementedExpenseServiceServer) CreateExpense(context.Context, *CreateExpenseRequest) (*CreateExpenseResponse, error) {
	return nil, nil
}
func (UnimplementedExpenseServiceServer) GetExpense(context.Context, *GetExpenseRequest) (*GetExpenseResponse, error) {
	return nil, nil
}
func (UnimplementedExpenseServiceServer) ListGroupExpenses(context.Context, *ListGroupExpensesRequest) (*ListGroupExpensesResponse, error) {
	return nil, nil
}
func (UnimplementedExpenseServiceServer) UpdateExpense(context.Context, *UpdateExpenseRequest) (*UpdateExpenseResponse, error) {
	return nil, nil
}
func (UnimplementedExpenseServiceServer) DeleteExpense(context.Context, *DeleteExpenseRequest) (*DeleteExpenseResponse, error) {
	return nil, nil
}
func (UnimplementedExpenseServiceServer) GetGroupBalances(context.Context, *GetGroupBalancesRequest) (*GetGroupBalancesResponse, error) {
	return nil, nil
}
func (UnimplementedExpenseServiceServer) mustEmbedUnimplementedExpenseServiceServer() {}

type ExpenseServiceClient interface {
	CreateExpense(ctx context.Context, in *CreateExpenseRequest, opts ...grpc.CallOption) (*CreateExpenseResponse, error)
	GetExpense(ctx context.Context, in *GetExpenseRequest, opts ...grpc.CallOption) (*GetExpenseResponse, error)
	ListGroupExpenses(ctx context.Context, in *ListGroupExpensesRequest, opts ...grpc.CallOption) (*ListGroupExpensesResponse, error)
	UpdateExpense(ctx context.Context, in *UpdateExpenseRequest, opts ...grpc.CallOption) (*UpdateExpenseResponse, error)
	DeleteExpense(ctx context.Context, in *DeleteExpenseRequest, opts ...grpc.CallOption) (*DeleteExpenseResponse, error)
	GetGroupBalances(ctx context.Context, in *GetGroupBalancesRequest, opts ...grpc.CallOption) (*GetGroupBalancesResponse, error)
}

type expenseServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewExpenseServiceClient(cc grpc.ClientConnInterface) ExpenseServiceClient {
	return &expenseServiceClient{cc}
}

func (c *expenseServiceClient) CreateExpense(ctx context.Context, in *CreateExpenseRequest, opts ...grpc.CallOption) (*CreateExpenseResponse, error) {
	out := new(CreateExpenseResponse)
	err := c.cc.Invoke(ctx, "/expense.ExpenseService/CreateExpense", in, out, opts...)
	return out, err
}

func (c *expenseServiceClient) GetExpense(ctx context.Context, in *GetExpenseRequest, opts ...grpc.CallOption) (*GetExpenseResponse, error) {
	out := new(GetExpenseResponse)
	err := c.cc.Invoke(ctx, "/expense.ExpenseService/GetExpense", in, out, opts...)
	return out, err
}

func (c *expenseServiceClient) ListGroupExpenses(ctx context.Context, in *ListGroupExpensesRequest, opts ...grpc.CallOption) (*ListGroupExpensesResponse, error) {
	out := new(ListGroupExpensesResponse)
	err := c.cc.Invoke(ctx, "/expense.ExpenseService/ListGroupExpenses", in, out, opts...)
	return out, err
}

func (c *expenseServiceClient) UpdateExpense(ctx context.Context, in *UpdateExpenseRequest, opts ...grpc.CallOption) (*UpdateExpenseResponse, error) {
	out := new(UpdateExpenseResponse)
	err := c.cc.Invoke(ctx, "/expense.ExpenseService/UpdateExpense", in, out, opts...)
	return out, err
}

func (c *expenseServiceClient) DeleteExpense(ctx context.Context, in *DeleteExpenseRequest, opts ...grpc.CallOption) (*DeleteExpenseResponse, error) {
	out := new(DeleteExpenseResponse)
	err := c.cc.Invoke(ctx, "/expense.ExpenseService/DeleteExpense", in, out, opts...)
	return out, err
}

func (c *expenseServiceClient) GetGroupBalances(ctx context.Context, in *GetGroupBalancesRequest, opts ...grpc.CallOption) (*GetGroupBalancesResponse, error) {
	out := new(GetGroupBalancesResponse)
	err := c.cc.Invoke(ctx, "/expense.ExpenseService/GetGroupBalances", in, out, opts...)
	return out, err
}

func RegisterExpenseServiceServer(s grpc.ServiceRegistrar, srv ExpenseServiceServer) {
	s.RegisterService(&ExpenseService_ServiceDesc, srv)
}

var ExpenseService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "expense.ExpenseService",
	HandlerType: (*ExpenseServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{MethodName: "CreateExpense", Handler: _ExpenseService_CreateExpense_Handler},
		{MethodName: "GetExpense", Handler: _ExpenseService_GetExpense_Handler},
		{MethodName: "ListGroupExpenses", Handler: _ExpenseService_ListGroupExpenses_Handler},
		{MethodName: "UpdateExpense", Handler: _ExpenseService_UpdateExpense_Handler},
		{MethodName: "DeleteExpense", Handler: _ExpenseService_DeleteExpense_Handler},
		{MethodName: "GetGroupBalances", Handler: _ExpenseService_GetGroupBalances_Handler},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "proto/expense/expense.proto",
}

func _ExpenseService_CreateExpense_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateExpenseRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ExpenseServiceServer).CreateExpense(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/expense.ExpenseService/CreateExpense"}
	return interceptor(ctx, in, info, func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ExpenseServiceServer).CreateExpense(ctx, req.(*CreateExpenseRequest))
	})
}

func _ExpenseService_GetExpense_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetExpenseRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ExpenseServiceServer).GetExpense(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/expense.ExpenseService/GetExpense"}
	return interceptor(ctx, in, info, func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ExpenseServiceServer).GetExpense(ctx, req.(*GetExpenseRequest))
	})
}

func _ExpenseService_ListGroupExpenses_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListGroupExpensesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ExpenseServiceServer).ListGroupExpenses(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/expense.ExpenseService/ListGroupExpenses"}
	return interceptor(ctx, in, info, func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ExpenseServiceServer).ListGroupExpenses(ctx, req.(*ListGroupExpensesRequest))
	})
}

func _ExpenseService_UpdateExpense_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateExpenseRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ExpenseServiceServer).UpdateExpense(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/expense.ExpenseService/UpdateExpense"}
	return interceptor(ctx, in, info, func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ExpenseServiceServer).UpdateExpense(ctx, req.(*UpdateExpenseRequest))
	})
}

func _ExpenseService_DeleteExpense_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteExpenseRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ExpenseServiceServer).DeleteExpense(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/expense.ExpenseService/DeleteExpense"}
	return interceptor(ctx, in, info, func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ExpenseServiceServer).DeleteExpense(ctx, req.(*DeleteExpenseRequest))
	})
}

func _ExpenseService_GetGroupBalances_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetGroupBalancesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ExpenseServiceServer).GetGroupBalances(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/expense.ExpenseService/GetGroupBalances"}
	return interceptor(ctx, in, info, func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ExpenseServiceServer).GetGroupBalances(ctx, req.(*GetGroupBalancesRequest))
	})
}
