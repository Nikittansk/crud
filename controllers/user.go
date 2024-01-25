package controllers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/Nikittansk/crud/model"
	"github.com/Nikittansk/crud/server"
	"github.com/julienschmidt/httprouter"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserControllers struct {
	db *mongo.Database
	ctx context.Context
}

func NewUserComtroller(db *mongo.Database, ctx context.Context) *UserControllers {
	return &UserControllers{
		db: db,
		ctx: ctx,
	}
}

func (us *UserControllers) col() *mongo.Collection {
	return us.db.Collection("users")
}

func (us *UserControllers) GetUsers(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
    var users []*model.User

	cur, err := us.col().Find(us.ctx, bson.M{})
	if err != nil {
		server.JSONResponse(w, model.Response{
			StatusCode: http.StatusNotFound,
		})
		return
	}

	for cur.Next(us.ctx) {
		u := model.User{}
		if err := cur.Decode(&u); err != nil {
			server.JSONResponse(w, model.Response{
				StatusCode: http.StatusNotFound,
			})
			return
		}

		users = append(users, &u)
	}

	if err := cur.Err(); err != nil {
		log.Fatal(err)
		server.JSONResponse(w, model.Response{
			StatusCode: http.StatusNotFound,
		})
		return
	}

	cur.Close(us.ctx)

	if len(users) == 0 {
		server.JSONResponse(w, model.Response{
			StatusCode: http.StatusNotFound,
		})
		return
	}

	server.JSONResponse(w, model.Response{
		StatusCode: http.StatusOK,
		Data: users,
	})
}

func (us *UserControllers) GetUserById(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	objectId, err :=  primitive.ObjectIDFromHex(p.ByName("id"))

	if err != nil{
		server.JSONResponse(w, model.Response{
			StatusCode: http.StatusNotFound,
		})
		return
	}

	u := model.User{}

	if err := us.col().FindOne(us.ctx, bson.M{"_id": objectId}).Decode(&u); err != nil {
		server.JSONResponse(w, model.Response{
			StatusCode: http.StatusNotFound,
		})
		return
	}

	server.JSONResponse(w, model.Response{
		StatusCode: http.StatusOK,
		Data: u,
	})
}

func (us *UserControllers) CreateUser(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	u := model.User{}

	json.NewDecoder(r.Body).Decode(&u)

	u.Id = primitive.NewObjectID()

	res, err := us.col().InsertOne(us.ctx, u)
	if err != nil{
		server.JSONResponse(w, model.Response{
			StatusCode: http.StatusNotFound,
		})
		return
	}

	server.JSONResponse(w, model.Response{
		StatusCode: http.StatusOK,
		Data: res,
	})
}

func (us *UserControllers) UpdateUserById(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	objectId, err :=  primitive.ObjectIDFromHex(p.ByName("id"))
	if err != nil{
		server.JSONResponse(w, model.Response{
			StatusCode: http.StatusNotFound,
		})
		return
	}
	u := model.User{}

	json.NewDecoder(r.Body).Decode(&u)
	if u.Name == "" || u.Gender == "" || u.Age == 0 {
		server.JSONResponse(w, model.Response{
			StatusCode: http.StatusBadRequest,
			Data: "Убедитесь, что все обязательные поля заполнены. Если запрос требует определенных параметров, удостоверьтесь, что они указаны и корректны.",
		})
		return
	}

	res, err := us.col().UpdateOne(us.ctx, bson.M{"_id": objectId}, bson.M{"$set": bson.M{
		"name": u.Name,
		"gender": u.Gender,
		"age": u.Age,
	}})
	if err != nil{
		server.JSONResponse(w, model.Response{
			StatusCode: http.StatusNotFound,
		})
		return
	}

	server.JSONResponse(w, model.Response{
		StatusCode: http.StatusOK,
		Data: res,
	})
}

func (us *UserControllers) DeleteUserById(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	objectId, err :=  primitive.ObjectIDFromHex(p.ByName("id"))

	if err != nil{
		server.JSONResponse(w, model.Response{
			StatusCode: http.StatusNotFound,
		})
		return
	}

	res, err := us.col().DeleteOne(us.ctx, bson.M{"_id": objectId})
	if err != nil {
		server.JSONResponse(w, model.Response{
			StatusCode: http.StatusNotFound,
		})
		return
	}

	if res.DeletedCount == 0 {
		server.JSONResponse(w, model.Response{
			StatusCode: http.StatusNotFound,
		})
		return
	}

	server.JSONResponse(w, model.Response{
		StatusCode: http.StatusOK,
		Data: "Запись была успешно удалена!",
	})
}