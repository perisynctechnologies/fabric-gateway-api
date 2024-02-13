package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/Subskribo-BV/dnn-fabric-api/service"
)

type IHandler interface {
	HandleAddAsset(w http.ResponseWriter, r *http.Request)
	HandleVoidAsset(w http.ResponseWriter, r *http.Request)
	HandleReleaseAsset(w http.ResponseWriter, r *http.Request)
	HandleExpireAsset(w http.ResponseWriter, r *http.Request)
	HandleGetAsset(w http.ResponseWriter, r *http.Request)
	HandleGetAllAssets(w http.ResponseWriter, r *http.Request)
}

type RequestBody struct {
	Data string `json:"data"`
}

type Map map[string]any

type Handler struct {
	s service.IService
}

func New(s service.IService) IHandler {
	return &Handler{
		s: s,
	}
}

func (h *Handler) HandleAddAsset(w http.ResponseWriter, r *http.Request) {
	body := new(RequestBody)
	if err := json.NewDecoder(r.Body).Decode(body); err != nil {
		writeJson(w, http.StatusBadRequest, Map{"code": 1, "err": "malformed request"})
		return
	}

	res, err := h.s.CreateAsset(body.Data)
	if err != nil {
		writeJson(w, http.StatusBadRequest, Map{"code": 2, "err": "unable to create contract"})
		return
	}

	writeJson(w, http.StatusOK, res)
}

func (h *Handler) HandleVoidAsset(w http.ResponseWriter, r *http.Request) {
	body := new(RequestBody)
	if err := json.NewDecoder(r.Body).Decode(body); err != nil {
		writeJson(w, http.StatusBadRequest, Map{"code": 1, "err": "malformed request"})
		return
	}

	res, err := h.s.VoidAsset(body.Data)
	if err != nil {
		writeJson(w, http.StatusBadRequest, Map{"code": 2, "err": "unable to void contract"})
		return
	}

	writeJson(w, http.StatusOK, res)
}

func (h *Handler) HandleReleaseAsset(w http.ResponseWriter, r *http.Request) {
	body := new(RequestBody)
	if err := json.NewDecoder(r.Body).Decode(body); err != nil {
		writeJson(w, http.StatusBadRequest, Map{"code": 1, "err": "malformed request"})
		return
	}

	res, err := h.s.ReleaseAsset(body.Data)
	if err != nil {
		writeJson(w, http.StatusBadRequest, Map{"code": 2, "err": "unable to release contract"})
		return
	}

	writeJson(w, http.StatusOK, res)
}

func (h *Handler) HandleExpireAsset(w http.ResponseWriter, r *http.Request) {
	body := new(RequestBody)
	if err := json.NewDecoder(r.Body).Decode(body); err != nil {
		writeJson(w, http.StatusBadRequest, Map{"code": 1, "err": "malformed request"})
		return
	}

	res, err := h.s.ExpireAsset(body.Data)
	if err != nil {
		writeJson(w, http.StatusBadRequest, Map{"code": 2, "err": "unable to expire contract"})
		return
	}

	writeJson(w, http.StatusOK, res)
}

func (h *Handler) HandleGetAsset(w http.ResponseWriter, r *http.Request) {
	body := new(RequestBody)
	if err := json.NewDecoder(r.Body).Decode(body); err != nil {
		writeJson(w, http.StatusBadRequest, Map{"code": 1, "err": "malformed request"})
		return
	}

	id := r.URL.Query().Get("id")
	if id == "" {
		writeJson(w, http.StatusBadRequest, Map{"code": 3, "err": "id is empty"})
		return
	}

	res, err := h.s.ReadAssetByID(id)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(Map{"code": 2, "err": fmt.Sprintf("unable to read contract with id: %s", id)})
		return
	}

	writeJson(w, http.StatusOK, res)
}

func (h *Handler) HandleGetAllAssets(w http.ResponseWriter, r *http.Request) {
	body := new(RequestBody)
	if err := json.NewDecoder(r.Body).Decode(body); err != nil {
		writeJson(w, http.StatusBadRequest, Map{"code": 1, "err": "malformed request"})
		return
	}

	res, err := h.s.GetAllAssets()
	if err != nil {
		writeJson(w, http.StatusBadRequest, Map{"code": 2, "err": "unable to get all contracts"})
		return
	}

	writeJson(w, http.StatusOK, res)
}

func writeJson(w http.ResponseWriter, code int, data any) error {
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(code)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		return err
	}

	return nil
}
