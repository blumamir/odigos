package memoryinfoextension

import (
	"context"
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

func (m *MemoryInfo) emitMemoryInfo() {

	// calculate time since last report
	timeSinceLastReport := time.Since(m.lastReportTime)
	// if there is no CPU throttle, the last report delay should be 0.
	// e.g. the loop should run after the loop interval and it does without any delay.
	// any delay here means that events takes time to be processed by the collector.
	lastReportDelay := timeSinceLastReport - m.loopInterval

	delayIsProblematic := lastReportDelay > 500*time.Millisecond
	if delayIsProblematic {
		m.logger.Warn("memory diagnostic reported with noticeable delay, CPU might be throttled", zap.String("loop interval", m.loopInterval.String()), zap.String("lastReportDelay", lastReportDelay.String()), zap.String("timeSinceLastReport", timeSinceLastReport.String()))
	}

	stats := rtml.GetMemLimitRelatedStats()
	isMemoryLimitReached := rtml.IsMemLimitReached()

	m.logger.Info("memory diagnostic info",
		zap.Bool("isMemoryLimitReached", isMemoryLimitReached),
		zap.Float64("memoryLimitMiB", bytesToMiB(stats.MemoryLimit)),
		zap.Float64("heapGoalMiB", bytesToMiB(stats.HeapGoal)),
		zap.Float64("heapLiveMiB", bytesToMiB(stats.HeapLive)),
		zap.Float64("mappedReadyMiB", bytesToMiB(stats.MappedReady)),
		zap.Float64("heapFreeMiB", bytesToMiB(stats.HeapFree)),
		zap.Float64("totalAllocMiB", bytesToMiB(stats.TotalAlloc)),
		zap.Float64("totalFreeMiB", bytesToMiB(stats.TotalFree)),
	)

	m.lastReportTime = time.Now()
}
