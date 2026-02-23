package dtos

import (
	"fmt"
	"slices"
	"strings"
)

var ValidDDLObjectTypes = []string{
	"TABLE",
	"VIEW",
	"INDEX",
	"SEQUENCE",
	"TRIGGER",
	"FUNCTION",
	"PROCEDURE",
	"PACKAGE",
}

type getDdlDto struct {
	Type  string
	Name  string
	Owner string
}

func GetDdlDto() getDdlDto {
	return getDdlDto{
		Owner: "",
		Name:  "",
		Type:  "",
	}
}

func (gd *getDdlDto) IsValid() error {
	gd.Name = strings.ToUpper(strings.ToUpper(gd.Name))
	gd.Type = strings.ToUpper(strings.ToUpper(gd.Type))
	gd.Owner = strings.ToUpper(strings.ToUpper(gd.Owner))
	if gd.Name == "" {
		return fmt.Errorf("Object name not specified")
	}

	if gd.Owner == "" {
		return fmt.Errorf("Object owner not specified")
	}

	if !slices.Contains(ValidDDLObjectTypes, gd.Type) {
		return fmt.Errorf("Object Type is not valid")
	}

	return nil
}
