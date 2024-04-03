package vm

import (
	"math"
)

// NOTE: Maybe change on Set structure?
var handlerNumToOpcode = map[OpCode]map[OpCode]bool{
	CALL: {
		RUNSDKMSG: true,
	},
	STATICCALL: {
		PPFD: true,
	},
}

func wrapWithKeeper(op OpCode, interpreter *EVMInterpreter, scope *ScopeContext) *operation {
	handlers, ok := handlerNumToOpcode[op]
	// only ALLOWED for solpatching
	if !ok {
		return interpreter.getOperation(op)
	}

	stack := scope.Stack
	addr := stack.Back(1)
	opKey := addr.ToBig().Uint64()

	// skip if it is not an available opcode
	if opKey > math.MaxUint8 {
		return interpreter.getOperation(op)
	}

	opEx := OpCode(opKey)

	// determine opcode with reserved address for amplification
	if ok := handlers[opEx]; !ok {
		return interpreter.getOperation(op)
	}

	// swaping callGasTemp with addr and replace with gas to decrase params
	stack.swap(2)
	stack.pop()

	return interpreter.getOperation(opEx)
}
