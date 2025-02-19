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
var globalGenIDTask = &GenID{}

func GenerateID() string {
	globalGenID.idMutex.Lock()
	defer globalGenID.idMutex.Unlock()
	globalGenID.idCounter++
	return fmt.Sprintf("expression-%d", globalGenID.idCounter)
}

func GenerateIDTask() string {
	globalGenIDTask.idMutex.Lock()
	defer globalGenIDTask.idMutex.Unlock()
	globalGenIDTask.idCounter++
	return fmt.Sprintf("task-%d", globalGenIDTask.idCounter)
}
