package variety

import (
	"context"
	"encoding/json"

	"github.com/samurenkoroma/agro-platform/internal/application/queries/agronomy/variety"
	"github.com/samurenkoroma/agro-platform/internal/application/uow"
)

type projection struct {
	db uow.DB
}

func New(db uow.DB) variety.Projection {
	return &projection{
		db: db,
	}
}

// storedProfile — форма JSON, которую реально пишет Variety.Save()
// (Profile{Maturity, Spacing} с тегами json:"maturity"/"spacing").
type storedProfile struct {
	Maturity struct {
		DaysToHarvest *int     `json:"daysToHarvest"`
		GDDToHarvest  *float64 `json:"gddToHarvest"`
	} `json:"maturity"`
	Spacing struct {
		PlantDistanceCM      *float64 `json:"plantDistanceCm"`
		RowDistanceCM        *float64 `json:"rowDistanceCm"`
		PlantsPerSquareMeter *float64 `json:"plantsPerSquareMeter"`
		RecommendedDensity   *float64 `json:"recommendedDensity"`
	} `json:"spacing"`
}

func (p projection) Get(ctx context.Context, id string) (*variety.Detail, error) {
	// Вся модель, которая реально есть у сорта в БД: собственные поля
	// agronomy_varieties + имя культуры через JOIN (SpeciesName).
	query := `
SELECT
	v.id,
	v.crop_id,
	v.name,
	v.breeder,
	v.image,
	v.profile,
	c.name AS species_name
FROM agronomy_varieties v
LEFT JOIN agronomy_crops c ON c.id = v.crop_id
WHERE v.id = $1
`

	var (
		result      variety.Detail
		cropID      string
		breeder     *string
		image       *string
		profileRaw  []byte
		speciesName *string
	)

	err := p.db.QueryRow(ctx, query, id).Scan(
		&result.ID,
		&cropID,
		&result.Name,
		&breeder,
		&image,
		&profileRaw,
		&speciesName,
	)
	if err != nil {
		return nil, err
	}

	result.SpeciesKey = cropID
	if speciesName != nil {
		result.SpeciesName = *speciesName
	}
	if image != nil {
		result.Image = *image
	}

	if profileRaw != nil {
		var profile storedProfile
		if err := json.Unmarshal(profileRaw, &profile); err != nil {
			return nil, err
		}
		if profile.Maturity.DaysToHarvest != nil {
			result.DaysToMaturity = *profile.Maturity.DaysToHarvest
		}
	}
	result.GrowingTypes = []string{}
	result.LightRequirement.CriticalPhases = []string{}
	result.WaterRequirement.CriticalPhases = []string{}

	// Намеренно НЕ заполняются (нет источника данных в Variety вообще —
	// эти поля по смыслу принадлежат Crop.AgronomyProfile, а не сорту):
	//   BaseTemperature, MaxTemperature, PhenophaseGDD,
	//   WaterRequirement, LightRequirement, YieldPotential,
	//   GrowingTypes, Description.
	// breeder (v.breeder) — прочитан из БД, но в variety.Detail сейчас
	// нет поля под него; если понадобится — надо сначала добавить
	// Breeder в сам Detail.

	return &result, nil
}

func (p projection) List(ctx context.Context, filter variety.ListFilter) ([]variety.ListItem, error) {

	query := `SELECT v.id, v.name, c.id, c.name,
       COALESCE((v.profile#>>'{maturity,daysToHarvest}')::int, 0) as maturity
FROM agronomy_varieties v 
    LEFT JOIN  
    agronomy_crops c on v.crop_id = c.id
WHERE v.crop_id = $1
ORDER BY v.name`
	rows, err := p.db.Query(ctx, query, filter.CropID)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	result := make([]variety.ListItem, 0)

	for rows.Next() {

		var item variety.ListItem

		err = rows.Scan(&item.ID, &item.Name, &item.CropId, &item.CropName, &item.DaysToHarvest)

		if err != nil {
			return nil, err
		}

		result = append(result, item)
	}

	return result, nil
}
