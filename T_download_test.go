package biligo

import (
	"io"
	"testing"
)

func TestDownloadVideoMp4(t *testing.T) {
	const aid = "114385129378451"
	vd := NewDownloadVideoMp4(t.Context(), aid, "", VIDEO_QN_720)
	size, err := vd.Init()
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("size: %d", size)
	n, err := vd.Start(io.Discard)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("written: %d", n)
	if n != size {
		t.Errorf("written %d != size %d", n, size)
		t.FailNow()
	}
}
