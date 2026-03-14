# Genetics - Punnett Squares

> Tool for predicting genetic cross outcomes

## Key Concepts
1. Punnett Square Basics
2. Monohybrid Cross
3. Dihybrid Cross
4. Genotypic/Phenotypic Ratios
5. Probability

## Monohybrid Cross

### Example: Flower color
- $P$: $RR$ (red) × $rr$ (white)
- $F_1$: All $Rr$ (red)
- $F_2$: $RR$, $2Rr$, $rr$ = 1:2:1 (genotype), 3:1 (phenotype)

```
       R      R
   ┌──────┬──────┐
r  │  Rr  │  Rr  │
   ├──────┼──────┤
r  │  Rr  │  Rr  │
   └──────┴──────┘
```

## Dihybrid Cross

Example: Seed color (Y/y) and texture (R/r)

```
    YR   Yr   yR   yr
┌─────────────────────┐
│YYRR|YYRr|YyRR|YyRr │
├─────────────────────┤
│YYRr|YYrr|YyRr|Yyrr │
├─────────────────────┤
│YyRR|YyRr|yyRR|yyRr │
├─────────────────────┤
│YyRr|Yyrr|yyRr|yyrr │
└─────────────────────┘
```

## Visualizations
```bash
python /home/med/Documents/bac/assets/manim/Biology/PunnettSquare.py --output 01-Concepts/Natural-Sciences/Genetics/assets/PunnettSquare.png
```

## Related
- [[Genetics - Mendel Laws]]
- [[Genetics - DNA Structure]]
- [[Genetics MOC]]