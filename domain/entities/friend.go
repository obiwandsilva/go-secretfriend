package entities

// Friend ...
type Friend struct {
	Name string `json:"name"`
	PhoneNumber string `json:"phoneNumber"`
}

// Pair ...
type Pairs map[Friend]Friend
