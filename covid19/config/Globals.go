package config

var (
	Days             int
	Duration         int
	Incubation       int
	Capacity         int
	Dim              int
	Mortality        float64
	Immunity         float64
	Rate             float64
	ActiveDistancing bool
	/*
		For stats: number of dead, infected, tested, and positive tested cells per day
	*/
	DeadCells           int
	InfectedCells       int
	TestedCells         int
	TestedPositiveCells int
	Infected            []int
	TestedPositive      []int
	Tested              []int
	Quarantine          []int
	Dead                []int
)
