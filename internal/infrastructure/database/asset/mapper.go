package asset

import (
	"github.com/simesaba80/toybox-back/internal/domain/entity"
	"github.com/simesaba80/toybox-back/internal/infrastructure/database/dto"
)

func ToEntity(d *dto.Asset) entity.Asset {
	return entity.Asset{
		ID:        d.ID,
		WorkID:    d.WorkID,
		AssetType: d.AssetType,
		UserID:    d.UserID,
		Extension: d.Extension,
		URL:       d.URL,
		CreatedAt: d.CreatedAt,
		UpdatedAt: d.UpdatedAt,
	}
}

func ToEntities(ds []*dto.Asset) []entity.Asset {
	if ds == nil {
		return nil
	}
	es := make([]entity.Asset, len(ds))
	for i, d := range ds {
		es[i] = ToEntity(d)
	}
	return es
}

func ToDTO(e *entity.Asset) *dto.Asset {
	return &dto.Asset{
		ID:        e.ID,
		WorkID:    e.WorkID,
		AssetType: e.AssetType,
		UserID:    e.UserID,
		Extension: e.Extension,
		URL:       e.URL,
		CreatedAt: e.CreatedAt,
		UpdatedAt: e.UpdatedAt,
	}
}

func ToDTOs(es []entity.Asset) []*dto.Asset {
	if es == nil {
		return nil
	}
	ds := make([]*dto.Asset, len(es))
	for i, e := range es {
		ds[i] = ToDTO(&e)
	}
	return ds
}