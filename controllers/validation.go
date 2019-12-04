package controllers

import (
	"github.com/astaxie/beego/validation"
	"fmt"
)

func (this MainController)FormValidation()  {
	valid := validation.Validation{}
	fmt.Printf("valid:%v\n",valid)
}



