package types

type Meal struct {
	MealName      string   `json:"name"`
	ImageFilePath string   `json:"image_link"`
	Ingredients   []string `json:"ingredients"`
	Directions    []string `json:"directions"`
	PrepTime      string   `json:"prep_time"`
	CookTime      string   `json:"cook_time"`
	Servings      int      `json:"servings"`
}
