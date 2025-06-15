package api

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/LuisDavid01/femProject/internal/store"
	"github.com/go-chi/chi/v5"
)

type WorkoutHandler struct {
	workoutStore store.WorkoutStore
}

func NewWorkoutHandler(workoutStore store.WorkoutStore) *WorkoutHandler {
	return &WorkoutHandler{
		workoutStore: workoutStore,
	}
}

func (wh *WorkoutHandler) HandlerWorkoutById(w http.ResponseWriter, r *http.Request) {
	paramsWorkoutId := chi.URLParam(r, "id")
	if paramsWorkoutId == "" {
		http.NotFound(w, r)
		return
	}
	workoutId, err := strconv.ParseInt(paramsWorkoutId, 10, 64)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	workout, err := wh.workoutStore.GetWorkoutById(workoutId)
	if err != nil {
		fmt.Println("Error getting workout by ID:", err)
		http.Error(w, "Failed to get workout", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(workout)
}

func (wh *WorkoutHandler) HandlerCreateWorkout(w http.ResponseWriter, r *http.Request) {
	var workout store.Workout
	err := json.NewDecoder(r.Body).Decode(&workout)
	if err != nil {
		fmt.Println("Error decoding request body:", err)
		http.Error(w, "Invalid request body", http.StatusInternalServerError)
		return
	}
	createdWorkout, err := wh.workoutStore.CreateWorkout(&workout)
	if err != nil {
		fmt.Println("Error creating workout:", err)
		http.Error(w, "Failed to create workout", http.StatusInternalServerError)
		return

	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(createdWorkout)

}
func (wh *WorkoutHandler) HandleUpdateWorkoutById(w http.ResponseWriter, r *http.Request) {
	paramsWorkoutId := chi.URLParam(r, "id")
	if paramsWorkoutId == "" {
		http.NotFound(w, r)
		return
	}
	workoutId, err := strconv.ParseInt(paramsWorkoutId, 10, 64)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	existingWorkout, err := wh.workoutStore.GetWorkoutById(workoutId)
	if err != nil {
		fmt.Println("Error getting workout by ID:", err)
		http.Error(w, "Failed to get workout", http.StatusNotFound)
		return
	}
	if existingWorkout == nil {
		http.NotFound(w, r)
		return
	}
	//we assume we already have the workout

	var UpdatedWorkout struct {
		Title            *string              `json:"title"`
		Description      *string              `json:"description"`
		Duration_minutes *int                 `json:"duration_minutes"`
		CaloriesBurned   *int                 `json:"calories_burned"`
		Entries          []store.WorkoutEntry `json:"entries"`
	}
	err = json.NewDecoder(r.Body).Decode(&UpdatedWorkout)
	if err != nil {
		fmt.Println("Error decoding request body:", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	if UpdatedWorkout.Title != nil {
		existingWorkout.Title = *UpdatedWorkout.Title
	}
	if UpdatedWorkout.Description != nil {
		existingWorkout.Description = *UpdatedWorkout.Description
	}
	if UpdatedWorkout.Duration_minutes != nil {
		existingWorkout.Duration_minutes = *UpdatedWorkout.Duration_minutes
	}
	if UpdatedWorkout.CaloriesBurned != nil {
		existingWorkout.CaloriesBurned = *UpdatedWorkout.CaloriesBurned
	}
	if UpdatedWorkout.Entries != nil {
		existingWorkout.Entries = UpdatedWorkout.Entries
	}
	existingWorkout.Entries = UpdatedWorkout.Entries
	err = wh.workoutStore.UpdateWorkout(existingWorkout)
	if err != nil {
		fmt.Println("Error updating workout:", err)
		http.Error(w, "Failed to update workout", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(existingWorkout)

}

func (wh *WorkoutHandler) HandlerDeleteWorkoutById(w http.ResponseWriter, r *http.Request) {
	paramsWorkoutId := chi.URLParam(r, "id")
	if paramsWorkoutId == "" {
		http.NotFound(w, r)
		return
	}
	workoutId, err := strconv.ParseInt(paramsWorkoutId, 10, 64)
	if err == sql.ErrNoRows {
		http.Error(w, "Workout not found", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, "Invalid workout ID", http.StatusInternalServerError)
		return
	}
	err = wh.workoutStore.DeleteWorkout(workoutId)
	if err != nil {
		fmt.Println("Error deleting workout by ID:", err)
		http.Error(w, "Failed to delete workout", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)

}
