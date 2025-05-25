// touched for cleanup
package model

type GenerationLog struct {
	ID        int
	Timestamp string

	RefImg1URL string
	RefImg2URL string

	Prompt    string
	NegPrompt string
	Seed      int64
	Steps     int
	Guidance  float64

	RefTask1   string
	RefTask2   string
	RefWeight1 float64
	RefWeight2 float64

	Width  int
	Height int
	RefRes int

	TrueCFG           float64
	CFGStartStep      int
	CFGEndStep        int
	NegGuidance       float64
	FirstStepGuidance float64

	GenerationTime float64
	ResponseStatus int
	ErrorMessage   string

	Comment   string
	Effective float64
}
