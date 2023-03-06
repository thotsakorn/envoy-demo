package global

import "os"

const (
	dburi       = "mongodb+srv://sunmongo:passw0rd@cluster0.r2cfi.mongodb.net/?retryWrites=true&w=majority"
	dbname      = "store"
	performance = 100
)

var (
	jwtSecret = []byte(os.Getenv("JWT_SECRET_KEY"))
)
