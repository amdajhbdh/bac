---
tags: [math, number-theory, divisibility]
difficulty: easy
status: active
---
> [!summary] 📊 Note Summary
> 
> | Property | Value |
> |----------|-------|
> | **Difficulty** | `easy` #difficulty/easy |
> | **Formulas** | 0 |
> | **Concepts** | 0 |
> | **Related Notes** | 10 |
> | **Word Count** | 316 |
> | **Last Enhanced** | 2026-03-10 |



## 📊 Note Summary

| Property | Value |
|----------|-------|
| Difficulty | Easy |
| Formulas | 0 |
| Concepts | 0 |
| Related Notes | 10 |
| Word Count | 243 |
| Last Enhanced | 2026-03-10 |



# Number Theory - Divisibility

## Definition
Integer a divides integer b (written a|b) if there exists an integer k such that b = ka

Example: 3|12 because 12 = 3×4

## Divisibility Rules

### By 2
Last digit is even (0, 2, 4, 6, 8)

### By 3
Sum of digits is divisible by 3
Example: 123  ->  1+2+3=6, 3|6 

### By 4
Last two digits form a number divisible by 4

### By 5
Last digit is 0 or 5

### By 6
Divisible by both 2 and 3

### By 9
Sum of digits is divisible by 9

### By 10
Last digit is 0

### By 11
Alternating sum of digits is divisible by 11
Example: 1234  ->  1-2+3-4=-2, not divisible by 11

## Properties

### Transitivity
If a|b and b|c, then a|c

### Linear Combination
If a|b and a|c, then a|(bx + cy) for any integers x, y

### Division Algorithm
For any integers a and b (b>0), there exist unique integers q and r such that:
a = bq + r, where 0 ≤ r < b

## Examples

**Example 1**: Does 7|91?
91 = 7×13, so yes, 7|91

**Example 2**: Find remainder when 47 is divided by 5
47 = 5×9 + 2, remainder = 2

## Related
- [[Number Theory - Primes]]
- [[Number Theory - GCD and LCM]]
- [[Number Theory - Practice Easy]]


## 🔗 Related Notes

- [[QUICK-REFERENCE.md|QUICK-REFERENCE]]
- [[Resource Links.md|Resource Links]]
- [[VAULT-COMPLETION-REPORT.md|VAULT-COMPLETION-REPORT]]
- [[QUICK-REFERENCE.md|QUICK-REFERENCE]]
- [[Resource Links.md|Resource Links]]
- [[VAULT-COMPLETION-REPORT.md|VAULT-COMPLETION-REPORT]]
- [[VAULT-COMPLETION-REPORT.md|VAULT-COMPLETION-REPORT]]
- [[00-Meta/QUICK-START.md|QUICK-START]]
- [[00-Meta/MOCs/Math MOC.md|Math MOC]]
- [[Animations/README.md|README]]


> [!related] 🔗 Related Notes
> 
> - [[QUICK-REFERENCE.md|QUICK-REFERENCE]]
> - [[Resource Links.md|Resource Links]]
> - [[ANIMATION-SYSTEM-COMPLETE.md|ANIMATION-SYSTEM-COMPLETE]]
> - [[QUICK-REFERENCE.md|QUICK-REFERENCE]]
> - [[ANIMATION-SYSTEM-COMPLETE.md|ANIMATION-SYSTEM-COMPLETE]]
> - [[Animations/ALL-EXERCISES-COVERED.md|ALL-EXERCISES-COVERED]]
> - [[00-Meta/DEEP-CONTENT-STATUS.md|DEEP-CONTENT-STATUS]]
> - [[00-Meta/MOCs/Chemistry MOC.md|Chemistry MOC]]
> - [[01-Concepts/Math/Complex-Numbers/Complex Numbers - Operations.md|Complex Numbers - Operations]]
> - [[Animations/ANIMATION_SPEC.md|ANIMATION_SPEC]]
