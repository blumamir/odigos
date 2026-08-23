package rollout

import (
	"github.com/odigos-io/odigos/common"
)

type RolloutOptions struct {
	MaxConcurrentRollouts int
}

func getRolloutOptions(conf *common.OdigosConfiguration) (isAutomaticRolloutDisabled bool, rolloutOptions RolloutOptions) {
	isAutomaticRolloutDisabled = conf.Rollout != nil && conf.Rollout.AutomaticRolloutDisabled != nil && *conf.Rollout.AutomaticRolloutDisabled

	maxConcurrentRollouts := NoConcurrencyLimiting
	if conf.Rollout != nil && conf.Rollout.MaxConcurrentRollouts > 0 {
		maxConcurrentRollouts = conf.Rollout.MaxConcurrentRollouts
	}

	return isAutomaticRolloutDisabled, RolloutOptions{
		MaxConcurrentRollouts: maxConcurrentRollouts,
	}
}
