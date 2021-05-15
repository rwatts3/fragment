package rest

import (
	"github.com/nunchistudio/blacksmith/flow/source"
)

var _ source.Trigger = Identify{}
var _ source.TriggerHTTP = Identify{}
