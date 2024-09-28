# 🧅 Meal

Meal is a dumb and simple Domain-Specific Language (DSL) implemented in Go for creating shopping lists based on daily meal plans. It allows you to define meals for different days of the week and automatically generates a shopping list with the required ingredients.

## Features

- Define meals for each day of the week
- Create reusable meal components
- Automatically generate shopping lists
- Simple and intuitive syntax

## Installation

To use Meal DSL, you need to have Go installed on your system. Then, you can install it using:

```
go get github.com/rodrigo-picanco/meal-dsl
```

## Usage

1. Create a file with your meal plan using the Meal DSL syntax.
2. Run the Meal DSL parser on your file.
3. Get the generated shopping list.

### Syntax

```
day_of_week (
    meal1
    meal2
    ...
)

meal_name (
    ingredient1 quantity unit
    ingredient2 quantity unit
    ...
)
```

- Days of the week: monday, tuesday, wednesday, thursday, friday, saturday, sunday
- Quantities can be integers or decimals
- Units can be g (grams), kg (kilograms), ml (milliliters), l (liters), u (units)

### Example

```
monday (
    toast
    rice_and_beans
)
tuesday (
    toast
    egg_fried_rice
)

rice_and_beans (
    laurel 1u
    rice 100g
    beans 100g
    garlic 10g
)
egg_fried_rice (
    rice 100g
    egg 1u
    ham 100g
    scallion 1u
)
toast (
    bread 1u
    butter 25g
)
```

## Running the Parser

```
meal plan.meal
```


This will generate a shopping list based on your meal plan.

```
Shopping list:
- garlic 10g
- ham 100g
- laurel 1u
- beans 100g
- rice 200g
- egg 1u
- scallion 1u
- bread 2u
- butter 50g
```

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the MIT License - see the LICENSE file for details.
