package Terminal

import (
	"fmt"
	"strings"

	"github.com/gdamore/tcell/v2"
)

type InputBox struct {
	// The current line position
	Position SimpleCoordinates
	// The last line length
	LastLength int
	Messages   []string
}

func (i *InputBox) Init(t Terminal) {
	i.Position = t.MemoryBlocks[len(t.MemoryBlocks)-1].Position.End
	i.Position.X += 2
	i.Position.Y -= 1
	i.LastLength = 0
}

func (i *InputBox) flash(t Terminal, s tcell.Screen, cell String) {
	X := i.Position.X
	Y := i.Position.Y

	// Move back and wipe the buffer
	messages := i.Messages
	for P := len(i.Messages); P > 0; P-- {
		emitStr(s, X, Y-(2+(P)), t.Style.Default, strings.Repeat(" ", len(messages[P])))
	}

	// Wipe the current line
	emptyBoxes := strings.Repeat(" ", i.LastLength)
	emitStr(s, i.Position.X, i.Position.Y, t.Style.Default, emptyBoxes)

	// Compile the current string
	fullString := fmt.Sprintf("> %s", cell.Content)
	i.Messages = append(i.Messages, fullString)
	emitStr(s, i.Position.X, i.Position.Y, t.Style.Default, fullString)
	i.LastLength = len(fullString)
}
