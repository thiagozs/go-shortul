package initialize

import (
	"fmt"

	"github.com/thiagozs/go-shorturl/api"
	"github.com/thiagozs/go-shorturl/config"
	"github.com/thiagozs/go-shorturl/handler"
	"github.com/thiagozs/go-shorturl/infra/database"
	"github.com/thiagozs/go-shorturl/middleware"
)

type Initialize struct {
	params *InitializeParams
}

func NewInitialize(opts ...Options) (*Initialize, error) {
	params, err := newInitializeParams(opts...)
	if err != nil {
		return nil, err
	}

	return &Initialize{params: params}, nil
}

func (i *Initialize) Init(reload ...string) error {

	if i.params.GetLogger() == nil {
		return fmt.Errorf("logger is required")
	}

	// Load configuration
	cfg, err := config.NewConfig()
	if err != nil {
		return err
	}

	if len(reload) > 0 {
		i.params.GetLogger().Info("Reloading configuration")
		cfg = i.params.GetConfig()
		i.params.SetConfig(cfg)
	} else {
		i.params.SetConfig(cfg)
	}

	// Create a database connection
	db, err := database.NewDatabase(database.Memory, i.params.GetLogger())
	if err != nil {
		return err
	}

	i.params.SetDB(db)

	// Load handlers
	handlerOpts := []handler.Options{
		handler.WithStore(db),
		handler.WithLogger(i.params.GetLogger()),
		handler.WithPort(cfg.GetPort()),
		handler.WithDomain(cfg.GetDomain()),
		handler.WithHost(cfg.GetHost()),
		handler.WithLocal(cfg.GetLocal()),
		handler.WithHTTPS(cfg.GetHTTPS()),
		handler.WithConfig(cfg),
	}

	hd, err := handler.NewHandler(handlerOpts...)
	if err != nil {
		return err
	}

	i.params.SetHandler(hd)

	// Load Midlewares
	optsMid := []middleware.Options{
		middleware.WithLogger(i.params.GetLogger()),
		middleware.WithToken(cfg.GetToken()),
	}

	md, err := middleware.NewMiddleware(optsMid...)
	if err != nil {
		return err
	}

	i.params.SetMiddleware(md)

	// Load API with handlers and middlewares
	apiOpts := []api.Options{
		api.WithLogger(i.params.GetLogger()),
		api.WithPort(cfg.GetPort()),
		api.WithDomain(cfg.GetDomain()),
		api.WithHost(cfg.GetHost()),
		api.WithHTTPS(cfg.GetHTTPS()),
		api.WithDB(db),
		api.WithMiddleware(md),
		api.WithHandlers(hd),
		api.WithConfig(cfg),
	}

	ap, err := api.NewApi(apiOpts...)
	if err != nil {
		return err
	}

	i.params.SetAPI(ap)

	return nil
}

func (i *Initialize) GetParams() *InitializeParams {
	return i.params
}

func (i *Initialize) SetConfigByFlags(cfg *config.Config) {
	i.params.SetConfig(cfg)
}

func (i *Initialize) ReloadInit() error {
	if i.params.GetLogger() == nil {
		return fmt.Errorf("logger is required")
	}

	return i.Init("reload")
}
