package service

import "strings"

type summerClient struct {
}

func NewSummer() Summer {
	return &summerClient{}
}

func (s *summerClient) Sum(a float64, b float64) float64 {
	return a + b
}

func (s *summerClient) Uppercase(str string) string {
	return strings.ToUpper(str)
}
