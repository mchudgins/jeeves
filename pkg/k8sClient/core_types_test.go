package k8sClient

import (
	"log"
	"testing"

	"github.com/technosophos/kubelite/v1"
)

func TestCreatePods(t *testing.T) {
	pod := &v1.Pod{
		ObjectMeta: v1.ObjectMeta{
			Name: "golang-builder",
			/*Labels: map[string]string{
				{"fubar": "gorf"},
			},*/
		},
		Spec: v1.PodSpec{
			Containers: []v1.Container{
				{
					Name:  "golang-builder",
					Image: "registry.dstresearch.com/jump-box:latest",
					Env: []v1.EnvVar{
						{
							Name:  "GIT_REPO",
							Value: "https://github.com/mchudgins/testRepo.git",
						},
						{
							Name:  "Ref",
							Value: "ref/heads/master",
						},
					},
				},
			},
		},
	}

	k8s := NewClientOrDie()
	resp, err := k8s.LaunchPod(k8s.CurrentNamespace, pod)
	if err != nil {
		t.Fatal(err)
	}
	log.Printf("response: %v", resp)
}

func TestPodList(t *testing.T) {
	k8s := NewClientOrDie()
	pods, err := k8s.PodList(k8s.CurrentNamespace)
	if err != nil {
		t.Fatal(err)
	}

	for _, pod := range pods.Items {
		log.Printf("pod: %s %s", pod.Name, pod.Status.Phase)
	}

}
