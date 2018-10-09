package positions

import (
	"github.com/bernardjkim/ptrade-api/pkg/db"
)

// Positions is a list of Position objects
type Positions []Position

// Position represents a users position for a specific stock
type Position db.PositionTable

// TableName simply returns the table name
func (h *Position) TableName() string {
	return "positions"
}
