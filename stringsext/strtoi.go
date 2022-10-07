package stringsext

//Parses the given string interpreting its content as a base 10 integral number
//which is returned as a int value. The function also returns the index the first
//character after the number.
func Strtoi(s string) (int, int) {
	j := 0
	i := 0
	for ; j < len(s); j++ {
		c := s[j]
		if c <= '0' || c >= '9' {
			break
		}
		i = i*10 + int((byte(c) - byte('0')))
	}
	return i, j
}
