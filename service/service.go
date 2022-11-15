package service

// Service 尝试用反射生成代码
type Service interface {
	Sum(a float64, b float64, c string) (float64, string, error)
}
