<p align="center"><img width="150px" src="https://blog.golang.org/go-brand/Go-Logo/PNG/Go-Logo_Blue.png"></p>

[![GoDoc](https://godoc.org/github.com/wilsontwm/user-registration?status.svg)](https://godoc.org/github.com/wilsontwm/user-registration)

## About

The user-registration module is a module that uses Golang on user registration and authentication. This is created to simplify the codes to sign up and login users. By integrating the module, it takes out the pain of development by easing common tasks used in majority of web projects, such as:

- User registration
- User authentication
- User forget password
- User reset password
- User activation

## Example

1. Database setup
```go
dbConfig := userreg.DBConfig{
  Driver:   os.Getenv("db_type"),
  Username: os.Getenv("db_user"),
  Password: os.Getenv("db_pass"),
  Host:     os.Getenv("db_host"),
  DBName:   os.Getenv("db_name"),
}

tableName := "tests"
userreg.Initialize(dbConfig)
userreg.Config(userreg.TableName(tableName))
```

2. User registration
```go
in := &userreg.User{}
in.Name = input.Name
in.Email = input.Email
in.Password = input.Password

// Signup the account
user, err := userreg.Signup(in)
```

3. User login
```go
in := &userreg.User{}
in.Email = input.Email
in.Password = input.Password

// Login the account
user, err := userreg.Login(in)
```

4. User forget password
```go
in := &userreg.User{}
in.Email = input.Email
user, err := userreg.ForgetPassword(in)
```

5. User reset password
```go
in := &userreg.User{}
in.ResetPasswordCode = &input.ResetPasswordCode
in.Password = input.Password
user, err := userreg.ResetPassword(in)
```

6. User activate account
```go
in := &userreg.User{}
in.ActivationCode = &input.ActivationCode
user, err := userreg.ActivateAccount(in)
```

## License

The user-registration module is open-sourced software licensed under the [MIT license](http://opensource.org/licenses/MIT).
