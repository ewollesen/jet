//
// Code generated by go-jet DO NOT EDIT.
// Generated at Thursday, 26-Sep-19 12:02:13 CEST
//
// WARNING: Changes to this file may cause incorrect behavior
// and will be lost if the code is regenerated
//

package model

import (
	"time"
)

type Category struct {
	CategoryID int32 `sql:"primary_key"`
	Name       string
	LastUpdate time.Time
}
