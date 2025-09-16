package services

// TemperatureService interface para operações de temperatura
type TemperatureService interface {
	ConvertTemperatures(celsius float64) (fahrenheit, kelvin float64)
}

type temperatureService struct{}

// NewTemperatureService cria uma nova instância do serviço de temperatura
func NewTemperatureService() TemperatureService {
	return &temperatureService{}
}

// ConvertTemperatures converte temperatura de Celsius para Fahrenheit e Kelvin
func (s *temperatureService) ConvertTemperatures(celsius float64) (fahrenheit, kelvin float64) {
	return convertTemperatures(celsius)
}

// convertTemperatures converte temperatura de Celsius para Fahrenheit e Kelvin
func convertTemperatures(celsius float64) (fahrenheit, kelvin float64) {
	fahrenheit = convertCelsiusToFahrenheit(celsius)
	kelvin = convertCelsiusToKelvin(celsius)
	return
}

// convertCelsiusToFahrenheit converte Celsius para Fahrenheit
// Fórmula: F = C * 1.8 + 32
func convertCelsiusToFahrenheit(celsius float64) float64 {
	return celsius*1.8 + 32
}

// convertCelsiusToKelvin converte Celsius para Kelvin
// Fórmula: K = C + 273
func convertCelsiusToKelvin(celsius float64) float64 {
	return celsius + 273
}
