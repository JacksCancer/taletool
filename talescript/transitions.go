package talescript


var tokenizerContexts = [1][]int {
	// Context 0
	{ TokEND, TokWs, TokNl, TokIdentifier, TokOp, TokColon, TokComma, TokRule, Tok_list_rule, TokScript, TokCondition, Tok_list_condition_comma, TokAction, Tok_list_action_comma, 15, TokArg, Tok_list_arg_comma, 18, TokRefarg, Tok_list_refarg_comma },
}


var transitionTable = []struct {
	potential TokenSet
	transitions []transition
} {
	// State 0
	{ TokenSet{ TokWs, TokNl, TokString, TokIdentifier, TokOp, TokColon, TokComma, 15, 18 }, []transition{
		{ '\t', 1, TokenSet{ TokWs } },
		{ '\n', 2, TokenSet{ TokNl } },
		{ '\v', 1, TokenSet{ TokWs } },
		{ '\r', 3, TokenSet{ TokNl } },
		{ '\x0E', -1, TokenSet{ } },
		{ ' ', 1, TokenSet{ TokWs } },
		{ '!', -1, TokenSet{ } },
		{ '"', 4, TokenSet{ } },
		{ '#', -1, TokenSet{ } },
		{ '(', 5, TokenSet{ 15 } },
		{ ')', 5, TokenSet{ 18 } },
		{ '*', -1, TokenSet{ } },
		{ '+', 5, TokenSet{ TokOp } },
		{ ',', 5, TokenSet{ TokComma } },
		{ '-', 5, TokenSet{ TokOp } },
		{ '.', -1, TokenSet{ } },
		{ ':', 5, TokenSet{ TokColon } },
		{ ';', -1, TokenSet{ } },
		{ 'A', 6, TokenSet{ TokIdentifier } },
		{ '[', -1, TokenSet{ } },
		{ 'a', 6, TokenSet{ TokIdentifier } },
		{ '{', -1, TokenSet{ } },
	} },
	// State 1
	{ TokenSet{ TokWs }, []transition{
		{ '\t', 1, TokenSet{ TokWs } },
		{ '\n', -1, TokenSet{ } },
		{ '\v', 1, TokenSet{ TokWs } },
		{ '\r', -1, TokenSet{ } },
		{ ' ', 1, TokenSet{ TokWs } },
		{ '!', -1, TokenSet{ } },
	} },
	// State 2
	{ TokenSet{ TokNl }, []transition{
		{ '\r', 5, TokenSet{ TokNl } },
		{ '\x0E', -1, TokenSet{ } },
	} },
	// State 3
	{ TokenSet{ TokNl }, []transition{
		{ '\n', 5, TokenSet{ TokNl } },
		{ '\v', -1, TokenSet{ } },
	} },
	// State 4
	{ TokenSet{ TokString }, []transition{
		{ '\x00', 7, TokenSet{ } },
		{ '"', -1, TokenSet{ } },
		{ '#', 7, TokenSet{ } },
		{ '\\', 8, TokenSet{ } },
		{ ']', 7, TokenSet{ } },
	} },
	// State 5
	{ TokenSet{ }, []transition{
	} },
	// State 6
	{ TokenSet{ TokIdentifier }, []transition{
		{ '0', 6, TokenSet{ TokIdentifier } },
		{ ':', -1, TokenSet{ } },
		{ 'A', 6, TokenSet{ TokIdentifier } },
		{ '[', -1, TokenSet{ } },
		{ '_', 6, TokenSet{ TokIdentifier } },
		{ '`', -1, TokenSet{ } },
		{ 'a', 6, TokenSet{ TokIdentifier } },
		{ '{', -1, TokenSet{ } },
	} },
	// State 7
	{ TokenSet{ TokString }, []transition{
		{ '\x00', 7, TokenSet{ } },
		{ '"', 5, TokenSet{ TokString } },
		{ '#', 7, TokenSet{ } },
		{ '\\', 8, TokenSet{ } },
		{ ']', 7, TokenSet{ } },
	} },
	// State 8
	{ TokenSet{ TokString }, []transition{
		{ '\x01', 7, TokenSet{ } },
	} },
}



