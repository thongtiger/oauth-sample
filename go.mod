module oauth-sample

go 1.13

require (
	github.com/labstack/echo v3.3.10+incompatible
	github.com/local v0.0.0-00010101000000-000000000000
	github.com/thongtiger/gostack v0.0.0-20200420022901-9bb9b97ff9d0
	github.com/thongtiger/oauth-rfc6749 v0.0.0-20200422074005-7949868687ab
	go.mongodb.org/mongo-driver v1.3.2
	gopkg.in/go-playground/validator.v9 v9.31.0
)

replace github.com/local => ./
