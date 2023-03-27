package _2_gock

import (
	"github.com/stretchr/testify/assert"
	"gopkg.in/h2non/gock.v1"
	"io/ioutil"
	"net/http"
	"testing"
)

func TestSimple(t *testing.T) {
	defer gock.Off()

	gock.New("https://foo.com").
		Get("/bar").
		Reply(200).
		JSON(map[string]string{"foo": "bar"})

	res, err := http.Get("https://foo.com/bar")
	assert.Nil(t, err)
	assert.Equal(t, 200, res.StatusCode)

	body, _ := ioutil.ReadAll(res.Body)
	assert.Equal(t, `{"foo":"bar"}`, string(body)[:13])

	// Verify that we don't have pending mocks
	assert.True(t, gock.IsDone())
}
