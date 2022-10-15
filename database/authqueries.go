package database

import (
	"database/sql"
	"fmt"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func ValidateLogIn(username string, password string) (int, bool) {
	sqlStatement := `SELECT userid, password FROM users WHERE username = $1;`
	var pw string
	var id int

	row := db.QueryRow(sqlStatement, username)
	switch err := row.Scan(&id, &pw); err {
	case sql.ErrNoRows:
		fmt.Println("Invalid credentials")
	case nil:
		if bcrypt.CompareHashAndPassword([]byte(pw), []byte(password)) == nil {
			return id, true
		} else {
			fmt.Println("Invalid credentials")
		}
	default:
		panic(err)
	}

	return 0, false
}

func SaveSession(s string, id int) bool {
	sqlStatement := `INSERT INTO sessions(session, userid, startdate) VALUES ($1, $2, NOW());`

	_, err := db.Exec(sqlStatement, s, id)
	if err != nil {
		panic(err)
	}

	return true
}

func ValidateSession(c *gin.Context) *User {
	var user User

	cookie, err := c.Cookie("session")
	if err != nil {
		fmt.Println(err)
		return nil
	}

	sqlStatement := `SELECT users.userid, username, creationdate, permissionlevel, profilepicture FROM users JOIN sessions ON users.userid = sessions.userid WHERE session=$1;`

	row := db.QueryRow(sqlStatement, cookie)
	switch err := row.Scan(&user.Userid, &user.Username, &user.CreationDate, &user.PermissionLevel, &user.ProfilePicture); err {
	case sql.ErrNoRows:
		fmt.Println("session not found")
		return nil
	case nil:
		fmt.Println("session found")
		return &user
	default:
		panic(err)
	}
}

func DeleteSession(s string) bool {
	sqlStatement := `DELETE FROM sessions WHERE session=$1;`

	_, err := db.Exec(sqlStatement, s)
	if err != nil {
		fmt.Println(err)
		return false
	}

	return true
}

func RegisterAccount(username string, password string) bool {
	sqlStatement := `INSERT INTO users(username, password, creationdate, permissionlevel) VALUES ($1, $2, NOW(), 0);`

	_, err := db.Exec(sqlStatement, username, password)
	if err != nil {
		fmt.Println(err)
		return false
	}

	return true
}

func CanAccessThread(user *User, threadid string) bool {
	sqlStatement := `SELECT permissionlevel FROM subcategories JOIN threads ON subcategories.subcategoryid = threads.subcategoryid
					WHERE threadid=$1;`

	var p int
	err := db.QueryRow(sqlStatement, threadid).Scan(&p)
	if err != nil {
		fmt.Println(err)
		return false
	}

	if user.PermissionLevel < p {
		return false
	}

	return true
}
