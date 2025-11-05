package memoryinfoextension

import (
	"context"
	"os"
	"strconv"
	"strings"
	"time"

	rtml "github.com/odigos-io/go-rtml"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/extension"
	"go.uber.org/zap"
)

type MemoryInfo struct {
	cfg          *Config
	logger       *zap.Logger
	loopInterval time.Duration

	lastReportTime time.Time

	stopCh chan struct{}
}

var _ extension.Extension = (*MemoryInfo)(nil)

func NewMemoryInfo(cfg *Config, logger *zap.Logger) (*MemoryInfo, error) {

	loopInterval, err := time.ParseDuration(cfg.LoopInterval)
	if err != nil {
		return nil, err
	}

	return &MemoryInfo{
		cfg:          cfg,
		logger:       logger,
		loopInterval: loopInterval,
		stopCh:       make(chan struct{}),
	}, nil
}

func (m *MemoryInfo) Start(ctx context.Context, _ component.Host) error {
	if !m.cfg.Enabled {
		m.logger.Info("memory info extension is disabled")
		return nil
	}

	go m.run(ctx)

	return nil
}

func (m *MemoryInfo) Shutdown(ctx context.Context) error {
	close(m.stopCh)
	return nil
}

func (m *MemoryInfo) run(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case <-m.stopCh:
			m.logger.Info("memory info extension is stopped")
			return
		case <-time.After(m.loopInterval):
			m.emitMemoryInfo()
		}
	}
}

func bytesToMiB(bytes uint64) float64 {
	return float64(bytes) / 1024 / 1024
}

// cgroupMemoryInfo contains cgroup memory statistics
type cgroupMemoryInfo struct {
	limitBytes uint64
	usageBytes uint64
	hasLimit   bool
	hasUsage   bool
}

// readCgroupMemoryInfo reads cgroup memory information from both v1 and v2
func readCgroupMemoryInfo() cgroupMemoryInfo {
	info := cgroupMemoryInfo{}

	// Try cgroup v2 first (newer systems)
	if limit, err := readCgroupV2MemoryLimit(); err == nil {
		info.limitBytes = limit
		info.hasLimit = true
	} else if limit, err := readCgroupV1MemoryLimit(); err == nil {
		// Fallback to cgroup v1
		info.limitBytes = limit
		info.hasLimit = true
	}

	// Try cgroup v2 for usage
	if usage, err := readCgroupV2MemoryUsage(); err == nil {
		info.usageBytes = usage
		info.hasUsage = true
	} else if usage, err := readCgroupV1MemoryUsage(); err == nil {
		// Fallback to cgroup v1
		info.usageBytes = usage
		info.hasUsage = true
	}

	return info
}

// readCgroupV2MemoryLimit reads memory limit from cgroup v2
func readCgroupV2MemoryLimit() (uint64, error) {
	data, err := os.ReadFile("/sys/fs/cgroup/memory.max")
	if err != nil {
		return 0, err
	}

	limitStr := strings.TrimSpace(string(data))
	// "max" means no limit
	if limitStr == "max" {
		return 0, nil
	}

	return strconv.ParseUint(limitStr, 10, 64)
}

// readCgroupV2MemoryUsage reads current memory usage from cgroup v2
func readCgroupV2MemoryUsage() (uint64, error) {
	data, err := os.ReadFile("/sys/fs/cgroup/memory.current")
	if err != nil {
		return 0, err
	}

	return strconv.ParseUint(strings.TrimSpace(string(data)), 10, 64)
}

// readCgroupV1MemoryLimit reads memory limit from cgroup v1
func readCgroupV1MemoryLimit() (uint64, error) {
	data, err := os.ReadFile("/sys/fs/cgroup/memory/memory.limit_in_bytes")
	if err != nil {
		return 0, err
	}

	limit, err := strconv.ParseUint(strings.TrimSpace(string(data)), 10, 64)
	if err != nil {
		return 0, err
	}

	// In cgroup v1, a very large number (like 9223372036854771712) means no limit
	// This is typically 0x7FFFFFFFFFFFF000 or similar
	if limit > (1 << 62) {
		return 0, nil
	}

	return limit, nil
}

// readCgroupV1MemoryUsage reads current memory usage from cgroup v1
func readCgroupV1MemoryUsage() (uint64, error) {
	data, err := os.ReadFile("/sys/fs/cgroup/memory/memory.usage_in_bytes")
	if err != nil {
		return 0, err
	}

	return strconv.ParseUint(strings.TrimSpace(string(data)), 10, 64)
}

func (m *MemoryInfo) emitMemoryInfo() {

	// calculate time since last report
	timeSinceLastReport := time.Since(m.lastReportTime)
	// if there is no CPU throttle, the last report delay should be 0.
	// e.g. the loop should run after the loop interval and it does without any delay.
	// any delay here means that events takes time to be processed by the collector.
	lastReportDelay := timeSinceLastReport - m.loopInterval

	delayIsProblematic := lastReportDelay > 100*time.Millisecond
	if delayIsProblematic {
		m.logger.Warn("memory diagnostic reported with noticeable delay, CPU might be throttled", zap.String("loop interval", m.loopInterval.String()), zap.String("lastReportDelay", lastReportDelay.String()), zap.String("timeSinceLastReport", timeSinceLastReport.String()))
	}
	m.lastReportTime = time.Now()

	isMemoryLimitReached := rtml.IsMemLimitReached()
	m.LogFullMemoryInfo("periodic interval", isMemoryLimitReached)
}

// log the go runtime memory info and cgroup memory info (which is expensive to read)
func (m *MemoryInfo) LogFullMemoryInfo(scope string, memoryLimitReached bool) {
	cgroupInfo := readCgroupMemoryInfo()
	stats := rtml.GetMemLimitRelatedStats()

	logFields := []zap.Field{
		zap.String("scope", scope),
		zap.Bool("isMemoryLimitReached", memoryLimitReached),
		zap.Float64("memoryLimitMiB", bytesToMiB(stats.MemoryLimit)),
		zap.Float64("heapGoalMiB", bytesToMiB(stats.HeapGoal)),
		zap.Float64("heapLiveMiB", bytesToMiB(stats.HeapLive)),
		zap.Float64("mappedReadyMiB", bytesToMiB(stats.MappedReady)),
		zap.Float64("heapFreeMiB", bytesToMiB(stats.HeapFree)),
		zap.Float64("totalAllocMiB", bytesToMiB(stats.TotalAlloc)),
		zap.Float64("totalFreeMiB", bytesToMiB(stats.TotalFree)),
	}

	// Add cgroup memory information if available
	if cgroupInfo.hasLimit {
		if cgroupInfo.limitBytes == 0 {
			logFields = append(logFields, zap.String("cgroupMemoryLimitMiB", "unlimited"))
		} else {
			logFields = append(logFields, zap.Float64("cgroupMemoryLimitMiB", bytesToMiB(cgroupInfo.limitBytes)))
		}
	}

	if cgroupInfo.hasUsage {
		logFields = append(logFields, zap.Float64("cgroupMemoryUsageMiB", bytesToMiB(cgroupInfo.usageBytes)))

		// Calculate usage percentage if we have both limit and usage
		if cgroupInfo.hasLimit && cgroupInfo.limitBytes > 0 {
			usagePercent := float64(cgroupInfo.usageBytes) / float64(cgroupInfo.limitBytes) * 100
			logFields = append(logFields, zap.Float64("cgroupMemoryUsagePercent", usagePercent))
		}
	}

	m.logger.Info("full memory diagnostic info", logFields...)
}

// log only go runtime memory info. avoid the cgroup info which is expensive to read.
func (m *MemoryInfo) LogCheapMemoryInfo(scope string, memoryLimitReached bool) {
	stats := rtml.GetMemLimitRelatedStats()

	m.logger.Info("basic memory diagnostic info",
		zap.String("scope", scope),
		zap.Bool("isMemoryLimitReached", memoryLimitReached),
		zap.Float64("heapGoalMiB", bytesToMiB(stats.HeapGoal)),
		zap.Float64("heapLiveMiB", bytesToMiB(stats.HeapLive)),
		zap.Float64("mappedReadyMiB", bytesToMiB(stats.MappedReady)),
		zap.Float64("heapFreeMiB", bytesToMiB(stats.HeapFree)),
		zap.Float64("totalAllocMiB", bytesToMiB(stats.TotalAlloc)),
		zap.Float64("totalFreeMiB", bytesToMiB(stats.TotalFree)),
	)
}
