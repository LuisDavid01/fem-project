package app

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/LuisDavid01/femProject/internal/api"
	"github.com/LuisDavid01/femProject/internal/middleware"
	"github.com/LuisDavid01/femProject/internal/store"
	"github.com/LuisDavid01/femProject/migrations"
)

type Application struct {
	Logger        *log.Logger
	WorkoutHander *api.WorkoutHandler
	UserHandler   *api.UserHandler
	TokenHandler  *api.TokenHandler
	Middleware    middleware.UserMiddleware
	DB            *sql.DB
}

func NewApplication() (*Application, error) {
	pgDb, err := store.Open()
	if err != nil {
		return nil, err
	}
	err = store.MigrateFS(pgDb, migrations.FS, ".")
	if err != nil {
		panic(err)
	}
	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)
	// our handler will go here
	workoutStore := store.NewPostgresWorkoutStore(pgDb)
	userStore := store.NewPostgresUserStore(pgDb)
	tokenStore := store.NewPostgresTokenStore(pgDb)
	workoutHandler := api.NewWorkoutHandler(workoutStore, logger)
	userHandler := api.NewUserHandler(userStore, logger)
	tokenHandler := api.NewTokenHandler(tokenStore, userStore, logger)
	middlewareHandler := middleware.UserMiddleware{UserStore: userStore}
	app := &Application{
		Logger:        logger,
		WorkoutHander: workoutHandler,
		UserHandler:   userHandler,
		TokenHandler:  tokenHandler,
		Middleware:    middlewareHandler,
		DB:            pgDb,
	}
	return app, nil
}

func (a *Application) HealthCheck(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Status is avaliable\n")
}
