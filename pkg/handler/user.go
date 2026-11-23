package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"microservices/gocourse_user/internal/user"
	"net/http"
	"strconv"

	"github.com/go-kit/kit/endpoint"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"github.com/michael870917/go_lib_response/response"
)

func NewUserHTTPServer(ctx context.Context, s user.ServiceInterface, config user.Config) *mux.Router {
	endpoints := user.MakeEndpoints(s, config)
	router := mux.NewRouter()

	opts := []httptransport.ServerOption{
		httptransport.ServerErrorEncoder(encodeError),
	}
	router.Handle("/users", httptransport.NewServer(
		endpoint.Endpoint(endpoints.Create),
		decodeCreateUser, encodeResponse,
		opts...,
	)).Methods("POST")
	router.Handle("/users", httptransport.NewServer(
		endpoint.Endpoint(endpoints.GetAll),
		decodeGetAllUser, encodeResponse,
		opts...,
	)).Methods("GET")
	router.Handle("/users/{id}", httptransport.NewServer(
		endpoint.Endpoint(endpoints.Get),
		decodeGetUser, encodeResponse,
		opts...,
	)).Methods("GET")
	router.Handle("/users/{id}", httptransport.NewServer(
		endpoint.Endpoint(endpoints.Update),
		decodeUpdateUser, encodeResponse,
		opts...,
	)).Methods("PATCH")
	router.Handle("/users/{id}", httptransport.NewServer(
		endpoint.Endpoint(endpoints.Delete),
		decodeDeleteUser, encodeResponse,
		opts...,
	)).Methods("DELETE")

	return router
}

func decodeCreateUser(_ context.Context, r *http.Request) (interface{}, error) {
	var req user.CreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, response.BadRequest(fmt.Sprintf("invalid request format: '%v'", err.Error()))
	}

	return req, nil
}

func decodeGetUser(_ context.Context, r *http.Request) (interface{}, error) {
	params := mux.Vars(r)
	var req user.GetRequest
	req.ID = params["id"]

	return req, nil
}

func decodeGetAllUser(_ context.Context, r *http.Request) (interface{}, error) {
	v := r.URL.Query()
	limit, _ := strconv.Atoi(v.Get("limit"))
	page, _ := strconv.Atoi(v.Get("page"))

	var req user.GetAllRequest
	req.FirstName = v.Get("first_name")
	req.LastName = v.Get("last_name")
	req.Page = page
	req.Limit = limit

	return req, nil
}

func decodeUpdateUser(_ context.Context, r *http.Request) (interface{}, error) {
	var req user.UpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, response.BadRequest(fmt.Sprintf("invalid request format: '%v'", err.Error()))
	}

	params := mux.Vars(r)
	req.ID = params["id"]

	return req, nil
}

func decodeDeleteUser(_ context.Context, r *http.Request) (interface{}, error) {
	params := mux.Vars(r)
	var req user.DeleteRequest
	req.ID = params["id"]

	return req, nil
}

func encodeError(_ context.Context, err error, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	r := err.(response.Response)
	w.WriteHeader(r.StatusCode())
	_ = json.NewEncoder(w).Encode(r)
}

func encodeResponse(_ context.Context, w http.ResponseWriter, resp interface{}) error {
	r := resp.(response.Response)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(r.StatusCode())

	return json.NewEncoder(w).Encode(r)
}
