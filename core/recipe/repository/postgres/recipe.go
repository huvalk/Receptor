package postgres

import (
	baseModels "GOSecretProject/core/model/base"
	"database/sql"
	"github.com/kataras/golog"
	"github.com/lib/pq"
)

type recipeRepository struct {
	db *sql.DB
}

func NewRecipeRepository(db *sql.DB) *recipeRepository {
	return &recipeRepository{db: db}
}

func (r *recipeRepository) CreateRecipe(recipe *baseModels.Recipe) (err error) {
	var id uint64

	query := `
		INSERT INTO recipe (user_id, title, cooking_time, ingredients, steps)
		VALUES ($1, $2, $3, $4, $5) RETURNING id`
	err = r.db.QueryRow(query, &recipe.Author, &recipe.Title, &recipe.CookingTime,
		pq.Array(recipe.Ingredients), pq.Array(&recipe.Steps)).Scan(&id)
	if err != nil {
		return err
	}

	golog.Infof("Created recipe with id %d", id)
	return nil
}

func (r *recipeRepository) GetRecipe(id uint64) (*baseModels.Recipe, error) {
	var recipe baseModels.Recipe

	query := `
		SELECT re.id, re.user_id, re.title, re.cooking_time, re.ingredients, re.steps,
			COALESCE(SUM(ra.stars)::numeric/COUNT(ra.stars), 0) stars
		FROM recipe re
		LEFT JOIN rating ra ON re.id = ra.recipe_id
		WHERE re.id = $1
		GROUP BY re.id, re.user_id, re.title, re.cooking_time, re.ingredients, re.steps`
	err := r.db.QueryRow(query, id).Scan(&recipe.Id, &recipe.Author, &recipe.Title, &recipe.CookingTime,
		pq.Array(&recipe.Ingredients), pq.Array(&recipe.Steps), &recipe.Rating)
	if err != nil {
		return nil, err
	}

	return &recipe, nil
}

func (r *recipeRepository) GetRecipes(authorId uint64) (recipes []baseModels.Recipe, err error) {
	query := `
		SELECT re.id, re.user_id, re.title, re.cooking_time, re.ingredients, re.steps,
			COALESCE(SUM(ra.stars)::numeric/COUNT(ra.stars), 0) stars
		FROM recipe re
		LEFT JOIN rating ra ON re.id = ra.recipe_id
		WHERE re.user_id = $1
		GROUP BY re.id, re.user_id, re.title, re.cooking_time, re.ingredients, re.steps`
	rows, err := r.db.Query(query, authorId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var recipe baseModels.Recipe

		err = rows.Scan(&recipe.Id, &recipe.Author, &recipe.Title, &recipe.CookingTime,
			pq.Array(&recipe.Ingredients), pq.Array(&recipe.Steps), &recipe.Rating)
		if err != nil {
			return nil, err
		}

		recipes = append(recipes, recipe)
	}

	return recipes, nil
}

func (r *recipeRepository) AddToFavorites(userId, recipeId uint64) (err error) {
	query := "INSERT INTO favorites (user_id, recipe_id) VALUES ($1, $2)"
	_, err = r.db.Exec(query, userId, recipeId)
	if err != nil {
		return err
	}

	return nil
}

func (r *recipeRepository) GetFavorites(userId uint64) (recipes []baseModels.Recipe, err error) {
	query := `
		SELECT re.id, re.user_id, re.title, re.cooking_time, re.ingredients, re.steps,
			COALESCE(SUM(ra.stars)::numeric/COUNT(ra.stars), 0) stars
		FROM favorites f
		LEFT JOIN recipe re ON f.recipe_id = re.id
		LEFT JOIN rating ra ON f.recipe_id = ra.recipe_id
		WHERE f.user_id = $1
		GROUP BY re.id, re.user_id, re.title, re.cooking_time, re.ingredients, re.steps`
	rows, err := r.db.Query(query, userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var recipe baseModels.Recipe

		err = rows.Scan(&recipe.Id, &recipe.Author, &recipe.Title, &recipe.CookingTime,
			pq.Array(&recipe.Ingredients), pq.Array(&recipe.Steps), &recipe.Rating)
		if err != nil {
			return nil, err
		}

		recipes = append(recipes, recipe)
	}

	return recipes, nil
}

func (r *recipeRepository) VoteRecipe(userId, recipeId, stars uint64) (rating float64, err error) {
	query := `
		INSERT INTO rating (user_id, recipe_id, stars) VALUES ($1, $2, $3)
		ON CONFLICT ON CONSTRAINT rating_user_id_recipe_id_key DO UPDATE SET stars = $3`
	_, err = r.db.Exec(query, userId, recipeId, stars)
	if err != nil {
		return 0, err
	}

	query = `
		SELECT COALESCE(SUM(stars)::numeric/COUNT(stars), 0) stars
		FROM rating WHERE recipe_id = $1 GROUP BY recipe_id`
	err = r.db.QueryRow(query, recipeId).Scan(&rating)
	if err != nil {
		return 0, err
	}

	return rating, nil
}

func (r *recipeRepository) FindRecipes(searchString string) (recipes []baseModels.Recipe, err error) {
	query := `
		SELECT re.id, re.user_id, re.title, re.cooking_time, re.ingredients, re.steps,
			COALESCE(SUM(ra.stars)::numeric/COUNT(ra.stars), 0) stars
		FROM recipe re
		LEFT JOIN rating ra ON re.id = ra.recipe_id
		WHERE LOWER(re.title) LIKE LOWER('%' || $1 || '%')
		GROUP BY re.id, re.user_id, re.title, re.cooking_time, re.ingredients, re.steps`
	rows, err := r.db.Query(query, searchString)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var recipe baseModels.Recipe

		err = rows.Scan(&recipe.Id, &recipe.Author, &recipe.Title, &recipe.CookingTime,
			pq.Array(&recipe.Ingredients), pq.Array(&recipe.Steps), &recipe.Rating)
		if err != nil {
			if err == sql.ErrNoRows {
				return recipes, nil
			}
			return nil, err
		}

		recipes = append(recipes, recipe)
	}

	return recipes, nil
}
