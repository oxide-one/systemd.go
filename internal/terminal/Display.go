package Terminal

import (
	"fmt"
	"os"

	"github.com/gdamore/tcell/v2"
	"github.com/gdamore/tcell/v2/encoding"
	"github.com/mattn/go-runewidth"
	"github.com/oxide-one/systemd.go/pkg/clear"
)

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

		s.Show()
		x += w
	}

}

func displayHeader(terminal Terminal, s tcell.Screen) {
	//position := terminal.Header.Position
	content := terminal.Header.Content
	//x_start := position.Start.X
	//y_start := position.Start.Y
	for _, line := range content {
		linePosition := line.Position
		lineContent := line.Content
		emitStr(s, linePosition.Start.X, linePosition.Start.Y, terminal.Style.Default, lineContent)
	}
}

func displayBlocks(terminal Terminal, s tcell.Screen) {
	// Display the header
	addressBlocks := terminal.AddressBlocks
	memoryBlocks := terminal.MemoryBlocks
	for column := 0; column < terminal.Settings.Columns; column++ {
		addressBlock := addressBlocks[column]
		memoryBlock := memoryBlocks[column]
		_ = memoryBlock
		for line := 0; line < terminal.Settings.Lines; line++ {
			// Pulling address information
			addressLine := addressBlock.Content[line]
			addressLinePosition := addressLine.Position
			// Display the Addressblock
			emitStr(s, addressLinePosition.Start.X, addressLinePosition.Start.Y, terminal.Style.Default, addressLine.Content)
			// Display the Line block
			memoryLine := memoryBlock.Content[line]
			for stringSet := range memoryLine.Content {
				myString := memoryLine.Content[stringSet]
				myStringPosition := myString.Position
				emitStr(s, myStringPosition.Start.X, myStringPosition.Start.Y, tcell.StyleDefault, myString.Content)
			}
		}
	}
}

func Display(terminal Terminal) {
	// Clear the TTY
	clear.ClearTTY()
	//vaultTec.EarlyBoot()

	s, e := tcell.NewScreen()
	if e != nil {
		fmt.Fprintf(os.Stderr, "%v\n", e)
		os.Exit(1)
	}
	if e := s.Init(); e != nil {
		fmt.Fprintf(os.Stderr, "%v\n", e)
		os.Exit(1)
	}

	s.SetStyle(terminal.Style.Default)

	encoding.Register()
	displayHeader(terminal, s)
	for {
		switch ev := s.PollEvent().(type) {
		case *tcell.EventResize:
			s.Sync()
			displayBlocks(terminal, s)
		case *tcell.EventKey:

			if ev.Key() == tcell.KeyEscape {
				s.Fini()
				os.Exit(0)
			}
			// //refreshSelection(s, terminal, false)
			// switch ev.Key() {

			// case tcell.KeyUp:
			// 	terminal = moveUp(terminal)
			// case tcell.KeyDown:
			// 	terminal = moveDown(terminal)
			// case tcell.KeyRight:
			// 	terminal = moveRight(terminal)
			// case tcell.KeyLeft:
			// 	terminal = moveLeft(terminal)
			// }

			//emitStr(s, 70, 0, defStyle, fmt.Sprintf("LINE: %d, LINEPOS: %d, CURXSTART: %d, CURXEND %d, CURY %d, FINALX %d, FINALY %d", cursor.line, cursor.linePos, cursor.curXStart, cursor.curXEnd, cursor.curY, cursor.finalX, cursor.finalY))
			//emitStr(s, 0, 0, defStyle, "")
			//refreshSelection(s, terminal, true)

		}
	}

}
