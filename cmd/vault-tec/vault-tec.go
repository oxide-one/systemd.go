package main

import (
	"embed"
	"fmt"
	"math/rand"
	"os"
	"sort"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/gdamore/tcell/v2/encoding"
	"github.com/mattn/go-runewidth"
	"github.com/oxide-one/systemd.go/internal/vaultTec"
	readembedded "github.com/oxide-one/systemd.go/pkg/readEmbedded"
	"github.com/oxide-one/systemd.go/pkg/sleeper"
)

//go:embed wordlist
var wordList embed.FS

func generatePasswordBlock(lines int) {
	wordlist := readembedded.File(wordList, "wordlist")
	fmt.Println(wordlist)
	const punctuationList string = ",.<>/?@':;~#}]{[+=-_)(*&^%$£\"!0123456789"
	var matches = []string{"{}", "[]", "()", "<>"}
	_ = matches

}

func generateAddressBlock(numberOfAddresses int) []string {
	// The deviation between the minimum and max is 396, always.
	const standardDeviation int = 396

	// Declare the minimum address values (000-FFF)
	const addrMin int = 0
	const addrMax int = 4095 - standardDeviation

	// Generate the start and end ranges
	var startAddress int = rand.Intn(addrMax-addrMin) + addrMin
	var endAddress int = startAddress + standardDeviation

	// Create a slice of addresses
	var addrSlice []int
	// Create a checkmap to see if the value already exists
	chkMap := make(map[int]bool)

	// Add the start and end address to the slice
	addrSlice = append(addrSlice, startAddress)
	addrSlice = append(addrSlice, endAddress)

	// Iterate over the number of addresses min 2, to create an address map
	for i := 0; i < numberOfAddresses-2; i++ {
		for {
			randAddress := startAddress + 1 + rand.Intn(standardDeviation)
			if _, ok := chkMap[randAddress]; !ok {
				addrSlice = append(addrSlice, randAddress)
				chkMap[randAddress] = true
				break
			}
		}
	}

	sort.Ints(addrSlice)

	var addrBlock []string
	for _, i := range addrSlice {
		addr := fmt.Sprintf("0xF%0X", i)
		addrBlock = append(addrBlock, addr)

	}

	return addrBlock
}

func emitStr(s tcell.Screen, x, y int, style tcell.Style, str string) {
	for _, c := range str {
		var comb []rune
		w := runewidth.RuneWidth(c)
		if w == 0 {
			comb = []rune{c}
			c = ' '
			w = 1
		}
		s.SetContent(x, y, c, comb, style)
		sleeper.Sleep(1)
		s.Show()
		x += w
	}
}

func blockDisplay(s tcell.Screen, x int, y int, style tcell.Style, strArr []string, cols int, offset int) {
	strArrLen := len(strArr)
	colXLength := int(strArrLen / cols)
	var columns [][]string
	for i := 0; i <= cols-1; i++ {
		if len(strArr) < colXLength {
			columns = append(columns, strArr)
			break
		} else {
			columns = append(columns, strArr[:colXLength])
			strArr = strArr[colXLength:]
		}

	}
	initialY := y
	for _, column := range columns {
		for _, str := range column {
			emitStr(s, x, y, style, str)
			y++
		}
		y = initialY
		x += offset

	}
}

func displayScreen(s tcell.Screen, addrBlock []string) {
	w, h := s.Size()
	_ = h
	_ = w
	s.Clear()
	style := tcell.StyleDefault.Foreground(tcell.ColorGreen.TrueColor()).Background(tcell.ColorBlack.TrueColor()).Bold(true)
	blockDisplay(s, 10, 10, style, addrBlock, 2, 56)
	s.Show()
}

func main() {
	rand.Seed(time.Now().UnixNano())
	generatePasswordBlock(200)

	// Clear the TTY
	//clear.ClearTTY()
	vaultTec.EarlyBoot()

	s, e := tcell.NewScreen()
	if e != nil {
		fmt.Fprintf(os.Stderr, "%v\n", e)
		os.Exit(1)
	}
	if e := s.Init(); e != nil {
		fmt.Fprintf(os.Stderr, "%v\n", e)
		os.Exit(1)
	}

	defStyle := tcell.StyleDefault.
		Background(tcell.ColorBlack).
		Foreground(tcell.ColorWhite)
	s.SetStyle(defStyle)

	encoding.Register()
	addrBlock := generateAddressBlock(30)
	for {
		switch ev := s.PollEvent().(type) {
		case *tcell.EventResize:
			s.Sync()
			displayScreen(s, addrBlock)
		case *tcell.EventKey:
			if ev.Key() == tcell.KeyEscape {
				s.Fini()
				os.Exit(0)
			}
		}
	}
}
