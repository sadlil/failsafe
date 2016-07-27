package failsafe

import (
	"errors"
	"fmt"
	"testing"
)

func TestBasic(t *testing.T) {
	f := New()
	f.AddHandler(basicRun, basicErr).Run()
	f.AddHandler(basicRunReturnsErr, basicErr).Run()
	f.AddHandler(basicRunPanics, basicErr).Run()
	f.AddHandler(basicParamTest, basicParamErr).WithParam("hello", 5).Run()
}

func basicRun(c *FailsafeContext) error {
	fmt.Println("running run")
	return nil
}

func basicRunReturnsErr(c *FailsafeContext) error {
	fmt.Println("running run with errors")
	return errors.New("new errors")
}

func basicRunPanics(c *FailsafeContext) error {
	panic("hello this is paniced")
	return errors.New("new errors")
}

func basicErr(c *FailsafeContext) error {
	fmt.Println("running err")
	return nil
}

func basicParamTest(c *FailsafeContext) error {
	v, ok := c.Params.Get("hello")
	if ok {
		fmt.Println(v.Int())
	}
	return nil
}

func basicParamErr(c *FailsafeContext) error {
	return nil
}
