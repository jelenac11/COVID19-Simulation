package main

import (
	"covid19/config"
	"covid19/simulation"
	"covid19/util"
	"fmt"
	"math"
	"os"
	"strconv"
)

func checkArguments(lowerBound, upperBound, value float64, error string) {
	if value < lowerBound || value > upperBound {
		fmt.Println(error)
		os.Exit(1)
	}
}

func fillStats() {
	config.Infected = append(config.Infected, config.InfectedCells)
	config.Dead = append(config.Dead, config.DeadCells)
	config.Tested = append(config.Tested, config.TestedCells)
	config.TestedPositive = append(config.TestedPositive, config.TestedPositiveCells)
	config.InfectedCells = 0
	config.TestedPositiveCells = 0
	config.TestedCells = 0
	config.DeadCells = 0
}

func main() {
	var err1, err2, err3, err4, err5, err6, err7, err8 error

	if len(os.Args) < 8 {
		Scaling()
		return
	}

	config.Days, err1 = strconv.Atoi(os.Args[1])
	if err1 != nil {
		fmt.Println(err1)
		os.Exit(2)
	}
	checkArguments(0, math.MaxInt32, float64(config.Days), "Number of days must be positive integer")

	config.Dim, err2 = strconv.Atoi(os.Args[2])
	if err1 != nil {
		fmt.Println(err2)
		os.Exit(2)
	}
	checkArguments(0, math.MaxInt32, float64(config.Dim), "Dimension must be positive integer")

	config.Duration, err3 = strconv.Atoi(os.Args[3])
	if err2 != nil {
		fmt.Println(err3)
		os.Exit(2)
	}
	checkArguments(0, math.MaxInt32, float64(config.Duration), "Duration must be positive integer")

	config.Incubation, err4 = strconv.Atoi(os.Args[4])
	if err3 != nil {
		fmt.Println(err4)
		os.Exit(2)
	}
	checkArguments(0, math.MaxInt32, float64(config.Incubation), "Incubation must be positive integer")

	config.Rate, err5 = strconv.ParseFloat(os.Args[5], 64)
	if err4 != nil {
		fmt.Println(err5)
		os.Exit(2)
	}
	checkArguments(0, 1, config.Rate, "Rate must be number between 0 and 1!")

	config.Mortality, err6 = strconv.ParseFloat(os.Args[6], 64)
	if err5 != nil {
		fmt.Println(err6)
		os.Exit(2)
	}
	checkArguments(0, 1, config.Mortality, "Mortality must be number between 0 and 1!")

	config.Immunity, err7 = strconv.ParseFloat(os.Args[7], 64)
	if err6 != nil {
		fmt.Println(err7)
		os.Exit(2)
	}
	checkArguments(0, 1, config.Immunity, "Immunity must be number between 0 and 1!")

	tasks, err8 := strconv.Atoi(os.Args[8])
	if err8 != nil {
		fmt.Println(err7)
		os.Exit(2)
	}
	checkArguments(0, math.MaxInt32, float64(tasks), "Number of tasks must be number between 0 and 1!")

	config.ActiveDistancing = false

	population := simulation.Population{Dim: config.Dim}
	population.CreatePopulation()
	population.InfectPatientZero()
	config.Capacity = int(math.Pow(float64(config.Dim), 2) / 4)

	for i := 0; i < config.Days; i++ {
		population.PrintMesh()
		if population.GetHospitalCount() > int(float64(config.Capacity) * 0.8) {
			config.ActiveDistancing = true
			config.Quarantine = append(config.Quarantine, 1)
		} else {
			config.ActiveDistancing = false
			config.Quarantine =  append(config.Quarantine, 0)
		}
		if tasks == 0 {
			population.Save(i, "serial")
			population.UpdateSerial()
		} else {
			population.Save(i, "parallel")
			population.UpdateParallel(tasks)
		}
		population.RunTests()
		fillStats()
	}
	util.SaveStats()
}
func  StrongAmdahlScaling() {
	config.Dim = 1000
	config.Duration = 14
	config.Incubation = 5
	config.Rate = 0.8
	config.Mortality = 0.1
	config.Immunity = 0.7
	config.Capacity = int(math.Pow(float64(config.Dim), 2) / 4)

	threads := [4]int{1, 2, 4, 8}
	population := simulation.Population{Dim: 1000}
	population.CreatePopulation()
	population.InfectPatientZero()
	mesh := population.Mesh
	for _, s := range threads {
		population.UpdateParallel(s)
		population.RunTests()
		population.Mesh = mesh
	}
}

func WeakGustafsonScaling() {
	threads := [4]int{1, 2, 4, 8}
	for _, s := range threads {
		config.Dim = s * 1000
		config.Capacity = int(math.Pow(float64(config.Dim), 2) / 4)
		population := simulation.Population{Dim: s * 1000}
		population.CreatePopulation()
		population.InfectPatientZero()
		population.UpdateParallel(s)
		population.RunTests()
	}
}

func Scaling() {
	fmt.Println("Strong scaling (Amdahl):")
	StrongAmdahlScaling()
	fmt.Println("Weak scaling (Gustafson):")
	WeakGustafsonScaling()
}
