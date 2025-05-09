package biligo

import "testing"

func TestBvConvert(t *testing.T) {
	const av int64 = 660188595
	const bv = `BV1vh4y1U71j`
	if BV2AV(bv) != av {
		t.Fail()
	}
	if AV2BV(av) != bv {
		t.Fail()
	}
}
