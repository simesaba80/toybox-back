package migration

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/simesaba80/toybox-back/internal/infrastructure/database/dto"
	"github.com/simesaba80/toybox-back/internal/infrastructure/database/types"
	"github.com/uptrace/bun"
)

func MigrateAssets(ctx context.Context, sourceDB, targetDB bun.IDB) error {
	fmt.Println("Migrating assets...")
	var oldAssets []OldAsset
	if err := sourceDB.NewSelect().Model(&oldAssets).Scan(ctx); err != nil {
		return fmt.Errorf("failed to select old assets: %w", err)
	}

	newAssets := make([]*dto.Asset, 0, len(oldAssets))
	for _, old := range oldAssets {

		parsedOldID, err := uuid.Parse(old.ID)
		if err != nil {
			return fmt.Errorf("failed to parse asset UUID %s: %w", old.ID, err)
		}
		parsedOldUserID, err := uuid.Parse(old.UserID)
		if err != nil {
			return fmt.Errorf("failed to parse user UUID %s for asset %s: %w", old.UserID, old.ID, err)
		}
		parsedOldWorkID, err := parseNullableUUID(old.WorkID)
		if err != nil {
			return fmt.Errorf("failed to parse work UUID %s for asset %s: %w", *old.WorkID, old.ID, err)
		}

		newAsset := &dto.Asset{
			ID:        parsedOldID,
			WorkID:    parsedOldWorkID,
			AssetType: types.AssetType(old.AssetType),
			UserID:    parsedOldUserID,
			Extension: old.Extension,
			URL:       old.URL,
			CreatedAt: old.CreatedAt,
			UpdatedAt: old.UpdatedAt,
		}
		newAssets = append(newAssets, newAsset)
	}

	if len(newAssets) > 0 {
		if _, err := targetDB.NewInsert().Model(&newAssets).Exec(ctx); err != nil {
			return fmt.Errorf("failed to insert new assets: %w", err)
		}
	}

	fmt.Printf("-> Migrated %d assets.\n", len(newAssets))
	return nil
}
