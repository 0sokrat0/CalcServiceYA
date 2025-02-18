package genid

import (
	"fmt"
	"sync"
)

type GenID struct {
	idCounter int
	idMutex   sync.Mutex
}

var globalGenID = &GenID{}

func GenerateID() string {
	globalGenID.idMutex.Lock()
	defer globalGenID.idMutex.Unlock()
	globalGenID.idCounter++
	return fmt.Sprintf("expression-%d", globalGenID.idCounter)
}
