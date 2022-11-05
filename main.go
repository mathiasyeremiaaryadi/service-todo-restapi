package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"service-todo-restapi/db"
	"service-todo-restapi/middleware"
	"service-todo-restapi/model"
	"time"

	"github.com/google/uuid"
)

func Register(w http.ResponseWriter, r *http.Request) {
	var credential model.Credentials

	err := json.NewDecoder(r.Body).Decode(&credential)
	if err != nil {
		resp := model.ErrorResponse{Error: "Internal Server Error"}
		jsonResp, _ := json.Marshal(resp)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write(jsonResp)

		return
	}

	if credential.Username == "" || credential.Password == "" {
		resp := model.ErrorResponse{Error: "Username or Password empty"}
		jsonResp, _ := json.Marshal(resp)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write(jsonResp)

		return
	}

	_, ok := db.Users[credential.Username]
	if ok {
		resp := model.ErrorResponse{Error: "Username already exist"}
		jsonResp, _ := json.Marshal(resp)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusConflict)
		w.Write(jsonResp)

		return
	}

	db.Users[credential.Username] = credential.Password
	resp := model.SuccessResponse{
		Username: credential.Username,
		Message:  "Register Success",
	}

	jsonResp, _ := json.Marshal(resp)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResp)

}

func Login(w http.ResponseWriter, r *http.Request) {
	var credential model.Credentials

	err := json.NewDecoder(r.Body).Decode(&credential)
	if err != nil {
		resp := model.ErrorResponse{Error: "Internal Server Error"}
		jsonResp, _ := json.Marshal(resp)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write(jsonResp)

		return
	}

	if credential.Username == "" || credential.Password == "" {
		resp := model.ErrorResponse{Error: "Username or Password empty"}
		jsonResp, _ := json.Marshal(resp)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write(jsonResp)

		return
	}

	v, ok := db.Users[credential.Username]
	if !ok || v != credential.Password {
		resp := model.ErrorResponse{Error: "Wrong User or Password!"}
		jsonResp, _ := json.Marshal(resp)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		w.Write(jsonResp)

		return
	}

	var c *http.Cookie

	c.Name = "session_token"
	c.Value = uuid.NewString()
	c.Expires = time.Now().Add(time.Hour * 5)

	session := model.Session{
		Username: credential.Username,
		Expiry:   c.Expires,
	}

	db.Sessions[c.Value] = session
	resp := model.SuccessResponse{
		Username: credential.Username,
		Message:  "Login Success",
	}

	jsonResp, _ := json.Marshal(resp)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResp)
}

func AddToDo(w http.ResponseWriter, r *http.Request) {
	var todo model.Todo
	username := r.Context().Value("usernameContext").(string)

	err := json.NewDecoder(r.Body).Decode(&todo)
	if err != nil {
		resp := model.ErrorResponse{Error: "Internal Server Error"}
		jsonResp, _ := json.Marshal(resp)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write(jsonResp)

		return
	}

	todo.Id = uuid.NewString()
	db.Task[username] = append(db.Task[username], todo)
	resp := model.SuccessResponse{
		Username: username,
		Message:  fmt.Sprintf("Task %s added!", todo.Task),
	}

	jsonResp, _ := json.Marshal(resp)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResp)
}

func ListToDo(w http.ResponseWriter, r *http.Request) {
	var todo []model.Todo
	username := r.Context().Value("usernameContext").(string)

	err := json.NewDecoder(r.Body).Decode(&todo)
	if err != nil {
		resp := model.ErrorResponse{Error: "Internal Server Error"}
		jsonResp, _ := json.Marshal(resp)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write(jsonResp)

		return
	}

	if len(db.Task[username]) == 0 {
		resp := model.ErrorResponse{Error: "Todolist not found!"}
		jsonResp, _ := json.Marshal(resp)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		w.Write(jsonResp)

		return
	}

	todo = append(todo, db.Task[username]...)
	jsonResp, _ := json.Marshal(todo)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResp)
}

func ClearToDo(w http.ResponseWriter, r *http.Request) {
	username := r.Context().Value("usernameContext").(string)

	resp := model.SuccessResponse{
		Username: username,
		Message:  "Clear ToDo Success",
	}

	db.Task[username] = []model.Todo{}
	jsonResp, _ := json.Marshal(resp)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResp)
}

func Logout(w http.ResponseWriter, r *http.Request) {
	username := r.Context().Value("usernameContext").(string)

	resp := model.SuccessResponse{
		Username: username,
		Message:  "Logout Success",
	}

	db.Sessions = map[string]model.Session{}

	jsonResp, _ := json.Marshal(resp)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResp)
}

func ResetToDo(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

type API struct {
	mux *http.ServeMux
}

func NewAPI() API {
	mux := http.NewServeMux()
	api := API{
		mux,
	}

	mux.Handle("/user/register", middleware.Post(http.HandlerFunc(Register)))
	mux.Handle("/user/login", middleware.Post(http.HandlerFunc(Login)))
	mux.Handle("/user/logout", middleware.Get(middleware.Auth(http.HandlerFunc(Logout))))

	mux.Handle("/todo/create", middleware.Post(middleware.Auth(http.HandlerFunc(AddToDo))))
	mux.Handle("/todo/read", middleware.Post(middleware.Auth(http.HandlerFunc(ListToDo))))
	mux.Handle("/todo/read", middleware.Post(middleware.Auth(http.HandlerFunc(ListToDo))))

	mux.Handle("/todo/reset", http.HandlerFunc(ResetToDo))

	return api
}

func (api *API) Handler() *http.ServeMux {
	return api.mux
}

func (api *API) Start() {
	fmt.Println("starting web server at http://localhost:8080")
	http.ListenAndServe(":8080", api.Handler())
}

func main() {
	mainAPI := NewAPI()
	mainAPI.Start()
}
