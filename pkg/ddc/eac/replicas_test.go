/*
  Copyright 2022 The Fluid Authors.

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

package eac

import (
	"testing"

	datav1alpha1 "github.com/fluid-cloudnative/fluid/api/v1alpha1"
	"github.com/fluid-cloudnative/fluid/pkg/common"
	"github.com/fluid-cloudnative/fluid/pkg/ctrl"
	v1 "k8s.io/api/core/v1"
	utilpointer "k8s.io/utils/pointer"

	"github.com/fluid-cloudnative/fluid/pkg/utils/fake"
	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/record"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/fluid-cloudnative/fluid/pkg/ddc/base"
	cruntime "github.com/fluid-cloudnative/fluid/pkg/runtime"
)

func newEACEngineREP(client client.Client, name string, namespace string) *EACEngine {
	runTimeInfo, _ := base.BuildRuntimeInfo(name, namespace, common.EACRuntimeType, datav1alpha1.TieredStore{})
	engine := &EACEngine{
		runtime:     &datav1alpha1.EACRuntime{},
		name:        name,
		namespace:   namespace,
		Client:      client,
		runtimeInfo: runTimeInfo,
		Log:         fake.NullLogger(),
	}
	engine.Helper = ctrl.BuildHelper(runTimeInfo, client, engine.Log)
	return engine
}

func TestSyncReplicas(t *testing.T) {
	nodeInputs := []*v1.Node{
		{
			ObjectMeta: metav1.ObjectMeta{
				Name: "test-node-spark",
				Labels: map[string]string{
					"fluid.io/dataset-num":           "1",
					"fluid.io/s-eac-fluid-spark":     "true",
					"fluid.io/s-fluid-spark":         "true",
					"fluid.io/s-h-eac-d-fluid-spark": "5B",
					"fluid.io/s-h-eac-m-fluid-spark": "1B",
					"fluid.io/s-h-eac-t-fluid-spark": "6B",
					"fluid_exclusive":                "fluid_spark",
				},
			},
		},
		{
			ObjectMeta: metav1.ObjectMeta{
				Name: "test-node-share",
				Labels: map[string]string{
					"fluid.io/dataset-num":            "2",
					"fluid.io/s-eac-fluid-hadoop":     "true",
					"fluid.io/s-fluid-hadoop":         "true",
					"fluid.io/s-h-eac-d-fluid-hadoop": "5B",
					"fluid.io/s-h-eac-m-fluid-hadoop": "1B",
					"fluid.io/s-h-eac-t-fluid-hadoop": "6B",
					"fluid.io/s-eac-fluid-hbase":      "true",
					"fluid.io/s-fluid-hbase":          "true",
					"fluid.io/s-h-eac-d-fluid-hbase":  "5B",
					"fluid.io/s-h-eac-m-fluid-hbase":  "1B",
					"fluid.io/s-h-eac-t-fluid-hbase":  "6B",
				},
			},
		},
		{
			ObjectMeta: metav1.ObjectMeta{
				Name: "test-node-hadoop",
				Labels: map[string]string{
					"fluid.io/dataset-num":            "1",
					"fluid.io/s-eac-fluid-hadoop":     "true",
					"fluid.io/s-fluid-hadoop":         "true",
					"fluid.io/s-h-eac-d-fluid-hadoop": "5B",
					"fluid.io/s-h-eac-m-fluid-hadoop": "1B",
					"fluid.io/s-h-eac-t-fluid-hadoop": "6B",
					"node-select":                     "true",
				},
			},
		},
	}
	runtimeInputs := []*datav1alpha1.EACRuntime{
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "hbase",
				Namespace: "fluid",
			},
			Spec: datav1alpha1.EACRuntimeSpec{
				Replicas: 3, // 2
			},
			Status: datav1alpha1.RuntimeStatus{
				CurrentWorkerNumberScheduled: 2,
				DesiredWorkerNumberScheduled: 2,
				Conditions:                   []datav1alpha1.RuntimeCondition{},
			},
		},
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "hadoop",
				Namespace: "fluid",
			},
			Spec: datav1alpha1.EACRuntimeSpec{
				Replicas: 2,
			},
			Status: datav1alpha1.RuntimeStatus{
				CurrentWorkerNumberScheduled: 3,
				DesiredWorkerNumberScheduled: 3,
				Conditions:                   []datav1alpha1.RuntimeCondition{},
			},
		},
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "obj",
				Namespace: "fluid",
			},
			Spec: datav1alpha1.EACRuntimeSpec{
				Replicas: 2,
			},
			Status: datav1alpha1.RuntimeStatus{
				CurrentWorkerNumberScheduled: 2,
				DesiredWorkerNumberScheduled: 2,
				Conditions:                   []datav1alpha1.RuntimeCondition{},
			},
		},
	}
	workersInputs := []*appsv1.StatefulSet{
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "hbase-worker",
				Namespace: "fluid",
			},
			Spec: appsv1.StatefulSetSpec{
				Replicas: utilpointer.Int32Ptr(2),
			},
		},
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "hadoop-worker",
				Namespace: "fluid",
			},
			Spec: appsv1.StatefulSetSpec{
				Replicas: utilpointer.Int32Ptr(3),
			},
		},
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "obj-worker",
				Namespace: "fluid",
			},
			Spec: appsv1.StatefulSetSpec{
				Replicas: utilpointer.Int32Ptr(2),
			},
		},
	}
	dataSetInputs := []*datav1alpha1.Dataset{
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "hbase",
				Namespace: "fluid",
			},
		},
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "hadoop",
				Namespace: "fluid",
			},
		},
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "obj",
				Namespace: "fluid",
			},
		},
	}
	fuseInputs := []*appsv1.DaemonSet{
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "hbase-fuse",
				Namespace: "fluid",
			},
		},
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "hadoop-fuse",
				Namespace: "fluid",
			},
		},
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "obj-fuse",
				Namespace: "fluid",
			},
		},
	}

	objs := []runtime.Object{}
	for _, nodeInput := range nodeInputs {
		objs = append(objs, nodeInput.DeepCopy())
	}
	for _, runtimeInput := range runtimeInputs {
		objs = append(objs, runtimeInput.DeepCopy())
	}
	for _, workerInput := range workersInputs {
		objs = append(objs, workerInput.DeepCopy())
	}
	for _, fuseInput := range fuseInputs {
		objs = append(objs, fuseInput.DeepCopy())
	}
	for _, dataSetInput := range dataSetInputs {
		objs = append(objs, dataSetInput.DeepCopy())
	}

	fakeClient := fake.NewFakeClientWithScheme(testScheme, objs...)
	testCases := []struct {
		testName       string
		name           string
		namespace      string
		Type           datav1alpha1.RuntimeConditionType
		isErr          bool
		condtionLength int
	}{
		{
			testName:       "scaleout",
			name:           "hbase",
			namespace:      "fluid",
			Type:           datav1alpha1.RuntimeWorkerScaledOut,
			isErr:          false,
			condtionLength: 1,
		},
		{
			testName:       "scalein",
			name:           "hadoop",
			namespace:      "fluid",
			Type:           datav1alpha1.RuntimeWorkerScaledIn,
			isErr:          false,
			condtionLength: 1,
		},
		{
			testName:       "noscale",
			name:           "obj",
			namespace:      "fluid",
			Type:           "",
			isErr:          false,
			condtionLength: 0,
		},
	}
	for _, testCase := range testCases {
		engine := newEACEngineREP(fakeClient, testCase.name, testCase.namespace)
		err := engine.SyncReplicas(cruntime.ReconcileRequestContext{
			Log:      fake.NullLogger(),
			Recorder: record.NewFakeRecorder(300),
		})
		if err != nil {
			t.Errorf("sync replicas failed,err:%s", err.Error())
		}

		if testCase.condtionLength == 0 {
			return
		}

		rt, _ := engine.getRuntime()
		found := false

		for _, cond := range rt.Status.Conditions {
			if cond.Type == testCase.Type {
				found = true
				break
			}
		}

		if !found {
			t.Errorf("testCase: %s runtime condition want conditionType %v, got  conditions %v", testCase.testName, testCase.Type, rt.Status.Conditions)
		}
	}
}
