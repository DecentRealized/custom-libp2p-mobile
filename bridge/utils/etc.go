package butils

func ToString(input []byte) string {
	if input == nil {
		return ""
	}

	return string(input)
}
