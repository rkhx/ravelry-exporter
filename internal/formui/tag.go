package formui

import "strings"

func ParseFyneFormTag(tag string) map[string]string {
	if tag == "" {
		return nil
	}

	result := make(map[string]string)
	parts := strings.Split(tag, ",")
	for _, part := range parts {
		if strings.Contains(part, "=") {
			kv := strings.SplitN(part, "=", 2)
			result[strings.TrimSpace(kv[0])] = strings.TrimSpace(kv[1])
		} else {
			result[strings.TrimSpace(part)] = "true"
		}
	}
	return result
}
