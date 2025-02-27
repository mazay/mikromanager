package internal

import "strconv"

func parseFloat32(s string) (float32, error) {
	val, err := strconv.ParseFloat(s, 32)
	if err != nil {
		return 0, err
	}
	return float32(val), nil
}

func parseInt(s string) (int, error) {
	return strconv.Atoi(s)
}
