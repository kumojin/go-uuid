# go-uuid

A package providing UUID v1, v4 and OrderedUUID with utility interface and function.

## What is an OrderedUUID?

It is a UUID v1 where the timestamp is positioned at the first part of the uuid.
It is useful as a primary key or index in a database like MySQL.
This way, it can be sorted by timestamp, and new elements should be added on top of the index.

## Utility function?

Yes, the UUID struct has a Scan/Value function that follows the `driver.Valuer` / `driver.Scanner` interface.
This allows you to scan your struct with a UUID right from the database!
Moreover, it is implementing the `encoding.TextUnmarshaler` and the `encoding.TextMarshaler` for encoding and decoding features.

### Supported UUID formats for decoding
```
   "6ba7b810-9dad-11d1-80b4-00c04fd430c8"
   "6ba7b8109dad11d180b400c04fd430c8"
```

### Supported UUID text representations
```
   uuid := canonical | hashlike
   plain := canonical | hashlike
   canonical := 4hexoct '-' 2hexoct '-' 2hexoct '-' 6hexoct
   hashlike := 12hexoct
```
