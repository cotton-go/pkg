package slicex

import "testing"

func Test_Slicex_For(t *testing.T) {
	hook := func(i int) (bool, error) {
		if i == 4 {
			return true, nil
		}

		t.Log("for print i =", i)
		return false, nil
	}

	if err := For(5, hook); err != nil {
		t.Fatal(err)
	}

	t.Log("循环便利完毕")
}

func Test_Slicex_EmptyAtrray(t *testing.T) {
	for i := range EmptyArray(10) {
		t.Log("Empty Item =", i)
	}

	t.Log("循环便利完毕")
}
