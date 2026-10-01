package services

import (
	"context"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/odigos-io/odigos/api/k8sconsts"
	odigosv1 "github.com/odigos-io/odigos/api/odigos/v1alpha1"
	"github.com/odigos-io/odigos/common"
	actionsapi "github.com/odigos-io/odigos/common/api/actions"
	"github.com/odigos-io/odigos/common/urltemplate"
	"github.com/odigos-io/odigos/frontend/graph/model"
	"github.com/odigos-io/odigos/frontend/kube"
	"github.com/odigos-io/odigos/k8sutils/pkg/env"
	"github.com/odigos-io/odigos/k8sutils/pkg/workload"
	"github.com/redis/go-redis/v9"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	unmatchedRedisKeyPrefixServer = "odigos:urltempl:unmatched:server"
	unmatchedRedisKeyPrefixClient = "odigos:urltempl:unmatched:client"
)

var (
	unmatchedRedisMu sync.Mutex
	unmatchedRedis   *redis.Client
)

func unmatchedRedisClient() (*redis.Client, error) {
	unmatchedRedisMu.Lock()
	defer unmatchedRedisMu.Unlock()
	if unmatchedRedis != nil {
		return unmatchedRedis, nil
	}
	ns := env.GetCurrentNamespace()
	addr := k8sconsts.OdigosCacheEndpoint(ns)
	rdb := redis.NewClient(&redis.Options{
		Addr:         addr,
		DialTimeout:  2 * time.Second,
		ReadTimeout:  2 * time.Second,
		WriteTimeout: 2 * time.Second,
	})
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	if err := rdb.Ping(ctx).Err(); err != nil {
		_ = rdb.Close()
		return nil, fmt.Errorf("connect to cache at %s: %w", addr, err)
	}
	unmatchedRedis = rdb
	return unmatchedRedis, nil
}

func emptyUnmatchedUrlPaths() *model.UnmatchedURLPaths {
	return &model.UnmatchedURLPaths{
		Server:                 []*model.UnmatchedURLPath{},
		Client:                 []*model.UnmatchedURLPath{},
		ServerRecommendedRules: []*model.URLTemplatizationRecommendedRule{},
		ClientRecommendedRules: []*model.URLTemplatizationRecommendedRule{},
		ExistingConfigs:        []*model.URLTemplatizationExistingConfig{},
	}
}

// GetUnmatchedUrlPaths returns client/server unmatched HTTP path counts for a workload
// from cacheDb Redis, and recommended templatization rules computed on the fly from
// those current counts (shared algorithm in common/urltemplate).
// Paths that already match an accepted URL templatization rule for the workload are
// omitted from both the path lists and the recommended-rule input.
// ExistingConfigs is the resolved UrlTemplatization config from InstrumentationConfig.
func GetUnmatchedUrlPaths(ctx context.Context, namespace, kind, name string) (*model.UnmatchedURLPaths, error) {
	cfg, err := getOdigosConfiguration(ctx)
	if err != nil {
		return nil, err
	}
	if cfg == nil || !common.UrlTemplatizationLiveTrafficLearningActive(cfg.CardinalityControl) {
		return emptyUnmatchedUrlPaths(), nil
	}

	existingConfigs, existingRules, err := loadExistingUrlTemplatizationFromIC(ctx, namespace, kind, name)
	if err != nil {
		return nil, err
	}

	rdb, err := unmatchedRedisClient()
	if err != nil {
		return nil, err
	}

	workloadPrefix := fmt.Sprintf("%s/%s/%s/", namespace, kind, name)
	server, err := loadUnmatchedPaths(ctx, rdb, unmatchedRedisKeyPrefixServer, workloadPrefix)
	if err != nil {
		return nil, err
	}
	client, err := loadUnmatchedPaths(ctx, rdb, unmatchedRedisKeyPrefixClient, workloadPrefix)
	if err != nil {
		return nil, err
	}

	server = filterPathsNotMatchingRules(server, existingRules)
	client = filterPathsNotMatchingRules(client, existingRules)

	return &model.UnmatchedURLPaths{
		Server:                 server,
		Client:                 client,
		ServerRecommendedRules: recommendedRulesFromPaths(server),
		ClientRecommendedRules: recommendedRulesFromPaths(client),
		ExistingConfigs:        existingConfigs,
	}, nil
}

// loadExistingUrlTemplatizationFromIC reads the resolved UrlTemplatization config from
// InstrumentationConfig (collector WorkloadCollectorConfig, else agent Traces).
// Returns GraphQL models and parsed PathRules used to filter already-covered paths.
func loadExistingUrlTemplatizationFromIC(ctx context.Context, namespace, kind, name string) ([]*model.URLTemplatizationExistingConfig, []urltemplate.PathRule, error) {
	if kube.DefaultClient == nil || kube.DefaultClient.OdigosClient == nil {
		return []*model.URLTemplatizationExistingConfig{}, nil, nil
	}

	icName := workload.CalculateWorkloadRuntimeObjectName(name, kind)
	ic, err := kube.DefaultClient.OdigosClient.InstrumentationConfigs(namespace).Get(ctx, icName, metav1.GetOptions{})
	if err != nil {
		if apierrors.IsNotFound(err) {
			return []*model.URLTemplatizationExistingConfig{}, nil, nil
		}
		return nil, nil, fmt.Errorf("get instrumentationconfig %s/%s: %w", namespace, icName, err)
	}

	byContainer := collectUrlTemplatizationConfigsByContainer(ic)
	configs := make([]*model.URLTemplatizationExistingConfig, 0, len(byContainer))
	seenTemplates := make(map[string]struct{})
	var templates []string

	names := make([]string, 0, len(byContainer))
	for containerName := range byContainer {
		names = append(names, containerName)
	}
	sort.Strings(names)

	for _, containerName := range names {
		cfg := byContainer[containerName]
		configs = append(configs, urlTemplatizationConfigToExistingModel(containerName, cfg))
		for _, template := range cfg.Templates {
			template = strings.TrimSpace(template)
			if template == "" {
				continue
			}
			if _, ok := seenTemplates[template]; ok {
				continue
			}
			seenTemplates[template] = struct{}{}
			templates = append(templates, template)
		}
	}

	return configs, parseTemplatizationPathRules(templates), nil
}

func collectUrlTemplatizationConfigsByContainer(ic *odigosv1.InstrumentationConfig) map[string]*actionsapi.UrlTemplatizationConfig {
	byContainer := make(map[string]*actionsapi.UrlTemplatizationConfig)

	for i := range ic.Spec.WorkloadCollectorConfig {
		c := &ic.Spec.WorkloadCollectorConfig[i]
		if c.UrlTemplatization != nil {
			byContainer[c.ContainerName] = c.UrlTemplatization
		}
	}
	for i := range ic.Spec.Containers {
		c := &ic.Spec.Containers[i]
		if _, ok := byContainer[c.ContainerName]; ok {
			continue
		}
		if c.Traces != nil && c.Traces.UrlTemplatization != nil {
			byContainer[c.ContainerName] = c.Traces.UrlTemplatization
		}
	}
	return byContainer
}

func urlTemplatizationConfigToExistingModel(containerName string, cfg *actionsapi.UrlTemplatizationConfig) *model.URLTemplatizationExistingConfig {
	templates := cfg.Templates
	if templates == nil {
		templates = []string{}
	}
	out := &model.URLTemplatizationExistingConfig{
		ContainerName: containerName,
		Templates:     templates,
	}
	if cfg.Default != nil {
		out.Default = &model.URLTemplatizationExistingDefault{
			Disabled: cfg.Default.Disabled,
		}
		if cfg.Default.SkipPolicy != nil {
			skipForNonSuccess := cfg.Default.SkipPolicy.SkipForNonSuccessCodes
			out.Default.SkipPolicy = &model.URLTemplatizationDefaultSkipPolicy{
				SkipForNonSuccessCodes: &skipForNonSuccess,
				SkipHTTPStatusCodes:    cfg.Default.SkipPolicy.SkipHttpStatusCodes,
			}
		}
	}
	return out
}

func parseTemplatizationPathRules(templates []string) []urltemplate.PathRule {
	rules := make([]urltemplate.PathRule, 0, len(templates))
	for _, template := range templates {
		rule, err := urltemplate.ParseUserInputRuleString(template, false)
		if err != nil {
			continue
		}
		rules = append(rules, rule)
	}
	return rules
}

func filterPathsNotMatchingRules(paths []*model.UnmatchedURLPath, rules []urltemplate.PathRule) []*model.UnmatchedURLPath {
	if len(paths) == 0 || len(rules) == 0 {
		return paths
	}
	out := make([]*model.UnmatchedURLPath, 0, len(paths))
	for _, item := range paths {
		if item == nil {
			continue
		}
		if pathMatchesAnyRule(item.Path, rules) {
			continue
		}
		out = append(out, item)
	}
	return out
}

func pathMatchesAnyRule(path string, rules []urltemplate.PathRule) bool {
	for _, rule := range rules {
		if rule.IsPathMatching(path) {
			return true
		}
	}
	return false
}

func recommendedRulesFromPaths(paths []*model.UnmatchedURLPath) []*model.URLTemplatizationRecommendedRule {
	counts := make(map[string]int64, len(paths))
	for _, item := range paths {
		if item == nil || item.Path == "" {
			continue
		}
		counts[item.Path] += int64(item.Count)
	}
	learned := urltemplate.FindLiveTrafficLearningRules(urltemplate.BuildPathTrieFromCounts(counts))
	return toRecommendedRuleModels(learned)
}

func toRecommendedRuleModels(rules []actionsapi.URLTemplatizationLearnedRule) []*model.URLTemplatizationRecommendedRule {
	out := make([]*model.URLTemplatizationRecommendedRule, 0, len(rules))
	for _, rule := range rules {
		segments := make([]*model.URLTemplatizationRecommendedSegment, 0, len(rule.Segments))
		for _, seg := range rule.Segments {
			examples := seg.Examples
			if examples == nil {
				examples = []string{}
			}
			segments = append(segments, &model.URLTemplatizationRecommendedSegment{
				TemplateName: seg.TemplateName,
				Certainty:    string(seg.Certainty),
				Examples:     examples,
			})
		}
		out = append(out, &model.URLTemplatizationRecommendedRule{
			Template: rule.Template,
			Reason:   rule.Reason,
			Segments: segments,
		})
	}
	return out
}

func loadUnmatchedPaths(ctx context.Context, rdb *redis.Client, kindPrefix, workloadPrefix string) ([]*model.UnmatchedURLPath, error) {
	pattern := fmt.Sprintf("%s:%s*", kindPrefix, workloadPrefix)
	var cursor uint64
	aggregated := map[string]*model.UnmatchedURLPath{}

	for {
		keys, next, err := rdb.Scan(ctx, cursor, pattern, 100).Result()
		if err != nil {
			return nil, fmt.Errorf("scan unmatched paths (%s): %w", pattern, err)
		}
		for _, key := range keys {
			containerName := containerFromUnmatchedKey(key, kindPrefix, workloadPrefix)
			fields, err := rdb.HGetAll(ctx, key).Result()
			if err != nil {
				return nil, fmt.Errorf("hgetall %s: %w", key, err)
			}
			for path, countStr := range fields {
				count, err := strconv.ParseInt(countStr, 10, 64)
				if err != nil {
					continue
				}
				aggKey := path + "\x00" + containerName
				if existing, ok := aggregated[aggKey]; ok {
					existing.Count += int(count)
					continue
				}
				cn := containerName
				aggregated[aggKey] = &model.UnmatchedURLPath{
					Path:          path,
					Count:         int(count),
					ContainerName: &cn,
				}
			}
		}
		cursor = next
		if cursor == 0 {
			break
		}
	}

	out := make([]*model.UnmatchedURLPath, 0, len(aggregated))
	for _, item := range aggregated {
		out = append(out, item)
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].Count != out[j].Count {
			return out[i].Count > out[j].Count
		}
		return out[i].Path < out[j].Path
	})
	return out, nil
}

func containerFromUnmatchedKey(key, kindPrefix, workloadPrefix string) string {
	// key = odigos:urltempl:unmatched:{server|client}:{ns}/{kind}/{name}/{container}
	prefix := kindPrefix + ":" + workloadPrefix
	if !strings.HasPrefix(key, prefix) {
		return ""
	}
	return strings.TrimPrefix(key, prefix)
}
