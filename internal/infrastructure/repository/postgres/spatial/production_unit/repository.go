package productionunit

import (
	"context"
	"encoding/json"

	"github.com/samurenkoroma/agro-platform/internal/application/uow"
	vo "github.com/samurenkoroma/agro-platform/internal/domain/shared/valueobject"
	pu "github.com/samurenkoroma/agro-platform/internal/domain/spatial/aggregate/production_unit"
	spatial "github.com/samurenkoroma/agro-platform/internal/domain/spatial/repository"
)

type productionUnitRepository struct {
	db uow.DB
}

func New(db uow.DB) spatial.ProductionUnitRepository {
	return &productionUnitRepository{db: db}
}

func (r *productionUnitRepository) GetNextSequence(ctx context.Context, orgID vo.ID, parentID *vo.ID, unitType pu.ProductionUnitType) (int, error) {
	sql := `SELECT COALESCE(MAX(sequence),0)+1 
				FROM production_units
			WHERE owner_id = $1
			  AND type = $2
			  AND parent_id IS NOT DISTINCT FROM $3
`

	var next int
	err := r.db.QueryRow(ctx, sql, orgID, unitType, parentID).Scan(&next)

	return next, err
}

func (r *productionUnitRepository) ListByOwner(ctx context.Context, ownerID vo.ID) ([]*pu.ProductionUnit, error) {
	query := `SELECT id, owner_id, parent_id, type, status, code, area, properties, created_at, updated_at, archived_at
				FROM production_units
				WHERE owner_id = $1
				ORDER BY code`

	rows, err := r.db.Query(ctx, query, ownerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make([]*pu.ProductionUnit, 0)
	for rows.Next() {
		var item pu.ProductionUnit
		var propertiesRaw []byte

		if err := rows.Scan(
			&item.ID, &item.OwnerID, &item.ParentID, &item.Type, &item.Status,
			&item.Code, &item.Area, &propertiesRaw,
			&item.CreatedAt, &item.UpdatedAt, &item.ArchivedAt,
		); err != nil {
			return nil, err
		}
		if propertiesRaw != nil {
			if err := json.Unmarshal(propertiesRaw, &item.Properties); err != nil {
				return nil, err
			}
		}
		result = append(result, &item)
	}

	return result, rows.Err()
}
