package utils

func ValidationTranslate(validatedData map[string]string, translate map[string]string) map[string]string {
	for key, _ := range validatedData {
		validatedData[key] = translate[key]
	}

	return validatedData
}
