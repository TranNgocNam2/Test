package post

import (
	"Backend/api/internal/platform/db/sqlc"
	"Backend/kit/enum"
	"Backend/kit/web"
	"github.com/gin-gonic/gin"
)

func GetAllPost(queries *db.Queries) gin.HandlerFunc {
	return func(context *gin.Context) {

		pageNumber := 1
		pageSize := 10
		keyWord := ""
		sortBy := "title"
		sortOrder := "DESC"
		post, err := queries.GetPostsByUserAndTags(context.Request.Context(), db.GetPostsByUserAndTagsParams{
			Size:       int32(pageSize),
			Page:       int32(pageNumber),
			SortBy:     sortBy,
			SortOrder:  sortOrder,
			UserID:     nil,
			SearchTerm: keyWord,
			TagIds:     nil,
			IsBuyable:  false,
		})
		if err != nil {
			web.SystemError(context, err)
			return
		}

		response := web.BaseResponse{
			ResultCode:    enum.SuccessCode,
			ResultMessage: enum.SuccessMessage,
			Data:          post,
		}
		context.JSON(200, response)
	}
}
