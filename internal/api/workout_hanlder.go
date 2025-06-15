package api

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	"github.com/LuisDavid01/femProject/internal/store"
	"github.com/LuisDavid01/femProject/internal/utils"
)

type WorkoutHandler struct {
	workoutStore store.WorkoutStore
	logger       *log.Logger
}

func NewWorkoutHandler(workoutStore store.WorkoutStore, logger *log.Logger) *WorkoutHandler {
	return &WorkoutHandler{
		workoutStore: workoutStore,
		logger:       logger,
	}
}

func (wh *WorkoutHandler) HandlerWorkoutById(w http.ResponseWriter, r *http.Request) {
	workoutId, err := utils.ReadIDParam(r)
	if err != nil {
		wh.logger.Printf("Error reading workout ID: %v", err)
		utils.WriteJson(w, http.StatusBadRequest, utils.Envelope{"error": "Invalid workout ID"})
		return
	}
	workout, err := wh.workoutStore.GetWorkoutById(workoutId)
	if err != nil {

		wh.logger.Printf("Error getting the workout: %v", err)
		utils.WriteJson(w, http.StatusNotFound, utils.Envelope{"error": "Internal server Error ID"})
		return
	}
	utils.WriteJson(w, http.StatusOK, utils.Envelope{"workout": workout})
}

func (wh *WorkoutHandler) HandlerCreateWorkout(w http.ResponseWriter, r *http.Request) {
	var workout store.Workout
	err := json.NewDecoder(r.Body).Decode(&workout)
	if err != nil {

		wh.logger.Printf("Error decoding the body: %v", err)
		utils.WriteJson(w, http.StatusBadRequest, utils.Envelope{"error": "Invalid request"})
		return
	}
	createdWorkout, err := wh.workoutStore.CreateWorkout(&workout)
	if err != nil {

		wh.logger.Printf("Error creating the workout: %v", err)
		utils.WriteJson(w, http.StatusInternalServerError, utils.Envelope{"error": "Error creating the workout"})
		return

	}
	utils.WriteJson(w, http.StatusCreated, utils.Envelope{"workout": createdWorkout})
}
func (wh *WorkoutHandler) HandleUpdateWorkoutById(w http.ResponseWriter, r *http.Request) {

	workoutId, err := utils.ReadIDParam(r)
	if err != nil {
		wh.logger.Printf("Error reading workout ID: %v", err)
		utils.WriteJson(w, http.StatusBadRequest, utils.Envelope{"error": "Invalid workout ID"})
		return
	}
	existingWorkout, err := wh.workoutStore.GetWorkoutById(workoutId)
	if err != nil {

		wh.logger.Printf("Error getting the workout by ID: %v", err)
		utils.WriteJson(w, http.StatusNotFound, utils.Envelope{"error": "Not found"})
		return
	}
	if existingWorkout == nil {

		wh.logger.Printf("Error getting the workout by ID: %v", err)
		utils.WriteJson(w, http.StatusNotFound, utils.Envelope{"error": "Not found"})
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

		wh.logger.Printf("Error decoding the body: %v", err)
		utils.WriteJson(w, http.StatusBadRequest, utils.Envelope{"error": "Invalid request"})
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

		wh.logger.Printf("Error updating the workout: %v", err)
		utils.WriteJson(w, http.StatusInternalServerError, utils.Envelope{"error": "workout could not be updated"})
		return
	}
	utils.WriteJson(w, http.StatusOK, utils.Envelope{"workout": existingWorkout})
}

func (wh *WorkoutHandler) HandlerDeleteWorkoutById(w http.ResponseWriter, r *http.Request) {

	workoutId, err := utils.ReadIDParam(r)
	if err == sql.ErrNoRows {

		wh.logger.Printf("Workout not foung: %v", err)
		utils.WriteJson(w, http.StatusNotFound, utils.Envelope{"error": "workout could not be found"})
		return
	}
	if err != nil {

		wh.logger.Printf("Error not valid id: %v", err)
		utils.WriteJson(w, http.StatusInternalServerError, utils.Envelope{"error": "workout could not a valid id"})
		return
	}
	err = wh.workoutStore.DeleteWorkout(workoutId)
	if err != nil {

		wh.logger.Printf("Error Deleting the workout: %v", err)
		utils.WriteJson(w, http.StatusInternalServerError, utils.Envelope{"error": "workout could not be Deleted"})
		return
	}
	utils.WriteJson(w, http.StatusNoContent, utils.Envelope{"message": "Workout deleted successfully"})

}
