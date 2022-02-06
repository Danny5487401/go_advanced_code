package product_srv

import "time"

/********商品服务**********/

type Product struct {
	Title string
	Price uint32
}

func GetProductList() ([]Product, error) {
	time.Sleep(400 * time.Millisecond)
	var list []Product
	list = append(list, Product{
		Title: "SHib",
		Price: 10,
	})
	return list, nil
}
