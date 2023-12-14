package dao

import (
	"context"
	"github.com/bbdshow/bkit/tests"
	"testing"
)

func TestDao_FindMenuConfigAsRootOrChildren(t *testing.T) {
	root, children, err := d.FindMenuConfigAsRootOrChildren(context.Background(), []int64{8, 12})
	if err != nil {
		return
	}
	tests.PrintBeautifyJSON(root)
	tests.PrintBeautifyJSON(children)
}
