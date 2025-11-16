package main

import (
	"fmt"
	"net/http"

	"github.com/CosmicBDry/k8s-admissionReview/common"
	"github.com/CosmicBDry/k8s-admissionReview/webhookmutate"
	"github.com/CosmicBDry/k8s-admissionReview/webhookvalidate"
	"github.com/gin-gonic/gin"
)

var ValidInterfaceIns common.AdmissionReviewValidInterface

var MutateInterfaceIns common.AdmissionReviewMutateInterface

func main() {
	//gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "success",
			"status":  http.StatusOK,
		})
	})

	r.POST("/validate", func(c *gin.Context) {

		ValidInterfaceIns = webhookvalidate.NewAdmissionReviewValidate()

		ReponseAdmissionReview, err := ValidInterfaceIns.Validating(c)
		if err != nil {
			fmt.Println(err)
		}
		c.JSON(http.StatusOK, ReponseAdmissionReview)

	})

	r.POST("/mutate", func(c *gin.Context) {
		MutateInterfaceIns := webhookmutate.NewAdmissionReviewMutate()
		ReponseAdmissionReview, err := MutateInterfaceIns.Mutating(c)
		if err != nil {
			fmt.Println(err)
		}
		c.JSON(http.StatusOK, ReponseAdmissionReview)

	})

	//r.Run(":9090")
	r.RunTLS(":8080", "./certs/webserver.crt", "./certs/webserver.key")

}
