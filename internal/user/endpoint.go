package user

import (
	"context"
	"errors"

	"github.com/michael870917/go_lib_response/response"
	"github.com/michael870917/gocourse_meta/meta"
)

type (
	//Controller func(w http.ResponseWriter, r *http.Request)
	Controller func(ctx context.Context, request interface{}) (interface{}, error)
	Endpoints  struct {
		Create Controller
		Get    Controller
		GetAll Controller
		Update Controller
		Delete Controller
	}

	CreateRequest struct {
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		Email     string `json:"email"`
		Phone     string `json:"phone"`
	}

	GetRequest struct {
		ID string
	}

	GetAllRequest struct {
		FirstName string
		LastName  string
		Page      int
		Limit     int
	}

	UpdateRequest struct {
		ID        string
		FirstName *string `json:"first_name"`
		LastName  *string `json:"last_name"`
		Email     *string `json:"email"`
		Phone     *string `json:"phone"`
	}

	DeleteRequest struct {
		ID string
	}

	Config struct {
		LimitPageDefault string
	}
)

func MakeEndpoints(s ServiceInterface, config Config) Endpoints {
	return Endpoints{
		Create: createEndpointHandler(s),
		Get:    getEndpointHandler(s),
		GetAll: getAllEndpointHandler(s, config),
		Update: updateEndpointHandler(s),
		Delete: deleteEndpointHandler(s),
	}
}

func createEndpointHandler(s ServiceInterface) Controller {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(CreateRequest)

		if req.FirstName == "" {
			return nil, response.BadRequest(ErrFirstNameRequired.Error())
		}

		if req.LastName == "" {
			return nil, response.BadRequest(ErrLastNameRequired.Error())
		}

		user, err := s.Create(ctx, req.FirstName, req.LastName, req.Email, req.Phone)

		if err != nil {
			return nil, response.InternalServerError(err.Error())
		}

		return response.Created("success", user, nil), nil
	}
}

func getAllEndpointHandler(s ServiceInterface, config Config) Controller {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(GetAllRequest)

		filters := Filters{
			FirstName: req.FirstName,
			LastName:  req.LastName,
		}

		count, err := s.Count(ctx, filters)
		if err != nil {
			return nil, response.InternalServerError(err.Error())
		}

		metaResponse, err := meta.New(req.Page, req.Limit, int(count), config.LimitPageDefault)
		if err != nil {
			return nil, response.InternalServerError(err.Error())
		}

		users, err := s.GetAll(ctx, filters, metaResponse.Offset(), metaResponse.Limit())
		if err != nil {
			return nil, response.InternalServerError(err.Error())
		}

		return response.OK("success", users, metaResponse), nil
	}
}

func getEndpointHandler(s ServiceInterface) Controller {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(GetRequest)
		user, err := s.GetById(ctx, req.ID)
		if err != nil {
			if errors.As(err, &ErrNotFound{}) {
				return nil, response.NotFound(err.Error())
			}

			return nil, response.InternalServerError(err.Error())
		}
		return response.OK("success", user, nil), nil
	}
}

func updateEndpointHandler(s ServiceInterface) Controller {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(UpdateRequest)

		if req.FirstName != nil && *req.FirstName == "" {
			return nil, response.BadRequest(ErrFirstNameRequired.Error())
		}

		if req.LastName != nil && *req.LastName == "" {
			return nil, response.BadRequest(ErrLastNameRequired.Error())
		}

		err := s.Update(ctx, req.ID, req.FirstName, req.LastName, req.Email, req.Phone)
		if err != nil {
			if errors.As(err, &ErrNotFound{}) {
				return nil, response.NotFound(err.Error())
			}

			return nil, response.InternalServerError(err.Error())
		}
		return response.OK("success", nil, nil), nil
	}
}

func deleteEndpointHandler(s ServiceInterface) Controller {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(DeleteRequest)
		err := s.Delete(ctx, req.ID)
		if err != nil {
			if errors.As(err, &ErrNotFound{}) {
				return nil, response.NotFound(err.Error())
			}

			return nil, response.BadRequest(err.Error())
		}

		return response.OK("success", nil, nil), nil
	}
}
