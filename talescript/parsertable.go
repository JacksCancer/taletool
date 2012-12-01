package talescript

var parsertable = [33]struct {
	context int
	parsermap
}{
	// state 0
	//[_list_rule ->._list_rule rule], (END identifier condition _list_condition_comma)
	//[_list_rule ->.], (END identifier condition _list_condition_comma)
	//[script ->._list_rule], (END)
	//[END -> END.script], (END)
	//<ignore> [END ->.nl], (END ws nl identifier rule _list_rule condition _list_condition_comma)
	//<ignore> [END ->.ws], (END ws nl identifier rule _list_rule condition _list_condition_comma)
	{0, parsermap{
		TokWs:                    shift(1),
		TokNl:                    shift(2),
		Tok_list_rule:            shift(3),
		TokScript:                shift(4),
		TokEND:                   reduce(Tok_list_rule, 0), // _list_rule ->
		TokIdentifier:            reduce(Tok_list_rule, 0), // _list_rule ->
		TokCondition:             reduce(Tok_list_rule, 0), // _list_rule ->
		Tok_list_condition_comma: reduce(Tok_list_rule, 0), // _list_rule ->
	}},
	// state 1
	//<IGNORABLE>
	//<ignore> [END -> ws.], (END ws nl identifier op colon comma rule _list_rule condition _list_condition_comma action _list_action_comma "(" arg _list_arg_comma ")" refarg _list_refarg_comma)
	{0, parsermap{
		TokEND: reduce(TokEND, 1), // END -> ws
	}},
	// state 2
	//<IGNORABLE>
	//<ignore> [END -> nl.], (END ws nl identifier op colon comma rule _list_rule condition _list_condition_comma action _list_action_comma "(" arg _list_arg_comma ")" refarg _list_refarg_comma)
	{0, parsermap{
		TokEND: reduce15, // END -> nl
	}},
	// state 3
	//[_list_rule -> _list_rule.rule], (END identifier condition _list_condition_comma)
	//[script -> _list_rule.], (END)
	//[_list_condition_comma ->._list_condition_comma condition comma], (identifier)
	//[_list_condition_comma ->.], (identifier)
	//[rule ->._list_condition_comma condition colon _list_action_comma action], (END identifier condition _list_condition_comma)
	//<ignore> [END ->.nl], (END ws nl identifier condition _list_condition_comma)
	//<ignore> [END ->.ws], (END ws nl identifier condition _list_condition_comma)
	{0, parsermap{
		TokWs:                    shift(1),
		TokNl:                    shift(2),
		TokRule:                  shift(5),
		Tok_list_condition_comma: shift(6),
		TokEND:                   reduce(TokScript, 1),                // script -> _list_rule
		TokIdentifier:            reduce(Tok_list_condition_comma, 0), // _list_condition_comma ->
	}},
	// state 4
	//[END -> END script.], (END)
	{0, parsermap{
		TokEND: reduce(TokEND, 2), // END -> END script
	}},
	// state 5
	//[_list_rule -> _list_rule rule.], (END identifier condition _list_condition_comma)
	{0, parsermap{
		TokEND: reduce(Tok_list_rule, 2), // _list_rule -> _list_rule rule
	}},
	// state 6
	//[_list_condition_comma -> _list_condition_comma.condition comma], (identifier)
	//[rule -> _list_condition_comma.condition colon _list_action_comma action], (END identifier condition _list_condition_comma)
	//[condition ->.identifier "(" _list_arg_comma arg ")"], (colon comma)
	//<ignore> [END ->.nl], (ws nl identifier)
	//<ignore> [END ->.ws], (ws nl identifier)
	{0, parsermap{
		TokWs:         shift(1),
		TokNl:         shift(2),
		TokIdentifier: pushIdentifier(7),
		TokCondition:  shift(8),
	}},
	// state 7
	//[condition -> identifier."(" _list_arg_comma arg ")"], (colon comma)
	//<ignore> [END ->.nl], (ws nl "(")
	//<ignore> [END ->.ws], (ws nl "(")
	{0, parsermap{
		TokWs: shift(1),
		TokNl: shift(2),
		15:    shift(9),
	}},
	// state 8
	//[_list_condition_comma -> _list_condition_comma condition.comma], (identifier)
	//[rule -> _list_condition_comma condition.colon _list_action_comma action], (END identifier condition _list_condition_comma)
	//<ignore> [END ->.nl], (ws nl colon comma)
	//<ignore> [END ->.ws], (ws nl colon comma)
	{0, parsermap{
		TokWs:    shift(1),
		TokNl:    shift(2),
		TokColon: shift(10),
		TokComma: shift(11),
	}},
	// state 9
	//[_list_arg_comma ->._list_arg_comma arg comma], (identifier)
	//[_list_arg_comma ->.], (identifier)
	//[condition -> identifier "("._list_arg_comma arg ")"], (colon comma)
	//<ignore> [END ->.nl], (ws nl identifier arg _list_arg_comma)
	//<ignore> [END ->.ws], (ws nl identifier arg _list_arg_comma)
	{0, parsermap{
		TokWs:              shift(1),
		TokNl:              shift(2),
		Tok_list_arg_comma: shift(12),
		TokIdentifier:      reduce(Tok_list_arg_comma, 0), // _list_arg_comma ->
	}},
	// state 10
	//[_list_action_comma ->._list_action_comma action comma], (identifier op)
	//[_list_action_comma ->.], (identifier op)
	//[rule -> _list_condition_comma condition colon._list_action_comma action], (END identifier condition _list_condition_comma)
	//<ignore> [END ->.nl], (ws nl identifier op action _list_action_comma)
	//<ignore> [END ->.ws], (ws nl identifier op action _list_action_comma)
	{0, parsermap{
		TokWs: shift(1),
		TokNl: shift(2),
		Tok_list_action_comma: shift(13),
		TokIdentifier:         reduce(Tok_list_action_comma, 0), // _list_action_comma ->
		TokOp:                 reduce(Tok_list_action_comma, 0), // _list_action_comma ->
	}},
	// state 11
	//[_list_condition_comma -> _list_condition_comma condition comma.], (identifier)
	{0, parsermap{
		TokEND: reduce(Tok_list_condition_comma, 3), // _list_condition_comma -> _list_condition_comma condition comma
	}},
	// state 12
	//[_list_arg_comma -> _list_arg_comma.arg comma], (identifier)
	//[condition -> identifier "(" _list_arg_comma.arg ")"], (colon comma)
	//[arg ->.identifier], (comma ")")
	//<ignore> [END ->.nl], (ws nl identifier)
	//<ignore> [END ->.ws], (ws nl identifier)
	{0, parsermap{
		TokWs:         shift(1),
		TokNl:         shift(2),
		TokIdentifier: pushIdentifier(14),
		TokArg:        shift(15),
	}},
	// state 13
	//[_list_action_comma -> _list_action_comma.action comma], (identifier op)
	//[rule -> _list_condition_comma condition colon _list_action_comma.action], (END identifier condition _list_condition_comma)
	//[action ->.op identifier "(" _list_refarg_comma refarg ")"], (END identifier comma condition _list_condition_comma)
	//[action ->.identifier "(" _list_refarg_comma refarg ")"], (END identifier comma condition _list_condition_comma)
	//<ignore> [END ->.nl], (ws nl identifier op)
	//<ignore> [END ->.ws], (ws nl identifier op)
	{0, parsermap{
		TokWs:         shift(1),
		TokNl:         shift(2),
		TokIdentifier: pushIdentifier(16),
		TokOp:         pushOp(17),
		TokAction:     shift(18),
	}},
	// state 14
	//[arg -> identifier.], (comma ")")
	{0, parsermap{
		TokEND: reduce13, // arg -> identifier
	}},
	// state 15
	//[_list_arg_comma -> _list_arg_comma arg.comma], (identifier)
	//[condition -> identifier "(" _list_arg_comma arg.")"], (colon comma)
	//<ignore> [END ->.nl], (ws nl comma ")")
	//<ignore> [END ->.ws], (ws nl comma ")")
	{0, parsermap{
		TokWs:    shift(1),
		TokNl:    shift(2),
		TokComma: shift(19),
		18:       shift(20),
	}},
	// state 16
	//[action -> identifier."(" _list_refarg_comma refarg ")"], (END identifier comma condition _list_condition_comma)
	//<ignore> [END ->.nl], (ws nl "(")
	//<ignore> [END ->.ws], (ws nl "(")
	{0, parsermap{
		TokWs: shift(1),
		TokNl: shift(2),
		15:    shift(21),
	}},
	// state 17
	//[action -> op.identifier "(" _list_refarg_comma refarg ")"], (END identifier comma condition _list_condition_comma)
	//<ignore> [END ->.nl], (ws nl identifier)
	//<ignore> [END ->.ws], (ws nl identifier)
	{0, parsermap{
		TokWs:         shift(1),
		TokNl:         shift(2),
		TokIdentifier: pushIdentifier(22),
	}},
	// state 18
	//[_list_action_comma -> _list_action_comma action.comma], (identifier op)
	//[rule -> _list_condition_comma condition colon _list_action_comma action.], (END identifier condition _list_condition_comma)
	//<ignore> [END ->.nl], (END ws nl identifier comma condition _list_condition_comma)
	//<ignore> [END ->.ws], (END ws nl identifier comma condition _list_condition_comma)
	{0, parsermap{
		TokWs:    shift(1),
		TokNl:    shift(2),
		TokComma: shift(23),
		TokEND:   reduce7, // rule -> _list_condition_comma condition colon _list_action_comma action
	}},
	// state 19
	//[_list_arg_comma -> _list_arg_comma arg comma.], (identifier)
	{0, parsermap{
		TokEND: reduce(Tok_list_arg_comma, 3), // _list_arg_comma -> _list_arg_comma arg comma
	}},
	// state 20
	//[condition -> identifier "(" _list_arg_comma arg ")".], (colon comma)
	{0, parsermap{
		TokEND: reduce10, // condition -> identifier "(" _list_arg_comma arg ")"
	}},
	// state 21
	//[_list_refarg_comma ->._list_refarg_comma refarg comma], (identifier)
	//[_list_refarg_comma ->.], (identifier)
	//[action -> identifier "("._list_refarg_comma refarg ")"], (END identifier comma condition _list_condition_comma)
	//<ignore> [END ->.nl], (ws nl identifier refarg _list_refarg_comma)
	//<ignore> [END ->.ws], (ws nl identifier refarg _list_refarg_comma)
	{0, parsermap{
		TokWs: shift(1),
		TokNl: shift(2),
		Tok_list_refarg_comma: shift(24),
		TokIdentifier:         reduce(Tok_list_refarg_comma, 0), // _list_refarg_comma ->
	}},
	// state 22
	//[action -> op identifier."(" _list_refarg_comma refarg ")"], (END identifier comma condition _list_condition_comma)
	//<ignore> [END ->.nl], (ws nl "(")
	//<ignore> [END ->.ws], (ws nl "(")
	{0, parsermap{
		TokWs: shift(1),
		TokNl: shift(2),
		15:    shift(25),
	}},
	// state 23
	//[_list_action_comma -> _list_action_comma action comma.], (identifier op)
	{0, parsermap{
		TokEND: reduce(Tok_list_action_comma, 3), // _list_action_comma -> _list_action_comma action comma
	}},
	// state 24
	//[_list_refarg_comma -> _list_refarg_comma.refarg comma], (identifier)
	//[refarg ->.identifier], (comma ")")
	//[action -> identifier "(" _list_refarg_comma.refarg ")"], (END identifier comma condition _list_condition_comma)
	//<ignore> [END ->.nl], (ws nl identifier)
	//<ignore> [END ->.ws], (ws nl identifier)
	{0, parsermap{
		TokWs:         shift(1),
		TokNl:         shift(2),
		TokIdentifier: pushIdentifier(26),
		TokRefarg:     shift(27),
	}},
	// state 25
	//[_list_refarg_comma ->._list_refarg_comma refarg comma], (identifier)
	//[_list_refarg_comma ->.], (identifier)
	//[action -> op identifier "("._list_refarg_comma refarg ")"], (END identifier comma condition _list_condition_comma)
	//<ignore> [END ->.nl], (ws nl identifier refarg _list_refarg_comma)
	//<ignore> [END ->.ws], (ws nl identifier refarg _list_refarg_comma)
	{0, parsermap{
		TokWs: shift(1),
		TokNl: shift(2),
		Tok_list_refarg_comma: shift(28),
		TokIdentifier:         reduce(Tok_list_refarg_comma, 0), // _list_refarg_comma ->
	}},
	// state 26
	//[refarg -> identifier.], (comma ")")
	{0, parsermap{
		TokEND: reduce14, // refarg -> identifier
	}},
	// state 27
	//[_list_refarg_comma -> _list_refarg_comma refarg.comma], (identifier)
	//[action -> identifier "(" _list_refarg_comma refarg.")"], (END identifier comma condition _list_condition_comma)
	//<ignore> [END ->.nl], (ws nl comma ")")
	//<ignore> [END ->.ws], (ws nl comma ")")
	{0, parsermap{
		TokWs:    shift(1),
		TokNl:    shift(2),
		TokComma: shift(29),
		18:       shift(30),
	}},
	// state 28
	//[_list_refarg_comma -> _list_refarg_comma.refarg comma], (identifier)
	//[refarg ->.identifier], (comma ")")
	//[action -> op identifier "(" _list_refarg_comma.refarg ")"], (END identifier comma condition _list_condition_comma)
	//<ignore> [END ->.nl], (ws nl identifier)
	//<ignore> [END ->.ws], (ws nl identifier)
	{0, parsermap{
		TokWs:         shift(1),
		TokNl:         shift(2),
		TokIdentifier: pushIdentifier(26),
		TokRefarg:     shift(31),
	}},
	// state 29
	//[_list_refarg_comma -> _list_refarg_comma refarg comma.], (identifier)
	{0, parsermap{
		TokEND: reduce(Tok_list_refarg_comma, 3), // _list_refarg_comma -> _list_refarg_comma refarg comma
	}},
	// state 30
	//[action -> identifier "(" _list_refarg_comma refarg ")".], (END identifier comma condition _list_condition_comma)
	{0, parsermap{
		TokEND: reduce18, // action -> identifier "(" _list_refarg_comma refarg ")"
	}},
	// state 31
	//[_list_refarg_comma -> _list_refarg_comma refarg.comma], (identifier)
	//[action -> op identifier "(" _list_refarg_comma refarg.")"], (END identifier comma condition _list_condition_comma)
	//<ignore> [END ->.nl], (ws nl comma ")")
	//<ignore> [END ->.ws], (ws nl comma ")")
	{0, parsermap{
		TokWs:    shift(1),
		TokNl:    shift(2),
		TokComma: shift(29),
		18:       shift(32),
	}},
	// state 32
	//[action -> op identifier "(" _list_refarg_comma refarg ")".], (END identifier comma condition _list_condition_comma)
	{0, parsermap{
		TokEND: reduce17, // action -> op identifier "(" _list_refarg_comma refarg ")"
	}},
}
