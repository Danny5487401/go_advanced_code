package main

import (
	. "github.com/agiledragon/gomonkey/v2"
	. "github.com/smartystreets/goconvey/convey"
	"go_advanced_code/chapter28_test/01_gomonkey/fake"
	"testing"
)

func TestApplyInterfaceReused(t *testing.T) {
	e := &fake.Etcd{}

	Convey("TestApplyInterface", t, func() {
		patches := ApplyFunc(fake.NewDb, func(_ string) fake.Db {
			return e
		})
		defer patches.Reset()
		db := fake.NewDb("mysql")

		Convey("TestApplyInterface", func() {
			info := "hello interface"
			patches.ApplyMethod(e, "Retrieve",
				func(_ *fake.Etcd, _ string) (string, error) {
					return info, nil
				})
			output, err := db.Retrieve("")
			So(err, ShouldEqual, nil)
			So(output, ShouldEqual, info)
		})
	})
}
