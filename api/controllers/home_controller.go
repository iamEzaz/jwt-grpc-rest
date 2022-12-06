package controllers

import (
	"fmt"
	"net/http"
	"context"

	

	"github.com/iamEzaz/jwt-grpc-rest/api/proto"
	"github.com/iamEzaz/jwt-grpc-rest/api/responses"
	


	
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)


 

func (server *Server) Home(w http.ResponseWriter, r *http.Request) {

		dialOption := grpc.WithTransportCredentials(insecure.NewCredentials())
		serviceConnection, err := grpc.Dial("localhost:3005", dialOption)
		if err != nil{
			panic(err)
		}

		serviceClient := proto.NewServiceClient(serviceConnection)

		res, err :=serviceClient.Home(context.Background(), &proto.Req{Email: "Name: Ezazul Islam ", Password: "Pass: ezaz@0987 "})
		if err !=nil{
			panic(err)
		}
		fmt.Fprint(w, res.Email, res.Password)
	
	responses.JSON(w, http.StatusOK, " Welcome To This Awesome API ")
}