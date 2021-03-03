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

func CalSliceLeftRight(limit int32, offset int32, capacity int32) (left int32, right int32) {
	if offset < capacity {
		left = offset
	} else {
		left = capacity
	}
	if limit+offset < capacity {
		right = limit + offset
	} else {
		right = capacity
	}
	return
}
