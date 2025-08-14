package model

type Readme struct {
	Name     string           `json:"name"`
	Desc     string           `json:"description"`
	Upgrades []*ReadmeUpgrade `json:"upgrades"`
}

type ReadmeUpgrade struct {
	Id      string `json:"id"`
	Version string `json:"version"`
	Created string `json:"created"`
	Notes   string `json:"notes"`
}

