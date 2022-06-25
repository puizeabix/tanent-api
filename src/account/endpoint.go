package account

import (
	"context"
	"time"

	"github.com/go-kit/kit/endpoint"
)

type Endpoints struct {
	CreateAccountEndpoint   endpoint.Endpoint
	GetAccountEndpoint      endpoint.Endpoint
	UpdateAccountEndpoint   endpoint.Endpoint
	ListAccountsEndpoint    endpoint.Endpoint
	ActivateAccountEndpoint endpoint.Endpoint
	DeactiveAccountEndpoint endpoint.Endpoint
}

func MakeEndpoints(s Service) Endpoints {
	return Endpoints{
		CreateAccountEndpoint:   makeCreateAccountEndpoint(s),
		GetAccountEndpoint:      makeGetAccountEndpoint(s),
		UpdateAccountEndpoint:   makeUpdateAccountEndpoint(s),
		ListAccountsEndpoint:    makeListAccountsEndpoint(s),
		ActivateAccountEndpoint: makeActiveAccountEndpoint(s),
		DeactiveAccountEndpoint: makeDeactivateAccountEndpoint(s),
	}

}

type createAccountRequest struct {
	Name string
}

type createAccountResponse struct {
	Id  string `json:"id,omitempty"`
	Err error  `json:"err,omitempty"`
}

func (r createAccountResponse) error() error      { return r.Err }
func (r createAccountResponse) data() interface{} { return r }

func makeCreateAccountEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(createAccountRequest)
		acc := Account{
			Name:     req.Name,
			IsActive: true,
			Created:  time.Now(),
			Modified: time.Now(),
		}
		id, err := s.CreateAccount(ctx, acc)
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

func (r getAccountResponse) error() error      { return r.Err }
func (r getAccountResponse) data() interface{} { return r.Account }

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

func (r updateAccountResponse) error() error      { return r.Err }
func (r updateAccountResponse) data() interface{} { return r.Account }

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

func (r listAccountsResponse) error() error      { return r.Err }
func (r listAccountsResponse) data() interface{} { return r.Accounts }

func makeListAccountsEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, _ interface{}) (interface{}, error) {
		accs, err := s.ListAccounts(ctx)
		return listAccountsResponse{Accounts: accs, Err: err}, err
	}
}

type activateAccountRequest struct {
	Id string
}

type activateAccountResponse struct {
	Err error `json:"err,omitempty"`
}

func (r activateAccountResponse) error() error      { return r.Err }
func (r activateAccountResponse) data() interface{} { return nil }

func makeActiveAccountEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(activateAccountRequest)
		err := s.ActivateAccount(ctx, req.Id)
		return activateAccountResponse{Err: err}, err
	}
}

type deactivateAccountRequest struct {
	Id string
}

type deactivateAccountResponse struct {
	Err error `json:"err,omitempty"`
}

func (r deactivateAccountResponse) error() error      { return r.Err }
func (r deactivateAccountResponse) data() interface{} { return nil }

func makeDeactivateAccountEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(deactivateAccountRequest)
		err := s.DeactivateAccount(ctx, req.Id)
		return deactivateAccountResponse{Err: err}, err
	}
}
