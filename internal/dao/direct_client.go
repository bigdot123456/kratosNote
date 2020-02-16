package dao //modify here, copy from api/client.go
import (
	//callapi "callServer/api"
	"context"
	//"fmt"

	"github.com/bilibili/kratos/pkg/net/rpc/warden"

	svrapi "callServer/smallapi" // <-direct use server proto api client.go content

	"google.golang.org/grpc"
)

// AppID .
const AppID = "TODO: ADD APP ID"

const target = "direct://default/127.0.0.1:9000"

// NewClient new grpc client
func NewClient(cfg *warden.ClientConfig, opts ...grpc.DialOption) (svrapi.AsmallsClient, error) {
	client := warden.NewClient(cfg, opts...)
	//	cc, err := client.Dial(context.Background(), fmt.Sprintf("discovery://default/%s", AppID))
	cc, err := client.Dial(context.Background(), target)
	if err != nil {
		return nil, err
	}
	return svrapi.NewAsmallsClient(cc), nil
}

// 生成 gRPC 代码
//go:generate kratos tool protoc --grpc --bm api.proto
