package services

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/hashicorp/go-version"
	"github.com/odigos-io/odigos/api/k8sconsts"
	"github.com/odigos-io/odigos/frontend/graph/model"
	"github.com/odigos-io/odigos/frontend/kube"
	"github.com/odigos-io/odigos/k8sutils/pkg/env"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// GetGoEnterpriseOffsets returns parsed go_offset_results.json from the
// odigos-go-offsets ConfigMap in the Odigos namespace.
func GetGoEnterpriseOffsets(ctx context.Context) (*model.GoEnterpriseOffsets, error) {
	content, err := getGoEnterpriseOffsetsRaw(ctx)
	if err != nil {
		return nil, err
	}

	parsed, err := parseGoEnterpriseOffsetsContent(content)
	if err != nil {
		return nil, err
	}

	return goEnterpriseOffsetsToModel(parsed), nil
}

func getGoEnterpriseOffsetsRaw(ctx context.Context) (string, error) {
	ns := env.GetCurrentNamespace()

	var cm corev1.ConfigMap
	err := kube.CacheClient.Get(ctx, client.ObjectKey{
		Namespace: ns,
		Name:      k8sconsts.GoOffsetsConfigMap,
	}, &cm)
	if err != nil {
		return "", fmt.Errorf("failed to get ConfigMap %s/%s: %w", ns, k8sconsts.GoOffsetsConfigMap, err)
	}

	content, ok := cm.Data[k8sconsts.GoOffsetsFileName]
	if !ok {
		return "", fmt.Errorf("key %q not found in ConfigMap %s/%s", k8sconsts.GoOffsetsFileName, ns, k8sconsts.GoOffsetsConfigMap)
	}

	return content, nil
}

func parseGoEnterpriseOffsetsContent(content string) (*versionedModules, error) {
	trimmed := strings.TrimSpace(content)
	if trimmed == "" {
		return &versionedModules{Mods: []*jsonModule{}}, nil
	}

	// ConfigMap data is a JSON-encoded string (see odigos pro update-offsets / Helm).
	var inner string
	if err := json.Unmarshal([]byte(trimmed), &inner); err != nil {
		return nil, fmt.Errorf("invalid go enterprise offsets JSON: %w", err)
	}
	inner = strings.TrimSpace(inner)
	if inner == "" {
		return &versionedModules{Mods: []*jsonModule{}}, nil
	}
	inner = stripGoEnterpriseOffsetsSignature(inner)

	var parsed versionedModules
	if err := json.Unmarshal([]byte(inner), &parsed); err != nil {
		return nil, fmt.Errorf("invalid go enterprise offsets JSON: %w", err)
	}

	if parsed.Mods == nil {
		parsed.Mods = []*jsonModule{}
	}

	return &parsed, nil
}

const goEnterpriseOffsetsSignatureDelimiter = "---SIGNATURE---"

func stripGoEnterpriseOffsetsSignature(inner string) string {
	parts := strings.SplitN(inner, goEnterpriseOffsetsSignatureDelimiter, 2)
	return strings.TrimSpace(parts[0])
}

func goEnterpriseOffsetsToModel(parsed *versionedModules) *model.GoEnterpriseOffsets {
	mods := make([]*model.GoEnterpriseOffsetModule, 0, len(parsed.Mods))
	for _, mod := range parsed.Mods {
		if mod == nil {
			continue
		}
		byMinorVersions, minVersion, maxVersion := moduleVersions(mod)
		minorVersionsModel := make([]*model.GoEnterpriseOffsetMinorVersionEnumeration, 0, len(byMinorVersions))
		minorVersions := make([]string, 0, len(byMinorVersions))
		for majorMinor := range byMinorVersions {
			minorVersions = append(minorVersions, majorMinor)
		}
		sort.Slice(minorVersions, func(i, j int) bool {
			return compareVersions(minorVersions[i], minorVersions[j])
		})

		for _, minorVersion := range minorVersions {
			versionsForThisMinor := byMinorVersions[minorVersion]
			sort.Slice(versionsForThisMinor, func(i, j int) bool {
				return compareVersions(versionsForThisMinor[i], versionsForThisMinor[j])
			})
			minorVersionsModel = append(minorVersionsModel, &model.GoEnterpriseOffsetMinorVersionEnumeration{
				MinorVersion: minorVersion,
				Versions:     versionsForThisMinor,
			})
		}
		mods = append(mods, &model.GoEnterpriseOffsetModule{
			Module:        mod.Module,
			MinVersion:    minVersion,
			MaxVersion:    maxVersion,
			MinorVersions: minorVersionsModel,
		})
	}

	timestamp := ""
	if !parsed.Timestamp.IsZero() {
		timestamp = parsed.Timestamp.UTC().Format(time.RFC3339Nano)
	}

	return &model.GoEnterpriseOffsets{
		Timestamp: timestamp,
		Mods:      mods,
	}
}

func moduleVersions(mod *jsonModule) (byMinorVersions map[string][]string, minVersion, maxVersion string) {
	seen := make(map[string]struct{})

	var min, max *version.Version
	var minorToVersions map[string][]string

	consider := func(raw string) {
		if _, ok := seen[raw]; ok {
			return
		}
		v, err := version.NewVersion(raw)
		if err != nil {
			return
		}
		seen[raw] = struct{}{}
		if min == nil || v.LessThan(min) {
			min = v
		}
		if max == nil || v.GreaterThan(max) {
			max = v
		}
		seg := v.Segments()
		majorMinor := fmt.Sprintf("%d.%d", seg[0], seg[1])
		if minorToVersions == nil {
			minorToVersions = make(map[string][]string)
		}
		minorToVersions[majorMinor] = append(minorToVersions[majorMinor], raw)
	}

	for _, pkg := range mod.Packages {
		if pkg == nil {
			continue
		}
		for _, s := range pkg.Structs {
			if s == nil {
				continue
			}
			for _, f := range s.Fields {
				if f == nil {
					continue
				}
				for _, o := range f.Offsets {
					if o == nil {
						continue
					}
					for _, v := range o.Versions {
						consider(v)
					}
				}
			}
		}
	}

	if min == nil || max == nil {
		return minorToVersions, "", ""
	}
	return minorToVersions, min.String(), max.String()
}

func compareVersions(a, b string) bool {
	va, errA := version.NewVersion(a)
	vb, errB := version.NewVersion(b)
	if errA != nil || errB != nil {
		return a < b
	}
	return va.LessThan(vb)
}

// UpdateGoEnterpriseOffsets writes the provided offsets file content to the
// odigos-go-offsets ConfigMap (same encoding as `odigos pro update-offsets`).
func UpdateGoEnterpriseOffsets(ctx context.Context, content string) error {
	if GetTier(ctx) == model.TierCommunity {
		return fmt.Errorf("custom offsets support is only available in Odigos pro tier")
	}

	return writeGoEnterpriseOffsetsConfigMap(ctx, []byte(content))
}

func writeGoEnterpriseOffsetsConfigMap(ctx context.Context, data []byte) error {
	ns := env.GetCurrentNamespace()

	cm, err := kube.DefaultClient.CoreV1().ConfigMaps(ns).Get(ctx, k8sconsts.GoOffsetsConfigMap, metav1.GetOptions{})
	if err != nil {
		return fmt.Errorf("failed to get ConfigMap %s/%s: %w", ns, k8sconsts.GoOffsetsConfigMap, err)
	}

	if cm.Data == nil {
		cm.Data = make(map[string]string)
	}

	if len(data) == 0 {
		cm.Data[k8sconsts.GoOffsetsFileName] = ""
	} else {
		encoded, err := json.Marshal(string(data))
		if err != nil {
			return fmt.Errorf("failed to encode go enterprise offsets: %w", err)
		}
		cm.Data[k8sconsts.GoOffsetsFileName] = string(encoded)
	}

	_, err = kube.DefaultClient.CoreV1().ConfigMaps(ns).Update(ctx, cm, metav1.UpdateOptions{})
	if err != nil {
		return fmt.Errorf("failed to update ConfigMap %s/%s: %w", ns, k8sconsts.GoOffsetsConfigMap, err)
	}

	return nil
}
