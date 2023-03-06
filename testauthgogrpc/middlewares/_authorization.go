package middlewares

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
)

var validateClient rpc_auth.ValidateTokenServiceClient

func InitMiddleware() *data.errorDetail {
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithInsecure())
	conn, err := grpc.Dial("localhost:9000", opts...)
	if err != nil {
		return &data.errorDetail{
			errorMessage: fmt.Sprintf("error in dial %+v", err),
			errorType:    data.ErrorTyeFatal,
		}
	}
	validateClient = rpc_auth.NewValidateTokenServiceClient(conn)
	return nil
}

func AuthenticationMiddleware() gin.HandlerFunc {
	return func(gContext *gin.Context) {
		tokenString := gContext.Request.Header.Get("apikey")

		if tokenString == "" {
			returnUnauthorized(gContext)
		}

		fmt.Printf("token --%s--", tokenString)
		response, err := validateClient.simpleValidate(context.Background(), &rpc_auth.ValidateTokenRequest{
			Token: tokenString,
		})

		if err != nil {
			returnUnauthorized(gContext)
			return
		}
		gContext.Keys["CompanyId"] = response.CompanyId
		gContext.Keys["Username"] = response.Username
		gContext.Keys["Roles"] = response.Roles
	}
}

func returnUnauthorized(context *gin.Context) {
	context.AbortWithStatusJSON(http.StatusUnauthorized, data.Response{
		Error: []data.errorDetail{
			{
				errorType:    data.ErrorTypeUnauthorized,
				errorMessage: "You are not authorized to access this path",
			},
		},
		Status:  http.StatusUnauthorized,
		Message: "Unauthorized access",
	})
}
