package rest

import (
	"github.com/nunchistudio/blacksmith/flow/source"
)

var _ source.Trigger = Page{}
var _ source.TriggerHTTP = Page{}
