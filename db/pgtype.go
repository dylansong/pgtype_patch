package db

import (
	"time"
)

type InfinityModifier int8

// Bits represents the PostgreSQL bit and varbit types.
type Bits struct {
	Bytes []byte
	Len   int32 // Number of bits
	Valid bool
}

type Bool struct {
	Bool  bool
	Valid bool
}

type Date struct {
	Time             time.Time
	InfinityModifier InfinityModifier
	Valid            bool
}

type Vec2 struct {
	X float64
	Y float64
}

type Box struct {
	P     [2]Vec2
	Valid bool
}

type Circle struct {
	P     Vec2
	R     float64
	Valid bool
}

type Float4 struct {
	Float32 float32
	Valid   bool
}

type Float8 struct {
	Float64 float64
	Valid   bool
}

type Int2 struct {
	Int16 int16
	Valid bool
}

type Interval struct {
	Microseconds int64
	Days         int32
	Months       int32
	Valid        bool
}

type JSONCodec struct {
	Marshal   func(v any) ([]byte, error)
	Unmarshal func(data []byte, v any) error
}

type JSONBCodec struct {
	Marshal   func(v any) ([]byte, error)
	Unmarshal func(data []byte, v any) error
}

type Line struct {
	A, B, C float64
	Valid   bool
}

type Lseg struct {
	P     [2]Vec2
	Valid bool
}

type Path struct {
	P      []Vec2
	Closed bool
	Valid  bool
}

type Polygon struct {
	P     []Vec2
	Valid bool
}

type Text struct {
	String string
	Valid  bool
}

type TID struct {
	BlockNumber  uint32
	OffsetNumber uint16
	Valid        bool
}

type Time struct {
	Microseconds int64 // Number of microseconds since midnight
	Valid        bool
}

// Timestamp represents the PostgreSQL timestamp type.
type Timestamp struct {
	Time             time.Time // Time zone will be ignored when encoding to PostgreSQL.
	InfinityModifier InfinityModifier
	Valid            bool
}

// Timestamptz represents the PostgreSQL timestamptz type.
type Timestamptz struct {
	Time             time.Time
	InfinityModifier InfinityModifier
	Valid            bool
}

// Uint32 is the core type that is used to represent PostgreSQL types such as OID, CID, and XID.
type Uint32 struct {
	Uint32 uint32
	Valid  bool
}

type UUID struct {
	Bytes [16]byte
	Valid bool
}

type XMLCodec struct {
	Marshal   func(v any) ([]byte, error)
	Unmarshal func(data []byte, v any) error
}