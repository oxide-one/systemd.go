package vaultTec

import (
	"fmt"
	"math/rand"
	"sort"
)

// func GeneratePasswordBlock(lineCount int, lineWidth int, columnCount int, passwordCount int, passwords map[string]passStruct, passwordList []string) [][][]string {
func GenerateAddressBlock(terminalSettings TerminalSettings) [][]string {
	// The deviation between the minimum and max is 396, always.
	const standardDeviation int = 396

	// Declare the minimum address values (000-FFF)
	const addrMin int = 0
	const addrMax int = 4095

	// Generate the start and end ranges
	var startAddress int = rand.Intn(addrMax-addrMin) + addrMin
	var endAddress int = startAddress + standardDeviation

	// Create a checkmap to see if the value already exists
	chkMap := make(map[int]bool)
	// Determine the total line count
	totalLineCount := (terminalSettings.LineCount * terminalSettings.ColumnCount) - 2

	// Create a list of addresses
	var addrList []int
	addrList = append(addrList, startAddress)
	// Iterate over the number of addresses min 2, to create an address map
	for line := 0; line < totalLineCount; line++ {
		for {
			randAddress := startAddress + 1 + rand.Intn(standardDeviation)
			if _, ok := chkMap[randAddress]; !ok {
				addrList = append(addrList, randAddress)
				chkMap[randAddress] = true
				break
			}
		}
	}
	addrList = append(addrList, endAddress)
	sort.Ints(addrList)

	// Create a slice of addresses
	var addrSlice [][]string
	for column := 0; column < terminalSettings.ColumnCount; column++ {
		var addrColumn []string
		for line := 0; line < terminalSettings.LineCount; line++ {
			addrColumn = append(addrColumn, fmt.Sprintf("0xF%0X", addrList[line]))
		}
		addrList = addrList[terminalSettings.LineCount:]
		addrSlice = append(addrSlice, addrColumn)
	}

	return addrSlice
}
