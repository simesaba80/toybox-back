package migration

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/simesaba80/toybox-back/internal/infrastructure/database/dto"
	"github.com/uptrace/bun"
)

func MigrateFavorites(ctx context.Context, sourceDB, targetDB bun.IDB) error {
	fmt.Println("Migrating favorites...")
	var oldFavorites []OldFavorite
	if err := sourceDB.NewSelect().Model(&oldFavorites).Scan(ctx); err != nil {
		return fmt.Errorf("failed to select old favorites: %w", err)
	}

	newFavorites := make([]*dto.Favorite, 0, len(oldFavorites))
	for _, old := range oldFavorites {
		parsedOldWorkID, err := uuid.Parse(old.WorkID)
		if err != nil {
			return fmt.Errorf("failed to parse work UUID %s for favorite: %w", old.WorkID, err)
		}
		parsedOldUserID, err := uuid.Parse(old.UserID)
		if err != nil {
			return fmt.Errorf("failed to parse user UUID %s for favorite: %w", old.UserID, err)
		}
		newFavorite := &dto.Favorite{
			WorkID:    parsedOldWorkID,
			UserID:    parsedOldUserID,
			CreatedAt: old.CreatedAt,
		}
		newFavorites = append(newFavorites, newFavorite)
	}

	if len(newFavorites) > 0 {
		if _, err := targetDB.NewInsert().Model(&newFavorites).Exec(ctx); err != nil {
			return fmt.Errorf("failed to insert new favorites: %w", err)
		}
	}

	fmt.Printf("-> Migrated %d favorites.\n", len(newFavorites))
	return nil
}
