package vm

var handlerNumToOpcode = map[int64]OpCode{
	0xf1: PREPAY,
}

func wrapWithKeeper(op OpCode, interpreter *EVMInterpreter, scope *ScopeContext) *operation {
	// only CALL for solpatching
	if op != CALL {
		return interpreter.getOperation(op)
	}

	stack := scope.Stack
	addr := stack.Back(1)
	opKey := addr.ToBig().Int64()

	// determine opcode with reserved address for amplification
	opToOverride, ok := handlerNumToOpcode[opKey]
	if !ok {
		return interpreter.getOperation(op)
	}

	// swaping callGasTemp with addr and replace with gas to decrase params
	stack.swap(2)
	stack.pop()

	return interpreter.getOperation(opToOverride)
}
