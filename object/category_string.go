// Code generated by "stringer -type=Category"; DO NOT EDIT.

package object

import "fmt"

const _Category_name = "NONEWEAPONWANDSTAFFTREASUREARMORFOODOTHER"

var _Category_index = [...]uint8{0, 4, 10, 14, 19, 27, 32, 36, 41}

func (i Category) String() string {
	if i < 0 || i >= Category(len(_Category_index)-1) {
		return fmt.Sprintf("Category(%d)", i)
	}
	return _Category_name[_Category_index[i]:_Category_index[i+1]]
}
