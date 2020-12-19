package raffleservice

import (
	"github.com/obiwandsilva/go-secretfriend/domain/entities"
	"math/rand"
	"time"
)

func Draw(friendsList []entities.Friend, pool []entities.Friend, pairs entities.Pairs) entities.Pairs {
	if len(friendsList) == 0 {
		return pairs
	}

	// Get the first of the list
	picker := friendsList[0]
	// Remove the picker from the list
	newFriendsList := append([]entities.Friend{}, friendsList[1:]...)

	var i int
	picked := picker

	// Picker friend cannot be equal the picked friend
	for picked == picker {
		rand.Seed(time.Now().UnixNano())
		i = rand.Intn(len(pool))
		picked = pool[i]
	}

	pairs[picker] = picked

	// Remove picked participant from the pool
	newPool := append(pool[:i], pool[i+1:]...)

	return Draw(newFriendsList, newPool, pairs)
}

func GeneratePool(friendsList []entities.Friend) (pool []entities.Friend) {
	for i := len(friendsList)-1; i > -1; i-- {
		pool = append(pool, friendsList[i])
	}

	return pool
}
