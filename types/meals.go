package types


type Meal struct {
    MealName        string   `json:"mealName"`
    ImageFilePath   string   `json:"imageFilePath"`
    Ingredients     []string `json:"ingredients"`
    Procedure       []string `json:"procedure"`
    FriendlyComments string   `json:"friendlyComments"`
}