package main

import (
	"flag"
	"log"
	"net"
	"net/http"

	pb "github.com/JunLang-7/tag-service/proto"
	"github.com/JunLang-7/tag-service/server"
	"github.com/soheilhy/cmux"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var (
	port string
)

func init() {
	flag.StringVar(&port, "port", "8003", "server port")
	flag.Parse()
}

func main() {
	l, err := RunTCPServer(port)
	if err != nil {
		log.Fatalf("Run TCP Server err: %v", err)
	}

	m := cmux.New(l)
	grpcL := m.MatchWithWriters(cmux.HTTP2MatchHeaderFieldPrefixSendSettings("content-type", "application/grpc"))
	httpL := m.Match(cmux.HTTP1Fast())

	grpcS := RunGrpcServer()
	httpS := RunHttpServer(port)
	go func() {
		err := grpcS.Serve(grpcL)
		if err != nil {
			log.Fatalf("Run grpc server err: %v", err)
		}
	}()
	go func() {
		err := httpS.Serve(httpL)
		if err != nil {
			log.Fatalf("Run http server err: %v", err)
		}
	}()

	err = m.Serve()
	if err != nil {
		log.Fatalf("Run TCP Server err: %v", err)
	}
}

func RunTCPServer(port string) (net.Listener, error) {
	return net.Listen("tcp", ":"+port)
}

func RunHttpServer(port string) *http.Server {
	serveMux := http.NewServeMux()
	serveMux.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("pong"))
	})
	return &http.Server{
		Addr:    ":" + port,
		Handler: serveMux,
	}
}

func RunGrpcServer() *grpc.Server {
	s := grpc.NewServer()
	pb.RegisterTagServiceServer(s, server.NewTagServer())
	reflection.Register(s)
	return s
}
