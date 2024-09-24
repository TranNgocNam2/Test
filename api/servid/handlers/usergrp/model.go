package usergrp

//type WebUser struct {
//	ID       string               `json:"id"`
//	FullName string               `json:"fullName"`
//	Email    string               `json:"email"`
//	Phone    string               `json:"phone"`
//	Gender   string               `json:"gender"`
//	Role     string               `json:"-"`
//	Photo    string               `json:"photo"`
//	School   *schoolgrp.WebSchool `json:"school,omitempty"`
//}

//func toWebUser(dbUser sqlc.User, dbSchool) WebUser {
//	return WebUser{
//		ID:       dbUser.ID,
//		FullName: dbUser.FullName,
//		Email:    dbUser.Email,
//		Phone:    dbUser.Phone,
//		Gender:   user.GetGenderStr(dbUser.Gender),
//		Role:     user.GetRoleName(dbUser.Role),
//		Photo:    dbUser.ProfilePhoto,
//		School:   schoolgrp.ToWebSchool(dbSchool),
//	}
//}
//
//func toWebUsers(dbUsers []sqlc.User) []WebUser {
//	users := make([]WebUser, len(dbUsers))
//	for i, dbUser := range dbUsers {
//		users[i] = toWebUser(dbUser)
//	}
//	return users
//}
//
//type WebNewUser struct {
//	ID       string `json:"id" validate:"required"`
//	FullName string `json:"fullName" validate:"required"`
//}
//
//func toCoreNewUser(webNewUser WebNewUser) sqlc.User {
//	return sqlc.User{
//		ID:       webNewUser.ID,
//		FullName: webNewUser.FullName,
//	}
//}
