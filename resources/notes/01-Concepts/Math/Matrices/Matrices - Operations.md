---
tags: [math, linear-algebra, matrices, operations]
difficulty: medium
---
> [!summary] 📊 Note Summary
> 
> | Property | Value |
> |----------|-------|
> | **Difficulty** | `easy` #difficulty/easy |
> | **Formulas** | 0 |
> | **Concepts** | 0 |
> | **Related Notes** | 10 |
> | **Word Count** | 1005 |
> | **Last Enhanced** | 2026-03-10 |



## 📊 Note Summary

| Property | Value |
|----------|-------|
| Difficulty | Easy |
| Formulas | 0 |
| Concepts | 0 |
| Related Notes | 10 |
| Word Count | 918 |
| Last Enhanced | 2026-03-10 |



# Matrix Operations

## Addition and Subtraction

### Requirements
- Matrices must have same dimensions (m×n)
- Add/subtract corresponding elements

### Formula
(A + B)_i_j = A_i_j + B_i_j

### Example
A = | 1  2 |    B = | 5  6 |
    | 3  4 |        | 7  8 |

A + B = | 6   8  |
        | 10  12 |

A - B = | -4  -4 |
        | -4  -4 |

### Properties
- Commutative: A + B = B + A
- Associative: (A + B) + C = A + (B + C)
- Identity: A + O = A (O = zero matrix)
- Inverse: A + (-A) = O

## Scalar Multiplication

### Definition
Multiply every element by scalar k

### Formula
(kA)_i_j = k·A_i_j

### Example
k = 3, A = | 1  2 |
           | 3  4 |

3A = | 3   6  |
     | 9   12 |

### Properties
- k(A + B) = kA + kB
- (k + m)A = kA + mA
- k(mA) = (km)A
- 1·A = A

## Matrix Multiplication

### Requirements
- Number of columns in A = number of rows in B
- If A is m×n and B is n×p, then AB is m×p

### Formula
(AB)_i_j = Σ(k=1 to n) A_i_k·B_k_j

Row i of A * Column j of B

### Example 1: 2×2 matrices
A = | 1  2 |    B = | 5  6 |
    | 3  4 |        | 7  8 |

AB = | 1·5+2·7  1·6+2·8 | = | 19  22 |
     | 3·5+4·7  3·6+4·8 |   | 43  50 |

### Example 2: Different dimensions
A = | 1  2  3 |    B = | 1 |
    | 4  5  6 |        | 2 |
                       | 3 |

AB = | 1·1+2·2+3·3 | = | 14 |
     | 4·1+5·2+6·3 |   | 32 |

### Properties
- **NOT commutative**: AB ≠ BA (usually)
- Associative: (AB)C = A(BC)
- Distributive: A(B + C) = AB + AC
- Identity: AI = IA = A
- (AB)ᵀ = B**TA**T

### Example: Non-commutativity
A = | 1  2 |    B = | 0  1 |
    | 0  1 |        | 1  0 |

AB = | 2  1 |    BA = | 0  1 |
     | 1  0 |         | 1  2 |

AB ≠ BA

## Matrix Powers

### Definition
A**n = A·A·...·A (n times)

Only defined for square matrices

### Example
A = | 1  1 |
    | 0  1 |

A**2 = | 1  2 |
     | 0  1 |

A**3 = | 1  3 |
     | 0  1 |

Pattern: A**n = | 1  n |
              | 0  1 |

## Special Products

### Diagonal Matrix * Matrix
D = | d_1  0   0  |
    | 0   d_2  0  |
    | 0   0   d_3 |

DA: Multiply row i of A by d_i
AD: Multiply column j of A by d_j

### Identity Matrix
I·A = A·I = A

For 2×2: I = | 1  0 |
             | 0  1 |

For 3×3: I = | 1  0  0 |
             | 0  1  0 |
             | 0  0  1 |

## Trace

### Definition
Sum of diagonal elements (square matrices only)

tr(A) = Σ A_i_i

### Example
A = | 1  2  3 |
    | 4  5  6 |
    | 7  8  9 |

tr(A) = 1 + 5 + 9 = 15

### Properties
- tr(A + B) = tr(A) + tr(B)
- tr(kA) = k·tr(A)
- tr(AB) = tr(BA)
- tr(A**T) = tr(A)

## Block Multiplication

For partitioned matrices:
| A  B | | E  F | = | AE+BG  AF+BH |
| C  D | | G  H |   | CE+DG  CF+DH |

## Kronecker Product (⊗)

A ⊗ B = | a_1_1B  a_1_2B  ... |
        | a_2_1B  a_2_2B  ... |
        | ...   ...   ... |

### Example
A = | 1  2 |    B = | 0  5 |
    | 3  4 |        | 6  7 |

A ⊗ B = | 0   5   0   10 |
        | 6   7   12  14 |
        | 0   15  0   20 |
        | 18  21  24  28 |

## Hadamard Product (element-wise)

(A ∘ B)_i_j = A_i_j·B_i_j

### Example
A = | 1  2 |    B = | 5  6 |
    | 3  4 |        | 7  8 |

A ∘ B = | 5   12 |
        | 21  32 |

## Practice Problems

### Problem 1
A = | 2  1 |    B = | 1  0 |
    | 0  3 |        | 2  1 |

Find AB and BA

<details>
<summary>Solution</summary>
AB = | 4  1 |
     | 6  3 |

BA = | 2  1 |
     | 4  5 |
</details>

### Problem 2
A = | 1  2  3 |
    | 4  5  6 |

B = | 1  0 |
    | 0  1 |
    | 1  1 |

Find AB

<details>
<summary>Solution</summary>
AB = | 1·1+2·0+3·1  1·0+2·1+3·1 | = | 4  5 |
     | 4·1+5·0+6·1  4·0+5·1+6·1 |   | 10 11 |
</details>

### Problem 3
Find A**3 where A = | 1  1 |
                   | 0  1 |

<details>
<summary>Solution</summary>
A**2 = | 1  2 |
     | 0  1 |

A**3 = A**2·A = | 1  3 |
            | 0  1 |
</details>

## Related
- [[Matrices - Basics]]
- [[Matrices - Determinants]]
- [[Matrices - Inverse]]


## 🔗 Related Notes

- [[Resource Links.md|Resource Links]]
- [[ANIMATION-SYSTEM-COMPLETE.md|ANIMATION-SYSTEM-COMPLETE]]
- [[Animations/README.md|README]]
- [[Resource Links.md|Resource Links]]
- [[VAULT-COMPLETION-REPORT.md|VAULT-COMPLETION-REPORT]]
- [[ANIMATION-SYSTEM-COMPLETE.md|ANIMATION-SYSTEM-COMPLETE]]
- [[Animations/ANIMATION_SPEC.md|ANIMATION_SPEC]]
- [[01-Concepts/Math/Complex-Numbers/Complex Numbers - Operations.md|Complex Numbers - Operations]]
- [[01-Concepts/Chemistry/Organic/Organic - Hydrocarbons.md|Organic - Hydrocarbons]]
- [[01-Concepts/Math/Complex-Numbers/Complex Numbers - Operations.md|Complex Numbers - Operations]]


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
