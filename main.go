package main

import (
	"net/http"
	"time"
	"fmt"
	"flag"
	"github.com/LuisDavid01/femProject/internal/app"
	"github.com/LuisDavid01/femProject/internal/routes"
)
func main() {
	var port int
	flag.IntVar(&port, "port", 8080, "go backend server port")
	flag.Parse()
	app, err := app.NewApplication()
	if err != nil {
		panic(err)
	}
	defer app.DB.Close()
	r := routes.SetupRoutes(app)
	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: r,
		IdleTimeout: time.Minute,
		ReadTimeout: 10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}
	app.Logger.Printf("Starting application... on port: %d", port)
	err = server.ListenAndServe()
	if err != nil {
		app.Logger.Fatal(err)
	}


}




