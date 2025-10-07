package usecase_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/simesaba80/toybox-back/internal/domain/entity"
	"github.com/simesaba80/toybox-back/internal/usecase"
)

func TestUserUseCase_CreateUser(t *testing.T) {
	user := &entity.User{
		ID:                  uuid.New(),
		Name:                "test",
		Email:               "test@test.com",
		PasswordHash:        "password",
		DisplayName:         "test",
		AvatarURL:           "https://test.com/avatar.png",
		DiscordToken:        "test",
		DiscordRefreshToken: "test",
		DiscordUserID:       "test",
		Profile:             "test",
		TwitterID:           "test",
		GithubID:            "test",
		CreatedAt:           time.Now(),
		UpdatedAt:           time.Now(),
	}
	userUseCase := usecase.NewUserUseCase(nil, 30*time.Second)
	user, err := userUseCase.CreateUser(context.Background(), user.Name, user.Email, user.PasswordHash, user.DisplayName, user.AvatarURL)
	if err != nil {
		t.Fatalf("failed to create user: %v", err)
	}
	if user.ID != user.ID {
		t.Fatalf("user ID mismatch: %v != %v", user.ID, user.ID)
	}
	if user.Name != user.Name {
		t.Fatalf("user Name mismatch: %v != %v", user.Name, user.Name)
	}
	if user.Email != user.Email {
		t.Fatalf("user Email mismatch: %v != %v", user.Email, user.Email)
	}
	if user.DisplayName != user.DisplayName {
		t.Fatalf("user DisplayName mismatch: %v != %v", user.DisplayName, user.DisplayName)
	}
	if user.AvatarURL != user.AvatarURL {
		t.Fatalf("user AvatarURL mismatch: %v != %v", user.AvatarURL, user.AvatarURL)
	}
	if user.CreatedAt.IsZero() {
		t.Fatalf("user CreatedAt is zero")
	}
	if user.UpdatedAt.IsZero() {
		t.Fatalf("user UpdatedAt is zero")
	}
	if user.PasswordHash != user.PasswordHash {
		t.Fatalf("user PasswordHash mismatch: %v != %v", user.PasswordHash, user.PasswordHash)
	}
	if user.DiscordToken != user.DiscordToken {
		t.Fatalf("user DiscordToken mismatch: %v != %v", user.DiscordToken, user.DiscordToken)
	}
	if user.DiscordRefreshToken != user.DiscordRefreshToken {
		t.Fatalf("user DiscordRefreshToken mismatch: %v != %v", user.DiscordRefreshToken, user.DiscordRefreshToken)
	}
	if user.DiscordUserID != user.DiscordUserID {
		t.Fatalf("user DiscordUserID mismatch: %v != %v", user.DiscordUserID, user.DiscordUserID)
	}
	if user.Profile != user.Profile {
		t.Fatalf("user Profile mismatch: %v != %v", user.Profile, user.Profile)
	}
	if user.TwitterID != user.TwitterID {
		t.Fatalf("user TwitterID mismatch: %v != %v", user.TwitterID, user.TwitterID)
	}
	if user.GithubID != user.GithubID {
		t.Fatalf("user GithubID mismatch: %v != %v", user.GithubID, user.GithubID)
	}
	if user.CreatedAt.IsZero() {
		t.Fatalf("user CreatedAt is zero")
	}
	if user.UpdatedAt.IsZero() {
		t.Fatalf("user UpdatedAt is zero")
	}
}
