package common

import (
	"github.com/Jeffail/gabs"
	"github.com/PiotrFerenc/mash2/cmd/worker/actions"
	"github.com/PiotrFerenc/mash2/internal/types"
	"log"
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

func (d *forEachLoop) Execute(process types.Pipeline) (types.Pipeline, error) {
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

	for key, child := range items {
		for _, s := range forEachBody.Steps {
			currentIndex = currentIndex + 1
			s.Sequence = currentIndex
			st := types.MapToStep(s)
			process.Steps = append(process.Steps, *st)
			if currentIndex == process.CurrentStep.Sequence+1 {
				process.CurrentStep = *st
			}
		}

		log.Printf("%s %s", key, child.Data())
	}
	for i, s := range process.Steps {
		if s.Sequence > currentIndex {
			s.Sequence = currentIndex + i
		}
	}

	return process, nil
}
