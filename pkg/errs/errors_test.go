package errs

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestGenError_AddMsg(t *testing.T) {
	Convey("TestGenError_AddMsg", t, func() {
		msg := "some error message"
		ge := GenError{}
		So(ge.AddMsg(msg), ShouldResemble, &GenError{Messages: []string{msg}})
	})
}

func TestGenError_Error(t *testing.T) {
	Convey("TestGenError_Error", t, func() {
		msg1 := "message one"
		msg2 := "message two"
		ge := GenError{Messages: []string{msg1, msg2}}
		So(ge.Error(), ShouldEqual, msg1+"->"+msg2)
	})
}
