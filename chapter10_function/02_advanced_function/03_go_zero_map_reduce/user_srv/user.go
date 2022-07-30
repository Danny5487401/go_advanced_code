package user_srv

import "time"

/********用户服务**********/

type User struct {
	Name string
	Age  uint8
}

func GetUser() (*User, error) {
	time.Sleep(500 * time.Millisecond)
	var u User
	u.Name = "Danny"
	u.Age = 18
	return &u, nil
}
