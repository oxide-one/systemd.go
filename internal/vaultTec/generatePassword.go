package vaultTec

import (
	"math/rand"
	"strings"

	readembedded "github.com/oxide-one/systemd.go/pkg/readEmbedded"
)

func GeneratePassword(passwordCount int) (map[string]PassStruct, []string) {
	// Grab the embedded Wordlist
	wordList := readembedded.File(wordList, "wordlist")
	// Calculate the length of the wordlist
	wordListLen := len(wordList)

	// Make the passwordList and the checkmap
	var passwordList []string
	chkMap := make(map[string]bool)

	// Iterate through the count of passwords and
	for i := 0; i < passwordCount; i++ {
		for {
			selectedPassword := strings.ToUpper(wordList[rand.Intn(wordListLen)])
			if _, ok := chkMap[selectedPassword]; !ok {
				passwordList = append(passwordList, selectedPassword)
				chkMap[selectedPassword] = true
				break
			}
		}
	}

	chosenPassword := passwordList[rand.Intn(passwordCount)]
	passwords := make(map[string]PassStruct)
	for _, selectedPassword := range passwordList {

		var correctPassword bool
		if selectedPassword == chosenPassword {
			correctPassword = true
		} else {
			correctPassword = false
		}

		passwords[selectedPassword] = PassStruct{
			Password:   selectedPassword,
			Correct:    correctPassword,
			Length:     len(selectedPassword),
			Similarity: calculateSimilarity(chosenPassword, selectedPassword),
		}
	}
	return passwords, passwordList
}
