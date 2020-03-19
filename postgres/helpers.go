package postgres

import (
	"database/sql"
	"strings"
)

// ToNullString is a utility method to convert a string into a sql null string.
// If the input string is len 0 it will set a sql.NullString with Valid false.
func ToNullString(s string) sql.NullString {
	return sql.NullString{String: s, Valid: s != ""}
}

// ToNullInt64 is a utility method to convert a int64 into an sql null int64.
// This func never sets the valid field to false.
func ToNullInt64(i int64) sql.NullInt64 {
	return sql.NullInt64{Int64: i, Valid: true}
}

// ToNullFloat64 is a utility method to convert a float64 into an sql null float64.
// This func never sets the valid field to false.
func ToNullFloat64(f float64) sql.NullFloat64 {
	if f == 0 {
		return sql.NullFloat64{Valid: false}
	}
	return sql.NullFloat64{Float64: f, Valid: true}
}

// SanitizeString removes \u0000 (null byte) characters, as PG cannot handle
// them.
//
// Relevant Thread: https://www.postgresql.org/message-id/CAE37PpOn%3DMcGeokmny4tm4FTHmXSG4KydgUJemKqT9XxkrrTmg%40mail.gmail.com
// TLDR: Postgres is written in C and uses null bytes to terminate strings. It
// would be too much work to change across the whole codebase.
func SanitizeString(s string) string {
	s = strings.Replace(s, "\u0000", "", -1)  // actual character
	s = strings.Replace(s, "\\u0000", "", -1) // encoded character

	return s
}
