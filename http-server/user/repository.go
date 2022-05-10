package user

import (
	"fmt"
	"sync"
)

type UserId string

var (
	UserNotFound      = fmt.Errorf("user could not be found")
	UserAlreadyExists = fmt.Errorf("user alread exists")
)

type User struct {
	Id               UserId
	IdentityVerified bool
	PassbaseKey      *string
}

type userRepository struct {
	mutex sync.RWMutex
	users map[UserId]User
}

func NewUserRepository() userRepository {
	return userRepository{
		mutex: sync.RWMutex{},
		users: make(map[UserId]User),
	}
}

func (u *userRepository) CreateIfNotExist(userId UserId) (User, error) {
	newUser := User{
		Id:               userId,
		IdentityVerified: false,
		PassbaseKey:      nil,
	}
	u.mutex.Lock()
	defer u.mutex.Unlock()

	if user, exists := u.users[userId]; exists {
		return user, nil
	}

	u.users[userId] = newUser

	return newUser, nil
}

func (u *userRepository) GetUserById(userId UserId) (*User, error) {
	u.mutex.RLock()
	defer u.mutex.RUnlock()
	if user, exists := u.users[userId]; exists {
		return &user, nil
	}
	return nil, UserNotFound
}

func (u *userRepository) GetUserFromPassbaseKey(passbaseKey string) (*User, error) {
	u.mutex.RLock()
	defer u.mutex.RUnlock()
	for _, user := range u.users {
		if user.PassbaseKey != nil && *user.PassbaseKey == passbaseKey {
			return &user, nil
		}
	}
	return nil, UserNotFound
}

func (u *userRepository) RegisterUserVerified(userId UserId) error {
	u.mutex.RLock()
	defer u.mutex.RUnlock()
	if user, exists := u.users[userId]; exists {
		user.IdentityVerified = true
		u.users[userId] = user
		return nil
	}
	return UserNotFound
}

func (u *userRepository) AssociatePassbaseKey(userId UserId, passbaseKey string) error {
	u.mutex.Lock()
	defer u.mutex.Unlock()
	if user, exists := u.users[userId]; exists {
		user.PassbaseKey = &passbaseKey
		u.users[userId] = user
	} else {
		return UserNotFound
	}

	return nil
}

func (u *userRepository) GetAllUsers() []User {
	u.mutex.RLock()
	defer u.mutex.RUnlock()

	users := []User{}

	for _, user := range u.users {
		users = append(users, user)
	}

	return users
}
