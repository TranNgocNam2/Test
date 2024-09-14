package middleware

//
//import (
//	"Backend/internal/platform/enum"
//	"fmt"
//	"github.com/gin-gonic/gin"
//	"gitlab.com/innovia69420/kit/web"
//)
//
//func RecoverPanic() gin.HandlerFunc {
//	return func(c *gin.Context) {
//		defer func() {
//			if err := recover(); err != nil {
//				c.Header(enum.HttpContentType, enum.HttpJson)
//				web.SystemError(c, fmt.Errorf("%v", err))
//			}
//		}()
//
//		c.Next()
//	}
//}
