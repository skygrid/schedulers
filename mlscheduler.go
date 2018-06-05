package schedulers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"sync"
)

// ML Scheduler
type MlScheduler struct {
	GeneralScheduler
}

func (mlsc *MlScheduler) Schedule(jobs []ResourceVolume, workers []ResourceVolume) []Decision {
	var d []Decision
	var Jindecies map[string]uint64
	Jindecies = make(map[string]uint64)
	var Windecies map[int]uint64
	Windecies = make(map[int]uint64)
	var maxcores uint32
	if len(workers) > 0 {
		maxcores = workers[0].GetCPU()
	}
	var workerscopy []ResourceVolume
	for i, w := range workers {
		workerscopy = append(workerscopy, w)
		Windecies[i] = w.Id
	}
	var intasks []uint64
	for _, j := range jobs {
		intasks = append(intasks, j.Param)
		Jindecies[strconv.Itoa(int(j.Param))] = j.Id
	}
	_ = os.Mkdir("temp", os.ModePerm)
	names := []string{"temp/intasks1.json", "temp/machines1.json", "temp/maxcores1.json"}

	mlsc.DumpData(intasks, len(workerscopy), maxcores, names)
	mlsc.Packing()
	ans := mlsc.LoadRes("temp/out1.json")
	for ix, dict1 := range ans {
		for key, value := range dict1 {
			d = append(d, Decision{CoresNum: uint32(value), JobIdx: Jindecies[key], WorkerIdx: Windecies[ix]})
		}
	}
	return d
}

func (mlsc *MlScheduler) DumpData(intasks []uint64, machines int, maxcores uint32, names []string) {
	in1, err := json.Marshal(intasks)
	if err != nil {
		fmt.Printf("%s", err)
	}
	err = ioutil.WriteFile(names[0], in1, 0644)
	if err != nil {
		fmt.Printf("%s", err)
	}
	in3, err := json.Marshal(machines)
	if err != nil {
		fmt.Printf("%s", err)
	}
	err = ioutil.WriteFile(names[1], in3, 0644)
	if err != nil {
		fmt.Printf("%s", err)
	}
	in4, err := json.Marshal(maxcores)
	if err != nil {
		fmt.Printf("%s", err)
	}
	err = ioutil.WriteFile(names[2], in4, 0644)
	if err != nil {
		fmt.Printf("%s", err)
	}
}

func (mlsc *MlScheduler) Packing() {
	wg := new(sync.WaitGroup)
	wg.Add(1)
	exeCmd("python3 packing.py", wg)
}

func exeCmd(cmd string, wg *sync.WaitGroup) string {
	//fmt.Println("command is ", cmd)
	// splitting head => g++ parts => rest of the command
	parts := strings.Fields(cmd)
	head := parts[0]
	parts = parts[1:]

	out, err := exec.Command(head, parts...).Output()
	if err != nil {
		fmt.Printf("%s", err)
	}
	fmt.Printf("%s", out)
	wg.Done()
	return string(out)
	// Need to signal to waitgroup that this goroutine is done
}

func (mlsc *MlScheduler) LoadRes(name string) []map[string]int {
	raw, err := ioutil.ReadFile(name)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	var c []map[string]int
	json.Unmarshal(raw, &c)
	return c
}
