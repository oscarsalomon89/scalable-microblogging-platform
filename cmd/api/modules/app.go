package modules

import (
	"github.com/oscarsalomon89/go-hexagonal/internal/platform/config"
	"github.com/oscarsalomon89/go-hexagonal/internal/platform/environment"
	db "github.com/oscarsalomon89/go-hexagonal/internal/platform/pg"
	clonctx "github.com/oscarsalomon89/go-hexagonal/pkg/context"
	"github.com/oscarsalomon89/go-hexagonal/pkg/validator"
	"go.uber.org/fx"
)

func NewApp() *fx.App {
	cfg, err := config.NewConfig()
	if err != nil {
		panic(err)
	}

	return NewAppWithConfig(cfg)
}

func NewAppWithConfig(cfg config.Configuration) *fx.App {
	options := []fx.Option{
		fx.Provide(func() config.Configuration { return cfg }),
		fx.Provide(func() config.Database { return cfg.Database }),
		fx.Provide(func() config.Cache { return cfg.Cache }),
		internalModule,
		userModule,
		tweetModule,
	}

	return fx.New(
		fx.Options(options...),
		fx.Invoke(validator.RegisterValidation),
		fx.Invoke(clonctx.NewLogger),
		fx.Invoke(runMigrations),
	)
}

func runMigrations(cfg config.Configuration) error {
	if environment.GetFromString(cfg.Scope) != environment.Local {
		return nil
	}

	return db.RunMigrations(cfg.Database)
}
