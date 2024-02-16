package main

import (
	"context"
	"fmt"
	"time"

	"github.com/myfantasy/mfctx"
	"github.com/myfantasy/mfctx/jsonify"
	"github.com/myfantasy/pgmfx"
	"github.com/myfantasy/pgmfx/oncfgsrc"
	"gitlab.com/myfantasy.ru/tools/oncfg"
)

func main() {
	mfctx.SetAppName("d.2.5.6")

	cfg := &pgmfx.ConnectConfig{}

	oncfg.SetDefault(cfg)
	oncfg.SetDefaultPwd(cfg)

	pool, err := cfg.MakePool(context.Background())

	if err != nil {
		panic(err)
	}

	ctx := mfctx.FromCtx(context.Background())

	ctx, cancel := ctx.WithTimeout(time.Second * 10)
	defer cancel()

	cfgSrc := oncfgsrc.ConfigSource{Pool: pool}

	cfgRaw, exists, err := cfgSrc.GetConfig(ctx, "n1.fun")

	if err != nil {
		panic(err)
	}

	if !exists {
		panic("config not exists")
	}

	fmt.Println(jsonify.JsonifySLn(cfgRaw))
}
