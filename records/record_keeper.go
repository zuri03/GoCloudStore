package records

import (
	"fmt"
	"time"
)

//TODO: Make record keeper thread safe
//TODO: Replace record keeper with database
type RecordKeeper struct {
	records map[string]Record
}

type Record struct {
	MetaData *FileMetaData `json:"file"`
	//Location points to the node and path on that node that the file is stored on
	//For now files are simply stored in a single chunk in a single location
	Location     string   `json:"location"`
	CreatedAt    string   `json:"createdAt"`
	IsPublic     bool     `json:"isPublic"`
	Owner        string   `json:"owner"`
	AllowedUsers []string `json:"allowedUsers"`
}

type FileMetaData struct {
	Size int64
	Name string
}

//Ensure this is only called once in the main function
func InitRecordKeeper() RecordKeeper {
	return RecordKeeper{
		records: make(map[string]Record),
	}
}

func (keeper *RecordKeeper) Get(key string, id string) (*Record, error) {
	record, ok := keeper.records[key]
	if !ok {
		return nil, fmt.Errorf("Not Found")
	}

	if record.Owner == id {
		return &record, nil
	}

	if record.Owner != id {
		for _, allowedUser := range record.AllowedUsers {
			if allowedUser == id {
				return &record, nil
			}
		}
	}
	return nil, fmt.Errorf("Unauthorized")
}

func (keeper *RecordKeeper) New(key string, id string, size int64, name string) (*Record, error) {
	_, ok := keeper.records[key]
	if ok {
		return nil, fmt.Errorf("Record %s already exists", key)
	}

	now := time.Now()
	creationTime := fmt.Sprintf("%d-%02d-%02dT%02d:%02d:%02d",
		now.Year(), now.Month(), now.Day(),
		now.Hour(), now.Minute(), now.Second())
	record := Record{
		MetaData: &FileMetaData{
			Size: size,
			Name: name,
		},
		Location:     fmt.Sprintf("%s/%s", id, name), //Ex: user/foo.txt
		CreatedAt:    creationTime,
		IsPublic:     false,
		Owner:        id,
		AllowedUsers: make([]string, 0),
	}
	keeper.records[key] = record
	return &record, nil
}

func (keeper *RecordKeeper) Remove(key string, id string) error {
	record, ok := keeper.records[key]
	if !ok {
		return fmt.Errorf("Not found")
	}

	if record.Owner != id {
		return fmt.Errorf("Unathorized")
	}
	delete(keeper.records, key)
	return nil
}

func (keeper *RecordKeeper) AddAllowedUser(key string, ownerId string, allowedId string) error {
	record, ok := keeper.records[key]
	if !ok {
		return fmt.Errorf("Record %s does not exist", key)
	}

	if record.Owner != ownerId {
		return fmt.Errorf("Unathorized")
	}

	if record.AllowedUsers == nil {
		record.AllowedUsers = []string{allowedId}
	} else {
		record.AllowedUsers = append(record.AllowedUsers, allowedId)
	}

	keeper.records[key] = record
	return nil
}

func (keeper *RecordKeeper) RemoveAllowedUser(key string, ownerId string, removedId string) error {
	record, ok := keeper.records[key]
	if !ok {
		return fmt.Errorf("Record %s does not exist", key)
	}

	index, err := findItemIndex(removedId, record.AllowedUsers)
	if err != nil {
		return fmt.Errorf("User %s is not in the allowed list", removedId)
	}

	if ownerId != record.Owner {
		return fmt.Errorf("Unathorized")
	}

	record.AllowedUsers[index] = record.AllowedUsers[len(record.AllowedUsers)-1]
	record.AllowedUsers = record.AllowedUsers[:len(record.AllowedUsers)-1]
	keeper.records[key] = record
	return nil
}

func (keeper *RecordKeeper) Exists(key string) bool {
	_, ok := keeper.records[key]
	return ok
}

func (keeper *RecordKeeper) Authorized(key string, id string) (bool, error) {
	record, ok := keeper.records[key]
	if !ok {
		return false, fmt.Errorf("%s not found \n", key)
	}

	if record.Owner == id {
		return true, nil
	}

	for _, user := range record.AllowedUsers {
		if id == user {
			return true, nil
		}
	}
	return false, nil
}

func (keeper *RecordKeeper) IsOnwer(key string, id string) (bool, error) {
	record, ok := keeper.records[key]
	if !ok {
		return false, fmt.Errorf("%s not found \n", key)
	}

	if record.Owner == id {
		return true, nil
	}

	return false, nil
}

func findItemIndex(item string, arr []string) (int, error) {
	for i := 0; i < len(arr); i++ {
		if arr[i] == item {
			return i, nil
		}
	}
	return -1, fmt.Errorf("Could not find %s in array", item)
}
