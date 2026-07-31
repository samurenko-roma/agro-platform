package productionunit

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/jackc/pgx/v5"
	vo "github.com/samurenkoroma/agro-platform/internal/domain/shared/valueobject"
	pu "github.com/samurenkoroma/agro-platform/internal/domain/spatial/aggregate/production_unit"
)

func (r *productionUnitRepository) GetByID(ctx context.Context, id vo.ID, orgId vo.ID) (*pu.ProductionUnit, error) {
	query := `SELECT id,owner_id,parent_id,type,status,code,area,ST_AsGeoJSON(geometry),properties,created_at,updated_at,archived_at
				FROM production_units
				WHERE id=$1 AND owner_id=$2`

	row := r.db.QueryRow(ctx, query, id, orgId)

	var root pu.ProductionUnit
	var propertiesRaw []byte
	var geometryRaw *string

	err := row.Scan(
		&root.ID,
		&root.OwnerID,
		&root.ParentID,
		&root.Type,
		&root.Status,
		&root.Code,
		&root.Area,
		&geometryRaw,
		&propertiesRaw,
		&root.CreatedAt,
		&root.UpdatedAt,
		&root.ArchivedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	if propertiesRaw != nil {
		if err := json.Unmarshal(propertiesRaw, &root.Properties); err != nil {
			return nil, err
		}
	}

	if geometryRaw != nil {
		geom, err := vo.GeometryFromGeoJSON([]byte(*geometryRaw))
		if err != nil {
			return nil, err
		}
		root.Geometry = geom
	}

	return &root, nil
}
