package main

import (
	"github.com/ckatsak/actischedulertest01/acti"

	"k8s.io/klog/v2"
	sched "k8s.io/kubernetes/cmd/kube-scheduler/app"
)

func main() {
	cmd := sched.NewSchedulerCommand(
		sched.WithPlugin(acti.Name, acti.New),
	)
	if err := cmd.Execute(); err != nil {
		klog.Fatalf("failed to execute %q: %v", acti.Name, err)
	}
}
