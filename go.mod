module oauth-sample

go 1.13

require (
	github.com/labstack/echo v3.3.10+incompatible
	github.com/local v0.0.0-00010101000000-000000000000
	github.com/thongtiger/oauth-rfc6749 v0.0.0-20200420032149-d7fd7e9f270d
	go.mongodb.org/mongo-driver v1.3.2
	gopkg.in/go-playground/validator.v9 v9.31.0
)

replace github.com/local => ./
