package util

//隐藏昵称后几位，用*替代
func FormatNickname(nickname string, length int) (formatName string) {
	rs := []rune(nickname)
	rl := len(rs)
	end := length

	if end < 0 {
		end = 0
	}
	if end > rl {
		end = rl
	}
	formatName = string(rs[0:end]) + "*"
	return formatName
}

//手机号中间打*号
func FormatPhoneNo(phoneNo string) (formatString string) {
	if len(phoneNo) < 11{
		return phoneNo
	}
	rs := []rune(phoneNo)

	formatString = string(rs[0:3]) + "****" +  string(rs[7:11])
	return formatString
}

//截取字符串 start 起点下标 length 需要截取的长度
func SubstrByLength(str string, start int, length int) string {
	rs := []rune(str)
	rl := len(rs)
	end := 0

	if start < 0 {
		start = rl - 1 + start
	}
	end = start + length

	if start > end {
		start, end = end, start
	}

	if start < 0 {
		start = 0
	}
	if start > rl {
		start = rl
	}
	if end < 0 {
		end = 0
	}
	if end > rl {
		end = rl
	}

	return string(rs[start:end])
}

//截取字符串 start 起点下标 end 终点下标(不包括)
func SubstrByEnd(str string, start int, end int) string {
	rs := []rune(str)
	length := len(rs)

	if start < 0 || start > length {
		panic("start is wrong")
	}

	if end < 0 || end > length {
		panic("end is wrong")
	}

	return string(rs[start:end])
}

