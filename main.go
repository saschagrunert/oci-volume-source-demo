package main

import (
	"log"

	. "github.com/saschagrunert/demo"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	r := NewRun(
		"KEP-4639: OCI VolumeSource",
		"",
		"https://github.com/kubernetes/enhancements/issues/4639",
	)

	r.Step(S(
		"Using OCI artifacts is a great way to store any kind of data",
		"in container registries. For example this artifact contains",
		"a single layer as tarball",
	), S(
		"oras manifest fetch quay.io/saschagrunert/artifact:v1 | jq .",
	))

	r.Step(S(
		"When pulling the artifact",
	), S(
		"oras pull quay.io/saschagrunert/artifact:v1",
	))

	r.Step(S(
		"And untar the layer",
	), S(
		"mkdir -p artifact &&",
		"tar xf layer.tar -C artifact",
	))

	r.Step(S(
		"Then we can see that it just contains a file as well as",
		"a directory containing another file",
	), S(
		"ls -T artifact &&",
		"cat artifact/dir/file artifact/file",
	))

	r.Step(S(
		"Let's consume this artifact as volume directly in Kubernetes!",
	), nil)

	r.Step(S(
		"Assuming we have a local Kubernetes cluster up and running which is",
		"based on PR: https://github.com/kubernetes/kubernetes/pull/125663",
		"and with the `OCIVolume` feature gate enabled",
	), S(
		"kubectl get nodes",
	))

	r.Step(S(
		"The cluster is running on top of a CRI-O work in progress PR, too:",
		"https://github.com/cri-o/cri-o/pull/8317",
	), S(
		"sudo crictl version",
	))

	r.Step(S(
		"The new OCI volume source can be used directly from the pod manifest",
	), S(
		"cat pod.yml",
	))

	r.Step(S(
		"If we apply the pod",
	), S(
		"kubectl apply -f pod.yml &&",
		"kubectl wait --for=condition=ready pod/pod",
	))

	r.Step(S(
		"Then we can see that the OCI object got pulled side by side to the",
		"container image",
	), S(
		"kubectl get pods,events",
	))

	r.Step(S(
		"And the artifact is available within the container",
	), S(
		"kubectl exec pod -- sh -c",
		"'set -x && ls -l /volume && cat /volume/dir/file && cat /volume/file'",
	))

	r.Step(S(
		"When removing the pod",
	), S(
		"kubectl delete pod pod --grace-period 1 &&",
		"kubectl delete events --all",
	))

	r.Step(S(
		"And applying it again",
	), S(
		"kubectl apply -f pod.yml &&",
		"kubectl wait --for=condition=ready pod/pod",
	))

	r.Step(S(
		"Then the runtime can re-use the existing mount",
	), S(
		"kubectl get pods,events",
	))

	r.StepCanFail(S(
		"If we now try to remove the artifact, then CRI-O",
		"will block that and mention the mounted reference",
	), S(
		"sudo crictl rmi quay.io/saschagrunert/artifact:v1",
	))

	r.Step(S(
		"A work in progress version of crictl is able to utilize",
		"the new CRI API, which we can see in the set `mountRef`",
	), S(
		"sudo ./crictl inspecti quay.io/saschagrunert/artifact:v1 | jq .",
	))

	r.Step(S(
		"crictl will be able to pull an OCI object and indicate that",
		"it should get mounted",
	), S(
		"sudo ./crictl pull --mount quay.io/saschagrunert/artifact:v1",
	))

	r.Step(nil, S(
		"sudo ls -la (sudo ./crictl inspecti quay.io/saschagrunert/artifact:v1 | jq -r .status.mountRef)",
	))

	r.Step(S(
		"This also works for any other container image and",
		"opens up many new use cases for Kubernetes",
	), S(
		"sudo ./crictl pull --mount alpine",
	))

	r.Step(nil, S(
		"sudo ls -la (sudo ./crictl inspecti alpine | jq -r .status.mountRef)",
	))

	return r.RunWithOptions(Options{Shell: "fish"})
}
