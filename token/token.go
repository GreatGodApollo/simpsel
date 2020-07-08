package token

type TokenType string

const (
	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"
	IDENT   = "IDENT"

	// Identifiers & Literals
	INT      = "INT"      // #10, #2, #30
	REGISTER = "REGISTER" // $10, $1, $0

	// Opcodes
	LOAD = "LOAD"
	ADD  = "ADD"
	SUB  = "SUB"
	MUL  = "MUL"
	DIV  = "DIV"
	HLT  = "HLT"
	JMP  = "JMP"
	JMPF = "JMPF"
	JMPB = "JMPB"
	EQ   = "EQ"
	NEQ  = "NEQ"
	GT   = "GT"
	LT   = "LT"
	GTE  = "GTE"
	LTE  = "LTE"
	JMPE = "JMPE"
	NOP  = "NOP"
)

type Token struct {
	Type    TokenType
	Literal string
	Line    int
}

var keywords = map[string]TokenType{
	"load": LOAD,
	"add":  ADD,
	"sub":  SUB,
	"mul":  MUL,
	"div":  DIV,
	"hlt":  HLT,
	"jmp":  JMP,
	"jmpf": JMPF,
	"jmpb": JMPB,
	"eq":   EQ,
	"neq":  NEQ,
	"gt":   GT,
	"lt":   LT,
	"gte":  GTE,
	"lte":  LTE,
	"jmpe": JMPE,
	"nop":  NOP,
}

func LookupIdent(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return ILLEGAL
}