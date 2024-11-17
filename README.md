# Paramrep

`paramrep` is a utility for replacing parameter values in URLs. It can process input from a file or `stdin` and output results to a file or `stdout`.

## Installation

To install `paramrep`, run the following command:

```
go install github.com/aleksey-vi/paramrep@latest
```


## Usage

The utility accepts arguments through flags. Here are the key options:

- `-p <payload>` (required): The value to replace parameter values in the URL.
- `-i <input_file>`: Input file with URLs (optional; defaults to `stdin` if not provided).
- `-o <output_file>`: Output file for results (optional; defaults to `stdout` if not provided).

### Examples

#### Process URLs from a file:

Process a list of URLs from the `input.txt` file and save the results to `output.txt`:

```
paramrep -p PAYLOAD -i input.txt -o output.txt
```

#### Process URLs from a pipe line:

```
cat url_list | paramrep -p PAYLOAD
```

#### Process URLs from a string:

```
echo "https://example.com/?param1=1&param2=2" | paramrep -p PAYLOAD
```


- Contents of `input.txt`:<br>
  https://example.com?param1=1&param2=2  <br>
  https://test.com?foo=bar&baz=qux  <br> <br>
- Contents of `output.txt` after execution:<br>
  https://example.com/?param1=PAYLOAD&param2=2  <br>
  https://example.com/?param1=1&param2=PAYLOAD  <br>
  https://test.com/?foo=PAYLOAD&baz=qux  <br>
  https://test.com/?foo=bar&baz=PAYLOAD  <br>

## Help

To display help, use the `-h` flag:

```
paramrep -h
```

## Support

not support )))
