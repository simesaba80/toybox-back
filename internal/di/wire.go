//go:generate wire
//go:build wireinject
// +build wireinject

package di

import (
	"time"

	"github.com/google/wire"
	"github.com/labstack/echo/v4"
	"github.com/uptrace/bun"

	"github.com/simesaba80/toybox-back/internal/domain/repository"
	"github.com/simesaba80/toybox-back/internal/infrastructure/database/user"
	"github.com/simesaba80/toybox-back/internal/infrastructure/database/work"
	"github.com/simesaba80/toybox-back/internal/infrastructure/router"
	"github.com/simesaba80/toybox-back/internal/interface/controller"
	"github.com/simesaba80/toybox-back/internal/usecase"
	"github.com/simesaba80/toybox-back/pkg/db"
)

var RepositorySet = wire.NewSet(
	user.NewUserRepository,
	wire.Bind(new(repository.UserRepository), new(*user.UserRepository)),
	work.NewWorkRepository,
	wire.Bind(new(repository.WorkRepository), new(*work.WorkRepository)),
)

var UseCaseSet = wire.NewSet(
	ProvideUserUseCase,
	ProvideWorkUseCase,
)

var ControllerSet = wire.NewSet(
	controller.NewUserController,
	controller.NewWorkController,
)

var InfrastructureSet = wire.NewSet(
	ProvideDatabase,
	router.NewRouter,
	ProvideEcho,
)

// ProviderSet は依存関係を定義します
var ProviderSet = wire.NewSet(
	RepositorySet,
	UseCaseSet,
	ControllerSet,
	InfrastructureSet,
	NewApp, // App構造体のコンストラクタを追加
)

// ProvideDatabase はデータベース接続を提供します
func ProvideDatabase() *bun.DB {
	db.Init()
	return db.DB
}

// ProvideUserUseCase はUserUseCaseを提供します
func ProvideUserUseCase(repo repository.UserRepository) *usecase.UserUseCase {
	return usecase.NewUserUseCase(repo, 30*time.Second)
}

// ProvideWorkUseCase はWorkUseCaseを提供します
func ProvideWorkUseCase(repo repository.WorkRepository) *usecase.WorkUseCase {
	return usecase.NewWorkUseCase(repo, 30*time.Second)
}

// ProvideEcho はEchoインスタンスを提供します
func ProvideEcho() *echo.Echo {
	return echo.New()
}

// NewApp はAppインスタンスを作成します
func NewApp(router *router.Router, database *bun.DB) *App {
	return &App{
		Router:   router,
		Database: database,
	}
}

// InitializeApp はアプリケーションを初期化します
func InitializeApp() (*App, func(), error) {
	wire.Build(ProviderSet)
	return nil, nil, nil
}

type App struct {
	Router   *router.Router
	Database *bun.DB
}

// Start アプリケーションの開始
func (app *App) Start() *echo.Echo {
	return app.Router.Setup()
}

// Cleanup アプリケーションのクリーンアップ
func (app *App) Cleanup() {
	if app.Database != nil {
		app.Database.Close()
	}
}
