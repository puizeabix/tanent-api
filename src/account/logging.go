package account

import (
	"context"
	"time"

	"github.com/go-kit/kit/log"
)

type loggingMiddleware struct {
	logger log.Logger
	s      Service
}

func NewLoggingMiddleware(logger log.Logger, s Service) Service {
	return &loggingMiddleware{
		logger: logger,
		s:      s,
	}
}

func (l *loggingMiddleware) CreateAccount(ctx context.Context, acc Account) (id string, err error) {
	defer func(begin time.Time) {
		l.logger.Log(
			"method", "create_account",
			"id", id,
			"took", time.Since(begin),
			"error", err,
		)
	}(time.Now())
	return l.s.CreateAccount(ctx, acc)
}

func (l *loggingMiddleware) UpdateAccount(ctx context.Context, id string, acc Account) (_ *Account, err error) {
	defer func(begin time.Time) {
		l.logger.Log(
			"method", "update_account",
			"id", id,
			"took", time.Since(begin),
			"error", err,
		)
	}(time.Now())
	return l.s.UpdateAccount(ctx, id, acc)
}

func (l *loggingMiddleware) GetAccount(ctx context.Context, id string) (_ *Account, err error) {
	defer func(begin time.Time) {
		l.logger.Log(
			"method", "get_account",
			"id", id,
			"took", time.Since(begin),
			"error", err,
		)
	}(time.Now())
	return l.s.GetAccount(ctx, id)
}

func (l *loggingMiddleware) ListAccounts(ctx context.Context) (accs []Account, err error) {
	defer func(begin time.Time) {
		l.logger.Log(
			"method", "list_account",
			"count", len(accs),
			"took", time.Since(begin),
			"error", err,
		)
	}(time.Now())
	return l.s.ListAccounts(ctx)
}

func (l *loggingMiddleware) ActivateAccount(ctx context.Context, id string) (err error) {
	defer func(begin time.Time) {
		l.logger.Log(
			"method", "activate_account",
			"id", id,
			"took", time.Since(begin),
			"error", err,
		)
	}(time.Now())
	return l.s.ActivateAccount(ctx, id)
}

func (l *loggingMiddleware) DeactivateAccount(ctx context.Context, id string) (err error) {
	defer func(begin time.Time) {
		l.logger.Log(
			"method", "deactivate_account",
			"id", id,
			"took", time.Since(begin),
			"error", err,
		)
	}(time.Now())
	return l.s.DeactivateAccount(ctx, id)
}
