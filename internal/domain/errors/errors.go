package errors

import "errors"

// 認証関連のエラー定義
var (
	ErrUserNotAllowedGuild   = errors.New("user is not in an allowed discord guild")
	ErrRefreshTokenExpired   = errors.New("refresh token is expired")
	ErrRefreshTokenInvalid   = errors.New("refresh token is invalid")
	ErrFaileRequestToDiscord = errors.New("failed to request to discord")
	ErrClientIDNotSet        = errors.New("client ID is not set")
	ErrRedirectURLNotSet     = errors.New("redirect URL is not set")
)

// ユーザー関連のエラー定義
var (
	ErrUserNotFound      = errors.New("user not found")
	ErrUserAlreadyExists = errors.New("user already exists")
)
