package urltemplate

import (
	"fmt"
	"sort"
	"strings"

	actionsapi "github.com/odigos-io/odigos/common/api/actions"
)

const (
	// minNodeTotalObservations is the minimum TotalObservations for a node to be
	// considered (#1). Below this we wait for more data.
	minNodeTotalObservations int64 = 100

	// Leaf observation bands (#2).
	minLeafModerateObservations int64 = 25
	minLeafHighObservations     int64 = 100

	// Children must collectively exceed this (TotalObservations - ObservationCount).
	minChildrenTotalObservations int64 = 100

	// Branching bands (#3): 1 → recurse; 2–5 → recurse each;
	// 6–10 → "*" Moderate; 11–25 → "*" High; >25 → "{id}" High.
	maxRecurseEachChildren  = 5
	maxWildcardLowChildren  = 10
	maxWildcardHighChildren = 25

	templatedSegmentName = "id"
	wildcardSegmentName  = "*"
	maxTemplatedExamples = 5
)

// FindLiveTrafficLearningRules walks each entry node (first path segment) and builds
// rule suggestions recursively.
func FindLiveTrafficLearningRules(trie *PathTrie) []actionsapi.URLTemplatizationLearnedRule {
	if trie == nil || trie.Root == nil {
		return nil
	}

	var rules []actionsapi.URLTemplatizationLearnedRule
	for _, entry := range sortedChildList(trie.Root) {
		rules = append(rules, walkNode(entry, "/"+entry.Segment, nil)...)
	}
	return dedupeRulesByTemplate(rules)
}

// walkNode implements the recursive live traffic learning walk for one trie node.
func walkNode(node *PathTrieNode, path string, segments []actionsapi.URLTemplatizationSegment) []actionsapi.URLTemplatizationLearnedRule {
	if node == nil {
		return nil
	}

	// #1 Enough observations on this node (including all descendants).
	if node.TotalObservations <= minNodeTotalObservations {
		return nil
	}

	var rules []actionsapi.URLTemplatizationLearnedRule

	// #2 Leaf: an observed path terminated here (may still have children).
	if node.Leaf {
		if rule, ok := leafRule(node, path, segments); ok {
			rules = append(rules, rule)
		}
	}

	// #3 Continue to children when they collectively have enough observations.
	childrenObs := node.TotalObservations - node.ObservationCount
	if childrenObs <= minChildrenTotalObservations {
		return rules
	}

	children := sortedChildList(node)
	n := len(children)
	if n == 0 {
		return rules
	}

	// If every distinct child value matches a built-in ID/date/email heuristic,
	// treat this segment as templated with high certainty regardless of branching.
	if templateName, ok := childrenRegexTemplateName(children); ok {
		rules = append(rules, walkMergedDynamic(
			children, path, segments,
			"{"+templateName+"}",
			actionsapi.URLTemplatizationRuleCertaintyHigh,
		)...)
		return rules
	}

	switch {
	case n == 1:
		child := children[0]
		rules = append(rules, walkNode(child, joinPath(path, child.Segment), segments)...)

	case n <= maxRecurseEachChildren:
		for _, child := range children {
			rules = append(rules, walkNode(child, joinPath(path, child.Segment), segments)...)
		}

	case n <= maxWildcardLowChildren:
		rules = append(rules, walkMergedDynamic(children, path, segments, wildcardSegmentName, actionsapi.URLTemplatizationRuleCertaintyModerate)...)

	case n <= maxWildcardHighChildren:
		rules = append(rules, walkMergedDynamic(children, path, segments, wildcardSegmentName, actionsapi.URLTemplatizationRuleCertaintyHigh)...)

	default:
		rules = append(rules, walkMergedDynamic(children, path, segments, "{"+templatedSegmentName+"}", actionsapi.URLTemplatizationRuleCertaintyHigh)...)
	}

	return rules
}

// childrenRegexTemplateName returns a template name when every child segment matches
// a built-in heuristic. If they agree on one name, that name is used; otherwise "id".
func childrenRegexTemplateName(children []*PathTrieNode) (string, bool) {
	if len(children) == 0 {
		return "", false
	}
	first := SegmentTemplateName(children[0].Segment)
	if first == "" {
		return "", false
	}
	same := true
	for _, child := range children[1:] {
		name := SegmentTemplateName(child.Segment)
		if name == "" {
			return "", false
		}
		if name != first {
			same = false
		}
	}
	if same {
		return first, true
	}
	return templatedSegmentName, true
}

func walkMergedDynamic(
	children []*PathTrieNode,
	path string,
	segments []actionsapi.URLTemplatizationSegment,
	segmentLiteral string,
	certainty actionsapi.URLTemplatizationRuleCertainty,
) []actionsapi.URLTemplatizationLearnedRule {
	merged := mergeChildren(children, segmentLiteral)
	nextSegments := appendSegment(segments, segmentLiteral, children, certainty)
	nextPath := joinPath(path, segmentLiteral)
	rules := walkNode(merged, nextPath, nextSegments)
	for i := range rules {
		rules[i].Reason = annotateDynamicReason(rules[i].Reason, segmentLiteral, len(children))
	}
	return rules
}

func leafRule(node *PathTrieNode, path string, segments []actionsapi.URLTemplatizationSegment) (actionsapi.URLTemplatizationLearnedRule, bool) {
	obs := node.ObservationCount
	switch {
	case obs > minLeafHighObservations:
		return actionsapi.URLTemplatizationLearnedRule{
			Template: path,
			Reason: fmt.Sprintf(
				"Static endpoint with %d observations at this path (threshold > %d)",
				obs, minLeafHighObservations,
			),
			Segments: copySegments(segments),
		}, true
	case obs > minLeafModerateObservations:
		return actionsapi.URLTemplatizationLearnedRule{
			Template: path,
			Reason: fmt.Sprintf(
				"Static endpoint with %d observations at this path (threshold > %d and ≤ %d)",
				obs, minLeafModerateObservations, minLeafHighObservations,
			),
			Segments: copySegments(segments),
		}, true
	default:
		return actionsapi.URLTemplatizationLearnedRule{}, false
	}
}

func annotateDynamicReason(reason, segmentLiteral string, childCount int) string {
	kind := "templated segment"
	if segmentLiteral == wildcardSegmentName {
		kind = "wildcard segment"
	}
	return fmt.Sprintf("%s; %s %q from %d distinct child values", reason, kind, segmentLiteral, childCount)
}

func mergeChildren(children []*PathTrieNode, segmentName string) *PathTrieNode {
	merged := &PathTrieNode{
		Segment:  segmentName,
		Children: make(map[string]*PathTrieNode),
	}
	for _, child := range children {
		if child.Leaf {
			merged.Leaf = true
			merged.ObservationCount += child.ObservationCount
		}
		for name, grandchild := range child.Children {
			if existing, ok := merged.Children[name]; ok {
				merged.Children[name] = mergeNodes(existing, grandchild)
			} else {
				merged.Children[name] = cloneNode(grandchild)
			}
		}
	}
	merged.recompute()
	return merged
}

func mergeNodes(a, b *PathTrieNode) *PathTrieNode {
	out := &PathTrieNode{
		Segment:          a.Segment,
		Children:         make(map[string]*PathTrieNode),
		Leaf:             a.Leaf || b.Leaf,
		ObservationCount: a.ObservationCount + b.ObservationCount,
	}
	for name, child := range a.Children {
		out.Children[name] = cloneNode(child)
	}
	for name, child := range b.Children {
		if existing, ok := out.Children[name]; ok {
			out.Children[name] = mergeNodes(existing, child)
		} else {
			out.Children[name] = cloneNode(child)
		}
	}
	out.recompute()
	return out
}

func cloneNode(n *PathTrieNode) *PathTrieNode {
	if n == nil {
		return nil
	}
	out := &PathTrieNode{
		Segment:          n.Segment,
		Children:         make(map[string]*PathTrieNode, len(n.Children)),
		Leaf:             n.Leaf,
		ObservationCount: n.ObservationCount,
	}
	for name, child := range n.Children {
		out.Children[name] = cloneNode(child)
	}
	out.recompute()
	return out
}

func appendSegment(
	segments []actionsapi.URLTemplatizationSegment,
	segmentLiteral string,
	children []*PathTrieNode,
	certainty actionsapi.URLTemplatizationRuleCertainty,
) []actionsapi.URLTemplatizationSegment {
	name := segmentLiteral
	if strings.HasPrefix(segmentLiteral, "{") && strings.HasSuffix(segmentLiteral, "}") {
		name = strings.TrimSuffix(strings.TrimPrefix(segmentLiteral, "{"), "}")
	}
	examples := make([]string, 0, len(children))
	for _, child := range children {
		examples = append(examples, child.Segment)
	}
	sort.Strings(examples)
	if len(examples) > maxTemplatedExamples {
		examples = examples[:maxTemplatedExamples]
	}
	out := make([]actionsapi.URLTemplatizationSegment, len(segments), len(segments)+1)
	copy(out, segments)
	out = append(out, actionsapi.URLTemplatizationSegment{
		TemplateName: name,
		Certainty:    certainty,
		Examples:     examples,
	})
	return out
}

func copySegments(segments []actionsapi.URLTemplatizationSegment) []actionsapi.URLTemplatizationSegment {
	if len(segments) == 0 {
		return nil
	}
	out := make([]actionsapi.URLTemplatizationSegment, len(segments))
	copy(out, segments)
	return out
}

func sortedChildList(n *PathTrieNode) []*PathTrieNode {
	if n == nil || len(n.Children) == 0 {
		return nil
	}
	names := make([]string, 0, len(n.Children))
	for name := range n.Children {
		names = append(names, name)
	}
	sort.Strings(names)
	out := make([]*PathTrieNode, 0, len(names))
	for _, name := range names {
		out = append(out, n.Children[name])
	}
	return out
}

func joinPath(prefix, segment string) string {
	if prefix == "" {
		return "/" + segment
	}
	return prefix + "/" + segment
}

func dedupeRulesByTemplate(rules []actionsapi.URLTemplatizationLearnedRule) []actionsapi.URLTemplatizationLearnedRule {
	seen := make(map[string]struct{}, len(rules))
	out := make([]actionsapi.URLTemplatizationLearnedRule, 0, len(rules))
	for _, rule := range rules {
		if _, ok := seen[rule.Template]; ok {
			continue
		}
		seen[rule.Template] = struct{}{}
		out = append(out, rule)
	}
	return out
}
