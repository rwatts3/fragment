package rest

import (
	"github.com/nunchistudio/blacksmith/source"
)

var _ source.Trigger = Batch{}
var _ source.TriggerHTTP = Batch{}
