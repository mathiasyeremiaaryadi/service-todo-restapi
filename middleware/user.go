package middleware

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"service-todo-restapi/db"
	"service-todo-restapi/model"
)

func isExpired(s model.Session) bool {
	return s.Expiry.Before(time.Now())
}

func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c, err := r.Cookie("session_token")
		if err != nil {
			resp := model.ErrorResponse{Error: "http: named cookie not present"}

			jsonResp, _ := json.Marshal(resp)

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			w.Write(jsonResp)

			return
		}

		type UsernameContext string
		const usernameContext UsernameContext = "username"

		ctx := context.WithValue(r.Context(), usernameContext, db.Sessions[c.Value].Username)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
