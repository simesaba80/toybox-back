package migration

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/simesaba80/toybox-back/internal/infrastructure/database/dto"
	"github.com/uptrace/bun"
)

func MigrateTokens(ctx context.Context, sourceDB, targetDB bun.IDB) error {
	fmt.Println("Migrating tokens...")

	var oldTokens []OldToken
	if err := sourceDB.NewSelect().Model(&oldTokens).Scan(ctx); err != nil {
		return fmt.Errorf("failed to select old tokens: %w", err)
	}

	newTokens := make([]*dto.Token, 0, len(oldTokens))
	for _, old := range oldTokens {
		parsedOldRefreshToken, err := uuid.Parse(old.RefreshToken)
		if err != nil {
			return fmt.Errorf("failed to parse token UUID %s: %w", old.RefreshToken, err)
		}
		parsedOldUserID, err := uuid.Parse(old.UserID)
		if err != nil {
			return fmt.Errorf("failed to parse user UUID %s: %w", old.UserID, err)
		}
		newToken := &dto.Token{
			RefreshToken: parsedOldRefreshToken,
			UserID:       parsedOldUserID,
			ExpiredAt:    old.ExpiredAt,
			CreatedAt:    old.CreatedAt,
			UpdatedAt:    old.UpdatedAt,
		}
		newTokens = append(newTokens, newToken)
	}
	if len(newTokens) > 0 {
		if _, err := targetDB.NewInsert().Model(&newTokens).Exec(ctx); err != nil {
			return fmt.Errorf("failed to insert new tokens: %w", err)
		}
	}
	fmt.Printf("-> Migrated %d tokens.\n", len(newTokens))
	return nil
}
