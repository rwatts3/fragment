package rest

import (
	"github.com/nunchistudio/blacksmith/source"
)

var _ source.Trigger = Group{}
var _ source.TriggerHTTP = Group{}
