package kmgYaml

// Flush the buffer if needed.
func flush(emitter *yaml_emitter_t) bool {
	if emitter.buffer_pos+5 >= len(emitter.buffer) {
		return yaml_emitter_flush(emitter)
	}
	return true
}

// Put a character to the output buffer.
func put(emitter *yaml_emitter_t, value byte) bool {
	if emitter.buffer_pos+5 >= len(emitter.buffer) && !yaml_emitter_flush(emitter) {
		return false
	}
	emitter.buffer[emitter.buffer_pos] = value
	emitter.buffer_pos++
	emitter.column++
	return true
}

// Put a line break to the output buffer.
func put_break(emitter *yaml_emitter_t) bool {
	if emitter.buffer_pos+5 >= len(emitter.buffer) && !yaml_emitter_flush(emitter) {
		return false
	}
	switch emitter.line_break {
	case yaml_CR_BREAK:
		emitter.buffer[emitter.buffer_pos] = '\r'
		emitter.buffer_pos += 1
	case yaml_LN_BREAK:
		emitter.buffer[emitter.buffer_pos] = '\n'
		emitter.buffer_pos += 1
	case yaml_CRLN_BREAK:
		emitter.buffer[emitter.buffer_pos+0] = '\r'
		emitter.buffer[emitter.buffer_pos+1] = '\n'
		emitter.buffer_pos += 2
	default:
		panic("unknown line break setting")
	}
	emitter.column = 0
	emitter.line++
	return true
}

// Copy a character from a string into buffer.
func write(emitter *yaml_emitter_t, s []byte, i *int) bool {
	if emitter.buffer_pos+5 >= len(emitter.buffer) && !yaml_emitter_flush(emitter) {
		return false
	}
	p := emitter.buffer_pos
	w := width(s[*i])
	switch w {
	case 4:
		emitter.buffer[p+3] = s[*i+3]
		fallthrough
	case 3:
		emitter.buffer[p+2] = s[*i+2]
		fallthrough
	case 2:
		emitter.buffer[p+1] = s[*i+1]
		fallthrough
	case 1:
		emitter.buffer[p+0] = s[*i+0]
	default:
		panic("unknown character width")
	}
	emitter.column++
	emitter.buffer_pos += w
	*i += w
	return true
}

// Write a whole string into buffer.
func write_all(emitter *yaml_emitter_t, s []byte) bool {
	for i := 0; i < len(s); {
		if !write(emitter, s, &i) {
			return false
		}
	}
	return true
}

// Copy a line break character from a string into buffer.
func write_break(emitter *yaml_emitter_t, s []byte, i *int) bool {
	if s[*i] == '\n' {
		if !put_break(emitter) {
			return false
		}
		*i++
	} else {
		if !write(emitter, s, i) {
			return false
		}
		emitter.column = 0
		emitter.line++
	}
	return true
}

