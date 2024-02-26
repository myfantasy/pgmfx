package main

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/myfantasy/mfctx"
	"github.com/myfantasy/pgmfx"
	"gitlab.com/myfantasy.ru/tools/oncfg"
)

func main() {
	cfg := &pgmfx.ConnectConfig{}

	oncfg.SetDefault(cfg)
	oncfg.SetDefaultPwd(cfg)

	bbb, _ := json.MarshalIndent(cfg, "", " ")
	fmt.Println(string(bbb))

	pool, err := cfg.MakePool(context.Background(), nil)

	if err != nil {
		panic(err)
	}

	ctx := mfctx.FromCtx(context.Background())

	ctx, cancel := ctx.WithTimeout(time.Second * 10)
	defer cancel()

	err = pool.Query(ctx,
		func(rows pgx.Rows) error {
			for rows.Next() {
				if rows.Err() != nil {
					return rows.Err()
				}

				fmt.Println("SUPER!")
			}

			if rows.Err() != nil {
				return rows.Err()
			}

			return nil
		},
		"test",
		"select 1 a",
	)

	if err != nil {
		panic(err)
	}
}
