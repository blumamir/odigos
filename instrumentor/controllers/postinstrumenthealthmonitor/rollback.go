package postinstrumenthealthmonitor

import (
	"fmt"
	"time"

	"github.com/odigos-io/odigos/common"
	"github.com/odigos-io/odigos/common/consts"
)

type autoRollbackConfig struct {
	disabled        bool
	graceTime       time.Duration
	stabilityWindow time.Duration
}

// getAutoRollbackConfig returns the effective auto-rollback settings from OdigosConfiguration,
// applying defaults when unset. Returns an error if a configured duration is invalid.
func getAutoRollbackConfig(o *common.OdigosConfiguration) (autoRollbackConfig, error) {
	c := autoRollbackConfig{
		graceTime:       consts.DefaultAutoRollbackGraceTime,
		stabilityWindow: consts.DefaultAutoRollbackStabilityWindow,
	}
	if o == nil {
		return c, nil
	}

	c.disabled = o.RollbackDisabled != nil && *o.RollbackDisabled

	if o.RollbackGraceTime != "" {
		parsed, err := time.ParseDuration(o.RollbackGraceTime)
		if err != nil {
			return autoRollbackConfig{}, fmt.Errorf("invalid RollbackGraceTime %q: %w", o.RollbackGraceTime, err)
		}
		c.graceTime = parsed
	}

	if o.RollbackStabilityWindow != "" {
		parsed, err := time.ParseDuration(o.RollbackStabilityWindow)
		if err != nil {
			return autoRollbackConfig{}, fmt.Errorf("invalid RollbackStabilityWindow %q: %w", o.RollbackStabilityWindow, err)
		}
		c.stabilityWindow = parsed
	}

	return c, nil
}
