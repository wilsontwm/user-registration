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

- User registration
```go
in := &userreg.User{}
in.Name = input.Name
in.Email = input.Email
in.Password = input.Password

// Signup the account
user, err := userreg.Signup(in)
```

## License

The user-registration module is open-sourced software licensed under the [MIT license](http://opensource.org/licenses/MIT).