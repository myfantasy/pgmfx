package oncfgsrc

import (
	"context"
	"errors"

	"github.com/myfantasy/mfctx"
	"github.com/myfantasy/pgmfx"
	"gitlab.com/myfantasy.ru/tools/oncfg"
)

type ConfigSource struct {
	Pool *pgmfx.PgxMFPool
}

var _ oncfg.ConfigSource = &ConfigSource{}

func (cs *ConfigSource) GetConfig(ctxIn context.Context, name string) (cfgRaw *oncfg.ConfigRaw, exists bool, err error) {
	ctx := mfctx.FromCtx(ctxIn).Start("oncfgsrc.ConfigSource.GetConfig")
	defer func() { ctx.Complete(err) }()

	cfgs := []oncfg.ConfigRaw{}

	err = cs.Pool.Query(ctx, pgmfx.RowsProcessStructFronJson(&cfgs, &exists),
		"config.dynamic_config_get_one_by_name_and_app_name",
		"select config.dynamic_config_get_one_by_name_and_app_name(_data_center => $1, _app_name => $2, _name => $3)",
		mfctx.GetDataCenter(), mfctx.GetAppName(), name)

	if err != nil {
		return nil, false, errors.Join(ErrCallSourceFail, err)
	}

	if !exists {
		return nil, false, nil
	}

	if len(cfgs) == 0 {
		return nil, false, nil
	}

	cfgRes0, err := oncfg.ConfigRawDictUinion(cfgs...)
	if err != nil {
		return nil, false, errors.Join(ErrUnionConfigFail, err)
	}

	return &cfgRes0, true, err
}

// func (cs *ConfigSource) GetConfig111(ctxIn context.Context, name string) (cfgRaw *oncfg.ConfigRaw, exists bool, err error) {
// 	ctx := mfctx.FromCtx(ctxIn).Start("oncfgsrc.ConfigSource.GetConfig")
// 	defer func() { ctx.Complete(err) }()

// 	var js json.RawMessage

// 	err = cs.Pool.Query(ctx, pgmfx.RowsProcessJson(&js, &exists),
// 		"config.dynamic_config_get_one_by_name_and_app_name",
// 		"select config.dynamic_config_get_one_by_name_and_app_name(_app_name => $1, _name => $2)",
// 		mfctx.GetAppName(), name)

// 	if err != nil {
// 		return nil, false, errors.Join(ErrCallSourceFail, err)
// 	}

// 	if !exists {
// 		return nil, false, nil
// 	}

// 	cfgs := []oncfg.ConfigRaw{}

// 	err = json.Unmarshal(js, &cfgs)
// 	if err != nil {
// 		return nil, false, errors.Join(ErrUnmarshalJsonFail, err)
// 	}

// 	if len(cfgs) == 0 {
// 		return nil, false, nil
// 	}

// 	cfgRes0, err := oncfg.ConfigRawDictUinion(cfgs[0], cfgs[1:]...)
// 	if err != nil {
// 		return nil, false, errors.Join(ErrUnionConfigFail, err)
// 	}

// 	return &cfgRes0, true, err
// }
