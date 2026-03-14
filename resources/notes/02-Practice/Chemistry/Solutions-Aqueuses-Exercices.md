# Exercices Solutions Aqueuses - Pratique Avancée

*Basé sur la série EL MOUHAMEDIA 2024-2025*

## Exercice 1: Mélanges Acide-Base

### Énoncé
On mélange $V_a = 200\text{mL}$ de solution d'acide chlorhydrique de molarité $C_a = 0,08\text{mol/L}$ et $V_b = 300\text{mL}$ de solution de soude de molarité $C_b$. Le pH du mélange a pour valeur 5. Calculer la valeur de $C_b$.

### Méthode de résolution
1. **Quantités initiales** :
   - $n_{HCl} = C_a \times V_a = 0,08 \times 0,2 = 0,016\text{ mol}$
   - $n_{NaOH} = C_b \times V_b = C_b \times 0,3\text{ mol}$

2. **Réaction** : $HCl + NaOH \rightarrow NaCl + H_2O$

3. **Analyse du pH = 5** :
   - $[H_3O^+] = 10^{-5}\text{ mol/L}$
   - Solution acide → HCl en excès

4. **Bilan** :
   - $n_{HCl,\text{excès}} = 0,016 - C_b \times 0,3$
   - $[H_3O^+] = \frac{n_{HCl,\text{excès}}}{V_{\text{total}}} = \frac{0,016 - 0,3C_b}{0,5}$

5. **Résolution** :
   $$10^{-5} = \frac{0,016 - 0,3C_b}{0,5}$$
   $$C_b = \frac{0,016 - 5 \times 10^{-6}}{0,3} = 0,0533\text{ mol/L}$$

## Exercice 2: Préparation et Dilution $HNO_3$

### Partie A: Solution mère
**Données** : 1,26 g $HNO_3$ dans 1L d'eau

1. **Concentration** :
   $$C_1 = \frac{m}{M \times V} = \frac{1,26}{63 \times 1} = 0,02\text{ mol/L}$$

2. **pH** (acide fort) :
   $$pH = -\log(0,02) = 1,7$$

### Partie B: Dilution au 1/10
**Objectif** : Préparer 200mL de solution $S_2$

1. **Volume à prélever** :
   $$C_1V_1 = C_2V_2 \Rightarrow V_1 = \frac{C_2 \times V_2}{C_1} = \frac{0,002 \times 200}{0,02} = 20\text{mL}$$

2. **pH de $S_2$** :
   $$pH = -\log(0,002) = 2,7$$

### Partie C: Solution commerciale
**Données** : 76% masse, densité 1,2, $V = 1\text{L}$

1. **Masse d'acide** :
   $$m_{\text{acide}} = 1000 \times 1,2 \times 0,76 = 912\text{ g}$$

2. **Concentration** :
   $$C = \frac{912}{63 \times 1} = 14,5\text{ mol/L}$$

## Exercice 3: Couple $NH_4^+$/$NH_3$

### Données
- $pK_a = 9,2$
- Solution $NH_3$ : $pH = 10,2$

### Résolution

1. **Équilibre** : $NH_3 + H_2O \rightleftharpoons NH_4^+ + OH^-$

2. **Concentrations** :
   - $pOH = 14 - 10,2 = 3,8$
   - $[OH^-] = 10^{-3,8} = 1,58 \times 10^{-4}\text{ mol/L}$

3. **Concentration initiale** :
   $$K_b = \frac{K_e}{K_a} = \frac{10^{-14}}{10^{-9,2}} = 10^{-4,8}$$
   $$K_b = \frac{[OH^-]^2}{C_0 - [OH^-]} \approx \frac{[OH^-]^2}{C_0}$$
   $$C_0 = \frac{(1,58 \times 10^{-4})^2}{10^{-4,8}} = 0,158\text{ mol/L}$$

4. **Coefficient d'ionisation** :
   $$\alpha = \frac{[OH^-]}{C_0} = \frac{1,58 \times 10^{-4}}{0,158} = 10^{-3} = 0,1\%$$

## Exercice 4: Acide Éthanoïque vs HCl

### Données
- Solution A : $CH_3COOH$, $C = 0,1\text{ mol/L}$, $pH = 2,9$
- Solution B : HCl, $C = 0,01\text{ mol/L}$

### Analyse

1. **Démonstration acide faible** :
   - Si fort : $pH = -\log(0,1) = 1$
   - Observé : $pH = 2,9$ → acide faible

2. **Coefficient d'ionisation** :
   $$[H_3O^+] = 10^{-2,9} = 1,26 \times 10^{-3}\text{ mol/L}$$
   $$\alpha = \frac{1,26 \times 10^{-3}}{0,1} = 1,26\%$$

3. **Constante d'acidité** :
   $$K_a = \frac{[H_3O^+]^2}{C_0 - [H_3O^+]} = \frac{(1,26 \times 10^{-3})^2}{0,1 - 1,26 \times 10^{-3}} = 1,6 \times 10^{-5}$$

### Effet de dilution
**Dilution 1:3** → $C' = 0,033\text{ mol/L}$, $pH = 3,1$

$$\alpha' = \frac{10^{-3,1}}{0,033} = 2,4\%$$

**Conclusion** : La dilution augmente le taux d'ionisation (loi de modération).

## Formules de Référence

### pH et concentrations
$$pH = -\log[H_3O^+]$$
$$pOH = -\log[OH^-]$$
$$pH + pOH = 14$$

### Acides faibles
$$K_a = \frac{[H_3O^+][A^-]}{[AH]}$$
$$pH = \frac{1}{2}(pK_a - \log C)$$ (si $\alpha << 1$)

### Dilutions
$$C_1V_1 = C_2V_2$$
$$\text{Facteur de dilution} = \frac{V_f}{V_i}$$
