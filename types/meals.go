package types

type Meal struct {
	MealName         string   `json:"name"`
	ImageFilePath    string   `json:"imageFilePath"`
	Ingredients      []string `json:"ingredients"`
	Directions       []string `json:"directions"`
	PrepTime         string   `json:"prep_time"`
	CookTime         string   `json:"cook_time"`
	Servings         string   `json:"servings"`
	FriendlyComments string   `json:"friendlyComments"`
}
