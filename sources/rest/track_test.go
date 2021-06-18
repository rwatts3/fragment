package rest

import (
	"github.com/nunchistudio/blacksmith/source"
)

var _ source.Trigger = Track{}
var _ source.TriggerHTTP = Track{}
