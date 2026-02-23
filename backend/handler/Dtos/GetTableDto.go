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
	if len(gt.Owner) > 0 {
		return fmt.Sprintf("%s.%s", gt.Owner, gt.Table)
	}

	return gt.Table
}

func (gt *getTableDto) IsValid() error {
	gt.Owner = strings.TrimSpace(strings.ToUpper(gt.Owner))
	gt.Table = strings.TrimSpace(strings.ToUpper(gt.Table))
	if gt.Table == "" {
		return fmt.Errorf("Table not specified")
	}

	return nil
}
