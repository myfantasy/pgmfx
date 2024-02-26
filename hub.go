package pgmfx

import (
	"context"
	"fmt"
	"sync"

	"github.com/myfantasy/mfctx"
	"github.com/myfantasy/poh"
	"gitlab.com/myfantasy.ru/tools/oncfg"
)

var ErrPgxMFHubClosePoint = fmt.Errorf("close point in hub")

type PgxMFHub struct {
	ctxBase  context.Context
	ctxClose context.CancelCauseFunc

	cfg *oncfg.ListSourceConfig
	hub *poh.Hub[string, *PgxMFPool]

	mx sync.Mutex
}

func MakePgxMFHub(ctxBase context.Context) *PgxMFHub {
	res := &PgxMFHub{
		ctxBase: ctxBase,
		cfg:     oncfg.ListSourceConfigCreate(),
	}

	res.hub = poh.MakeHub[string, *PgxMFPool](
		ctxBase,
		func(ctx context.Context) (keys []string, err error) {
			return res.cfg.Keys(ctx)
		},
		func(ctx context.Context, key string, point *PgxMFPool) (err error) {
			point.Pool.Close(ErrPgxMFHubClosePoint)
			return nil
		},
		func(ctxIn context.Context, key string) (point *PgxMFPool, err error) {
			if res.ctxBase.Err() != nil {
				return nil, res.ctxBase.Err()
			}

			cfg := &ConnectConfig{}

			oncfg.SetDefault(cfg)
			oncfg.SetDefaultPwd(cfg)

			ctx := mfctx.FromCtx(ctxIn).Copy().With("key", key)

			rCfg, err := res.cfg.Config(ctx, key)
			if err != nil {
				return nil, err
			}

			err = cfg.FeelCfg(ctx, &rCfg)
			if err != nil {
				return nil, err
			}

			ctx, cancel := ctx.WithCancelCause()

			mpl, err := cfg.MakePool(ctx, cancel)
			if err != nil {
				return nil, err
			}

			return mpl, nil
		},
		func(ctxIn context.Context, key string, point *PgxMFPool) (err error) {
			cfg := point.Cfg()

			ctx := mfctx.FromCtx(ctxIn).Copy().With("key", key)

			rCfg, err := res.cfg.Config(ctx, key)
			if err != nil {
				return err
			}

			err = cfg.FeelCfg(ctx, &rCfg)
			if err != nil {
				return err
			}

			return nil
		},
	)

	res.cfg.SetCallAfterUpdateFunc(func(ctx context.Context) {
		res.hub.Refresh()
	})

	return res
}

func (pgxhub *PgxMFHub) Cfg() *oncfg.ListSourceConfig {
	pgxhub.mx.Lock()
	defer pgxhub.mx.Unlock()

	return pgxhub.cfg
}

func (pgxhub *PgxMFHub) Get(key string) (point *PgxMFPool, ok bool) {
	pgxhub.mx.Lock()
	defer pgxhub.mx.Unlock()

	return pgxhub.hub.Get(key)
}

// Close - calls ctxClose if its set (error should be set)
func (pgxhub *PgxMFHub) Close(err error) {
	pgxhub.mx.Lock()
	defer pgxhub.mx.Unlock()

	if pgxhub.ctxClose != nil {
		pgxhub.ctxClose(err)
	}
}
