package mergeRouter

import (
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v4/pgxpool"

	loginRouter "sample-backend-go/internal/interface/routers/login"
	postRouters "sample-backend-go/internal/interface/routers/posts"
	userRouters "sample-backend-go/internal/interface/routers/users"
	vehicleRouters "sample-backend-go/internal/interface/routers/vehicles"
)

func SetupMergedRouter(pgPool *pgxpool.Pool) *mux.Router {
	router := mux.NewRouter()

	v1Router := router.PathPrefix("/v1").Subrouter()

	loginRouter.SetupLoginRouter(v1Router, pgPool)
	vehicleRouters.SetupVehicleRouter(v1Router, pgPool)
	userRouters.SetupUserRouter(v1Router, pgPool)
	postRouters.SetupPostRouter(v1Router, pgPool)

	return router
}
