package common

import (
	"fmt"
	"github.com/Jeffail/gabs"
	"github.com/PiotrFerenc/mash2/cmd/worker/actions"
	"github.com/PiotrFerenc/mash2/internal/types"
	"strings"
)

func CreateForEachLoop() actions.Action {
	return &forEachLoop{
		collection: actions.Property{
			Name:        "collection",
			Type:        actions.Loop,
			Description: "The collection to loop through",
			DisplayName: "Collection",
			Validation:  "required",
		}, collectionKeyName: actions.Property{
			Name:        "collectionKeyName",
			Type:        actions.Loop,
			Description: "The key name of the collection to loop through",
			DisplayName: "Collection Key Name",
			Validation:  "required",
		},
	}

}

type forEachLoop struct {
	collection        actions.Property
	collectionKeyName actions.Property
}

func (d *forEachLoop) GetCategoryName() string {
	return "for-each"
}

func (d *forEachLoop) Inputs() []actions.Property {
	return []actions.Property{
		d.collection,
	}
}
func (d *forEachLoop) Outputs() []actions.Property {
	return []actions.Property{}
}

func (d *forEachLoop) Execute(process types.Process) (types.Process, error) {
	payload, err := d.collection.GetStringFrom(&process)
	if err != nil {
		return process, err
	}
	key, err := d.collectionKeyName.GetStringFrom(&process)
	if err != nil {
		return process, err
	}
	jsonParsed, err := gabs.ParseJSON([]byte(payload))
	if err != nil {
		return process, err
	}

	items, err := jsonParsed.S(key).Children()
	if err != nil {
		return process, err
	}

	forEachBody := process.CurrentStep.ForeachBody
	currentIndex := process.CurrentStep.Sequence

	for _, _ = range items {
		for _, s := range forEachBody.Steps {
			currentIndex = currentIndex + 1
			s.Sequence = currentIndex
			st := types.MapToStep(s)
			for key, value := range forEachBody.Parameters {
				if strings.Contains(key, st.Name) {
					newName := fmt.Sprintf("%s%d", st.Name, st.Sequence)
					keyName := strings.Replace(key, st.Name, newName, 1)
					st.Name = newName
					process.Parameters[keyName] = value
				}
			}
			process.Steps = append(process.Steps, *st)
		}

	}
	for i, s := range process.Steps {
		if s.Sequence > currentIndex {
			s.Sequence = currentIndex + i
		}
	}

	process.CurrentStep = process.Steps[0]

	return process, nil
}
