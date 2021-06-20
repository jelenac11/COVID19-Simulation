package simulation

import (
	"covid19/config"
	"covid19/util"
	"fmt"
	"math"
	"math/rand"
	"os"
	"strconv"
	"sync"
)

var firstMove = true

type Population struct {
	Mesh []Cell
	Dim  int
}

func createCell(x, y int) (c Cell) {
	c = Cell{
		X:              x,
		Y:              y,
		Goal:           [2]int{x, y},
		Incubation:     0,
		Infected:       0,
		Duration:       0,
		Immunity:       0.0,
		Hospitalized:   false,
		Quarantined:    false,
		Mobile:         false,
		Dead:           false,
		SurvivalChance: 1 - config.Mortality,
	}
	return
}

func (population *Population) CreatePopulation() {
	population.Mesh = make([]Cell, config.Dim*config.Dim)
	n := 0
	for i := 1; i <= config.Dim; i++ {
		for j := 1; j <= config.Dim; j++ {
			population.Mesh[n] = createCell(i, j)
			n++
		}
	}
}

func (population *Population) InfectPatientZero() {
	i := (population.Dim * (population.Dim / 2)) + (population.Dim / 2)
	population.Mesh[i].Infected = config.INFECTED
	population.Mesh[i].Duration = config.Duration
}

func (population *Population) UpdateSerial() {
	defer util.CalculateTime("Serial", nil, population.Dim)()
	newMash := make([]Cell, int(math.Pow(float64(population.Dim), 2)))
	for i := 0; i < population.Dim; i++ {
		for j := 0; j < population.Dim; j++ {
			population.moveEntity(newMash, i, j)
			population.updateCell(newMash, i, j)
		}
	}
	population.Mesh = newMash
}

func (population *Population) moveEntity(newMesh []Cell, i, j int) {
	cell := population.Mesh[i*population.Dim+j]

	if !cell.Mobile || cell.Quarantined || cell.Hospitalized {
		return
	}

	rng := rand.Float64()
	if (cell.Goal[0] == cell.X && cell.Goal[1] == cell.Y) || firstMove {
		if rng < config.CROUD && !config.ActiveDistancing {
			cell.Goal = [2]int{population.Dim / 2, population.Dim / 2}
		} else {
			goal := rand.Intn(population.Dim)
			if goal%population.Dim != 0 {
				cell.Goal = [2]int{goal/population.Dim + 1, goal%population.Dim + 1}
			} else {
				cell.Goal = [2]int{goal / population.Dim, population.Dim}
			}
		}
	}
	cell.moveCell()

	firstMove = false
}

func (population *Population) UpdateParallel(tasks int) {
	defer util.CalculateTime("Parallel", &tasks, population.Dim)()
	newMash := make([]Cell, int(math.Pow(float64(population.Dim), 2)))
	var waitgroup sync.WaitGroup
	taskSize := population.Dim / tasks
	for i := 0; i < tasks; i++ {
		waitgroup.Add(1)
		coef := population.Dim
		if i < tasks-1 {
			coef = (i + 1) * taskSize
		}
		go population.updateMatrix(&waitgroup, newMash, i*taskSize, 0, coef, population.Dim)
	}
	waitgroup.Wait()
	population.Mesh = newMash
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
	cell := population.Mesh[i*population.Dim+j]
	if cell.Dead {
		newMesh[i*population.Dim+j] = cell
		return
	}
	cell.runInfection()
	newMesh[i*population.Dim+j] = cell
	if cell.Infected == 0 || cell.Incubation > 0 {
		return
	}
	if !cell.Quarantined {
		neighbours := population.findNeighbours(cell.X, cell.Y)

		for _, neighbour := range neighbours {
			if rand.Float64() > population.Mesh[neighbour].Immunity {
				if rand.Float64() < config.Rate {
					population.Mesh[neighbour].infectionEvent()
					newMesh[neighbour] = population.Mesh[neighbour]
				}
			}
		}
	}
}

func (population *Population) findNeighbours(x, y int) []int {
	var neighbours []int
	for i := range population.Mesh {
		if population.Mesh[i].Infected == config.NOT_INFECTED && !population.Mesh[i].Dead {
			if checkPositions(population.Mesh[i].X, population.Mesh[i].Y, x, y) {
				neighbours = append(neighbours, i)
			}
		}
	}
	return neighbours
}

func (population *Population) RunTests() {
	for i := 0; i < population.Dim*3/2; i++ {
		idx := rand.Intn(len(population.Mesh))
		if population.Mesh[idx].Infected == config.INFECTED && !population.Mesh[idx].Hospitalized {
			config.TestedCells++
			rng := rand.Float64()
			if rng <= config.TEST_ACCURACY {
				config.TestedPositiveCells++
				if population.GetHospitalCount() == config.Capacity {
					population.Mesh[idx].hospitalizeEvent()
				} else {
					population.Mesh[idx].quarantinedEvent()
				}
			}
		}
	}
}

func (population *Population) GetHospitalCount() int {
	number := 0
	for i := range population.Mesh {
		if population.Mesh[i].Hospitalized {
			number++
		}
	}
	return number
}

func checkPositions(x1, y1, x, y int) bool {
	if (x1 == x && (y1 == (y-1) || y1 == (y+1))) || (y1 == y && (x1 == (x-1) || x1 == (x+1))) || (x1 == (x+1) &&
		(y1 == (y+1) || y1 == (y-1))) || (x1 == (x-1) && (y1 == (y+1) || y1 == (y-1))) {
		return true
	}
	return false
}

func Sgn(a float64) int {
	switch {
	case a < 0:
		return -1
	case a > 0:
		return +1
	}
	return 0
}

func (population *Population) PrintMesh() {
	forPrinting := make([]int, config.Dim*config.Dim)
	for i := range population.Mesh {
		if !population.Mesh[i].Dead {
			forPrinting[(population.Mesh[i].X-1)*population.Dim+population.Mesh[i].Y-1] = population.Mesh[i].Infected + 1
		} else {
			forPrinting[(population.Mesh[i].X-1)*population.Dim+population.Mesh[i].Y-1] = -1
		}
	}
	for i := 0; i < population.Dim; i++ {
		for j := 0; j < population.Dim; j++ {
			fmt.Print(forPrinting[i*population.Dim+j])
		}
		fmt.Println()
	}
	fmt.Println()
}

func (population *Population) Save(iter int, fileName string) {
	forPrinting := make([]int, config.Dim*config.Dim)
	file, err := os.Create(fmt.Sprintf("%s/%s%v.txt", config.PATH, fileName, iter))
	if err != nil {
		panic(err)
	}
	defer file.Close()
	for i := range population.Mesh {
		if !population.Mesh[i].Dead {
			forPrinting[(population.Mesh[i].X-1)*population.Dim+population.Mesh[i].Y-1] = population.Mesh[i].Infected + 1
		} else {
			forPrinting[(population.Mesh[i].X-1)*population.Dim+population.Mesh[i].Y-1] = -1
		}
	}
	for i := 0; i < population.Dim; i++ {
		for j := 0; j < population.Dim; j++ {
			if j != population.Dim - 1 {
				file.WriteString(strconv.Itoa(forPrinting[i*population.Dim+j]) + " ")
			} else {
				file.WriteString(strconv.Itoa(forPrinting[i*population.Dim+j]))
			}
		}
		file.WriteString("\n")
	}
}
