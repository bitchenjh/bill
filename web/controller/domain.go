package controller

type User struct {
	UserName string
	Password string

	CmId string
	Acct string
}

var Users []User

func init() {
	u := User{"admin", "123456", "AAMID", "管理员"}
	u2 := User{"alice", "123456", "ACMID", "A公司"}
	u3 := User{"bob", "123456", "BCMID", "B公司"}
	u4 := User{"jack", "123456", "CCMID", "C公司"}

	Users = append(Users, u)
	Users = append(Users, u2)
	Users = append(Users, u3)
	Users = append(Users, u4)
}
