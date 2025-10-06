package service

// RATING_TO_PERCENT_MODIFICATOR converts rating (0-5) to percentage (0-100)
const RATING_TO_PERCENT_MODIFICATOR = 20

// CalculateCategoryScore calculates the weighted score for a category
// Formula: AvgPercent * CategoryWeight * RATING_TO_PERCENT_MODIFICATOR
func CalculateCategoryScore(avgPercent float64, categoryWeight float64) float64 {
	return avgPercent * categoryWeight * RATING_TO_PERCENT_MODIFICATOR
}
