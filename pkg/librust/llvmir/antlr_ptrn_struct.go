package llvmir

type LLVMPatternStruct struct {
	ID    string
	IsMut bool
}

func (ps LLVMPatternStruct) String() string {
	return ps.ID
}
