package handlers

import (
	"fmt"
	"github.com/PiotrFerenc/mash2/cmd/worker/actions"
	"github.com/PiotrFerenc/mash2/web/repositories"
	"github.com/PiotrFerenc/mash2/web/types"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"net/http"
)

func CreateParametersHandler(parametersRepository repositories.ParametersRepository, parameters map[string]actions.Action) func(c echo.Context) error {
	return func(c echo.Context) error {
		idParam := c.Param("id")
		actionName := c.Param("action")
		id, err := uuid.Parse(idParam)
		if err != nil {
			return c.Render(http.StatusBadRequest, "error.html", map[string]interface{}{
				"error": err.Error(),
			})
		}
		values := parametersRepository.GetParameters(id)
		if err != nil {
			return c.Render(http.StatusBadRequest, "error.html", map[string]interface{}{
				"error": err.Error(),
			})
		}
		params, ok := parameters[actionName]
		if !ok {
			return c.Render(http.StatusBadRequest, "error.html", map[string]interface{}{
				"error": fmt.Sprintf("Action %s not found", actionName),
			})
		}

		inputs := mapPropertiesToInputs(params.Inputs(), values)
		output := mapPropertiesToOutputs(params.Outputs())
		data := map[string]interface{}{
			"inputs":  inputs,
			"outputs": output,
		}
		return c.Render(http.StatusOK, "action-form.html", data)
	}
}

func mapPropertiesToInputs(properties []actions.Property, values []types.Parameters) []types.Input {
	var inputs []types.Input
	for _, property := range properties {
		value := getParameterValue(values, property.Name)
		input := types.Input{
			Name:        property.Name,
			Type:        property.Type,
			Description: property.Description,
			Validation:  property.Validation,
			Value:       value.Value,
			Id:          value.ID,
		}
		inputs = append(inputs, input)
	}
	return inputs
}
func mapPropertiesToOutputs(properties []actions.Property) []types.Input {
	var inputs []types.Input
	for _, property := range properties {
		input := types.Input{
			Name:        property.Name,
			Type:        property.Type,
			Description: property.Description,
			Validation:  property.Validation,
		}
		inputs = append(inputs, input)
	}
	return inputs
}
func getParameterValue(values []types.Parameters, name string) types.Parameters {
	for _, value := range values {
		if value.Key == name {
			return value
		}
	}
	return types.Parameters{}
}
