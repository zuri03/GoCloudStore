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
	Owner        string   `json:"owner"` //For now Onwer is just username:password
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

func (keeper *RecordKeeper) GetRecord(key string, username string, password string) (*Record, error) {
	record, ok := keeper.records[key]
	if !ok {
		return nil, fmt.Errorf("Not Found")
	}

	user := fmt.Sprintf("%s:%s", username, password)

	if record.Owner == user {
		return &record, nil
	}
	if record.Owner != user {
		for _, allowedUser := range record.AllowedUsers {
			if allowedUser == user {
				return &record, nil
			}
		}
	}
	return nil, fmt.Errorf("Unathorized")
}

func (keeper *RecordKeeper) SetRecord(key string, username string, password string, size int64, name string) (*Record, error) {
	_, ok := keeper.records[key]
	if ok {
		return nil, fmt.Errorf("Record %s already exists", key)
	}
	creationTime := time.Now().Format("YYYY-MM-DD") //TODO:Determine formatting
	now := time.Now()
	owner := fmt.Sprintf("%d-%02d-%02dT%02d:%02d:%02d",
		now.Year(), now.Month(), now.Day(),
		now.Hour(), now.Minute(), now.Second())
	record := Record{
		MetaData: &FileMetaData{
			Size: size,
			Name: name,
		},
		Location:     fmt.Sprintf("%s/%s", username, name), //Ex: user/foo.txt
		CreatedAt:    creationTime,
		IsPublic:     false,
		Owner:        owner,
		AllowedUsers: make([]string, 0),
	}
	keeper.records[key] = record
	return &record, nil
}

func (keeper *RecordKeeper) RemoveRecord(key string, username string, password string) error {
	record, ok := keeper.records[key]
	if !ok {
		return fmt.Errorf("Not found")
	}

	user := fmt.Sprintf("%s:%s", username, password)
	if record.Owner != user {
		return fmt.Errorf("Unathorized")
	}
	delete(keeper.records, key)
	return nil
}

func (keeper *RecordKeeper) AddAllowedUser(key string, ownerUsername string, ownerPassword, username string) error {
	record, ok := keeper.records[key]
	if !ok {
		return fmt.Errorf("Record %s does not exist", key)
	}

	owner := fmt.Sprintf("%s:%s", ownerUsername, ownerPassword)
	if record.Owner != owner {
		return fmt.Errorf("Unathorized")
	}

	for i, c := range record.AllowedUsers {
		fmt.Printf("%d => %s\n", i, c)
	}
	fmt.Println("After")
	record.AllowedUsers = append(record.AllowedUsers, username)
	for i, c := range record.AllowedUsers {
		fmt.Printf("%d => %s\n", i, c)
	}
	keeper.records[key] = record
	return nil
}

func (keeper *RecordKeeper) RemoveAllowedUser(key string, ownerUsername string, ownerPassword string, user string) error {
	record, ok := keeper.records[key]
	if !ok {
		return fmt.Errorf("Record %s does not exist", key)
	}

	index, err := findItemIndex(user, record.AllowedUsers)
	if err != nil {
		return fmt.Errorf("User %s is not in the allowed list", user)
	}

	if fmt.Sprintf("%s:%s", ownerUsername, ownerPassword) != record.Owner {
		return fmt.Errorf("Unathorized")
	}

	record.AllowedUsers[index] = record.AllowedUsers[len(record.AllowedUsers)-1]
	record.AllowedUsers = record.AllowedUsers[:len(record.AllowedUsers)-1]
	keeper.records[key] = record
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
