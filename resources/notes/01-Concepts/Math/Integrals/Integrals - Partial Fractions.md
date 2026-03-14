---
tags: [concept, math, integrals, partial-fractions]
subject: Mathematics
topic: Integrals
difficulty: Medium
---

# Integrals - Partial Fractions

## Introduction

Partial fraction decomposition is a technique used to integrate rational functions (fractions where both numerator and denominator are polynomials). The idea is to break down a complex fraction into simpler fractions that can be integrated more easily.

## Method

To decompose a rational function $\frac{P(x)}{Q(x)}$ where the degree of P(x) is less than the degree of Q(x):

1. **Factor the denominator** Q(x) completely into linear and irreducible quadratic factors
2. **Set up the partial fraction decomposition** with unknown coefficients
3. **Solve for the coefficients** by equating numerators or substituting convenient values
4. **Integrate each term** separately

## Types of Factors

### Linear Factors

For each linear factor $(ax + b)^n$, include terms:
$$\frac{A_1}{ax + b} + \frac{A_2}{(ax + b)^2} + \cdots + \frac{A_n}{(ax + b)^n}$$

### Quadratic Factors

For each irreducible quadratic factor $(ax^2 + bx + c)^m$, include terms:
$$\frac{B_1x + C_1}{ax^2 + bx + c} + \frac{B_2x + C_2}{(ax^2 + bx + c)^2} + \cdots + \frac{B_mx + C_m}{(ax^2 + bx + c)^m}$$

## Examples

### Example 1: Simple Linear Factors

Integrate $\int \frac{3x + 2}{(x + 1)(x + 2)} dx$

Set up decomposition:
$$\frac{3x + 2}{(x + 1)(x + 2)} = \frac{A}{x + 1} + \frac{B}{x + 2}$$

Multiply both sides by $(x + 1)(x + 2)$:
$$3x + 2 = A(x + 2) + B(x + 1)$$

Method 1 - Substitution:
- When $x = -1$: $3(-1) + 2 = A(-1 + 2) \Rightarrow -1 = A$
- When $x = -2$: $3(-2) + 2 = B(-2 + 1) \Rightarrow -4 = -B \Rightarrow B = 4$

Method 2 - Equating coefficients:
$$3x + 2 = Ax + 2A + Bx + B = (A + B)x + (2A + B)$$
- Coefficient of $x$: $3 = A + B$
- Constant term: $2 = 2A + B$

Solving: $A = -1$, $B = 4$

Therefore:
$$\int \frac{3x + 2}{(x + 1)(x + 2)} dx = \int \left(\frac{-1}{x + 1} + \frac{4}{x + 2}\right) dx$$
$$= -\ln|x + 1| + 4\ln|x + 2| + C$$
$$= \ln\left|\frac{(x + 2)^4}{x + 1}\right| + C$$

### Example 2: Repeated Linear Factors

Integrate $\int \frac{x^2 + 1}{(x - 1)^3} dx$

Set up decomposition:
$$\frac{x^2 + 1}{(x - 1)^3} = \frac{A}{x - 1} + \frac{B}{(x - 1)^2} + \frac{C}{(x - 1)^3}$$

Multiply both sides by $(x - 1)^3$:
$$x^2 + 1 = A(x - 1)^2 + B(x - 1) + C$$

Differentiate both sides:
$$2x = 2A(x - 1) + B$$

Differentiate again:
$$2 = 2A \Rightarrow A = 1$$

Substituting back:
- When $x = 1$: $1 + 1 = C \Rightarrow C = 2$
- When $x = 0$: $1 = A - B + C = 1 - B + 2 \Rightarrow B = 2$

Therefore:
$$\int \frac{x^2 + 1}{(x - 1)^3} dx = \int \left(\frac{1}{x - 1} + \frac{2}{(x - 1)^2} + \frac{2}{(x - 1)^3}\right) dx$$
$$= \ln|x - 1| + 2 \cdot \frac{(x - 1)^{-1}}{-1} + 2 \cdot \frac{(x - 1)^{-2}}{-2} + C$$
$$= \ln|x - 1| - \frac{2}{x - 1} - \frac{1}{(x - 1)^2} + C$$

### Example 3: Irreducible Quadratic Factor

Integrate $\int \frac{2x + 3}{x(x^2 + 1)} dx$

Set up decomposition:
$$\frac{2x + 3}{x(x^2 + 1)} = \frac{A}{x} + \frac{Bx + C}{x^2 + 1}$$

Multiply both sides by $x(x^2 + 1)$:
$$2x + 3 = A(x^2 + 1) + (Bx + C)x$$
$$2x + 3 = Ax^2 + A + Bx^2 + Cx$$
$$2x + 3 = (A + B)x^2 + Cx + A$$

Equating coefficients:
- $x^2$ coefficient: $0 = A + B$
- $x$ coefficient: $2 = C$
- Constant term: $3 = A$

Therefore: $A = 3$, $B = -3$, $C = 2$

$$\int \frac{2x + 3}{x(x^2 + 1)} dx = \int \left(\frac{3}{x} + \frac{-3x + 2}{x^2 + 1}\right) dx$$
$$= 3\ln|x| + \int \frac{-3x}{x^2 + 1} dx + \int \frac{2}{x^2 + 1} dx$$
$$= 3\ln|x| - \frac{3}{2}\ln(x^2 + 1) + 2\arctan(x) + C$$

## Integration Results for Standard Forms

After decomposition, the following integrals commonly appear:

1. $\int \frac{1}{x + a} dx = \ln|x + a| + C$

2. $\int \frac{1}{(x + a)^n} dx = \frac{1}{1-n}(x + a)^{1-n} + C$ for $n \neq 1$

3. $\int \frac{bx + c}{x^2 + px + q} dx$ where $x^2 + px + q$ is irreducible:
   - Split as $\int \frac{bx + c}{x^2 + px + q} dx = \frac{b}{2}\int \frac{2x + p}{x^2 + px + q} dx + \int \frac{c - bp/2}{x^2 + px + q} dx$
   - First term integrates to $\frac{b}{2}\ln(x^2 + px + q)$
   - Second term integrates to a multiple of $\arctan$

## Practical Tips

1. **Always check** that the degree of the numerator is less than the degree of the denominator
2. **Factor completely** before setting up partial fractions
3. **Use substitution shortcuts** when possible to find coefficients
4. **Separate logarithmic terms** appropriately when combining results

## Advanced Applications

Partial fractions are essential in:
- Laplace transforms
- Solving differential equations
- Complex analysis
- Numerical integration methods

## Visualizations

![[Integrals - Partial Fractions.excalidraw]]

---
Back to: [[Integrals MOC]] | [[Math MOC]]