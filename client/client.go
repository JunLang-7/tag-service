package main

import (
	"context"
	"log"

	pb "github.com/JunLang-7/tag-service/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	ctx := context.Background()
	clientConn, err := GetClientConn(ctx, "localhost:8001", nil)
	if err != nil {
		log.Fatal(err)
	}
	defer clientConn.Close()

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
