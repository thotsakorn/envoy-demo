package main

import (
	"context"
	"log"

	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	pb "github.com/grpc-ecosystem/go-grpc-middleware/testing/testproto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func parseToken(token string) (struct{}, error) {
	return struct{}{}, nil
}

func userClaimFromToken(struct{}) string {
	return "foobar"
}

// exampleAuthFunc is used by a middleware to authenticate requests
func exampleAuthFunc(ctx context.Context) (context.Context, error) {
	println("@@@@@@@@@@")
	token, err := grpc_auth.AuthFromMD(ctx, "bearer")
	log.Println("token  ", token)
	if err != nil {
		return nil, err
	}

	tokenInfo, err := parseToken(token)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "invalid auth token: %v", err)
	}
	log.Println("testAuth  ===>>>>   ", tokenInfo)

	grpc_ctxtags.Extract(ctx).Set("auth.sub", userClaimFromToken(tokenInfo))

	newCtx := context.WithValue(ctx, "tokenInfo", tokenInfo)
	return newCtx, nil
}

type server struct {
	pb.UnimplementedTestServiceServer
	message string
}

// SayHello only can be called by client when authenticated by exampleAuthFunc
func (g *server) Ping(ctx context.Context, request *pb.PingRequest) (*pb.PingResponse, error) {
	return &pb.PingResponse{Value: g.message}, nil
}

// AuthFuncOverride is called instead of exampleAuthFunc
func (g *server) AuthFuncOverride(ctx context.Context, fullMethodName string) (context.Context, error) {
	log.Println("client is calling method:", fullMethodName)
	return ctx, nil
}

func main() {
	svr := grpc.NewServer(
		grpc.StreamInterceptor(grpc_auth.StreamServerInterceptor(exampleAuthFunc)),
		grpc.UnaryInterceptor(grpc_auth.UnaryServerInterceptor(exampleAuthFunc)),
	)
	log.Println("start ===>", exampleAuthFunc)
	overrideActive := true

	if overrideActive {
		pb.RegisterTestServiceServer(svr, &server{message: "pong unauthenticated"})
	} else {
		pb.RegisterTestServiceServer(svr, &server{message: "pong authenticated"})
	}
}
