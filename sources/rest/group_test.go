package rest

import (
	"github.com/nunchistudio/blacksmith/flow/source"
)

var _ source.Trigger = Group{}
var _ source.TriggerHTTP = Group{}
