package util

import (
	"covid19/config"
	"fmt"
	"os"
	"strconv"
	"time"
)

func CalculateTime(version string, tasks *int, dim int) func() {
	start := time.Now()
	return func() {
		str := ""
		if tasks != nil {
			str = fmt.Sprintf(" with %v threads ", *tasks)
		}
		fmt.Printf("%s took %v%s for dimensions: %vx%v\n", version, time.Since(start), str, dim, dim)
	}
}

func SaveStats() {
	toWrite := ""
	for i := 0; i < config.Days; i++ {
		data := strconv.Itoa(i+1)+ " " + strconv.Itoa(config.Infected[i]) + " " + strconv.Itoa(config.Dead[i]) + " " + strconv.Itoa(config.Quarantine[i]) +
			" " + strconv.Itoa(config.Tested[i]) + " " + strconv.Itoa(config.TestedPositive[i])
		if i == 0 {
			toWrite = data
		} else {
			toWrite = toWrite + "|" + data
		}
	}

	f, err := os.Create(config.STATS_PATH)
	if err != nil {
		fmt.Println(err)
		return
	}
	l, err := f.WriteString(toWrite)
	if err != nil {
		fmt.Println(err)
		f.Close()
		return
	}
	fmt.Println(l, "bytes written successfully")
	err = f.Close()
	if err != nil {
		fmt.Println(err)
		return
	}
}
