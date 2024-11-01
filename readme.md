## Data format for communication between a client and a DB server.

This repository contains the encoding protocol used for serializing and deserializing a recursive data structure in Go. The protocol is designed to efficiently encode different data types, including strings, integers, and nested arrays, into a binary format that can be reliably decoded back into the original data.

### Encoding / Decoding idea

The initial data type (type Data = Array<string | int32 | Data>) contains three different data types. We encode each element as a sequence of bytes, adding a data type byte and a length of the sequence or array if necessary. During decoding, we read a data type byte and based on that information decide the number of bytes to read to form an element.
### Time and Space complexity

Time Complexity: O(E + S)

As we monotonously go through the string from left to right having constant number of copies, the asymptotics of the encoding/decoding function will be O(N), where N is the length of the string. Which may be divided more precisely into the following:
- E: Total number of elements (including nested elements).
- S: Total length of all DataString values.

Space Complexity:

During encoding/decoding we copy the elements of the string, so the space complexity is O(N). 

- O(N) for the data structures created, _where N is the length of the string._
- O(D) for the call stack during recursion, _where _MaximumDepth(D)_: The deepest level of nested arrays._