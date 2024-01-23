package internal

func BasicResponse(message string) map[string]string {
	return map[string]string{
		"message": message,
	}
}

func UserErrorResponse(message string, example string) map[string]string {
	return map[string]string{
		"message": message,
		"example": example,
	}
}
