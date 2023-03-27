


## h2non/gock

Gock(https://github.com/h2non/gock)是一个用于Golang的HTTP服务器模拟和期望库，它在很大程度上受到了NodeJs的流行和较早的同类库的启发，称为Nock


### How it mocks
1. Intercepts any HTTP outgoing request via http.DefaultTransport or custom http.Transport used by any http.Client.
2. Matches outgoing HTTP requests against a pool of defined HTTP mock expectations in FIFO declaration order.
3. If at least one mock matches, it will be used in order to compose the mock HTTP response.
4. If no mock can be matched, it will resolve the request with an error, unless real networking mode is enable, in which case a real HTTP request will be performed.


