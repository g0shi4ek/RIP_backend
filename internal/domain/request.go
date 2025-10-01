package domain

type PhoneRequest struct {
	Phone string `json:"phone" binding:"required"`
}

type ChargingOrderRequest struct {
	BatteryCapacity float32 `json:"battery_capacity" binding:"required"`
	CurrentPercent  int     `json:"current_percent" binding:"required"`
	StartTime       string  `json:"start_time" binding:"required"`
}

type UserLoginRequest struct {
	Login    string `json:"login" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type UserUpdateRequest struct {
	Login    string `json:"login" binding:"required"`
	Password string `json:"password"`
}
