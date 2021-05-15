package rest

import (
	"github.com/nunchistudio/blacksmith/flow/source"
)

var _ source.Trigger = Track{}
var _ source.TriggerHTTP = Track{}
