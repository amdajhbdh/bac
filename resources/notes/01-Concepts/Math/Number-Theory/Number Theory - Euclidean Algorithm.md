# Number Theory - Euclidean Algorithm

> Efficient method for computing GCD

## Key Concepts
1. Division Algorithm
2. Euclidean Algorithm Steps
3. Extended Euclidean Algorithm
4. Applications to Diophantine Equations

## Mathematical Formulas
- **Division**: $a = bq + r$ where $0 \leq r < b$
- **Euclidean**: $\gcd(a,b) = \gcd(b, r)$ where $r = a \bmod b$

## Algorithm
```
function gcd(a, b):
    while b ≠ 0:
        (a, b) = (b, a mod b)
    return a
```

## Related
- [[Number Theory - GCD and LCM]]
- [[Number Theory - Divisibility]]
- [[Number Theory MOC]]