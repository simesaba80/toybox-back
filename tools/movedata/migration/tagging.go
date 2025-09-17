package migration

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/simesaba80/toybox-back/internal/infrastructure/database/dto"
	"github.com/uptrace/bun"
)

func MigrateTaggings(ctx context.Context, sourceDB, targetDB bun.IDB) error {
	fmt.Println("Migrating taggings...")
	var oldTaggings []OldTagging
	if err := sourceDB.NewSelect().Model(&oldTaggings).Scan(ctx); err != nil {
		return fmt.Errorf("failed to select old taggings: %w", err)
	}

	newTaggings := make([]*dto.Tagging, 0, len(oldTaggings))
	for _, old := range oldTaggings {
		parsedOldWorkID, err := uuid.Parse(old.WorkID)
		if err != nil {
			return fmt.Errorf("failed to parse work UUID %s for tagging: %w", old.WorkID, err)
		}
		parsedOldTagID, err := uuid.Parse(old.TagID)
		if err != nil {
			return fmt.Errorf("failed to parse tag UUID %s for tagging: %w", old.TagID, err)
		}
		newTagging := &dto.Tagging{
			WorkID: parsedOldWorkID,
			TagID:  parsedOldTagID,
		}
		newTaggings = append(newTaggings, newTagging)
	}

	if len(newTaggings) > 0 {
		if _, err := targetDB.NewInsert().Model(&newTaggings).Exec(ctx); err != nil {
			return fmt.Errorf("failed to insert new taggings: %w", err)
		}
	}

	fmt.Printf("-> Migrated %d taggings.\n", len(newTaggings))
	return nil
}
