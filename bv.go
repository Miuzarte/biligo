package biligo

import (
	"strings"
)

func isNumber(s string) bool {
	for _, c := range []byte(s) {
		if c < '0' || c > '9' {
			return false
		}
	}
	return true
}

// AnyToAid 将任意 id 转换为 aid
func AnyToAid(id string) (string, error) {
	if len(id) > 2 && id[:2] == "av" {
		return id[2:], nil
	}
	if isNumber(id) { // aid
		return id, nil
	}
	// bvid
	if len(id) != len(_BV_PREFIX)+_BV_LENGTH || !strings.HasPrefix(id, _BV_PREFIX) {
		return id, wrapErr(ErrInvalidId, id)
	}
	return itoa(BV2AV(id)), nil
}

const (
	_XOR_CODE    = 0x1552356C4CDB
	_MASK_CODE   = 0x7FFFFFFFFFFFF
	_MAX_AID     = 1 << 51
	_BV_PREFIX   = "BV1"
	_BV_LENGTH   = 9
	_BV_ALPHABET = "FcwAPNKTMug3GV5Lj7EJnHpWsx4tb8haYeviqBz6rkCy12mUSDQX9RdoZf"
	_BV_BASE     = 58
)

var (
	encodeMap = [_BV_LENGTH]int{8, 7, 0, 5, 1, 3, 2, 4, 6}
	decodeMap = [_BV_LENGTH]int{6, 4, 2, 3, 1, 5, 0, 7, 8}
)

func AV2BV(aid int64) string {
	var bvid [_BV_LENGTH]byte
	tmp := (_MAX_AID | aid) ^ _XOR_CODE
	for i := range _BV_LENGTH {
		bvid[encodeMap[i]] = _BV_ALPHABET[tmp%_BV_BASE]
		tmp /= _BV_BASE
	}
	return _BV_PREFIX + string(bvid[:])
}

func BV2AV(bvid string) int64 {
	if bvid[:3] != _BV_PREFIX {
		return 0
	}
	bvid = bvid[3:]
	var tmp int64
	for i := range _BV_LENGTH {
		tmp = tmp*_BV_BASE + int64(index(_BV_ALPHABET, bvid[decodeMap[i]]))
	}
	return (tmp & _MASK_CODE) ^ _XOR_CODE
}

func index(s string, c byte) int {
	for i := range s {
		if s[i] == c {
			return i
		}
	}
	return -1
}
