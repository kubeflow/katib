package main

import (
	"fmt"
	"strconv"
	"strings"

	batchv1 "k8s.io/api/batch/v1"
	v1 "k8s.io/api/core/v1"
	resource "k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Generate Job Template
func genJobTemplate(lt *learningTask) *batchv1.Job {
	//construct entry point nad parameter
	cmd := lt.ltc.EntryPoint
	args := lt.ltc.Parameters

	template := &batchv1.Job{
		TypeMeta: metav1.TypeMeta{
			Kind: "Job",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: "", // must be filled later
			Labels: map[string]string{
				"priority": strconv.Itoa(lt.ltc.Priority),
			},
		},
		Spec: batchv1.JobSpec{
			Template: v1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{}, // "type", "app", "learning-task", "nrPSes", "nrWorkers" must be filled
				},

				Spec: v1.PodSpec{
					SchedulerName: lt.ltc.Scheduler,

					Containers: []v1.Container{
						{
							Name:            "",
							Command:         strings.Fields(cmd),
							Args:            strings.Fields(args),
							ImagePullPolicy: v1.PullAlways,
							Ports: []v1.ContainerPort{
								v1.ContainerPort{
									ContainerPort: 2222,
								},
							},
						},
					},
					RestartPolicy: v1.RestartPolicyOnFailure,
				},
			},
		},
	}

	// Specified pvc is mounted to both PS and Worker Pods
	if lt.pvc != nil {
		if lt.ltc.Pvc != "" {
			template.Spec.Template.Spec.Volumes = []v1.Volume{
				v1.Volume{
					Name: "pvc-mount-point",
					VolumeSource: v1.VolumeSource{
						PersistentVolumeClaim: &v1.PersistentVolumeClaimVolumeSource{
							ClaimName: lt.ltc.Pvc,
						},
					},
				},
			}

			template.Spec.Template.Spec.Containers[0].VolumeMounts = []v1.VolumeMount{
				v1.VolumeMount{
					Name:      "pvc-mount-point",
					MountPath: lt.ltc.MountPath,
				},
			}
		}
	}

	return template
}

// PS Job
type psJob struct {
	name string
	job  *batchv1.Job
}

// Worker Job
type workerJob struct {
	name string
	job  *batchv1.Job
}

// Create PS Jobs
func newPSJobs(lt *learningTask, memberLists []string) []*psJob {
	ret := make([]*psJob, lt.ltc.NrPS)
	for i := 0; i < lt.ltc.NrPS; i++ {
		template := genJobTemplate(lt)
		psname := fmt.Sprintf("%s-ps-%d", lt.name, i) // FIXME: should be shared with services
		template.ObjectMeta.Name = psname
		template.Spec.Template.ObjectMeta.Labels["app"] = psname
		template.Spec.Template.ObjectMeta.Labels["learning-task"] = lt.name
		template.Spec.Template.ObjectMeta.Labels["nrPSes"] = fmt.Sprintf("%d", lt.ltc.NrPS)
		template.Spec.Template.ObjectMeta.Labels["nrWorkers"] = fmt.Sprintf("%d", lt.ltc.NrWorker)
		template.Spec.Template.ObjectMeta.Labels["type"] = "PS"
		for _, e := range lt.ltc.Envs {
			template.Spec.Template.Spec.Containers[0].Env = append(template.Spec.Template.Spec.Containers[0].Env, v1.EnvVar{Name: e.Name, Value: e.Value})
		}
		template.Spec.Template.Spec.ImagePullSecrets = []v1.LocalObjectReference{
			v1.LocalObjectReference{
				Name: lt.ltc.PullSecret,
			},
		}

		template.Spec.Template.Spec.Containers[0].Name = psname
		template.Spec.Template.Spec.Containers[0].Args =
			append(template.Spec.Template.Spec.Containers[0].Args, "--job_name=ps")
		template.Spec.Template.Spec.Containers[0].Args =
			append(template.Spec.Template.Spec.Containers[0].Args, fmt.Sprintf("--task_index=%d", i))
		template.Spec.Template.Spec.Containers[0].Args =
			append(template.Spec.Template.Spec.Containers[0].Args, memberLists...)
		template.Spec.Template.Spec.Containers[0].Image = lt.ltc.PsImage

		if template.Spec.Template.Spec.Volumes == nil {
			template.Spec.Template.Spec.Volumes = []v1.Volume{}
		}
		hostPathDirectoryOrCreate := v1.HostPathDirectoryOrCreate
		vol := v1.Volume{
			Name: "nvidialib",
			VolumeSource: v1.VolumeSource{
				HostPath: &v1.HostPathVolumeSource{
					Path: "/usr/lib/nvidia",
					Type: &hostPathDirectoryOrCreate,
				},
			},
		}
		template.Spec.Template.Spec.Volumes = append(template.Spec.Template.Spec.Volumes, vol)

		if template.Spec.Template.Spec.Volumes == nil {
			template.Spec.Template.Spec.Containers[0].VolumeMounts = []v1.VolumeMount{}
		}
		volMnt := v1.VolumeMount{
			Name:      "nvidialib",
			MountPath: "/usr/local/nvidia",
			ReadOnly:  true,
		}
		template.Spec.Template.Spec.Containers[0].VolumeMounts = append(template.Spec.Template.Spec.Containers[0].VolumeMounts, volMnt)

		ret[i] = &psJob{
			name: psname,
			job:  template,
		}
	}

	return ret
}

// Create Worker Jobs
func newWorkerJobs(lt *learningTask, memberLists []string) []*workerJob {
	ret := make([]*workerJob, lt.ltc.NrWorker)
	for i := 0; i < lt.ltc.NrWorker; i++ {
		template := genJobTemplate(lt)
		workername := fmt.Sprintf("%s-worker-%d", lt.name, i) // FIXME: should be shared with services
		template.ObjectMeta.Name = workername
		template.Spec.Template.ObjectMeta.Labels["type"] = "worker"
		template.Spec.Template.ObjectMeta.Labels["app"] = workername
		template.Spec.Template.ObjectMeta.Labels["learning-task"] = lt.name
		template.Spec.Template.ObjectMeta.Labels["nrPSes"] = fmt.Sprintf("%d", lt.ltc.NrPS)
		template.Spec.Template.ObjectMeta.Labels["nrWorkers"] = fmt.Sprintf("%d", lt.ltc.NrWorker)
		template.Spec.Template.Spec.Containers[0].Name = workername
		template.Spec.Template.Spec.Containers[0].Image = lt.ltc.WorkerImage
		for _, e := range lt.ltc.Envs {
			template.Spec.Template.Spec.Containers[0].Env = append(template.Spec.Template.Spec.Containers[0].Env, v1.EnvVar{Name: e.Name, Value: e.Value})
		}
		template.Spec.Template.Spec.ImagePullSecrets = []v1.LocalObjectReference{
			v1.LocalObjectReference{
				Name: lt.ltc.PullSecret,
			},
		}

		if lt.ltc.NrPS > 0 {
			template.Spec.Template.Spec.Containers[0].Args =
				append(template.Spec.Template.Spec.Containers[0].Args, "--job_name=worker")
			template.Spec.Template.Spec.Containers[0].Args =
				append(template.Spec.Template.Spec.Containers[0].Args, fmt.Sprintf("--task_index=%d", i))
			template.Spec.Template.Spec.Containers[0].Args =
				append(template.Spec.Template.Spec.Containers[0].Args, memberLists...)
		}
		if lt.ltc.Gpu > 0 {
			gpuReq, err := resource.ParseQuantity(strconv.Itoa(lt.ltc.Gpu))
			if err != nil {
				return nil
			}
			template.Spec.Template.Spec.Containers[0].Resources =
				v1.ResourceRequirements{
					Limits: v1.ResourceList{"nvidia.com/gpu": gpuReq},
					//					Limits:   v1.ResourceList{"alpha.kubernetes.io/nvidia-gpu": gpuReq},
					//					Requests: v1.ResourceList{"alpha.kubernetes.io/nvidia-gpu": gpuReq},
				}

			//			if template.Spec.Template.Spec.Volumes == nil {
			//				template.Spec.Template.Spec.Volumes = []v1.Volume{}
			//			}
			//			hostPathDirectoryOrCreate := v1.HostPathDirectoryOrCreate
			//			vol := v1.Volume{
			//				Name: "nvidialib",
			//				VolumeSource: v1.VolumeSource{
			//					HostPath: &v1.HostPathVolumeSource{
			//						Path: "/usr/lib/nvidia",
			//						Type: &hostPathDirectoryOrCreate,
			//					},
			//				},
			//			}
			//			template.Spec.Template.Spec.Volumes = append(template.Spec.Template.Spec.Volumes, vol)
			//
			//			if template.Spec.Template.Spec.Volumes == nil {
			//				template.Spec.Template.Spec.Containers[0].VolumeMounts = []v1.VolumeMount{}
			//			}
			//			volMnt := v1.VolumeMount{
			//				Name:      "nvidialib",
			//				MountPath: "/usr/local/nvidia",
			//				ReadOnly:  true,
			//			}
			//			template.Spec.Template.Spec.Containers[0].VolumeMounts = append(template.Spec.Template.Spec.Containers[0].VolumeMounts, volMnt)
		}

		ret[i] = &workerJob{
			name: workername,
			job:  template,
		}
	}

	return ret
}
