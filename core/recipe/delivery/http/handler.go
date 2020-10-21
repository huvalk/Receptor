package recipeHttp

import (
	"GOSecretProject/core/model/base"
	"GOSecretProject/core/recipe/interfaces"
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/kataras/golog"
	"io/ioutil"
	"net/http"
	"strconv"
)

type recipeHandler struct {
	useCase recipeInterfaces.RecipeUseCase
}

func NewRecipeHandler(useCase recipeInterfaces.RecipeUseCase) *recipeHandler {
	return &recipeHandler{useCase: useCase}
}

func (h *recipeHandler) CreateRecipe(w http.ResponseWriter, r *http.Request) {
	recipeByte, err := ioutil.ReadAll(r.Body)
	if err != nil {
		golog.Error(err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	var recipe baseModels.Recipe
	err = json.Unmarshal(recipeByte, &recipe)
	if err != nil {
		golog.Error(err.Error())
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	err = h.useCase.CreateRecipe(&recipe)
	if err != nil {
		golog.Error(err.Error())
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *recipeHandler) GetRecipe(w http.ResponseWriter, r *http.Request) {
	idString := mux.Vars(r)["id"]
	id, err := strconv.ParseUint(idString, 10, 64)
	if err != nil {
		golog.Error(err.Error())
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	recipe, err := h.useCase.GetRecipe(id)
	if err != nil {
		golog.Error(err.Error())
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	recipeJson, err := json.Marshal(recipe)
	if err != nil {
		golog.Error(err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(recipeJson)
}

func (h *recipeHandler) GetRecipes(w http.ResponseWriter, r *http.Request) {
	authorIdString := mux.Vars(r)["id"]
	authorId, err := strconv.ParseUint(authorIdString, 10, 64)
	if err != nil {
		golog.Error(err.Error())
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	recipes, err := h.useCase.GetRecipes(authorId)
	if err != nil {
		golog.Error(err.Error())
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	recipesJson, err := json.Marshal(recipes)
	if err != nil {
		golog.Error(err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(recipesJson)
}
