package dao

import (
	"context"
	"fmt"
	"testing"
)

func TestDao_GroupRolesMenuId(t *testing.T) {
	menuId, err := d.GroupRolesMenuId(context.Background(), []int64{2})
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(menuId)
}
