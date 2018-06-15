package main

import (
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

var (
	autos = 200.0
	timer = 3.0

	cmdName     = "cliclick"
	positionCmd = "p"
	moveCmd     = "m:"
	clickCmd    = "c:."
	colorCmd    = "cp:"
	easeCmd     = "-e 500"

	expected = []string{"240", "241", "242", "243", "244", "245", "246", "247",
		"248", "249", "250", "251", "252", "253", "254", "255"}

	clickArgs = []string{clickCmd}
	cmdOut    []byte
	err       error
)

var (
	crafting = false
)

func main() {
	forever := make(chan bool)
	go startRefresh()
	go startCrafting()
	<-forever
}

func startRefresh() {
	time.Sleep(10 * time.Second)
	fmt.Println("replenish starting...")
	for {
		if crafting {
			fmt.Println("Waiting 30 seconds for crafting")
			time.Sleep(30 * time.Second)
		}
		fmt.Println("checking for captcha...")
		out := output("581,374")
		valid := checkOutput(out)

		if !valid {
			motd := checkForMOTD()
			if !motd {
				fmt.Println("Failed", out)
				fmt.Println(time.Now())
				os.Exit(1)
			}
		}

		fmt.Println("click")
		click()
		calcWaitTime()
	}
}

func startCrafting() {
	fmt.Println("crafting started...")
	for {
		crafting = true
		fmt.Println("crafting")
		move("392,527")
		click()
		move("949,566")
		click()
		move("754,566")
		click()
		move("131,727")
		click()
		click()
		move("131,350")
		move("680,365")
		time.Sleep(5 * time.Second)
		fmt.Println("done")
		crafting = false
		time.Sleep(25 * time.Minute)
	}
}

func output(coords string) []string {
	args := []string{colorCmd + coords}
	cmdOut, err = exec.Command(cmdName, args...).Output()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error checking for captcha ", err)
		fmt.Println(time.Now())
		os.Exit(1)
	}
	output := strings.Split(string(cmdOut), " ")
	output[2] = strings.TrimSuffix(output[2], "\n\n")
	return output
}

func checkOutput(output []string) bool {
	for i := 0; i < 3; i++ {
		index := indexOf(output[i], expected)
		if index == -1 {
			return false
		}
	}
	return true
}

func move(position string) {
	moveArgs := []string{easeCmd, moveCmd + position}
	moveCmd := exec.Command(cmdName, moveArgs...)
	err = moveCmd.Start()
	if err != nil {
		fmt.Println(os.Stderr, "Error moving to original mouse position ", err.Error())
		fmt.Println(time.Now())
		os.Exit(1)
	}
	time.Sleep(3 * time.Second)
}

func click() {
	cmd := exec.Command(cmdName, clickArgs...)
	err = cmd.Start()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error starting Cmd", err.Error())
		fmt.Println(time.Now())
		os.Exit(1)
	}
	time.Sleep(3 * time.Second)
}

func calcWaitTime() {
	rand.Seed(time.Now().UnixNano())
	r := .9 + rand.Float64()/10
	waitTime := strconv.FormatFloat((((r * autos) * timer) * 1000), 'f', 6, 64)
	fmt.Println("waiting " + waitTime + "ms to click again")
	t, _ := time.ParseDuration(waitTime + "ms")
	time.Sleep(t)
}

func checkForMOTD() bool {
	output := output("1164,334")
	if (output[0] != "40") || (output[0] != "97") || (output[0] != "94") {
		if (output[0] == "243") && (output[0] == "156") && (output[0] == "57") {
			move("1164,334")
			click()
			move("685,369")
			return true
		}
	}
	return false
}

func indexOf(word string, data []string) int {
	for k, v := range data {
		if word == v {
			return k
		}
	}
	return -1
}
