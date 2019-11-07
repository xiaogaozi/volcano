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
	vcbatch "volcano.sh/volcano/pkg/apis/batch/v1alpha1"
	"volcano.sh/volcano/pkg/controllers/apis"
)

type terminatingState struct {
	job *apis.JobInfo
}

func (ps *terminatingState) Execute(action vcbatch.Action) error {
	return KillJob(ps.job, PodRetainPhaseSoft, func(status *vcbatch.JobStatus) bool {
		// If any "alive" pods, still in Terminating phase
		if status.Terminating != 0 || status.Pending != 0 || status.Running != 0 {
			return false
		}
		status.State.Phase = vcbatch.Terminated
		return true

	})
}
