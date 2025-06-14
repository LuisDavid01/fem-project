package routes

import (
	"github.com/go-chi/chi/v5"
	"github.com/LuisDavid01/femProject/internal/app"
)

func SetupRoutes(app *app.Application) *chi.Mux {
	r := chi.NewRouter()

	// Health check route
	r.Get("/health", app.HealthCheck)
	r.Get("/workouts/{id}", app.WorkoutHander.HandlerWorkoutById)

	r.Post("/workouts", app.WorkoutHander.HandlerCreateWorkout)
	return r
}
