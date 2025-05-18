package structs

type Persona struct {
	Name              string
	Active            bool
	UserID            int
	PersonaID         int
	DateCreated       string
	FirstLoggedInAt   string
	LastLoggedInAt    string
	Cash              float64
	Boost             float64
	IconIndex         int
	Level             int
	Motto             string
	PercentToLevel    float32
	Rating            float64
	Rep               float64
	RepAtCurrentLevel int
	Score             int
	NumCarSlots       int
	CurrentCarID      int
	Badges            []string
	Presence          int
}
