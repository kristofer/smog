# Smog Standard Library Index

Quick reference guide for all standard library classes and their methods.

## Collections

### Set - Unordered collection of unique elements

**File:** `stdlib/collections/Set.smog`

**Methods:**
- `initialize` - Initialize empty set
- `add: anElement` - Add element (returns false if duplicate or full)
- `remove: anElement` - Remove element (returns true if found)
- `includes: anElement` - Test membership (returns boolean)
- `size` - Number of elements (returns integer)
- `isEmpty` - Test if empty (returns boolean)
- `do: aBlock` - Iterate over elements (returns self)
- `union: anotherSet` - Set union (returns new Set)
- `intersection: anotherSet` - Set intersection (returns new Set)
- `difference: anotherSet` - Set difference (returns new Set)

**Example:**
```smog
| colors |
colors := Set new.
colors initialize.
colors add: 'red'.
colors add: 'blue'.
(colors includes: 'red') println.  " true "
```

---

### OrderedCollection - Growable, ordered collection

**File:** `stdlib/collections/OrderedCollection.smog`

**Methods:**
- `initialize` - Initialize empty collection
- `add: anElement` - Add at end (returns element or nil)
- `addFirst: anElement` - Add at beginning (returns element or nil)
- `addLast: anElement` - Add at end (same as add:)
- `at: index` - Get element (1-based, returns nil if out of bounds)
- `at: index put: value` - Set element (returns value or nil)
- `removeAt: index` - Remove at index (returns element or nil)
- `removeFirst` - Remove first (returns element or nil)
- `removeLast` - Remove last (returns element or nil)
- `first` - Get first element (returns element or nil)
- `last` - Get last element (returns element or nil)
- `size` - Number of elements (returns integer)
- `isEmpty` - Test if empty (returns boolean)
- `do: aBlock` - Iterate over elements (returns self)
- `collect: aBlock` - Transform elements (returns new OrderedCollection)
- `select: aBlock` - Filter elements (returns new OrderedCollection)
- `reject: aBlock` - Inverse filter (returns new OrderedCollection)
- `detect: aBlock` - Find first match (returns element or nil)
- `anySatisfy: aBlock` - Test if any matches (returns boolean)
- `allSatisfy: aBlock` - Test if all match (returns boolean)

**Example:**
```smog
| numbers evens |
numbers := OrderedCollection new.
numbers initialize.
numbers add: 1.
numbers add: 2.
numbers add: 3.
evens := numbers select: [ :n | (n - ((n / 2) * 2)) = 0 ].
```

---

### Bag - Unordered collection with occurrences

**File:** `stdlib/collections/Bag.smog`

**Methods:**
- `initialize` - Initialize empty bag
- `add: anElement` - Add one occurrence (returns boolean)
- `add: anElement withOccurrences: n` - Add n occurrences (returns boolean)
- `remove: anElement` - Remove one occurrence (returns boolean)
- `occurrencesOf: anElement` - Count occurrences (returns integer)
- `includes: anElement` - Test membership (returns boolean)
- `size` - Total elements including duplicates (returns integer)
- `uniqueSize` - Number of unique elements (returns integer)
- `isEmpty` - Test if empty (returns boolean)
- `do: aBlock` - Iterate over unique elements (returns self)
- `doWithOccurrences: aBlock` - Iterate over all occurrences (returns self)

**Example:**
```smog
| wordCount |
wordCount := Bag new.
wordCount initialize.
wordCount add: 'hello'.
wordCount add: 'hello'.
(wordCount occurrencesOf: 'hello') println.  " 2 "
```

---

## Core Utilities

### Math - Mathematical functions

**File:** `stdlib/core/Math.smog`

**Constants:**
- `pi` - π ≈ 3.14159265359
- `e` - Euler's number ≈ 2.71828182846

**Methods:**
- `abs: n` - Absolute value
- `max: a and: b` - Maximum of two numbers
- `min: a and: b` - Minimum of two numbers
- `sqrt: n` - Square root (Newton's method)
- `power: base to: exponent` - Exponentiation
- `floor: n` - Round down (placeholder)
- `ceiling: n` - Round up (placeholder)
- `round: n` - Round to nearest (placeholder)
- `sign: n` - Sign (-1, 0, or 1)
- `isEven: n` - Test if even
- `isOdd: n` - Test if odd
- `gcd: a and: b` - Greatest common divisor
- `lcm: a and: b` - Least common multiple
- `factorial: n` - Factorial (n!)
- `fibonacci: n` - nth Fibonacci number

**Example:**
```smog
| math |
math := Math new.
(math sqrt: 16) println.     " 4 "
(math factorial: 5) println. " 120 "
(math fibonacci: 10) println." 55 "
```

---

### Stream - Sequential data access

**File:** `stdlib/core/Stream.smog`

**Stream (abstract):**
- `initialize` - Initialize stream
- `next` - Read next element
- `next: n` - Read n elements
- `nextPut: element` - Write element
- `atEnd` - Test if at end
- `peek` - Look at next without consuming
- `position` - Get position
- `position: n` - Set position
- `reset` - Reset to beginning
- `skip: n` - Skip n elements

**ReadStream:**
- `on: aCollection` - Create stream on collection
- All Stream methods for reading

**WriteStream:**
- `nextPut: element` - Write element
- `nextPutAll: aCollection` - Write all elements
- `contents` - Get written contents
- All Stream methods for writing

**Example:**
```smog
| stream element |
stream := ReadStream new.
stream on: #(1 2 3 4 5).
element := stream next.
element println.  " 1 "
```

---

## I/O

### HTTP - HTTP client

**File:** `stdlib/io/HTTP.smog`

**Note:** Requires VM primitive support (not yet implemented)

**HTTP:**
- `initialize` - Initialize client
- `get: url` - HTTP GET request (returns body)
- `post: url body: body` - HTTP POST request (returns body)
- `getStatus` - Last response status code
- `getBody` - Last response body

**HTTPRequest:**
- `initialize` - Initialize request
- `url: aUrl` - Set URL
- `method: aMethod` - Set method (GET, POST, etc.)
- `body: aBody` - Set body
- `addHeader: name value: value` - Add header
- `execute` - Execute request (returns HTTPResponse)

**HTTPResponse:**
- `status` - Status code
- `body` - Response body
- `header: name` - Get header value

**Example:**
```smog
| http response |
http := HTTP new.
http initialize.
response := http get: 'http://example.com'.
response println.
```

---

## Cryptography

### AES - Advanced Encryption Standard

**File:** `stdlib/crypto/AES.smog`

**Note:** Requires VM primitive support (not yet implemented)

**Methods:**
- `initialize` - Initialize AES
- `encrypt: data key: key` - Encrypt data (AES-256)
- `decrypt: data key: key` - Decrypt data
- `generateKey` - Generate random 32-byte key
- `isValidKeyLength: key` - Validate key length

**Example:**
```smog
| aes key encrypted decrypted |
aes := AES new.
aes initialize.
key := 'this-is-a-32-byte-secret-key!!'.
encrypted := aes encrypt: 'secret' key: key.
decrypted := aes decrypt: encrypted key: key.
```

---

### Hash - Cryptographic hashes

**File:** `stdlib/crypto/AES.smog`

**Note:** Requires VM primitive support (not yet implemented)

**Methods:**
- `sha256: data` - SHA-256 hash (returns hex string)
- `sha512: data` - SHA-512 hash (returns hex string)
- `md5: data` - MD5 hash (deprecated, returns hex string)

---

### Base64 - Base64 encoding

**File:** `stdlib/crypto/AES.smog`

**Note:** Requires VM primitive support (not yet implemented)

**Methods:**
- `encode: data` - Encode to base64
- `decode: encodedData` - Decode from base64

---

## Compression

### ZIP - ZIP compression

**File:** `stdlib/compression/ZIP.smog`

**Note:** Requires VM primitive support (not yet implemented)

**ZIP:**
- `initialize` - Initialize compressor
- `compress: data` - Compress data
- `decompress: data` - Decompress data

**ZIPArchive:**
- `initialize` - Initialize archive
- `addFile: filename content: content` - Add file
- `writeToData` - Write archive to binary data
- `readFromData: data` - Read archive from data
- `fileNames` - List filenames
- `extractFile: filename` - Extract file content
- `extractAll` - Extract all files

**GZIP:**
- `compress: data` - GZIP compress
- `decompress: data` - GZIP decompress

---

## Status Legend

- ✅ **Fully Implemented** - Set, OrderedCollection, Bag, Math, Stream
- ⚠️ **Interface Ready, Needs VM Primitives** - HTTP, AES, Hash, Base64, ZIP, GZIP, StringUtilities

## Loading Standard Library Classes

To use a standard library class, you can either:

1. **Copy the class definition into your program** (current approach)
2. **Wait for module system** (planned for v0.6.0) which will allow:
   ```smog
   " Future syntax when module system is available "
   import: 'stdlib/collections/Set'.
   ```

## Contributing

See `stdlib/README.md` for contribution guidelines.
