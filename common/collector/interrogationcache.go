package collector

import "go.opentelemetry.io/collector/pdata/pcommon"

// InterrogationCacheExtension is implemented by the odigos_interrogation_cache collector
// extension. The interrogation profiles exporter stores probe/uprobe sample stacks keyed
// by span link; the traces exporter reads them for the bounding join.
type InterrogationCacheExtension interface {
	// RecordSample stores the resolved stack frames for a profile sample linked to the
	// given trace and span. Frames are ordered leaf-first (same order as the profile
	// location indices). Callers must pass non-empty trace and span IDs. Multiple
	// samples for the same key are retained as separate stacks.
	RecordSample(traceID pcommon.TraceID, spanID pcommon.SpanID, frames []string)

	// GetSamples returns the stack frames recorded for the given trace and span.
	// Each inner slice is one sample's frames (leaf-first). ok is false when no
	// non-expired entry exists. The returned slices are copies and safe to mutate.
	GetSamples(traceID pcommon.TraceID, spanID pcommon.SpanID) (samples [][]string, ok bool)
}
