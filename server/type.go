package server

type Context struct {
	Title string
	Name  string
	Count int
}
type Status struct {
	PlaerName     string  `json:"player_name"`
	PlayerID      string  `json:"player"`
	LegalState    string  `json:"legal_state"`
	FuelMain      float64 `json:"fuel_main"`
	FuelReservior float64 `json:"fuel_reservoir"`
	CargoMass     int     `json:"cargo_mass"`
	System        int     `json:"Sys"`
	Weapon        int     `json:"Wep"`
	Engine        int     `json:"Eng"`
	Flags         int     `json:"flags"`
	Docked        bool    `json:"docked"`
	Landed        bool    `json:"landed"`
	LandGear      bool    `json:"land_gear"`
	Shields       bool    `json:"shields"`
	Supercruise   bool    `json:"supercruise"`
	FAOff         bool    `json:"flight_assist"`
	WeaponOn      bool    `json:"weapons"`
	CargoScope    bool    `json:"cargo_scoope"`
	InWing        bool    `json:"in_wing"`
	Lights        bool    `json:"lights"`
	Silents       bool    `json:"silents"`
	FuelScoope    bool    `json:"fuel_scope"`
	FSDMassLock   bool    `json:"mass_lock"`
	FSDCharge     bool    `json:"fsd_charge"`
	FSDCooldown   bool    `json:"fsd_cooldown"`
	LowFuel       bool    `json:"low_fuel"`
	OverHeating   bool    `json:"over_heating"`
	AnalysMode    bool    `json:"a_mode"`
	NVisin        bool    `json:"n_vision"`
	FSDJump       bool    `json:"fsd_jump"`
	IsInDanger    bool    `json:"is_in_danger"`
}

type Shutdown struct {
	PlayerID string `json:"player"`
	Delete   bool   `json:"delete"`
}
