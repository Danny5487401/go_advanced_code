package _1_symmetric

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
		jwtInfo := NewJWT([]byte("231321dsad"))
		token, err := jwtInfo.CreateToken(userInfo)
		So(err, ShouldBeNil)
		claim, err := jwtInfo.ParseToken(token)
		So(err, ShouldBeNil)
		So(claim.User, ShouldResemble, userInfo)
	})

}
