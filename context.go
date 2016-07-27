package failsafe

type FailsafeContext struct {
	Params valueInterface
	Return setterInterface
}

type FailsafeChainContext struct {
	*FailsafeContext
	GlobalParams valueInterface
}

func newFailsafeContext() *FailsafeContext {
	return &FailsafeContext{
		Params: newValue(),
		Return: newValue(),
	}
}

func newFailsafeChainContext() *FailsafeChainContext {
	return &FailsafeChainContext{
		FailsafeContext: &FailsafeContext{
			Params: newValue(),
			Return: newValue(),
		},
		GlobalParams: newValue(),
	}
}
