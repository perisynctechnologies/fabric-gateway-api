package router

import (
	"net/http"

	"github.com/Subskribo-BV/dnn-fabric-api/api/handler"
	"github.com/Subskribo-BV/dnn-fabric-api/utils/auth"

	"github.com/gorilla/mux"
)

func BuildRouter(h handler.IHandler) *mux.Router {
	r := mux.NewRouter()

	api := r.PathPrefix("/api/v1/cc").Subrouter()

	api.HandleFunc("", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("content-type", "text/text")
		w.Write([]byte("blockchain gateway!"))
	})

	secure := api.PathPrefix("").Subrouter()
	secure.Use(auth.AuthMiddleware)

	secure.HandleFunc("/asset", h.HandleAddAsset).Methods(http.MethodPost)
	secure.HandleFunc("/assets", h.HandleGetAllAssets).Methods(http.MethodPost)
	secure.HandleFunc("/asset", h.HandleGetAsset).Methods(http.MethodGet)
	secure.HandleFunc("/asset/void", h.HandleVoidAsset).Methods(http.MethodPut)
	secure.HandleFunc("/asset/release", h.HandleReleaseAsset).Methods(http.MethodPut)
	secure.HandleFunc("/asset/expire", h.HandleExpireAsset).Methods(http.MethodPut)

	return r
}
