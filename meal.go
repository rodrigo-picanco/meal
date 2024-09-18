package main

import (
        "fmt"
        "gopkg.in/yaml.v3"
        "log"
        "errors"
        "os"
)

type Ingredient struct {
        Name string
        Amount int
}

type Recipe struct {
        Name string
        Ingredients []Ingredient
}

type Day struct {
        Day string
        Recipes []string
}

type DSL struct {
        Plan []Day
        Recipes []Recipe
}

func main() {
        data, err := os.ReadFile(os.Args[1])
        if err != nil {
                log.Fatalf("reading error: %v", err)
        }

        result := DSL{}
        err = yaml.Unmarshal([]byte(data), &result)
        if err != nil {
                log.Fatalf("parsing error: %v", err)
        }

        fmt.Println("Shopping List:")
        list := make(map[string]int)
        notFound := make([]string, 0)
        for _, day := range result.Plan {
                for _, recipe := range day.Recipes {
                        recipe, err := findRecipe(result.Recipes, recipe)

                        if err != nil{
                                notFound = append(notFound, recipe.Name)
                                continue
                        }

                        for _, ingredient := range recipe.Ingredients {
                                list[ingredient.Name] += ingredient.Amount
                        }
                }
        }

        for ingredient, amount := range list {
                fmt.Println("  ", ingredient, amount)
        }

        if len(notFound) > 0 {
                fmt.Println("Recipes not found:")
                for _, recipe := range notFound {
                        fmt.Println("  ", recipe)
                }
        }


}

func findRecipe(recipes []Recipe, name string) (*Recipe, error) {
        for _, recipe := range recipes {
                if recipe.Name == name {
                        return &recipe, nil
                }
        }
        return &Recipe{Name: name}, errors.New("Recipe not found: " + name)
}
