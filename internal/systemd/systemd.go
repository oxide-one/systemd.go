package systemd

import (
	"fmt"
	"math/rand"
	"strings"

	"github.com/oxide-one/systemd.go/pkg/color"
	"github.com/oxide-one/systemd.go/pkg/sleeper"
)

func EarlyBoot(startHooks []string) {
	println(":: running early hook [udev]")
	sleeper.Sleep(1000)
	println("starting system version 253")
	println(":: running early hook [initiso_system]")
	sleeper.Sleep(500)
	for lineNo, line := range startHooks {
		if lineNo == 0 {
			fmt.Println(":: Triggering uevents...")
			sleeper.Sleep(2000)
		} else if lineNo == 4 {
			println(":: Mounting '/dev/disk/by-label/rootfs' to '/'")
			sleeper.Sleep(2)
			println(":: Device '/dev/disk/by-label/rootfs' mounted successfully.")
			sleeper.Sleep(2)
			println(":: Mounting '/dev/disk/by-label/efipart' to '/boot/efi'")
			sleeper.Sleep(2)
			println(":: Device '/dev/disk/by-label/efipart' mounted successfully.")
			sleeper.Sleep(2)
			println(":: Mounting (tmpfs) filesystem, size=32m...")
		}
		sleeper.Sleep(2)
		fmt.Printf(":: running hook [%s]\n", line)
	}
	println(":: running cleanup hook [udev]")
	println("Welcome to " + color.Bold + color.Green + "okami.dev" + color.Reset + "!")
	sleeper.Sleep(1000)
}

func SystemdStart(startUnits []string) {
	startUnitsLen := len(startUnits)

	delayOne := rand.Intn(startUnitsLen-1) + 1

	delayTwo := rand.Intn(startUnitsLen-1) + 1

	for lineNo, line := range startUnits {
		if lineNo == delayOne || lineNo == delayTwo {
			sleeper.Sleep(500)
		} else {
			sleeper.Sleep(10)
		}
		fmt.Println(strings.Replace(line, "[  OK  ]", "[  "+color.Green+color.Bold+"OK"+color.Reset+"  ]", -1))
	}
}
