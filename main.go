package main

import (
	"gitlab.com/merrittcorp/fspop/command"
	"gitlab.com/merrittcorp/fspop/message"
)

func main() {
	message.PrintTitle()

	command.Run()
}
