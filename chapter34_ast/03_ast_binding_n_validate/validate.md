

## 需求

```go
func createAccount(w http.ResponseWriter, req *http.Request) {
    req.ParseForm()
    var createRequest =  createAccountRequest{}
    if createRequest.Age, ok = req.Form["age"];!ok {
        return someErrorResponse
    }
    if createRequest.Name, ok = req.Form["name"]; !ok {
        return someErrorResponse
    }
    ......
}
```
在每一个 web 开发的入口层 api，一般会做统一的参数绑定和校验，实际上这些工作大多数情况下都是让人比较烦躁的重复劳动。

生活真是不美好，如果这个 Request 有 50 个字段，我真的有点想跳楼。

### 问题

```go
type req struct {
    Age int
}
```
那么当 Age == 0 的时候，你就再也没有办法判断这个 0 到底是 go 的类型默认值还是调用方压根儿没有传这个参数了。
解决方法倒也是有，你可以把 Age 变成 *int，然后在每次取 req.Age 的时候都去解引用