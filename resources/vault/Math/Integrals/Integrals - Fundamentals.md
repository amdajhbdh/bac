---
tags: [math, calculus, integrals]
difficulty: medium
status: active
---

# Integrals - Fundamentals

## Definition
Integration is the reverse process of differentiation. It finds the area under a curve.

## Indefinite Integral
∫ f(x) dx = F(x) + C

Where F'(x) = f(x) and C is the constant of integration

## Definite Integral
∫[a to b] f(x) dx = F(b) - F(a)

Represents the net area between x=a and x=b

## Basic Rules

### Power Rule
∫ x^n dx = (x^n⁺¹)/(n+1) + C  (n ≠ -1)

### Constant Multiple
∫ k·f(x) dx = k·∫ f(x) dx

### Sum/Difference
∫ [f(x) ± g(x)] dx = ∫ f(x) dx ± ∫ g(x) dx

## Common Integrals

| Function | Integral |
|----------|----------|
| ∫ 1 dx | x + C |
| ∫ x^n dx | x^n⁺¹/(n+1) + C |
| ∫ 1/x dx | ln\|x\| + C |
| ∫ e^x dx | e^x + C |
| ∫ a^x dx | a^x/ln(a) + C |
| ∫ sin(x) dx | -cos(x) + C |
| ∫ cos(x) dx | sin(x) + C |
| ∫ sec^2(x) dx | tan(x) + C |

## Fundamental Theorem of Calculus

**Part 1**: If F(x) = ∫[a to x] f(t) dt, then F'(x) = f(x)

**Part 2**: ∫[a to b] f(x) dx = F(b) - F(a)

## Examples

### Example 1
∫ (3x^2 + 2x - 5) dx = x^3 + x^2 - 5x + C

### Example 2
∫[0 to 2] x^2 dx = [x^3/3] from 0 to 2 = 8/3 - 0 = 8/3

## Related
- [[Integrals - Techniques]]
- [[Integrals - Substitution]]
- [[Integrals - Practice Easy]]
