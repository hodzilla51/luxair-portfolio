package routers

import (
	controllers "sample-backend-go/internal/interface/controllers/login"
	repositories "sample-backend-go/internal/interface/repositories/users"
	usecase "sample-backend-go/internal/usecase/interactors/users"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v4/pgxpool"
)

func SetupLoginRouter(router *mux.Router, pool *pgxpool.Pool) {
	repo := repositories.NewUserRepository(pool)
	uu := usecase.NewUserUsecase(repo)
	lc := controllers.NewUserController(uu)

	router.HandleFunc("/login", lc.LoginEndpoint).Methods("POST")
	router.HandleFunc("/login/refresh", lc.RefreshToken).Methods("POST")
}
