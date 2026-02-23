package dtos

type pagination struct {
	Skip int
	Take int
}

func Pagination() pagination {
	return pagination{
		Skip: 0,
		Take: 25,
	}
}
