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
	"github.com/odigos-io/odigos/common"
	"github.com/odigos-io/odigos/frontend/graph/model"
	"github.com/odigos-io/odigos/k8sutils/pkg/env"
	"github.com/redis/go-redis/v9"
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

// GetUnmatchedUrlPaths returns client/server unmatched HTTP path counts for a workload
// from cacheDb Redis (populated when cardinalityControl.urlTemplatization.autoComputeRules is enabled).
func GetUnmatchedUrlPaths(ctx context.Context, namespace, kind, name string) (*model.UnmatchedURLPaths, error) {
	cfg, err := getOdigosConfiguration(ctx)
	if err != nil {
		return nil, err
	}
	if cfg == nil || !common.UrlTemplatizationAutoComputeRulesActive(cfg.CardinalityControl) {
		return &model.UnmatchedURLPaths{Server: []*model.UnmatchedURLPath{}, Client: []*model.UnmatchedURLPath{}}, nil
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
	return &model.UnmatchedURLPaths{Server: server, Client: client}, nil
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
