package enum

const (
	SuccessCode           = "00"
	BadRequestCode        = "01"
	ResourceErrorCode     = "02"
	UnAuthorizedErrorCode = "03"
	NotFoundErrorCode     = "04"
	SystemErrorCode       = "99"

	SuccessMessage           = "Thành công"
	SystemErrorMessage       = "Lỗi hệ thống"
	BadRequestMessage        = "Yêu cầu không hợp lệ"
	ResourceErrorMessage     = "Could not get resource !"
	UnAuthorizedErrorMessage = "Người dùng không có quyền sử dụng chức năng này"
)
