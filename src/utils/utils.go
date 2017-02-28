package utils

import (
	"database/sql"
	"strings"
	"github.com/markbates/pop/nulls"
)

func ToNullInt64(id uint64) nulls.Int64 {
	return nulls.Int64{Int64 : int64(id), Valid : id != 0}
}

func ToUInt32(id uint32) nulls.UInt32 {
	return nulls.UInt32{UInt32: id, Valid : id != 0}
}

func ToNullString(s string) sql.NullString {
	return sql.NullString{String : s, Valid : strings.TrimSpace(s) == ""}
}