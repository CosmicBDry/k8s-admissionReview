package common

import (
	"github.com/gin-gonic/gin"
	admissionV1 "k8s.io/api/admission/v1"
	coreV1 "k8s.io/api/core/v1"
	res "k8s.io/apimachinery/pkg/api/resource"
)

type AdmissionReviewValidInterface interface {
	Validating(*gin.Context) (*admissionV1.AdmissionReview, error)
}

type AdmissionReviewMutateInterface interface {
	InjectContainer(coreV1.Pod) ([]map[string]interface{}, error)
	Mutating(*gin.Context) (*admissionV1.AdmissionReview, error)
}

func CreateSiderCarContainer(pod coreV1.Pod) *coreV1.Container {
	var (
		imageUrl  = "swr.cn-north-4.myhuaweicloud.com/ddn-k8s/docker.io/ikubernetes/filebeat:5.6.7-alpine"
		mountName = "log"
		envSlice  []coreV1.EnvVar
	)

	for _, ins := range pod.Spec.Containers {

		if len(ins.Env) > 0 {

			for _, v := range ins.Env {
				if v.Name == "FILEBEAT_IMAGE_URL" {
					imageUrl = v.Value
				} else if v.Name == "LOG_VolumeMountName" {
					mountName = v.Value
				} else if v.Name == "REDIS_HOST" || v.Name == "REDIS_PASS" || v.Name == "REDIS_TOPIC" {
					envSlice = append(envSlice, v)
				}

			}

		}

	}

	return &coreV1.Container{
		Name:  "filebeat",
		Image: imageUrl,
		Env:   envSlice,
		Resources: coreV1.ResourceRequirements{
			Limits: coreV1.ResourceList{
				"memory": *res.NewQuantity(300*1024*1024, res.BinarySI),
				"cpu":    *res.NewMilliQuantity(200, res.DecimalSI),
			},
		},
		VolumeMounts: []coreV1.VolumeMount{
			{
				Name:      mountName,
				MountPath: "/logs",
			},
		},
	}

}

// func CreateVolumes() []coreV1.Volume {

// 	FileMountType := coreV1.HostPathDirectoryOrCreate
// 	//var configMapDefaultMode int32 = 420
// 	return []coreV1.Volume{
// 		{
// 			Name: "log",
// 			VolumeSource: coreV1.VolumeSource{
// 				HostPath: &coreV1.HostPathVolumeSource{
// 					Path: "/data/log",
// 					Type: &FileMountType,
// 				},
// 			},
// 		},
// 	}
// }

// func InjectVolume(pod coreV1.Pod) []map[string]interface{} {

// 	var patchOps []map[string]interface{}
// 	requireVoluems := CreateVolumes()

// 	for _, vol := range requireVoluems {

// 		patchOp := map[string]interface{}{
// 			"op":    "add",
// 			"path":  "/spec/volumes/-",
// 			"value": vol,
// 		}
// 		if len(pod.Spec.Volumes) == 0 {

// 			patchOp["path"] = "/spec/volumes"
// 			patchOp["value"] = []coreV1.Volume{vol}

// 		}else{

// 			for _,v :=range pod.Spec.Volumes{
// 				if  v.HostPath.Path == vol.HostPath.Path {
// 					return nil
// 				}
// 			}
// 		}
// 		patchOps = append(patchOps, patchOp)

// 	}
// 	return patchOps

// }
