package growingcycle

import (
	"context"

	growingcycle "github.com/samurenkoroma/agro-platform/internal/application/queries/production/growing_cycle"
	"github.com/samurenkoroma/agro-platform/internal/application/shared"
	"github.com/samurenkoroma/agro-platform/internal/application/uow"
	vo "github.com/samurenkoroma/agro-platform/internal/domain/shared/valueobject"
)

type projection struct {
	db uow.DB
}

func New(db uow.DB) growingcycle.Projection {
	return &projection{db: db}
}

func (p *projection) List(ctx context.Context, filter growingcycle.FilterCycle) ([]*growingcycle.DTO, error) {

	sql := `
SELECT
    crop.id,
    crop.name,
       COALESCE(SUM(a.area), 0) allocated_area,
       NULL::integer            tasks_count,
       0                        progress,
       count(crop.id)           count
FROM production_growing_cycles cycle
         INNER JOIN agronomy_crops crop ON crop.id = cycle.crop_id
         LEFT JOIN production_allocations a ON a.cycle_id = cycle.id
WHERE cycle.farm_id = $1

GROUP BY crop.name, crop.id
ORDER BY crop.name DESC;
`

	rows, err := p.db.Query(ctx, sql, filter.OwnerId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make([]*growingcycle.DTO, 0)
	cropIDs := make([]vo.ID, 0)
	index := make(map[vo.ID]*growingcycle.DTO)

	for rows.Next() {
		item := &growingcycle.DTO{}
		var cropId vo.ID
		if err := rows.Scan(
			&cropId,
			&item.CropName,
			&item.AllocatedArea,
			&item.TasksCount,
			&item.Progress,
			&item.Count,
		); err != nil {
			return nil, err
		}

		item.Allocations = make([]growingcycle.AllocationDTO, 0)

		result = append(result, item)
		cropIDs = append(cropIDs, cropId)
		index[cropId] = item
	}
	if len(cropIDs) == 0 {
		return result, nil
	}

	sql = `
SELECT c.id,
       c.name,
       crop.id,
       v.name,
       c.status,
       c.stage,
       a.production_unit_id,
       pu.code,
       a.area,
       (CURRENT_DATE - a.started_at::date)::int growing_days,
       v.profile#>'{Maturity,DaysToHarvest}'  varietyDays,
       crop.agronomy#>'{Maturity,DaysToHarvest}'  cropDays,
       a.started_at,
       a.ended_at
FROM production_allocations a
         LEFT JOIN production_units pu ON pu.id = a.production_unit_id
         RIGHT JOIN production_growing_cycles c on c.id = a.cycle_id
         LEFT JOIN agronomy_crops crop on crop.id = c.crop_id
         LEFT JOIN agronomy_varieties v on v.id = c.variety_id
WHERE c.crop_id = ANY($1)
ORDER BY pu.code
`

	rows2, err := p.db.Query(ctx, sql, cropIDs)
	if err != nil {
		return nil, err
	}
	defer rows2.Close()

	for rows2.Next() {
		var cropId vo.ID
		var growingDays, varietyDays, cropDays *int
		var allocation growingcycle.AllocationDTO

		if err := rows2.Scan(
			&allocation.CycleId,
			&allocation.CycleName,
			&cropId,
			&allocation.VarietyName,
			&allocation.Status,
			&allocation.Stage,
			&allocation.ProductionUnitId,
			&allocation.ProductionUnitName,
			&allocation.Area,
			&growingDays, &varietyDays, &cropDays,
			&allocation.StartDate,
			&allocation.EndDate,
		); err != nil {
			return nil, err
		}

		dto, ok := index[cropId]
		if !ok {
			continue
		}

		// Цикл без единой аллокации (обычная ситуация для только что
		// созданного цикла) — a.* колонки NULL из-за RIGHT JOIN. Это не
		// ошибка, просто нечего показать по прогрессу для этой записи.
		if allocation.ProductionUnitId == nil {
			dto.Allocations = append(dto.Allocations, allocation)
			continue
		}

		if daysToMaturity := shared.Override[int](cropDays, varietyDays); daysToMaturity != nil && *daysToMaturity > 0 && growingDays != nil {
			allocation.DaysToMaturity = *daysToMaturity
			allocation.Progress = int(float64(*growingDays) / float64(*daysToMaturity) * 100)
		}
		// Иначе DaysToMaturity/Progress остаются нулевыми: нет данных о
		// зрелости ни у сорта, ни у культуры — честное "неизвестно" вместо
		// паники или мусорного значения от деления на ноль/NaN.

		dto.Progress += allocation.Progress
		dto.Allocations = append(dto.Allocations, allocation)
	}

	for _, dto := range result {
		if len(dto.Allocations) > 0 {
			dto.Progress = dto.Progress / len(dto.Allocations)
		}
	}

	return result, nil
}

func (p *projection) Summary(ctx context.Context, ownerId vo.ID, cycleId vo.ID) (*growingcycle.SummaryDTO, error) {
	sql := `SELECT
    cycle.id,
    cycle.name,
    cycle.status,

    COALESCE((
        SELECT SUM(area)
        FROM production_allocations
        WHERE cycle_id = cycle.id
          AND ended_at IS NULL
    ),0) allocated_area,

    COALESCE((
        SELECT SUM(quantity)
        FROM production_plantings
        WHERE cycle_id = cycle.id
    ),0) planted_quantity,

    COALESCE((
        SELECT SUM(quantity)
        FROM public.production_harvest_batch
        WHERE cycle_id = cycle.id
    ),0) harvested_quantity
FROM production_growing_cycles cycle WHERE  cycle.farm_id = $1 and cycle.id = $2`
	row := p.db.QueryRow(ctx, sql, ownerId, cycleId)

	var sum growingcycle.SummaryDTO
	if err := row.Scan(&sum.ID, &sum.Name, &sum.Status, &sum.AllocatedArea, &sum.PlantedQuantity, &sum.HarvestedQuantity); err != nil {
		return nil, err
	}

	return &sum, nil
}

type scanner interface {
	Scan(dest ...any) error
}

func scanDTO(row scanner) (*growingcycle.DTO, error) {
	var item growingcycle.DTO

	if err := row.Scan(
		&item.CropName,
		&item.AllocatedArea,
		&item.TasksCount,
		&item.Progress,
		&item.Count,
	); err != nil {
		return nil, err
	}

	return &item, nil
}
