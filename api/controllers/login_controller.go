package controllers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"fmt"
	"context"

	"github.com/iamEzaz/jwt-grpc-rest/api/auth"
	"github.com/iamEzaz/jwt-grpc-rest/api/models"
	"github.com/iamEzaz/jwt-grpc-rest/api/responses"
	"github.com/iamEzaz/jwt-grpc-rest/api/utils/formaterror"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"github.com/iamEzaz/jwt-grpc-rest/api/proto"
)

func (server *Server) Login(w http.ResponseWriter, r *http.Request) {

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
		return
	}
	user := models.User{}
	err = json.Unmarshal(body, &user)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	user.Prepare()
	err = user.Validate("login")
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	token, err := server.SignIn(user.Email, user.Password)
	if err != nil {
		formattedError := formaterror.FormatError(err.Error())
		responses.ERROR(w, http.StatusUnprocessableEntity, formattedError)
		return
	}
	responses.JSON(w, http.StatusOK, token)
}

func (server *Server) SignIn(email, password string) (string, error) {

	var err error

	user := models.User{}

	err = server.DB.Debug().Model(models.User{}).Where("email = ?", email).Take(&user).Error
	if err != nil {
		return "", err
	}
	err = models.VerifyPassword(user.Password, password)
	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
		return "", err
	}
	return auth.CreateToken(user.ID)
}
