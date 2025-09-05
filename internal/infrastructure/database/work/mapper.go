package work

import (
	"github.com/simesaba80/toybox-back/internal/domain/entity"
	"github.com/simesaba80/toybox-back/internal/infrastructure/database/asset"
	"github.com/simesaba80/toybox-back/internal/infrastructure/database/dto"
)

func ToEntity(d *dto.Work) *entity.Work {
	return &entity.Work{
		ID:              d.ID,
		Title:           d.Title,
		Description:     d.Description,
		DescriptionHTML: d.DescriptionHTML,
		UserID:          d.UserID,
		Visibility:      d.Visibility,
		Assets:          asset.ToEntities(d.Assets),
		CreatedAt:       d.CreatedAt,
		UpdatedAt:       d.UpdatedAt,
	}
}

func ToEntities(ds []*dto.Work) []*entity.Work {
	if ds == nil {
		return nil
	}
	es := make([]*entity.Work, len(ds))
	for i, d := range ds {
		es[i] = ToEntity(d)
	}
	return es
}

func ToDTO(e *entity.Work) *dto.Work {
	return &dto.Work{
		ID:              e.ID,
		Title:           e.Title,
		Description:     e.Description,
		DescriptionHTML: e.DescriptionHTML,
		UserID:          e.UserID,
		Visibility:      e.Visibility,
		Assets:          asset.ToDTOs(e.Assets),
		CreatedAt:       e.CreatedAt,
		UpdatedAt:       e.UpdatedAt,
	}
}