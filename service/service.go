package service

// Service 尝试用反射生成代码
// 类型限制：支持float64、string、int64
// 返回值：支持多返回值，尾返回值必须为error
type Service interface {
	Sum(a float64, b float64) (float64, error)
	Upper(s string) (string, error)
}
