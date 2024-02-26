package main

import (
	"context"
	"fmt"
	"time"

	"github.com/myfantasy/mfctx"
	"github.com/myfantasy/pgmfx"
	"github.com/myfantasy/pgmfx/oncfgsrc"
	"gitlab.com/myfantasy.ru/tools/oncfg"
)

func main() {
	mfctx.SetAppName("d.2.5.6")

	cfg := &pgmfx.ConnectConfig{}

	oncfg.SetDefault(cfg)
	oncfg.SetDefaultPwd(cfg)

	pool, err := cfg.MakePool(context.Background(), nil)

	if err != nil {
		panic(err)
	}

	// ctx := mfctx.FromCtx(context.Background())

	// ctx, cancel := ctx.WithTimeout(time.Second * 10)
	// defer cancel()

	cfgSrc := oncfgsrc.ConfigSource{Pool: pool}

	ctxBase := context.Background()

	pgHub := pgmfx.MakePgxMFHub(ctxBase)

	fmt.Println("Hub created")

	srcCfgDbl := oncfg.DoubleSourceConfigCreate(pgHub.Cfg().FeelCfg)

	go oncfg.LoadPermonent(ctxBase, &cfgSrc, "db.pg.test_all.cfg", func() time.Duration { return time.Second * 5 }, srcCfgDbl.FeelSecondCfg)
	go oncfg.LoadPermonent(ctxBase, &cfgSrc, "db.pg.test_all.pwd", func() time.Duration { return time.Second * 5 }, srcCfgDbl.FeelFirstCfg)

	fmt.Println("update start. Let`s sleep 2 sec")

	time.Sleep(2 * time.Second)

	// fmt.Println(ok)

	// fmt.Println(jsonify.JsonifySLn(pgHub.Cfg()))

	pool, ok := pgHub.Get("test1")
	if !ok {
		panic("should be OK")
	}

	err = pool.Query(ctxBase, nil, "test_sql", "select 1")
	if err != nil {
		panic(err)
	}

	fmt.Println("WOW hub is working")
}
