package gcsv

// Escape escapes the , and " char from a byte slice
// If the byte slice contains a comma, it will be surrounded by double quotes
// If the byte slice contains a double quote, it will be escaped with two double quotes
func Escape(b []byte) []byte {
	extraBytes := 0
	containsComma := false
	for i := 0; i < len(b); i++ {
		if b[i] == '"' {
			extraBytes++
		} else if b[i] == ',' {
			containsComma = true
		}
	}

	// If containsComma is true, we need to surround the whole string with quotes
	if containsComma {
		extraBytes += 2
	}

	if extraBytes == 0 {
		return b // No changes required
	}

	newB := make([]byte, len(b)+extraBytes)
	writePos := 0

	// If containsComma, start with a quote
	if containsComma {
		newB[writePos] = '"'
		writePos++
	}

	for readPos := 0; readPos < len(b); readPos++ {
		if b[readPos] == '"' {
			newB[writePos] = '"'
			writePos++
			newB[writePos] = '"'
		} else {
			newB[writePos] = b[readPos]
		}
		writePos++
	}

	// If containsComma, end with a quote
	if containsComma {
		newB[writePos] = '"'
	}

	return newB
}
