package _1_symmetric

import (
	"github.com/Danny5487401/go_advanced_code/chapter12_net/08_middleware/models"
	"testing"

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
