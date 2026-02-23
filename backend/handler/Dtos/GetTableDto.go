package dtos

import (
	"fmt"
	"strings"
)

type getTableDto struct {
	Table string
	Owner string

	Pagination pagination
}

func GetTableDto() getTableDto {
	return getTableDto{
		Pagination: Pagination(),	
	}
}

func (gt *getTableDto) TableName() string {
	owner := strings.TrimSpace(gt.Owner)
	if len(owner) > 0 {
		return fmt.Sprintf("%s.%s", owner, gt.Table)
	}

	return gt.Table
}

func (gt *getTableDto) IsValid() error {
	if gt.Table == "" {
		return fmt.Errorf("Table not specified")
	}

	return nil
}
