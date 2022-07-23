---
title: `encoding/json.Unmarshal` merges structs but not map values
description: Golang package JSON does not merge map values when decoding.  
---

# Introduction

Let's say we have a struct type with two fields

```go
type foo struct {
	A int `json:"a"`
	B int `json:"b"`
}
```

and we try to unmarshal two JSON payloads `{"a": 1}` and `{"b": 1}` on it. 

[The resulting value would be](https://go.dev/play/p/1ofWE7uAvP3) `{A: 1, B: 1}`: 
i.e. we can say that `json.Unmarshal` would _merge_ the existing value with the one coming from the JSON payload.

However, if our type is `map[string]foo` and we try to unmarshal two JSON payloads:

```json
{
    "key": {"a": 1}
}
```

and

```json
{
    "key": {"b": 1}
}
```

Then [the result would be](https://go.dev/play/p/8fIcGSzkMBM) just `map[key:{A:0 B:1}]`:
i.e., the second JSON payload would completely overwrite the first one, instead of merging them.

# More examples

## Struct with struct field

When the outer type is defined as a specific struct instead of a map:

```go
type foo struct {
	A int `json:"a"`
	B int `json:"b"`
}

type bar struct {
	Key foo `json:"key"`
}
```

Then [the result is properly merged](https://go.dev/play/p/bqZtbm6TMVZ), producing the _expected_ `{Key:{A:1 B:1}}` as the result.

## Different map keys

On the other hand, when different map keys are unmarshaled, like `{"key1": {"a": 1}}` and `{"key2": {"b": 1}}`,
then [the outer map is merged](https://go.dev/play/p/rrtGiAADTlZ) too, resulting in `map[key1:{A:1 B:0} key2:{A:0 B:1}]`.

## Merging slice values

Finally, when two JSON arrays of one element, `[{"a": 1}]` and `[{"b": 1}]` [are unmarshaled on the same slice](https://go.dev/play/p/UuT5JOrG7Rc) `[]foo`,
then the resulting first element of the slice is merged again, producing `[{A:1 B:1}]`.

# Why?

The [issue golang/go#33487](https://github.com/golang/go/issues/33487) has a discussion around it, 
and the main argument is that `json` package does not perform recursive merging of values, like map values.

This doesn't seem to apply to slice values, however, which _are_ merged as was shown before.

# YAML and friends

This issue can also be reproduced with other format unmarshaling libraries, like `gopkg.in/yaml.v3`,
and it obviously depends on the implementation of each one, but they seem to be largerly consistent on this inconsistency.
