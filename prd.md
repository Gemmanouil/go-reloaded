# PRD - go-reloaded (FSM Architecture)

---

## 1. Problem Statement

Build a Go command-line tool that reads a text file, applies a defined set of formatting and transformation rules, and writes the transformed result to an output file.

---

## 2. User / Use Case

Primary user: Student/Customer  
Use case: Run a CLI command to transform a text file into a normalized output based on the project rules.

---

## 3. CLI Contract

Command:  
`go run . <inputFile> <outputFile>`

Inputs:
1. `<inputFile>`: path to an existing readable text file.  
2. `<outputFile>`: path where the corrected text will be written.

Output:  
The tool writes the transformed text into `<outputFile>`.

Error handling:
1. Missing or extra arguments  
2. Input file missing/unreadable  
3. Output file cannot be created/written  
4. I/O errors during processing  
5. Markers are case-insensitive; malformed markers are ignored  
6. Any invalid invocation prints usage + exits  

---

## 4. Functional Requirements (Rules)

### 4.1 Marker Conversions

`(hex)` / `(bin)` convert the previous word:

- `FF (hex)` → `255`  
- `2A (hex)` → `42`  
- `1110 (bin)` → `14`  
- `110010 (bin)` → `50`

---

### 4.2 Case Transformations

1. `(up)` / `(low)` / `(cap)` modify the previous word  
   - `hello (up)` → `HELLO`  
   - `LOUD (low)` → `loud`  
   - `robot (cap)` → `Robot`

2. `(up, N)` / `(low, N)` / `(cap, N)` modify the previous N words  
   - `this is wild (up, 3)` → `THIS IS WILD`  
   - `MAKE IT QUIET (low, 2)` → `MAKE it quiet`  
   - `jump high now (cap, 2)` → `jump High Now`

---

### 4.3 Punctuation Formatting

Rules:
1. `. , ! ? : ;` attach to the previous word  
2. Space must follow punctuation  
3. Groups like `...` or `!?` stay grouped  

Examples:
- `Are you sure  ? Yes !` → `Are you sure? Yes!`  
- `This is strange . . . really` → `This is strange... really`  
- `Stop ! ? Please` → `Stop!? Please`

---

### 4.4 Quotes With Apostrophes

Rules:
1. `'` always appears in pairs  
2. Remove spaces inside quotes  

Examples:
- `He called it '   magic   '` → `He called it 'magic'`  
- `They said ' this   place ' is haunted` → `They said 'this place' is haunted`

---

### 4.5 `a` → `an`

Replace `a` with `an` if next word starts with vowel or `h`.

Examples:
- `She found a umbrella` → `She found an umbrella`  
- `It was a honor to meet you` → `It was an honor to meet you`

---

## 5. Non-Goals (Out of Scope)

1. No external libraries — only Go standard library  
2. No grammar correction beyond `a → an`  
3. No typography enhancements  
4. Malformed markers are ignored  
5. No semantic/contextual understanding  
6. No automatic formatting (line wrapping, indentation, etc.)

---

## 6. Acceptance Criteria

### 6.1 Audit Cases

- [ ] Marker conversion  
  Input: `2A (hex)` → `42`

- [ ] Case transformation  
  Input: `this is wild (up, 3)` → `THIS IS WILD`

- [ ] Punctuation formatting  
  Input: `Are you sure  ? Yes !` → `Are you sure? Yes!`

- [ ] Quotes  
  Input: `He called it '   magic   '` → `He called it 'magic'`

- [ ] a → an  
  Input: `She found a umbrella` → `She found an umbrella`

---

### 6.2 Additional Golden Tests

- [ ] Marker with no previous word  
  `(up)` → `(up)`

- [ ] N larger than available words  
  `Hello world (low, 10)` → `hello world`

- [ ] Invalid count  
  `Try this (up, X)` → `Try this`

- [ ] Sequential markers  
  `run fast (up) now (cap)` → `RUN fast Now`

- [ ] Complex punctuation  
  `Wait , stop ! This is odd ... right ?` → `Wait, stop! This is odd... right?`

---

## 7. Implementation Approach (High Level)

### Decision Summary

Chosen architecture: **FSM**

Reason:
- Handles markers, quotes, and punctuation reliably  
- Allows state transitions for complex parsing  
- Supports backtracking when needed  

Tradeoffs:
- More states = more complexity  
- Adding new rules requires updating transitions  

---

### FSM Diagram

+-----------+
|  Normal   |
+-----------+
|   ^
|   | punctuation / quote / EOF
v   |
+-------------+
| MarkerStart |
+-------------+
|
v
+-------------+
| MarkerBody  |
+-------------+
|
v
+-------------+
| MarkerEnd   |
+-------------+
|
v
(apply marker)
|
v
+-----------+
|  Normal   |
+-----------+
|
v
sees quote (')
|
v
+-------------+
| QuoteStart  |
+-------------+
|
v
+-------------+
| QuoteBody   |
+-------------+
|
v
+-------------+
| QuoteEnd    |
+-------------+
|
v
+-----------+
|  Normal   |
+-----------+

---

## 8. Milestones (Next Week Plan)

**Milestone 1 — File I/O**  
Read input, write changed output.

**Milestone 2 — Tokenization & Marker Parsing**  
Detect `(hex)`, `(bin)`, `(up)`, `(low)`, `(cap)`.

**Milestone 3 — Punctuation & Quotes**  
Attach punctuation, fix spacing, trim quotes.

**Milestone 4 — a → an & Combined Rules**  
Apply article rule and test mixed cases.

**Milestone 5 — Error Handling & Malformed Markers**  
Handle invalid markers, missing args, file errors.

---

## 9. Risks / Open Questions

- Malformed markers may require extra validation  
- Large files may slow down multi-pass parsing  
- Should `( up )` with spaces be accepted?  
- Should markers inside quotes be applied or ignored?

