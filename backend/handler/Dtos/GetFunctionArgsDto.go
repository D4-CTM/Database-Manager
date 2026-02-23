package dtos

import (
	"fmt"
	"strings"
)

type getFunctionArgs struct {
	ObjectName string
	Owner      string
}

func GetFunctionArgs() getFunctionArgs {
	return getFunctionArgs{
		Owner:      "",
		ObjectName: "",
	}
}

func (gf *getFunctionArgs) IsValid() error {
	gf.ObjectName = strings.ToUpper(strings.ToUpper(gf.ObjectName))
	gf.Owner = strings.ToUpper(strings.ToUpper(gf.Owner))
	if gf.ObjectName == "" {
		return fmt.Errorf("Object name not specified")
	}

	return nil
}
