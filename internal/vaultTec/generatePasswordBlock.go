package vaultTec

import (
	"fmt"
	"math/rand"
	"sort"
)

func GeneratePasswordBlock(terminal Terminal, passwordList []string) [][][]MemoryBlock {

	// List of allowed punctuation
	var punctuationList = []string{",", ",", ".", "<", ">", "/", "?", "@", "'", ":", ";", "~", "#", "}", "]", "{", "[", "+", "=", "-", "_", ")", "(", "*", "&", "^", "%", "$", "\"", "!", "0", "1", "2", "3", "4", "5", "6", "7", "8", "9"}
	var punctuationListLen int = len(punctuationList)
	// List of matched pairs
	var matches = []string{"{}", "[]", "()", "<>"}
	var matchLen int = len(matches)

	var matchLines []int
	var matchCols []int
	// Checkmap so we don't get conflicts
	matchMap := make(map[string]bool)
	matchCount := rand.Intn(terminal.Settings.LineCount-10) + 10

	// Iterate through the number of passwords needed
	for i := 0; i < matchCount; i++ {
		for {
			// Generate a random number between 0 and the lineCount
			lineNo := rand.Intn(terminal.Settings.LineCount)

			// Generate a random number between 0 and the columns
			colNo := rand.Intn(terminal.Settings.ColumnCount)

			// Generate a 'hash' of sorts to ensure no column and line collissions
			colLineHash := fmt.Sprintf("%d%d", colNo, lineNo)

			if _, ok := matchMap[colLineHash]; !ok {
				// Append the line number and column number
				matchLines = append(matchLines, lineNo)
				matchCols = append(matchCols, colNo)

				// Set the hash value
				matchMap[colLineHash] = true
				break
			}
		}
	}
	// Sort the matches
	sort.Ints(matchLines)

	// The Line that passwords will sit on
	var passwordLines []int
	var passwordCols []int
	// Checkmap so we don't get conflicts
	passMap := make(map[string]bool)

	// Iterate through the number of passwords needed
	for i := 0; i < terminal.Settings.PasswordCount; i++ {
		for {
			// Generate a random number between 0 and the lineCount
			lineNo := rand.Intn(terminal.Settings.LineCount)

			// Generate a random number between 0 and the columns
			colNo := rand.Intn(terminal.Settings.ColumnCount)

			// Generate a 'hash' of sorts to ensure no column and line collissions
			colLineHash := fmt.Sprintf("%d%d", colNo, lineNo)

			// Evaluate the hashes
			_, passMapChk := passMap[colLineHash]
			_, matchMapChk := matchMap[colLineHash]

			// If There are no collissions, add it in.
			if !matchMapChk && !passMapChk {
				// Append the line number and column number
				passwordLines = append(passwordLines, lineNo)
				passwordCols = append(passwordCols, colNo)

				// Set the hash value
				passMap[colLineHash] = true
				break
			}
		}
	}
	// Sort the arrays
	sort.Ints(passwordLines)

	var passwordColumns [][][]MemoryBlock
	for column := 0; column < terminal.Settings.ColumnCount; column++ {
		var passwordLines [][]MemoryBlock
		for line := 0; line < terminal.Settings.LineCount; line++ {
			// Generate a 'hash' of sorts to ensure no column and line collissions
			colLineHash := fmt.Sprintf("%d%d", column, line)

			// Evaluate the hashes
			_, passMapChk := passMap[colLineHash]
			_, matchMapChk := matchMap[colLineHash]

			lineCharsLeft := terminal.Settings.LineWidth
			var currentBuffer []string
			if passMapChk {
				currentBuffer = append(currentBuffer, passwordList[0])
				lineCharsLeft -= len(passwordList[0])
				passwordList = passwordList[1:]

			}

			if matchMapChk {
				matchIndex := rand.Intn(matchLen)
				currentBuffer = append(currentBuffer, matches[matchIndex])
				lineCharsLeft -= len(matches[matchIndex])
			}
			// Fill the remaining matches with random punctuation
			for fillChars := lineCharsLeft; fillChars > 0; fillChars-- {
				puncIndex := rand.Intn(punctuationListLen)
				currentBuffer = append(currentBuffer, punctuationList[puncIndex])
			}

			passwordBlockLine := make([]MemoryBlock, len(currentBuffer))
			perm := rand.Perm(len(currentBuffer))
			for i, v := range perm {
				passwordBlockLine[v].Value = currentBuffer[i]
				passwordBlockLine[v].Length = len(currentBuffer[i])
			}
			passwordLines = append(passwordLines, passwordBlockLine)
		}
		passwordColumns = append(passwordColumns, passwordLines)
	}

	//fmt.Println(passwordBlockLines)

	//fmt.Println(passwordLines, passwordCols, matchLines, matchCols)
	fmt.Println(passwordColumns)
	return passwordColumns
}
