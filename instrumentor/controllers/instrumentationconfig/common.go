package instrumentationconfig

import (
	"context"

	odigosv1alpha1 "github.com/odigos-io/odigos/api/odigos/v1alpha1"
	rulesv1alpha1 "github.com/odigos-io/odigos/api/rules/v1alpha1"
	"github.com/odigos-io/odigos/common"
	"github.com/odigos-io/odigos/k8sutils/pkg/workload"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type instrumentationRules struct {
	payloadCollection *rulesv1alpha1.PayloadCollectionList
}

func getAllInstrumentationRules(ctx context.Context, client client.Client) (instrumentationRules, error) {

	payloadCollectionRules := &rulesv1alpha1.PayloadCollectionList{}
	err := client.List(ctx, payloadCollectionRules)
	if err != nil {
		return instrumentationRules{}, err
	}

	return instrumentationRules{
		payloadCollection: payloadCollectionRules,
	}, nil
}

func updateInstrumentationConfigForWorkload(ic *odigosv1alpha1.InstrumentationConfig, ia *odigosv1alpha1.InstrumentedApplication, rules instrumentationRules) error {

	workloadName, workloadKind, err := workload.ExtractWorkloadInfoFromRuntimeObjectName(ia.Name)
	if err != nil {
		return err
	}
	workload := workload.PodWorkload{
		Name:      workloadName,
		Namespace: ia.Namespace,
		Kind:      workloadKind,
	}

	// delete all existing sdk configs to re-calculate them
	ic.Spec.SdkConfigs = []odigosv1alpha1.SdkConfig{}

	// create an empty sdk config for each detected programming language
	for _, container := range ia.Spec.RuntimeDetails {
		containerLanguage := container.Language
		if containerLanguage == common.IgnoredProgrammingLanguage || containerLanguage == common.UnknownProgrammingLanguage {
			continue
		}
		createDefaultSdkConfig(ic, containerLanguage)
	}

	// iterate over all the payload collection rules, and update the instrumentation config accordingly
	if rules.payloadCollection != nil {
		for i := range rules.payloadCollection.Items {
			rule := &rules.payloadCollection.Items[i]
			if rule.Spec.Disabled {
				continue
			}
			// filter out rules where the workload does not match
			participating := isWorkloadParticipatingInRule(workload, rule)
			if !participating {
				continue
			}

			for i := range ic.Spec.SdkConfigs {
				if rule.Spec.InstrumentationLibraries == nil { // nil means a rule in SDK level, that applies unless overridden by library level rule
					ic.Spec.SdkConfigs[i].DefaultHttpRequestPayloadCollection = mergeHttpPayloadCollectionRules(ic.Spec.SdkConfigs[i].DefaultHttpRequestPayloadCollection, rule.Spec.HttpRequest)
					ic.Spec.SdkConfigs[i].DefaultHttpResponsePayloadCollection = mergeHttpPayloadCollectionRules(ic.Spec.SdkConfigs[i].DefaultHttpResponsePayloadCollection, rule.Spec.HttpResponse)
					ic.Spec.SdkConfigs[i].DefaultDbQueryPayloadCollection = mergeDbPayloadCollectionRules(ic.Spec.SdkConfigs[i].DefaultDbQueryPayloadCollection, rule.Spec.DbQuery)
				} else {
					for _, library := range *rule.Spec.InstrumentationLibraries {
						if library.Language != ic.Spec.SdkConfigs[i].Language {
							continue
						}
						libraryConfig := findOrCreateSdkLibraryConfig(&ic.Spec.SdkConfigs[i], library)
						libraryConfig.HttpRequestPayloadCollection = mergeHttpPayloadCollectionRules(libraryConfig.HttpRequestPayloadCollection, rule.Spec.HttpRequest)
						libraryConfig.HttpResponsePayloadCollection = mergeHttpPayloadCollectionRules(libraryConfig.HttpResponsePayloadCollection, rule.Spec.HttpResponse)
						libraryConfig.DbQueryPayloadCollection = mergeDbPayloadCollectionRules(libraryConfig.DbQueryPayloadCollection, rule.Spec.DbQuery)
					}
				}
			}
		}
	}

	return nil
}

// returns a pointer to the instrumentation library config, creating it if it does not exist
// the pointer can be used to modify the config
func findOrCreateSdkLibraryConfig(sdkConfig *odigosv1alpha1.SdkConfig, library rulesv1alpha1.InstrumentationLibraryId) *odigosv1alpha1.InstrumentationLibraryConfig {
	if library.Language != sdkConfig.Language {
		return nil
	}

	for i, libConfig := range sdkConfig.InstrumentationLibraryConfigs {
		if libConfig.InstrumentationLibraryId.InstrumentationLibraryName == library.Name &&
			libConfig.InstrumentationLibraryId.SpanKind == library.SpanKind {

			// if already present, return a pointer to it which can be modified by the caller
			return &sdkConfig.InstrumentationLibraryConfigs[i]
		}
	}
	newLibConfig := odigosv1alpha1.InstrumentationLibraryConfig{
		InstrumentationLibraryId: odigosv1alpha1.InstrumentationLibraryId{
			InstrumentationLibraryName: library.Name,
			SpanKind:                   library.SpanKind,
		},
	}
	sdkConfig.InstrumentationLibraryConfigs = append(sdkConfig.InstrumentationLibraryConfigs, newLibConfig)
	return &sdkConfig.InstrumentationLibraryConfigs[len(sdkConfig.InstrumentationLibraryConfigs)-1]
}

func createDefaultSdkConfig(ic *odigosv1alpha1.InstrumentationConfig, containerLanguage common.ProgrammingLanguage) {
	// if the language is already present, do nothing
	for _, sdkConfig := range ic.Spec.SdkConfigs {
		if sdkConfig.Language == containerLanguage {
			return
		}
	}
	ic.Spec.SdkConfigs = append(ic.Spec.SdkConfigs, odigosv1alpha1.SdkConfig{
		Language: containerLanguage,
	})
}

// naive implementation, can be optimized.
// assumption is that the list of workloads is small
func isWorkloadParticipatingInRule(workload workload.PodWorkload, rule *rulesv1alpha1.PayloadCollection) bool {
	// nil means all workloads are participating
	if rule.Spec.Workloads == nil {
		return true
	}
	for _, allowedWorkload := range *rule.Spec.Workloads {
		if allowedWorkload == workload {
			return true
		}
	}
	return false
}

func mergeHttpPayloadCollectionRules(rule1 *rulesv1alpha1.HttpPayloadCollectionRule, rule2 *rulesv1alpha1.HttpPayloadCollectionRule) *rulesv1alpha1.HttpPayloadCollectionRule {

	// nil means a rules has not yet been set, so return the other rule
	if rule1 == nil {
		return rule2
	} else if rule2 == nil {
		return rule1
	}

	// merge of the 2 non nil rules
	mergedRules := rulesv1alpha1.HttpPayloadCollectionRule{}

	// MimeTypes is extended to include both. nil means "all" so treat it as such
	if rule1.MimeTypes == nil || rule2.MimeTypes == nil {
		mergedRules.MimeTypes = nil
	} else {
		mergeMimeTypeMap := make(map[string]struct{})
		for _, mimeType := range *rule1.MimeTypes {
			mergeMimeTypeMap[mimeType] = struct{}{}
		}
		for _, mimeType := range *rule2.MimeTypes {
			mergeMimeTypeMap[mimeType] = struct{}{}
		}
		mergedMimeTypeSlice := make([]string, 0, len(mergeMimeTypeMap))
		for mimeType := range mergeMimeTypeMap {
			mergedMimeTypeSlice = append(mergedMimeTypeSlice, mimeType)
		}
		mergedRules.MimeTypes = &mergedMimeTypeSlice
	}

	// MaxPayloadLength - choose the smallest value, as this is the maximum allowed
	if rule1.MaxPayloadLength == nil {
		mergedRules.MaxPayloadLength = rule2.MaxPayloadLength
	} else if rule2.MaxPayloadLength == nil {
		mergedRules.MaxPayloadLength = rule1.MaxPayloadLength
	} else {
		if *rule1.MaxPayloadLength < *rule2.MaxPayloadLength {
			mergedRules.MaxPayloadLength = rule1.MaxPayloadLength
		} else {
			mergedRules.MaxPayloadLength = rule2.MaxPayloadLength
		}
	}

	// DropPartialPayloads - if any of the rules is set to drop, the merged rule will drop
	if rule1.DropPartialPayloads == nil {
		mergedRules.DropPartialPayloads = rule2.DropPartialPayloads
	} else if rule2.DropPartialPayloads == nil {
		mergedRules.DropPartialPayloads = rule1.DropPartialPayloads
	} else {
		mergedRules.DropPartialPayloads = boolPtr(*rule1.DropPartialPayloads || *rule2.DropPartialPayloads)
	}

	return &mergedRules
}

func mergeDbPayloadCollectionRules(rule1 *rulesv1alpha1.DbQueryPayloadCollectionRule, rule2 *rulesv1alpha1.DbQueryPayloadCollectionRule) *rulesv1alpha1.DbQueryPayloadCollectionRule {
	if rule1 == nil {
		return rule2
	} else if rule2 == nil {
		return rule1
	}

	mergedRules := rulesv1alpha1.DbQueryPayloadCollectionRule{}

	// MaxPayloadLength - choose the smallest value, as this is the maximum allowed
	if rule1.MaxPayloadLength == nil {
		mergedRules.MaxPayloadLength = rule2.MaxPayloadLength
	} else if rule2.MaxPayloadLength == nil {
		mergedRules.MaxPayloadLength = rule1.MaxPayloadLength
	} else {
		if *rule1.MaxPayloadLength < *rule2.MaxPayloadLength {
			mergedRules.MaxPayloadLength = rule1.MaxPayloadLength
		} else {
			mergedRules.MaxPayloadLength = rule2.MaxPayloadLength
		}
	}

	// DropPartialPayloads - if any of the rules is set to drop, the merged rule will drop
	if rule1.DropPartialPayloads == nil {
		mergedRules.DropPartialPayloads = rule2.DropPartialPayloads
	} else if rule2.DropPartialPayloads == nil {
		mergedRules.DropPartialPayloads = rule1.DropPartialPayloads
	} else {
		mergedRules.DropPartialPayloads = boolPtr(*rule1.DropPartialPayloads || *rule2.DropPartialPayloads)
	}

	return &mergedRules
}

func boolPtr(b bool) *bool {
	return &b
}