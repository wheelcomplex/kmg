package kmgYaml

// Write an achor.
func yaml_emitter_process_anchor(emitter *yaml_emitter_t) bool {
	if emitter.anchor_data.anchor == nil {
		return true
	}
	c := []byte{'&'}
	if emitter.anchor_data.alias {
		c[0] = '*'
	}
	if !yaml_emitter_write_indicator(emitter, c, true, false, false) {
		return false
	}
	return yaml_emitter_write_anchor(emitter, emitter.anchor_data.anchor)
}

// Write a tag.
func yaml_emitter_process_tag(emitter *yaml_emitter_t) bool {
	if len(emitter.tag_data.handle) == 0 && len(emitter.tag_data.suffix) == 0 {
		return true
	}
	if len(emitter.tag_data.handle) > 0 {
		if !yaml_emitter_write_tag_handle(emitter, emitter.tag_data.handle) {
			return false
		}
		if len(emitter.tag_data.suffix) > 0 {
			if !yaml_emitter_write_tag_content(emitter, emitter.tag_data.suffix, false) {
				return false
			}
		}
	} else {
		// [Go] Allocate these slices elsewhere.
		if !yaml_emitter_write_indicator(emitter, []byte("!<"), true, false, false) {
			return false
		}
		if !yaml_emitter_write_tag_content(emitter, emitter.tag_data.suffix, false) {
			return false
		}
		if !yaml_emitter_write_indicator(emitter, []byte{'>'}, false, false, false) {
			return false
		}
	}
	return true
}

// Write a scalar.
func yaml_emitter_process_scalar(emitter *yaml_emitter_t) bool {
	switch emitter.scalar_data.style {
	case yaml_PLAIN_SCALAR_STYLE:
		return yaml_emitter_write_plain_scalar(emitter, emitter.scalar_data.value, !emitter.simple_key_context)

	case yaml_SINGLE_QUOTED_SCALAR_STYLE:
		return yaml_emitter_write_single_quoted_scalar(emitter, emitter.scalar_data.value, !emitter.simple_key_context)

	case yaml_DOUBLE_QUOTED_SCALAR_STYLE:
		return yaml_emitter_write_double_quoted_scalar(emitter, emitter.scalar_data.value, !emitter.simple_key_context)

	case yaml_LITERAL_SCALAR_STYLE:
		return yaml_emitter_write_literal_scalar(emitter, emitter.scalar_data.value)

	case yaml_FOLDED_SCALAR_STYLE:
		return yaml_emitter_write_folded_scalar(emitter, emitter.scalar_data.value)
	}
	panic("unknown scalar style")
}
