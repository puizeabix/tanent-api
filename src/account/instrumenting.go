package account

import (
	"context"
	"time"

	"github.com/go-kit/kit/metrics"
)

type instrumentingMiddleware struct {
	requestCount   metrics.Counter
	requestLatency metrics.Histogram
	s              Service
}

func NewInstrumentingMiddleware(counter metrics.Counter, latency metrics.Histogram, s Service) Service {
	return &instrumentingMiddleware{
		requestCount:   counter,
		requestLatency: latency,
		s:              s,
	}
}

func (mw *instrumentingMiddleware) CreateAccount(ctx context.Context, acc Account) (string, error) {
	defer func(begin time.Time) {
		mw.requestCount.With("method", "create_account").Add(1)
		mw.requestLatency.With("method", "create_account").Observe(time.Since(begin).Seconds())
	}(time.Now())
	return mw.s.CreateAccount(ctx, acc)
}

func (mw *instrumentingMiddleware) UpdateAccount(ctx context.Context, id string, acc Account) (*Account, error) {
	defer func(begin time.Time) {
		mw.requestCount.With("method", "update_account").Add(1)
		mw.requestLatency.With("method", "update_account").Observe(time.Since(begin).Seconds())
	}(time.Now())
	return mw.s.UpdateAccount(ctx, id, acc)
}

func (mw *instrumentingMiddleware) GetAccount(ctx context.Context, id string) (*Account, error) {
	defer func(begin time.Time) {
		mw.requestCount.With("method", "get_account").Add(1)
		mw.requestLatency.With("method", "get_account").Observe(time.Since(begin).Seconds())
	}(time.Now())
	return mw.s.GetAccount(ctx, id)
}

func (mw *instrumentingMiddleware) ListAccounts(ctx context.Context) ([]Account, error) {
	defer func(begin time.Time) {
		mw.requestCount.With("method", "list_account").Add(1)
		mw.requestLatency.With("method", "list_account").Observe(time.Since(begin).Seconds())
	}(time.Now())
	return mw.s.ListAccounts(ctx)
}

func (mw *instrumentingMiddleware) ActivateAccount(ctx context.Context, id string) error {
	defer func(begin time.Time) {
		mw.requestCount.With("method", "activate_account").Add(1)
		mw.requestLatency.With("method", "activate_account").Observe(time.Since(begin).Seconds())
	}(time.Now())
	return mw.s.ActivateAccount(ctx, id)
}

func (mw *instrumentingMiddleware) DeactivateAccount(ctx context.Context, id string) error {
	defer func(begin time.Time) {
		mw.requestCount.With("method", "deactivate_account").Add(1)
		mw.requestLatency.With("method", "deactivate_account").Observe(time.Since(begin).Seconds())
	}(time.Now())
	return mw.s.DeactivateAccount(ctx, id)
}
