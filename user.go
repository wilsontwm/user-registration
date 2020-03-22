package user

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
}

// Setting the table name
func (User) TableName() string {
	return tableName
}

// Sign up of the user
func (user *User) Signup() error {
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	user.Password = string(hashedPassword)

	db := getDB()
	defer db.Close()

	// Create the user
	db.Create(user)

	if user.ID == uuid.Nil {
		return fmt.Errorf("User is not created.")
	}

	// Store the activation code to the user
	hash := md5.New()
	hash.Write([]byte(fmt.Sprint(user.ID)))
	activationCode := hex.EncodeToString(hash.Sum(nil))

	db.Model(&user).Update("ActivationCode", activationCode)

	user.Password = "" // delete the password

	return nil
}

// Resend the activation code to the user
func (user *User) ResendActivation() error {
	user = getUserByEmail(user.Email)

	if user == nil {
		return fmt.Errorf("Invalid email address.")
	} else if user.ActivationCode == nil {
		return fmt.Errorf("User has already been activated.")
	}

	// TODO: Send email to resend activation code

	return nil
}

// Set forgotten password code
func (user *User) ForgetPassword() error {
	user = getUserByEmail(user.Email)

	if user == nil {
		return fmt.Errorf("Invalid email address.")
	} else if user.ActivationCode != nil {
		return fmt.Errorf("Please activate your account first.")
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

	return nil
}

// Activate user account
func (user *User) ActivateAccount() error {
	user = getUserByActivationCode(*user.ActivationCode)

	if user == nil {
		return fmt.Errorf("Invalid activation link.")
	}

	db := getDB()
	defer db.Close()

	db.Model(&user).Update("ActivationCode", nil)

	return nil
}

// Reset user password
func (user *User) ResetPassword() error {
	user = getUserByResetPasswordCode(*user.ResetPasswordCode)

	if user == nil {
		return fmt.Errorf("Invalid reset password link.")
	}

	// Reset the password of the user
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	db := getDB()
	defer db.Close()

	db.Model(&user).Update(map[string]interface{}{
		"ResetPasswordCode":      nil,
		"ResetPasswordExpiredAt": nil,
		"Password":               string(hashedPassword),
	})

	return nil
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
