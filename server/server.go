package server

import (
	"backend/model"
	"backend/service"
	"encoding/json"
	"log"
	"net/http"
)

type handler struct{}

const basePath string = "/world/"

func (handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	router(r, func(statusCode int, response any) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(statusCode)
		json.NewEncoder(w).Encode(response)
	})
}

func Start() {
	log.Fatal(http.ListenAndServe(":3000", handler{}))
}

func router(r *http.Request, callback func(int, any)) {
	switch true {
	case len(r.URL.Path) >= len(basePath+model.ContinentPath) && r.URL.Path[:len(basePath+model.ContinentPath)] == basePath+model.ContinentPath:
		service.Start(r, &model.Continents{}, callback)
	case len(r.URL.Path) >= len(basePath+model.CountryPath) && r.URL.Path[:len(basePath+model.CountryPath)] == basePath+model.CountryPath:
		service.Start(r, &model.Countries{}, callback)
	case len(r.URL.Path) >= len(basePath+model.CityPath) && r.URL.Path[:len(basePath+model.CityPath)] == basePath+model.CityPath:
		service.Start(r, &model.Cities{}, callback)
	}
}
