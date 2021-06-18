package rest

import (
	"github.com/nunchistudio/blacksmith/source"
)

var _ source.Trigger = Page{}
var _ source.TriggerHTTP = Page{}
