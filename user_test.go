package userreg

import (
	"github.com/satori/go.uuid"
	"strings"
	"testing"
)

// Initialize step to set the database configuration
func initialize() {
	tableName := "tests"
	// Mysql
	dbConfig := DBConfig{
		Driver:   Mysql,
		Username: "admin",
		Password: "password",
		Host:     "localhost",
		DBName:   "test",
	}

	Initialize(dbConfig)

	Config(TableName(tableName))
}

// Test signing up the user
func TestSignup(t *testing.T) {
	initialize()
	Config(UserActivation(true))

	input := &User{
		Email:    "test@gmail.com",
		Name:     "test",
		Password: "password",
	}

	user, err := Signup(input)

	if err != nil {
		t.Error(err)
	} else if user.ID == uuid.Nil {
		t.Error("User 1 ID is empty, user is not created.")
	} else if user.ActivationCode == nil {
		t.Error("User 1 activation code is not set.")
	}

	input = &User{
		Email:    "test@gmail.com",
		Name:     "test 2",
		Password: "password",
	}

	user, err = Signup(input)

	if err == nil {
		t.Error("User with same email has been created.")
	}

	input = &User{
		Email:    "test99@gmail.com",
		Name:     "test99",
		Password: "password",
	}

	user, err = Signup(input)

	if err != nil {
		t.Error(err)
	} else if user.ID == uuid.Nil {
		t.Error("User 99 ID is empty, user is not created.")
	} else if user.ActivationCode == nil {
		t.Error("User 99 activation code is not set.")
	}

	Config(UserActivation(false))

	input = &User{
		Email:    "test2@gmail.com",
		Name:     "test2",
		Password: "password",
	}

	user, err = Signup(input)

	if err != nil {
		t.Error(err)
	} else if user.ID == uuid.Nil {
		t.Error("User 2 ID is empty, user is not created.")
	} else if user.ActivationCode != nil {
		t.Error("User 2 activation code is set.")
	}
}

// Test user forggoten password
func TestForgetPassword(t *testing.T) {
	initialize()
	input := &User{
		Email: "test@gmail.com",
	}

	user, err := ForgetPassword(input)

	if err == nil {
		t.Error("User 1 is not activated. No error is prompted")
	}

	input = &User{
		Email: "test2@gmail.com",
	}

	user, err = ForgetPassword(input)

	if err != nil {
		t.Error(err)
	} else if user.ResetPasswordCode == nil {
		t.Error("User 2 does not have reset password code set.")
	} else if user.ResetPasswordExpiredAt == nil {
		t.Error("User 2 does not have reset password code expiry set.")
	}
}

// Test user activate account
func TestActivateAccount(t *testing.T) {
	initialize()
	code := "this_is_random_code"
	input := &User{
		ActivationCode: &code,
	}

	user, err := ActivateAccount(input)

	if err == nil {
		t.Error("Activation code should not be found.")
	}

	db := getDB()
	defer db.Close()

	// Randomly get an user with activation code
	temp := &User{}
	db.Where("activation_code <> ?", "").Where("email = ?", "test@gmail.com").First(temp)

	if temp == nil {
		t.Error("User with activation code cannot be found.")
	} else {
		user, err = ActivateAccount(temp)

		if err != nil {
			t.Error(err)
		} else if user.ActivationCode != nil {
			t.Error("User is not activated.")
		}
	}
}

// Test user reset password
func TestResetPassword(t *testing.T) {
	initialize()
	db := getDB()
	defer db.Close()

	// Randomly get an user with reset password code
	temp := &User{}
	db.Where("reset_password_code <> ?", "").First(temp)

	if temp == nil {
		t.Error("User with reset password code cannot be found.")
		return
	}

	input := &User{
		Password:          "newpassword",
		ResetPasswordCode: temp.ResetPasswordCode,
	}

	user, err := ResetPassword(input)

	if err != nil {
		t.Error(err)
	} else if user.ResetPasswordCode != nil {
		t.Error("User reset password code is still active.")
	}
}

// Test get activation code
func TestGetActivationCode(t *testing.T) {
	initialize()
	db := getDB()
	defer db.Close()

	// Randomly get an user with activation code
	temp := &User{}
	db.Where("activation_code <> ?", "").First(temp)

	if temp == nil {
		t.Error("User with activation code cannot be found.")
		return
	}

	input := &User{
		Email: temp.Email,
	}

	user, err := GetActivationCode(input)

	if err != nil {
		t.Error(err)
	} else if user.ActivationCode == nil {
		t.Error("User activation code is not retrieved.")
	}

	// Randomly get an user without activation code
	temp = &User{}
	db.Where("activation_code = ?", nil).First(temp)

	if temp == nil {
		t.Error("User with activation code cannot be found.")
		return
	}

	input = &User{
		Email: temp.Email,
	}

	user, err = GetActivationCode(input)

	if err == nil {
		t.Error("User has already been activated and error should be prompted.")
	}
}

// Test logging in the user
func TestLogin(t *testing.T) {
	initialize()
	input := &User{
		Email:    "test@gmail.com",
		Password: "password",
	}

	user, err := Login(input)
	if err != nil {
		t.Error(err)
	} else if user == nil {
		t.Error("User login fails.")
	} else if user.Token == "" {
		t.Error("JWT Token is not set.")
	}

	input = &User{
		Email:    "test@gmail.com",
		Password: "wrongpassword",
	}

	user, err = Login(input)

	if err == nil {
		t.Error("Invalid login with incorrect password.")
	}

	input = &User{
		Email:    "wrongemail@gmail.com",
		Password: "wrongpassword",
	}

	user, err = Login(input)

	if err == nil {
		t.Error("Invalid login with incorrect password.")
	}

	input = &User{
		Email:    "test99@gmail.com",
		Password: "password",
	}

	user, err = Login(input)

	if err == nil {
		t.Error("Invalid login to inactive user account.")
	}
}

// Test authenticate the user
func TestAuthenticate(t *testing.T) {
	initialize()
	input := &User{
		Email:    "test@gmail.com",
		Password: "password",
	}

	user, err := Login(input)
	if err != nil {
		t.Error(err)
	} else if user == nil {
		t.Error("User login fails.")
	}

	token := user.Token

	authenticatedToken, err := Authenticate(token)

	if err != nil {
		t.Error(err)
	} else if user.ID != authenticatedToken.UserID {
		t.Error("Invalid authenticated user ID.")
	}

	// Authenticate with incorrect token
	token = token[:len(token)-5]
	values := []string{token, "test"}
	token = strings.Join(values, "")
	authenticatedToken, err = Authenticate(token)

	if err == nil {
		t.Error("Invalid successful authentication of user.")
	}
}
