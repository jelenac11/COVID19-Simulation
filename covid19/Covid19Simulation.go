package main

import (
	"fmt"
	"math"
	"os"
	"time"
)

var dim, duration, incubation int
var rate, mortality, immunity float64
var path string
var dead, infections []int
var capacity int

func checkArguments(lowerBound, upperBound, value float64, error string) {
	if value < lowerBound || value > upperBound {
		fmt.Println(error)
		os.Exit(1)
	}
}

func calculateTime(version string, tasks *int, dim int) func() {
	start := time.Now()
	return func() {
		str := ""
		if tasks != nil {
			str = fmt.Sprintf(" with %v threads ", *tasks)
		}
		fmt.Printf("%s took %v%s for dimensions: %vx%v\n", version, time.Since(start), str, dim, dim)
	}
}

func main() {
	// simulacija epidemije za 30 dana
	// uspesnost testa je 90%
	// kapacitet bolnice je 1/3 ukupnog broja jedinki
	// broj jedinki je broj celija u gridu
	// krece se 1/2 jedinki, dok ostatak miruje
	// karantin se uvodi kada broj zarazenih bude odredjeni procenat ili odmah to utvriditi
	// da li je parametre bolje unositi preko komandne ili preko konzole
	// kako cu ja u fajlu predstaviti ako su dve jedinke na istoj poziciji
	// posto se sporo siri zaraza mozda dodati vise nultih pacijenata
	//var err1, err2, err3, err4, err5, err6 error

	/*if len(os.Args) == 6 {
		//Scaling()
		return
	}*/

	//path := os.Args[1]
	/*
	dim, err1 = strconv.Atoi(os.Args[2])
	if err1 != nil {
		fmt.Println(err1)
		os.Exit(2)
	}

	duration, err2 = strconv.Atoi(os.Args[3])
	if err2 != nil {
		fmt.Println(err2)
		os.Exit(2)
	}
	if duration <= 0 {
		fmt.Println("Duration must be positive integer!")
		os.Exit(1)
	}

	incubation, err3 = strconv.Atoi(os.Args[4])
	if err3 != nil {
		fmt.Println(err3)
		os.Exit(2)
	}
	if incubation <= 0 {
		fmt.Println("Incubation must be positive integer!")
		os.Exit(1)
	}

	rate, err4 = strconv.ParseFloat(os.Args[5], 64)
	if err4 != nil {
		fmt.Println(err4)
		os.Exit(2)
	}
	checkArguments(0, 1, rate, "Rate must be number between 0 and 1!")

	mortality, err5 = strconv.ParseFloat(os.Args[6], 64)
	if err5 != nil {
		fmt.Println(err5)
		os.Exit(2)
	}
	checkArguments(0, 1, mortality, "Mortality must be number between 0 and 1!")

	immunity, err6 = strconv.ParseFloat(os.Args[7], 64)
	if err6 != nil {
		fmt.Println(err6)
		os.Exit(2)
	}
	checkArguments(0, 1, immunity, "Immunity must be number between 0 and 1!")
	*/

	dim = 100
	duration = 5
	incubation = 2
	rate = 0.8
	mortality = 0.2
	immunity = 0.6
	tasks := 4

	population := Population{dim: dim}
	population.createPopulation()
	population.infectPatientZero()

	capacity =int(math.Pow(float64(dim), 2) / 4)
	for i:= 0; i < 30; i++ {
		population.PrintMesh()
		/*if path != "-" {
			population.save(path, i)
		}*/
		if tasks == 0 {
			population.UpdateSerial()
		} else {
			population.UpdateParallel(tasks)
		}
		population.runTests()
		// ovde treba da ide statistika

	}
}

func StrongAmdahlScaling() {
	dim = 1000
	duration = 5
	incubation = 2
	rate = 0.8
	mortality = 0.2
	immunity = 0.6
	capacity =int(math.Pow(float64(dim), 2) / 4)

	threads := [4]int{1, 2, 4, 8}
	population := Population{dim: 1000}
	population.createPopulation()
	population.infectPatientZero()
	mesh := population.mesh
	for _, s := range threads {
		population.UpdateParallel(s)
		population.runTests()
		population.mesh = mesh
	}
}

func WeakGustafsonScaling() {
	threads := [4]int{1, 2, 4, 8}
	for _, s := range threads {
		dim = s * 1000
		capacity =int(math.Pow(float64(dim), 2) / 4)
		population := Population{dim: s * 1000}
		population.createPopulation()
		population.infectPatientZero()
		population.UpdateParallel(s)
		population.runTests()
	}
}

func Scaling() {
	fmt.Println("Strong scaling (Amdahl):")
	StrongAmdahlScaling()
	fmt.Println("Weak scaling (Gustafson):")
	WeakGustafsonScaling()
}