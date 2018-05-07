package model

import (
	"fmt"
)

// Building is the type containing information about a building
type Building struct {
	Street      string
	Building    string
	URL         string
	IsBrick     bool
	IsApartment bool
}

func (b Building) String() string {
	return fmt.Sprintf("Street: %s; Building: %s; URL: %s", b.Street, b.Building, b.URL)
}
