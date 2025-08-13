package service

import (
	"backend/db"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"reflect"
	"strings"
)

type responseAlert struct {
	Error   string `json:"error,omitempty"`
	Warning string `json:"warning,omitempty"`
}

func Start(r *http.Request, models any, callback func(int, any)) {
	switch r.Method {
	case http.MethodDelete:
		del(r, models, callback)
	case http.MethodGet:
		get(r, models, callback)
	case http.MethodPatch:
		patch(r, models, callback)
	case http.MethodPost:
		post(r, models, callback)
	default:
		callback(http.StatusMethodNotAllowed, responseAlert{Warning: "HTTP method not allowed"})
	}
}

func del(r *http.Request, models any, callback func(int, any)) {
	rvms := reflect.ValueOf(models).Elem()
	rvm := reflect.New(rvms.Field(0).Type().Elem())
	middlePath := rvm.Elem().Type().Name()
	query := strings.Trim(r.URL.Path[strings.Index(r.URL.Path, middlePath)+len(middlePath):], "/")
	data, err := base64.RawURLEncoding.DecodeString(query)
	if err != nil {
		callback(http.StatusBadRequest, responseAlert{Warning: "invalid content"})
		return
	}
	err = json.Unmarshal(data, rvm.Interface())
	if err != nil {
		callback(http.StatusBadRequest, responseAlert{Warning: "invalid content"})
		return
	}
	err = db.Delete(rvms, rvm.Elem())
	if err != nil {
		callback(http.StatusInternalServerError, responseAlert{Error: "could not delete content"})
		return
	}
	callback(http.StatusNoContent, nil)
}

func get(r *http.Request, models any, callback func(int, any)) {
	rvms := reflect.ValueOf(models).Elem()
	rvm := reflect.New(rvms.Field(0).Type().Elem())
	middlePath := rvm.Elem().Type().Name()
	query := strings.Trim(r.URL.Path[strings.Index(r.URL.Path, middlePath)+len(middlePath):], "/")
	if len(query) > 0 {
		data, err := base64.RawURLEncoding.DecodeString(query)
		if err != nil {
			callback(http.StatusBadRequest, responseAlert{Warning: "invalid content"})
			return
		}
		err = json.Unmarshal(data, rvm.Interface())
		if err != nil {
			callback(http.StatusBadRequest, responseAlert{Warning: "invalid content"})
			return
		}
	}
	err := db.Select(rvms, rvm.Elem())
	if err != nil {
		callback(http.StatusInternalServerError, responseAlert{Error: "could not select content"})
		println(err.Error())
		return
	}
	if rvms.Field(0).Len() == 1 {
		callback(http.StatusOK, rvms.Field(0).Index(0).Interface())
		return
	}
	if rvm.Elem().IsZero() || rvms.Field(0).Len() > 1 {
		callback(http.StatusOK, rvms.Interface())
		return
	}
	callback(http.StatusNotFound, responseAlert{Warning: "content not found"})
}

func patch(r *http.Request, models any, callback func(int, any)) {
	rvms := reflect.ValueOf(models).Elem()
	rvm := reflect.New(rvms.Field(0).Type().Elem())
	middlePath := rvm.Elem().Type().Name()
	query := strings.Trim(r.URL.Path[strings.Index(r.URL.Path, middlePath)+len(middlePath):], "/")
	data, err := base64.RawURLEncoding.DecodeString(query)
	if err != nil {
		callback(http.StatusBadRequest, responseAlert{Warning: "invalid content"})
		return
	}
	err = json.Unmarshal(data, rvm.Interface())
	if err != nil {
		callback(http.StatusBadRequest, responseAlert{Warning: "invalid content"})
		return
	}
	err = json.NewDecoder(r.Body).Decode(rvm.Interface())
	if err != nil {
		callback(http.StatusBadRequest, responseAlert{Warning: "invalid content"})
		return
	}
	err = db.Update(rvms, rvm.Elem())
	if err != nil {
		callback(http.StatusInternalServerError, responseAlert{Error: "could not update content"})
		return
	}
	callback(http.StatusNoContent, nil)
}

func post(r *http.Request, models any, callback func(int, any)) {
	rvms := reflect.ValueOf(models).Elem()
	rvm := reflect.New(rvms.Field(0).Type().Elem())
	err := json.NewDecoder(r.Body).Decode(rvm.Interface())
	if err != nil {
		callback(http.StatusBadRequest, responseAlert{Warning: "invalid JSON format or invalid value field type"})
		return
	}
	for f := range rvm.Elem().NumField() {
		if rvm.Elem().Type().Field(f).Tag.Get("db") == "pk" && rvm.Elem().Field(f).IsZero() {
			callback(http.StatusBadRequest, responseAlert{Warning: "empty required field"})
			return
		}
	}
	err = db.Insert(rvms, rvm.Elem())
	if err != nil {
		callback(http.StatusInternalServerError, responseAlert{Error: "could not insert content"})
		return
	}
	callback(http.StatusCreated, rvm.Interface())
}
