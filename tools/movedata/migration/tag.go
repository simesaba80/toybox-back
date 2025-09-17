package migration

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/simesaba80/toybox-back/internal/infrastructure/database/dto"
	"github.com/uptrace/bun"
)

func MigrateTags(ctx context.Context, sourceDB, targetDB bun.IDB) error {
	fmt.Println("Migrating tags...")

	var oldTags []OldTag
	if err := sourceDB.NewSelect().Model(&oldTags).Scan(ctx); err != nil {
		return fmt.Errorf("failed to select old tags: %w", err)
	}

	newTags := make([]*dto.Tag, 0, len(oldTags))
	for _, old := range oldTags {
		parsedOldID, err := uuid.Parse(old.ID)
		if err != nil {
			return fmt.Errorf("failed to parse tag UUID %s: %w", old.ID, err)
		}
		newTag := &dto.Tag{
			ID:        parsedOldID,
			Name:      old.Name,
			CreatedAt: old.CreatedAt,
			UpdatedAt: old.UpdatedAt,
		}
		newTags = append(newTags, newTag)
	}

	if len(newTags) > 0 {
		if _, err := targetDB.NewInsert().Model(&newTags).Exec(ctx); err != nil {
			return fmt.Errorf("failed to insert new tags: %w", err)
		}
	}

	fmt.Printf("-> Migrated %d tags.\n", len(newTags))
	return nil
}
