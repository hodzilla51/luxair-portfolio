package routers

import (
	"net/http"
	"sample-backend-go/internal/domain/entities/config"
	"sample-backend-go/internal/infrastructure/middleware"
	controllers "sample-backend-go/internal/interface/controllers/posts"
	repositories "sample-backend-go/internal/interface/repositories/posts"
	usecase "sample-backend-go/internal/usecase/interactors/posts"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v4/pgxpool"
)

func SetupPostRouter(router *mux.Router, pool *pgxpool.Pool) {

	repo := repositories.NewPostRepository(pool)
	postUseCase := usecase.NewPostUsecase(repo)
	postController := controllers.NewPostController(postUseCase)

	lr := repositories.NewLikeRepository(pool)
	likeUseCase := usecase.NewLikeUseCase(lr)
	likeController := controllers.NewLikeController(likeUseCase)

	cfg, _ := config.NewConfig()                            // configインスタンスを生成
	authMiddleware := middleware.AuthMiddlewareFactory(cfg) // ファクトリ関数を使用してMiddlewareを生成

	router.Handle("/posts", authMiddleware(http.HandlerFunc(postController.PostSubmitController))).Methods("POST")
	router.Handle("/posts/{postId}/likes", authMiddleware(http.HandlerFunc(likeController.ToggleLikeController))).Methods("POST")

	router.HandleFunc("/posts", postController.PostGetController).Methods("GET")
	router.HandleFunc("/posts?userID={userID}", postController.PostGetUserController).Methods("GET")
	router.HandleFunc("/posts?follow=true", postController.PostGetFollowController).Methods("GET")
	router.HandleFunc("/posts/{postID}/likes", likeController.GetLikes).Methods("GET")
}
