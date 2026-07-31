package valueobject

type Coordinates struct {
	X float64  `json:"x"`
	Y float64  `json:"y"`
	Z *float64 `json:"z,omitempty"`
}

func NewCoordinates(x, y float64, z *float64) Coordinates {
	return Coordinates{X: x, Y: y, Z: z}
}
