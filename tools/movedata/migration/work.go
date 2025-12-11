package migration

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/simesaba80/toybox-back/internal/infrastructure/database/dto"
	"github.com/simesaba80/toybox-back/internal/infrastructure/database/types"
	"github.com/uptrace/bun"
)

func MigrateWorks(ctx context.Context, sourceDB, targetDB bun.IDB) error {
	fmt.Println("Migrating works...")

	var oldWorks []OldWork
	if err := sourceDB.NewSelect().Model(&oldWorks).Scan(ctx); err != nil {
		return fmt.Errorf("failed to select old works: %w", err)
	}

	newWorks := make([]*dto.Work, 0, len(oldWorks))
	for _, old := range oldWorks {

		parsedOldID, err := uuid.Parse(old.ID)
		if err != nil {
			return fmt.Errorf("failed to parse work UUID %s: %w", old.ID, err)
		}

		parsedOldUserID, err := parseNullableUUID(old.UserID)
		if err != nil {
			return fmt.Errorf("failed to parse user UUID %s for work %s: %w", *old.UserID, old.ID, err)
		}

		newWork := &dto.Work{
			ID:          parsedOldID,
			Title:       old.Title,
			Description: old.Description,
			UserID:      parsedOldUserID,
			Visibility:  types.Visibility(old.Visibility),
			CreatedAt:   old.CreatedAt,
			UpdatedAt:   old.UpdatedAt,
		}
		newWorks = append(newWorks, newWork)
	}

	if len(newWorks) > 0 {
		if _, err := targetDB.NewInsert().Model(&newWorks).Exec(ctx); err != nil {
			return fmt.Errorf("failed to insert new works: %w", err)
		}
	}

	fmt.Printf("-> Migrated %d works.\n", len(newWorks))
	return nil
}
