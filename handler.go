package main

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo"
	"github.com/thongtiger/gostack/util"
)

// Handlers : user
func createUser(c echo.Context) (err error) {
	payload := struct {
		Username string   `json:"username" validate:"required"`
		Password string   `json:"password" validate:"required"`
		Role     string   `json:"role"`
		Scope    []string `json:"scope"`
	}{}
	if err = c.Bind(&payload); err != nil {
		return c.JSON(http.StatusUnsupportedMediaType, echo.Map{"message": "invalid format"})
	}
	if err = c.Validate(payload); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"message": "invalid parameter"})
	}
	if !util.UsernameValid(payload.Username) {
		return c.JSON(http.StatusBadRequest, echo.Map{"message": "username is not correct format"})
	}
	if !util.PasswordValid(payload.Password) {
		return c.JSON(http.StatusBadRequest, echo.Map{"message": "password is not correct format"})
	}
	if user, _ := st.GetUser(payload.Username); user != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"message": fmt.Sprintf("%s already exists", payload.Username)})
	}
	currentUser, err := st.NewUser(payload.Username, payload.Password, payload.Role, payload.Scope)
	if err != nil {
		return
	}
	// (*currentUser).Password = ""
	return c.JSON(http.StatusCreated, echo.Map{"message": "successfully", "result": currentUser})
}

func findUser(c echo.Context) error {
	users := st.FindUser()
	return c.JSON(200, users)
}

func updateUser(c echo.Context) error {
	return c.NoContent(200)
}

func deleteUser(c echo.Context) error {
	return c.NoContent(200)
}
