---
tags: [math, linear-algebra, matrices]
difficulty: medium
status: active
---

# Matrices - Basics

## Definition
A matrix is a rectangular array of numbers arranged in rows and columns.

## Notation
A = [a_i_j] where i = row, j = column

Example 2×3 matrix:
```
A = | 1  2  3 |
    | 4  5  6 |
```

## Types of Matrices

### Square Matrix
Number of rows = number of columns (n×n)

### Identity Matrix (I)
Square matrix with 1s on diagonal, 0s elsewhere
```
I_3 = | 1  0  0 |
     | 0  1  0 |
     | 0  0  1 |
```

### Zero Matrix (O)
All elements are 0

### Diagonal Matrix
Non-zero elements only on main diagonal

### Symmetric Matrix
A = A^T (equals its transpose)

### Upper/Lower Triangular
All elements below/above diagonal are 0

## Matrix Dimensions
m × n matrix has m rows and n columns

## Transpose
A^T: Swap rows and columns

If A = | 1  2 |, then A^T = | 1  3 |
       | 3  4 |            | 2  4 |

## Equality
Two matrices are equal if:
- Same dimensions
- Corresponding elements are equal

## Related
- [[Matrices - Operations]]
- [[Matrices - Determinants]]
- [[Matrices - Practice Easy]]
