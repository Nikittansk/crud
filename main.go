package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/Nikittansk/crud/controllers"
	"github.com/Nikittansk/crud/db"
	"github.com/julienschmidt/httprouter"
)

func main() {
	dbCtx, dbCancel := context.WithTimeout(context.Background(), 1 * time.Second)
	defer dbCancel()

	mgoClient, err := db.ConnectToMongoDB(dbCtx, "mongodb://localhost:27017")
	if err != nil {
		log.Fatal(err)
	}
	mgoDB := mgoClient.Client().Database("crud")
	
	userController := controllers.NewUserComtroller(mgoDB, context.TODO())
	
	router := httprouter.New()

    router.GET("/users", userController.GetUsers)
	router.GET("/users/:id", userController.GetUserById)
	router.DELETE("/users/:id", userController.DeleteUserById)
	router.POST("/user", userController.CreateUser)
	router.PUT("/users/:id", userController.UpdateUserById)
	
	log.Fatal(http.ListenAndServe(":8080", router))
}