package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/Nikittansk/crud/controllers"
	"github.com/Nikittansk/crud/db"
	"github.com/Nikittansk/crud/router"
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
	
	log.Fatal(http.ListenAndServe(":8080", router.Init(userController)))
}