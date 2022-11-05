package db

import "service-todo-restapi/model"

var Users = map[string]string{}
var Task = map[string][]model.Todo{}

var Sessions = map[string]model.Session{}
