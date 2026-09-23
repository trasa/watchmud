package rules

import "time"

type MudTime struct {
	// Mobiles determine their activity at this duration
	Mobile time.Duration `json:"mobile"`

	// Violence determines the basic speed of combat
	Violence time.Duration `json:"violence"`

	// Zone reset times are evaluated at this duration
	Zone time.Duration `json:"zone"`
}
