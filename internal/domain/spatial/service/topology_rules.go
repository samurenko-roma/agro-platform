package service

import pu "github.com/samurenkoroma/agro-platform/internal/domain/spatial/aggregate/production_unit"

type TopologyRules interface {
	CanAttach(parent pu.ProductionUnitType, child pu.ProductionUnitType) bool
}

type DefaultTopology struct{}

func (DefaultTopology) CanAttach(parent pu.ProductionUnitType, child pu.ProductionUnitType) bool {
	allowed := map[pu.ProductionUnitType][]pu.ProductionUnitType{
		// Открытый грунт
		pu.Field: {pu.Plot, pu.Block},
		pu.Plot:  {pu.Block},
		pu.Block: {pu.Bed},
		pu.Bed:   {pu.Row},

		// Теплица
		pu.Greenhouse: {pu.GreenhouseZone},
		pu.GreenhouseZone: {
			pu.Bed, pu.Rack, pu.NFTChannel, pu.DWCTank, pu.AeroChamber, pu.Slot,
		},

		// Контейнерные / стеллажные
		pu.Rack:          {pu.Shelf},
		pu.VerticalTower: {pu.Shelf},
		pu.Shelf:         {pu.Slot, pu.Tray, pu.Pot},

		// Гидропоника
		pu.NFTChannel:  {pu.Slot},
		pu.DWCTank:     {pu.Slot},
		pu.AeroChamber: {pu.Slot},
	}

	for _, v := range allowed[parent] {
		if v == child {
			return true
		}
	}

	return false
}

// RootTypes — типы, которые ОСМЫСЛЕННО создавать без родителя (верхний
// уровень дерева). Используется как рекомендация, не как жёсткое
// ограничение (Create не запрещает root-узлы других типов, чтобы не
// блокировать служебные узлы вроде RESERVOIR/STORAGE/CONTAINER, у которых
// по домену нет естественного родителя).
var RootTypes = map[pu.ProductionUnitType]bool{
	pu.Field:         true,
	pu.Greenhouse:    true,
	pu.Rack:          true,
	pu.VerticalTower: true,
	pu.Reservoir:     true,
	pu.Storage:       true,
	pu.Container:     true,
}
