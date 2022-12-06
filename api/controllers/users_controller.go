package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"context"

	
	"github.com/gorilla/mux"
	"github.com/iamEzaz/jwt-grpc-rest/api/auth"
	"github.com/iamEzaz/jwt-grpc-rest/api/models"
	"github.com/iamEzaz/jwt-grpc-rest/api/responses"
	"github.com/iamEzaz/jwt-grpc-rest/api/utils/formaterror"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"github.com/iamEzaz/jwt-grpc-rest/api/proto"
)

func (server *Server) CreateUser(w http.ResponseWriter, r *http.Request) {

	dialOption := grpc.WithTransportCredentials(insecure.NewCredentials())
		serviceConnection, err := grpc.Dial("localhost:3005", dialOption)
		if err != nil{
			panic(err)
		}

		serviceClient := proto.NewServiceClient(serviceConnection)

		res, err :=serviceClient.Home(context.Background(), &proto.Req{Email: "ezaz@gmail.com", Password: "ezaz@1234"})
		if err !=nil{
			panic(err)
		}
		fmt.Fprint(w, res.Email, res.Password)

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
	}
	user := models.User{}
	err = json.Unmarshal(body, &user)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	user.Prepare()
	err = user.Validate("")
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	userCreated, err := user.SaveUser(server.DB)

	if err != nil {

		formattedError := formaterror.FormatError(err.Error())

		responses.ERROR(w, http.StatusInternalServerError, formattedError)
		return
	}
	w.Header().Set("Location", fmt.Sprintf("%s%s/%d", r.Host, r.RequestURI, userCreated.ID))
	responses.JSON(w, http.StatusCreated, userCreated)
}

func (server *Server) GetUsers(w http.ResponseWriter, r *http.Request) {

	dialOption := grpc.WithTransportCredentials(insecure.NewCredentials())
		serviceConnection, err := grpc.Dial("localhost:3005", dialOption)
		if err != nil{
			panic(err)
		}

		serviceClient := proto.NewServiceClient(serviceConnection)

		res, err :=serviceClient.Home(context.Background(), &proto.Req{Email: "ezaz@gmail.com", Password: "ezaz@1234"})
		if err !=nil{
			panic(err)
		}
		fmt.Fprint(w, res.Email, res.Password)

	user := models.User{}

	users, err := user.FindAllUsers(server.DB)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	responses.JSON(w, http.StatusOK, users)
}

func (server *Server) GetUser(w http.ResponseWriter, r *http.Request) {

	dialOption := grpc.WithTransportCredentials(insecure.NewCredentials())
		serviceConnection, err := grpc.Dial("localhost:3005", dialOption)
		if err != nil{
			panic(err)
		}

		serviceClient := proto.NewServiceClient(serviceConnection)

		res, err :=serviceClient.Home(context.Background(), &proto.Req{Email: "ezaz@gmail.com", Password: "ezaz@1234"})
		if err !=nil{
			panic(err)
		}
		fmt.Fprint(w, res.Email, res.Password)

	vars := mux.Vars(r)
	uid, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}
	user := models.User{}
	userGotten, err := user.FindUserByID(server.DB, uint32(uid))
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}
	responses.JSON(w, http.StatusOK, userGotten)
}

func (server *Server) UpdateUser(w http.ResponseWriter, r *http.Request) {

	dialOption := grpc.WithTransportCredentials(insecure.NewCredentials())
		serviceConnection, err := grpc.Dial("localhost:3005", dialOption)
		if err != nil{
			panic(err)
		}

		serviceClient := proto.NewServiceClient(serviceConnection)

		res, err :=serviceClient.Home(context.Background(), &proto.Req{Email: "ezaz@gmail.com", Password: "ezaz@1234"})
		if err !=nil{
			panic(err)
		}
		fmt.Fprint(w, res.Email, res.Password)

	vars := mux.Vars(r)
	uid, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	user := models.User{}
	err = json.Unmarshal(body, &user)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	tokenID, err := auth.ExtractTokenID(r)
	if err != nil {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}
	if tokenID != uint32(uid) {
		responses.ERROR(w, http.StatusUnauthorized, errors.New(http.StatusText(http.StatusUnauthorized)))
		return
	}
	user.Prepare()
	err = user.Validate("update")
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	updatedUser, err := user.UpdateAUser(server.DB, uint32(uid))
	if err != nil {
		formattedError := formaterror.FormatError(err.Error())
		responses.ERROR(w, http.StatusInternalServerError, formattedError)
		return
	}
	responses.JSON(w, http.StatusOK, updatedUser)
}

func (server *Server) DeleteUser(w http.ResponseWriter, r *http.Request) {

	dialOption := grpc.WithTransportCredentials(insecure.NewCredentials())
		serviceConnection, err := grpc.Dial("localhost:3005", dialOption)
		if err != nil{
			panic(err)
		}

		serviceClient := proto.NewServiceClient(serviceConnection)

		res, err :=serviceClient.Home(context.Background(), &proto.Req{Email: "ezaz@gmail.com", Password: "ezaz@1234"})
		if err !=nil{
			panic(err)
		}
		fmt.Fprint(w, res.Email, res.Password)

	vars := mux.Vars(r)

	user := models.User{}

	uid, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}
	tokenID, err := auth.ExtractTokenID(r)
	if err != nil {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}
	if tokenID != 0 && tokenID != uint32(uid) {
		responses.ERROR(w, http.StatusUnauthorized, errors.New(http.StatusText(http.StatusUnauthorized)))
		return
	}
	_, err = user.DeleteAUser(server.DB, uint32(uid))
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	w.Header().Set("Entity", fmt.Sprintf("%d", uid))
	responses.JSON(w, http.StatusNoContent, "")
}
