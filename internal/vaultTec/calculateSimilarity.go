package vaultTec

func calculateSimilarity(chosenPassword string, selectedPassword string) int {
	// Calculate the similarity of words by the number of shared letters in them
	// Generate a checkMap to append each rune to
	chkMap := make(map[rune]bool)
	for _, char := range chosenPassword {
		chkMap[char] = true
	}
	var matchingLetters int
	for _, char := range selectedPassword {
		if _, ok := chkMap[char]; ok {
			matchingLetters++
		}
	}
	return matchingLetters
}
