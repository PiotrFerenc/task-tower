*DRAFT*
---

**Opis**

<p align="center">
  <img src="https://github.com/PiotrFerenc/mash2/assets/30370747/0d288f65-cb91-4770-88bc-2329fd9d52bb" alt="logo" width="200"/>
</p>
<div style="text-align: justify;">

[![Go](https://github.com/PiotrFerenc/mash2/actions/workflows/go.yml/badge.svg?branch=main)](https://github.com/PiotrFerenc/mash2/actions/workflows/go.yml)

[![Quality Gate Status](https://sonarcloud.io/api/project_badges/measure?project=PiotrFerenc_mash2&metric=alert_status)](https://sonarcloud.io/summary/new_code?id=PiotrFerenc_mash2) [![Bugs](https://sonarcloud.io/api/project_badges/measure?project=PiotrFerenc_mash2&metric=bugs)](https://sonarcloud.io/summary/new_code?id=PiotrFerenc_mash2) [![Code Smells](https://sonarcloud.io/api/project_badges/measure?project=PiotrFerenc_mash2&metric=code_smells)](https://sonarcloud.io/summary/new_code?id=PiotrFerenc_mash2) [![Coverage](https://sonarcloud.io/api/project_badges/measure?project=PiotrFerenc_mash2&metric=coverage)](https://sonarcloud.io/summary/new_code?id=PiotrFerenc_mash2) [![Duplicated Lines (%)](https://sonarcloud.io/api/project_badges/measure?project=PiotrFerenc_mash2&metric=duplicated_lines_density)](https://sonarcloud.io/summary/new_code?id=PiotrFerenc_mash2) [![Lines of Code](https://sonarcloud.io/api/project_badges/measure?project=PiotrFerenc_mash2&metric=ncloc)](https://sonarcloud.io/summary/new_code?id=PiotrFerenc_mash2) [![Reliability Rating](https://sonarcloud.io/api/project_badges/measure?project=PiotrFerenc_mash2&metric=reliability_rating)](https://sonarcloud.io/summary/new_code?id=PiotrFerenc_mash2) [![Security Rating](https://sonarcloud.io/api/project_badges/measure?project=PiotrFerenc_mash2&metric=security_rating)](https://sonarcloud.io/summary/new_code?id=PiotrFerenc_mash2)  [![Technical Debt](https://sonarcloud.io/api/project_badges/measure?project=PiotrFerenc_mash2&metric=sqale_index)](https://sonarcloud.io/summary/new_code?id=PiotrFerenc_mash2) [![Maintainability Rating](https://sonarcloud.io/api/project_badges/measure?project=PiotrFerenc_mash2&metric=sqale_rating)](https://sonarcloud.io/summary/new_code?id=PiotrFerenc_mash2) [![Vulnerabilities](https://sonarcloud.io/api/project_badges/measure?project=PiotrFerenc_mash2&metric=vulnerabilities)](https://sonarcloud.io/summary/new_code?id=PiotrFerenc_mash2)

![100 - commitów](https://img.shields.io/badge/100-commitów-2ea44f?logo=go)

**Distributed Work Automation System (DWAS)** is a system designed to automate and streamline processes in distributed
work environments, integrating various tools and platforms to coordinate tasks across multiple locations.
</div>
<p align="center">
  <img src="https://github.com/PiotrFerenc/mash2/assets/30370747/7e4f24c1-1a14-4840-a7af-1713b6c958d2" alt="logo" width="200"/>

  <img src="https://github.com/PiotrFerenc/mash2/assets/30370747/105cd7c5-bccb-435c-aa4d-fb33930ab2f8" alt="logo" width="200"/>
</p>



**Build**
-----------------------

- clone repo

```git
git clone https://github.com/PiotrFerenc/mash2
```

- impl interface Action 

```go
package actions

import (
	"github.com/PiotrFerenc/mash2/api/types"
)

type addnumbers struct {
}

// CreateAddNumbers This is a function that initializes an instance of the addnumbers struct.
// It returns a pointer to the addnumbers instance.
// This is useful when we don't want to pass the struct by value in subsequent calls.
func CreateAddNumbers() Action {
	return &addnumbers{}
}

// Inputs The Inputs() method returns a slice of Property structure.
// The Property structure includes two fields: Name and Type, both of which are strings.
// These property structures are created for two inputs, 'a' and 'b', of 'number' type.
// It then returns these properties.
func (action *addnumbers) Inputs() []Property {
	output := make([]Property, 2)
	output[0] = Property{
		Name: "a",
		Type: "number",
	}
	output[1] = Property{
		Name: "b",
		Type: "number",
	}
	return output
}

// Outputs The Outputs() method returns a slice of Property structure.
// It creates a property structure for an output, 'c', of 'number' type and returns it.
func (action *addnumbers) Outputs() []Property {
	output := make([]Property, 1)
	output[0] = Property{
		Name: "c",
		Type: "number",
	}
	return output
}

// Execute The Execute() method receives a parameter of types.Message type and returns (types.Message, error).
func (action *addnumbers) Execute(message types.Message) (types.Message, error) {

	// In this Execute() method, first, it retrieves the integer values 'a' and 'b' from the message.
	a, _ := message.GetInt("a")
	b, _ := message.GetInt("b")

	//Then it adds them together and sets the resulting 'c' back into the message.
	c := a + b

	//After performing these operations, it returns the updated message and nil for the error value.
	_, _ = message.SetInt("c", c)

	return message, nil

}

```

register: 

https://github.com/PiotrFerenc/mash2/blob/main/internal/executor/map-executor.go#L20

