package migration

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/simesaba80/toybox-back/internal/infrastructure/database/dto"
	"github.com/simesaba80/toybox-back/internal/infrastructure/database/types"
	"github.com/uptrace/bun"
)

func MigrateURLInfos(ctx context.Context, sourceDB, targetDB bun.IDB) error {
	fmt.Println("Migrating url_infos...")
	var oldURLInfos []OldURLInfo
	if err := sourceDB.NewSelect().Model(&oldURLInfos).Scan(ctx); err != nil {
		return fmt.Errorf("failed to select old url_infos: %w", err)
	}

	newURLInfos := make([]*dto.URLInfo, 0, len(oldURLInfos))
	for _, old := range oldURLInfos {

		parsedOldID, err := uuid.Parse(old.ID)
		if err != nil {
			return fmt.Errorf("failed to parse url_info UUID %s: %w", old.ID, err)
		}
		parsedOldUserID, err := uuid.Parse(old.UserID)
		if err != nil {
			return fmt.Errorf("failed to parse user UUID %s for url_info %s: %w", old.UserID, old.ID, err)
		}
		parsedOldWorkID, err := parseNullableUUID(old.WorkID)
		if err != nil {
			return fmt.Errorf("failed to parse work UUID %s for url_info %s: %w", *old.WorkID, old.ID, err)
		}

		newURLInfo := &dto.URLInfo{
			ID:        parsedOldID,
			WorkID:    parsedOldWorkID,
			URL:       old.URL,
			URLType:   types.URLType(old.URLType),
			UserID:    parsedOldUserID,
			CreatedAt: old.CreatedAt,
			UpdatedAt: old.UpdatedAt,
		}
		newURLInfos = append(newURLInfos, newURLInfo)
	}

	if len(newURLInfos) > 0 {
		if _, err := targetDB.NewInsert().Model(&newURLInfos).Exec(ctx); err != nil {
			return fmt.Errorf("failed to insert new url_infos: %w", err)
		}
	}

	fmt.Printf("-> Migrated %d url_infos.\n", len(newURLInfos))
	return nil
}
