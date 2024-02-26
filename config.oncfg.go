package pgmfx

import (
	"context"
	"encoding/json"

	"github.com/myfantasy/mfctx"
	"gitlab.com/myfantasy.ru/tools/oncfg"
)

func (cfg *ConnectConfig) FeelCfg(ctxIn context.Context, cfgRaw *oncfg.ConfigRaw) (err error) {

	ctx := mfctx.FromCtx(ctxIn).Start("ConnectConfig.FeelCfg")
	defer func() { ctx.Complete(ctxIn.Err()) }()

	oncfg.SetDefault(cfg)
	oncfg.SetDefaultPwd(cfg)

	err = json.Unmarshal(cfgRaw.Cfg, cfg)
	if err != nil {
		return err
	}

	return nil
}
