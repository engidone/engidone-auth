package users

type User struct {
	ID       string `yaml:"id"`
	Username string    `yaml:"username"`
	Email    string    `yaml:"email"`
}

type UserConfig struct {
	Users []User `yaml:"users"`
}