package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/gdamore/tcell/v2/encoding"
	"github.com/mattn/go-runewidth"
	"github.com/oxide-one/systemd.go/internal/vaultTec"
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

func refreshSelection(s tcell.Screen, terminal vaultTec.Terminal, highlighted bool) {
	column := terminal.Position.Column
	line := terminal.Position.Line
	linePos := terminal.Position.LinePos
	curX := terminal.CursorLocation.Xstart
	curY := terminal.CursorLocation.Ystart

	if highlighted {
		emitStr(s, curX, curY, terminal.Style.Highlight, terminal.MemoryBlock[column][line][linePos].Value)
	} else {
		emitStr(s, curX, curY, terminal.Style.Regular, terminal.MemoryBlock[column][line][linePos].Value)
	}

}

func calculateLocations(terminal vaultTec.Terminal) {
	//
	//0  [HeaderPaddingTop]
	//	 OKAMIDASH INDUSTIRES
	//	 ENTER PASSWORD NOW
	//
	//
	//	 [HeaderPaddingBottom]
	//	 [AddressWidtrh][ADDR_PADDING][MEMORY_WIDTH][MEMORY_PADDING]
	//   0x0000000fff              £W)£Q(A)TESTST
	//
	//
	//
	//
	//
	//
	//
	//
	//
	AddressWidth := terminal.Settings.
	
}
func displayScreen(s tcell.Screen, terminal vaultTec.Terminal) {
	//s.Clear()
	// Set the style

	x := terminal.Settings.Origin.Xstart
	y := terminal.Settings.Origin.Ystart

	// Print out the Prompts
	emitStr(s, x, y, terminal.Style.Regular, "OKAMIDASH INDUSTRIES (TM) TERMLINK PROTOCOL")
	y += 1
	emitStr(s, x, y, terminal.Style.Regular, "ENTER PASSWORD NOW")
	y += 2

	// Save the location of the attempts block
	terminal.AttemptsBlockLocation = vaultTec.Coordinates{
		Xstart: x,
		Ystart: y,
		Yend:   y,
		Xend:   x + 20,
	}
	//emitStr(s, x, y, defaultStyle, fmt.Sprintf("%d ATTEMPTS LEFT: %s", terminal.attempts.remaining, )
	y += terminal.Style.HeaderPaddingBottom

	// terminal.AddressBlockLocation.Xstart = x
	// terminal.AddressBlockLocation.Xend = x + terminal.Settings.LineCount

	// terminal.AddressBlockLocation.Ystart = y
	terminal.MemoryBlockLocation.Xstart = x + terminal.Style.AddressPadding + terminal.Style.ColumnInterPadding
	terminal.MemoryBlockLocation.Ystart = y

	// For each column
	for column := 0; column < terminal.Settings.ColumnCount; column++ {
		// Reset the position to 0,0 (relative to the memory blocks)
		x = terminal.AddressBlockLocation.Xstart + (column * (
		// |ADDRESS|| ||LINEWIDTH ||OUTERPADDING
		terminal.Style.AddressPadding + terminal.Style.ColumnInterPadding + terminal.Settings.LineWidth + terminal.Style.ColumnOuterPadding))
		y = terminal.AddressBlockLocation.Ystart

		for line := 0; line < terminal.Settings.LineCount; line++ {

			// Display the Address Block first
			address := terminal.AddressBlock[column][line]
			emitStr(s, x, y, terminal.Style.Regular, address)

			// Now display the Memory block
			curX := x + terminal.Style.AddressPadding + terminal.Style.ColumnInterPadding
			for linePos, myLine := range terminal.MemoryBlock[column][line] {

				lineLength := myLine.Length

				if terminal.Position.Column == column && terminal.Position.Line == line && terminal.Position.LinePos == linePos {
					emitStr(s, curX, y, terminal.Style.Highlight, myLine.Value)

					terminal.CursorLocation.Xstart = curX
					terminal.CursorLocation.Xend = curX + lineLength
					terminal.CursorLocation.Ystart = y
					terminal.CursorLocation.Yend = y

				} else {
					emitStr(s, curX, y, terminal.Style.Regular, myLine.Value)
				}
				terminal.MemoryBlock[column][line][linePos].Position.Xstart = curX
				terminal.MemoryBlock[column][line][linePos].Position.Xend = curX + lineLength
				terminal.MemoryBlock[column][line][linePos].Position.Ystart = y
				terminal.MemoryBlock[column][line][linePos].Position.Yend = y
				curX += lineLength
			}

			y++
		}
	}
	emitStr(s, finalX+2, finalY, defaultStyle, ">")
	emitStr(s, finalX+3, finalY, defaultStyle, "")
	return terminal, passwordBlock
}

func moveUp(terminal Terminal, lineCount int, passwordBlock [][][]vaultTec.MemoryBlock) Terminal {
	if terminal.position.line == 0 {
		terminal.position.line = lineCount - 1
		terminal.cursor.Ystart = terminal.finalY
		terminal.cursor.Yend = terminal.finalY
	} else {
		terminal.position.line--
		terminal.cursor.Ystart--
		terminal.cursor.Yend--
	}
	if terminal.position.linePos >= len(passwordBlock[terminal.column][terminal.line]) {
		terminal.linePos = len(passwordBlock[terminal.column][terminal.line]) - 1
	}
	midPoint := terminal.curXStart + ((terminal.curXEnd - terminal.curXStart) / 2)

	for linePos, myLine := range passwordBlock[terminal.column][terminal.line] {
		if myLine.StartX <= midPoint && myLine.EndX >= midPoint {
			terminal.linePos = linePos
		}
	}
	terminal.curXStart = passwordBlock[terminal.column][terminal.line][terminal.linePos].StartX
	terminal.curXEnd = passwordBlock[terminal.column][terminal.line][terminal.linePos].EndX
	terminal.curY = passwordBlock[terminal.column][terminal.line][terminal.linePos].StartY

	return terminal
}

func moveDown(cursor Terminal, lineCount int, passwordBlock [][][]vaultTec.MemoryBlock) Terminal {
	if cursor.line == lineCount-1 {
		cursor.line = 0
		cursor.curY = cursor.startY
	} else {
		cursor.line++
		cursor.curY++
	}
	if cursor.linePos >= len(passwordBlock[cursor.column][cursor.line]) {
		cursor.linePos = len(passwordBlock[cursor.column][cursor.line]) - 1
	}
	midPoint := cursor.curXStart + ((cursor.curXEnd - cursor.curXStart) / 2)

	for linePos, myLine := range passwordBlock[cursor.column][cursor.line] {
		if myLine.StartX <= midPoint && myLine.EndX >= midPoint {
			cursor.linePos = linePos
		}
	}
	cursor.curXStart = passwordBlock[cursor.column][cursor.line][cursor.linePos].StartX
	cursor.curXEnd = passwordBlock[cursor.column][cursor.line][cursor.linePos].EndX
	cursor.curY = passwordBlock[cursor.column][cursor.line][cursor.linePos].StartY
	return cursor
}

func moveRight(cursor Terminal, columnCount int, passwordBlock [][][]vaultTec.MemoryBlock) Terminal {
	// IF We are at the Right edge of a password block
	if cursor.linePos == len(passwordBlock[cursor.column][cursor.line])-1 {
		// Set the line position to zero
		cursor.linePos = 0
		// If we are at the furthest right column, wrap around to 0
		if cursor.column == columnCount-1 {
			cursor.column = 0
		} else { // If not, move to the next column across
			cursor.column += 1
		}
	} else { // If we are not at the furthest right edge, increment to the next word.
		cursor.linePos += 1
	}
	cursor.curXStart = passwordBlock[cursor.column][cursor.line][cursor.linePos].StartX
	cursor.curXEnd = passwordBlock[cursor.column][cursor.line][cursor.linePos].EndX
	return cursor
}

func moveLeft(terminal vaultTec.Terminal, columnCount int, passwordBlock [][][]vaultTec.MemoryBlock) vaultTec.Terminal {
	// IF We are at the left edge of a password block
	if cursor.linePos == 0 {
		// Set the line position to the max
		cursor.linePos = len(passwordBlock[cursor.column][cursor.line]) - 1
		// If we are at the furthest left column, wrap around to the right most
		if cursor.column == 0 {
			cursor.column = columnCount - 1
		} else { // If not, move to the next column across
			cursor.column -= 1
		}
	} else { // If we are not at the furthest right edge, increment to the next word.
		cursor.linePos -= 1
	}
	cursor.curXStart = passwordBlock[cursor.column][cursor.line][cursor.linePos].StartX
	cursor.curXEnd = passwordBlock[cursor.column][cursor.line][cursor.linePos].EndX
	return cursor
}

func main() {
	rand.Seed(time.Now().UnixNano())

	terminal := vaultTec.Terminal{}

	// Set the terminal style
	terminal.Style = vaultTec.TerminalStyle{
		Regular:             tcell.StyleDefault.Foreground(tcell.ColorGreen.TrueColor()).Background(tcell.ColorBlack.TrueColor()).Bold(true),
		Highlight:           tcell.StyleDefault.Foreground(tcell.ColorWhite.TrueColor()).Background(tcell.ColorGreen.TrueColor()).Bold(true),
		ColumnInterPadding:  5,
		ColumnOuterPadding:  5,
		HeaderPaddingBottom: 2,
		AddressPadding:      5,
	}

	// Set the settings of the terminal
	terminal.Settings = vaultTec.TerminalSettings{
		PasswordCount: 15,
		ColumnCount:   2,
		LineCount:     30,
		LineWidth:     30,
		Origin: vaultTec.Coordinates{
			Xstart: 4,
			Ystart: 2,
		},
	}

	// Generate the address Block
	terminal.AddressBlock = vaultTec.GenerateAddressBlock(terminal.Settings)

	var passwordList []string
	// Generate a list of passwords
	terminal.Passwords, passwordList = vaultTec.GeneratePassword(terminal.Settings.PasswordCount)

	// Generate a memory Block
	terminal.MemoryBlock = vaultTec.GeneratePasswordBlock(terminal, passwordList)

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

	s.SetStyle(terminal.Style.Regular)

	encoding.Register()

	terminal.CursorLocation = vaultTec.Coordinates{
		Xstart: 0,
		Xend:   0,
		Ystart: 0,
		Yend:   0,
	}
	for {
		switch ev := s.PollEvent().(type) {
		case *tcell.EventResize:
			s.Sync()
			terminal = displayScreen(s, terminal)
		case *tcell.EventKey:

			if ev.Key() == tcell.KeyEscape {
				s.Fini()
				os.Exit(0)
			}
			refreshSelection(s, terminal, false)
			switch ev.Key() {

			case tcell.KeyUp:
				terminal = moveUp(terminal)
			case tcell.KeyDown:
				terminal = moveDown(terminal)
			case tcell.KeyRight:
				terminal = moveRight(terminal)
			case tcell.KeyLeft:
				terminal = moveLeft(terminal)
			}

			//emitStr(s, 70, 0, defStyle, fmt.Sprintf("LINE: %d, LINEPOS: %d, CURXSTART: %d, CURXEND %d, CURY %d, FINALX %d, FINALY %d", cursor.line, cursor.linePos, cursor.curXStart, cursor.curXEnd, cursor.curY, cursor.finalX, cursor.finalY))
			//emitStr(s, 0, 0, defStyle, "")
			refreshSelection(s, terminal, true)

		}
	}
}
