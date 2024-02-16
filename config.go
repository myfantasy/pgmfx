package pgmfx

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/myfantasy/mfctx"
	"github.com/myfantasy/poh"
)

type RowsProcess func(rows pgx.Rows) error

type PgxMFPool struct {
	Pool                    *poh.ConnectionPool[*pgx.Conn]
	PoolName                func() string
	MetricsProvider         MetricsProvider
	GetNewConnectionTimeout poh.ExpireDurationFunc
}

// ConnectConfig connection config and pool config
type ConnectConfig struct {
	ConnString     string `envconfig:"conn_string" default:"user=postgres host=localhost port=5432 dbname=postgres" json:"conn_string"`
	Password       string `envconfig:"password" default_pwd:"postgres" json:"password"`
	MaxConnections int    `envconfig:"max_connections" default:"5" json:"max_connections"`
	MinConnections int    `envconfig:"min_connections" default:"0" json:"min_connections"`
	// OpenTimeout time in open state; 0 mean forever open
	OpenTimeout              time.Duration `envconfig:"open_timeout" default:"600s" json:"open_timeout"`
	IdleTimeout              time.Duration `envconfig:"idle_timeout" default:"600s" json:"idle_timeout"`
	PoolGetConnectionTimeout time.Duration `envconfig:"get_connection_timeout" default:"5s" json:"get_connection_timeout"`
	PoolName                 string        `envconfig:"pool_name" json:"pool_name"`
}

func (cfg *ConnectConfig) MakeConnection(ctxIn context.Context) (cnct *poh.Connection[*pgx.Conn], err error) {
	ctx := mfctx.FromCtx(ctxIn).Start("pgmfx.ConnectConfig.MakeConnection")
	defer func() { ctx.Complete(err) }()

	pgxCfg, err := pgx.ParseConfig(cfg.ConnString)
	if err != nil {
		return nil, err
	}

	if cfg.Password != "" {
		pgxCfg.Password = cfg.Password
	}

	pgxConn, err := pgx.ConnectConfig(ctx, pgxCfg)

	if err != nil {
		return nil, err
	}

	cnct = poh.MakeConnection(pgxConn,
		func(ctx context.Context, conn *pgx.Conn) error { return conn.Close(ctx) },
		func() time.Duration { return cfg.OpenTimeout },
		func() time.Duration { return cfg.IdleTimeout },
	)

	return cnct, nil
}

func (cfg *ConnectConfig) MakePool(ctxBase context.Context) (pool *PgxMFPool, err error) {
	ctx := mfctx.FromCtx(ctxBase).Start("pgmfx.ConnectConfig.MakeConnection")
	defer func() { ctx.Complete(err) }()
	cp := poh.MakeConnectionPool[*pgx.Conn](ctxBase,
		cfg.MakeConnection,
		func() int { return cfg.MaxConnections },
		func() int { return cfg.MinConnections },
	)

	return &PgxMFPool{
		Pool: cp,
		PoolName: func() string {
			return cfg.PoolName
		},
		GetNewConnectionTimeout: func() time.Duration {
			return cfg.PoolGetConnectionTimeout
		},
	}, nil
}

func RowsProcessNoRows(rows pgx.Rows) error {
	for rows.Next() {
		if rows.Err() != nil {
			return rows.Err()
		}
	}

	if rows.Err() != nil {
		return rows.Err()
	}

	return nil
}

func RowsProcessStructFronJson[T any](objs *[]T, exists *bool) RowsProcess {
	var js json.RawMessage
	rf := RowsProcessJson(&js, exists)
	return func(rows pgx.Rows) error {
		err := rf(rows)
		if err != nil {
			return errors.Join(ErrGetJsonFail, err)
		}

		if !*exists {
			return nil
		}

		err = json.Unmarshal(js, objs)
		if err != nil {
			return errors.Join(ErrUnmarshalJsonFail, err)
		}

		return nil
	}
}

func RowsProcessJson(js *json.RawMessage, exists *bool) RowsProcess {
	return func(rows pgx.Rows) error {
		for rows.Next() {
			if rows.Err() != nil {
				return rows.Err()
			}

			vals := rows.RawValues()

			if len(vals) >= 0 {
				jsV := vals[0]
				*js = jsV
				*exists = true
			}
		}

		if rows.Err() != nil {
			return rows.Err()
		}

		return nil
	}
}

func (pool PgxMFPool) Query(ctxIn context.Context, rp RowsProcess, sqlName string, sql string, args ...any) (err error) {
	ctx := mfctx.FromCtx(ctxIn).Start("pgmfx.PgxMFPool.Run:" + pool.PoolName() + ":" + sqlName)
	defer func() { ctx.Complete(err) }()

	tStart := time.Now()
	defer func() {
		if errors.Is(err, ErrQueryPoolGet) {
			WriteMetricResponce(pool.MetricsProvider, tStart, pool.PoolName(), sqlName, ErrorPoolStatus)
		} else if errors.Is(err, ErrQueryCall) {
			WriteMetricResponce(pool.MetricsProvider, tStart, pool.PoolName(), sqlName, ErrorQueryStatus)
		} else if errors.Is(err, ErrRowsProcess) {
			WriteMetricResponce(pool.MetricsProvider, tStart, pool.PoolName(), sqlName, ErrorRowsReadStatus)
		} else if err != nil {
			WriteMetricResponce(pool.MetricsProvider, tStart, pool.PoolName(), sqlName, ErrorStatus)
		} else {
			WriteMetricResponce(pool.MetricsProvider, tStart, pool.PoolName(), sqlName, OKStatus)
		}
	}()

	WriteMetricRequest(pool.MetricsProvider, pool.PoolName(), sqlName)

	ctxPool := ctx

	if pool.GetNewConnectionTimeout != nil {
		var cancel context.CancelFunc
		ctxPool, cancel = ctxPool.WithTimeout(pool.GetNewConnectionTimeout())
		defer cancel()
	}

	conn, free, err := pool.Pool.GetWait(ctxPool)
	defer free()

	if err != nil {
		WriteMetricGetPool(pool.MetricsProvider, tStart, pool.PoolName(), sqlName, ErrorStatus)
		return errors.Join(ErrQueryPoolGet, err)
	}

	WriteMetricGetPool(pool.MetricsProvider, tStart, pool.PoolName(), sqlName, OKStatus)

	ctx.With(poh.ConnectionIDLogParam, conn.ID)
	qStart := time.Now()

	rows, err := conn.Conn.Query(ctx, sql, args...)
	if err != nil {
		WriteMetricDoRequest(pool.MetricsProvider, qStart, pool.PoolName(), sqlName, ErrorStatus)
		return errors.Join(ErrQueryCall, err)
	}
	defer rows.Close()
	WriteMetricDoRequest(pool.MetricsProvider, qStart, pool.PoolName(), sqlName, OKStatus)

	rStart := time.Now()

	if rp == nil {
		rp = RowsProcessNoRows
	}

	err = rp(rows)
	if err != nil {
		WriteMetricDoReadResponce(pool.MetricsProvider, rStart, pool.PoolName(), sqlName, ErrorStatus)
		return errors.Join(ErrRowsProcess, err)
	}
	WriteMetricDoReadResponce(pool.MetricsProvider, rStart, pool.PoolName(), sqlName, OKStatus)

	return nil
}
