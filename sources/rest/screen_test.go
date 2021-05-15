package rest

import (
	"github.com/nunchistudio/blacksmith/flow/source"
)

var _ source.Trigger = Screen{}
var _ source.TriggerHTTP = Screen{}
