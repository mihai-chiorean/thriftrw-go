// @generated Code generated by yacc

// Copyright (c) 2016 Uber Technologies, Inc.
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.



//line thrift.y:2
package internal

import __yyfmt__ "fmt"

//line thrift.y:2
import "github.com/thriftrw/thriftrw-go/ast"

//line thrift.y:7
type yySymType struct {
	yys int
	// Used to record line numbers when the line number at the start point is
	// required.
	line int

	// Holds the final AST for the file.
	prog *ast.Program

	// Other intermediate variables:

	bul bool
	str string
	i64 int64
	dub float64

	fieldType     ast.Type
	structType    ast.StructureType
	baseTypeID    ast.BaseTypeID
	fieldRequired ast.Requiredness

	field  *ast.Field
	fields []*ast.Field

	header  ast.Header
	headers []ast.Header

	function  *ast.Function
	functions []*ast.Function

	enumItem  *ast.EnumItem
	enumItems []*ast.EnumItem

	definition  ast.Definition
	definitions []ast.Definition

	typeAnnotations []*ast.Annotation

	constantValue    ast.ConstantValue
	constantValues   []ast.ConstantValue
	constantMapItems []ast.ConstantMapItem
}

const IDENTIFIER = 57346
const LITERAL = 57347
const INTCONSTANT = 57348
const DUBCONSTANT = 57349
const NAMESPACE = 57350
const INCLUDE = 57351
const VOID = 57352
const BOOL = 57353
const BYTE = 57354
const I8 = 57355
const I16 = 57356
const I32 = 57357
const I64 = 57358
const DOUBLE = 57359
const STRING = 57360
const BINARY = 57361
const MAP = 57362
const LIST = 57363
const SET = 57364
const ONEWAY = 57365
const TYPEDEF = 57366
const STRUCT = 57367
const UNION = 57368
const EXCEPTION = 57369
const EXTENDS = 57370
const THROWS = 57371
const SERVICE = 57372
const ENUM = 57373
const CONST = 57374
const REQUIRED = 57375
const OPTIONAL = 57376
const TRUE = 57377
const FALSE = 57378

var yyToknames = [...]string{
	"$end",
	"error",
	"$unk",
	"IDENTIFIER",
	"LITERAL",
	"INTCONSTANT",
	"DUBCONSTANT",
	"NAMESPACE",
	"INCLUDE",
	"VOID",
	"BOOL",
	"BYTE",
	"I8",
	"I16",
	"I32",
	"I64",
	"DOUBLE",
	"STRING",
	"BINARY",
	"MAP",
	"LIST",
	"SET",
	"ONEWAY",
	"TYPEDEF",
	"STRUCT",
	"UNION",
	"EXCEPTION",
	"EXTENDS",
	"THROWS",
	"SERVICE",
	"ENUM",
	"CONST",
	"REQUIRED",
	"OPTIONAL",
	"TRUE",
	"FALSE",
	"'*'",
	"'='",
	"'{'",
	"'}'",
	"':'",
	"'('",
	"')'",
	"'<'",
	"','",
	"'>'",
	"'['",
	"']'",
	"';'",
}
var yyStatenames = [...]string{}

const yyEofCode = 1
const yyErrCode = 2
const yyMaxDepth = 200

//line yacctab:1
var yyExca = [...]int{
	-1, 1,
	1, -1,
	-2, 0,
	-1, 2,
	8, 70,
	9, 70,
	-2, 8,
	-1, 3,
	1, 1,
	-2, 70,
}

const yyNprod = 74
const yyPrivate = 57344

var yyTokenNames []string
var yyStates []string

const yyLast = 191

var yyAct = [...]int{

	30, 67, 48, 5, 7, 11, 69, 66, 117, 12,
	74, 70, 71, 10, 119, 11, 74, 70, 71, 12,
	82, 81, 80, 52, 25, 51, 50, 154, 146, 145,
	121, 78, 152, 49, 49, 49, 139, 92, 126, 40,
	72, 73, 122, 86, 77, 92, 72, 73, 83, 112,
	77, 115, 76, 113, 134, 56, 55, 64, 76, 68,
	75, 79, 89, 58, 59, 129, 85, 88, 149, 131,
	132, 9, 8, 106, 57, 61, 62, 63, 74, 70,
	71, 127, 24, 142, 99, 100, 101, 22, 21, 104,
	44, 133, 107, 103, 97, 94, 75, 75, 102, 93,
	54, 105, 114, 116, 108, 98, 120, 53, 72, 73,
	123, 118, 77, 47, 124, 23, 111, 46, 45, 43,
	76, 42, 128, 14, 18, 19, 20, 75, 125, 17,
	15, 13, 137, 135, 41, 148, 109, 140, 91, 60,
	96, 136, 95, 3, 88, 144, 75, 143, 6, 141,
	150, 151, 147, 65, 88, 138, 84, 90, 2, 4,
	153, 110, 31, 32, 33, 34, 35, 36, 37, 38,
	39, 27, 28, 29, 31, 32, 33, 34, 35, 36,
	37, 38, 39, 27, 28, 29, 87, 16, 130, 26,
	1,
}
var yyPact = [...]int{

	-1000, -1000, -1000, -1000, -1000, 63, -40, 99, 83, 78,
	-1000, -1000, -1000, 163, 163, 130, 117, 115, -1000, -1000,
	-1000, -1000, 85, 114, 113, 109, -7, -18, -19, -21,
	103, -1000, -1000, -1000, -1000, -1000, -1000, -1000, -1000, -1000,
	96, 17, 16, 35, -1000, -1000, -1000, 26, -1000, -1000,
	163, 163, 163, -1000, -7, -1000, -1000, -1000, -1000, 73,
	-12, -23, -25, -26, -1000, 8, 3, 22, 95, -1000,
	-1000, -1000, -1000, -1000, -1000, 91, -1000, -1000, -1000, 90,
	163, -7, -7, -7, -40, 89, -7, -40, 67, -7,
	-40, 151, -1000, 10, -1000, 5, 11, -30, -32, -1000,
	-1000, -1000, -1000, -8, -1000, -1000, 1, -1000, -1000, -1000,
	-1000, -1000, -1000, -1000, -40, -1000, -3, 76, -1000, -7,
	-1000, 59, 36, 87, 14, -1000, 73, -40, -1000, -7,
	163, -1000, -1000, -6, -7, -40, -1000, -1000, 79, -1000,
	-1000, -1000, -9, -15, -1000, 73, 39, -7, -7, -10,
	-1000, -1000, -1000, -16, -1000,
}
var yyPgo = [...]int{

	0, 0, 190, 24, 189, 188, 187, 186, 7, 159,
	158, 157, 1, 156, 153, 148, 143, 6, 142, 140,
	139, 2, 13, 138, 136, 135,
}
var yyR1 = [...]int{

	0, 2, 10, 10, 9, 9, 9, 9, 16, 16,
	15, 15, 15, 15, 15, 15, 6, 6, 6, 14,
	14, 13, 13, 8, 8, 7, 7, 5, 5, 5,
	12, 12, 11, 23, 23, 24, 24, 25, 25, 3,
	3, 3, 3, 3, 4, 4, 4, 4, 4, 4,
	4, 4, 4, 17, 17, 17, 17, 17, 17, 17,
	17, 18, 18, 19, 19, 21, 21, 20, 20, 20,
	1, 22, 22, 22,
}
var yyR2 = [...]int{

	0, 2, 0, 2, 3, 4, 4, 4, 0, 3,
	6, 5, 7, 7, 7, 10, 1, 1, 1, 0,
	3, 3, 5, 0, 3, 7, 9, 1, 1, 0,
	0, 3, 9, 1, 0, 1, 1, 0, 4, 2,
	7, 5, 5, 2, 1, 1, 1, 1, 1, 1,
	1, 1, 1, 1, 1, 1, 1, 1, 2, 3,
	3, 0, 3, 0, 5, 0, 3, 0, 6, 4,
	0, 1, 1, 0,
}
var yyChk = [...]int{

	-1000, -2, -10, -16, -9, -1, -15, -1, 9, 8,
	-22, 45, 49, 32, 24, 31, -6, 30, 25, 26,
	27, 5, 4, 37, 4, -3, -4, 20, 21, 22,
	-1, 11, 12, 13, 14, 15, 16, 17, 18, 19,
	-3, 4, 4, 4, 5, 4, 4, 4, -21, 42,
	44, 44, 44, 4, 4, 39, 39, 39, 28, 38,
	-20, -3, -3, -3, -21, -14, -8, -12, -1, -17,
	6, 7, 35, 36, 5, -1, 47, 39, 43, -1,
	45, 46, 46, 40, -13, -1, 40, -7, -1, 40,
	-11, -23, 23, 4, 4, -18, -19, 4, -3, -21,
	-21, -21, -22, 4, -21, -22, 6, -21, -22, -24,
	10, -3, 39, 48, -17, 40, -17, 38, -22, 46,
	-21, 38, 41, -1, -12, -22, 41, 5, -21, 6,
	-5, 33, 34, 4, 40, -17, -22, -21, -3, 42,
	-21, -22, 4, -8, -21, 38, 43, -17, -25, 29,
	-21, -21, 42, -8, 43,
}
var yyDef = [...]int{

	2, -2, -2, -2, 3, 0, 73, 0, 0, 0,
	9, 71, 72, 70, 70, 0, 0, 0, 16, 17,
	18, 4, 0, 0, 0, 0, 65, 0, 0, 0,
	0, 44, 45, 46, 47, 48, 49, 50, 51, 52,
	0, 0, 0, 0, 5, 6, 7, 0, 39, 67,
	70, 70, 70, 43, 65, 19, 23, 30, 70, 70,
	70, 0, 0, 0, 11, 70, 70, 34, 0, 10,
	53, 54, 55, 56, 57, 0, 61, 63, 66, 0,
	70, 65, 65, 65, 73, 0, 65, 73, 0, 65,
	73, 70, 33, 0, 58, 70, 70, 73, 0, 41,
	42, 12, 20, 65, 13, 24, 0, 14, 31, 70,
	35, 36, 30, 59, 73, 60, 0, 0, 69, 65,
	21, 0, 29, 0, 34, 62, 70, 73, 40, 65,
	70, 27, 28, 0, 65, 73, 68, 22, 0, 23,
	15, 64, 65, 70, 25, 70, 37, 65, 65, 0,
	26, 32, 23, 70, 38,
}
var yyTok1 = [...]int{

	1, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	42, 43, 37, 3, 45, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 41, 49,
	44, 38, 46, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 47, 3, 48, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 39, 3, 40,
}
var yyTok2 = [...]int{

	2, 3, 4, 5, 6, 7, 8, 9, 10, 11,
	12, 13, 14, 15, 16, 17, 18, 19, 20, 21,
	22, 23, 24, 25, 26, 27, 28, 29, 30, 31,
	32, 33, 34, 35, 36,
}
var yyTok3 = [...]int{
	0,
}

var yyErrorMessages = [...]struct {
	state int
	token int
	msg   string
}{}

//line yaccpar:1

/*	parser for yacc output	*/

var (
	yyDebug        = 0
	yyErrorVerbose = false
)

type yyLexer interface {
	Lex(lval *yySymType) int
	Error(s string)
}

type yyParser interface {
	Parse(yyLexer) int
	Lookahead() int
}

type yyParserImpl struct {
	lookahead func() int
}

func (p *yyParserImpl) Lookahead() int {
	return p.lookahead()
}

func yyNewParser() yyParser {
	p := &yyParserImpl{
		lookahead: func() int { return -1 },
	}
	return p
}

const yyFlag = -1000

func yyTokname(c int) string {
	if c >= 1 && c-1 < len(yyToknames) {
		if yyToknames[c-1] != "" {
			return yyToknames[c-1]
		}
	}
	return __yyfmt__.Sprintf("tok-%v", c)
}

func yyStatname(s int) string {
	if s >= 0 && s < len(yyStatenames) {
		if yyStatenames[s] != "" {
			return yyStatenames[s]
		}
	}
	return __yyfmt__.Sprintf("state-%v", s)
}

func yyErrorMessage(state, lookAhead int) string {
	const TOKSTART = 4

	if !yyErrorVerbose {
		return "syntax error"
	}

	for _, e := range yyErrorMessages {
		if e.state == state && e.token == lookAhead {
			return "syntax error: " + e.msg
		}
	}

	res := "syntax error: unexpected " + yyTokname(lookAhead)

	// To match Bison, suggest at most four expected tokens.
	expected := make([]int, 0, 4)

	// Look for shiftable tokens.
	base := yyPact[state]
	for tok := TOKSTART; tok-1 < len(yyToknames); tok++ {
		if n := base + tok; n >= 0 && n < yyLast && yyChk[yyAct[n]] == tok {
			if len(expected) == cap(expected) {
				return res
			}
			expected = append(expected, tok)
		}
	}

	if yyDef[state] == -2 {
		i := 0
		for yyExca[i] != -1 || yyExca[i+1] != state {
			i += 2
		}

		// Look for tokens that we accept or reduce.
		for i += 2; yyExca[i] >= 0; i += 2 {
			tok := yyExca[i]
			if tok < TOKSTART || yyExca[i+1] == 0 {
				continue
			}
			if len(expected) == cap(expected) {
				return res
			}
			expected = append(expected, tok)
		}

		// If the default action is to accept or reduce, give up.
		if yyExca[i+1] != 0 {
			return res
		}
	}

	for i, tok := range expected {
		if i == 0 {
			res += ", expecting "
		} else {
			res += " or "
		}
		res += yyTokname(tok)
	}
	return res
}

func yylex1(lex yyLexer, lval *yySymType) (char, token int) {
	token = 0
	char = lex.Lex(lval)
	if char <= 0 {
		token = yyTok1[0]
		goto out
	}
	if char < len(yyTok1) {
		token = yyTok1[char]
		goto out
	}
	if char >= yyPrivate {
		if char < yyPrivate+len(yyTok2) {
			token = yyTok2[char-yyPrivate]
			goto out
		}
	}
	for i := 0; i < len(yyTok3); i += 2 {
		token = yyTok3[i+0]
		if token == char {
			token = yyTok3[i+1]
			goto out
		}
	}

out:
	if token == 0 {
		token = yyTok2[1] /* unknown char */
	}
	if yyDebug >= 3 {
		__yyfmt__.Printf("lex %s(%d)\n", yyTokname(token), uint(char))
	}
	return char, token
}

func yyParse(yylex yyLexer) int {
	return yyNewParser().Parse(yylex)
}

func (yyrcvr *yyParserImpl) Parse(yylex yyLexer) int {
	var yyn int
	var yylval yySymType
	var yyVAL yySymType
	var yyDollar []yySymType
	_ = yyDollar // silence set and not used
	yyS := make([]yySymType, yyMaxDepth)

	Nerrs := 0   /* number of errors */
	Errflag := 0 /* error recovery flag */
	yystate := 0
	yychar := -1
	yytoken := -1 // yychar translated into internal numbering
	yyrcvr.lookahead = func() int { return yychar }
	defer func() {
		// Make sure we report no lookahead when not parsing.
		yystate = -1
		yychar = -1
		yytoken = -1
	}()
	yyp := -1
	goto yystack

ret0:
	return 0

ret1:
	return 1

yystack:
	/* put a state and value onto the stack */
	if yyDebug >= 4 {
		__yyfmt__.Printf("char %v in %v\n", yyTokname(yytoken), yyStatname(yystate))
	}

	yyp++
	if yyp >= len(yyS) {
		nyys := make([]yySymType, len(yyS)*2)
		copy(nyys, yyS)
		yyS = nyys
	}
	yyS[yyp] = yyVAL
	yyS[yyp].yys = yystate

yynewstate:
	yyn = yyPact[yystate]
	if yyn <= yyFlag {
		goto yydefault /* simple state */
	}
	if yychar < 0 {
		yychar, yytoken = yylex1(yylex, &yylval)
	}
	yyn += yytoken
	if yyn < 0 || yyn >= yyLast {
		goto yydefault
	}
	yyn = yyAct[yyn]
	if yyChk[yyn] == yytoken { /* valid shift */
		yychar = -1
		yytoken = -1
		yyVAL = yylval
		yystate = yyn
		if Errflag > 0 {
			Errflag--
		}
		goto yystack
	}

yydefault:
	/* default state action */
	yyn = yyDef[yystate]
	if yyn == -2 {
		if yychar < 0 {
			yychar, yytoken = yylex1(yylex, &yylval)
		}

		/* look through exception table */
		xi := 0
		for {
			if yyExca[xi+0] == -1 && yyExca[xi+1] == yystate {
				break
			}
			xi += 2
		}
		for xi += 2; ; xi += 2 {
			yyn = yyExca[xi+0]
			if yyn < 0 || yyn == yytoken {
				break
			}
		}
		yyn = yyExca[xi+1]
		if yyn < 0 {
			goto ret0
		}
	}
	if yyn == 0 {
		/* error ... attempt to resume parsing */
		switch Errflag {
		case 0: /* brand new error */
			yylex.Error(yyErrorMessage(yystate, yytoken))
			Nerrs++
			if yyDebug >= 1 {
				__yyfmt__.Printf("%s", yyStatname(yystate))
				__yyfmt__.Printf(" saw %s\n", yyTokname(yytoken))
			}
			fallthrough

		case 1, 2: /* incompletely recovered error ... try again */
			Errflag = 3

			/* find a state where "error" is a legal shift action */
			for yyp >= 0 {
				yyn = yyPact[yyS[yyp].yys] + yyErrCode
				if yyn >= 0 && yyn < yyLast {
					yystate = yyAct[yyn] /* simulate a shift of "error" */
					if yyChk[yystate] == yyErrCode {
						goto yystack
					}
				}

				/* the current p has no shift on "error", pop stack */
				if yyDebug >= 2 {
					__yyfmt__.Printf("error recovery pops state %d\n", yyS[yyp].yys)
				}
				yyp--
			}
			/* there is no state on the stack with an error shift ... abort */
			goto ret1

		case 3: /* no shift yet; clobber input char */
			if yyDebug >= 2 {
				__yyfmt__.Printf("error recovery discards %s\n", yyTokname(yytoken))
			}
			if yytoken == yyEofCode {
				goto ret1
			}
			yychar = -1
			yytoken = -1
			goto yynewstate /* try again in the same state */
		}
	}

	/* reduction by production yyn */
	if yyDebug >= 2 {
		__yyfmt__.Printf("reduce %v in:\n\t%v\n", yyn, yyStatname(yystate))
	}

	yynt := yyn
	yypt := yyp
	_ = yypt // guard against "declared and not used"

	yyp -= yyR2[yyn]
	// yyp is now the index of $0. Perform the default action. Iff the
	// reduced production is ε, $1 is possibly out of range.
	if yyp+1 >= len(yyS) {
		nyys := make([]yySymType, len(yyS)*2)
		copy(nyys, yyS)
		yyS = nyys
	}
	yyVAL = yyS[yyp+1]

	/* consult goto table to find next state */
	yyn = yyR1[yyn]
	yyg := yyPgo[yyn]
	yyj := yyg + yyS[yyp].yys + 1

	if yyj >= yyLast {
		yystate = yyAct[yyg]
	} else {
		yystate = yyAct[yyj]
		if yyChk[yystate] != -yyn {
			yystate = yyAct[yyg]
		}
	}
	// dummy call; replaced with literal code
	switch yynt {

	case 1:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line thrift.y:92
		{
			yyVAL.prog = &ast.Program{Headers: yyDollar[1].headers, Definitions: yyDollar[2].definitions}
			yylex.(*lexer).program = yyVAL.prog
			return 0
		}
	case 2:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line thrift.y:104
		{
			yyVAL.headers = nil
		}
	case 3:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line thrift.y:105
		{
			yyVAL.headers = append(yyDollar[1].headers, yyDollar[2].header)
		}
	case 4:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line thrift.y:110
		{
			yyVAL.header = &ast.Include{
				Path: yyDollar[3].str,
				Line: yyDollar[1].line,
			}
		}
	case 5:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line thrift.y:117
		{
			yyVAL.header = &ast.Include{
				Name: yyDollar[3].str,
				Path: yyDollar[4].str,
				Line: yyDollar[1].line,
			}
		}
	case 6:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line thrift.y:125
		{
			yyVAL.header = &ast.Namespace{
				Scope: "*",
				Name:  yyDollar[4].str,
				Line:  yyDollar[1].line,
			}
		}
	case 7:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line thrift.y:133
		{
			yyVAL.header = &ast.Namespace{
				Scope: yyDollar[3].str,
				Name:  yyDollar[4].str,
				Line:  yyDollar[1].line,
			}
		}
	case 8:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line thrift.y:147
		{
			yyVAL.definitions = nil
		}
	case 9:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line thrift.y:148
		{
			yyVAL.definitions = append(yyDollar[1].definitions, yyDollar[2].definition)
		}
	case 10:
		yyDollar = yyS[yypt-6 : yypt+1]
		//line thrift.y:155
		{
			yyVAL.definition = &ast.Constant{
				Name:  yyDollar[4].str,
				Type:  yyDollar[3].fieldType,
				Value: yyDollar[6].constantValue,
				Line:  yyDollar[1].line,
			}
		}
	case 11:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line thrift.y:165
		{
			yyVAL.definition = &ast.Typedef{
				Name:        yyDollar[4].str,
				Type:        yyDollar[3].fieldType,
				Annotations: yyDollar[5].typeAnnotations,
				Line:        yyDollar[1].line,
			}
		}
	case 12:
		yyDollar = yyS[yypt-7 : yypt+1]
		//line thrift.y:174
		{
			yyVAL.definition = &ast.Enum{
				Name:        yyDollar[3].str,
				Items:       yyDollar[5].enumItems,
				Annotations: yyDollar[7].typeAnnotations,
				Line:        yyDollar[1].line,
			}
		}
	case 13:
		yyDollar = yyS[yypt-7 : yypt+1]
		//line thrift.y:183
		{
			yyVAL.definition = &ast.Struct{
				Name:        yyDollar[3].str,
				Type:        yyDollar[2].structType,
				Fields:      yyDollar[5].fields,
				Annotations: yyDollar[7].typeAnnotations,
				Line:        yyDollar[1].line,
			}
		}
	case 14:
		yyDollar = yyS[yypt-7 : yypt+1]
		//line thrift.y:194
		{
			yyVAL.definition = &ast.Service{
				Name:        yyDollar[3].str,
				Functions:   yyDollar[5].functions,
				Annotations: yyDollar[7].typeAnnotations,
				Line:        yyDollar[1].line,
			}
		}
	case 15:
		yyDollar = yyS[yypt-10 : yypt+1]
		//line thrift.y:204
		{
			parent := &ast.ServiceReference{
				Name: yyDollar[6].str,
				Line: yyDollar[5].line,
			}

			yyVAL.definition = &ast.Service{
				Name:        yyDollar[3].str,
				Functions:   yyDollar[8].functions,
				Parent:      parent,
				Annotations: yyDollar[10].typeAnnotations,
				Line:        yyDollar[1].line,
			}
		}
	case 16:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line thrift.y:221
		{
			yyVAL.structType = ast.StructType
		}
	case 17:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line thrift.y:222
		{
			yyVAL.structType = ast.UnionType
		}
	case 18:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line thrift.y:223
		{
			yyVAL.structType = ast.ExceptionType
		}
	case 19:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line thrift.y:227
		{
			yyVAL.enumItems = nil
		}
	case 20:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line thrift.y:228
		{
			yyVAL.enumItems = append(yyDollar[1].enumItems, yyDollar[2].enumItem)
		}
	case 21:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line thrift.y:233
		{
			yyVAL.enumItem = &ast.EnumItem{Name: yyDollar[2].str, Annotations: yyDollar[3].typeAnnotations, Line: yyDollar[1].line}
		}
	case 22:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line thrift.y:235
		{
			value := int(yyDollar[4].i64)
			yyVAL.enumItem = &ast.EnumItem{
				Name:        yyDollar[2].str,
				Value:       &value,
				Annotations: yyDollar[5].typeAnnotations,
				Line:        yyDollar[1].line,
			}
		}
	case 23:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line thrift.y:247
		{
			yyVAL.fields = nil
		}
	case 24:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line thrift.y:248
		{
			yyVAL.fields = append(yyDollar[1].fields, yyDollar[2].field)
		}
	case 25:
		yyDollar = yyS[yypt-7 : yypt+1]
		//line thrift.y:254
		{
			yyVAL.field = &ast.Field{
				ID:           int(yyDollar[2].i64),
				Name:         yyDollar[6].str,
				Type:         yyDollar[5].fieldType,
				Requiredness: yyDollar[4].fieldRequired,
				Annotations:  yyDollar[7].typeAnnotations,
				Line:         yyDollar[1].line,
			}
		}
	case 26:
		yyDollar = yyS[yypt-9 : yypt+1]
		//line thrift.y:266
		{
			yyVAL.field = &ast.Field{
				ID:           int(yyDollar[2].i64),
				Name:         yyDollar[6].str,
				Type:         yyDollar[5].fieldType,
				Requiredness: yyDollar[4].fieldRequired,
				Default:      yyDollar[8].constantValue,
				Annotations:  yyDollar[9].typeAnnotations,
				Line:         yyDollar[1].line,
			}
		}
	case 27:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line thrift.y:280
		{
			yyVAL.fieldRequired = ast.Required
		}
	case 28:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line thrift.y:281
		{
			yyVAL.fieldRequired = ast.Optional
		}
	case 29:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line thrift.y:282
		{
			yyVAL.fieldRequired = ast.Unspecified
		}
	case 30:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line thrift.y:286
		{
			yyVAL.functions = nil
		}
	case 31:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line thrift.y:287
		{
			yyVAL.functions = append(yyDollar[1].functions, yyDollar[2].function)
		}
	case 32:
		yyDollar = yyS[yypt-9 : yypt+1]
		//line thrift.y:293
		{
			yyVAL.function = &ast.Function{
				Name:        yyDollar[4].str,
				Parameters:  yyDollar[6].fields,
				ReturnType:  yyDollar[2].fieldType,
				Exceptions:  yyDollar[8].fields,
				OneWay:      yyDollar[1].bul,
				Annotations: yyDollar[9].typeAnnotations,
				Line:        yyDollar[3].line,
			}
		}
	case 33:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line thrift.y:307
		{
			yyVAL.bul = true
		}
	case 34:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line thrift.y:308
		{
			yyVAL.bul = false
		}
	case 35:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line thrift.y:312
		{
			yyVAL.fieldType = nil
		}
	case 36:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line thrift.y:313
		{
			yyVAL.fieldType = yyDollar[1].fieldType
		}
	case 37:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line thrift.y:317
		{
			yyVAL.fields = nil
		}
	case 38:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line thrift.y:318
		{
			yyVAL.fields = yyDollar[3].fields
		}
	case 39:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line thrift.y:327
		{
			yyVAL.fieldType = ast.BaseType{ID: yyDollar[1].baseTypeID, Annotations: yyDollar[2].typeAnnotations}
		}
	case 40:
		yyDollar = yyS[yypt-7 : yypt+1]
		//line thrift.y:331
		{
			yyVAL.fieldType = ast.MapType{KeyType: yyDollar[3].fieldType, ValueType: yyDollar[5].fieldType, Annotations: yyDollar[7].typeAnnotations}
		}
	case 41:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line thrift.y:333
		{
			yyVAL.fieldType = ast.ListType{ValueType: yyDollar[3].fieldType, Annotations: yyDollar[5].typeAnnotations}
		}
	case 42:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line thrift.y:335
		{
			yyVAL.fieldType = ast.SetType{ValueType: yyDollar[3].fieldType, Annotations: yyDollar[5].typeAnnotations}
		}
	case 43:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line thrift.y:337
		{
			yyVAL.fieldType = ast.TypeReference{Name: yyDollar[2].str, Line: yyDollar[1].line}
		}
	case 44:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line thrift.y:341
		{
			yyVAL.baseTypeID = ast.BoolTypeID
		}
	case 45:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line thrift.y:342
		{
			yyVAL.baseTypeID = ast.I8TypeID
		}
	case 46:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line thrift.y:343
		{
			yyVAL.baseTypeID = ast.I8TypeID
		}
	case 47:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line thrift.y:344
		{
			yyVAL.baseTypeID = ast.I16TypeID
		}
	case 48:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line thrift.y:345
		{
			yyVAL.baseTypeID = ast.I32TypeID
		}
	case 49:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line thrift.y:346
		{
			yyVAL.baseTypeID = ast.I64TypeID
		}
	case 50:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line thrift.y:347
		{
			yyVAL.baseTypeID = ast.DoubleTypeID
		}
	case 51:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line thrift.y:348
		{
			yyVAL.baseTypeID = ast.StringTypeID
		}
	case 52:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line thrift.y:349
		{
			yyVAL.baseTypeID = ast.BinaryTypeID
		}
	case 53:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line thrift.y:357
		{
			yyVAL.constantValue = ast.ConstantInteger(yyDollar[1].i64)
		}
	case 54:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line thrift.y:358
		{
			yyVAL.constantValue = ast.ConstantDouble(yyDollar[1].dub)
		}
	case 55:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line thrift.y:359
		{
			yyVAL.constantValue = ast.ConstantBoolean(true)
		}
	case 56:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line thrift.y:360
		{
			yyVAL.constantValue = ast.ConstantBoolean(false)
		}
	case 57:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line thrift.y:361
		{
			yyVAL.constantValue = ast.ConstantString(yyDollar[1].str)
		}
	case 58:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line thrift.y:363
		{
			yyVAL.constantValue = ast.ConstantReference{Name: yyDollar[2].str, Line: yyDollar[1].line}
		}
	case 59:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line thrift.y:365
		{
			yyVAL.constantValue = ast.ConstantList{Items: yyDollar[2].constantValues}
		}
	case 60:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line thrift.y:366
		{
			yyVAL.constantValue = ast.ConstantMap{Items: yyDollar[2].constantMapItems}
		}
	case 61:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line thrift.y:370
		{
			yyVAL.constantValues = nil
		}
	case 62:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line thrift.y:372
		{
			yyVAL.constantValues = append(yyDollar[1].constantValues, yyDollar[2].constantValue)
		}
	case 63:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line thrift.y:376
		{
			yyVAL.constantMapItems = nil
		}
	case 64:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line thrift.y:378
		{
			yyVAL.constantMapItems = append(yyDollar[1].constantMapItems, ast.ConstantMapItem{Key: yyDollar[2].constantValue, Value: yyDollar[4].constantValue})
		}
	case 65:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line thrift.y:386
		{
			yyVAL.typeAnnotations = nil
		}
	case 66:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line thrift.y:387
		{
			yyVAL.typeAnnotations = yyDollar[2].typeAnnotations
		}
	case 67:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line thrift.y:391
		{
			yyVAL.typeAnnotations = nil
		}
	case 68:
		yyDollar = yyS[yypt-6 : yypt+1]
		//line thrift.y:393
		{
			yyVAL.typeAnnotations = append(yyDollar[1].typeAnnotations, &ast.Annotation{Name: yyDollar[3].str, Value: yyDollar[5].str, Line: yyDollar[2].line})
		}
	case 69:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line thrift.y:395
		{
			yyVAL.typeAnnotations = append(yyDollar[1].typeAnnotations, &ast.Annotation{Name: yyDollar[3].str, Line: yyDollar[2].line})
		}
	case 70:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line thrift.y:412
		{
			yyVAL.line = yylex.(*lexer).line
		}
	}
	goto yystack /* stack new state and value */
}
