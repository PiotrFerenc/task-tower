*DRAFT*
---

**Opis**

<p align="center">
  <img src="https://github.com/PiotrFerenc/mash2/assets/30370747/0d288f65-cb91-4770-88bc-2329fd9d52bb" alt="logo" width="200"/>
</p>
<div style="text-align: justify;">

[![Go](https://github.com/PiotrFerenc/mash2/actions/workflows/go.yml/badge.svg?branch=main)](https://github.com/PiotrFerenc/mash2/actions/workflows/go.yml)

[![Quality Gate Status](https://sonarcloud.io/api/project_badges/measure?project=PiotrFerenc_mash2&metric=alert_status)](https://sonarcloud.io/summary/new_code?id=PiotrFerenc_mash2)

![100 - commitów](https://img.shields.io/badge/100-commitów-2ea44f?logo=go)

# DWAS

------------------------

**D**istributed **W**ork **A**utomation **S**ystem is a system designed to automate and streamline processes in
distributed
work environments, integrating various tools and platforms to coordinate tasks across multiple locations.
</div>

## Installation

----------------------

### requirements

- [MAKE](https://www.gnu.org/software/make/)
- [docker](https://docs.docker.com/engine/install/)
- [docker-compose](https://docs.docker.com/compose/install/)
- [GO](https://go.dev/doc/install)


```git
git clone https://github.com/PiotrFerenc/mash2
```
```shell
cd mash2
```
serve

```makefile
make docker-rebuild
```
or

```shell
docker build -t dwas/controller -f docker/Dockerfile-controller .
docker build -t dwas/worker -f docker/Dockerfile-worker .
docker-compose -f docker/docker-compose.yml up
```

## Usage

---------


This is a cURL command which sends a POST request to the URL "http://localhost:5000/execute". It sends a payload of JSON data where it sets a series of parameters and stages. The parameters and actions in the stages seem to indicate a sequence of operations to be performed.

```shell
curl -X POST --location "http://localhost:5000/execute" \
    -d '{
    "Parameters": {
        "numbers.a" : "1",
        "numbers.b": "2",
        "wynik.text": "{{numbers.a}} + {{numbers.b}} = {{numbers.c}}"
    },
    "Stages": [
        {
            "Order": 1,
            "Name": "numbers",
            "Action": "add-numbers"
        },{
            "Order": 2,
            "Name": "wynik",
            "Action": "console"
        }
    ]
}'
```

`Parameters`- In the "Parameters" section, we can declare variables that will be used in the "Stages" section. 

`numbers.a`- `numbers`- proper name of action, `a`- name of input argument [link](https://github.com/PiotrFerenc/mash2/blob/main/cmd/worker/actions/add-numbers.go#L24)


`Stages`- In the "Stages" section, we define the next steps.

`Order`- means the sequence of performing actions

`Name`- the proper name of the action

`Action` - the name of the action

## Todo

- [ ] User Interface (UI)
- [x] Application Programming Interface (API)
- [ ] Scheduling
- [x] Triggering Tasks
- [x] Handling Distributed Tasks
- [ ] Dependency Management of Tasks
- [ ] Monitoring and Logging
- [x] Error Management
- [ ] Retry Mechanisms
- [ ] Security
- [ ] Integration with External Services and Applications
- [x] Extensibility
- [x] Configuration

## STATISTICS
------------------------
[![Bugs](https://sonarcloud.io/api/project_badges/measure?project=PiotrFerenc_mash2&metric=bugs)](https://sonarcloud.io/summary/new_code?id=PiotrFerenc_mash2) [![Code Smells](https://sonarcloud.io/api/project_badges/measure?project=PiotrFerenc_mash2&metric=code_smells)](https://sonarcloud.io/summary/new_code?id=PiotrFerenc_mash2) [![Coverage](https://sonarcloud.io/api/project_badges/measure?project=PiotrFerenc_mash2&metric=coverage)](https://sonarcloud.io/summary/new_code?id=PiotrFerenc_mash2) [![Duplicated Lines (%)](https://sonarcloud.io/api/project_badges/measure?project=PiotrFerenc_mash2&metric=duplicated_lines_density)](https://sonarcloud.io/summary/new_code?id=PiotrFerenc_mash2) [![Lines of Code](https://sonarcloud.io/api/project_badges/measure?project=PiotrFerenc_mash2&metric=ncloc)](https://sonarcloud.io/summary/new_code?id=PiotrFerenc_mash2) [![Reliability Rating](https://sonarcloud.io/api/project_badges/measure?project=PiotrFerenc_mash2&metric=reliability_rating)](https://sonarcloud.io/summary/new_code?id=PiotrFerenc_mash2) [![Security Rating](https://sonarcloud.io/api/project_badges/measure?project=PiotrFerenc_mash2&metric=security_rating)](https://sonarcloud.io/summary/new_code?id=PiotrFerenc_mash2)  [![Technical Debt](https://sonarcloud.io/api/project_badges/measure?project=PiotrFerenc_mash2&metric=sqale_index)](https://sonarcloud.io/summary/new_code?id=PiotrFerenc_mash2) [![Maintainability Rating](https://sonarcloud.io/api/project_badges/measure?project=PiotrFerenc_mash2&metric=sqale_rating)](https://sonarcloud.io/summary/new_code?id=PiotrFerenc_mash2) [![Vulnerabilities](https://sonarcloud.io/api/project_badges/measure?project=PiotrFerenc_mash2&metric=vulnerabilities)](https://sonarcloud.io/summary/new_code?id=PiotrFerenc_mash2)

## License

--------------

MIT License

Copyright (c) 2024 DWAS

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
