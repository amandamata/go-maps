package main

import (
	"fmt"
	"go-maps/src"
	"go-maps/src/db"
	"go-maps/src/router"

	"net/http"
)

func init() {
	src.LoadEnvs()
}

func main() {
	done := make(chan bool)
	go func() {
		r := router.Generate()
		fmt.Println("api running, listening on port 8080")
		if err := http.ListenAndServe(":8080", r); err != nil {
			fmt.Println("Failed to start server:", err)
			done <- true
		}
	}()
	<-done

	db.Disconnect()
}
