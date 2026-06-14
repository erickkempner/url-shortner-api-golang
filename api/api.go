package api

import (
	"encoding/json"
	"log/slog"
	"math/rand"
	"net/http"
	"net/url"
	"strings"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

type DB map[string]string

type PostBody struct {
	URL string `json:"url"`
}

type Response struct {
	Error string `json:"error,omitempty"`
	Data  any    `json:"data,omitempty"`
}

func NewHandler(db DB) http.Handler {
	r := chi.NewMux()

	r.Use(middleware.RequestID)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Logger)

	r.Post("/api/short", handlePost(db))
	r.Get("/{code}", handleGet(db))

	return r
}

func handlePost(db DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var body PostBody
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			SendJSON(w, Response{Error: "invalid body"}, http.StatusUnprocessableEntity)
			return
		}

		if _, err := url.Parse(body.URL); err != nil {
			SendJSON(w, Response{Error: "invalid url passed"}, http.StatusBadRequest)
		}

		if !strings.HasPrefix(body.URL, "https") {
			SendJSON(w, Response{Error: "only https urls are allowed"}, http.StatusBadRequest)
			return
		}

		code := genCode()
		db[code] = body.URL

		SendJSON(w, Response{Data: code}, http.StatusCreated)
	}
}

func handleGet(db DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		codeStr := chi.URLParam(r, "code")
		if redirectURL, ok := db[codeStr]; ok {
			http.Redirect(w, r, redirectURL, http.StatusMovedPermanently)
		}

	}
}

func SendJSON(w http.ResponseWriter, resp Response, status int) {
	w.Header().Set("Content-Type", "application/json")
	data, err := json.Marshal(resp)
	if err != nil {
		slog.Error("failed to marshal json data", "error", err)
		SendJSON(w, Response{Error: "something went wrong"}, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(status)
	if _, err := w.Write(data); err != nil {
		slog.Error("failed to write response data", "error", err)
		return
	}

}

const characters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ123456789"

func genCode() string {
	const n = 8
	bytes := make([]byte, n)
	for i := range n {
		bytes[i] = characters[rand.Intn(len(characters))]
	}
	return string(bytes)
}
