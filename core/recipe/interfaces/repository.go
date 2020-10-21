package recipeInterfaces

import baseModels "GOSecretProject/core/model/base"

type RecipeRepository interface {
	CreateRecipe(recipe *baseModels.Recipe) (err error)
	GetRecipe(id uint64) (recipe *baseModels.Recipe, err error)
}
