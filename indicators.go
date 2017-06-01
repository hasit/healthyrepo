package main

// Indicator holds information about an indicator support by HealthyRepo.
type Indicator struct {
	Name        string `json:"name" bson:"name"`
	Key         string `json:"key" bson:"key"`
	Description string `json:"description" bson:"description"`
	Active      bool   `json:"active" bson:"active"`
}

// Indicators holds a list of indicators.
type Indicators []Indicator
