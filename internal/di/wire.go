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
	"github.com/simesaba80/toybox-back/internal/infrastructure/router"
	"github.com/simesaba80/toybox-back/internal/interface/controller"
	"github.com/simesaba80/toybox-back/internal/usecase"
	"github.com/simesaba80/toybox-back/pkg/db"
)

// ProviderSet は依存関係を定義します
var ProviderSet = wire.NewSet(
	// データベース
	ProvideDatabase,

	// リポジトリ
	user.NewUserRepository,
	wire.Bind(new(repository.UserRepository), new(*user.UserRepository)),

	// ユースケース
	NewUserUseCaseProvider,

	// コントローラー
	controller.NewUserController,
	wire.Bind(new(controller.UserController), new(*controller.UserController)),

	// ルーター
	router.NewRouter,
	ProvideEcho,
)

// ProvideDatabase はデータベース接続を提供します
func ProvideDatabase() *bun.DB {
	return db.DB
}

// NewUserUseCaseProvider はUserUseCaseを作成します
func NewUserUseCaseProvider(repo repository.UserRepository) *usecase.UserUseCase {
	return usecase.NewUserUseCase(repo, 30*time.Second)
}

// ProvideEcho はEchoインスタンスを提供します
func ProvideEcho() *echo.Echo {
	return echo.New()
}

// InitializeRouter はルーターを初期化します
func InitializeRouter() *router.Router {
	wire.Build(ProviderSet)
	return nil
}
