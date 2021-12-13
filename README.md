Build popular [jdam](https://gitlab.com/michenriksen/jdam) json fuzzer to linux so library

```bash
go build -o jdamso.so -buildmode=c-shared
```

jdam-so library exports function *jdam_fuzz*
parameters:
- config - string json format
- target - json string to fuzz

returns:
- fuzzed json string

Example of using library is in *test.py* file.

Config format you can see in sources of [jdam fuzzer](https://gitlab.com/michenriksen/jdam)
