package file

import (
	"bufio"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strings"
)

type Pairs map[string]string

func ReadFile(filePath string) {
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatalf("error when opening file %s: %v", filePath, err)
	}

	reader := bufio.NewReader(file)

	content := make([]byte, 1024)
	_, err = reader.Read(content)
	if err != nil {
		log.Fatalf("error when opening file %s: %v", filePath, err)
	}

	friendsList := strings.Split(string(content), "\n")
	if len(friendsList) == 0 {
		log.Println("Empty list")
		os.Exit(0)
	}

	log.Printf("Friends List: %d\n", len(friendsList))
	for _, friendName := range friendsList {
		fmt.Println(friendName)
	}

	if len(friendsList) % 2 != 0 {
		log.Println("adding an extra entry")
		friendsList = append(friendsList, "CHOOSE ONE, CHOSEN ONE!")
	}

	pairs := draw(friendsList, Pairs{})

	log.Println("List of Pairs:")
	for k, v := range pairs {
		fmt.Printf("%s <> %s\n", k, v)
	}
}

func draw(friendsList []string, pairs Pairs) Pairs {
	if len(friendsList) == 0 {
		return pairs
	}

	i := rand.Intn(len(friendsList))
	friendA := friendsList[i]

	newFriendsList := append(friendsList[:i], friendsList[i+1:]...)

	j := rand.Intn(len(newFriendsList))
	friendB := friendsList[j]

	newFriendsList = append(newFriendsList[:j], newFriendsList[j+1:]...)
	pairs[friendA] = friendB

	return draw(newFriendsList, pairs)
}
