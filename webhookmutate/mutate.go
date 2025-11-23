package webhookmutate

import (
	"encoding/json"
	"io"

	"github.com/CosmicBDry/k8s-admissionReview/common"
	"github.com/gin-gonic/gin"
	admissionV1 "k8s.io/api/admission/v1"
	coreV1 "k8s.io/api/core/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type AdmissionReviewMutate struct {
	Pod                coreV1.Pod
	K8sAdmissionReview admissionV1.AdmissionReview
}

func NewAdmissionReviewMutate() *AdmissionReviewMutate {
	return &AdmissionReviewMutate{}
}

func (AR *AdmissionReviewMutate) InjectContainer(pod coreV1.Pod) ([]map[string]interface{}, error) {

	var patchOps []map[string]interface{}

	siderCarcontainer := common.CreateSiderCarContainer(pod)

	patchOp := map[string]interface{}{
		"op":    "add",
		"path":  "/spec/containers/-",
		"value": siderCarcontainer,
	}

	for _, ins := range pod.Spec.Containers {
		if ins.Name == siderCarcontainer.Name {
			patchOp = nil
			break
		}

	}

	if len(patchOp) > 0 {
		patchOps = append(patchOps, patchOp)
	}

	// patchOpVolumes := common.InjectVolume(pod)
	// if len(patchOpVolumes) > 0 {
	// 	patchOps = append(patchOps, patchOpVolumes...)
	// }

	return patchOps, nil
}

func (AR *AdmissionReviewMutate) Mutating(c *gin.Context) (*admissionV1.AdmissionReview, error) {

	reqBytes, err := io.ReadAll(c.Request.Body)
	err = json.Unmarshal(reqBytes, &AR.K8sAdmissionReview)
	ReqInstance := AR.K8sAdmissionReview.Request
	podInstance := AR.Pod
	err = json.Unmarshal(ReqInstance.Object.Raw, &podInstance)

	patchOps, err := AR.InjectContainer(podInstance)

	patchOpsBytes, err := json.Marshal(patchOps)

	patch_type := admissionV1.PatchTypeJSONPatch

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
				Message: "Allowe created",
			},
		},
	}

	if len(patchOpsBytes) > 0 {

		if val, ok := podInstance.Annotations["myk8s.io/webhookmutate-plugin"]; ok && val == "enable" {

			ReponseAdmissionReview.Response.Result.Status = "code=200"
			ReponseAdmissionReview.Response.Result.Message = "allowed inject_sider_container"
			ReponseAdmissionReview.Response.PatchType = &patch_type
			ReponseAdmissionReview.Response.Patch = patchOpsBytes

		}

	}

	return ReponseAdmissionReview, err
}
