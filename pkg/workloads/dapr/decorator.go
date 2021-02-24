// ------------------------------------------------------------
// Copyright (c) Microsoft Corporation.
// Licensed under the MIT License.
// ------------------------------------------------------------

package dapr

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/Azure/radius/pkg/workloads"
	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

// Renderer is the WorkloadRenderer implementation for the dapr deployment decorator.
type Renderer struct {
	Inner workloads.WorkloadRenderer
}

// Allocate is the WorkloadRenderer implementation for the dapr deployment decorator.
func (r Renderer) Allocate(ctx context.Context, w workloads.InstantiatedWorkload, wrp []workloads.WorkloadResourceProperties, service workloads.WorkloadService) (map[string]interface{}, error) {
	return r.Inner.Allocate(ctx, w, wrp, service)
}

// Render is the WorkloadRenderer implementation for the dapr deployment decorator.
func (r Renderer) Render(ctx context.Context, w workloads.InstantiatedWorkload) ([]workloads.WorkloadResource, error) {
	// Let the inner renderer do its work
	resources, err := r.Inner.Render(ctx, w)
	if err != nil {
		return []workloads.WorkloadResource{}, err
	}

	trait, err := r.findDaprTrait(w.Traits)
	if err != nil {
		return []workloads.WorkloadResource{}, err
	}

	if trait == nil {
		// no dapr
		return resources, err
	}

	// dapr detected! update the deployment
	for _, res := range resources {
		if res.Type != "kubernetes" {
			// Not a kubernetes resource
			continue
		}

		o, ok := res.Resource.(runtime.Object)
		if !ok {
			return []workloads.WorkloadResource{}, errors.New("Found kubernetes resource with non-Kubernetes paylod")
		}

		annotations, ok := r.getAnnotations(o)
		if !ok {
			continue
		}

		// use the workload name
		if trait.Properties.AppID == "" {
			trait.Properties.AppID = w.Workload.GetName()
		}

		annotations["dapr.io/enabled"] = "true"
		annotations["dapr.io/app-id"] = trait.Properties.AppID
		if trait.Properties.AppPort != 0 {
			annotations["dapr.io/app-port"] = fmt.Sprintf("%d", trait.Properties.AppPort)
		}
		if trait.Properties.Config != "" {
			annotations["dapr.io/config"] = trait.Properties.Config
		}
		if trait.Properties.Protocol != "" {
			annotations["dapr.io/protocol"] = trait.Properties.Protocol
		}

		r.setAnnotations(o, annotations)

		// HACK: for Apps types, set the minimum replicas to 1.
		// The autoscaler implementation is not aware of Dapr traffic or bindings/pubsub.

		gvk := o.GetObjectKind().GroupVersionKind()
		appkind := schema.GroupVersionKind{Group: "k8se.microsoft.com", Version: "v1alpha1", Kind: "App"}
		if gvk == appkind {
			if un, ok := o.(*unstructured.Unstructured); ok {
				if obj, ok := un.Object["spec"]; ok {
					if spec, ok := obj.(map[string]interface{}); ok {

						var scaleOptions map[string]interface{}
						if obj, ok := spec["scaleOptions"]; ok {
							scaleOptions = obj.(map[string]interface{})
						}

						if scaleOptions == nil {
							scaleOptions = map[string]interface{}{}
							spec["scaleOptions"] = scaleOptions
						}

						scaleOptions["minReplicaCount"] = 1
					}
				}
			}
		}
	}

	return resources, err
}

type daprTrait struct {
	Kind       string              `json:"kind"`
	Properties daprTraitProperties `json:"properties"`
}

type daprTraitProperties struct {
	AppID    string `json:"appId"`
	AppPort  int    `json:"appPort"`
	Config   string `json:"config"`
	Protocol string `json:"protocol"`
}

func (r Renderer) findDaprTrait(traits []workloads.WorkloadTrait) (*daprTrait, error) {
	var match *workloads.WorkloadTrait
	for _, t := range traits {
		if t.Kind == "dapr.io/App@v1alpha1" {
			match = &t
			break
		}
	}

	if match == nil {
		return nil, nil
	}

	val := map[string]interface{}{
		"kind":       match.Kind,
		"properties": match.Properties,
	}
	b, err := json.Marshal(val)
	if err != nil {
		return nil, fmt.Errorf("error reading trait '%v': %w", match.Kind, err)
	}

	trait := &daprTrait{}
	err = json.Unmarshal(b, trait)
	if err != nil {
		return nil, fmt.Errorf("error reading trait '%v': %w", match.Kind, err)
	}

	return trait, nil
}

func (r Renderer) getAnnotations(o runtime.Object) (map[string]string, bool) {
	dep, ok := o.(*appsv1.Deployment)
	if ok {
		if dep.Spec.Template.Annotations == nil {
			dep.Spec.Template.Annotations = map[string]string{}
		}

		return dep.Spec.Template.Annotations, true
	}

	un, ok := o.(*unstructured.Unstructured)
	if ok {
		if a := un.GetAnnotations(); a != nil {
			return a, true
		}

		return map[string]string{}, true
	}

	return nil, false
}

func (r Renderer) setAnnotations(o runtime.Object, annotations map[string]string) {
	un, ok := o.(*unstructured.Unstructured)
	if ok {
		un.SetAnnotations(annotations)
	}
}
