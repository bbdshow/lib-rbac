package model

import (
	"github.com/bbdshow/bkit/typ"
	"time"
)

type TableAt struct {
	UpdatedAt time.Time `bson:"updated_at" gorm:"type:timestamp;NOT NULL;DEFAULT:CURRENT_TIMESTAMP;COMMENT:创建时间"`
	CreatedAt time.Time `bson:"created_at" gorm:"type:timestamp;NOT NULL;DEFAULT:CURRENT_TIMESTAMP;COMMENT:创建时间"`
}

type TableID struct {
	ID  int64        `bson:"-" gorm:"column:id;type:BIGINT(20);NOT NULL;AUTO_INCREMENT;PRIMARY_KEY"`
	OID typ.ObjectID `bson:"_id,omitempty" gorm:"column:oid;type:VARCHAR(24);NOT NULL;index:oid,unique;COMMENT:OID"`
}
