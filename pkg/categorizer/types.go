package categorizer

type Category string

const (
	CategoryGrocery       Category = "grocery"
	CategoryRestaurant    Category = "restaurant"
	CategoryFoodDelivery  Category = "food_delivery"
	CategoryRideHailing   Category = "ride_hailing"
	CategoryTransfer      Category = "transfer"
	CategoryShopping      Category = "shopping"
	CategorySubscription  Category = "subscription"
	CategoryUtilities     Category = "utilities"
	CategoryHealth        Category = "health"
	CategoryEducation     Category = "education"
	CategoryTravel        Category = "travel"
	CategoryEntertainment Category = "entertainment"
	CategorySalary        Category = "salary"
	CategoryFees          Category = "fees"
	CategoryOther         Category = "other"
)

var CategoryList = []Category{
	CategoryGrocery,
	CategoryRestaurant,
	CategoryFoodDelivery,
	CategoryRideHailing,
	CategoryTransfer,
	CategoryShopping,
	CategorySubscription,
	CategoryUtilities,
	CategoryHealth,
	CategoryEducation,
	CategoryTravel,
	CategoryEntertainment,
	CategorySalary,
	CategoryFees,
	CategoryOther,
}

type Direction string

const (
	DirectionIn  Direction = "in"
	DirectionOut Direction = "out"
)

type TransactionRecord struct {
	Original               string    `json:"original"`
	Category               Category  `json:"category"`
	MerchantOrCounterparty string    `json:"merchant_or_counterparty"`
	Direction              Direction `json:"direction"`
	Amount                 float64   `json:"amount"`
	Currency               string    `json:"currency"`
	Confidence             float64   `json:"confidence"`
}
