// Code generated by "stringer -type=Location"; DO NOT EDIT.

package slot

import "fmt"

const _Location_name = "NoneWieldHoldHeadNeckBodyAbout_BodyLegsFeetArmsWristHandsFingersWaist"

var _Location_index = [...]uint8{0, 4, 9, 13, 17, 21, 25, 35, 39, 43, 47, 52, 57, 64, 69}

func (i Location) String() string {
	if i < 0 || i >= Location(len(_Location_index)-1) {
		return fmt.Sprintf("Location(%d)", i)
	}
	return _Location_name[_Location_index[i]:_Location_index[i+1]]
}
