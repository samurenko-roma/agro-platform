package productionunit

import (
	"context"
	"encoding/json"
	"fmt"

	pu "github.com/samurenkoroma/agro-platform/internal/domain/spatial/aggregate/production_unit"
)

func (r *productionUnitRepository) Save(ctx context.Context, unit *pu.ProductionUnit) error {
	propsJSON, err := json.Marshal(unit.Properties)
	if err != nil {
		return fmt.Errorf("failed to marshal properties: %w", err)
	}

	var geomArg any
	if unit.Geometry.Type != "" {
		geoJSON, err := unit.Geometry.ToGeoJSON()
		if err != nil {
			return fmt.Errorf("failed to marshal geometry: %w", err)
		}
		geomArg = string(geoJSON)
	}

	// area == 0 означает "домен не вычислил площадь ни из Dimensions, ни из
	// ручного ввода" — передаём NULL, чтобы сработал COALESCE на геометрию.
	var areaArg any
	if unit.Area != 0 {
		areaArg = unit.Area
	}

	query := `INSERT INTO 
    production_units(
                     id,owner_id,parent_id, code, sequence, area,
                     status,type,properties,geometry,
                     created_at,updated_at,archived_at
                     ) VALUES(
                         $1,$2,$3,$4,$5,
                         COALESCE($6, ST_Area(geography(ST_GeomFromGeoJSON($10)))),
                         $7,$8,$9,ST_GeomFromGeoJSON($10),
                         $11,$12,$13
                     )
				ON CONFLICT(id) 
				DO UPDATE SET
					parent_id=excluded.parent_id,
					code=excluded.code,
					sequence=excluded.sequence,
					area=COALESCE(excluded.area, ST_Area(geography(excluded.geometry))),
					status=excluded.status,
					type=excluded.type,
					properties=excluded.properties,
					geometry=excluded.geometry,
					updated_at=excluded.updated_at,
					archived_at=excluded.archived_at`

	_, err = r.db.Exec(ctx, query,
		unit.ID, unit.OwnerID, unit.ParentID, unit.Code, unit.Sequence, areaArg,
		unit.Status, unit.Type, propsJSON, geomArg,
		unit.CreatedAt, unit.UpdatedAt, unit.ArchivedAt,
	)

	return err
}
