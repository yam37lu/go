package utils

import (
	"fmt"
	"sort"
)

type compareObject struct {
	No     int
	Offset int
}

type compareObjects []compareObject

func (t compareObjects) Len() int {
	return len(t)
}
func (t compareObjects) Less(i, j int) bool {
	return t[i].Offset > t[j].Offset
}
func (t compareObjects) Swap(i, j int) {
	t[i], t[j] = t[j], t[i]
}

func SYEncodeByXYZ(src string, l, x, y int) string {
	a := l % 8
	b := x % 8
	c := y % 8
	segALen := len(src) / 3
	segBLen := segALen
	key := MD5V([]byte(src))
	segA := key[0:a] + src[0:segALen]
	segB := key[a:a+b] + src[segALen:segALen+segBLen]
	segC := key[a+b:a+b+c] + src[segALen+segBLen:]
	lenSegMap := make(map[int]string)
	lenSegMap[0] = segA
	lenSegMap[1] = segB
	lenSegMap[2] = segC
	segArr := []compareObject{{No: 0, Offset: a}, {No: 1, Offset: b}, {No: 2, Offset: c}}
	sort.Sort(compareObjects(segArr))
	return fmt.Sprintf("%s%s%s",
		lenSegMap[segArr[2].No], lenSegMap[segArr[1].No], lenSegMap[segArr[0].No])
}
