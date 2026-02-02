Agent Specification — go-reloaded (FSM Text Processor)

This document defines the behavior, responsibilities, and internal logic of the agent that processes text according to the rules described in prd.md. The agent operates as a deterministic Finite State Machine (FSM).
1. Agent Purpose

The agent is responsible for:

    Reading a stream of tokens

    Tracking context (markers, quotes, punctuation, article rules)

    Applying transformations exactly as defined in the PRD

    Producing a deterministic output sequence of tokens

The agent must never guess, infer meaning, or apply grammar rules beyond those explicitly defined.
2. Core Principles

    Determinism
    Same input → same output, always.

    Single-pass FSM logic
    The agent processes tokens sequentially, transitioning between states.

    Strict rule application
    Only the rules in the PRD are applied.
    No interpretation, no semantic reasoning.

    Graceful handling of malformed input
    Invalid markers, unmatched quotes, or unexpected tokens must not crash the agent.

3. FSM States
Normal

Default state.
The agent reads words, punctuation, or markers.

Transitions:

    "(" → MarkerStart

    "'" → QuoteStart

    punctuation → apply punctuation rule

    word → store token and continue

MarkerStart

Triggered when encountering "(".

Responsibilities:

    Begin reading marker type

    Normalize case (e.g., "(UP)" → "up")

Transitions:

    alphanumeric → MarkerBody

    anything else → treat as normal text

MarkerBody

Reads the marker content:

    hex

    bin

    up

    low

    cap

    optional ", N"

Transitions:

    ")" → MarkerEnd

    invalid structure → ignore marker

MarkerEnd

Marker is complete.

Responsibilities:

    Apply the marker to the appropriate number of previous tokens

    Remove the marker token from output

    Return to Normal

QuoteStart

Triggered by "'".

Responsibilities:

    Begin capturing quoted content

    Do not output the "'" yet

Transition:

    any token → QuoteBody

QuoteBody

Captures all tokens inside quotes.

Responsibilities:

    Store tokens exactly as they appear

    Do not output them yet

    Ignore extra spaces

Transition:

    "'" → QuoteEnd

QuoteEnd

Responsibilities:

    Trim spaces inside quotes

    Output "'content'" as a single unit

    Return to Normal

4. Rule Application Logic
Marker Rules
(hex)

Convert previous token from hex → decimal.
(bin)

Convert previous token from binary → decimal.
(up) / (low) / (cap)

Modify previous token’s case.
(up, N) / (low, N) / (cap, N)

Modify previous N tokens.
If N > available tokens → apply to all.

Malformed markers:
Ignored safely.
Punctuation Rules

    . , ! ? : ; attach to previous token

    A space must follow punctuation

    Groups like ... or !? remain intact

    Multiple punctuation tokens are merged

Quote Rules

    Quotes always appear in pairs

    Remove inner spaces

    Output "'content'" as a single token

    Multi‑word quotes are allowed

Article Rule (a → an)

If the agent sees:

    token = "a"

    next token starts with vowel or "h"

Then replace "a" → "an".

Case-insensitive.
5. Token Processing Order

    Tokenize input

    For each token:

        Determine state

        Apply state logic

        Apply rule if applicable

    Build final output string

6. Error Handling

The agent must:

    Ignore malformed markers

    Ignore unmatched quotes

    Ignore invalid punctuation sequences

    Never panic or crash

    Always produce output

7. Output Construction

The agent reconstructs the final text by:

    Joining tokens with spaces

    Ensuring punctuation spacing rules

    Ensuring quotes are correctly placed

    Ensuring markers are removed after application

8. Example Walkthrough

Input:
this is wild (up, 3) ... right ?

FSM steps:

    Normal: "this"

    Normal: "is"

    Normal: "wild"

    MarkerStart: "("

    MarkerBody: "up, 3"

    MarkerEnd → apply (up, 3)

    Normal: "..." → punctuation group

    Normal: "right"

    Normal: "?" → attach to "right"

Output:
THIS IS WILD... right?
9. Agent Guarantees

    Deterministic output

    No silent failures

    No semantic interpretation

    Full compliance with PRD rules

    FSM-driven processing
    

