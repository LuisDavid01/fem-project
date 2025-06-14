package app

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/LuisDavid01/femProject/internal/api"
)
type Application struct {
	Logger *log.Logger
	WorkoutHander *api.WorkoutHandler
}

func NewApplication() (*Application, error) {
	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)
	// our handler will go here
	workoutHandler := api.NewWorkoutHandler()
	app := &Application{
		Logger: logger,
		WorkoutHander: workoutHandler,
	}
	return app, nil
}


func (a *Application) HealthCheck(w http.ResponseWriter, r *http.Request){
	fmt.Fprintf(w, "Status is avaliable\n")
}
