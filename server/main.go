package main

import (
	"context"
	"fmt"

	"github.com/pulumi/pulumi/sdk/v3/go/common/util/rpcutil"
	pulumirpc "github.com/pulumi/pulumi/sdk/v3/proto/go"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

func main() {
	resmon := &resmon{}

	handle, err := rpcutil.ServeWithOptions(rpcutil.ServeOptions{
		Port: 8888,
		Init: func(srv *grpc.Server) error {
			pulumirpc.RegisterResourceMonitorServer(srv, resmon)
			return nil
		},
	})
	if err != nil {
		panic(err)
	}

	fmt.Printf("port = %d\n", handle.Port)

	<-handle.Done
}

type resmon struct {
	pulumirpc.UnsafeResourceMonitorServer
}

func (rm *resmon) SupportsFeature(ctx context.Context,
	req *pulumirpc.SupportsFeatureRequest,
) (*pulumirpc.SupportsFeatureResponse, error) {
	return &pulumirpc.SupportsFeatureResponse{
		HasSupport: true,
	}, nil
}

func (rm *resmon) Call(ctx context.Context, req *pulumirpc.ResourceCallRequest) (*pulumirpc.CallResponse, error) {
	panic("not implemented")
}

func (rm *resmon) Invoke(ctx context.Context, req *pulumirpc.ResourceInvokeRequest) (*pulumirpc.InvokeResponse, error) {
	panic("not implemented")
}

func (rm *resmon) ReadResource(ctx context.Context,
	req *pulumirpc.ReadResourceRequest,
) (*pulumirpc.ReadResourceResponse, error) {
	panic("not implemented")
}

func (rm *resmon) RegisterResource(ctx context.Context,
	req *pulumirpc.RegisterResourceRequest,
) (*pulumirpc.RegisterResourceResponse, error) {
	panic("not implemented")
}

func (rm *resmon) RegisterResourceOutputs(ctx context.Context,
	req *pulumirpc.RegisterResourceOutputsRequest,
) (*emptypb.Empty, error) {
	panic("not implemented")
}

func (rm *resmon) RegisterStackTransform(ctx context.Context,
	req *pulumirpc.Callback,
) (*emptypb.Empty, error) {
	panic("not implemented")
}

func (rm *resmon) StreamInvoke(
	req *pulumirpc.ResourceInvokeRequest, stream pulumirpc.ResourceMonitor_StreamInvokeServer,
) error {
	panic("not implemented")
}
