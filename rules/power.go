package rules

// PowerBand is the range of power a zone is built for -- how content says
// "this is the 10-15 area". Zero value is the bottom band. See LEVELS.md.
type PowerBand struct {
	Min int `json:"min"`
	Max int `json:"max"`
}
