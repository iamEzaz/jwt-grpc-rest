package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"context"

	"jwt-grpc-rest/api/controllers"
	"jwt-grpc-rest/api/proto"
	"jwt-grpc-rest/api/seed"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var server = controllers.Server{}

func init() {
	// loads values from .env into the system
	if err := godotenv.Load(); err != nil {
		log.Print("sad .env file found")
	}
}

func main() {

	router := mux.NewRouter()
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request){
		dialOption := grpc.WithTransportCredentials(insecure.NewCredentials())
		serviceConnection, err := grpc.Dial("localhost:3005", dialOption)
		if err != nil{
			panic(err)
		}

		serviceClient := proto.NewServiceClient(serviceConnection)

		res, err :=serviceClient.Login(context.Background(), &proto.Req{Email: "ezaz@gmail.com", Password: "ezaz@1234"})
		if err !=nil{
			panic(err)
		}
		fmt.Fprint(w, res.Email, res.Password)
	})

	var err error
	err = godotenv.Load()
	if err != nil {
		log.Fatalf("Error getting env, %v", err)
	} else {
		fmt.Println("We are getting the env values")
	}

	server.Initialize(os.Getenv("DB_DRIVER"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_PORT"), os.Getenv("DB_HOST"), os.Getenv("DB_NAME"))

	seed.Load(server.DB)

	server.Run(":8080")
	
	fmt.Println("http running on 8080")
	if err := http.ListenAndServe(":8080", router); err != nil{
		panic(err)
	}

	

}
