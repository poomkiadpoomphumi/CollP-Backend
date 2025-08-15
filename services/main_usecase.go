package services

/* import (
	"context"
	"collp-backend/connection"
	"log"
) */

type MainMenu struct {
	ID   int
	Name string
}

var mainmenu = []MainMenu{
	{ID: 1, Name: "Heeeee"},
	{ID: 2, Name: "Kuyyyy"},
}

func GetAllMainMenu() []MainMenu {
	return mainmenu
}

/* type UserUsecase struct { Db *db.OracleDB }

func NewUserUsecase(d *db.OracleDB) *UserUsecase {
	return &UserUsecase{Db: d}
}
// GetAllUsers ดึงข้อมูลจาก Oracle โดยส่ง SQL command ไปที่ connection
func (u *UserUsecase) GetAllUsers() ([]User, error) {
	ctx := context.Background()
	rows, err := u.Db.QuerySQL(ctx, "SELECT id, name FROM users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var user User
		err := rows.Scan(&user.ID, &user.Name)
		if err != nil {
			log.Println("Scan error:", err)
			continue
		}
		users = append(users, user)
	}
	return users, nil
} */
