package records

import (
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/zuri03/GoCloudStore/common"
)

type scenario struct {
	Name           string
	id             string
	record         common.Record
	ExpectedResult error
}

func TestCanView(t *testing.T) {

	uuidStr := uuid.New().String()

	scenarios := []scenario{
		{
			Name: "Empty Id",
			id:   "",
			record: common.Record{
				Owner: "owner",
			},
			ExpectedResult: fmt.Errorf("user is not allowed to view this record"),
		},
		{
			Name:           "Empty record",
			id:             "id",
			record:         common.Record{},
			ExpectedResult: fmt.Errorf("user is not allowed to view this record"),
		},
		{
			Name: "Owner requests record",
			id:   uuidStr,
			record: common.Record{
				Owner: uuidStr,
			},
			ExpectedResult: nil,
		},
		{
			Name: "Allowed user requests record",
			id:   uuidStr,
			record: common.Record{
				Owner:        "",
				AllowedUsers: []string{uuidStr},
			},
			ExpectedResult: nil,
		},
		{
			Name: "Record is public",
			id:   uuidStr,
			record: common.Record{
				IsPublic: true,
			},
			ExpectedResult: nil,
		},
	}

	for _, scene := range scenarios {
		t.Run(scene.Name, func(t *testing.T) {
			testResult := canView(scene.id, scene.record)
			if scene.ExpectedResult == nil {
				if testResult != nil {
					t.Errorf("expected nil got %s\n", testResult.Error())
				}
			} else {
				if testResult.Error() != scene.ExpectedResult.Error() {
					t.Errorf("Expected \"%s\" got \"%s\"", scene.ExpectedResult, testResult.Error())
				}
			}
		})
	}
}

func TestCheckOwner(t *testing.T) {

	uuidStr := uuid.New().String()

	scenarios := []scenario{
		{
			Name: "Empty Id",
			id:   "",
			record: common.Record{
				Owner: uuidStr,
			},
			ExpectedResult: fmt.Errorf("User is not authorized"),
		},
		{
			Name:           "Empty record",
			id:             uuidStr,
			record:         common.Record{},
			ExpectedResult: fmt.Errorf("User is not authorized"),
		},
		{
			Name: "Owner requests record",
			id:   uuidStr,
			record: common.Record{
				Owner: uuidStr,
			},
			ExpectedResult: nil,
		},
		{
			Name: "Allowed user requests record",
			id:   uuidStr,
			record: common.Record{
				Owner:        "",
				AllowedUsers: []string{uuidStr},
			},
			ExpectedResult: fmt.Errorf("User is not authorized"),
		},
	}

	for _, scene := range scenarios {
		t.Run(scene.Name, func(t *testing.T) {
			testResult := checkOwner(scene.id, scene.record)
			if scene.ExpectedResult == nil {
				if testResult != nil {
					t.Errorf("expected nil got %s\n", testResult.Error())
				}
			} else {
				if testResult.Error() != scene.ExpectedResult.Error() {
					t.Errorf("Expected \"%s\" got \"%s\"", scene.ExpectedResult, testResult.Error())
				}
			}
		})
	}
}

type mockDB struct{}

func (m mockDB) GetRecord(key string) (*common.Record, error) {
	if key == "key" {
		return &common.Record{Key: "key"}, nil
	} else {
		return &common.Record{}, nil
	}
}

func (m mockDB) DeleteRecord(key string) error             { return nil }
func (m mockDB) CreateRecord(record common.Record) error   { return nil }
func (m mockDB) ReplaceRecord(record *common.Record) error { return nil }
func TestRecordExits(t *testing.T) {

	mock := mockDB{}

	scenarios := []struct {
		Name           string
		Key            string
		Db             recordDataBase
		ExpectedRecord *common.Record
		ExpectedError  error
	}{
		{
			Name:           "Record does not exist",
			Key:            "not key",
			Db:             mock,
			ExpectedRecord: &common.Record{},
			ExpectedError:  nil,
		},
		{
			Name:           "Record exists with matching key",
			Key:            "key",
			Db:             mock,
			ExpectedRecord: &common.Record{Key: "key"},
			ExpectedError:  nil,
		},
	}

	for _, scene := range scenarios {
		t.Run(scene.Name, func(t *testing.T) {
			testResult, testErr := recordExists(scene.Key, scene.Db)
			if testResult.Key != scene.ExpectedRecord.Key {
				t.Errorf("incorrect record: expected %s got %s \n", scene.ExpectedRecord.Key, testResult.Key)
			}

			if scene.ExpectedError == nil {
				if testErr != nil {
					t.Errorf("expected nil error got %s\n", testErr.Error())
				}
			} else {
				if testErr.Error() != scene.ExpectedError.Error() {
					t.Errorf("incorrect error: expected %s got %s\n", scene.ExpectedError.Error(), testErr.Error())
				}
			}
		})
	}
}
