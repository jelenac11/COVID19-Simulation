package simulation

import (
	"covid19/config"
	"math"
	"math/rand"
)

type Cell struct {
	X              int
	Y              int
	Goal           [2]int
	Incubation     int
	Infected       int
	Duration       int
	Immunity       float64
	Hospitalized   bool
	Quarantined    bool
	Cured          bool
	Mobile         bool
	Dead           bool
	SurvivalChance float64
}

func (cell *Cell) moveCell() {
	up := cell.X - cell.Goal[0]
	right := cell.Y - cell.Goal[1]
	if math.Abs(float64(up)) <= 3 {
		cell.X = cell.Goal[0]
	} else {
		cell.X += Sgn(float64(up)) * (-3)
	}
	if math.Abs(float64(right)) <= 3 {
		cell.Y = cell.Goal[1]
	} else {
		cell.Y += Sgn(float64(right)) * (-3)
	}
}

func (cell *Cell) runInfection() {
	if cell.Infected == 1 {
		if cell.Incubation > 0 {
			cell.Incubation = cell.Incubation - 1
		} else {
			if cell.Duration > 0 {
				cell.Duration = cell.Duration - 1
			} else {
				if rand.Float64() < cell.SurvivalChance {
					cell.cureEvent()
				} else {
					cell.deathEvent()
				}
				cell.dehospitalizeEvent()
				cell.leaveQuarantineEvent()
			}
		}
	}
}

func (cell *Cell) cureEvent() {
	cell.Infected = config.NOT_INFECTED
	cell.Incubation = 0
	cell.Duration = 0
	cell.Immunity = config.Immunity
}

func (cell *Cell) hospitalizeEvent() {
	cell.Hospitalized = true
	cell.SurvivalChance = 0.9
}

func (cell *Cell) quarantinedEvent() {
	cell.Quarantined = true
}

func (cell *Cell) dehospitalizeEvent() {
	cell.Hospitalized = false
	cell.SurvivalChance = 1 - config.Mortality
}

func (cell *Cell) infectionEvent() {
	cell.Infected = config.INFECTED
	cell.Incubation = config.Incubation
	cell.Duration = config.Duration - config.Incubation
	config.InfectedCells++
}

func (cell *Cell) leaveQuarantineEvent() {
	cell.Quarantined = false
}

func (cell *Cell) deathEvent() {
	cell.Infected = config.NOT_INFECTED
	cell.Duration = 0
	cell.Dead = true
	cell.Mobile = false
	config.DeadCells++
}
