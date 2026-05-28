package main

import (
	"context"
	"log"

	"github.com/JunLang-7/tag-service/internal/middleware"
	pb "github.com/JunLang-7/tag-service/proto"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	ctx := context.Background()
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithUnaryInterceptor(
		grpc_middleware.ChainUnaryClient(middleware.UnaryContextTimeout()),
	))
	opts = append(opts, grpc.WithStreamInterceptor(
		grpc_middleware.ChainStreamClient(middleware.StreamContextTimeout()),
	))
	clientConn, err := GetClientConn(ctx, "localhost:8004", opts)
	tagClient := pb.NewTagServiceClient(clientConn)
	resp, err := tagClient.GetTagList(ctx, &pb.GetTagListRequest{Name: "Golang"})
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("resp: %v", resp)
}

func GetClientConn(ctx context.Context, target string, opts []grpc.DialOption) (*grpc.ClientConn, error) {
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	return grpc.NewClient(target, opts...)
}
