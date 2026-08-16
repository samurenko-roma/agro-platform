package variety

import (
	"context"
)

type ListFilter struct {
	CropID string
}

type Detail struct {
	ID               string           `json:"id"`
	Name             string           `json:"name"`
	SpeciesKey       string           `json:"speciesKey"`
	SpeciesName      string           `json:"speciesName"`
	BaseTemperature  float64          `json:"baseTemperature"`
	MaxTemperature   float64          `json:"maxTemperature"`
	DaysToMaturity   int              `json:"daysToMaturity"`
	PhenophaseGDD    []Phenophase     `json:"phenophases"`
	WaterRequirement WaterRequirement `json:"waterRequirement"`
	LightRequirement LightRequirement `json:"lightRequirement"`
	YieldPotential   string           `json:"yieldPotential"`
	GrowingTypes     []string         `json:"growingTypes"`
	Description      string           `json:"description"`
	Image            string           `json:"image"`
}
type Phenophase struct {
	Code        string `json:"code"`
	Name        string `json:"name"`
	GddRequired int    `json:"gddRequired"`
	Description string `json:"description"`
	IsCritical  bool   `json:"isCritical"`
}
type WaterRequirement struct {
	DailyNeedMin   int      `json:"dailyNeedMin"`
	DailyNeedOpt   int      `json:"dailyNeedOpt"`
	CriticalPhases []string `json:"criticalPhases"`
}
type LightRequirement struct {
	PpfdMin         int      `json:"ppfdMin"`
	PpfdOpt         int      `json:"ppfdOpt"`
	DayLengthMin    int      `json:"dayLengthMin"`
	DayLengthOpt    int      `json:"dayLengthOpt"`
	PhotoperiodType string   `json:"photoperiodType"`
	CriticalPhases  []string `json:"criticalPhases"`
}

type ListItem struct {
	ID            string `json:"id"`
	Name          string `json:"name"`
	CropId        string `json:"cropId"`
	CropName      string `json:"cropName"`
	DaysToHarvest int    `json:"daysToMaturity"`
}

type Projection interface {
	Get(context.Context, string) (*Detail, error)
	List(context.Context, ListFilter) ([]ListItem, error)
}
