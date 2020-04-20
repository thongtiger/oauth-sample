package main

import (
	"github.com/local/store"
	"github.com/thongtiger/oauth-rfc6749/auth"
	"net/http"

	"github.com/labstack/echo/middleware"

	"github.com/labstack/echo"
	"gopkg.in/go-playground/validator.v9"

	"github.com/thongtiger/oauth-rfc6749/handle"
)

var validate *validator.Validate
var conf Config
var st store.Store

// Validator type of validator
type Validator struct{ validator *validator.Validate }

// Validate is a context validator
func (cv *Validator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}

func init() {
	conf.Read()
	st = store.NewMongoStore(conf.Mongodb.Host, conf.Mongodb.Port, conf.Mongodb.Username, conf.Mongodb.Password, conf.Mongodb.Database)
}
func main() {
	e := echo.New()
	// e.Validator = &Validator{validator: validator.New()}
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
	e.GET("/logout", handle.LogoutHandle, auth.JWTMiddleware())
	e.GET("/protected", func(c echo.Context) error { return c.String(http.StatusOK, "allow protected") }, auth.JWTMiddleware())
	e.Logger.Fatal(e.Start(":1323"))
}
