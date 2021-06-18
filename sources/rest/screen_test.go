package rest

import (
	"github.com/nunchistudio/blacksmith/source"
)

var _ source.Trigger = Screen{}
var _ source.TriggerHTTP = Screen{}
