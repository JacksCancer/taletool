package talescript

const (
	TokEND = 0
	TokWs = 1
	TokNl = 2
	TokString = 3
	TokIdentifier = 4
	TokOp = 5
	TokColon = 6
	TokComma = 7
	TokRule = 8
	Tok_list_rule = 9
	TokScript = 10
	TokCondition = 11
	Tok_list_condition_comma = 12
	TokAction = 13
	Tok_list_action_comma = 14
	TokArg = 16
	Tok_list_arg_comma = 17
	TokOperator = 19
	TokRefarg = 20
	Tok_list_refarg_comma = 21
)

var TokenNames = [22]string{
	"@",
	"ws",
	"nl",
	"string",
	"identifier",
	"op",
	"colon",
	"comma",
	"rule",
	"_list_rule",
	"script",
	"condition",
	"_list_condition_comma",
	"action",
	"_list_action_comma",
	"\"(\"",
	"arg",
	"_list_arg_comma",
	"\")\"",
	"operator",
	"refarg",
	"_list_refarg_comma",
}

