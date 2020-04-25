package main

import (
	"fmt"
	"net/http"

	"github.com/local/store"
	"github.com/thongtiger/oauth-rfc6749/auth"

	"github.com/labstack/echo/middleware"

	"github.com/labstack/echo"
	"github.com/thongtiger/gostack/util"
	"github.com/thongtiger/oauth-rfc6749/handle"
	"gopkg.in/go-playground/validator.v9"
)

var validate *validator.Validate
var conf Config
var st store.Store

// Validator type of validator
type CustomValidator struct{ validator *validator.Validate }

// Validate is a context validator
func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}

func init() {
	conf.Read()
	st = store.NewMongoStore(conf.Mongodb.Host, conf.Mongodb.Port, conf.Mongodb.Username, conf.Mongodb.Password, conf.Mongodb.Database)
}
func main() {
	e := echo.New()
	e.Validator = &CustomValidator{validator: validator.New()}
	e.Use(
		middleware.Recover(),
		middleware.Secure(),
		middleware.Logger(),
		middleware.Gzip(),
		middleware.BodyLimit("2M"),
		middleware.CORSWithConfig(middleware.CORSConfig{
			AllowOrigins: []string{"*"},
			AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentLength, echo.HeaderAcceptEncoding, echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization, echo.HeaderXRequestedWith},
			AllowMethods: []string{http.MethodGet, http.MethodHead, http.MethodPut, http.MethodPatch, http.MethodPost, http.MethodDelete, http.MethodOptions},
			MaxAge:       3600,
		}),
	)
	// ------------- public  -------------
	e.GET("/", func(c echo.Context) error { return c.String(http.StatusOK, "Hello, World!") })
	e.GET("/401", func(c echo.Context) error { return echo.ErrUnauthorized })

	// ------------- protected  -------------
	e.POST("/oauth2/token", func(c echo.Context) (err error) {
		oauth2 := auth.Oauth2{}
		if err = c.Bind(&oauth2); err != nil {
			return c.JSON(http.StatusUnsupportedMediaType, echo.Map{})
		}
		switch oauth2.GrantType {
		case "password":
			if ok, user := st.ValidateUser(oauth2.Username, oauth2.Password); ok {
				// generate token
				return handle.GenerateTK(c, user)
			}
		case "refresh_token":
			return handle.RefreshTK(c, oauth2)
		}
		return c.JSON(http.StatusUnauthorized, echo.Map{})
	})
	// for admin
	e.POST("/newuser", func(c echo.Context) (err error) {
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
			return c.JSON(http.StatusBadRequest, echo.Map{"message": fmt.Sprintf("%s already exists", user.Username)})
		}
		currentUser, err := st.NewUser(payload.Username, payload.Password, payload.Role, payload.Scope)
		if err != nil {
			return
		}
		(*currentUser).Password = ""
		return c.JSON(http.StatusCreated, echo.Map{"message": "successfully", "result": currentUser})
	}, auth.AcceptedRole("ADMIN"))

	e.GET("/logout", handle.LogoutHandle, auth.JWTMiddleware())
	e.GET("/protected", func(c echo.Context) error { return c.String(http.StatusOK, "allow protected") }, auth.JWTMiddleware())
	e.Logger.Fatal(e.Start(":1323"))

}
