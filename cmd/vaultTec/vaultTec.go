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

type colSet struct {
	start int
	end   int
}

type cursorPosition struct {
	column    int
	line      int
	linePos   int
	curXStart int
	curXEnd   int
	curY      int
	finalX    int
	finalY    int
	startY    int
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

		s.Show()
		x += w
	}

}

func refreshSelection(s tcell.Screen, cursor cursorPosition, passwordBlock [][][]string, defaultStyle tcell.Style, highlightStyle tcell.Style, highlighted bool) {
	column := cursor.column
	line := cursor.line
	linePos := cursor.linePos
	curX := cursor.curXStart
	curY := cursor.curY
	emitStr(s, 0, 0, highlightStyle, fmt.Sprint(column))
	emitStr(s, 0, 1, highlightStyle, fmt.Sprint(line))
	emitStr(s, 0, 2, highlightStyle, fmt.Sprint(linePos))
	if highlighted {
		emitStr(s, curX, curY, highlightStyle, passwordBlock[column][line][linePos])
	} else {
		emitStr(s, curX, curY, defaultStyle, passwordBlock[column][line][linePos])
	}

}

func displayScreen(s tcell.Screen, cursor cursorPosition, startX int, startY int, lineCount int, columnCount int, addrBlock [][]string, passwordBlock [][][]vaultTec.MemoryBlock, defaultStyle tcell.Style, highlightStyle tcell.Style) (cursorPosition, []colSet) {
	//s.Clear()
	// Set the style

	x := startX
	y := startY

	// Styling Vars
	columnInterPadding := 2  // Padding between the address and password lines
	columnOuterPadding := 2  // Padding between each column
	headerPaddingBottom := 5 // Padding between the header and the main blocks
	// Print out the Prompts
	emitStr(s, x, y, defaultStyle, "OKAMIDASH INDUSTRIES (TM) TERMLINK PROTOCOL")
	y += 1
	emitStr(s, x, y, defaultStyle, "ENTER PASSWORD NOW")
	y += 2
	emitStr(s, x, y, defaultStyle, "4 ATTEMPTS LEFT: ")
	columnStops := make([]colSet, columnCount)
	addressPadding := len(addrBlock[columnCount-1][lineCount-1])
	var totalLineWidth int
	for _, tlw := range passwordBlock[0][0] {
		totalLineWidth += len(tlw.value)
	}

	// Print out the Blocks
	totalColumnWidth := addressPadding + columnInterPadding + totalLineWidth + columnOuterPadding
	for column := 0; column < columnCount; column++ {
		columnStops[column].start = x + addressPadding + columnInterPadding
		y = startY + headerPaddingBottom
		cursor.startY = y
		var curX int
		for line := 0; line < lineCount; line++ {
			// Display the Address Block first
			address := addrBlock[column][line]
			emitStr(s, x, y, defaultStyle, address)

			// Now display the password block
			curX = x + addressPadding + columnInterPadding
			for linePos, myLine := range passwordBlock[column][line] {
				lineLength := len(myLine)
				if cursor.column == column && cursor.line == line && cursor.linePos == linePos {
					emitStr(s, curX, y, highlightStyle, myLine)
					cursor.curXStart = curX
					cursor.curXEnd = curX + lineLength
					cursor.curY = y

				} else {
					emitStr(s, curX, y, defaultStyle, myLine)
				}
				curX += lineLength
			}

			y++
		}
		columnStops[column].end = curX
		x += totalColumnWidth

	}
	var finalX int = x - columnOuterPadding
	var finalY int = y - 1
	cursor.finalX = finalX
	cursor.finalY = finalY
	emitStr(s, finalX, finalY, defaultStyle, ">")
	return cursor, columnStops
}

func moveUp(cursor cursorPosition, lineCount int) cursorPosition {
	if cursor.line == 0 {
		cursor.line = lineCount - 1
		cursor.curY = cursor.finalY
	} else {
		cursor.line--
		cursor.curY--
	}

	return cursor
}

func moveDown(cursor cursorPosition, lineCount int) cursorPosition {
	if cursor.line == lineCount-1 {
		cursor.line = 0
		cursor.curY = cursor.startY
	} else {
		cursor.line++
		cursor.curY++
	}

	return cursor
}

func moveRight(cursor cursorPosition, lineCount int, columnCount int, columnSet []colSet, passwordBlock [][][]string) cursorPosition {
	column := cursor.column
	line := cursor.line
	linePos := cursor.linePos
	if cursor.linePos == len(passwordBlock[column][line])-1 {
		cursor.linePos = 0
		if cursor.column == columnCount-1 {
			cursor.column = 0
			cursor.curXStart = columnSet[cursor.column].start
		} else {
			cursor.column += 1
			cursor.curXStart = columnSet[cursor.column].start
		}

	} else {
		cursor.curXStart = cursor.curXStart + len(passwordBlock[column][line][linePos])
		cursor.linePos += 1

	}

	cursor.curXEnd = cursor.curXStart + len(passwordBlock[column][line][linePos])
	return cursor
}

func moveLeft(cursor cursorPosition, lineCount int, columnCount int, columnSet []colSet, passwordBlock [][][]string) cursorPosition {
	if cursor.linePos == 0 {
		cursor.linePos = len(passwordBlock[cursor.column][cursor.line]) - 1
		if cursor.column == 0 {
			cursor.column = columnCount - 1
			cursor.curXStart = columnSet[cursor.column].end - len(passwordBlock[cursor.column][cursor.line][cursor.linePos]) + 1
			cursor.curXEnd = columnSet[cursor.column].end
		} else {
			cursor.column -= 1
			cursor.curXStart = columnSet[cursor.column].end - len(passwordBlock[cursor.column][cursor.line][cursor.linePos]) + 1
			cursor.curXEnd = columnSet[cursor.column].end
		}

	} else {
		cursor.linePos -= 1
		cursor.curXEnd = cursor.curXStart - 1
		cursor.curXStart = cursor.curXEnd - len(passwordBlock[cursor.column][cursor.line][cursor.linePos]) + 1

	}
	return cursor
}

func main() {
	rand.Seed(time.Now().UnixNano())
	passwords, passwordList := vaultTec.GeneratePassword(15)
	const columnCount int = 2
	const lineCount int = 30
	const lineWidth int = 30
	const startX int = 4
	const startY int = 2
	defaultStyle := tcell.StyleDefault.Foreground(tcell.ColorGreen.TrueColor()).Background(tcell.ColorBlack.TrueColor()).Bold(true)
	highlightStyle := tcell.StyleDefault.Foreground(tcell.ColorWhite.TrueColor()).Background(tcell.ColorGreen.TrueColor()).Bold(true)
	addrBlock := vaultTec.GenerateAddressBlock(lineCount, columnCount)
	passwordColumns := vaultTec.GeneratePasswordBlock(lineCount, lineWidth, columnCount, len(passwords), passwords, passwordList)
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

	defStyle := tcell.StyleDefault.
		Background(tcell.ColorBlack).
		Foreground(tcell.ColorWhite)
	s.SetStyle(defStyle)

	encoding.Register()

	cursor := cursorPosition{column: 0, line: 0, linePos: 0}
	columnSet := []colSet{}
	for {
		switch ev := s.PollEvent().(type) {
		case *tcell.EventResize:
			s.Sync()
			cursor, columnSet = displayScreen(s, cursor, startX, startY, lineCount, columnCount, addrBlock, passwordColumns, defaultStyle, highlightStyle)
		case *tcell.EventKey:

			if ev.Key() == tcell.KeyEscape {
				s.Fini()
				os.Exit(0)
			}
			refreshSelection(s, cursor, passwordColumns, defaultStyle, highlightStyle, false)
			switch ev.Key() {

			case tcell.KeyUp:
				cursor = moveUp(cursor, lineCount)
			case tcell.KeyDown:
				cursor = moveDown(cursor, lineCount)
			case tcell.KeyRight:
				cursor = moveRight(cursor, lineCount, columnCount, columnSet, passwordColumns)
			case tcell.KeyLeft:
				cursor = moveLeft(cursor, lineCount, columnCount, columnSet, passwordColumns)
			}

			emitStr(s, 70, 0, defStyle, fmt.Sprintf("LINE: %d, LINEPOS: %d, CURXSTART: %d, CURXEND %d, CURY %d, FINALX %d, FINALY %d", cursor.line, cursor.linePos, cursor.curXStart, cursor.curXEnd, cursor.curY, cursor.finalX, cursor.finalY))
			emitStr(s, 71, 1, defStyle, fmt.Sprint(columnSet))
			emitStr(s, 0, 0, defStyle, "")
			refreshSelection(s, cursor, passwordColumns, defaultStyle, highlightStyle, true)

		}
	}
}
