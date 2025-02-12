/*
Copyright 2023 The Radius Authors.

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

package resource_test

import (
	"testing"

	"github.com/radius-project/radius/test/functional/shared"
	"github.com/radius-project/radius/test/step"
	"github.com/radius-project/radius/test/testutil"
	"github.com/radius-project/radius/test/validation"
)

func Test_PersistentVolume(t *testing.T) {
	template := "testdata/corerp-resources-volume-azure-keyvault.bicep"
	name := "corerp-resources-volume-azure-keyvault"
	appNamespace := "corerp-resources-volume-azure-keyvault-app"

	test := shared.NewRPTest(t, name, []shared.TestStep{
		{
			Executor: step.NewDeployExecutor(template, testutil.GetMagpieImage(), testutil.GetOIDCIssuer()),
			RPResources: &validation.RPResourceSet{
				Resources: []validation.RPResource{
					{
						Name: "corerp-azure-workload-env",
						Type: validation.EnvironmentsResource,
					},
					{
						Name: name,
						Type: validation.ApplicationsResource,
					},
					{
						Name: "volume-azkv-ctnr",
						Type: validation.ContainersResource,
						App:  name,
					},
					{
						Name: "volume-azkv",
						Type: validation.VolumesResource,
						App:  name,
					},
				},
			},
			K8sObjects: &validation.K8sObjectSet{
				Namespaces: map[string][]validation.K8sObject{
					appNamespace: {
						validation.NewK8sPodForResource(name, "volume-azkv-ctnr"),
					},
				},
			},
		},
	})
	test.RequiredFeatures = []shared.RequiredFeature{shared.FeatureCSIDriver}

	test.Test(t)
}
