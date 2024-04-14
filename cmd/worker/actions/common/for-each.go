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
			Validation:  "Must be a valid collection",
		}, collectionKeyName: actions.Property{
			Name:        "collectionKeyName",
			Type:        actions.Loop,
			Description: "The key name of the collection to loop through",
			DisplayName: "Collection Key Name",
			Validation:  "Must be a valid collection key name",
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
	items, err := jsonParsed.S(key).ChildrenMap()
	if err != nil {
		return process, err
	}

	forEachBody := process.CurrentStep.ForeachBody
	currentIndex := process.CurrentStep.Sequence

	for key, child := range items {
		for i, s := range forEachBody.Steps {
			s.Sequence = currentIndex + i
			st := types.MapToStep(s)
			process.Steps = append(process.Steps, *st)
		}

		log.Printf("%s %s", key, child.Data())
	}
	for i, s := range process.Steps {
		if s.Sequence > currentIndex {
			s.Sequence = currentIndex + i
		}
	}

	process.CurrentStep = *types.MapToStep(forEachBody.Steps[0])
	return process, nil
}
