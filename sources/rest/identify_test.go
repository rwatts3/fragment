package rest

import (
	"github.com/nunchistudio/blacksmith/source"
)

var _ source.Trigger = Identify{}
var _ source.TriggerHTTP = Identify{}
