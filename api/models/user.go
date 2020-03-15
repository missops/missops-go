package models

//AddUserCredential : insert user to databases
func AddUserCredential(userName string, pwd string) error {
	stmtIn, err := dbConn.Prepare("INSERT INTO missops_user (user_name,user_pwd) VALUES (?,?)")
	if err != nil {
		return err
	}
	defer stmtIn.Close()
	stmtIn.Exec(userName, pwd)

	return nil
}

//GetUserCredential : select pwd from databases
func GetUserCredential(userName string) (string, error) {
	stmtOut, err := dbConn.Prepare("SELECT user_pwd FROM missops_user WHERE user_name = ?  ")
	if err != nil {
		return "", err
	}
	defer stmtOut.Close()
	var pwd string
	stmtOut.QueryRow(userName).Scan(&pwd)
	return pwd, nil
}

//DeleteUser : delete user
func DeleteUser(userName string, pwd string) error {
	stmtDel, err := dbConn.Prepare("DELETE FROM  missops_user WHERE user_name = ? and user_pwd = ?")
	if err != nil {
		return err
	}
	defer stmtDel.Close()
	stmtDel.Exec(userName, pwd)
	return nil
}
