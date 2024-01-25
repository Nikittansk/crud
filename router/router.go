package router

import (
	"github.com/Nikittansk/crud/controllers"
	"github.com/julienschmidt/httprouter"
)

func Init(us *controllers.UserControllers) *httprouter.Router {
	router := httprouter.New()

	// USERS
	router.GET(routerGroup("/users"), us.GetUsers)
	router.GET(routerGroup("/users/:id"), us.GetUserById)
	router.DELETE(routerGroup("/users/:id"), us.DeleteUserById)
	router.POST(routerGroup("/users"), us.CreateUser)
	router.PUT(routerGroup("/users/:id"), us.UpdateUserById)

	return router
}

func routerGroup(value string) string {
	return("/api" + value)
}