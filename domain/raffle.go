package domain

import (
	"math/rand"
)

func Draw(friendsList []Friend, pool []Friend, pairs Pairs) Pairs {
	if len(friendsList) == 0 {
		return pairs
	}

	// Get the first of the list
	picker := friendsList[0]
	// Remove the picker from the list
	newFriendsList := append([]Friend{}, friendsList[1:]...)

	var i int
	picked := picker

	// Picker friend cannot be equal the picked friend
	for picked == picker {
		i = rand.Intn(len(pool))
		picked = pool[i]
	}

	pairs[picker] = picked

	// Remove picked participant from the pool
	newPool := append(pool[:i], pool[i+1:]...)

	return Draw(newFriendsList, newPool, pairs)
}
