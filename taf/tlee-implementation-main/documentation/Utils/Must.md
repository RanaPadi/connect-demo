The public function is a utility function used for handling errors. It checks if the provided error `err` is not `nil`. If `err` is not `nil`, it calls the `Should` function to log the error and then raises a panic with the error message.

### Parameters:

- `err`: An error value to be checked.

### Flow of Execution:
1. The function takes an error `err` as input.

2. It checks if `err` is not `nil`.

3. If `err` is not `nil`, it calls the `Should` function to log the error.

4. Finally, it raises a panic with the error message, causing the program to terminate abruptly.

**Snippet**

```go
func Must(err error) {
	if err != nil {
		Should(err)
		panic(err)
	}
}
