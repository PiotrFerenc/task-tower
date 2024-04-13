package common

import (
	"fmt"
	"github.com/Jeffail/gabs"
	"github.com/PiotrFerenc/mash2/cmd/worker/actions"
	"github.com/PiotrFerenc/mash2/internal/types"
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
	return "docker"
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
	//
	for key, child := range items {
		// fmt.Printf("key: %v, value: %v\n", key, child.Path("123").String())
		// do akcji
		// dopasowanie parametrow
		// aktualizacja w procesie
		// aktualizacja aktualnego kroku
		fmt.Sprintf("%s %s", key, child.Data())
	}
	return process, nil
}
