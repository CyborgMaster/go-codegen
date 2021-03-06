package main

import "fmt"

//go:generate go-codegen $GOFILE

// cmdGen is a template.  Blank structs are good to use for targeting templates
// as they do not affect the compiled package.
type cmdGen struct{}

type cmd interface {
	Execute() (interface{}, error)
	MustExecute() interface{}
}

type HelloCommand struct {
	// HelloCommand needs to have the `cmd` template invoked upon it.
	// By mixing in cmd, we tell go-codegen so.
	cmdGen `codegen:""`
	Name   string
}

func (cmd *HelloCommand) Execute() (interface{}, error) {
	return "Hello, " + cmd.Name, nil
}

type GoodbyeCommand struct {
	cmdGen `codegen:""`
	Name   string
}

func (cmd *GoodbyeCommand) Execute() (interface{}, error) {
	return "Goodbye, " + cmd.Name, nil
}

func main() {
	var c cmd
	c = &HelloCommand{Name: "You"}
	fmt.Println(c.MustExecute())
	c = &GoodbyeCommand{Name: "You"}
	fmt.Println(c.MustExecute())
}
