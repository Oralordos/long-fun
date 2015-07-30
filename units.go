package game

type point struct {
	X int
	Y int
}

type unit struct {
	TileID   int
	Position point
}

func newUnit(tileID, x, y int) *unit {
	return &unit{
		TileID: tileID,
		Position: point{
			X: x,
			Y: y,
		},
	}
}

func (u *unit) confirmPath(p []point, m *gameMap) bool {
	// TODO Confirm that the input path is valid
	return false
}
