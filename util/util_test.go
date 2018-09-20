package util

import "testing"

func TestGenShortId(t *testing.T) {
	shortId, err := GenShortId()
	if shortId == "" || err != nil {
		t.Error("GenShortId failed!")
	}

	t.Log("GenShortId test pass")
}

func BenchmarkGenShortId(b *testing.B) {
	for i := 0; i < b.N; i++ {
		GenShortId()
	}
}

func BenchmarkGenShortIdTimeConsuming(b *testing.B) {
	b.StopTimer()

	// 停止压力测试时间后查看原程序能否正常运行，能正常运行才进行压力测试
	shortId, err := GenShortId()
	if shortId == "" || err != nil {
		b.Error(err)
	}
	// 重新开始时间
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		GenShortId()
	}
}
