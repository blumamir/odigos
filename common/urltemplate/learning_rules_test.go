package urltemplate

import (
	"fmt"
	"strconv"
	"testing"

	actionsapi "github.com/odigos-io/odigos/common/api/actions"
)

func TestFindLiveTrafficLearningRules_LeafBands(t *testing.T) {
	t.Run("high certainty leaf", func(t *testing.T) {
		rules := FindLiveTrafficLearningRules(BuildPathTrie(map[string]string{
			"/health": "101",
		}))
		assertTemplates(t, templatesOf(rules), []string{"/health"})
		if len(rules[0].Segments) != 0 {
			t.Fatalf("static leaf should have no templated segments, got %+v", rules[0].Segments)
		}
	})

	t.Run("skipped when total observations not above 100", func(t *testing.T) {
		rules := FindLiveTrafficLearningRules(BuildPathTrie(map[string]string{
			"/health": "100",
		}))
		if len(rules) != 0 {
			t.Fatalf("rules=%v, want none", templatesOf(rules))
		}
	})

	t.Run("moderate leaf when node total above 100 but leaf obs in band", func(t *testing.T) {
		rules := FindLiveTrafficLearningRules(BuildPathTrie(map[string]string{
			"/users":       "50",
			"/users/extra": "60",
		}))
		assertTemplates(t, templatesOf(rules), []string{"/users"})
	})

	t.Run("leaf under 25 not emitted even if total high", func(t *testing.T) {
		rules := FindLiveTrafficLearningRules(BuildPathTrie(map[string]string{
			"/users":       "20",
			"/users/extra": "90",
		}))
		if len(rules) != 0 {
			t.Fatalf("rules=%v, want none", templatesOf(rules))
		}
	})
}

func TestFindLiveTrafficLearningRules_SingleChildChain(t *testing.T) {
	rules := FindLiveTrafficLearningRules(BuildPathTrie(map[string]string{
		"/api/v1/health": "101",
	}))
	assertTemplates(t, templatesOf(rules), []string{"/api/v1/health"})
}

func TestFindLiveTrafficLearningRules_RecurseEachUpToFive(t *testing.T) {
	rules := FindLiveTrafficLearningRules(BuildPathTrie(map[string]string{
		"/a/x": "101",
		"/a/y": "101",
		"/a/z": "101",
	}))
	assertTemplates(t, templatesOf(rules), []string{"/a/x", "/a/y", "/a/z"})
}

func TestFindLiveTrafficLearningRules_WildcardAndTemplate(t *testing.T) {
	t.Run("wildcard moderate for 6-10 children", func(t *testing.T) {
		// Use non-ID segment names so branching bands apply (not regex short-circuit).
		paths := leafPathsWithAlpha("/items", 8, 20)
		rules := FindLiveTrafficLearningRules(BuildPathTrie(paths))
		assertTemplates(t, templatesOf(rules), []string{"/items/*"})
		if len(rules[0].Segments) != 1 || rules[0].Segments[0].TemplateName != "*" {
			t.Fatalf("segments=%+v", rules[0].Segments)
		}
		if rules[0].Segments[0].Certainty != actionsapi.URLTemplatizationRuleCertaintyModerate {
			t.Fatalf("segment certainty=%q, want Moderate", rules[0].Segments[0].Certainty)
		}
	})

	t.Run("wildcard high for 11-25 children", func(t *testing.T) {
		paths := leafPathsWithAlpha("/items", 12, 10)
		rules := FindLiveTrafficLearningRules(BuildPathTrie(paths))
		assertTemplates(t, templatesOf(rules), []string{"/items/*"})
		if rules[0].Segments[0].Certainty != actionsapi.URLTemplatizationRuleCertaintyHigh {
			t.Fatalf("segment certainty=%q, want High", rules[0].Segments[0].Certainty)
		}
	})

	t.Run("template id for more than 25 children", func(t *testing.T) {
		paths := leafPathsWithAlpha("/users", 30, 5)
		rules := FindLiveTrafficLearningRules(BuildPathTrie(paths))
		assertTemplates(t, templatesOf(rules), []string{"/users/{id}"})
		if len(rules[0].Segments) != 1 || rules[0].Segments[0].TemplateName != "id" {
			t.Fatalf("segments=%+v", rules[0].Segments)
		}
		if rules[0].Segments[0].Certainty != actionsapi.URLTemplatizationRuleCertaintyHigh {
			t.Fatalf("segment certainty=%q, want High", rules[0].Segments[0].Certainty)
		}
		if got := len(rules[0].Segments[0].Examples); got != maxTemplatedExamples {
			t.Fatalf("examples=%d, want %d", got, maxTemplatedExamples)
		}
	})
}

func TestFindLiveTrafficLearningRules_RegexHighCertainty(t *testing.T) {
	t.Run("numeric children templated as id with high certainty", func(t *testing.T) {
		rules := FindLiveTrafficLearningRules(BuildPathTrie(map[string]string{
			"/users/101": "40",
			"/users/202": "40",
			"/users/303": "40",
		}))
		assertTemplates(t, templatesOf(rules), []string{"/users/{id}"})
		if len(rules[0].Segments) != 1 || rules[0].Segments[0].TemplateName != "id" {
			t.Fatalf("segments=%+v", rules[0].Segments)
		}
		if rules[0].Segments[0].Certainty != actionsapi.URLTemplatizationRuleCertaintyHigh {
			t.Fatalf("segment certainty=%q, want High", rules[0].Segments[0].Certainty)
		}
	})

	t.Run("date children templated as date with high certainty", func(t *testing.T) {
		rules := FindLiveTrafficLearningRules(BuildPathTrie(map[string]string{
			"/reports/2025-01-01": "40",
			"/reports/2025-01-02": "40",
			"/reports/2025-01-03": "40",
		}))
		assertTemplates(t, templatesOf(rules), []string{"/reports/{date}"})
		if len(rules[0].Segments) != 1 || rules[0].Segments[0].TemplateName != "date" {
			t.Fatalf("segments=%+v", rules[0].Segments)
		}
		if rules[0].Segments[0].Certainty != actionsapi.URLTemplatizationRuleCertaintyHigh {
			t.Fatalf("segment certainty=%q, want High", rules[0].Segments[0].Certainty)
		}
	})

	t.Run("single regex-matching child templated with high certainty", func(t *testing.T) {
		rules := FindLiveTrafficLearningRules(BuildPathTrie(map[string]string{
			"/users/550e8400-e29b-41d4-a716-446655440000": "101",
		}))
		assertTemplates(t, templatesOf(rules), []string{"/users/{id}"})
		if rules[0].Segments[0].Certainty != actionsapi.URLTemplatizationRuleCertaintyHigh {
			t.Fatalf("segment certainty=%q, want High", rules[0].Segments[0].Certainty)
		}
	})
}

func TestFindLiveTrafficLearningRules_MergeContinues(t *testing.T) {
	paths := map[string]string{}
	for i := 0; i < 30; i++ {
		paths[fmt.Sprintf("/users/%d/orders", i)] = "5"
	}
	rules := FindLiveTrafficLearningRules(BuildPathTrie(paths))
	assertTemplates(t, templatesOf(rules), []string{"/users/{id}/orders"})
	if len(rules[0].Segments) != 1 || rules[0].Segments[0].Certainty != actionsapi.URLTemplatizationRuleCertaintyHigh {
		t.Fatalf("segments=%+v", rules[0].Segments)
	}
}

func TestFindLiveTrafficLearningRules_ChildrenNeedEnoughObservations(t *testing.T) {
	rules := FindLiveTrafficLearningRules(BuildPathTrie(map[string]string{
		"/api":    "101",
		"/api/v1": "50",
		"/api/v2": "50",
	}))
	assertTemplates(t, templatesOf(rules), []string{"/api"})
}

func leafPathsWithCount(prefix string, n int, count int64) map[string]string {
	paths := make(map[string]string, n)
	for i := 0; i < n; i++ {
		seg := strconv.Itoa(i)
		path := "/" + seg
		if prefix != "" {
			path = prefix + "/" + seg
		}
		paths[path] = strconv.FormatInt(count, 10)
	}
	if len(paths) != n {
		panic(fmt.Sprintf("expected %d paths, got %d", n, len(paths)))
	}
	return paths
}

// leafPathsWithAlpha builds n leaf paths with alphabetic segment names that do not
// match built-in ID/date/email heuristics (so branching bands apply).
func leafPathsWithAlpha(prefix string, n int, count int64) map[string]string {
	paths := make(map[string]string, n)
	for i := 0; i < n; i++ {
		seg := fmt.Sprintf("seg%c", 'a'+i%26)
		if i >= 26 {
			seg = fmt.Sprintf("seg%c%d", 'a'+i%26, i/26)
		}
		path := "/" + seg
		if prefix != "" {
			path = prefix + "/" + seg
		}
		paths[path] = strconv.FormatInt(count, 10)
	}
	if len(paths) != n {
		panic(fmt.Sprintf("expected %d paths, got %d", n, len(paths)))
	}
	return paths
}

func templatesOf(rules []actionsapi.URLTemplatizationLearnedRule) []string {
	out := make([]string, 0, len(rules))
	for _, r := range rules {
		out = append(out, r.Template)
	}
	return out
}

func assertTemplates(t *testing.T, got, want []string) {
	t.Helper()
	if len(got) != len(want) {
		t.Fatalf("templates=%v, want %v", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("templates=%v, want %v", got, want)
		}
	}
}
