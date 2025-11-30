package controller

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
	domainerrors "github.com/simesaba80/toybox-back/internal/domain/errors"
	"github.com/simesaba80/toybox-back/internal/interface/schema"
	"github.com/simesaba80/toybox-back/internal/usecase"
)

type AuthController struct {
	authUsecase usecase.IAuthUsecase
}

func NewAuthController(authUsecase usecase.IAuthUsecase) *AuthController {
	return &AuthController{
		authUsecase: authUsecase,
	}
}

// GetDiscordAuthURL godoc
// @Summary Get Discord authentication URL
// @Description Get Discord authentication URL
// @Tags auth
// @Produce json
// @Success 200 {object} schema.GetDiscordAuthURLResponse
// @Failure 500 {object} echo.HTTPError
// @Router /auth/discord/ [get]
func (ac *AuthController) GetDiscordAuthURL(c echo.Context) error {
	authURL, err := ac.authUsecase.GetDiscordAuthURL(c.Request().Context())
	if err != nil {
		return handleAuthError(c, err)
	}
	return c.JSON(http.StatusOK, schema.ToGetDiscordAuthURLResponse(authURL))
}

// GetDiscordToken godoc
// @Summary Get Discord token by code
// @Description Get Discord token by code
// @Tags auth
// @Produce json
// @Success 200 {object} schema.GetDiscordTokenResponse
// @Failure 400 {object} echo.HTTPError
// @Failure 403 {object} echo.HTTPError
// @Failure 500 {object} echo.HTTPError
// @Router /auth/discord/callback [get]
// @Param code query string true "Discord code"
func (ac *AuthController) AuthenticateUser(c echo.Context) error {
	code := c.QueryParam("code")
	if code == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "code is required")
	}
	appToken, refreshToken, err := ac.authUsecase.AuthenticateUser(c.Request().Context(), code)
	if err != nil {
		return handleAuthError(c, err)
	}
	return c.JSON(http.StatusOK, schema.ToGetDiscordTokenResponse(appToken, refreshToken))
}

// RegenerateToken godoc
// @Summary Regenerate token
// @Description Regenerate token
// @Tags auth
// @Produce json
// @Success 200 {object} schema.GetDiscordTokenResponse
// @Failure 400 {object} echo.HTTPError
// @Failure 500 {object} echo.HTTPError
// @Router /auth/refresh [post]
// @Param refresh_token query string true "Refresh token"
func (ac *AuthController) RegenerateToken(c echo.Context) error {
	var input schema.RegenerateTokenInput
	if err := c.Bind(&input); err != nil {
		c.Logger().Error("Bind error:", err)
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
	}
	if err := c.Validate(&input); err != nil {
		return err
	}

	appToken, err := ac.authUsecase.RegenerateToken(c.Request().Context(), input.RefreshToken)
	if err != nil {
		return handleAuthError(c, err)
	}
	return c.JSON(http.StatusOK, schema.ToRegenerateTokenResponse(appToken))
}

func handleAuthError(c echo.Context, err error) error {
	var httpErr *echo.HTTPError
	if errors.As(err, &httpErr) {
		return httpErr
	}

	switch {
	case errors.Is(err, domainerrors.ErrUserNotAllowedGuild):
		return echo.NewHTTPError(http.StatusForbidden, "ユーザーは許可されたDiscordギルドに所属していません")
	case errors.Is(err, domainerrors.ErrRefreshTokenExpired):
		return echo.NewHTTPError(http.StatusBadRequest, "リフレッシュトークンが期限切れです")
	case errors.Is(err, domainerrors.ErrRefreshTokenInvalid):
		return echo.NewHTTPError(http.StatusBadRequest, "リフレッシュトークンが無効です")
	case errors.Is(err, domainerrors.ErrFaileRequestToDiscord):
		return echo.NewHTTPError(http.StatusInternalServerError, "Discordへのリクエストに失敗しました")
	case errors.Is(err, domainerrors.ErrClientIDNotSet):
		return echo.NewHTTPError(http.StatusInternalServerError, "DiscordクライアントIDが設定されていません")
	case errors.Is(err, domainerrors.ErrRedirectURLNotSet):
		return echo.NewHTTPError(http.StatusInternalServerError, "リダイレクトURLが設定されていません")
	case errors.Is(err, domainerrors.ErrFailedToCreateUser):
		return echo.NewHTTPError(http.StatusInternalServerError, "ユーザーの作成に失敗しました")
	default:
		c.Logger().Error("Auth error:", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Internal server error")
	}
}
