package test

import (
	"fmt"
)

func test1() error {
	return nil
}

func test2() error {
	err := test1()
	a, b := 1, 2
	if a != b {
		return nil
	}
	if err != nil {
		return fmt.Errorf("this is err %s", err)
	}
	if err != nil {
		return fmt.Errorf("this is err %s", err.Error())
	}
	if err != nil {
		return err
	}
	if err != nil {
		return fmt.Errorf("number is %d, %d", a, b)
	}
	return nil
}