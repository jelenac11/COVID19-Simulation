package main

import (
	"fmt"
	"math"
	"math/rand"
	"os"
	"sync"
)

const INFECTED = 1
const NOT_INFECTED = 0

type Population struct {
	mesh    []Cell
	dim    int
}

type Cell struct {
	X           int
	Y           int
	Incubation  int
	Infected    int
	Duration    int
	Immunity    float64
	Hospitalized   bool
	Quarantined bool
	Cured bool
	Mobile bool
	Dead bool
	SurvivalChance float64
}

func createCell(x, y int) (c Cell) {
	c = Cell{
		X:              x,
		Y:              y,
		Incubation:     0,
		Infected:       0,
		Duration:       0,
		Immunity:       0.0,
		Hospitalized:   false,
		Quarantined:    false,
		Mobile:         false,
		Dead:           false,
		SurvivalChance: 1 - mortality,
	}
	return
}

func (population *Population) createPopulation() {
	population.mesh = make([]Cell, dim*dim)
	n := 0
	for i := 1; i <= dim; i++ {
		for j := 1; j <= dim; j++ {
			population.mesh[n] = createCell(i, j)
			n++
		}
	}
}

func (population *Population) infectPatientZero() {
	i := (population.dim * population.dim / 2) + (population.dim / 2)
	population.mesh[i].Infected = INFECTED
	population.mesh[i].Duration = duration
}

func (population *Population) PrintMesh() {
	for i := 0; i < population.dim; i++ {
		for j := 0; j < population.dim; j++ {
			if !population.mesh[i*population.dim+j].Dead {
				fmt.Print(population.mesh[i*population.dim+j].Infected)
			} else {
				fmt.Print("-1")
			}
		}
		fmt.Println()
	}
	fmt.Println()
}

func (population *Population) save(path string, iter int) {
	fileName := "serial"
	file, err := os.Create(fmt.Sprintf("%s/%s%v.mesh", path, fileName, iter))
	if err != nil {
		panic(err)
	}
	defer file.Close()
	for i := 0; i < population.dim; i++ {
		for j := 0; j < population.dim; j++ {
			file.WriteString(fmt.Sprint(population.mesh[i*population.dim+j], " "))
		}
		file.WriteString("\n")
	}
}

func (population *Population) UpdateSerial() {
	defer calculateTime("Serial", nil, population.dim)()
	newMash := make([]Cell, int(math.Pow(float64(population.dim), 2)))
	for i := 0; i < population.dim; i++ {
		for j := 0; j < population.dim; j++ {
			// ovde treba i population.moveEntity()
			population.updateCell(newMash, i, j)
		}
	}
	population.mesh = newMash
}

func (population *Population) UpdateParallel(tasks int) {
	defer calculateTime("Parallel", &tasks, population.dim)()
	newMash := make([]Cell, int(math.Pow(float64(population.dim), 2)))
	var waitgroup sync.WaitGroup
	taskSize := population.dim / tasks
	for i := 0; i < tasks; i++ {
		waitgroup.Add(1)
		coef := population.dim
		if i < tasks - 1 {
			coef = (i + 1) * taskSize
		}
		go population.updateMatrix(&waitgroup, newMash, i*taskSize, 0, coef, population.dim)
	}
	waitgroup.Wait()
	population.mesh = newMash
}

func (population *Population) updateMatrix(waitgroup *sync.WaitGroup, newMesh []Cell, from1, from2, to1, to2 int) {
	for i := from1; i < to1; i++ {
		for j := from2; j < to2; j++ {
			population.updateCell(newMesh, i, j)
		}
	}
	waitgroup.Done()
}

func (population *Population) updateCell(newMesh []Cell, i, j int) {
	cell := population.mesh[i*population.dim+j]
	newMesh[i * population.dim + j] = cell
	if cell.Dead {
		return
	}
	cell.runInfection()
	newMesh[i * population.dim + j] = cell
	if cell.Infected == 0 || cell.Incubation > 0 {
		return
	}
	if !cell.Quarantined {
		neighbours := population.findNeighbours(cell.X, cell.Y)

		for _, neighbour := range neighbours {
			if rand.Float64() > population.mesh[neighbour].Immunity {
				if rand.Float64() < rate {
					population.mesh[neighbour].infectionEvent()
					newMesh[neighbour] = population.mesh[neighbour]
				}
			}
		}
	}
}

func (population *Population) findNeighbours(x, y int) []int {
	var neighbours []int
	for i := range population.mesh {
		if population.mesh[i].Infected == NOT_INFECTED && !population.mesh[i].Dead {
			if checkPositions(population.mesh[i].X, population.mesh[i].Y, x, y){
				neighbours = append(neighbours, i)
			}
		}
	}
	return neighbours
}

func (population *Population) runTests() {
	for i := 0; i < population.dim * 3/2; i++ {
		idx := rand.Intn(len(population.mesh))
		if population.mesh[idx].Infected == INFECTED && !population.mesh[idx].Hospitalized {
			rng := rand.Float64()
			if rng <= 0.9 {
				if population.getHospitalCount() == capacity {
					population.mesh[idx].hospitalizeEvent()
				} else {
					population.mesh[idx].quarantinedEvent()
				}
			}
		}
	}
}

func (population *Population) getHospitalCount() int {
	number := 0
	for i := range population.mesh {
		if population.mesh[i].Hospitalized  {
			number++
		}
	}
	return number
}

func checkPositions (x1, y1, x, y int) bool {
	if (x1 == x && (y1 == (y - 1) || y1 == (y + 1))) || (y1 == y && (x1 == (x-1) || x1 == (x + 1))) || (x1 == (x+1) &&
		(y1 == (y+1) || y1 == (y - 1))) || (x1 == (x-1) && (y1 == (y+1) || y1 == (y - 1))){
		return true
	}
	return false
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
	cell.Infected = NOT_INFECTED
	cell.Incubation = 0
	cell.Duration = 0
	cell.Immunity = immunity
}

func (cell *Cell) hospitalizeEvent() {
	cell.Hospitalized = true
	cell.Mobile = false
	cell.SurvivalChance = 0.9
}

func (cell *Cell) quarantinedEvent() {
	cell.Quarantined = true
	cell.Mobile = false
}

func (cell *Cell) dehospitalizeEvent() {
	cell.Hospitalized = false
}

func (cell *Cell) infectionEvent() {
	cell.Infected = INFECTED
	cell.Incubation = incubation
	cell.Duration = duration
}

func (cell *Cell) leaveQuarantineEvent() {
	cell.Quarantined = false
}

func (cell *Cell) deathEvent() {
	cell.Infected = NOT_INFECTED
	cell.Duration = 0
	cell.Dead = true
	cell.Mobile = false
}