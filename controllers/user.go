package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/Nikittansk/crud/model"
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
		w.WriteHeader(http.StatusNotFound)
		return
	}

	for cur.Next(us.ctx) {
		u := model.User{}
		if err := cur.Decode(&u); err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		users = append(users, &u)
	}

	if err := cur.Err(); err != nil {
		log.Fatal(err)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	cur.Close(us.ctx)

	if len(users) == 0 {
		w.WriteHeader(http.StatusNotFound)
		return	
	}

	uj, err := json.Marshal(users)
	if err != nil{
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "%s\n", uj)
}

func (us *UserControllers) GetUserById(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	objectId, err :=  primitive.ObjectIDFromHex(p.ByName("id"))

	if err != nil{
		w.WriteHeader(http.StatusNotFound)
		return
	}

	u := model.User{}

	if err := us.col().FindOne(us.ctx, bson.M{"_id": objectId}).Decode(&u); err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	uj, err := json.Marshal(u)
	if err != nil{
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "%s\n", uj)
}

func (us *UserControllers) CreateUser(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	u := model.User{}

	json.NewDecoder(r.Body).Decode(&u)

	u.Id = primitive.NewObjectID()

	res, err := us.col().InsertOne(us.ctx, u)

	uj, err := json.Marshal(res)
	if err != nil{
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "%s\n", uj)
}

func (us *UserControllers) UpdateUserById(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	objectId, err :=  primitive.ObjectIDFromHex(p.ByName("id"))
	if err != nil{
		w.WriteHeader(http.StatusNotFound)
		return
	}
	u := model.User{}

	json.NewDecoder(r.Body).Decode(&u)

	res, err := us.col().UpdateOne(us.ctx, bson.M{"_id": objectId}, bson.M{"$set": bson.M{
		"name": u.Name,
		"gender": u.Gender,
		"age": u.Age,
	}})

	uj, err := json.Marshal(res)
	if err != nil{
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "%s\n", uj)
}

func (us *UserControllers) DeleteUserById(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	objectId, err :=  primitive.ObjectIDFromHex(p.ByName("id"))

	if err != nil{
		w.WriteHeader(http.StatusNotFound)
		return
	}

	res, err := us.col().DeleteOne(us.ctx, bson.M{"_id": objectId})
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	if res.DeletedCount == 0 {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintln(w, "Документ с таким id не существует!")
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Документ с %s удален успешно!", objectId)
}