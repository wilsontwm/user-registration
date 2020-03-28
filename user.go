package userreg

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
	"time"
)

const (
	formattedDate = "01/02/2006"
)

type Token struct {
	UserID uuid.UUID
	Name   string
	Email  string
	Expiry time.Time
	jwt.StandardClaims
}

type User struct {
	BaseModel
	Name                   string `gorm:"not null"`
	Email                  string `gorm:"unique;not null"`
	Password               string `gorm:"not null"`
	ActivationCode         *string
	ResetPasswordCode      *string
	ResetPasswordExpiredAt *time.Time
	Token                  string `gorm:"-"`
}

// Setting the table name
func (User) TableName() string {
	return tableName
}

// Login of the user
func Login(input *User) (*User, error) {
	db := getDB()
	defer db.Close()
	user := &User{}
	db.Table(tableName).Where("email = ?", input.Email).First(user)

	if user == nil {
		return nil, fmt.Errorf("Invalid email or password.")
	} else if user.ActivationCode != nil {
		return nil, fmt.Errorf("User account is not activated.")
	}

	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password))

	// If password does not match
	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
		return nil, fmt.Errorf("Invalid email or password.")
	} else if err != nil {
		return nil, fmt.Errorf(err.Error())
	}

	// Create new JWT token for the newly registered account
	expiry := time.Now().Add(time.Hour * 2) // Only valid for 2 hours
	tk := &Token{
		UserID: user.ID,
		Name:   user.Name,
		Email:  user.Email,
		Expiry: expiry,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, tk)
	tokenString, _ := token.SignedString([]byte(jwtKey))
	user.Token = tokenString

	return user, nil
}

// Sign up of the user
func Signup(input *User) (*User, error) {
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	input.Password = string(hashedPassword)
	user := input
	db := getDB()
	defer db.Close()

	temp := getUserByEmail(input.Email)
	if temp != nil {
		return nil, fmt.Errorf("Email has already been taken.")
	}

	// Create the user
	db.Create(user)

	if user.ID == uuid.Nil {
		return nil, fmt.Errorf("User is not created.")
	}

	// Store the activation code to the user
	if isUserActivationRequired {
		hash := md5.New()
		hash.Write([]byte(fmt.Sprint(user.ID)))
		activationCode := hex.EncodeToString(hash.Sum(nil))

		db.Model(&user).Update("ActivationCode", activationCode)
	}

	user.Password = "" // delete the password

	return user, nil
}

// Set forgotten password code
func ForgetPassword(input *User) (*User, error) {
	user := getUserByEmail(input.Email)

	if user == nil {
		return nil, fmt.Errorf("Invalid email address.")
	} else if user.ActivationCode != nil {
		return nil, fmt.Errorf("Please activate your account first.")
	}

	// Store the reset password code to the user
	hash := md5.New()
	hash.Write([]byte(fmt.Sprint(user.ID) + time.Now().String()))
	resetPasswordCode := hex.EncodeToString(hash.Sum(nil))
	// Add one hour to the expiry date for reseting the password
	resetPasswordExpiredAt := time.Now().Local().Add(time.Hour * 1)

	db := getDB()
	defer db.Close()

	db.Model(&user).Update(map[string]interface{}{
		"ResetPasswordCode":      resetPasswordCode,
		"ResetPasswordExpiredAt": resetPasswordExpiredAt,
	})

	return user, nil
}

// Activate user account
func ActivateAccount(input *User) (*User, error) {
	user := getUserByActivationCode(*input.ActivationCode)

	if user == nil {
		return nil, fmt.Errorf("Invalid activation link.")
	}

	db := getDB()
	defer db.Close()

	db.Model(&user).Update("ActivationCode", nil)

	return user, nil
}

// Reset user password
func ResetPassword(input *User) (*User, error) {
	user := getUserByResetPasswordCode(*input.ResetPasswordCode)

	if user == nil {
		return nil, fmt.Errorf("Invalid reset password link.")
	}

	// Reset the password of the user
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	db := getDB()
	defer db.Close()

	db.Model(&user).Update(map[string]interface{}{
		"ResetPasswordCode":      nil,
		"ResetPasswordExpiredAt": nil,
		"Password":               string(hashedPassword),
	})

	return user, nil
}

// Authenticate the user by token
func Authenticate(jwtToken string) (*Token, error) {
	tk := &Token{}
	token, err := jwt.ParseWithClaims(jwtToken, tk, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtKey), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, fmt.Errorf("Token is not valid.")
	}

	if time.Now().After(tk.Expiry) {
		return nil, fmt.Errorf("Token has expired. Please login again.")
	}

	return tk, nil
}

// Post processing of the user
func getUser(user *User) *User {
	if user.Email == "" {
		return nil
	}

	user.Password = ""
	return user
}

// Get the user by email
func getUserByEmail(email string) *User {
	user := &User{}
	db := getDB()
	defer db.Close()

	db.Select("*").
		Where("email = ?", email).
		First(user)

	return getUser(user)
}

// Get the user by activation code
func getUserByActivationCode(code string) *User {
	user := &User{}
	db := getDB()
	defer db.Close()

	db.Select("*").
		Where("activation_code = ?", code).
		First(user)

	return getUser(user)
}

// Get the user by reset password code (that has not expired)
func getUserByResetPasswordCode(resetPasswordCode string) *User {
	user := &User{}
	now := time.Now().Local()
	db := getDB()
	defer db.Close()

	db.Select("*").
		Where("reset_password_code = ?", resetPasswordCode).
		Where("reset_password_expired_at > ?", now).
		First(user)

	return getUser(user)
}
