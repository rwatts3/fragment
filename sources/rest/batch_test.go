package rest

import (
	"github.com/nunchistudio/blacksmith/flow/source"
)

var _ source.Trigger = Batch{}
var _ source.TriggerHTTP = Batch{}
