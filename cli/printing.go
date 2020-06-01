package cli

import (
	"fmt"

	"github.com/logrusorgru/aurora"
)

func printResult(e error) {
	if e != nil {
		fmt.Println(errorColour("An Error Occurred:"))
		fmt.Println(errorColour(e.Error()))
	} else {
		fmt.Println(successColour("Operation Completed"))
	}
}

func objectChangePreviewDecorator(s string) interface{} {
	return aurora.BrightYellow(s)
}

func errorColour(arg interface{}) aurora.Value {
	return aurora.Red(arg)
}

func warningColour(arg interface{}) aurora.Value {
	return aurora.Yellow(arg)
}

func successColour(arg interface{}) aurora.Value {
	return aurora.Green(arg)
}

func passiveColour(arg interface{}) aurora.Value {
	return aurora.BrightBlue(arg)
}

func expectedColour(arg interface{}) aurora.Value {
	return aurora.Cyan(arg)
}

func dataColour(arg interface{}) aurora.Value {
	return aurora.Yellow(arg)
}
