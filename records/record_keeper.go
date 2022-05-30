package records

import (
	"fmt"
	"time"
)

type RecordKeeper struct {
	records map[string]Record
}

type Record struct {
	MetaData *FileMetaData
	//Location points to the node and path on that node that the file is stored on
	//For now files are simply stored in a single chunk in a single location
	Location     string
	CreatedAt    string
	IsPublic     bool
	Owner        string
	AllowedUsers []string
}

type FileMetaData struct {
	Size int64
	Name string
}

func (keeper *RecordKeeper) getRecord(key string) *Record {
	record, ok := keeper.records[key]
	if !ok {
		return nil
	}
	return &record
}

func (keeper *RecordKeeper) setRecord(key string, owner string, meta *FileMetaData) error {
	_, ok := keeper.records[key]
	if ok {
		return fmt.Errorf("Record %s already exists", key)
	}
	creationTime := time.Now().Format("YYYY-MM-DD") //TODO:Determine formatting
	keeper.records[key] = Record{
		MetaData:     meta,
		Location:     fmt.Sprintf("%s/%s", owner, meta.Name), //Ex: user/foo.txt
		CreatedAt:    creationTime,
		IsPublic:     false,
		Owner:        owner,
		AllowedUsers: make([]string, 0),
	}
	return nil
}

func (keeper *RecordKeeper) removeRecord(key string) error {
	if record := keeper.getRecord(key); record == nil {
		return fmt.Errorf("Record %s does not exist", key)
	}
	delete(keeper.records, key)
	return nil
}

func (keeper *RecordKeeper) addAllowedUser(key string, user string) error {
	record := keeper.getRecord(key)
	if record == nil {
		return fmt.Errorf("Record %s does not exist", key)
	}
	record.AllowedUsers = append(record.AllowedUsers, user)
	return nil
}

func (keeper *RecordKeeper) removeAllowedUser(key string, user string) error {
	record := keeper.getRecord(key)
	if record == nil {
		return fmt.Errorf("Record %s does not exist", key)
	}

	index, err := findItemIndex(user, record.AllowedUsers)
	if err != nil {
		return fmt.Errorf("User %s is not in the allowed list", user)
	}

	record.AllowedUsers[index] = record.AllowedUsers[len(record.AllowedUsers)-1]
	record.AllowedUsers = record.AllowedUsers[:len(record.AllowedUsers)-1]
	return nil
}

func findItemIndex(item string, arr []string) (int, error) {
	for i := 0; i < len(arr); i++ {
		if arr[i] == item {
			return i, nil
		}
	}
	return -1, fmt.Errorf("Could not find %s in array", item)
}
