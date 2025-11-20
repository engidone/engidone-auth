package users

type User struct {
	ID       int32  `yaml:"id"`
	Username string `yaml:"username"`
	Email    string `yaml:"email"`
}
