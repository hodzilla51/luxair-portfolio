package routers

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v4/pgxpool"

	"sample-backend-go/internal/domain/entities/config"
	"sample-backend-go/internal/infrastructure/middleware"
	controllers "sample-backend-go/internal/interface/controllers/vehicles"
	repositories "sample-backend-go/internal/interface/repositories/vehicles"
	usecase "sample-backend-go/internal/usecase/interactors/vehicles"
)

func SetupVehicleRouter(router *mux.Router, pool *pgxpool.Pool) {

	repo := repositories.NewVehicleRepository(pool)
	vehicleUseCase := usecase.NewVehicleUseCase(repo)
	vehicleController := controllers.NewVehicleController(vehicleUseCase)

	cfg, _ := config.NewConfig()                            // configインスタンスを生成
	authMiddleware := middleware.AuthMiddlewareFactory(cfg) // ファクトリ関数を使用してMiddlewareを生成

	router.Handle("/vehicles", authMiddleware(http.HandlerFunc(vehicleController.VehicleRegister))).Methods("POST")
	router.HandleFunc("/vehicles/{vehicleID}/basic", vehicleController.VehicleGetBasic).Methods("GET")
	router.HandleFunc("/vehicles/{vehicleID}/detail", vehicleController.VehicleGetDetail).Methods("GET")
	router.HandleFunc("/vehicles/{vehicleID}/users", vehicleController.VehicleGetUsers).Methods("GET")
	router.HandleFunc("/vehicles/{vehicleID}/fuels", vehicleController.VehicleGetFuels).Methods("GET")
	router.HandleFunc("/vehicles/{vehicleID}/laptimes", vehicleController.VehicleGetLaptimes).Methods("GET")
}
