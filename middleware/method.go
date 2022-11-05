package middleware

import (
	"encoding/json"
	"net/http"
	"service-todo-restapi/model"
)

func Get(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			resp := model.ErrorResponse{Error: "Method is not allowed!"}
			jsonResp, _ := json.Marshal(resp)

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusMethodNotAllowed)
			w.Write(jsonResp)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func Post(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			resp := model.ErrorResponse{Error: "Method is not allowed!"}
			jsonResp, _ := json.Marshal(resp)

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusMethodNotAllowed)
			w.Write(jsonResp)
			return
		}

		next.ServeHTTP(w, r)

	})
}

func Delete(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			resp := model.ErrorResponse{Error: "Method is not allowed!"}
			jsonResp, _ := json.Marshal(resp)

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusMethodNotAllowed)
			w.Write(jsonResp)
			return
		}

		next.ServeHTTP(w, r)
	})
}
