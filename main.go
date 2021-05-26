package main

import (
	"embed"
	"fmt"
	"log"
	"math/rand"
	"strings"
	"time"
)

//go:embed assets/hooks
var myHooks embed.FS

//go:embed assets/units
var myUnits embed.FS

var Reset = "\033[0m"
var Red = "\033[31m"
var Green = "\033[32m"
var Yellow = "\033[33m"
var Blue = "\033[34m"
var Purple = "\033[35m"
var Cyan = "\033[36m"
var Gray = "\033[37m"
var White = "\033[97m"
var Bold = "\033[1m"

type secondParts struct {
	t500, t1000, t200, t2, t10 time.Duration
}

var mT = secondParts{
	t500:  500 * time.Millisecond,
	t1000: 1000 * time.Millisecond,
	t200:  200 * time.Millisecond,
	t10:   10 * time.Millisecond,
	t2:    2 * time.Millisecond,
}

func clearTTY() {
	println("\033[;H\033[2J")
}

func earlyBoot(startHooks []string) {

	println(":: running early hook [udev]")
	time.Sleep(mT.t1000)
	println("starting system version 253")
	println(":: running early hook [initiso_system]")
	time.Sleep(mT.t500)
	for lineNo, line := range startHooks {
		if lineNo == 0 {
			fmt.Println(":: Triggering uevents...")
			time.Sleep(mT.t200)
		} else if lineNo == 4 {
			println(":: Mounting '/dev/disk/by-label/rootfs' to '/'")
			time.Sleep(mT.t2)
			println(":: Device '/dev/disk/by-label/rootfs' mounted successfully.")
			time.Sleep(mT.t2)
			println(":: Mounting '/dev/disk/by-label/efipart' to '/boot/efi'")
			time.Sleep(mT.t2)
			println(":: Device '/dev/disk/by-label/efipart' mounted successfully.")
			time.Sleep(mT.t2)
			println(":: Mounting (tmpfs) filesystem, size=32m...")
		}
		time.Sleep(mT.t2)
		fmt.Printf(":: running hook [%s]\n", line)
	}
	println(":: running cleanup hook [udev]")
	println("Welcome to " + Bold + Green + "okami.dev" + Reset + "!")
	time.Sleep(mT.t1000)
}

func systemdStart(startUnits []string) {
	startUnitsLen := len(startUnits)

	delayOne := rand.Intn(startUnitsLen-1) + 1

	delayTwo := rand.Intn(startUnitsLen-1) + 1

	for lineNo, line := range startUnits {
		if lineNo == delayOne || lineNo == delayTwo {
			time.Sleep(mT.t500)
		} else {
			time.Sleep(mT.t10)
		}
		fmt.Println(strings.Replace(line, "[  OK  ]", "[  "+Green+Bold+"OK"+Reset+"  ]", -1))
	}
}

// Unit Startup shutdown types
var startTypes = [...]string{"Reached", "Started", "Listening", "Created"}

func readEmbeddedFile(fileObj embed.FS, path string) []string {
	byteArr, err := fileObj.ReadFile(path)
	if err != nil {
		log.Fatalf("Error %s", err)
	}
	strArr := strings.Split(string(byteArr), "\n")
	return strArr
}

func main() {
	// Clear the TTY
	clearTTY()

	// Pull the embedded files
	startUnits := readEmbeddedFile(myUnits, "assets/units")
	startHooks := readEmbeddedFile(myHooks, "assets/hooks")

	// Run the early boot process
	earlyBoot(startHooks)

	// Clear the screen again
	clearTTY()

	// Run the systemd start services
	systemdStart(startUnits)
}
