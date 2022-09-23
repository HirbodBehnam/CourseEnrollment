package AuthCore

import (
	"github.com/golang-jwt/jwt/v4"
	"time"
)

// signingMethod is the method which we sign all jwt tokens with
var signingMethod = jwt.SigningMethodHS256

// jwtTTL is the time interval which jwt token is valid
const jwtTTL = time.Minute * 5

// jwtIssuer which is "course enrollment auth"
const jwtIssuer = "cea"
