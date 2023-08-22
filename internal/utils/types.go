package utils

// usage:
//
//	for it.Next() {
//			elem := it.Elem()
//	}
//
// elem is guaranteed to be not nil if it.Next() returns true
type Iterator[E any] interface {
	Next() bool
	Elem() E
}
