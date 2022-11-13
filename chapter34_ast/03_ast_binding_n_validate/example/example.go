package example

type CollectRequest struct {
	Star int    `form:"star" validation:"gte=1,lte=5" doc:"formData"`
	Name string `form:"name" validation:"length>0"`
}
