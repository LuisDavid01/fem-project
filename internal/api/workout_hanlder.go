package api

import (
	"net/http"
	"strconv"
	"fmt"
	"github.com/go-chi/chi/v5"
)

type WorkoutHandler struct{}

func NewWorkoutHandler() *WorkoutHandler {
	return &WorkoutHandler{}
}

func (wh *WorkoutHandler) HandlerWorkoutById(w http.ResponseWriter, r *http.Request) {
	paramsWorkoutId := chi.URLParam(r, "id")
	if paramsWorkoutId == ""  {
		http.NotFound(w, r)
		return
	}
	workoutId, err := strconv.ParseInt(paramsWorkoutId, 10, 64)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	fmt.Fprintf(w, "this is the workout with id %d", workoutId)
}

func (wh *WorkoutHandler) HandlerCreateWorkout(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "workout created successfully\n")
}
