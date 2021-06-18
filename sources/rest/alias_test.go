package rest

import (
	"github.com/nunchistudio/blacksmith/source"
)

var _ source.Trigger = Alias{}
var _ source.TriggerHTTP = Alias{}
