package model

import "time"

const ContinentPath string = "continent"
const CountryPath string = "country"
const CityPath string = "city"

type Continents struct {
	Continents []continent `json:"continents"`
}

type continent struct {
	ID         uint8     `json:"id" db:"id"`
	Area       float32   `json:"area"`
	Name       string    `json:"name" db:"pk"`
	Population uint64    `json:"population"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
	Country    Countries `json:"-" db:"fk" fk:"ContinentID" refer:"ID"`
}

type Countries struct {
	Countries []country `json:"countries"`
}

type country struct {
	ID          uint8     `json:"id" db:"id"`
	ContinentID uint8     `json:"continent_id" db:"fk"`
	Area        float32   `json:"area"`
	Name        string    `json:"name" db:"pk"`
	Population  uint64    `json:"population"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	City        Cities    `json:"-" db:"fk" fk:"CountryID" refer:"ID"`
}

type Cities struct {
	Cities []city `json:"cities"`
}

type city struct {
	ID         uint8     `json:"id" db:"id"`
	CountryID  uint8     `json:"country_id" db:"fk"`
	Area       float32   `json:"area"`
	Name       string    `json:"name" db:"pk"`
	Population uint64    `json:"population"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}
