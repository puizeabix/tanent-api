package account

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/transport"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
)

func MakeHTTPServerHandler(s Service, logger log.Logger) http.Handler {
	r := mux.NewRouter()
	e := MakeEndpoints(s)
	options := []httptransport.ServerOption{
		httptransport.ServerErrorHandler(transport.NewLogErrorHandler(logger)),
		httptransport.ServerErrorEncoder(encodeError),
	}

	r.Methods("GET").Path("/tanent/v1/accounts").Handler(httptransport.NewServer(
		e.ListAccountsEndpoint,
		decodeListAccountsRequest,
		encodeResponse,
		options...,
	))

	r.Methods("POST").Path("/tanent/v1/accounts").Handler(httptransport.NewServer(
		e.CreateAccountEndpoint,
		decodeCreateAccountRequest,
		encodeResponse,
		options...,
	))

	r.Methods("GET").Path("/tanent/v1/accounts/{id}").Handler(httptransport.NewServer(
		e.GetAccountEndpoint,
		decodeGetAccountRequest,
		encodeResponse,
		options...,
	))

	r.Methods("PUT").Path("/tanent/v1/accounts/{id}").Handler(httptransport.NewServer(
		e.UpdateAccountEndpoint,
		decodeUpdateAccountRequest,
		encodeResponse,
		options...,
	))

	r.Methods("PUT").Path("/tanent/v1/accounts/{id}/activate").Handler(httptransport.NewServer(
		e.ActivateAccountEndpoint,
		decodeDeactivateAccountRequest,
		encodeResponse,
		options...,
	))

	r.Methods("PUT").Path("/tanent/v1/accounts/{id}/deactivate").Handler(httptransport.NewServer(
		e.DeactiveAccountEndpoint,
		decodeDeactivateAccountRequest,
		encodeResponse,
		options...,
	))

	return r
}

var (
	ErrBadRouting = errors.New("Bad routing")
)

func decodeCreateAccountRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req createAccountRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, err
	}
	return req, nil
}

func decodeGetAccountRequest(_ context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		return nil, ErrBadRouting
	}

	return getAccountRequest{Id: id}, nil
}

func decodeUpdateAccountRequest(_ context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		return nil, ErrBadRouting
	}

	var req updateAccountRequest
	if err := json.NewDecoder(r.Body).Decode(&req.Updated); err != nil {
		return nil, err
	}

	req.Id = id
	return req, nil
}

func decodeActivateAccountRequest(_ context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		return nil, ErrBadRouting
	}
	return activateAccountRequest{Id: id}, nil
}

func decodeDeactivateAccountRequest(_ context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		return nil, ErrBadRouting
	}
	return deactivateAccountRequest{Id: id}, nil
}

func decodeListAccountsRequest(_ context.Context, r *http.Request) (interface{}, error) {
	return listAccountsRequest{}, nil
}

type errorer interface {
	error() error
}

func encodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	if e, ok := response.(errorer); ok && e.error() != nil {
		encodeError(ctx, e.error(), w)
		return nil
	}
	w.Header().Set("Content-type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}

func encodeError(_ context.Context, err error, w http.ResponseWriter) {
	if err == nil {
		panic("encodeError with nil error")
	}
	w.Header().Set("Content-type", "application/json; charset=utf-8")
	w.WriteHeader(codeFrom(err))
	json.NewEncoder(w).Encode(map[string]interface{}{
		"error": err.Error(),
	})
}

func codeFrom(err error) int {
	switch err {
	case ErrAccountNotFound:
		return http.StatusNotFound
	default:
		return http.StatusInternalServerError
	}
}
