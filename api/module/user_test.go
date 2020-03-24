package module

import (
	"testing"
)

func clearTables() {
	dbConn.Exec("truncate missops_user")
}

func TestMain(m *testing.M) {
	clearTables()
	m.Run()
	clearTables()
}
func TestUserWork(t *testing.T) {
	t.Run("ADD", testAddUser)
	t.Run("GET", testGetUser)
	t.Run("DEL", testDelUser)
	t.Run("REGET", testRegetUser)
}

func testAddUser(t *testing.T) {
	err := AddUserCredential("missops", "123456")
	if err != nil {
		t.Errorf("Error of user add: %v ", err)
	}
}

func testGetUser(t *testing.T) {
	pwd, err := GetUserCredential("missops")
	if pwd != "123456" || err != nil {
		t.Errorf("Error of user get: %v ", err)
	}

}

func testDelUser(t *testing.T) {
	err := DeleteUser("missops", "123456")
	if err != nil {
		t.Errorf("Error of user del: %v ", err)
	}
}

func testRegetUser(t *testing.T) {
	_, err := GetUserCredential("missops")
	if err != nil {
		t.Errorf("Error of user get: %v ", err)
	}
}
