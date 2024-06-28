package container

import (
	"os"
	"time"

	"github.com/armiariyan/bepkg/logger"
	"github.com/armiariyan/synapsis/internal/config"
	"github.com/armiariyan/synapsis/internal/domain/entities"
	"github.com/armiariyan/synapsis/internal/domain/repositories"
	"github.com/armiariyan/synapsis/internal/infrastructure/postgresql"
	"github.com/armiariyan/synapsis/internal/pkg/log"
	"github.com/armiariyan/synapsis/internal/usecase/carts"
	"github.com/armiariyan/synapsis/internal/usecase/healthcheck"
	"github.com/armiariyan/synapsis/internal/usecase/products"
	"github.com/armiariyan/synapsis/internal/usecase/user"
	"gorm.io/gorm"
)

type Container struct {
	Config             *config.DefaultConfig
	PostgresqlDB       *config.PostgresqlDB
	SynapsisDB         *gorm.DB
	Logger             logger.Logger
	HealthCheckService healthcheck.Service
	UserService        user.Service
	ProductService     products.Service
	CartService        carts.Service
}

func (c *Container) Validate() *Container {
	if c.Config == nil {
		panic("Config is nil")
	}
	if c.SynapsisDB == nil {
		panic("SynapsisDB is nil")
	}
	if c.Logger == nil {
		panic("Logger is nil")
	}
	if c.HealthCheckService == nil {
		panic("HealthCheckService is nil")
	}
	if c.UserService == nil {
		panic("UserService is nil")
	}
	if c.ProductService == nil {
		panic("ProductService is nil")
	}
	if c.CartService == nil {
		panic("CartService is nil")
	}
	return c
}

func New() *Container {
	config.Load(os.Getenv("env"), ".env")

	fileLoc := config.GetString("logger.fileLocation")
	tdrFileLoc := config.GetString("logger.fileTdrLocation")
	maxAge := time.Duration(config.GetInt("logger.fileMaxAge"))
	stdOut := config.GetBool("logger.stdout")

	defLogger := logger.New(logger.Options{
		FileLocation:    fileLoc,
		FileTdrLocation: tdrFileLoc,
		FileMaxAge:      maxAge,
		Stdout:          stdOut,
	})

	defConfig := &config.DefaultConfig{
		Apps: config.Apps{
			Name:     config.GetString("app.name"),
			Address:  config.GetString("address"),
			HttpPort: config.GetString("port"),
		},
	}

	psqlConfig := &config.PostgresqlDB{
		Host:     config.GetString("postgresql.synapsis.host"),
		User:     config.GetString("postgresql.synapsis.user"),
		Password: config.GetString("postgresql.synapsis.password"),
		Name:     config.GetString("postgresql.synapsis.db"),
		Port:     config.GetInt("postgresql.synapsis.port"),
		SSLMode:  config.GetString("postgresql.synapsis.ssl"),
		Schema:   config.GetString("postgresql.synapsis.schema"),
		Debug:    config.GetBool("postgresql.synapsis.debug"),
	}

	log.New()

	synapsisDB := postgresql.NewDB(*psqlConfig)
	if config.GetString("env") == "development" {
		postgresql.CreateUUIDExtension(synapsisDB)
		synapsisDB.AutoMigrate(
			entities.ProductCategory{},
			entities.Product{},
			entities.User{},
			entities.Cart{},
		)
		// * for auto migrate feature
	}

	// * Repositories
	userRepository := repositories.NewUser(synapsisDB)
	productCategoryRepository := repositories.NewProductCategory(synapsisDB)
	productRepository := repositories.NewProduct(synapsisDB)
	cartRepository := repositories.NewCart(synapsisDB)

	// * Wrapper

	// * Services
	healthCheckService := healthcheck.NewService().Validate()
	userService := user.NewService().
		SetDB(synapsisDB).
		SetUserRepository(userRepository).
		Validate()
	productService := products.NewService().
		SetDB(synapsisDB).
		SetProductCategoryRepository(productCategoryRepository).
		SetProductRepository(productRepository).
		SetCartRepository(cartRepository).
		Validate()

	cartService := carts.NewService().
		SetDB(synapsisDB).
		SetProductCategoryRepository(productCategoryRepository).
		SetProductRepository(productRepository).
		SetCartRepository(cartRepository).
		Validate()

	// * Brokers

	// * Workers

	container := &Container{
		Config:             defConfig,
		Logger:             defLogger,
		PostgresqlDB:       psqlConfig,
		SynapsisDB:         synapsisDB,
		HealthCheckService: healthCheckService,
		UserService:        userService,
		ProductService:     productService,
		CartService:        cartService,
	}
	container.Validate()
	return container

}
