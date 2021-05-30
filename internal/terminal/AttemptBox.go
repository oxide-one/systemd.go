package Terminal

type AttemptBox struct {
	// The current lolumn
	TotalAttempts int
	// The current line
	RemainingAttempts int
	// The current line position
	Position MultiCoordinates
}

// func (t *Terminal) Init() {
// 	t.attemptBox.TotalAttempts = t.Settings.TotalAttempts
// 	t.attemptBox.RemainingAttempts = t.Settings.TotalAttempts
// 	t.attemptBox.Position = t.Header.Content[len(t.Header.Content)-1].Position
// }

// func (t *Terminal) flash(s tcell.Screen) {
// 	emitStr(s, t.attemptBox.Position.Start.X, t.attemptBox.Position.Start.Y, t.Style.Default, "HELLO")
// }
