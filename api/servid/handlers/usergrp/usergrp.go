package usergrp

//import (
//	"Backend/business/core/school"
//	"Backend/business/core/user"
//	"Backend/internal/web"
//	"github.com/gin-gonic/gin"
//	"net/http"
//)
//
//type Handlers struct {
//	user   *user.Core
//	school *school.Core
//}
//
//func (h *Handlers) CreateUser() gin.HandlerFunc {
//	return func(ctx *gin.Context) {
//		var webNewUser WebNewUser
//		if err := web.Decode(ctx, &webNewUser); err != nil {
//			web.Respond(ctx, nil, http.StatusBadRequest, err)
//			return
//		}
//
//		newSchool := toCoreNewUser(webNewUser)
//		err, status := toCore(ctx, newSchool)
//		if err != nil {
//			web.Respond(ctx, nil, status, err)
//			return
//		}
//
//		web.Respond(ctx, nil, http.StatusOK, nil)
//	}
//}
