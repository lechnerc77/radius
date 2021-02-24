// ------------------------------------------------------------
// Copyright (c) Microsoft Corporation.
// Licensed under the MIT License.
// ------------------------------------------------------------

package curp

import (
	"encoding/json"
	"testing"

	"github.com/Azure/radius/pkg/curp/db"
	"github.com/Azure/radius/pkg/curp/metadata"
	"github.com/Azure/radius/pkg/curp/revision"
	"github.com/Azure/radius/pkg/workloads/containerv1alpha1"
	"github.com/stretchr/testify/require"
)

func Test_DeploymentCreated_NoComponents(t *testing.T) {
	app := db.NewApplication()
	newer := db.NewDeployment()
	rp := rp{
		meta: metadata.NewRegistry(),
	}

	actions, err := rp.computeDeploymentActions(app, nil, newer)
	require.NoError(t, err)

	require.True(t, deploymentIsNoOp(actions))
	require.Empty(t, actions)
}

func Test_DeploymentCreated_ValidationError(t *testing.T) {
	app := db.NewApplication()
	newer := db.NewDeployment()
	newer.Properties.Components = []*db.DeploymentComponent{
		{
			ComponentName: "", // intentionally empty
		},
	}

	rp := rp{
		meta: metadata.NewRegistry(),
	}

	_, err := rp.computeDeploymentActions(app, nil, newer)
	require.Error(t, err)
}

func Test_DeploymentCreated_ErrMissingComponent(t *testing.T) {
	app := db.NewApplication()
	newer := db.NewDeployment()
	newer.Properties.Components = []*db.DeploymentComponent{
		{
			ComponentName: "A",
		},
	}

	rp := rp{
		meta: metadata.NewRegistry(),
	}

	_, err := rp.computeDeploymentActions(app, nil, newer)
	require.Error(t, err)
}

func Test_DeploymentCreated_ErrNoRevisions(t *testing.T) {
	app := db.NewApplication()
	app.Components["A"] = db.ComponentHistory{
		Revision:        revision.Revision(""),
		RevisionHistory: []db.ComponentRevision{},
	}
	newer := db.NewDeployment()
	newer.Properties.Components = []*db.DeploymentComponent{
		{
			ComponentName: "A",
			Revision:      "1",
		},
	}

	rp := rp{
		meta: metadata.NewRegistry(),
	}

	_, err := rp.computeDeploymentActions(app, nil, newer)
	require.Error(t, err)
}

func Test_DeploymentCreated_ErrMissingComponentRevision(t *testing.T) {
	app := db.NewApplication()
	app.Components["A"] = db.ComponentHistory{
		Revision: revision.Revision("1"),
		RevisionHistory: []db.ComponentRevision{
			{
				Kind:       "radius.dev/Container@v1alpha1",
				Revision:   revision.Revision("1"),
				Properties: *db.NewComponentProperties(),
			},
		},
	}
	newer := db.NewDeployment()
	newer.Properties.Components = []*db.DeploymentComponent{
		{
			ComponentName: "A",
			Revision:      "2",
		},
	}

	rp := rp{
		meta: metadata.NewRegistry(),
	}

	_, err := rp.computeDeploymentActions(app, nil, newer)
	require.Error(t, err)
}

func Test_DeploymentCreated_OneComponent_NoRevisionSpecified(t *testing.T) {
	app := db.NewApplication()
	app.Components["A"] = db.ComponentHistory{
		Revision: revision.Revision("2"),
		RevisionHistory: []db.ComponentRevision{
			{
				Kind:       "radius.dev/Container@v1alpha1",
				Revision:   revision.Revision("1"),
				Properties: *db.NewComponentProperties(),
			},
			{
				Kind:       "radius.dev/Container@v1alpha1",
				Revision:   revision.Revision("2"),
				Properties: *db.NewComponentProperties(),
			},
		},
	}

	newer := db.NewDeployment()
	newer.Properties.Components = []*db.DeploymentComponent{
		{
			ComponentName: "A",
			Revision:      revision.Revision(""), // Intentionally blank
		},
	}

	rp := rp{
		meta: metadata.NewRegistry(),
	}

	actions, err := rp.computeDeploymentActions(app, nil, newer)
	require.NoError(t, err)

	require.False(t, deploymentIsNoOp(actions))

	// Revision is updated in the deployment.
	require.Equal(t, revision.Revision("2"), newer.Properties.Components[0].Revision)

	// Updates to the components are in the actions
	require.Len(t, actions, 1)

	require.Contains(t, actions, "A")
	action := actions["A"]
	require.Equal(t, CreateWorkload, action.Operation)
	require.Equal(t, "A", action.ComponentName)
	require.Equal(t, app.Components["A"].RevisionHistory[1], *action.Definition)
	require.Equal(t, newer.Properties.Components[0], action.Instantiation)
	require.Nil(t, action.PreviousDefinition)
	require.Nil(t, action.PreviousInstanitation)
}

func Test_DeploymentCreated_OneComponent_RevisionSpecified(t *testing.T) {
	app := db.NewApplication()
	app.Components["A"] = db.ComponentHistory{
		Revision: revision.Revision("2"),
		RevisionHistory: []db.ComponentRevision{
			{
				Kind:       "radius.dev/Container@v1alpha1",
				Revision:   revision.Revision("1"),
				Properties: *db.NewComponentProperties(),
			},
			{
				Kind:       "radius.dev/Container@v1alpha1",
				Revision:   revision.Revision("2"),
				Properties: *db.NewComponentProperties(),
			},
		},
	}

	newer := db.NewDeployment()
	newer.Properties.Components = []*db.DeploymentComponent{
		{
			ComponentName: "A",
			Revision:      revision.Revision("1"),
		},
	}

	rp := rp{
		meta: metadata.NewRegistry(),
	}

	actions, err := rp.computeDeploymentActions(app, nil, newer)
	require.NoError(t, err)

	require.False(t, deploymentIsNoOp(actions))

	// Revision is updated in the deployment.
	require.Equal(t, revision.Revision("1"), newer.Properties.Components[0].Revision)

	// Updates to the components are in the actions
	require.Len(t, actions, 1)

	require.Contains(t, actions, "A")
	action := actions["A"]
	require.Equal(t, CreateWorkload, action.Operation)
	require.Equal(t, "A", action.ComponentName)
	require.Equal(t, app.Components["A"].RevisionHistory[0], *action.Definition)
	require.Equal(t, newer.Properties.Components[0], action.Instantiation)
	require.Nil(t, action.PreviousDefinition)
	require.Nil(t, action.PreviousInstanitation)
}

func Test_DeploymentCreated_MultipleComponents(t *testing.T) {
	app := db.NewApplication()
	app.Components["A"] = db.ComponentHistory{
		Revision: revision.Revision("2"),
		RevisionHistory: []db.ComponentRevision{
			{
				Kind:       "radius.dev/Container@v1alpha1",
				Revision:   revision.Revision("1"),
				Properties: *db.NewComponentProperties(),
			},
			{
				Kind:       "radius.dev/Container@v1alpha1",
				Revision:   revision.Revision("2"),
				Properties: *db.NewComponentProperties(),
			},
		},
	}
	app.Components["B"] = db.ComponentHistory{
		Revision: revision.Revision("2"),
		RevisionHistory: []db.ComponentRevision{
			{
				Kind:       "radius.dev/Container@v1alpha1",
				Revision:   revision.Revision("1"),
				Properties: *db.NewComponentProperties(),
			},
			{
				Kind:       "radius.dev/Container@v1alpha1",
				Revision:   revision.Revision("2"),
				Properties: *db.NewComponentProperties(),
			},
		},
	}

	app.Components["C"] = db.ComponentHistory{
		Revision: revision.Revision("1"),
		RevisionHistory: []db.ComponentRevision{
			{
				Kind:       "radius.dev/Container@v1alpha1",
				Revision:   revision.Revision("1"),
				Properties: *db.NewComponentProperties(),
			},
		},
	}

	newer := db.NewDeployment()
	newer.Properties.Components = []*db.DeploymentComponent{
		{
			ComponentName: "A",
			Revision:      revision.Revision(""), // Intentionally empty
		},
		{
			ComponentName: "B",
			Revision:      revision.Revision("1"),
		},
		{
			ComponentName: "C",
			Revision:      revision.Revision("1"),
		},
	}

	rp := rp{
		meta: metadata.NewRegistry(),
	}

	actions, err := rp.computeDeploymentActions(app, nil, newer)
	require.NoError(t, err)

	require.False(t, deploymentIsNoOp(actions))

	// Revision is updated in the deployment.
	require.Equal(t, revision.Revision("2"), newer.Properties.Components[0].Revision)
	require.Equal(t, revision.Revision("1"), newer.Properties.Components[1].Revision)
	require.Equal(t, revision.Revision("1"), newer.Properties.Components[2].Revision)

	// Updates to the components are in the actions
	require.Len(t, actions, 3)

	require.Contains(t, actions, "A")
	action := actions["A"]
	require.Equal(t, CreateWorkload, action.Operation)
	require.Equal(t, "A", action.ComponentName)
	require.Equal(t, app.Components["A"].RevisionHistory[1], *action.Definition)
	require.Equal(t, newer.Properties.Components[0], action.Instantiation)
	require.Nil(t, action.PreviousDefinition)
	require.Nil(t, action.PreviousInstanitation)

	require.Contains(t, actions, "B")
	action = actions["B"]
	require.Equal(t, CreateWorkload, action.Operation)
	require.Equal(t, "B", action.ComponentName)
	require.Equal(t, app.Components["B"].RevisionHistory[0], *action.Definition)
	require.Equal(t, newer.Properties.Components[1], action.Instantiation)
	require.Nil(t, action.PreviousDefinition)
	require.Nil(t, action.PreviousInstanitation)

	require.Contains(t, actions, "C")
	action = actions["C"]
	require.Equal(t, CreateWorkload, action.Operation)
	require.Equal(t, "C", action.ComponentName)
	require.Equal(t, app.Components["C"].RevisionHistory[0], *action.Definition)
	require.Equal(t, newer.Properties.Components[2], action.Instantiation)
	require.Nil(t, action.PreviousDefinition)
	require.Nil(t, action.PreviousInstanitation)
}

func Test_DeploymentUpdated_OneComponent_Deleted(t *testing.T) {
	app := db.NewApplication()
	app.Components["A"] = db.ComponentHistory{
		Revision: revision.Revision("2"),
		RevisionHistory: []db.ComponentRevision{
			{
				Kind:       "radius.dev/Container@v1alpha1",
				Revision:   revision.Revision("1"),
				Properties: *db.NewComponentProperties(),
			},
			{
				Kind:       "radius.dev/Container@v1alpha1",
				Revision:   revision.Revision("2"),
				Properties: *db.NewComponentProperties(),
			},
		},
	}

	older := db.NewDeployment()
	older.Properties.Components = []*db.DeploymentComponent{
		{
			ComponentName: "A",
			Revision:      revision.Revision("1"),
		},
	}

	newer := db.NewDeployment()
	newer.Properties.Components = []*db.DeploymentComponent{}

	rp := rp{
		meta: metadata.NewRegistry(),
	}

	actions, err := rp.computeDeploymentActions(app, older, newer)
	require.NoError(t, err)

	require.False(t, deploymentIsNoOp(actions))

	// Updates to the components are in the actions
	require.Len(t, actions, 1)

	require.Contains(t, actions, "A")
	action := actions["A"]
	require.Equal(t, DeleteWorkload, action.Operation)
	require.Equal(t, "A", action.ComponentName)
	require.Nil(t, action.Definition)
	require.Nil(t, action.Instantiation)
	require.Equal(t, app.Components["A"].RevisionHistory[0], *action.PreviousDefinition)
	require.Equal(t, older.Properties.Components[0], action.PreviousInstanitation)
}

func Test_DeploymentUpdated_OneComponent_NoAction(t *testing.T) {
	app := db.NewApplication()
	app.Components["A"] = db.ComponentHistory{
		Revision: revision.Revision("2"),
		RevisionHistory: []db.ComponentRevision{
			{
				Kind:       "radius.dev/Container@v1alpha1",
				Revision:   revision.Revision("1"),
				Properties: *db.NewComponentProperties(),
			},
			{
				Kind:       "radius.dev/Container@v1alpha1",
				Revision:   revision.Revision("2"),
				Properties: *db.NewComponentProperties(),
			},
		},
	}

	older := db.NewDeployment()
	older.Properties.Components = []*db.DeploymentComponent{
		{
			ComponentName: "A",
			Revision:      revision.Revision("2"),
		},
	}

	newer := db.NewDeployment()
	newer.Properties.Components = []*db.DeploymentComponent{
		{
			ComponentName: "A",
			Revision:      revision.Revision("2"),
		},
	}

	rp := rp{
		meta: metadata.NewRegistry(),
	}

	actions, err := rp.computeDeploymentActions(app, older, newer)
	require.NoError(t, err)

	require.True(t, deploymentIsNoOp(actions))

	// Updates to the components are in the actions
	require.Len(t, actions, 1)

	require.Contains(t, actions, "A")
	action := actions["A"]
	require.Equal(t, None, action.Operation)
	require.Equal(t, "A", action.ComponentName)
	require.Equal(t, app.Components["A"].RevisionHistory[1], *action.Definition)
	require.Equal(t, newer.Properties.Components[0], action.Instantiation)
	require.Equal(t, app.Components["A"].RevisionHistory[1], *action.PreviousDefinition)
	require.Equal(t, older.Properties.Components[0], action.PreviousInstanitation)
}

func Test_DeploymentUpdated_OneComponent_RevisionUpgraded(t *testing.T) {
	app := db.NewApplication()
	app.Components["A"] = db.ComponentHistory{
		Revision: revision.Revision("2"),
		RevisionHistory: []db.ComponentRevision{
			{
				Kind:       "radius.dev/Container@v1alpha1",
				Revision:   revision.Revision("1"),
				Properties: *db.NewComponentProperties(),
			},
			{
				Kind:       "radius.dev/Container@v1alpha1",
				Revision:   revision.Revision("2"),
				Properties: *db.NewComponentProperties(),
			},
		},
	}

	older := db.NewDeployment()
	older.Properties.Components = []*db.DeploymentComponent{
		{
			ComponentName: "A",
			Revision:      revision.Revision("1"),
		},
	}

	newer := db.NewDeployment()
	newer.Properties.Components = []*db.DeploymentComponent{
		{
			ComponentName: "A",
			Revision:      revision.Revision("2"),
		},
	}

	rp := rp{
		meta: metadata.NewRegistry(),
	}

	actions, err := rp.computeDeploymentActions(app, older, newer)
	require.NoError(t, err)

	require.False(t, deploymentIsNoOp(actions))

	// Updates to the components are in the actions
	require.Len(t, actions, 1)

	require.Contains(t, actions, "A")
	action := actions["A"]
	require.Equal(t, UpdateWorkload, action.Operation)
	require.Equal(t, "A", action.ComponentName)
	require.Equal(t, app.Components["A"].RevisionHistory[1], *action.Definition)
	require.Equal(t, newer.Properties.Components[0], action.Instantiation)
	require.Equal(t, app.Components["A"].RevisionHistory[0], *action.PreviousDefinition)
	require.Equal(t, older.Properties.Components[0], action.PreviousInstanitation)
}

func Test_DeploymentCreated_MultipleComponents_ServiceBinding(t *testing.T) {
	app := db.NewApplication()
	app.Components["A"] = db.ComponentHistory{
		Revision: revision.Revision("1"),
		RevisionHistory: []db.ComponentRevision{
			{
				Kind:     "radius.dev/Container@v1alpha1",
				Revision: revision.Revision("1"),
				Properties: db.ComponentProperties{
					Build: map[string]interface{}{},
					Run:   map[string]interface{}{},
					Provides: []db.ComponentProvides{
						{
							Name: "A",
							Kind: "http",
						},
					},
					DependsOn: []db.ComponentDependsOn{
						{
							Name: "B",
							Kind: "http",
						},
					},
				},
			},
		},
	}

	app.Components["B"] = db.ComponentHistory{
		Revision: revision.Revision("1"),
		RevisionHistory: []db.ComponentRevision{
			{
				Kind:     "radius.dev/Container@v1alpha1",
				Revision: revision.Revision("1"),
				Properties: db.ComponentProperties{
					Build: map[string]interface{}{},
					Run:   map[string]interface{}{},
					Provides: []db.ComponentProvides{
						{
							Name:          "B",
							Kind:          "http",
							ContainerPort: (func() *int { x := 80; return &x })(),
						},
					},
					DependsOn: []db.ComponentDependsOn{
						{
							Name: "A",
							Kind: "http",
						},
					},
				},
			},
		},
	}

	newer := db.NewDeployment()
	newer.Properties.Components = []*db.DeploymentComponent{
		{
			ComponentName: "A",
			Revision:      revision.Revision("1"),
		},
		{
			ComponentName: "B",
			Revision:      revision.Revision("1"),
		},
	}

	rp := rp{
		meta: metadata.NewRegistry(),
	}

	actions, err := rp.computeDeploymentActions(app, nil, newer)
	require.NoError(t, err)

	require.False(t, deploymentIsNoOp(actions))

	// Updates to the components are in the actions
	require.Len(t, actions, 2)

	require.Contains(t, actions, "A")
	action := actions["A"]
	require.Equal(t, CreateWorkload, action.Operation)
	require.Equal(t, "A", action.ComponentName)
	require.Equal(t, map[string]ServiceBinding{"B": {Name: "B", Kind: "http", Provider: "B"}}, action.ServiceBindings)

	require.Contains(t, actions, "B")
	action = actions["B"]
	require.Equal(t, CreateWorkload, action.Operation)
	require.Equal(t, "B", action.ComponentName)
	require.Equal(t, map[string]ServiceBinding{"A": {Name: "A", Kind: "http", Provider: "A"}}, action.ServiceBindings)
}

func Test_DeploymentUpdated_RenderRealisticContainerizedWorkload(t *testing.T) {
	app := db.NewApplication()
	app.Name = "radius/myapp"
	app.Components["A"] = db.ComponentHistory{
		Revision: revision.Revision("1"),
		RevisionHistory: []db.ComponentRevision{
			{
				Kind:     "radius.dev/Container@v1alpha1",
				Revision: revision.Revision("1"),
				Properties: db.ComponentProperties{
					Build: map[string]interface{}{},
					Run: map[string]interface{}{
						"container": map[string]interface{}{
							"image": "rynowak/frontend:0.5.0-dev",
							"env": []interface{}{
								map[string]interface{}{
									"name":  "SERVICE__BACKEND__HOST",
									"value": "backend",
								},
								map[string]interface{}{
									"name":  "SERVICE__BACKEND__PORT",
									"value": "80",
								},
							},
						},
					},
				},
			},
		},
	}

	newer := db.NewDeployment()
	newer.Properties.Components = []*db.DeploymentComponent{
		{
			ComponentName: "A",
			Revision:      revision.Revision("1"),
		},
	}

	rp := rp{
		meta: metadata.NewRegistry(),
	}

	actions, err := rp.computeDeploymentActions(app, nil, newer)
	require.NoError(t, err)

	require.False(t, deploymentIsNoOp(actions))

	// Updates to the components are in the actions
	require.Len(t, actions, 1)

	require.Contains(t, actions, "A")
	action := actions["A"]
	require.Equal(t, CreateWorkload, action.Operation)
	require.Equal(t, "A", action.ComponentName)
	require.Equal(t, app.Components["A"].RevisionHistory[0], *action.Definition)
	require.Equal(t, newer.Properties.Components[0], action.Instantiation)

	// unmarshall and validate it
	w := map[string]interface{}{}
	buf, err := action.Workload.MarshalJSON()
	require.NoError(t, err)

	err = json.Unmarshal(buf, &w)
	require.NoError(t, err)

	metadata := w["metadata"].(map[string]interface{})
	require.NotNil(t, metadata)
	require.Equal(t, "Container", w["kind"])
	require.Equal(t, "radius.dev/v1alpha1", w["apiVersion"])
	require.Equal(t, "myapp", metadata["namespace"])
	require.Equal(t, "A", metadata["name"])

	// Labels ignored

	spec := w["spec"]
	buf, err = json.Marshal(spec)
	require.NoError(t, err)

	workload := containerv1alpha1.ContainerWorkload{}
	err = json.Unmarshal(buf, &workload)
	require.NoError(t, err)

	cont := workload.Container
	require.Equal(t, "rynowak/frontend:0.5.0-dev", cont.Image)

	require.Len(t, cont.Environment, 2)
	require.Equal(t, "SERVICE__BACKEND__HOST", cont.Environment[0].Name)
	require.Equal(t, "backend", *cont.Environment[0].Value)
	require.Equal(t, "SERVICE__BACKEND__PORT", cont.Environment[1].Name)
	require.Equal(t, "80", *cont.Environment[1].Value)
}

func Test_DeploymentCreated_RenderContainerWithDapr(t *testing.T) {
	app := db.NewApplication()
	app.Name = "radius/myapp"
	app.Components["A"] = db.ComponentHistory{
		Revision: revision.Revision("1"),
		RevisionHistory: []db.ComponentRevision{
			{
				Kind:     "radius.dev/Container@v1alpha1",
				Revision: revision.Revision("1"),
				Properties: db.ComponentProperties{
					Build: map[string]interface{}{},
					Run: map[string]interface{}{
						"container": map[string]interface{}{
							"image": "rynowak/frontend:0.5.0-dev",
						},
					},
					Traits: []db.ComponentTrait{
						{
							Kind: "dapr.io/App@v1alpha1",
							Properties: map[string]interface{}{
								"appId":   "frontend",
								"appPort": 80,
							},
						},
					},
				},
			},
		},
	}

	newer := db.NewDeployment()
	newer.Properties.Components = []*db.DeploymentComponent{
		{
			ComponentName: "A",
			Revision:      revision.Revision("1"),
		},
	}

	rp := rp{
		meta: metadata.NewRegistry(),
	}

	actions, err := rp.computeDeploymentActions(app, nil, newer)
	require.NoError(t, err)

	require.False(t, deploymentIsNoOp(actions))

	// Updates to the components are in the actions
	require.Len(t, actions, 1)

	require.Contains(t, actions, "A")
	action := actions["A"]
	require.Equal(t, CreateWorkload, action.Operation)
	require.Equal(t, "A", action.ComponentName)
	require.Equal(t, app.Components["A"].RevisionHistory[0], *action.Definition)
	require.Equal(t, newer.Properties.Components[0], action.Instantiation)

	require.Len(t, action.Traits, 1)
	require.Equal(t, "dapr.io/App@v1alpha1", action.Traits[0].Kind)
	require.Equal(t, map[string]interface{}{
		"appId":   "frontend",
		"appPort": 80,
	}, action.Traits[0].Properties)
}
