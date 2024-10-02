package userdb

// import (
//
//	"Backend/business/core/user"
//	"Backend/business/db/sqlc"
//
// )
//func toCoreUser(dbUser sqlc.User) user.User {
//	return user.User{
//		ID:       dbUser.ID,
//		FullName: dbUser.FullName,
//		Email:    dbUser.Email,
//		Phone:    dbUser.Phone,
//		Gender:   dbUser.Gender,
//		Role:     dbUser.Role,
//		Photo:    dbUser.ProfilePhoto,
//	}
//}

//
//func toCoreUserSlice(dbUsers []sqlc.User) []user.User {
//	users := make([]user.User, len(dbUsers))
//	for i, dbUser := range dbUsers {
//		users[i] = toCoreUser(dbUser)
//	}
//	return users
//}
