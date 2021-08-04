package server

type Context struct {
	Title string
	Name  string
	Count int
}
type Status struct {
	PlayerID string `json:"player"`
	System   int    `json:"Sys"`
	Weapon   int    `json:"Wep"`
	Engine   int    `json:"Eng"`
}
