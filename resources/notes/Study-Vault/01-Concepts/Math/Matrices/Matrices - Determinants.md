---
tags: [math, linear-algebra, matrices, determinants]
difficulty: hard
---
> [!summary] 📊 Note Summary
> 
> | Property | Value |
> |----------|-------|
> | **Difficulty** | `easy` #difficulty/easy |
> | **Formulas** | 0 |
> | **Concepts** | 0 |
> | **Related Notes** | 10 |
> | **Word Count** | 966 |
> | **Last Enhanced** | 2026-03-10 |



## 📊 Note Summary

| Property | Value |
|----------|-------|
| Difficulty | Easy |
| Formulas | 0 |
| Concepts | 0 |
| Related Notes | 10 |
| Word Count | 893 |
| Last Enhanced | 2026-03-10 |



# Determinants

## Definition
A scalar value computed from square matrix elements. Denoted det(A) or |A|.

## 2×2 Determinant

### Formula
A = | a  b |
    | c  d |

det(A) = ad - bc

### Example
A = | 3  2 |
    | 1  4 |

det(A) = 3·4 - 2·1 = 12 - 2 = 10

## 3×3 Determinant

### Method 1: Cofactor Expansion (First Row)
A = | a  b  c |
    | d  e  f |
    | g  h  i |

det(A) = a·|e f| - b·|d f| + c·|d e|
           |h i|     |g i|     |g h|

= a(ei - fh) - b(di - fg) + c(dh - eg)

### Example
A = | 1  2  3 |
    | 0  4  5 |
    | 1  0  6 |

det(A) = 1·|4 5| - 2·|0 5| + 3·|0 4|
           |0 6|     |1 6|     |1 0|

= 1(24 - 0) - 2(0 - 5) + 3(0 - 4)
= 24 + 10 - 12 = 22

### Method 2: Rule of Sarrus (3×3 only)
Copy first two columns to the right:

| a  b  c | a  b
| d  e  f | d  e
| g  h  i | g  h

Positive diagonals: aei + bfg + cdh
Negative diagonals: ceg + afh + bdi

det(A) = (aei + bfg + cdh) - (ceg + afh + bdi)

## Properties

### Basic Properties
1. det(I) = 1
2. det(A**T) = det(A)
3. det(AB) = det(A)·det(B)
4. det(kA) = k**n·det(A) for n×n matrix
5. det(A**{-1}) = 1/det(A)

### Row Operations
1. **Swap rows**: det changes sign
2. **Multiply row by k**: det multiplied by k
3. **Add multiple of one row to another**: det unchanged

### Special Cases
- det(A) = 0 ⟺ A is singular (not invertible)
- det(A) ≠ 0 ⟺ A is invertible
- If row/column is all zeros: det(A) = 0
- If two rows/columns identical: det(A) = 0
- Triangular matrix: det = product of diagonal elements

## Cofactor Expansion (Any Row/Column)

### Cofactor
C_i_j = (-1)**(i+j)·M_i_j

Where M_i_j = minor (determinant of submatrix)

### Sign Pattern
| +  -  +  - |
| -  +  -  + |
| +  -  +  - |
| -  +  -  + |

### Expansion Along Row i
det(A) = Σ_j a_i_j·C_i_j

### Expansion Along Column j
det(A) = Σ_i a_i_j·C_i_j

### Example: Expand along column 1
A = | 2  1  3 |
    | 0  4  5 |
    | 0  2  1 |

det(A) = 2·|4 5| - 0·(...) + 0·(...)
           |2 1|

= 2(4 - 10) = 2(-6) = -12

## 4×4 and Higher

### Cofactor Expansion
Choose row/column with most zeros

### Example: 4×4
A = | 1  0  2  0 |
    | 3  1  0  4 |
    | 2  0  1  0 |
    | 0  2  0  1 |

Expand along row 1 (two zeros):
det(A) = 1·C_1_1 + 0 + 2·C_1_3 + 0

Calculate 3×3 determinants for C_1_1 and C_1_3

## Triangular Matrices

### Upper/Lower Triangular
det = product of diagonal elements

### Example
A = | 2  3  1 |
    | 0  4  5 |
    | 0  0  6 |

det(A) = 2·4·6 = 48

## Applications

### Area/Volume
- 2×2 determinant = area of parallelogram
- 3×3 determinant = volume of parallelepiped

### Cramer's Rule
For Ax = b:
x_i = det(A_i)/det(A)

Where A_i = A with column i replaced by b

### Example: Cramer's Rule
2x + y = 5
x + 3y = 8

A = | 2  1 |    det(A) = 6 - 1 = 5
    | 1  3 |

A_1 = | 5  1 |    det(A_1) = 15 - 8 = 7
     | 8  3 |

A_2 = | 2  5 |    det(A_2) = 16 - 5 = 11
     | 1  8 |

x = 7/5, y = 11/5

## Eigenvalues

det(A - λI) = 0 gives characteristic equation

### Example
A = | 2  1 |
    | 1  2 |

det(A - λI) = det| 2-λ   1  | = (2-λ)**2 - 1
                 | 1    2-λ |

= λ**2 - 4λ + 3 = 0
λ = 1 or λ = 3

## Computational Methods

### LU Decomposition
If A = LU (lower * upper triangular):
det(A) = det(L)·det(U) = product of diagonal elements

### Row Reduction
Reduce to triangular form, tracking:
- Row swaps (change sign)
- Row multiplications (divide by factor)

## Practice Problems

### Problem 1
Find det| 2  3 |
        | 1  4 |

<details>
<summary>Solution</summary>
det = 2·4 - 3·1 = 8 - 3 = 5
</details>

### Problem 2
Find det| 1  2  0 |
        | 3  1  4 |
        | 0  2  1 |

<details>
<summary>Solution</summary>
Expand along row 1:
= 1·|1 4| - 2·|3 4| + 0
    |2 1|     |0 1|
= 1(1-8) - 2(3-0)
= -7 - 6 = -13
</details>

### Problem 3
A = | 3  0  0 |
    | 2  5  0 |
    | 1  4  2 |

<details>
<summary>Solution</summary>
Lower triangular: det = 3·5·2 = 30
</details>

### Problem 4
If det(A) = 5 and det(B) = 3, find det(AB)

<details>
<summary>Solution</summary>
det(AB) = det(A)·det(B) = 5·3 = 15
</details>

## Related
- [[Matrices - Operations]]
- [[Matrices - Inverse]]
- [[Matrices - Eigenvalues]]


## 🔗 Related Notes

- [[Animations/ALL-EXERCISES-COVERED.md|ALL-EXERCISES-COVERED]]
- [[Animations/ANIMATION_SPEC.md|ANIMATION_SPEC]]
- [[00-Meta/QUICK-START.md|QUICK-START]]
- [[Animations/README.md|README]]
- [[Animations/ALL-EXERCISES-COVERED.md|ALL-EXERCISES-COVERED]]
- [[Animations/ANIMATION_SPEC.md|ANIMATION_SPEC]]
- [[01-Concepts/Math/Complex-Numbers/Complex Numbers - Basics.md|Complex Numbers - Basics]]
- [[Animations/ALL-EXERCISES-COVERED.md|ALL-EXERCISES-COVERED]]
- [[Animations/ANIMATION_SPEC.md|ANIMATION_SPEC]]
- [[00-Meta/QUICK-START.md|QUICK-START]]


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
