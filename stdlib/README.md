# Smog Standard Library

The Smog standard library provides common data structures, utilities, and functionality inspired by the Smalltalk-80 standard library. It follows the philosophy that "everything is an object" and provides clean, minimal APIs.

## Structure

The standard library is organized into the following modules:

### Collections (`stdlib/collections/`)

Data structures for organizing and manipulating groups of objects:

- **Set** - Unordered collection of unique elements
  - Operations: add, remove, includes, union, intersection, difference
  - Use when: You need to track unique items without caring about order

- **OrderedCollection** - Growable, ordered collection (like a dynamic array)
  - Operations: add, addFirst, addLast, removeAt, collect, select, reject
  - Use when: You need a flexible list that can grow and shrink

- **Bag** - Unordered collection that tracks element occurrences (multiset)
  - Operations: add, remove, occurrencesOf
  - Use when: You need to count how many times items appear

### Core Utilities (`stdlib/core/`)

Fundamental utilities for common operations:

- **Math** - Mathematical functions and constants
  - Functions: abs, max, min, sqrt, power, gcd, lcm, factorial, fibonacci
  - Constants: pi, e
  - Use when: You need mathematical operations beyond basic arithmetic

- **Stream** - Sequential access to data (ReadStream, WriteStream)
  - Operations: next, nextPut, peek, atEnd, position
  - Use when: Processing data sequentially or building output incrementally

- **StringUtilities** - String manipulation (placeholder for future)
  - Operations: startsWith, endsWith, contains, indexOf
  - Use when: Working with text data

### I/O (`stdlib/io/`)

Input/output operations:

- **HTTP** - HTTP client for web requests
  - Operations: get, post, getStatus, getBody
  - Use when: Making web API calls or fetching web content

### Cryptography (`stdlib/crypto/`)

Security and encryption:

- **AES** - Advanced Encryption Standard (AES-256)
  - Operations: encrypt, decrypt, generateKey
  - Use when: You need to encrypt sensitive data

- **Hash** - Cryptographic hash functions
  - Functions: sha256, sha512, md5
  - Use when: Verifying data integrity or creating fingerprints

- **Base64** - Base64 encoding/decoding
  - Operations: encode, decode
  - Use when: Converting binary data to text

### Compression (`stdlib/compression/`)

Data compression:

- **ZIP** - ZIP compression and archives
  - Operations: compress, decompress
  - Use when: Reducing data size or creating archives

- **GZIP** - GZIP compression
  - Operations: compress, decompress
  - Use when: Compressing single files (common for web content)

## Design Philosophy

The standard library follows these principles:

1. **Smalltalk Inspiration** - APIs and patterns drawn from Smalltalk-80
2. **Message Passing** - All operations use message sending
3. **Clean Code** - Clear, well-documented, and easy to understand
4. **Minimal Dependencies** - Each library is self-contained where possible
5. **Practical** - Focused on real-world use cases

## Usage

To use a standard library class in your Smog program, you would typically load the file and create instances:

```smog
" Example: Using Set "
| fruits |
fruits := Set new.
fruits initialize.
fruits add: 'apple'.
fruits add: 'banana'.
fruits add: 'cherry'.
fruits add: 'apple'.  " Duplicate ignored "

'Number of unique fruits: ' print.
fruits size println.  " Prints: 3 "
```

```smog
" Example: Using OrderedCollection "
| numbers evens |
numbers := OrderedCollection new.
numbers initialize.
numbers add: 1.
numbers add: 2.
numbers add: 3.
numbers add: 4.
numbers add: 5.

evens := numbers select: [ :n | (n - ((n / 2) * 2)) = 0 ].
'Even numbers: ' println.
evens do: [ :n | n println ].
```

```smog
" Example: Using Math "
| math result |
math := Math new.

result := math factorial: 5.
'5! = ' print.
result println.  " Prints: 120 "

result := math sqrt: 16.
'sqrt(16) = ' print.
result println.  " Prints: 4 "

result := math fibonacci: 10.
'fibonacci(10) = ' print.
result println.  " Prints: 55 "
```

## Implementation Notes

### Current Status

The standard library is in **initial implementation** phase. The APIs are designed and documented, but some features require VM primitive support:

- **Fully Implemented**: Set, OrderedCollection, Bag, Math, Stream
- **Requires VM Primitives**: HTTP, AES, Hash, Base64, ZIP, GZIP, StringUtilities

### VM Primitives

Some functionality requires native Go implementations exposed as VM primitives:

- HTTP operations (net/http)
- Cryptographic operations (crypto/aes, crypto/sha256)
- Compression operations (archive/zip, compress/gzip)
- Advanced string operations

These will be added to the VM in future versions.

### Future Enhancements

Planned additions to the standard library:

- File I/O (File, FileStream)
- Network sockets (Socket, ServerSocket)
- JSON parsing and generation
- Regular expressions
- Date and Time classes
- Random number generation
- More collection types (SortedCollection, IdentitySet)
- Process and Thread classes for concurrency

## Examples

See the `examples/stdlib/` directory for complete working examples demonstrating each library component.

## Contributing

When adding to the standard library:

1. Follow Smalltalk naming conventions
2. Document all public methods with comments
3. Provide clear examples in documentation
4. Keep code clean and readable
5. Write tests for new functionality
6. Consider performance implications

## References

- Smalltalk-80: The Language and its Implementation (Blue Book)
- Pharo by Example
- Squeak by Example
- GNU Smalltalk Library Reference

## License

Same as the Smog project license.
