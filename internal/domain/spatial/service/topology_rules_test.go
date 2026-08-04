package service_test

import (
	"testing"

	pu "github.com/samurenkoroma/agro-platform/internal/domain/spatial/aggregate/production_unit"
	"github.com/samurenkoroma/agro-platform/internal/domain/spatial/service"
)

func TestDefaultTopology_CanAttach_AllowedPairs(t *testing.T) {
	topo := service.DefaultTopology{}
	allowed := []struct{ parent, child pu.ProductionUnitType }{
		{pu.Field, pu.Block},
		{pu.Field, pu.Plot},
		{pu.Plot, pu.Block},
		{pu.Block, pu.Bed},
		{pu.Bed, pu.Row},
		{pu.Greenhouse, pu.GreenhouseZone},
		{pu.GreenhouseZone, pu.Bed},
		{pu.GreenhouseZone, pu.NFTChannel},
		{pu.Rack, pu.Shelf},
		{pu.Shelf, pu.Slot},
		{pu.NFTChannel, pu.Slot},
	}
	for _, c := range allowed {
		if !topo.CanAttach(c.parent, c.child) {
			t.Errorf("expected %s -> %s to be allowed", c.parent, c.child)
		}
	}
}

func TestDefaultTopology_CanAttach_DisallowedPairs(t *testing.T) {
	topo := service.DefaultTopology{}
	disallowed := []struct{ parent, child pu.ProductionUnitType }{
		{pu.Field, pu.Slot},
		{pu.Bed, pu.Field},
		{pu.Shelf, pu.Field},
		{pu.Reservoir, pu.Bed},
	}
	for _, c := range disallowed {
		if topo.CanAttach(c.parent, c.child) {
			t.Errorf("expected %s -> %s to be disallowed", c.parent, c.child)
		}
	}
}
