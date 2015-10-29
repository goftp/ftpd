package web

type User struct {
	Name string
	Pass string
}

type UserDB interface {
	GetUser(string) (string, error)
	AddUser(string, string) error
	DelUser(string) error
	ChgPass(string, string) error
	UserList(*[]User) error

	AddGroup(string) error
	DelGroup(string) error
	GroupList(*[]string) error

	AddUserGroup(string, string) error
	DelUserGroup(string, string) error
	GroupUser(string, *[]string) error
}
