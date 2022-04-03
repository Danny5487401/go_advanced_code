package main

// 主体

type server interface {
	handleRequest(string, string) (int, string)
}
