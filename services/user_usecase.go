package services

func CollPLoginUsecase(username, password string) map[string]interface{} {
	return map[string]interface{}{
		"username": username,
		"password": password,
		"status":   "Login success",
	}
}

func CollPRegisterUsecase(username, password, phone, address string) map[string]interface{} {
	return map[string]interface{}{
		"username": username,
		"password": password,
		"phone":    phone,
		"address":  address,
		"status":   "register success",
	}
}
