package rest

import (
	"github.com/nunchistudio/blacksmith/flow/source"
)

var _ source.Trigger = Alias{}
var _ source.TriggerHTTP = Alias{}
