package crypto

import (
	"fmt"
	"testing"
)

func TestRandomNumbers(t *testing.T) {

	tCnt := 0
	for i := 0; i < 100; i++ {
		b, err := RandBool()
		if err != nil {
			t.Error(err)
		}
		if b {
			tCnt++
		}
	}
	if tCnt > 65 || tCnt < 35 {
		t.Error("bool spread", tCnt)
	}

	chars, err := RandCharacters(100, true)
	if err != nil {
		t.Error(err)
	}
	fmt.Println("todo <check>", string(chars))

	chars, err = RandCharacters(100, false)
	if err != nil {
		t.Error(err)
	}
	fmt.Println("todo <check>", string(chars))

	f1 := RandFloat(30000)
	f2 := RandFloat(30000)
	if f1 > 30000 || f2 > 30000 {
		t.Error("float range")
	}
	if f1 == f2 {
		t.Error("float not random")
	}
	fmt.Println(f1, f2)

	i16_1, err := RandInt16(0, 17000)
	if err != nil {
		t.Error(err)
	}
	i16_2, err := RandInt16(0, 17000)
	if err != nil {
		t.Error(err)
	}
	if i16_1 > 17000 || i16_2 > 17000 {
		t.Error("int16 range")
	}
	if i16_1 == i16_2 {
		t.Error("int16 not random")
	}
	fmt.Println(i16_1, i16_2)
}
