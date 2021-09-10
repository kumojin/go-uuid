package uuid

import "testing"

func BenchmarkNewOrderedUUID(b *testing.B) {
	for n := 0; n < b.N; n++ {
		NewOrdered()
	}
}

func BenchmarkNewV1(b *testing.B) {
	for n := 0; n < b.N; n++ {
		NewV1()
	}
}

func BenchmarkNewV4(b *testing.B) {
	for n := 0; n < b.N; n++ {
		NewV4()
	}
}
