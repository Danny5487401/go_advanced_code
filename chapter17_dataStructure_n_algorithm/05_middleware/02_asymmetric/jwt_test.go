package _2_asymmetric

import (
	"testing"

	"fmt"

	"github.com/Danny5487401/go_advanced_code/chapter17_dataStructure_n_algorithm/05_middleware/models"

	. "github.com/smartystreets/goconvey/convey"
)

func TestName(t *testing.T) {
	var (
		privatePath = "chapter17_dataStructure_n_algorithm/06_certificate/pem_file/private.pem"
		pubPath     = "chapter17_dataStructure_n_algorithm/06_certificate/pem_file/public.pem"
		userInfo    = &models.User{
			UserID:   123,
			NickName: "danny",
		}
	)

	InitJWT(privatePath, pubPath)

	Convey("token 解析", t, func() {
		token, err := GenerateJWT(userInfo)
		So(err, ShouldBeNil)
		claim, err := ParseToken(token)
		fmt.Printf("%+v", claim.User)
		So(err, ShouldBeNil)
		So(claim.User, ShouldResemble, userInfo)
	})

}
