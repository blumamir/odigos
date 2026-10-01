package services

import (
	"testing"

	"github.com/odigos-io/odigos/common/api/actions"
	"github.com/odigos-io/odigos/frontend/graph/model"
	"github.com/stretchr/testify/require"
)

func TestFilterPathsNotMatchingRules(t *testing.T) {
	rules := parseTemplatizationPathRules([]string{"/users/{id}", "/health"})
	cn := "app"
	paths := []*model.UnmatchedURLPath{
		{Path: "/users/123", Count: 10, ContainerName: &cn},
		{Path: "/users/456/orders", Count: 5, ContainerName: &cn},
		{Path: "/health", Count: 3, ContainerName: &cn},
		{Path: "/api/v1/items", Count: 2, ContainerName: &cn},
	}

	filtered := filterPathsNotMatchingRules(paths, rules)
	require.Len(t, filtered, 2)
	require.Equal(t, "/users/456/orders", filtered[0].Path)
	require.Equal(t, "/api/v1/items", filtered[1].Path)
}

func TestFilterPathsNotMatchingRules_NoRules(t *testing.T) {
	paths := []*model.UnmatchedURLPath{{Path: "/users/123", Count: 1}}
	require.Equal(t, paths, filterPathsNotMatchingRules(paths, nil))
}

func TestUrlTemplatizationConfigToExistingModel(t *testing.T) {
	cfg := &actions.UrlTemplatizationConfig{
		Templates: []string{"/users/{id}", "/health"},
		Default: &actions.DefaultTemplatizationConfig{
			Disabled: false,
			SkipPolicy: &actions.DefaultTemplatizationSkipPolicyConfig{
				SkipForNonSuccessCodes: true,
				SkipHttpStatusCodes:    []int{404},
			},
		},
	}

	got := urlTemplatizationConfigToExistingModel("app", cfg)
	require.Equal(t, "app", got.ContainerName)
	require.Equal(t, []string{"/users/{id}", "/health"}, got.Templates)
	require.NotNil(t, got.Default)
	require.False(t, got.Default.Disabled)
	require.NotNil(t, got.Default.SkipPolicy)
	require.NotNil(t, got.Default.SkipPolicy.SkipForNonSuccessCodes)
	require.True(t, *got.Default.SkipPolicy.SkipForNonSuccessCodes)
	require.Equal(t, []int{404}, got.Default.SkipPolicy.SkipHTTPStatusCodes)
}

func TestUrlTemplatizationConfigToExistingModel_NilDefault(t *testing.T) {
	got := urlTemplatizationConfigToExistingModel("app", &actions.UrlTemplatizationConfig{
		Templates: []string{"/users/{id}"},
	})
	require.Equal(t, []string{"/users/{id}"}, got.Templates)
	require.Nil(t, got.Default)
}
