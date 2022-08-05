package repository

import "github.com/ditrit/badaas/persistence/models"

// This is the list which is used to store user sessions
var AuthenticatedUsers []*models.User

// Returns the list of all users who logged in. Be careful, there could be some users with expired tokens in this list
func GetUsers() []*models.User {
	return AuthenticatedUsers
}

// Add a user to the list of user
func AddUser(u *models.User) {
	AuthenticatedUsers = append(AuthenticatedUsers, u)
}

func ReplaceAllUsers(l []*models.User) {
	AuthenticatedUsers = l
}
