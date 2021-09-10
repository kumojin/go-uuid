# go-uuid

A package providing UUID v1, v4 and Ordered UUID with utility interface and fonction

## What is an OrderedUUID ?

It's an UUID v1 where the timestamp part of the uuid is move in the begening of the UUID.
It's mostly useful when you use them as a primary key or an index in a database like MySQL as the index we now be sorted and all new element should be added on top of the index.

## Utility function?

Yes, the UUID struct has the Scan/Value function that follow the `driver.Valuer` / `driver.Scanner` interface. This let you scan your struct with an UUID right from the database!
It also add implement the `encoding.TextUnmarshaler` and the `encoding.TextMarshaler` for encoding/decoding.

Following UUID formats for decoding are supported:
```
   "6ba7b810-9dad-11d1-80b4-00c04fd430c8",
   "6ba7b8109dad11d180b400c04fd430c8"
```

Supported UUID text representation follows:
```
   uuid := canonical | hashlike
   plain := canonical | hashlike
   canonical := 4hexoct '-' 2hexoct '-' 2hexoct '-' 6hexoct
   hashlike := 12hexoct
```