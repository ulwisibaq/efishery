package main

import (
	cmdauth "github.com/ulwisibaq/efishery/auth/cmd"
	cmdcomm "github.com/ulwisibaq/efishery/commodity/cmd"
)

func main() {
	go cmdauth.ExecuteAuth()
	cmdcomm.ExecuteCommodity()
}
