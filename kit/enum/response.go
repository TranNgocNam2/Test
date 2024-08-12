package enum

const (
	SuccessCode           = "00"
	ClientErrorCode       = "01"
	ResourceErrorCode     = "02"
	UnAuthorizedErrorCode = "03"
	NotFoundErrorCode     = "04"
	SystemErrorCode       = "99"

	SuccessMessage           = "Success"
	SystemErrorMessage       = "System Error"
	ClientErrorMessage       = "Invalid Request"
	ResourceErrorMessage     = "Could not get resource !"
	UnAuthorizedErrorMessage = "User is not authorized !"
)
