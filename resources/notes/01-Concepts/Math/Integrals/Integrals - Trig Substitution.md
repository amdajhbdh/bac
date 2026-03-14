---
tags: [concept, math, integrals, trig-substitution]
subject: Mathematics
topic: Integrals
difficulty: Hard
---

# Integrals - Trigonometric Substitution

## Introduction

Trigonometric substitution is a powerful technique for integrating functions involving radicals of the form $\sqrt{a^2 - x^2}$, $\sqrt{a^2 + x^2}$, or $\sqrt{x^2 - a^2}$. The method involves substituting a trigonometric function for the variable to simplify the radical expression.

## Standard Forms and Substitutions

### Form 1: $\sqrt{a^2 - x^2}$

When the integrand contains $\sqrt{a^2 - x^2}$, use the substitution:
$$x = a\sin\theta, \quad dx = a\cos\theta \, d\theta$$

This transforms the radical:
$$\sqrt{a^2 - x^2} = \sqrt{a^2 - a^2\sin^2\theta} = \sqrt{a^2(1 - \sin^2\theta)} = \sqrt{a^2\cos^2\theta} = a|\cos\theta|$$

Since we choose $\theta \in [-\frac{\pi}{2}, \frac{\pi}{2}]$, we have $\cos\theta \geq 0$, so $|\cos\theta| = \cos\theta$.

### Form 2: $\sqrt{a^2 + x^2}$

When the integrand contains $\sqrt{a^2 + x^2}$, use the substitution:
$$x = a\tan\theta, \quad dx = a\sec^2\theta \, d\theta$$

This transforms the radical:
$$\sqrt{a^2 + x^2} = \sqrt{a^2 + a^2\tan^2\theta} = \sqrt{a^2(1 + \tan^2\theta)} = \sqrt{a^2\sec^2\theta} = a|\sec\theta|$$

Since we choose $\theta \in (-\frac{\pi}{2}, \frac{\pi}{2})$, we have $\sec\theta > 0$, so $|\sec\theta| = \sec\theta$.

### Form 3: $\sqrt{x^2 - a^2}$

When the integrand contains $\sqrt{x^2 - a^2}$, use the substitution:
$$x = a\sec\theta, \quad dx = a\sec\theta\tan\theta \, d\theta$$

This transforms the radical:
$$\sqrt{x^2 - a^2} = \sqrt{a^2\sec^2\theta - a^2} = \sqrt{a^2(\sec^2\theta - 1)} = \sqrt{a^2\tan^2\theta} = a|\tan\theta|$$

We typically choose $\theta \in [0, \frac{\pi}{2})$ or $\theta \in [\pi, \frac{3\pi}{2})$ where $\tan\theta \geq 0$, so $|\tan\theta| = \tan\theta$.

## Examples

### Example 1: Form 1 ($\sqrt{a^2 - x^2}$)

Evaluate $\displaystyle \int \frac{dx}{\sqrt{9 - x^2}}$

Let $x = 3\sin\theta$, so $dx = 3\cos\theta \, d\theta$

The integral becomes:
$$\int \frac{3\cos\theta \, d\theta}{\sqrt{9 - 9\sin^2\theta}} = \int \frac{3\cos\theta \, d\theta}{\sqrt{9(1 - \sin^2\theta)}} = \int \frac{3\cos\theta \, d\theta}{3\cos\theta} = \int d\theta = \theta + C$$

To convert back to $x$, recall that $\sin\theta = \frac{x}{3}$, so $\theta = \arcsin\left(\frac{x}{3}\right)$

Therefore:
$$\int \frac{dx}{\sqrt{9 - x^2}} = \arcsin\left(\frac{x}{3}\right) + C$$

### Example 2: Form 2 ($\sqrt{a^2 + x^2}$)

Evaluate $\displaystyle \int \frac{dx}{x^2\sqrt{4 + x^2}}$

Let $x = 2\tan\theta$, so $dx = 2\sec^2\theta \, d\theta$

Then $\sqrt{4 + x^2} = \sqrt{4 + 4\tan^2\theta} = 2\sec\theta$

The integral becomes:
$$\int \frac{2\sec^2\theta \, d\theta}{4\tan^2\theta \cdot 2\sec\theta} = \int \frac{\sec^2\theta \, d\theta}{4\tan^2\theta \sec\theta} = \int \frac{\sec\theta \, d\theta}{4\tan^2\theta}$$

Since $\sec\theta = \frac{1}{\cos\theta}$ and $\tan\theta = \frac{\sin\theta}{\cos\theta}$:
$$\int \frac{\sec\theta \, d\theta}{4\tan^2\theta} = \int \frac{\frac{1}{\cos\theta} d\theta}{4\frac{\sin^2\theta}{\cos^2\theta}} = \int \frac{\cos\theta \, d\theta}{4\sin^2\theta} = \frac{1}{4}\int \csc\theta \cot\theta \, d\theta$$

Since $\frac{d}{d\theta}(\csc\theta) = -\csc\theta \cot\theta$:
$$\frac{1}{4}\int \csc\theta \cot\theta \, d\theta = -\frac{1}{4}\csc\theta + C$$

To convert back to $x$, recall that $\tan\theta = \frac{x}{2}$, so:
$$\sin\theta = \frac{x}{\sqrt{4 + x^2}}, \quad \csc\theta = \frac{\sqrt{4 + x^2}}{x}$$

Therefore:
$$\int \frac{dx}{x^2\sqrt{4 + x^2}} = -\frac{\sqrt{4 + x^2}}{4x} + C$$

### Example 3: Form 3 ($\sqrt{x^2 - a^2}$)

Evaluate $\displaystyle \int \frac{\sqrt{x^2 - 1}}{x} dx$

Let $x = \sec\theta$, so $dx = \sec\theta\tan\theta \, d\theta$

Then $\sqrt{x^2 - 1} = \sqrt{\sec^2\theta - 1} = \sqrt{\tan^2\theta} = \tan\theta$ (assuming $\theta \in [0, \frac{\pi}{2})$)

The integral becomes:
$$\int \frac{\tan\theta \cdot \sec\theta\tan\theta \, d\theta}{\sec\theta} = \int \tan^2\theta \, d\theta$$

Using the identity $\tan^2\theta = \sec^2\theta - 1$:
$$\int \tan^2\theta \, d\theta = \int (\sec^2\theta - 1) d\theta = \tan\theta - \theta + C$$

To convert back to $x$, recall that $\sec\theta = x$, so:
$$\tan\theta = \sqrt{\sec^2\theta - 1} = \sqrt{x^2 - 1}, \quad \theta = \arcsec(x)$$

Therefore:
$$\int \frac{\sqrt{x^2 - 1}}{x} dx = \sqrt{x^2 - 1} - \arcsec(x) + C$$

## Advanced Techniques

### Completing the Square

Sometimes, the expression under the radical isn't in standard form. In these cases, complete the square first.

Example: $\int \frac{dx}{\sqrt{-x^2 + 4x + 5}}$

First, complete the square:
$$-x^2 + 4x + 5 = -(x^2 - 4x) + 5 = -(x^2 - 4x + 4 - 4) + 5 = -(x - 2)^2 + 9$$

So the integral becomes:
$$\int \frac{dx}{\sqrt{9 - (x - 2)^2}}$$

Now substitute $u = x - 2$, $du = dx$:
$$\int \frac{du}{\sqrt{9 - u^2}} = \arcsin\left(\frac{u}{3}\right) + C = \arcsin\left(\frac{x - 2}{3}\right) + C$$

### Using Multiple Techniques

Some integrals require combining trigonometric substitution with other techniques like integration by parts.

## Important Trigonometric Identities

When working with trigonometric substitution, these identities are essential:

1. Pythagorean identities: $\sin^2\theta + \cos^2\theta = 1$
2. Tangent identities: $1 + \tan^2\theta = \sec^2\theta$
3. Secant relationships: $\sec^2\theta - 1 = \tan^2\theta$

## Common Mistakes

1. Forgetting to substitute $dx$ in terms of $d\theta$
2. Incorrectly converting back to the original variable
3. Choosing the wrong range for $\theta$, leading to sign errors
4. Not simplifying the radical correctly based on the chosen range of $\theta$

## Applications

Trigonometric substitution is especially useful in:

1. **Physics**: Calculating moments of inertia, center of mass
2. **Engineering**: Structural analysis, fluid dynamics calculations
3. **Probability and Statistics**: Normal distribution calculations
4. **Geometry**: Surface areas and volumes of revolution

## Visualizations

![[Integrals - Trig Substitution.excalidraw]]

---
Back to: [[Integrals MOC]] | [[Math MOC]]