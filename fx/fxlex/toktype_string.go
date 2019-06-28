// Code generated by "stringer -type TokType fxlex.go"; DO NOT EDIT.

package fxlex

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[TokId-1]
	_ = x[TokEOF-2]
	_ = x[TokFunc-3]
	_ = x[TokMain-4]
	_ = x[TokIter-5]
	_ = x[TokIf-6]
	_ = x[TokElse-7]
	_ = x[TokTypeID-8]
	_ = x[TokRecord-9]
	_ = x[TokInit-10]
	_ = x[TokBool-11]
	_ = x[TokCoord-12]
	_ = x[TokInt-13]
	_ = x[TokFloat-14]
	_ = x[TokAsig-15]
	_ = x[TokLPar-40]
	_ = x[TokRPar-41]
	_ = x[TokLCorch-91]
	_ = x[TokRCorch-93]
	_ = x[TokComma-44]
	_ = x[TokPC-59]
	_ = x[TokLKey-123]
	_ = x[TokRKey-125]
	_ = x[TokDot-46]
	_ = x[TokSum-43]
	_ = x[TokRest-45]
	_ = x[TokBar-47]
	_ = x[TokMin-60]
	_ = x[TokMax-62]
	_ = x[TokPorc-37]
	_ = x[TokMul-42]
	_ = x[TokPot-112]
	_ = x[TokOpInt-111]
	_ = x[TokOr-124]
	_ = x[TokAnd-38]
	_ = x[TokNot-33]
	_ = x[TokXOr-94]
	_ = x[TokEqual-61]
}

const (
	_TokType_name_0 = "TokIdTokEOFTokFuncTokMainTokIterTokIfTokElseTokTypeIDTokRecordTokInitTokBoolTokCoordTokIntTokFloatTokAsig"
	_TokType_name_1 = "TokNot"
	_TokType_name_2 = "TokPorcTokAnd"
	_TokType_name_3 = "TokLParTokRParTokMulTokSumTokCommaTokRestTokDotTokBar"
	_TokType_name_4 = "TokPCTokMinTokEqualTokMax"
	_TokType_name_5 = "TokLCorch"
	_TokType_name_6 = "TokRCorchTokXOr"
	_TokType_name_7 = "TokOpIntTokPot"
	_TokType_name_8 = "TokLKeyTokOrTokRKey"
)

var (
	_TokType_index_0 = [...]uint8{0, 5, 11, 18, 25, 32, 37, 44, 53, 62, 69, 76, 84, 90, 98, 105}
	_TokType_index_2 = [...]uint8{0, 7, 13}
	_TokType_index_3 = [...]uint8{0, 7, 14, 20, 26, 34, 41, 47, 53}
	_TokType_index_4 = [...]uint8{0, 5, 11, 19, 25}
	_TokType_index_6 = [...]uint8{0, 9, 15}
	_TokType_index_7 = [...]uint8{0, 8, 14}
	_TokType_index_8 = [...]uint8{0, 7, 12, 19}
)

func (i TokType) String() string {
	switch {
	case 1 <= i && i <= 15:
		i -= 1
		return _TokType_name_0[_TokType_index_0[i]:_TokType_index_0[i+1]]
	case i == 33:
		return _TokType_name_1
	case 37 <= i && i <= 38:
		i -= 37
		return _TokType_name_2[_TokType_index_2[i]:_TokType_index_2[i+1]]
	case 40 <= i && i <= 47:
		i -= 40
		return _TokType_name_3[_TokType_index_3[i]:_TokType_index_3[i+1]]
	case 59 <= i && i <= 62:
		i -= 59
		return _TokType_name_4[_TokType_index_4[i]:_TokType_index_4[i+1]]
	case i == 91:
		return _TokType_name_5
	case 93 <= i && i <= 94:
		i -= 93
		return _TokType_name_6[_TokType_index_6[i]:_TokType_index_6[i+1]]
	case 111 <= i && i <= 112:
		i -= 111
		return _TokType_name_7[_TokType_index_7[i]:_TokType_index_7[i+1]]
	case 123 <= i && i <= 125:
		i -= 123
		return _TokType_name_8[_TokType_index_8[i]:_TokType_index_8[i+1]]
	default:
		return "TokType(" + strconv.FormatInt(int64(i), 10) + ")"
	}
}
