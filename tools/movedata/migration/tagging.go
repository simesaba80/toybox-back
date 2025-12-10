package migration

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/simesaba80/toybox-back/internal/infrastructure/database/dto"
	"github.com/uptrace/bun"
)

func MigrateTaggings(ctx context.Context, sourceDB, targetDB bun.IDB, tagIdMap map[uuid.UUID]uuid.UUID) error {
	fmt.Println("Migrating taggings...")
	var oldTaggings []OldTagging
	if err := sourceDB.NewSelect().Model(&oldTaggings).Scan(ctx); err != nil {
		return fmt.Errorf("failed to select old taggings: %w", err)
	}

	newTaggings := make([]*dto.Tagging, 0, len(oldTaggings))
	processedTaggings := make(map[[2]uuid.UUID]struct{})

	for _, old := range oldTaggings {
		parsedOldWorkID, err := uuid.Parse(old.WorkID)
		if err != nil {
			return fmt.Errorf("failed to parse work UUID %s for tagging: %w", old.WorkID, err)
		}
		parsedOldTagID, err := uuid.Parse(old.TagID)
		if err != nil {
			return fmt.Errorf("failed to parse tag UUID %s for tagging: %w", old.TagID, err)
		}

		representativeTagID, exists := tagIdMap[parsedOldTagID]
		if !exists {
			// このケースは発生しないはずだが念のため確認を行う
			return fmt.Errorf("failed to find representative tag ID for old tag ID %s", old.TagID)
		}
		// 同一のUUIDに統合されたタグが一つの作品に複数紐づくのを防ぐため行う
		key := [2]uuid.UUID{parsedOldWorkID, representativeTagID}
		if _, exists := processedTaggings[key]; !exists {
			newTaggings = append(newTaggings, &dto.Tagging{
				WorkID: parsedOldWorkID,
				TagID:  representativeTagID,
			})
			processedTaggings[key] = struct{}{}
		}
	}

	if len(newTaggings) > 0 {
		if _, err := targetDB.NewInsert().Model(&newTaggings).Exec(ctx); err != nil {
			return fmt.Errorf("failed to insert new taggings: %w", err)
		}
	}

	fmt.Printf("-> Migrated %d taggings.\n", len(newTaggings))
	return nil
}
