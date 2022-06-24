package account

import (
	"context"

	"github.com/go-kit/kit/endpoint"
)

type Endpoints struct {
	CreateAccountEndpoint endpoint.Endpoint
	GetAccountEndpoint    endpoint.Endpoint
	UpdateAccountEndpoint endpoint.Endpoint
	ListAccountsEndpoint  endpoint.Endpoint
}

func MakeEndpoints(s Service) Endpoints {
	return Endpoints{
		CreateAccountEndpoint: makeCreateAccountEndpoint(s),
		GetAccountEndpoint:    makeGetAccountEndpoint(s),
		UpdateAccountEndpoint: makeUpdateAccountEndpoint(s),
		ListAccountsEndpoint:  makeListAccountsEndpoint(s),
	}

}

type createAccountRequest struct {
	Account
}

type createAccountResponse struct {
	Id  string `json:"id,omitempty"`
	Err error  `json:"err,omitempty"`
}

func (r createAccountResponse) error() error { return r.Err }

func makeCreateAccountEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(createAccountRequest)
		id, err := s.CreateAccount(ctx, req.Account)
		return createAccountResponse{Id: id, Err: err}, err
	}
}

type getAccountRequest struct {
	Id string
}

type getAccountResponse struct {
	Account `json:"account,omitempty"`
	Err     error `json:"err,omitempty"`
}

func (r getAccountResponse) error() error { return r.Err }

func makeGetAccountEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(getAccountRequest)
		acc, err := s.GetAccount(ctx, req.Id)
		return getAccountResponse{Account: *acc, Err: err}, err
	}
}

type updateAccountRequest struct {
	Id      string
	Updated Account
}

type updateAccountResponse struct {
	Account `json:"account,omitempty"`
	Err     error `json:"err,omitempty"`
}

func (r updateAccountResponse) error() error { return r.Err }

func makeUpdateAccountEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(updateAccountRequest)
		u, err := s.UpdateAccount(ctx, req.Id, req.Updated)
		return updateAccountResponse{Account: *u, Err: err}, err
	}
}

type listAccountsRequest struct {
}

type listAccountsResponse struct {
	Accounts []Account `json:"accounts,omitempty"`
	Err      error     `json:"err,omitempty"`
}

func (r listAccountsResponse) error() error { return r.Err }

func makeListAccountsEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, _ interface{}) (interface{}, error) {
		accs, err := s.ListAccounts(ctx)
		return listAccountsResponse{Accounts: accs, Err: err}, err
	}
}
