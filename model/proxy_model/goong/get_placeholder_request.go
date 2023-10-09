package goong

type GetPlaceHolderRequestDto struct {
	Input  string `json:"input"`
	Limit  int    `json:"limit"`
	Radius int    `json:"radius"`
}
