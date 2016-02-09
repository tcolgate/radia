package ghs

import "fmt"

type FragmentID Weight

func (f FragmentID) String() string {
	return fmt.Sprintf("F%v:%v:%v", f.float64, f.Msn, f.Lsn)
}
