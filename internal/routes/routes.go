package routes

import (
	"github.com/LuisDavid01/femProject/internal/app"
	"github.com/go-chi/chi/v5"
)

func SetupRoutes(app *app.Application) *chi.Mux {
	r := chi.NewRouter()

	// Health check route
	r.Get("/health", app.HealthCheck)
	r.Get("/workouts/{id}", app.WorkoutHander.HandlerWorkoutById)

	r.Delete("/workouts/{id}", app.WorkoutHander.HandlerDeleteWorkoutById)

	r.Put("/workouts/{id}", app.WorkoutHander.HandleUpdateWorkoutById)

	r.Post("/workouts", app.WorkoutHander.HandlerCreateWorkout)
	r.Post("/users/register", app.UserHandler.HandlerRegisterUser)
	return r
}
