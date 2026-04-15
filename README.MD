# ASCII-Art-Color

An extension of the ASCII-Art project written in Go that renders strings as large ASCII graphic representations, with the added ability to colorize specific substrings or the entire output using ANSI color codes.

---

## Table of Contents

- [Overview](#overview)
- [Project Structure](#project-structure)
- [How It Works](#how-it-works)
- [ANSI Colors](#ansi-colors)
- [Installation](#installation)
- [Usage](#usage)
- [Examples](#examples)
- [Edge Cases](#edge-cases)
- [Running Tests](#running-tests)
- [Go Concepts Used](#go-concepts-used)
- [Allowed Packages](#allowed-packages)

---

## Overview

ASCII-Art-Color takes a string as input and outputs a graphic representation using ASCII characters. It extends the base ASCII-Art project with a `--color` flag that lets you colorize any substring within the rendered output.

```
$ go run . --color=red kit "a king kitten have kit"
```

The above command renders the full string in ASCII art, with every occurrence of `kit` — including inside `kitten` — displayed in red.

---

## Project Structure

```
ascii-art-color/
├── go.mod                      # Go module definition
├── main.go                     # Entry point — flag validation, arg parsing, pipeline
├── banner/
│   ├── loader.go               # Reads banner file → map[rune][]string
│   └── loader_test.go          # Tests for banner loading
├── color/
│   ├── color.go                # ANSI color lookup + colored index finder
│   └── color_test.go           # Tests for color lookup and index finding
├── render/
│   ├── render.go               # Renders input with optional color support
│   └── render_test.go          # Tests for rendering logic
└── banners/
    ├── standard.txt            # Standard banner style
    ├── shadow.txt              # Shadow banner style
    └── thinkertoy.txt          # Thinkertoy banner style
```

### Why this structure?

Each package has exactly one responsibility:

- **`main.go`** — thin manager. Validates the flag format, parses all arguments, loads the banner, and prints the result.
- **`banner/`** — knows everything about reading and parsing banner files. Completely unchanged from the base project.
- **`color/`** — knows everything about colors. Maps color names to ANSI codes and finds which character indices need to be colored.
- **`render/`** — knows everything about rendering. Takes the input, banner map, color code and substring, and produces the final colored ASCII art string.

---

## How It Works

The program runs through a clear pipeline from execution to output:

### Stage 1 — Flag Format Validation (`main.go`)

Before any parsing happens, the program loops through `os.Args` and checks that any argument starting with `-` must start with exactly `--color=`. Anything else immediately triggers the usage message.

```
--color=red    ✅ valid
-color=red     ❌ invalid → usage message
--color red    ❌ invalid → usage message
--colour=red   ❌ invalid → usage message
```

### Stage 2 — Argument Parsing (`main.go`)

The `flag` package handles the `--color=` flag cleanly. After `flag.Parse()`, the remaining non-flag arguments are retrieved with `flag.Args()`.

The program maps argument counts to variables:

```
remaining has 1 element:
→ input=remaining[0], substring="", banner="standard"

remaining has 2 elements:
→ if remaining[1] is a valid banner name:
      input=remaining[0], banner=remaining[1]
  else:
      substring=remaining[0], input=remaining[1]

remaining has 3 elements:
→ substring=remaining[0], input=remaining[1], banner=remaining[2]

0 elements or 4+ elements:
→ usage message + exit(1)
```

### Stage 3 — Color Validation (`main.go` + `color/color.go`)

If a color name was provided, `color.GetAnsiCode()` looks it up in the ANSI map. If it returns an empty string, the color name is invalid and the usage message is printed.

### Stage 4 — Banner Loading (`banner/loader.go`)

Identical to the base project. Reads the banner file line by line using `bufio.Scanner`, builds and returns a `map[rune][]string` where each rune maps to its 8 lines of ASCII art.

Each character's starting line in the file is calculated as:
```
startLine = (c - 32) * 9 + 1
```

### Stage 5 — Rendering With Color (`render/render.go`)

This is the core of the new functionality:

1. Split input on `\n` into segments
2. Remove trailing empty segment
3. For each segment:
   - Call `color.FindColoredIndices(segment, substring)` → `map[int]bool` of which positions need coloring
   - Convert segment to `[]rune` for safe index-based access
   - For each of 8 rows, loop over every rune:
     - **Entering a colored run** (current colored, previous not) → prepend ANSI color code
     - **Always** → append the character's art line for this row
     - **Exiting a colored run** (current colored, next not or last) → append reset code `\033[0m`
   - Append completed line + newline to result

### Stage 6 — Output (`main.go`)

`fmt.Print(result)` prints the final string. `fmt.Print` not `fmt.Println` because the render function already appends newlines.

---

## ANSI Colors

The following color names are supported with the `--color` flag:

| Color Name | ANSI Code | Effect |
|---|---|---|
| `black` | `\033[30m` | Black text |
| `red` | `\033[31m` | Red text |
| `green` | `\033[32m` | Green text |
| `yellow` | `\033[33m` | Yellow text |
| `blue` | `\033[34m` | Blue text |
| `magenta` | `\033[35m` | Magenta text |
| `cyan` | `\033[36m` | Cyan text |
| `white` | `\033[37m` | White text |

Any other color name will trigger the usage message and exit with code 1.

---

## Installation

**Prerequisites:** Go 1.18 or higher

```bash
# Clone the repository
git clone https://zone01.git/ascii-art-color
cd ascii-art-color

# Initialize the module
go mod init ascii-art-color
```

Make sure the `banners/` folder contains `standard.txt`, `shadow.txt` and `thinkertoy.txt`.

---

## Usage

```bash
# No color — works exactly like base ascii-art
go run . "<input string>"
go run . "<input string>" "<banner>"

# With color, no substring — colors entire string
go run . --color=<color> "<input string>"
go run . --color=<color> "<input string>" "<banner>"

# With color and substring — colors every occurrence of substring
go run . --color=<color> <substring> "<input string>"
go run . --color=<color> <substring> "<input string>" "<banner>"
```

| Argument | Required | Description |
|---|---|---|
| `--color=<color>` | No | Color name from supported list |
| `<substring>` | No | Text to colorize within the input |
| `<input string>` | Yes | The text to render. Use `\n` for new lines |
| `<banner>` | No | `standard`, `shadow`, or `thinkertoy`. Defaults to `standard` |

---

## Examples

**No flag — works like base project:**
```bash
go run . "Hello"
go run . "Hello" shadow
```

**Color entire string:**
```bash
go run . --color=red "Hello"
```
Renders `Hello` in ASCII art, entire output in red.

**Color a substring:**
```bash
go run . --color=red kit "a king kitten have kit"
```
Renders the full string. Every occurrence of `kit` — including inside `kitten` — is colored red.

**Color with banner:**
```bash
go run . --color=green Hello "Hello World" shadow
```
Renders using the shadow banner with `Hello` colored green.

**With newlines:**
```bash
go run . --color=blue He "Hello\nThere"
```
Renders two lines. `He` is colored in each line where it appears.

**Empty input:**
```bash
go run . ""
# no output
```

**Invalid flag format:**
```bash
go run . -color=red "Hello"
# Usage: go run . [OPTION] [STRING]
# EX: go run . --color=<color> <substring to be colored> "something"
```

---

## Edge Cases

| Case | Behaviour |
|---|---|
| `""` empty input | Prints nothing, exits with code 0 |
| No `--color` flag | Renders normally like base project |
| `--color=red` no substring | Colors entire string |
| Substring not found in input | Renders normally, no color applied |
| Substring spans `\n` | Not found in any single segment, renders without color |
| Invalid color name | Usage message to stderr, exits with code 1 |
| Wrong flag format | Usage message to stderr, exits with code 1 |
| Too many or too few arguments | Usage message to stderr, exits with code 1 |
| `kit` inside `kitten` | Correctly colored as part of the longer word |
| Multiple occurrences | All occurrences colored independently |
| Consecutive colored characters | Wrapped in a single color run, not individually |

---

## Running Tests

```bash
# Run all tests
go test ./...

# Run with verbose output
go test ./... -v

# Run specific package tests
go test ascii-art-color/color
go test ascii-art-color/render
go test ascii-art-color/banner
```

Expected output:
```
?       ascii-art-color         [no test files]
ok      ascii-art-color/banner
ok      ascii-art-color/color
ok      ascii-art-color/render
```

### What is tested

**`banner/loader_test.go`**
- Map contains exactly 95 entries
- Each character has exactly 8 lines

**`color/color_test.go`**
- `GetAnsiCode` returns correct ANSI code for valid color names
- `GetAnsiCode` returns empty string for invalid color names
- `FindColoredIndices` marks all indices when substring is empty
- `FindColoredIndices` correctly finds all occurrences of a substring
- `FindColoredIndices` returns empty map when substring not found

**`render/render_test.go`**
- Empty input returns empty string
- Normal render without color produces correct output
- Render with color wraps correct characters with ANSI codes and reset code

---

## Go Concepts Used

| Concept | Where Used |
|---|---|
| `flag` package | Clean CLI flag parsing in `main.go` |
| `os.Args` | Manual format validation before `flag.Parse()` |
| `strings.HasPrefix` | Detecting flag format and banner names |
| `strings.Index` | Finding substring occurrences efficiently |
| `map[rune][]string` | Banner lookup table |
| `map[int]bool` | Set of colored character indices |
| `[]rune` conversion | Safe index-based character access |
| ANSI escape codes | Terminal color output without external packages |
| Colored run detection | Enter/exit logic for wrapping consecutive colored chars |
| `defer` | Auto-closing banner files |
| Error handling | `if err != nil` pattern throughout |
| Table-driven tests | Idiomatic Go test pattern in all test files |
| Package separation | Each package has one clear responsibility |
| Pointer dereferencing | `*colorFlag` to get value from `flag.String()` |
| `switch` statement | Clean banner name validation |

---

## Allowed Packages

Only standard Go packages are used:

| Package | Purpose |
|---|---|
| `os` | `os.Args`, `os.Open`, `os.Exit`, `os.Stderr` |
| `fmt` | `fmt.Print`, `fmt.Fprintf` for output and errors |
| `flag` | CLI flag parsing |
| `strings` | `strings.HasPrefix`, `strings.Split`, `strings.TrimRight`, `strings.Index` |
| `bufio` | `bufio.NewScanner` for line-by-line file reading |
| `testing` | Unit test files only |

---

## Relationship to Base Project

This project is a direct extension of `ascii-art`. The `banner/loader.go` file is completely unchanged. The new additions are:

- `color/color.go` — new package for ANSI color support
- `render/render.go` — updated to accept and apply color parameters
- `main.go` — updated with flag parsing and color validation

The program remains fully backward compatible — running it without the `--color` flag works exactly like the base `ascii-art` project.

---

## Author

Built as part of the Zone01 Go curriculum.