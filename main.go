package main

import (
	"context"
	"fmt"
	"net"


	"google.golang.org/grpc"
	"jwt-grpc-rest/api/proto"
)


type pingService struct{
	proto.UnimplementedServiceServer
} 

func (p *pingService) Login(ctx context.Context, in *proto.Req) (*proto.Res, error){
	return &proto.Res{
		Email: in.GetEmail(),
		Password: in.GetPassword(),
	}, nil
}


func main() {

	grpcServer := grpc.NewServer()
	lis, err := net.Listen("tcp", ":3005")

	if err !=nil{
		panic(err)
	}

	proto.RegisterServiceServer(grpcServer, &pingService{})

	fmt.Println("running on port 3005")
	if err := grpcServer.Serve(lis); err !=nil {
		panic (err)
	}
	// 	api.Run()

}
