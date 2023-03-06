package main

import (
	"context"
	"fmt"
	"log"
	"net"

	livescore "testgrpcenvoy/pb"
	livescore2 "testgrpcenvoy/pb/pbscore2"

	"google.golang.org/grpc"
)

var matches livescore.ListMatchesResponse

type liveScoreServer struct {
	livescore.UnimplementedScoreServiceServer
	livescore2.UnimplementedScore2ServiceServer
}

func (lss *liveScoreServer) ListMatches(ctx context.Context, req *livescore.ListMatchesRequest) (*livescore.ListMatchesResponse, error) {
	println("Connect List")
	cc, err := grpc.Dial("localhost:8084", grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}
	defer cc.Close()
	client := livescore2.NewScore2ServiceClient(cc)
	request := &livescore2.List2MatchesRequest{Country: "brian"}

	match2, err2 := client.List2Matches(ctx, request)

	if err2 != nil {
		log.Fatal(err2)
	}

	println("match2 :::::    ", match2.String())
	match := &livescore.MatchScoreResponse{
		Score: "4:1",
		Live:  true,
	}
	matches.Scores = matches.Scores[:0]
	matches.Scores = append(matches.Scores, match)
	return &matches, nil
}

const addr = ":50004"

func main() {
	//create tcp connection on port 50004
	conn, err := net.Listen("tcp", addr)

	if err != nil {
		log.Fatal("tcp connection err: ", err.Error())
	}
	//create grpc server
	grpcServer := grpc.NewServer()

	server := liveScoreServer{}
	//register our livescore service
	livescore.RegisterScoreServiceServer(grpcServer, &server)

	fmt.Println("Starting gRPC server at : ", addr)
	//serve our connection
	if err := grpcServer.Serve(conn); err != nil {
		log.Fatal(err)
	}
}
