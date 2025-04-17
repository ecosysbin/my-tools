package utils

func MergeMaps(m1 map[string]string, m2 map[string]string) map[string]string {
	for k, v := range m2 {
		m1[k] = v
	}
	return m1
}

func ReduceMaps(m1 map[string]string, reduceKeys []string) map[string]string {
	for _, k := range reduceKeys {
		delete(m1, k)
	}
	return m1
}
