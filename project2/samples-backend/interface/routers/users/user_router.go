package routers

import (
	"net/http"
	"sample-backend-go/internal/domain/entities/config"
	"sample-backend-go/internal/infrastructure/middleware"
	controllers "sample-backend-go/internal/interface/controllers/users"
	repositories "sample-backend-go/internal/interface/repositories/users"
	usecase "sample-backend-go/internal/usecase/interactors/users"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v4/pgxpool"
)

func SetupUserRouter(router *mux.Router, pool *pgxpool.Pool) {

	repo := repositories.NewUserRepository(pool)
	userUseCase := usecase.NewUserUsecase(repo)
	userController := controllers.NewUserController(userUseCase, repo)

	fr := repositories.NewFollowRepository(pool)
	followUsecase := usecase.NewFollowUseCase(fr)
	followController := controllers.NewFollowController(followUsecase)

	cfg, _ := config.NewConfig()                            // configインスタンスを生成
	authMiddleware := middleware.AuthMiddlewareFactory(cfg) // ファクトリ関数を使用してMiddlewareを生成

	// POST
	router.HandleFunc("/users", userController.UserRegister).Methods("POST")
	router.HandleFunc("/users/check/username", userController.UserCheckUsername).Methods("POST")
	router.HandleFunc("/users/check/email", userController.UserCheckEmail).Methods("POST")
	router.Handle("/users/follow", authMiddleware(http.HandlerFunc(followController.UserFollow))).Methods("POST")

	// GET
	router.HandleFunc("/users/{username}", userController.UserGetAll).Methods("GET")
	router.HandleFunc("/users/{userID}/basic", userController.UserGetBasic).Methods("GET")
	router.HandleFunc("/users/{userID}/detail", userController.UserGetDetail).Methods("GET")
	router.HandleFunc("/users/{userID}/vehicles", userController.UserGetVehicles).Methods("GET")
	router.HandleFunc("/users/{userID}/laptimes", userController.UserGetLaptimes).Methods("GET")
	router.HandleFunc("/users/{userID}/followings", followController.UserGetFollowings).Methods("GET")
	router.HandleFunc("/users/{userID}/followers", followController.UserGetFollowers).Methods("GET")

	// PATCH
	router.Handle("/users/me", authMiddleware(http.HandlerFunc(userController.UserUpdate))).Methods("PATCH")
	router.Handle("/users/me/profile", authMiddleware(http.HandlerFunc(userController.UserUpdateProfile))).Methods("PATCH")

}
