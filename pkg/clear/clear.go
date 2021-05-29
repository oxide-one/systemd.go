package clear

func ClearTTY() {
	println("\033[;H\033[2J")
}
