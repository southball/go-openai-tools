package tools_test

import "strconv"

type DemoToolRequest struct {
	Arg1 string `json:"arg1" required:"true" description:"First argument"`
	Arg2 int    `json:"arg2" description:"Second argument"`
}

type DemoToolResponse struct {
	Result string `json:"result" description:"Result of the tool call"`
}

func DemoTool(r DemoToolRequest) (DemoToolResponse, error) {
	return DemoToolResponse{Result: "Processed: " + r.Arg1 + " and " + strconv.Itoa(r.Arg2)}, nil
}

func ptr[T any](v T) *T {
	return &v
}
