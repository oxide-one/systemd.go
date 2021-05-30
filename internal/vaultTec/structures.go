package vaultTec

import "github.com/gdamore/tcell/v2"

type Coordinates struct {
	Xstart int
	Xend   int
	Ystart int
	Yend   int
}

type RelativePosition struct {
	Column  int
	Line    int
	LinePos int
}

type TerminalStyle struct {
	Regular   tcell.Style
	Highlight tcell.Style
	// Padding between the address and password lines
	ColumnInterPadding int
	// Padding between each column
	ColumnOuterPadding int
	// Padding between the header and main blocks
	HeaderPadding int
	// Padding for the ADDRESS specifically
	AddressPadding int
}

type MemoryBlock struct {
	Value    string
	Length   int
	Position Coordinates
}

type PassStruct struct {
	Password   string
	Correct    bool
	Length     int
	Similarity int
}

type Attempts struct {
	Total     int
	Remaining int
}

type Terminal struct {
	// Relative position in array
	Position RelativePosition
	// Absolute	position on screen
	CursorLocation Coordinates

	// Absolute position of password line
	MemoryBlockLocation Coordinates

	// Absolute position of the address block
	AddressBlockLocation Coordinates
	// Location of the "X ATTEMPTS LEFT" column
	AttemptsBlockLocation Coordinates
	// Location of the input block location

	//	Settings for the tertminal
	Settings TerminalSettings

	// The default style
	Style TerminalStyle

	// The address Block
	AddressBlock [][]string

	// The Memory (password) block
	MemoryBlock [][][]MemoryBlock

	// A map of passwords
	Passwords map[string]PassStruct

	// Attempts left
	Attempts Attempts
	// curXStart   int
	// curXEnd     int
	// curY        int
	// finalX      int
	// finalY      int
	// startY      int
	// passwordX   int
	// passwordY   int
	// inputXStart int
	// inputXEnd   int
	// inputY      int
}
