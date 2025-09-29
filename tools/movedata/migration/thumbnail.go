package migration

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/simesaba80/toybox-back/internal/infrastructure/database/dto"
	"github.com/uptrace/bun"
)

func MigrateThumbnails(ctx context.Context, sourceDB, targetDB bun.IDB) error {
	fmt.Println("Migrating thumbnails...")
	var oldThumbnails []OldThumbnail
	if err := sourceDB.NewSelect().Model(&oldThumbnails).Scan(ctx); err != nil {
		return fmt.Errorf("failed to select old thumbnails: %w", err)
	}

	newThumbnails := make([]*dto.Thumbnail, 0, len(oldThumbnails))
	for _, old := range oldThumbnails {
		parsedOldWorkID, err := uuid.Parse(old.WorkID)
		if err != nil {
			return fmt.Errorf("failed to parse work UUID %s for thumbnail: %w", old.WorkID, err)
		}
		parsedOldAssetID, err := uuid.Parse(old.AssetID)
		if err != nil {
			return fmt.Errorf("failed to parse asset UUID %s for thumbnail: %w", old.AssetID, err)
		}
		newThumbnail := &dto.Thumbnail{
			WorkID:  parsedOldWorkID,
			AssetID: parsedOldAssetID,
		}
		newThumbnails = append(newThumbnails, newThumbnail)
	}

	if len(newThumbnails) > 0 {
		if _, err := targetDB.NewInsert().Model(&newThumbnails).Exec(ctx); err != nil {
			return fmt.Errorf("failed to insert new thumbnails: %w", err)
		}
	}

	fmt.Printf("-> Migrated %d thumbnails.\n", len(newThumbnails))
	return nil
}
