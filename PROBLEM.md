<!--toc:start-->
- [**Exercise: Polynomials and Geometry in the Complex Plane**](#exercise-polynomials-and-geometry-in-the-complex-plane)
- [**Quick Hints for Solving** (Based on the Sources)](#quick-hints-for-solving-based-on-the-sources)
<!--toc:end-->

Based on the source material provided, here is a **medium-level exercise** on complex numbers, modeled after common exam problems found in the documents (specifically adapted from Exercise 2 in "JoussourMaths7C.pdf").

### **Exercise: Polynomials and Geometry in the Complex Plane**

**Part 1: Solving the Equation**
Consider the polynomial $P$ defined for any complex number $z$ by:
$$P(z) = z^3 - (5 + 7i)z^2 + (-6 + 26i)z + 24 - 24i$$

1. **Imaginary Root:** Show that the equation $P(z) = 0$ admits a **pure imaginary solution** $z_0 = i\alpha$, where $\alpha$ is a real number.
2. **Complete Resolution:**
    - Find the value of $z_0$.
    - Determine the complex numbers $a$ and $b$ such that $P(z) = (z - z_0)(z^2 + az + b)$.
    - Solve the equation $P(z) = 0$ in the set of complex numbers $\mathbb{C}$.

**Part 2: Geometry and Loci**
The complex plane is equipped with a direct orthonormal coordinate system $(O; \vec{u}, \vec{v})$. Let $A$, $B$, and $C$ be the points representing the solutions of $P(z) = 0$, such that their affixes satisfy $|z_A| < |z_C| < |z_B|$.

1. **Triangle Nature:**
    - Give the affixes $z_A$, $z_B$, and $z_C$ of points $A, B$, and $C$.
    - Calculate the ratio $\frac{z_C - z_A}{z_C - z_B}$ and deduce the **nature of the triangle $ABC$**.
2. **Barycenter:** Determine the affix of the point $G$, which is the **barycentre** of the system $\{(O, 5), (A, -3), (C, 4)\}$.
3. **Locus of Points:** Determine and represent the set $\Gamma$ of points $M$ with affix $z$ such that the number $\frac{z - 2}{z - 4i}$ is **pure imaginary**.

***

### **Quick Hints for Solving** (Based on the Sources)

- **For Part 1.2:** Use Horner's method or long division to factor the polynomial. The solutions are **$\{4i, 2, 3 + 3i\}$**.
- **For Part 1.1:** Substitute $z = i\alpha$ into $P(z)$, separate the real and imaginary parts, and solve for $\alpha$. You will find $\alpha = 4$, so **$z_0 = 4i$**.
- **For Part 2.1:** You will find that $z_A = 2$, $z_B = 4i$, and $z_C = 3 + 3i$. The ratio $\frac{z_C - z_A}{z_C - z_B} = i$, which means the triangle is **right-angled isosceles at $C$**.
- **For Part 2.3:** A ratio $\frac{z - z_A}{z - z_B}$ is pure imaginary if the point $M(z)$ belongs to the **circle with diameter $[AB]$**, excluding the point $B$ (where the denominator is zero).
