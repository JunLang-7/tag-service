package server

import (
	"context"
	"encoding/json"

	"github.com/JunLang-7/tag-service/pkg/bapi"
	"github.com/JunLang-7/tag-service/pkg/errcode"
	pb "github.com/JunLang-7/tag-service/proto"
	"google.golang.org/grpc/metadata"
)

type Auth struct{}

func (a *Auth) GetAppKey() string {
	return "go-programming-tour-book"
}

func (a *Auth) GetAppSecret() string {
	return "go-programming-tour-book"
}

func (a *Auth) Check(ctx context.Context) error {
	md, _ := metadata.FromIncomingContext(ctx)

	var appKey, appSecret string
	if value, ok := md["app_key"]; ok {
		appKey = value[0]
	}
	if value, ok := md["app_secret"]; ok {
		appSecret = value[0]
	}
	if appKey != a.GetAppKey() || appSecret != a.GetAppSecret() {
		return errcode.TogRPCError(errcode.Unauthorized)
	}
	return nil
}

type TagServer struct {
	pb.UnimplementedTagServiceServer
	auth *Auth
}

func NewTagServer() *TagServer {
	return &TagServer{}
}

func (t *TagServer) GetTagList(ctx context.Context, r *pb.GetTagListRequest) (*pb.GetTagListReply, error) {
	if err := t.auth.Check(ctx); err != nil {
		return nil, err
	}
	api := bapi.NewAPI("http://localhost:8000")
	body, err := api.GetTagList(ctx, r.GetName())
	if err != nil {
		return nil, errcode.TogRPCError(errcode.ErrorGetTagListFail)
	}

	var tagList pb.GetTagListReply
	err = json.Unmarshal(body, &tagList)
	if err != nil {
		return nil, errcode.TogRPCError(errcode.Fail)
	}

	return &tagList, nil
}
