package model

import (
	"github.com/bbdshow/bkit/typ"
	"strings"
)

type OIDSliceStr string

func (str OIDSliceStr) Unmarshal() []typ.ObjectID {
	ids := make([]typ.ObjectID, 0)
	for _, s := range strings.Split(string(str), ",") {
		oid, err := typ.ObjectIDFromHex(s)
		if err == nil {
			ids = append(ids, oid)
		}
	}
	return ids
}

func (str OIDSliceStr) Marshal(ids []typ.ObjectID) OIDSliceStr {
	hexes := make([]string, 0)
	for _, i := range ids {
		hexes = append(hexes, i.Hex())
	}
	if len(hexes) <= 0 {
		return ""
	}
	return OIDSliceStr(strings.Join(hexes, ","))
}

func (str OIDSliceStr) Set(i typ.ObjectID) (bool, OIDSliceStr) {
	hit := false
	if strings.Contains(string(str), i.Hex()) {
		hit = true
	}
	if !hit {
		return false, OIDSliceStr(string(str) + "," + i.Hex())
	}
	return true, str
}
