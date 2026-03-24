package controllers

// defaultPage 获取分页页码
func defaultPage(c interface{ Query(string) string }) int {
	return parseIntDefault(c.Query("page"), 1)
}

// defaultPageSize 获取分页大小
func defaultPageSize(c interface{ Query(string) string }) int {
	return parseIntDefault(c.Query("page_size"), 20)
}
