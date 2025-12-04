package utils

import "fmt"

func GetNested(data map[string]any, keys ...string) any {
	cur := data
	for i, k := range keys {
		v, ok := cur[k]
		if !ok || v == nil {
			return nil
		}

		// 最后一级：把真实值原样返回（float64/string/[]any 等）
		if i == len(keys)-1 {
			return v
		}

		// 中间级：继续深入
		switch next := v.(type) {
		case map[string]any:
			cur = next
		case []any:
			if len(next) > 0 {
				if m, ok := next[0].(map[string]any); ok {
					cur = m
				} else {
					return nil
				}
			} else {
				return nil
			}
		default:
			return nil
		}
	}
	return nil
}
func GetNestedString(data map[string]any, keys ...string) string {
	v := GetNested(data, keys...)
	switch val := v.(type) {
	case string:
		return val
	case []any: // Lark 文本列
		if len(val) > 0 {
			if str, ok := val[0].(string); ok {
				return str
			}
			if m, ok := val[0].(map[string]any); ok {
				return m["text"].(string)
			}
		}
	case float64: // 数值列
		return fmt.Sprintf("%.0f", val)
	}
	return ""
}

// GetNestedFloat64 专门拿数值
func GetNestedFloat64(data map[string]any, keys ...string) float64 {
	v := GetNested(data, keys...)
	switch val := v.(type) {
	case float64:
		return val
	case int:
		return float64(val)
	}
	return 0
}
