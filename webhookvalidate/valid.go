package webhookvalidate

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/gin-gonic/gin"
	admissionV1 "k8s.io/api/admission/v1"
	coreV1 "k8s.io/api/core/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type AdmissionReviewValidate struct {
	Pod                coreV1.Pod
	K8sAdmissionReview admissionV1.AdmissionReview
}

func NewAdmissionReviewValidate() *AdmissionReviewValidate {
	return &AdmissionReviewValidate{}
}

func (AR *AdmissionReviewValidate) Validating(c *gin.Context) (*admissionV1.AdmissionReview, error) {
	reqBytes, _ := io.ReadAll(c.Request.Body)
	err := json.Unmarshal(reqBytes, &AR.K8sAdmissionReview)
	//fmt.Println(string(res))
	PodInstance := AR.Pod
	ReqInstance := AR.K8sAdmissionReview.Request
	ReponseAdmissionReview := &admissionV1.AdmissionReview{
		TypeMeta: metaV1.TypeMeta{
			Kind:       "AdmissionReview",
			APIVersion: "admission.k8s.io/v1",
		},
		Response: &admissionV1.AdmissionResponse{
			UID:     ReqInstance.UID,
			Allowed: true,
			Result: &metaV1.Status{
				Status:  "code=200",
				Message: "Allowed creating container",
			},
		},
	}

	//ResInstance.UID = ReqInstance.UID
	err = json.Unmarshal(ReqInstance.Object.Raw, &PodInstance)
	AR.Pod = PodInstance
	//fmt.Println(PodInstance)
	for _, ins := range PodInstance.Spec.Containers {
		_, cpu_ok := ins.Resources.Limits["cpu"]
		_, mem_ok := ins.Resources.Limits["memory"]
		imageSlice := strings.Split(ins.Image, ":")
		imageTag := imageSlice[len(imageSlice)-1]

		if ins.Resources.Limits == nil || !cpu_ok || !mem_ok || imageTag == "latest" {
			fmt.Println("Resources.Limits is NULL")
			ReponseAdmissionReview.Response.Allowed = false

			ReponseAdmissionReview.Response.Result.Status = "code=403"
			ReponseAdmissionReview.Response.Result.Message = "Forbidden, resources.limit not set( Must include cpu and memory, all containers)"
			if imageTag == "latest" {
				ReponseAdmissionReview.Response.Result.Message = "Forbidden, imageTag is latest(not allowed latest)"
			}
			break
		}
		fmt.Println("Resources.Limits is NOT NULL")

	}

	return ReponseAdmissionReview, err

}
