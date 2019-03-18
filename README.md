# ledis

## Functionalities:
1. Strings:
   - GET (key): get a string value at key already
   - SET (key, value): set a string value at key

2. Lists: List is an ordered collection (duplicates allowed) of string values
   - LLEN (key): return length of a list 
   - RPUSH (key, []values): append 1 or more values to the list, create list if not exists, return length of list after operation.
   - LPOP (key): remove and return the first item of the list
   - RPOP (key): remove and return the last item of the list
   - LRANGE (key, start, stop): return a range of element from the list (zero-based, inclusive of
start and stop), start and stop are non-negative integers

3. Sets: Set is a unordered collection of unique string values (duplicates not allowed)
   - SADD (key, [values]): Add values to set stored at key.
   - SCARD (key): Return the number of elements of the set stored at key.
   - SMEMBERS (key): Return array of all members of set.
   - SREM (key, [values]): Remove values from set.
   - SINTER ([keys]): Set intersection among all set stored in specified keys. Return array of members of the result set.

4. Data expiration:
   - KEYS: List all available keys.
   - DEL (key): Delete a key.
   - FLUSHDB: Delete all keys.
   - EXPIRE (key, seconds): Set a timeout on a key, seconds is a positive integer. Return the number of seconds if the timeout is set.
   - TTL (key): Queries a key Time to Live in seconds
5. Snapshot:
   - SAVE: Save the db into "snapshot.db" file (windows).
   - RESTORE: Restore the db from "snapshot.db".

## Installation
`go install` and run holistic-ledis binary file.
## Usage
The program will expose a `POST` endpoint at `/`, put the commands above into the body to use.
