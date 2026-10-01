package odigosurltemplateprocessor

import (
	"strings"

	"github.com/odigos-io/odigos/common/urltemplate"
)

func attemptTemplateWithRule(pathSegments []string, rule urltemplate.PathRule) (string, bool) {
	// already verified that the len of the lists match pre calling this function
	if !rule.IsPathSegmentsMatching(pathSegments) {
		return "", false
	}

	result := make([]string, 0, len(rule.Segments))
	for i, segment := range rule.Segments {
		if segment.TemplateName != "" {
			result = append(result, "{"+segment.TemplateName+"}")
		} else if segment.Wildcard {
			// untemplated value, keep whatever is in the path (assumes low cardinality)
			result = append(result, pathSegments[i])
		} else {
			result = append(result, segment.StaticString)
		}
	}

	return strings.Join(result, "/"), true
}

// return the name to use for templatization "id" / "date" etc which will be embedded in the template
// as {id} / {date} etc
// empty string as return value means that the segment is not a templated id
func getSegmentTemplatizationString(segment string, customIds []internalCustomIdConfig) string {
	for _, customRegexp := range customIds {
		if customRegexp.Regexp.MatchString(segment) {
			return customRegexp.Name
		}
	}
	return urltemplate.SegmentTemplateName(segment)
}

// This function will replace all segments that matches a number or uuid with "{id}"
func defaultTemplatizeURLPath(pathSegments []string, customIdsRegexp []internalCustomIdConfig) (string, bool) {
	templated := false
	// avoid modifying the original segments slice
	templatizedSegments := make([]string, len(pathSegments))
	for i, segment := range pathSegments {
		if templateName := getSegmentTemplatizationString(segment, customIdsRegexp); templateName != "" {
			templatizedSegments[i] = "{" + templateName + "}"
			templated = true
		} else {
			templatizedSegments[i] = segment
		}
	}
	if !templated {
		return "", false
	}

	templatedPath := strings.Join(templatizedSegments, "/")
	return templatedPath, true
}

// check if specific path segments match any of the custom templatization rules
// if so, return the templated url and true
// if not, return false
func applyCustomRulesForTemplatization(pathSegments []string, rules map[int][]urltemplate.PathRule, hadLeadingSlash bool) (string, bool) {
	ruleList, found := rules[len(pathSegments)]
	if !found {
		return "", false
	}
	for _, rule := range ruleList {
		if templatedUrl, matched := attemptTemplateWithRule(pathSegments, rule); matched {
			if hadLeadingSlash {
				templatedUrl = "/" + templatedUrl
			}
			return templatedUrl, true
		}
	}
	return "", false
}
