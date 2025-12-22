package util

import (
	"fmt"
	"math/rand"
)

// GetRandomAvatar returns a random avatar URL from the randomuser.me API
func GetRandomAvatar(index int) string {
	return fmt.Sprintf("https://randomuser.me/api/portraits/lego/%d.jpg", index)
}

// Predefined routes for drivers (used for the gRPC Streaming module)
// (these are San Francisco routes, get these coordinates from Google Maps for example and build a custom route if you want)
var PredefinedRoutes = [][][]float64{
	// Центр → Дніпро
	{
		{49.444820, 32.059680},
		{49.443950, 32.062150},
		{49.442870, 32.064920},
		{49.441920, 32.067410},
	},

	// Центр → Парк / Набережна
	{
		{49.445120, 32.060010},
		{49.446210, 32.063420},
		{49.447080, 32.066880},
		{49.448250, 32.070210},
		{49.449370, 32.072940},
	},

	// Бульвар Шевченка
	{
		{49.443300, 32.055120},
		{49.444110, 32.057430},
		{49.444960, 32.059980},
		{49.445810, 32.062510},
		{49.446640, 32.065090},
	},

	// Зворотний маршрут (для імітації руху назад)
	{
		{49.446640, 32.065090},
		{49.445810, 32.062510},
		{49.444960, 32.059980},
		{49.444110, 32.057430},
		{49.443300, 32.055120},
	},
}

func GenerateRandomPlate() string {
	letters := "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	plate := ""
	for i := 0; i < 3; i++ {
		plate += string(letters[rand.Intn(len(letters))])
	}

	return plate
}
