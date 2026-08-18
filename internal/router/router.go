package router

import (
    "net/http"
    "chirik/internal/handlers"
    "chirik/internal/middleware"
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
) *Router {
    r := New()

    // Публичные
    r.HandleFunc("/api/auth/register", authHandler.Register)
    r.HandleFunc("/api/auth/login", authHandler.Login)

    // Защищенные
    r.HandleFunc("/api/users/profile", middleware.Auth(authHandler.Profile))
    r.HandleFunc("/api/users/update", middleware.Auth(authHandler.UpdateProfile))
    
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
    
    r.HandleFunc("/api/follow/", middleware.Auth(followHandler.Follow))
    r.HandleFunc("/api/unfollow/", middleware.Auth(followHandler.Unfollow))
    r.HandleFunc("/api/following/", followHandler.GetFollowing)
    r.HandleFunc("/api/followers/", followHandler.GetFollowers)

    return r
}
