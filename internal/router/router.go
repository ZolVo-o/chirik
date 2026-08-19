package router

import (
	"chirik/internal/handlers"
	"chirik/internal/middleware"
	"net/http"
)

type Router struct {
	mux *http.ServeMux
}

func New() *Router {
	return &Router{mux: http.NewServeMux()}
}

func (r *Router) HandleFunc(path string, handler http.HandlerFunc) {
	r.mux.HandleFunc(path, handler)
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.mux.ServeHTTP(w, req)
}

func Setup(
	authHandler *handlers.AuthHandler,
	tweetHandler *handlers.TweetHandler,
	followHandler *handlers.FollowHandler,
	messageHandler *handlers.MessageHandler,
	usersHandler *handlers.UsersHandler,
) *Router {
	r := New()

	// Публичные
	r.HandleFunc("/api/auth/register", authHandler.Register)
	r.HandleFunc("/api/auth/login", authHandler.Login)

	// Защищенные
	r.HandleFunc("/api/users/profile", func(w http.ResponseWriter, req *http.Request) {
		if req.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		middleware.Auth(authHandler.Profile)(w, req)
	})
	r.HandleFunc("/api/users/update", func(w http.ResponseWriter, req *http.Request) {
		if req.Method != http.MethodPut && req.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		middleware.Auth(authHandler.UpdateProfile)(w, req)
	})
	r.HandleFunc("/api/users/search", middleware.Auth(usersHandler.Search))

	r.HandleFunc("/api/tweets", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			tweetHandler.GetAllTweets(w, r)
		} else if r.Method == http.MethodPost {
			middleware.Auth(tweetHandler.CreateTweet)(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	r.HandleFunc("/api/tweets/user/", tweetHandler.GetTweetsByUser)
	r.HandleFunc("/api/tweets/like/", middleware.Auth(tweetHandler.LikeTweet))
	r.HandleFunc("/api/tweets/view/", middleware.Auth(tweetHandler.ViewTweet))

	r.HandleFunc("/api/follow/", middleware.Auth(followHandler.Follow))
	r.HandleFunc("/api/unfollow/", middleware.Auth(followHandler.Unfollow))
	r.HandleFunc("/api/following/", followHandler.GetFollowing)
	r.HandleFunc("/api/followers/", followHandler.GetFollowers)
	r.HandleFunc("/api/conversations", middleware.Auth(messageHandler.Conversations))
	r.HandleFunc("/api/conversations/", middleware.Auth(messageHandler.ConversationMessages))
	r.HandleFunc("/api/messages/", func(w http.ResponseWriter, req *http.Request) {
		if req.Method == http.MethodPut {
			middleware.Auth(messageHandler.UpdateMessage)(w, req)
			return
		}
		if req.Method == http.MethodDelete {
			middleware.Auth(messageHandler.DeleteMessage)(w, req)
			return
		}
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	})

	return r
}
