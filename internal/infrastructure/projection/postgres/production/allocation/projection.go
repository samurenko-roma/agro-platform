package allocation

import (
	"context"

	"github.com/samurenkoroma/agro-platform/internal/application/queries/production/allocation"
	"github.com/samurenkoroma/agro-platform/internal/application/uow"
	vo "github.com/samurenkoroma/agro-platform/internal/domain/shared/valueobject"
)

type projection struct {
	db uow.DB
}

func New(db uow.DB) allocation.Projection {
	return &projection{db: db}
}

func (p *projection) GetOccupancy(ctx context.Context, productionUnitID vo.ID) (*allocation.OccupancyDTO, error) {

	sql := `SELECT pu.id, pu.code,pu.area,COALESCE(SUM(a.area), 0) allocated_area
FROM production_units pu
LEFT JOIN production_allocations a
    ON a.production_unit_id = pu.id AND a.ended_at IS NULL
WHERE pu.id = $1
GROUP BY pu.id, pu.code,pu.area`

	var result allocation.OccupancyDTO
	var totalArea float64
	var allocatedArea float64

	err := p.db.QueryRow(ctx, sql, productionUnitID).Scan(
		&result.ProductionUnitID,
		&result.ProductionUnitName,
		&totalArea,
		&allocatedArea,
	)
	if err != nil {
		return nil, err
	}

	result.TotalArea = totalArea
	result.AllocatedArea = allocatedArea
	result.FreeArea = totalArea - allocatedArea

	return &result, nil
}

func (p *projection) ListByCycleID(ctx context.Context, cycleID vo.ID) ([]*allocation.DTO, error) {
	sql := `SELECT
    a.id,
    a.cycle_id,
    a.production_unit_id,
    pu.code,
    a.area,
    a.started_at,
    a.ended_at
FROM production_allocations a
LEFT JOIN production_units pu
    ON pu.id = a.production_unit_id
WHERE a.cycle_id = $1
ORDER BY a.started_at;`
	rows, err := p.db.Query(ctx, sql, cycleID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make([]*allocation.DTO, 0)
	for rows.Next() {
		item, err := scanDTO(rows)
		if err != nil {
			return nil, err
		}
		result = append(result, item)
	}
	return result, nil
}

func (p *projection) ListByProductionUnitID(ctx context.Context, productionUnitID vo.ID) ([]*allocation.DTO, error) {
	sql := `SELECT
    a.id,
    a.cycle_id,
    a.production_unit_id,
    pu.code,
    a.area,
    a.started_at,
    a.ended_at
FROM production_allocations a
LEFT JOIN production_units pu
    ON pu.id = a.production_unit_id
WHERE a.production_unit_id = $1
ORDER BY a.started_at;`
	rows, err := p.db.Query(ctx, sql, productionUnitID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make([]*allocation.DTO, 0)
	for rows.Next() {
		item, err := scanDTO(rows)
		if err != nil {
			return nil, err
		}
		result = append(result, item)
	}
	return result, nil
}

func (p *projection) Get(ctx context.Context, id vo.ID) (*allocation.DTO, error) {
	sql := `SELECT
    a.id,
    a.cycle_id,
    a.production_unit_id,
    unit.code,
    a.area,
    a.started_at,
    a.ended_at
FROM public.production_allocations a
         LEFT JOIN production_units unit ON unit.id = a.production_unit_id
WHERE a.id = $1`

	row := p.db.QueryRow(ctx, sql, id)
	return scanDTO(row)
}

func (p *projection) List(ctx context.Context, ownerId vo.ID) ([]*allocation.DTO, error) {
	sql := `SELECT
    a.id,
    a.cycle_id,
    a.production_unit_id,
    unit.code,
    a.area,
    a.started_at,
    a.ended_at
FROM public.production_allocations a
         LEFT JOIN production_units unit ON unit.id = a.production_unit_id
WHERE a.farm_id = $1
ORDER BY unit.code`
	rows, err := p.db.Query(ctx, sql, ownerId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make([]*allocation.DTO, 0)
	for rows.Next() {
		item, err := scanDTO(rows)
		if err != nil {
			return nil, err
		}
		result = append(result, item)
	}
	return result, nil
}

type scanner interface {
	Scan(dest ...any) error
}

func scanDTO(row scanner) (*allocation.DTO, error) {
	var result allocation.DTO

	if err := row.Scan(
		&result.ID, &result.CycleID, &result.ProductionUnitID, &result.ProductionUnitName,
		&result.Area, &result.StartedAt, &result.EndedAt,
	); err != nil {
		return nil, err
	}

	return &result, nil
}
