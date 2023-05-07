package vm

var handlerNumToOpcode = map[int64]OpCode{
	0xf1: PREPAY,
}

func wrapWithKeeperCall(execCall executionFunc) executionFunc {
	return func(pc *uint64, interpreter *EVMInterpreter, scope *ScopeContext) ([]byte, error) {
		stack := scope.Stack
		// swaping callGasTemp with addr
		stack.swap(2)

		addr := stack.peek()

		opKey := addr.ToBig().Int64()

		opToOverride, ok := handlerNumToOpcode[opKey]
		if !ok {
			// return back position
			stack.swap(2)
			return execCall(pc, interpreter, scope)
		}

		operation := interpreter.getOperation(opToOverride)
		// removing addr as it has not used in future, only for determitation scope
		stack.pop()
		return operation.execute(pc, interpreter, scope)
	}
}
