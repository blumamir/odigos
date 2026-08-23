package status

import (
	"github.com/odigos-io/odigos/api/odigos/v1alpha1"
	"github.com/odigos-io/odigos/frontend/graph/model"
	"github.com/odigos-io/odigos/status"
	postInstrumentHealthMonitorStatus "github.com/odigos-io/odigos/status/instrumentationconfig/generated"
)

const (
	RollbackStatus = "Rollback"
)

func CalculateAutoRollbackStatus(ic *v1alpha1.InstrumentationConfig) *model.DesiredConditionStatus {
	if ic == nil {
		return nil
	}

	for _, c := range ic.Status.Conditions {
		if c.Type != postInstrumentHealthMonitorStatus.PostInstrumentHealthMonitorType {
			continue
		}
		r, ok := postInstrumentHealthMonitorStatus.PostInstrumentHealthMonitorReasonByName(c.Reason)
		if !ok {
			return &model.DesiredConditionStatus{
				Name:       RollbackStatus,
				Status:     model.DesiredStateProgressUnknown,
				ReasonEnum: &c.Reason,
				Message:    c.Message,
			}
		}

		return postInstrumentHealthMonitorReasonToAutoRollbackStatus(r, c.Message)
	}

	return postInstrumentHealthMonitorReasonToAutoRollbackStatus(
		postInstrumentHealthMonitorStatus.PostInstrumentHealthMonitorNotApplicable,
		"",
	)
}

func postInstrumentHealthMonitorReasonToAutoRollbackStatus(r status.Reason, message string) *model.DesiredConditionStatus {
	if message == "" {
		message = r.Message
	}

	actionItems := make([]*model.DesiredConditionActionItem, 0, len(r.ActionItems))
	for _, actionItem := range r.ActionItems {
		actionItems = append(actionItems, &model.DesiredConditionActionItem{
			Type:       model.DesiredConditionActionItemType(actionItem.Type),
			ButtonText: actionItem.ButtonText,
		})
	}

	return &model.DesiredConditionStatus{
		Name:        RollbackStatus,
		Status:      model.DesiredStateProgress(r.OdigosSeverity),
		ReasonEnum:  &r.Title,
		Message:     message,
		ActionItems: actionItems,
	}
}
