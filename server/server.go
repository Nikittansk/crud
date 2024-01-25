package server

import (
	"encoding/json"
	"net/http"

	"github.com/Nikittansk/crud/model"
)

func JSONResponse(w http.ResponseWriter, r model.Response) {
	w.Header().Set("Content-Type", "application/json")

	switch r.StatusCode {
	case http.StatusNotFound:
		r.Data = "Кажется что-то пошло не так! Страница, которую вы запрашиваете, не существует. Возможно она устарела, была удалена, или был введен неверный адрес в адресной строке."
	}

	jsonResp, err := json.Marshal(r)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write(jsonResp)
}