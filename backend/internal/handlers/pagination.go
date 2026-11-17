package handlers

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

// Pagination bounds for list endpoints. Lists are capped so a growing
// collection cannot produce unbounded responses; clients page with
// ?limit=&offset=.
const (
	defaultListLimit = 50
	maxListLimit     = 100
)

// listPagination reads limit/offset query params, clamping them to safe
// bounds. Invalid values fall back to the defaults instead of erroring,
// keeping the public browsing endpoints forgiving.
func listPagination(c *gin.Context) (limit, offset int) {
	limit = defaultListLimit
	if v := c.Query("limit"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			limit = n
		}
	}
	if limit > maxListLimit {
		limit = maxListLimit
	}
	if v := c.Query("offset"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			offset = n
		}
	}
	return limit, offset
}

// pageSlice applies limit/offset to a slice, returning an empty (never
// nil) slice when the offset is past the end.
func pageSlice[T any](items []T, limit, offset int) []T {
	if offset >= len(items) {
		return []T{}
	}
	end := offset + limit
	if end > len(items) {
		end = len(items)
	}
	return items[offset:end]
}
