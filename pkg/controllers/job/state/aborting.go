/*
Copyright 2017 The Volcano Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package state

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	vcbatch "volcano.sh/volcano/pkg/apis/batch/v1alpha1"
	"volcano.sh/volcano/pkg/controllers/apis"
)

type abortingState struct {
	job *apis.JobInfo
}

func (ps *abortingState) Execute(action vcbatch.Action) error {
	switch action {
	case vcbatch.ResumeJobAction:
		return KillJob(ps.job, PodRetainPhaseSoft, func(status *vcbatch.JobStatus) bool {
			status.State.Phase = vcbatch.Restarting
			status.RetryCount++
			return true
		})
	default:
		return KillJob(ps.job, PodRetainPhaseSoft, func(status *vcbatch.JobStatus) bool {
			// If any "alive" pods, still in Aborting phase
			if status.Terminating != 0 || status.Pending != 0 || status.Running != 0 {
				return false
			}
			status.State.Phase = vcbatch.Aborted
			status.State.LastTransitionTime = metav1.Now()
			return true

		})
	}
}
