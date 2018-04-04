package datastore

import (
	"errors"
	"fmt"
	"reflect"
	"sync"

	log "github.com/sirupsen/logrus"
)

//LearningTaskMap is implementation of LearningTaskData interface
type LearningTaskMap struct {
	learningTaskMap map[string]*LearningTaskInfo
}

var pool LearningTaskMap
var mutex sync.RWMutex

//GetLearningTaskMap is singleton impl of LearningTaskMap instance
func GetLearningTaskMap() LearningTaskMap {
	//initialze map if it is not yet
	if pool.learningTaskMap == nil {
		pool = LearningTaskMap{}
		pool.learningTaskMap = make(map[string]*LearningTaskInfo)
		mutex = sync.RWMutex{}
		log.Debug("data store map initialized")
	}

	return pool
}

//Get method returns LearningTaskInfo coressponding to passed learningTask name
func (pool LearningTaskMap) Get(lt string) (LearningTaskInfo, error) {

	//get and return LearningTaskinfo
	mutex.RLock()
	rtn, ok := pool.learningTaskMap[lt]
	mutex.RUnlock()

	//error check if ok is false, then there is no such item
	if !ok {
		err := fmt.Errorf("Item not found. key:%s", lt)
		return LearningTaskInfo{}, err
	}

	return *rtn, nil
}

//Put method store passed struct into map
func (pool LearningTaskMap) Put(lt LearningTaskInfo) error {

	var err error
	key := lt.Name
	if key == "" {
		//basically dead code
		err = errors.New("learningTask has no name")
		return err
	}

	mutex.Lock()
	pool.learningTaskMap[key] = &lt
	mutex.Unlock()

	return nil
}

//Remove method store passed struct into map
func (pool LearningTaskMap) Remove(lt string) error {
	mutex.Lock()
	defer mutex.Unlock()

	//check key existstance
	_, exist := pool.learningTaskMap[lt]
	if !exist {
		return fmt.Errorf("Item not found. key:%s", lt)
	}

	delete(pool.learningTaskMap, lt)

	return nil
}

//GetAll method returns all learningTasks info from map
func (pool LearningTaskMap) GetAll() ([]LearningTaskInfo, error) {

	var err error
	var rtn []LearningTaskInfo
	//get and return LearningTaskinfo
	mutex.RLock()
	keys := reflect.ValueOf(pool.learningTaskMap).MapKeys()

	for _, key := range keys {
		rtn = append(rtn, *pool.learningTaskMap[key.String()])
	}
	mutex.RUnlock()

	return rtn, err
}

func (pool LearningTaskMap) UpdateState(lt string, state string, time string) error {
	mutex.Lock()
	defer mutex.Unlock()

	var info *LearningTaskInfo
	var ok bool
	if info, ok = pool.learningTaskMap[lt]; !ok {
		return fmt.Errorf("learningTask %s not found", lt)
	}

	info.State = state

	// set learning task exec time
	if time != "" {
		info.ExecTime = time
	}

	return nil
}

// UpdatePodState update pod's state value
func (pool LearningTaskMap) UpdatePodState(lt string, pod string, state string) error {
	mutex.Lock()
	defer mutex.Unlock()

	var info *LearningTaskInfo
	var ok bool
	if info, ok = pool.learningTaskMap[lt]; !ok {
		return fmt.Errorf("learningTask %s not found", lt)
	}

	info.PodState[pod] = state
	return nil
}
