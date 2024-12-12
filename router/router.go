package router

import (
	"otus-highload/handlers/friend"
	"otus-highload/handlers/post"
	"otus-highload/handlers/user"

	"github.com/gorilla/mux"
)

func NewRouter() *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/user/register", user.RegisterUserHandler).Methods("POST")
	router.HandleFunc("/user/get/{id}", user.GetUserHandler).Methods("GET")
	router.HandleFunc("/user/login", user.LoginHandler).Methods("POST")
	router.HandleFunc("/user/search", user.SearchUserHandler).Methods("GET")

	router.HandleFunc("/friend/set/{id}", friend.SetFriendHandler).Methods("PUT")
	router.HandleFunc("/friend/delete/{id}", friend.DeleteFriendHandler).Methods("PUT")

	router.HandleFunc("/post/create", post.CreatePostHandler).Methods("POST")
	router.HandleFunc("/post/update", post.UpdatePostHandler).Methods("PUT")
	router.HandleFunc("/post/delete/{id}", post.DeletePostHandler).Methods("PUT")
	router.HandleFunc("/post/get/{id}", post.GetPostHandler).Methods("GET")
	router.HandleFunc("/post/feed", post.FeedHandler).Methods("GET")

	return router
}
