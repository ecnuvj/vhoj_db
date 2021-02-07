package util

import "github.com/ecnuvj/vhoj_db/pkg/common"

func CalLimitOffset(pageNo int32, pageSize int32) (int32, int32) {
	if pageNo <= 0 {
		pageNo = 1
	}
	if pageSize <= 0 {
		pageSize = common.DEFAULT_PAGE_SIZE
	}
	return pageSize, (pageNo - 1) * pageSize
}
