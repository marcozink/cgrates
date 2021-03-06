/*
Real-time Online/Offline Charging System (OCS) for Telecom & ISP environments
Copyright (C) ITsysCOM GmbH

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program.  If not, see <http://www.gnu.org/licenses/>
*/

package engine

import (
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/cgrates/cgrates/utils"
)

func TestSetFailedPostCacheTTL(t *testing.T) {
	var1 := failedPostCache
	SetFailedPostCacheTTL(time.Duration(50 * time.Millisecond))
	var2 := failedPostCache
	if reflect.DeepEqual(var1, var2) {
		t.Error("Expecting to be different")
	}
}

func TestAddFailedPost(t *testing.T) {
	SetFailedPostCacheTTL(time.Duration(5 * time.Second))
	addFailedPost("path1", "format1", "module1", "1")
	x, ok := failedPostCache.Get(utils.ConcatenatedKey("path1", "format1", "module1"))
	if !ok {
		t.Error("Error reading from cache")
	}
	if x == nil {
		t.Error("Received an empty element")
	}

	failedPost, canCast := x.(*ExportEvents)
	if !canCast {
		t.Error("Error when casting")
	}
	eOut := &ExportEvents{
		Path:   "path1",
		Format: "format1",
		module: "module1",
		Events: []interface{}{"1"},
	}
	if !reflect.DeepEqual(eOut, failedPost) {
		t.Errorf("Expecting: %+v, received: %+v", utils.ToJSON(eOut), utils.ToJSON(failedPost))
	}
	addFailedPost("path1", "format1", "module1", "2")
	addFailedPost("path2", "format2", "module2", "3")
	x, ok = failedPostCache.Get(utils.ConcatenatedKey("path1", "format1", "module1"))
	if !ok {
		t.Error("Error reading from cache")
	}
	if x == nil {
		t.Error("Received an empty element")
	}
	failedPost, canCast = x.(*ExportEvents)
	if !canCast {
		t.Error("Error when casting")
	}
	eOut = &ExportEvents{
		Path:   "path1",
		Format: "format1",
		module: "module1",
		Events: []interface{}{"1", "2"},
	}
	if !reflect.DeepEqual(eOut, failedPost) {
		t.Errorf("Expecting: %+v, received: %+v", utils.ToJSON(eOut), utils.ToJSON(failedPost))
	}
	x, ok = failedPostCache.Get(utils.ConcatenatedKey("path2", "format2", "module2"))
	if !ok {
		t.Error("Error reading from cache")
	}
	if x == nil {
		t.Error("Received an empty element")
	}
	failedPost, canCast = x.(*ExportEvents)
	if !canCast {
		t.Error("Error when casting")
	}
	eOut = &ExportEvents{
		Path:   "path2",
		Format: "format2",
		module: "module2",
		Events: []interface{}{"3"},
	}
	if !reflect.DeepEqual(eOut, failedPost) {
		t.Errorf("Expecting: %+v, received: %+v", utils.ToJSON(eOut), utils.ToJSON(failedPost))
	}
}

func TestFileName(t *testing.T) {
	exportEvent := &ExportEvents{}
	rcv := exportEvent.FileName()
	if rcv[0] != '|' {
		t.Errorf("Expecting: '|', received: %+v", rcv[0])
	} else if rcv[8:] != ".gob" {
		t.Errorf("Expecting: '.gob', received: %+v", rcv[8:])
	}
	exportEvent = &ExportEvents{
		module: "module",
	}
	rcv = exportEvent.FileName()
	if rcv[:7] != "module|" {
		t.Errorf("Expecting: 'module|', received: %+v", rcv[:7])
	} else if rcv[14:] != ".gob" {
		fmt.Println(rcv)
		t.Errorf("Expecting: '.gob', received: %+v", rcv[14:])
	}

}

func TestSetModule(t *testing.T) {
	exportEvent := &ExportEvents{}
	eOut := &ExportEvents{
		module: "module",
	}
	exportEvent.SetModule("module")
	if !reflect.DeepEqual(eOut, exportEvent) {
		t.Errorf("Expecting: %+v, received: %+v", eOut, exportEvent)
	}
}

func TestAddEvent(t *testing.T) {
	exportEvent := &ExportEvents{}
	eOut := &ExportEvents{Events: []interface{}{"event1"}}
	exportEvent.AddEvent("event1")
	if !reflect.DeepEqual(eOut, exportEvent) {
		t.Errorf("Expecting: %+v, received: %+v", eOut, exportEvent)
	}
	exportEvent = &ExportEvents{}
	eOut = &ExportEvents{Events: []interface{}{"event1", "event2", "event3"}}
	exportEvent.AddEvent("event1")
	exportEvent.AddEvent("event2")
	exportEvent.AddEvent("event3")
	if !reflect.DeepEqual(eOut, exportEvent) {
		t.Errorf("Expecting: %+v, received: %+v", utils.ToJSON(eOut), utils.ToJSON(exportEvent))
	}
}
