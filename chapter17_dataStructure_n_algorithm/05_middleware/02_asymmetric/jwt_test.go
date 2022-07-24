package _2_asymmetric

import (
	"testing"

	"github.com/Danny5487401/go_advanced_code/chapter17_dataStructure_n_algorithm/05_middleware/models"

	. "github.com/smartystreets/goconvey/convey"
)

func TestName(t *testing.T) {
	userInfo := &models.User{
		UserID:   123,
		NickName: "danny",
	}

	Convey("token 解析", t, func() {
		token, err := GenerateJWT(userInfo)
		So(err, ShouldBeNil)
		claim, err := ParseToken(token)
		So(err, ShouldBeNil)
		So(claim.User, ShouldResemble, userInfo)
	})

}
