# Paramrep

`paramrep` is a utility for replacing parameter values in URLs and segments path. It can process input from a file or `stdin` and output results to a file or `stdout`.

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
- `-path`: Inject payloads into URL path segments.
- `-pl <payload list>`: File with multiple payloads (mutually exclusive with -p)

### Examples

#### Process URLs with flag `-p`:

```
echo "https://example.com/1/?param1=1&param2=2" | paramrep -p PAYLOAD
```


Contents after execution:

```
https://example.com/1/?param1=PAYLOAD&param2=2 
https://example.com/1/?param1=1&param2=PAYLOAD 
```


#### Process URLs from a file:

Process URLs from payloads the `payloads.txt` file and save the results to `output.txt`:

```
paramrep -pl payloads_list.txt -i input.txt -o output.txt
```

Contents after execution:

```
https://example.com:9443/1/2/?param1=1&param2=PAYLOAD1
https://example.com:9443/1/2/?param1=1&param2=PAYLOAD2
https://example.com:9443/1/2/?param1=PAYLOAD1&param2=2
https://example.com:9443/1/2/?param1=PAYLOAD2&param2=2
```


#### Process URLs with flag `-path`:

```
echo "https://example.com:9443/1/2/?param1=1&param2=2" | ./paramrep -p OLOLO -path
```


Contents after execution:

```
https://example.com:9443/1/2/?param1=1&param2=OLOLO
https://example.com:9443/1/2/?param1=OLOLO&param2=2
https://example.com:9443/OLOLO
https://example.com:9443/1/OLOLO
https://example.com:9443/1OLOLO/2/
https://example.com:9443/1/2OLOLO/
https://example.com:9443/1/2/OLOLO
```



## Help

To display help, use the `-h` flag:

```
paramrep -h
```

## Support

not support )))
