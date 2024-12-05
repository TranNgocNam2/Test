package certificategrp

import (
	"Backend/business/core/learner/certificate"
	"Backend/internal/app"
	"github.com/gin-gonic/gin"
)

func CertificateRoutes(router *gin.Engine, app *app.Application) {
	certificateCore := certificate.NewCore(app)
	handlers := New(certificateCore)
	certificates := router.Group("/certificates")
	{
		certificates.GET("/:id", handlers.GetCertificateById())
	}
}
