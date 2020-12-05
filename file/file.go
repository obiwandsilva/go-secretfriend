package file

import (
	"bufio"
	"fmt"
	"github.com/obiwandsilva/go-secretfriend/domain"
	"os"
	"strings"
)

func ReadFile(filePath string) ([]domain.Friend, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("error when opening file %s: %v", filePath, err)
	}

	reader := bufio.NewReader(file)

	content := make([]byte, 1024)
	n, err := reader.Read(content)
	if err != nil {
		return nil, fmt.Errorf("error when opening file %s: %v", filePath, err)
	}

	friendsData := strings.Split(string(content[:n]), "\n")
	if len(friendsData) < 3 {
		return nil, fmt.Errorf("friends list cannot have less than 3 friends")
	}

	friendsList := make([]domain.Friend, 0)
	for i, friendData := range friendsData {
		split := strings.Split(friendData, ";")

		if len(split) != 2 ||
			len(split[0]) == 0 ||
			len(split[1]) == 0 {
			return nil, fmt.Errorf("invalid format at line %d", i+1)
		}

		friendsList = append(friendsList, domain.Friend{
			Name:  split[0],
			PhoneNumber: split[1],
		})
	}

	return friendsList, nil
}
