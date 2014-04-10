package kmgYaml

// Emit an event.
func yaml_emitter_emit(emitter *yaml_emitter_t, event *yaml_event_t) bool {
	emitter.events = append(emitter.events, *event)
	for !yaml_emitter_need_more_events(emitter) {
		event := &emitter.events[emitter.events_head]
		if !yaml_emitter_analyze_event(emitter, event) {
			return false
		}
		if !yaml_emitter_state_machine(emitter, event) {
			return false
		}
		yaml_event_delete(event)
		emitter.events_head++
	}
	return true
}

// Expect STREAM-START.
func yaml_emitter_emit_stream_start(emitter *yaml_emitter_t, event *yaml_event_t) bool {
	if event.typ != yaml_STREAM_START_EVENT {
		return yaml_emitter_set_emitter_error(emitter, "expected STREAM-START")
	}
	if emitter.encoding == yaml_ANY_ENCODING {
		emitter.encoding = event.encoding
		if emitter.encoding == yaml_ANY_ENCODING {
			emitter.encoding = yaml_UTF8_ENCODING
		}
	}
	if emitter.best_indent < 2 || emitter.best_indent > 9 {
		emitter.best_indent = 2
	}
	if emitter.best_width >= 0 && emitter.best_width <= emitter.best_indent*2 {
		emitter.best_width = 80
	}
	if emitter.best_width < 0 {
		emitter.best_width = 1<<31 - 1
	}
	if emitter.line_break == yaml_ANY_BREAK {
		emitter.line_break = yaml_LN_BREAK
	}

	emitter.indent = -1
	emitter.line = 0
	emitter.column = 0
	emitter.whitespace = true
	emitter.indention = true

	if emitter.encoding != yaml_UTF8_ENCODING {
		if !yaml_emitter_write_bom(emitter) {
			return false
		}
	}
	emitter.state = yaml_EMIT_FIRST_DOCUMENT_START_STATE
	return true
}

// Expect DOCUMENT-START or STREAM-END.
func yaml_emitter_emit_document_start(emitter *yaml_emitter_t, event *yaml_event_t, first bool) bool {

	if event.typ == yaml_DOCUMENT_START_EVENT {

		if event.version_directive != nil {
			if !yaml_emitter_analyze_version_directive(emitter, event.version_directive) {
				return false
			}
		}

		for i := 0; i < len(event.tag_directives); i++ {
			tag_directive := &event.tag_directives[i]
			if !yaml_emitter_analyze_tag_directive(emitter, tag_directive) {
				return false
			}
			if !yaml_emitter_append_tag_directive(emitter, tag_directive, false) {
				return false
			}
		}

		for i := 0; i < len(default_tag_directives); i++ {
			tag_directive := &default_tag_directives[i]
			if !yaml_emitter_append_tag_directive(emitter, tag_directive, true) {
				return false
			}
		}

		implicit := event.implicit
		if !first || emitter.canonical {
			implicit = false
		}

		if emitter.open_ended && (event.version_directive != nil || len(event.tag_directives) > 0) {
			if !yaml_emitter_write_indicator(emitter, []byte("..."), true, false, false) {
				return false
			}
			if !yaml_emitter_write_indent(emitter) {
				return false
			}
		}

		if event.version_directive != nil {
			implicit = false
			if !yaml_emitter_write_indicator(emitter, []byte("%YAML"), true, false, false) {
				return false
			}
			if !yaml_emitter_write_indicator(emitter, []byte("1.1"), true, false, false) {
				return false
			}
			if !yaml_emitter_write_indent(emitter) {
				return false
			}
		}

		if len(event.tag_directives) > 0 {
			implicit = false
			for i := 0; i < len(event.tag_directives); i++ {
				tag_directive := &event.tag_directives[i]
				if !yaml_emitter_write_indicator(emitter, []byte("%TAG"), true, false, false) {
					return false
				}
				if !yaml_emitter_write_tag_handle(emitter, tag_directive.handle) {
					return false
				}
				if !yaml_emitter_write_tag_content(emitter, tag_directive.prefix, true) {
					return false
				}
				if !yaml_emitter_write_indent(emitter) {
					return false
				}
			}
		}

		if yaml_emitter_check_empty_document(emitter) {
			implicit = false
		}
		if !implicit {
			if !yaml_emitter_write_indent(emitter) {
				return false
			}
			if !yaml_emitter_write_indicator(emitter, []byte("---"), true, false, false) {
				return false
			}
			if emitter.canonical {
				if !yaml_emitter_write_indent(emitter) {
					return false
				}
			}
		}

		emitter.state = yaml_EMIT_DOCUMENT_CONTENT_STATE
		return true
	}

	if event.typ == yaml_STREAM_END_EVENT {
		if emitter.open_ended {
			if !yaml_emitter_write_indicator(emitter, []byte("..."), true, false, false) {
				return false
			}
			if !yaml_emitter_write_indent(emitter) {
				return false
			}
		}
		if !yaml_emitter_flush(emitter) {
			return false
		}
		emitter.state = yaml_EMIT_END_STATE
		return true
	}

	return yaml_emitter_set_emitter_error(emitter, "expected DOCUMENT-START or STREAM-END")
}

// Expect the root node.
func yaml_emitter_emit_document_content(emitter *yaml_emitter_t, event *yaml_event_t) bool {
	emitter.states = append(emitter.states, yaml_EMIT_DOCUMENT_END_STATE)
	return yaml_emitter_emit_node(emitter, event, true, false, false, false)
}

// Expect DOCUMENT-END.
func yaml_emitter_emit_document_end(emitter *yaml_emitter_t, event *yaml_event_t) bool {
	if event.typ != yaml_DOCUMENT_END_EVENT {
		return yaml_emitter_set_emitter_error(emitter, "expected DOCUMENT-END")
	}
	if !yaml_emitter_write_indent(emitter) {
		return false
	}
	if !event.implicit {
		// [Go] Allocate the slice elsewhere.
		if !yaml_emitter_write_indicator(emitter, []byte("..."), true, false, false) {
			return false
		}
		if !yaml_emitter_write_indent(emitter) {
			return false
		}
	}
	if !yaml_emitter_flush(emitter) {
		return false
	}
	emitter.state = yaml_EMIT_DOCUMENT_START_STATE
	emitter.tag_directives = emitter.tag_directives[:0]
	return true
}

// Expect a flow item node.
func yaml_emitter_emit_flow_sequence_item(emitter *yaml_emitter_t, event *yaml_event_t, first bool) bool {
	if first {
		if !yaml_emitter_write_indicator(emitter, []byte{'['}, true, true, false) {
			return false
		}
		if !yaml_emitter_increase_indent(emitter, true, false) {
			return false
		}
		emitter.flow_level++
	}

	if event.typ == yaml_SEQUENCE_END_EVENT {
		emitter.flow_level--
		emitter.indent = emitter.indents[len(emitter.indents)-1]
		emitter.indents = emitter.indents[:len(emitter.indents)-1]
		if emitter.canonical && !first {
			if !yaml_emitter_write_indicator(emitter, []byte{','}, false, false, false) {
				return false
			}
			if !yaml_emitter_write_indent(emitter) {
				return false
			}
		}
		if !yaml_emitter_write_indicator(emitter, []byte{']'}, false, false, false) {
			return false
		}
		emitter.state = emitter.states[len(emitter.states)-1]
		emitter.states = emitter.states[:len(emitter.states)-1]

		return true
	}

	if !first {
		if !yaml_emitter_write_indicator(emitter, []byte{','}, false, false, false) {
			return false
		}
	}

	if emitter.canonical || emitter.column > emitter.best_width {
		if !yaml_emitter_write_indent(emitter) {
			return false
		}
	}
	emitter.states = append(emitter.states, yaml_EMIT_FLOW_SEQUENCE_ITEM_STATE)
	return yaml_emitter_emit_node(emitter, event, false, true, false, false)
}

// Expect a flow key node.
func yaml_emitter_emit_flow_mapping_key(emitter *yaml_emitter_t, event *yaml_event_t, first bool) bool {
	if first {
		if !yaml_emitter_write_indicator(emitter, []byte{'{'}, true, true, false) {
			return false
		}
		if !yaml_emitter_increase_indent(emitter, true, false) {
			return false
		}
		emitter.flow_level++
	}

	if event.typ == yaml_MAPPING_END_EVENT {
		emitter.flow_level--
		emitter.indent = emitter.indents[len(emitter.indents)-1]
		emitter.indents = emitter.indents[:len(emitter.indents)-1]
		if emitter.canonical && !first {
			if !yaml_emitter_write_indicator(emitter, []byte{','}, false, false, false) {
				return false
			}
			if !yaml_emitter_write_indent(emitter) {
				return false
			}
		}
		if !yaml_emitter_write_indicator(emitter, []byte{'}'}, false, false, false) {
			return false
		}
		emitter.state = emitter.states[len(emitter.states)-1]
		emitter.states = emitter.states[:len(emitter.states)-1]
		return true
	}

	if !first {
		if !yaml_emitter_write_indicator(emitter, []byte{','}, false, false, false) {
			return false
		}
	}
	if emitter.canonical || emitter.column > emitter.best_width {
		if !yaml_emitter_write_indent(emitter) {
			return false
		}
	}

	if !emitter.canonical && yaml_emitter_check_simple_key(emitter) {
		emitter.states = append(emitter.states, yaml_EMIT_FLOW_MAPPING_SIMPLE_VALUE_STATE)
		return yaml_emitter_emit_node(emitter, event, false, false, true, true)
	}
	if !yaml_emitter_write_indicator(emitter, []byte{'?'}, true, false, false) {
		return false
	}
	emitter.states = append(emitter.states, yaml_EMIT_FLOW_MAPPING_VALUE_STATE)
	return yaml_emitter_emit_node(emitter, event, false, false, true, false)
}

// Expect a flow value node.
func yaml_emitter_emit_flow_mapping_value(emitter *yaml_emitter_t, event *yaml_event_t, simple bool) bool {
	if simple {
		if !yaml_emitter_write_indicator(emitter, []byte{':'}, false, false, false) {
			return false
		}
	} else {
		if emitter.canonical || emitter.column > emitter.best_width {
			if !yaml_emitter_write_indent(emitter) {
				return false
			}
		}
		if !yaml_emitter_write_indicator(emitter, []byte{':'}, true, false, false) {
			return false
		}
	}
	emitter.states = append(emitter.states, yaml_EMIT_FLOW_MAPPING_KEY_STATE)
	return yaml_emitter_emit_node(emitter, event, false, false, true, false)
}

// Expect a block item node.
func yaml_emitter_emit_block_sequence_item(emitter *yaml_emitter_t, event *yaml_event_t, first bool) bool {
	if first {
		if !yaml_emitter_increase_indent(emitter, false, emitter.mapping_context && !emitter.indention) {
			return false
		}
	}
	if event.typ == yaml_SEQUENCE_END_EVENT {
		emitter.indent = emitter.indents[len(emitter.indents)-1]
		emitter.indents = emitter.indents[:len(emitter.indents)-1]
		emitter.state = emitter.states[len(emitter.states)-1]
		emitter.states = emitter.states[:len(emitter.states)-1]
		return true
	}
	if !yaml_emitter_write_indent(emitter) {
		return false
	}
	if !yaml_emitter_write_indicator(emitter, []byte{'-'}, true, false, true) {
		return false
	}
	emitter.states = append(emitter.states, yaml_EMIT_BLOCK_SEQUENCE_ITEM_STATE)
	return yaml_emitter_emit_node(emitter, event, false, true, false, false)
}

// Expect a block key node.
func yaml_emitter_emit_block_mapping_key(emitter *yaml_emitter_t, event *yaml_event_t, first bool) bool {
	if first {
		if !yaml_emitter_increase_indent(emitter, false, false) {
			return false
		}
	}
	if event.typ == yaml_MAPPING_END_EVENT {
		emitter.indent = emitter.indents[len(emitter.indents)-1]
		emitter.indents = emitter.indents[:len(emitter.indents)-1]
		emitter.state = emitter.states[len(emitter.states)-1]
		emitter.states = emitter.states[:len(emitter.states)-1]
		return true
	}
	if !yaml_emitter_write_indent(emitter) {
		return false
	}
	if yaml_emitter_check_simple_key(emitter) {
		emitter.states = append(emitter.states, yaml_EMIT_BLOCK_MAPPING_SIMPLE_VALUE_STATE)
		return yaml_emitter_emit_node(emitter, event, false, false, true, true)
	}
	if !yaml_emitter_write_indicator(emitter, []byte{'?'}, true, false, true) {
		return false
	}
	emitter.states = append(emitter.states, yaml_EMIT_BLOCK_MAPPING_VALUE_STATE)
	return yaml_emitter_emit_node(emitter, event, false, false, true, false)
}

// Expect a block value node.
func yaml_emitter_emit_block_mapping_value(emitter *yaml_emitter_t, event *yaml_event_t, simple bool) bool {
	if simple {
		if !yaml_emitter_write_indicator(emitter, []byte{':'}, false, false, false) {
			return false
		}
	} else {
		if !yaml_emitter_write_indent(emitter) {
			return false
		}
		if !yaml_emitter_write_indicator(emitter, []byte{':'}, true, false, true) {
			return false
		}
	}
	emitter.states = append(emitter.states, yaml_EMIT_BLOCK_MAPPING_KEY_STATE)
	return yaml_emitter_emit_node(emitter, event, false, false, true, false)
}

// Expect a node.
func yaml_emitter_emit_node(emitter *yaml_emitter_t, event *yaml_event_t,
	root bool, sequence bool, mapping bool, simple_key bool) bool {

	emitter.root_context = root
	emitter.sequence_context = sequence
	emitter.mapping_context = mapping
	emitter.simple_key_context = simple_key

	switch event.typ {
	case yaml_ALIAS_EVENT:
		return yaml_emitter_emit_alias(emitter, event)
	case yaml_SCALAR_EVENT:
		return yaml_emitter_emit_scalar(emitter, event)
	case yaml_SEQUENCE_START_EVENT:
		return yaml_emitter_emit_sequence_start(emitter, event)
	case yaml_MAPPING_START_EVENT:
		return yaml_emitter_emit_mapping_start(emitter, event)
	default:
		return yaml_emitter_set_emitter_error(emitter,
			"expected SCALAR, SEQUENCE-START, MAPPING-START, or ALIAS")
	}
	return false
}

// Expect ALIAS.
func yaml_emitter_emit_alias(emitter *yaml_emitter_t, event *yaml_event_t) bool {
	if !yaml_emitter_process_anchor(emitter) {
		return false
	}
	emitter.state = emitter.states[len(emitter.states)-1]
	emitter.states = emitter.states[:len(emitter.states)-1]
	return true
}

// Expect SCALAR.
func yaml_emitter_emit_scalar(emitter *yaml_emitter_t, event *yaml_event_t) bool {
	if !yaml_emitter_select_scalar_style(emitter, event) {
		return false
	}
	if !yaml_emitter_process_anchor(emitter) {
		return false
	}
	if !yaml_emitter_process_tag(emitter) {
		return false
	}
	if !yaml_emitter_increase_indent(emitter, true, false) {
		return false
	}
	if !yaml_emitter_process_scalar(emitter) {
		return false
	}
	emitter.indent = emitter.indents[len(emitter.indents)-1]
	emitter.indents = emitter.indents[:len(emitter.indents)-1]
	emitter.state = emitter.states[len(emitter.states)-1]
	emitter.states = emitter.states[:len(emitter.states)-1]
	return true
}

// Expect SEQUENCE-START.
func yaml_emitter_emit_sequence_start(emitter *yaml_emitter_t, event *yaml_event_t) bool {
	if !yaml_emitter_process_anchor(emitter) {
		return false
	}
	if !yaml_emitter_process_tag(emitter) {
		return false
	}
	if emitter.flow_level > 0 || emitter.canonical || event.sequence_style() == yaml_FLOW_SEQUENCE_STYLE ||
			yaml_emitter_check_empty_sequence(emitter) {
		emitter.state = yaml_EMIT_FLOW_SEQUENCE_FIRST_ITEM_STATE
	} else {
		emitter.state = yaml_EMIT_BLOCK_SEQUENCE_FIRST_ITEM_STATE
	}
	return true
}

// Expect MAPPING-START.
func yaml_emitter_emit_mapping_start(emitter *yaml_emitter_t, event *yaml_event_t) bool {
	if !yaml_emitter_process_anchor(emitter) {
		return false
	}
	if !yaml_emitter_process_tag(emitter) {
		return false
	}
	if emitter.flow_level > 0 || emitter.canonical || event.mapping_style() == yaml_FLOW_MAPPING_STYLE ||
			yaml_emitter_check_empty_mapping(emitter) {
		emitter.state = yaml_EMIT_FLOW_MAPPING_FIRST_KEY_STATE
	} else {
		emitter.state = yaml_EMIT_BLOCK_MAPPING_FIRST_KEY_STATE
	}
	return true
}
