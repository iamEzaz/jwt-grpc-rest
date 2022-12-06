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

func (server *Server) CreatePost(w http.ResponseWriter, r *http.Request) {


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
	post := models.Post{}
	err = json.Unmarshal(body, &post)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	post.Prepare()
	err = post.Validate()
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	uid, err := auth.ExtractTokenID(r)
	if err != nil {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}
	if uid != post.AuthorID {
		responses.ERROR(w, http.StatusUnauthorized, errors.New(http.StatusText(http.StatusUnauthorized)))
		return
	}
	postCreated, err := post.SavePost(server.DB)
	if err != nil {
		formattedError := formaterror.FormatError(err.Error())
		responses.ERROR(w, http.StatusInternalServerError, formattedError)
		return
	}
	w.Header().Set("Location", fmt.Sprintf("%s%s/%d", r.Host, r.URL.Path, postCreated.ID))
	responses.JSON(w, http.StatusCreated, postCreated)
}

func (server *Server) GetPosts(w http.ResponseWriter, r *http.Request) {


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

	post := models.Post{}

	posts, err := post.FindAllPosts(server.DB)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	responses.JSON(w, http.StatusOK, posts)
}

func (server *Server) GetPost(w http.ResponseWriter, r *http.Request) {


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
	pid, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}
	post := models.Post{}

	postReceived, err := post.FindPostByID(server.DB, pid)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	responses.JSON(w, http.StatusOK, postReceived)
}

func (server *Server) UpdatePost(w http.ResponseWriter, r *http.Request) {


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

	// Check if the post id is valid
	pid, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	//CHeck if the auth token is valid and  get the user id from it
	uid, err := auth.ExtractTokenID(r)
	if err != nil {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}

	// Check if the post exist
	post := models.Post{}
	err = server.DB.Debug().Model(models.Post{}).Where("id = ?", pid).Take(&post).Error
	if err != nil {
		responses.ERROR(w, http.StatusNotFound, errors.New("Post not found"))
		return
	}

	// If a user attempt to update a post not belonging to him
	if uid != post.AuthorID {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}
	// Read the data posted
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	// Start processing the request data
	postUpdate := models.Post{}
	err = json.Unmarshal(body, &postUpdate)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	//Also check if the request user id is equal to the one gotten from token
	if uid != postUpdate.AuthorID {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}

	postUpdate.Prepare()
	err = postUpdate.Validate()
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	postUpdate.ID = post.ID //this is important to tell the model the post id to update, the other update field are set above

	postUpdated, err := postUpdate.UpdateAPost(server.DB)

	if err != nil {
		formattedError := formaterror.FormatError(err.Error())
		responses.ERROR(w, http.StatusInternalServerError, formattedError)
		return
	}
	responses.JSON(w, http.StatusOK, postUpdated)
}

func (server *Server) DeletePost(w http.ResponseWriter, r *http.Request) {

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

	// Is a valid post id given to us?
	pid, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	// Is this user authenticated?
	uid, err := auth.ExtractTokenID(r)
	if err != nil {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}

	// Check if the post exist
	post := models.Post{}
	err = server.DB.Debug().Model(models.Post{}).Where("id = ?", pid).Take(&post).Error
	if err != nil {
		responses.ERROR(w, http.StatusNotFound, errors.New("Unauthorized"))
		return
	}

	// Is the authenticated user, the owner of this post?
	if uid != post.AuthorID {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}
	_, err = post.DeleteAPost(server.DB, pid, uid)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}
	w.Header().Set("Entity", fmt.Sprintf("%d", pid))
	responses.JSON(w, http.StatusNoContent, "")
}
