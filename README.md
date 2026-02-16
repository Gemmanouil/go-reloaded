# go-reloaded — Documentation Index

This repository contains two core documents required for the Week 1 submission:

---

## 📘 1. prd-template.md  
A complete Product Requirements Document describing:

- Problem statement  
- CLI contract  
- Functional rules (hex/bin, case transforms, punctuation, quotes, a→an)  
- Non-goals  
- Acceptance criteria  
- Golden tests  
- Milestones  
- Risks & open questions  

This file defines **what** the program must do and **how success is measured**.

---

## 🤖 2. agent.md  
A technical document describing:

- The FSM architecture  
- All parser states (Normal, MarkerStart, MarkerBody, MarkerEnd, QuoteStart, QuoteBody, QuoteEnd)  
- State transitions  
- How the agent processes tokens  
- How rules are applied deterministically  
- Error-handling strategy  
- Data flow through the system  

This file defines **how** the program will work internally.

---

These two documents together form the complete specification and architectural plan for the go‑reloaded project.

## Usage & Tests

Run the CLI:

```
go run . <inputFile> <outputFile>
```

Quick check using provided golden tests:

```
go run . golden_test.txt out.txt
diff -u golden_expected.txt out.txt
```

Additional project artifacts (audit, tasks, golden tests, AI usage index) are present in the repository for Phase 2 requirements.
