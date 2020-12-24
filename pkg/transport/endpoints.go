package transport

import (
	"context"
	"time"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"github.com/macchid/learning/userService/pkg/user"
	"github.com/macchid/learning/userService/pkg/utils"
)

type FetchOneRequest struct {
	Email string
}

type FetchOneResponse struct {
	User  user.User `json:"user,omitempty"`
	Error error     `json:"error,omitempty"`
}

type FetchAllResponse struct {
	Users []user.User `json:"users,omitempty"`
	Error error       `json:"error,omitempty"`
}

type CreateRequest struct {
	Email    string `json:"email"`
	Name     string `json:"name"`
	LastName string `json:"lastname"`
}

type CreateResponse struct {
	User  user.User `json:"user,omitempty"`
	Error error     `json:"error,omitempty"`
}

type UpdateRequest struct {
	Email    string `json:"email,omitempty"`
	Name     string `json:"name,omitempty"`
	LastName string `json:"lastname,omitempty"`
}

type UpdateResponse struct {
	User  user.User `json:"user,omitempty"`
	Error error     `json:"error,omitempty"`
}

type DeleteRequest struct {
	Email string
}

type DeleteResponse struct {
	User  user.User `json:"user,omitempty"`
	Error error     `json:"error,omitempty"`
}

func MakeLoggerMiddleware(logger *log.Logger, method string) func(endpoint.Endpoint) endpoint.Endpoint {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (interface{}, error) {
			_, logEnd := utils.LogStart(logger, method)
			defer logEnd(time.Now())

			return next(ctx, request)
		}
	}
}

func MakeFetchOneEndpoint(srv user.UserService) endpoint.Endpoint {
	_, logEnd := utils.LogStart(srv.Logger(), "MakeEndpoint::FetchOne")
	defer logEnd(time.Now())

	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(FetchOneRequest)
		user, err := srv.FetchOne(req.Email)
		return FetchOneResponse{User: user, Error: err}, err
	}

}

func MakeFetchAllEndpoint(srv user.UserService) endpoint.Endpoint {
	_, logEnd := utils.LogStart(srv.Logger(), "MakeEndpoint::FetchAll")
	defer logEnd(time.Now())

	return func(ctx context.Context, request interface{}) (interface{}, error) {
		users, err := srv.FetchAll()
		return FetchAllResponse{Users: users, Error: err}, err
	}
}

func MakeCreateEndpoint(srv user.UserService) endpoint.Endpoint {
	_, logEnd := utils.LogStart(srv.Logger(), "MakeEndpoint::Create")
	defer logEnd(time.Now())

	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(CreateRequest)
		user, err := srv.Create(user.User{Email: req.Email, FirstName: req.Name, LastName: req.LastName})
		return CreateResponse{User: user, Error: err}, err
	}
}

func MakeUpdateEndpoint(srv user.UserService) endpoint.Endpoint {
	_, logEnd := utils.LogStart(srv.Logger(), "MakeEndpoint::Update")
	defer logEnd(time.Now())

	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(UpdateRequest)
		user, err := srv.Update(user.User{Email: req.Email, FirstName: req.Name, LastName: req.LastName})
		return UpdateResponse{User: user, Error: err}, err
	}
}

func MakeDeleteEndpoint(srv user.UserService) endpoint.Endpoint {
	_, logEnd := utils.LogStart(srv.Logger(), "MakeEndpoint::Delete")
	defer logEnd(time.Now())

	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(DeleteRequest)
		user, err := srv.Delete(req.Email)
		return DeleteResponse{User: user, Error: err}, err
	}
}
