package routes

import (
	"github.com/LuisDavid01/femProject/internal/app"
	"github.com/go-chi/chi/v5"
)

func SetupRoutes(app *app.Application) *chi.Mux {
	r := chi.NewRouter()
	r.Group(func(r chi.Router) {
		r.Use(app.Middleware.Authenticate)
		r.Get("/workouts/{id}", app.Middleware.RequireUser(app.WorkoutHander.HandlerWorkoutById))

		r.Delete("/workouts/{id}", app.Middleware.RequireUser(app.WorkoutHander.HandlerDeleteWorkoutById))

		r.Put("/workouts/{id}", app.Middleware.RequireUser(app.WorkoutHander.HandleUpdateWorkoutById))

		r.Post("/workouts", app.Middleware.RequireUser(app.WorkoutHander.HandlerCreateWorkout))

	})
	// Health check route
	r.Get("/health", app.HealthCheck)
	r.Post("/users/register", app.UserHandler.HandlerRegisterUser)
	r.Post("/tokens/authentication", app.TokenHandler.HandleCreateToken)
	return r
}
