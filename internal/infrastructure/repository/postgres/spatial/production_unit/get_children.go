package productionunit

import (
	"context"
	"encoding/json"

	vo "github.com/samurenkoroma/agro-platform/internal/domain/shared/valueobject"
	pu "github.com/samurenkoroma/agro-platform/internal/domain/spatial/aggregate/production_unit"
)

func (r *productionUnitRepository) GetChildren(ctx context.Context, parentID vo.ID) ([]*pu.ProductionUnit, error) {
	query := `SELECT id, owner_id, parent_id, type, status, code, area, properties, created_at, updated_at, archived_at
				FROM production_units
				WHERE parent_id=$1
				ORDER BY created_at`

	rows, err := r.db.Query(ctx, query, parentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make([]*pu.ProductionUnit, 0)

	for rows.Next() {
		var item pu.ProductionUnit
		var propertiesRaw []byte

		err = rows.Scan(
			&item.ID,
			&item.OwnerID,
			&item.ParentID,
			&item.Type,
			&item.Status,
			&item.Code,
			&item.Area,
			&propertiesRaw,
			&item.CreatedAt,
			&item.UpdatedAt,
			&item.ArchivedAt,
		)
		if err != nil {
			return nil, err
		}
		if propertiesRaw != nil {
			if err := json.Unmarshal(propertiesRaw, &item.Properties); err != nil {
				return nil, err
			}
		}

		result = append(result, &item)
	}

	return result, nil
}
