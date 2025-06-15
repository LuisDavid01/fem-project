package store

import (
	"database/sql"
	"testing"

	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupTestDB(t *testing.T) *sql.DB {
	db, err := sql.Open("pgx", "host=localhost user=postgres password=postgres dbname=postgres port=5433 sslmode=disable")
	if err != nil {
		t.Fatalf("Failed to open test database: %v", err)
	}
	err = Migrate(db, "../../migrations/")
	if err != nil {
		t.Fatalf("Failed to migrate test database: %v", err)
	}
	_, err = db.Exec("TRUNCATE TABLE workouts, workout_entries CASCADE")
	if err != nil {
		t.Fatalf("Failed to truncate test database: %v", err)
	}
	return db

}

func TestCreateWorkout(t *testing.T) {
	db := setupTestDB(t)
	store := NewPostgresWorkoutStore(db)

	tests := []struct {
		name    string
		workout *Workout
		wantErr bool
	}{
		{
			name: "Valid Workout",
			workout: &Workout{
				Title:            "Test Workout",
				Description:      "This is a test workout",
				Duration_minutes: 30,
				CaloriesBurned:   300,
				Entries: []WorkoutEntry{
					{
						ExerciseName: "Push Up",
						Sets:         3,
						Reps:         IntPtr(10),
						Weight:       FloatPtr(10.5),
						OrderIndex:   1,
					},
					{
						ExerciseName: "Squat",
						Sets:         3,
						Reps:         IntPtr(15),
						Weight:       FloatPtr(255.0),
						OrderIndex:   2,
					},
				},
			},
			wantErr: true,
		},
		{
			name: "Invalid Workout (invalid Entries)",
			workout: &Workout{
				Title:            "Test Workout",
				Description:      "This is a test workout",
				Duration_minutes: 30,
				CaloriesBurned:   300,
				Entries: []WorkoutEntry{
					{
						ExerciseName: "Push Up",
						Sets:         3,
						Reps:         IntPtr(50),
						Notes:        "Too many reps",
						OrderIndex:   1,
					},
					{
						ExerciseName:    "Squat",
						Sets:            4,
						Reps:            IntPtr(15),
						DurationSeconds: IntPtr(60),
						OrderIndex:      2,
					},
				},
			},
			wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			createdWorkout, err := store.CreateWorkout(tt.workout)
			if tt.wantErr {

				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.workout.Title, createdWorkout.Title)
			assert.Equal(t, tt.workout.Description, createdWorkout.Description)
			assert.Equal(t, tt.workout.Duration_minutes, createdWorkout.Duration_minutes)

			retrived, err := store.GetWorkoutById(int64(createdWorkout.ID))
			require.NoError(t, err)
			assert.Equal(t, createdWorkout.ID, retrived.ID)
			assert.Equal(t, len(tt.workout.Entries), len(retrived.Entries))
			for i := range retrived.Entries {
				assert.Equal(t, tt.workout.Entries[i].ExerciseName, retrived.Entries[i].ExerciseName)
				assert.Equal(t, tt.workout.Entries[i].Sets, retrived.Entries[i].Sets)
			}
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateWorkout() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
func IntPtr(i int) *int {
	return &i
}
func FloatPtr(f float64) *float64 {
	return &f
}
