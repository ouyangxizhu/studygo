package main

import (
	"errors"
	"fmt"
)

func main() {
	err := testErr()
	fmt.Println(`=======% v========`)
	fmt.Printf("%v\n", err)
	fmt.Println(`=======% #v========`)
	fmt.Printf("%#v\n", err)
	fmt.Println(`=======% +v========`)
	fmt.Printf("%+v\n", err)
	fmt.Println(`=======% s========`)
	fmt.Printf("%s\n", err)
	fmt.Println(`=======% s err.Error()========`)
	fmt.Printf("%s\n", err.Error())
	fmt.Println(`=======% v err.Error()========`)
	fmt.Printf("%v\n", err.Error())
}
func testErr() error {
	return testErr1()
}
func testErr1() error {
	return errors.New("testErr1 error")
}
