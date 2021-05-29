package systemd

import (
	"embed"
	"log"
	"strings"
	"time"

	"github.com/oxide-one/systemd.go/internal/systemd"
	"github.com/oxide-one/systemd.go/pkg/clear"
	readembedded "github.com/oxide-one/systemd.go/pkg/readEmbedded"
)

//go:embed ./..//assets/hooks
var myHooks embed.FS

//go:embed ./assets/units
var myUnits embed.FS

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

func systemdRun() {
	// Clear the TTY
	clear.ClearTTY()

	// Pull the embedded files
	startUnits := readembedded.File(myUnits, "assets/units")
	startHooks := readembedded.File(myHooks, "assets/hooks")

	// Run the early boot process
	systemd.EarlyBoot(startHooks)

	// Clear the screen again
	clear.ClearTTY()

	// Run the systemd start services
	systemd.SystemdStart(startUnits)
}
