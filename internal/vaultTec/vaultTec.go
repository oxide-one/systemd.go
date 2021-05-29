package vaultTec

import (
	"fmt"

	"github.com/oxide-one/systemd.go/pkg/clear"
	"github.com/oxide-one/systemd.go/pkg/color"
	gradualtype "github.com/oxide-one/systemd.go/pkg/gradualType"
	"github.com/oxide-one/systemd.go/pkg/sleeper"
)

// echo -e -n "\x1b[\x30 q" # changes to blinking block
// echo -e -n "\x1b[\x31 q" # changes to blinking block also
// echo -e -n "\x1b[\x32 q" # changes to steady block
// echo -e -n "\x1b[\x33 q" # changes to blinking underline
// echo -e -n "\x1b[\x34 q" # changes to steady underline
// echo -e -n "\x1b[\x35 q" # changes to blinking bar
// echo -e -n "\x1b[\x36 q" # changes to steady bar

func EarlyBoot() {
	clear.ClearTTY()
	fmt.Print(color.BG)
	sleeper.Sleep(100)
	fmt.Println(color.Green + color.Bold + "SECURITY RESET...")
	fmt.Println()
	sleeper.Sleep(500)
	gradualtype.GradualType("WELCOME TO OKAMIDASH INDUSTRIES (TM) TERMLINK", 10, color.Green)
	fmt.Println()
	fmt.Print(color.Green + color.Bold + ">")
	sleeper.Sleep(500)
	gradualtype.GradualType("SET TERMINAL/INQUIRE", 30, color.Green, color.Bold)
	fmt.Println()
	sleeper.Sleep(300)
	gradualtype.GradualType("OXIDE.ONE-V300", 10, color.Green, color.Bold)
	fmt.Println()
	fmt.Print(color.Green + color.Bold + ">")
	sleeper.Sleep(500)
	gradualtype.GradualType("SET FILE/PROTECTION=OWNER:RWED ACCOUNTS.F", 30, color.Green, color.Bold)
	fmt.Print(color.Green + color.Bold + ">")
	sleeper.Sleep(500)
	gradualtype.GradualType("SET HALT RESTART/MAIN", 30, color.Green, color.Bold)
	sleeper.Sleep(200)
}
