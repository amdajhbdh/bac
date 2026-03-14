---
tags: [practice, energy, medium]
---

# Energy Practice - Medium Difficulty

## Exercices d'application sur l'énergie, le travail et la puissance

### Exercice 1: Travail d'une force variable
Une force qui varie avec la position selon la loi $F(x) = 5x$ N (où x est en mètres) agit sur un objet qui se déplace de $x = 0$ à $x = 4$ m selon l'axe des x.

1. Quel est le travail réalisé par cette force ?
2. Si l'objet a une masse de $m = 2$ kg et part du repos, quelle est sa vitesse à $x = 4$ m ?
3. Quelle est la puissance moyenne développée par cette force sur le déplacement si le trajet prend 2 secondes ?

### Exercice 2: Energie mécanique et frottement
Un bloc de masse $m = 5$ kg glisse sur un plan incliné faisant un angle $\theta = 30^\circ$ avec l'horizontale. Le coefficient de frottement cinétique est $\mu_k = 0.2$. La hauteur verticale du plan est $h = 3$ m.

1. Quel est le travail réalisé par la force de gravité pendant la descente ?
2. Quel est le travail réalisé par la force de frottement ?
3. Quelle est la variation de l'énergie mécanique du bloc ?
4. Quelle est la vitesse du bloc en bas du plan ?

### Exercice 3: Puissance et efficacité
Un moteur électrique soulève une charge de masse $m = 100$ kg à une vitesse constante de $v = 0.5$ m/s. Le moteur consomme une puissance électrique de $P_{elec} = 600$ W.

1. Quelle est la puissance mécanique utile développée par le moteur ?
2. Quelle est l'efficacité du moteur ?
3. Si le moteur fonctionne pendant 1 heure, quelle quantité d'énergie électrique consomme-t-il ?
4. Quelle est l'augmentation de l'énergie potentielle de la charge après cette heure ?

## Corrigés

### Exercice 1
1. Le travail d'une force variable est donné par $W = \int_{x_i}^{x_f} F(x) dx$
   $W = \int_0^4 5x dx = 5 \left[ \frac{x^2}{2} \right]_0^4 = 5 \times \frac{16}{2} = 5 \times 8 = 40$ J
   
2. Par le théorème de l'énergie cinétique : $W = \Delta E_c = \frac{1}{2}mv_f^2 - \frac{1}{2}mv_i^2$
   Puisque $v_i = 0$ : $40 = \frac{1}{2} \times 2 \times v_f^2 = v_f^2$
   Donc : $v_f = \sqrt{40} = 6.32$ m/s
   
3. Puissance moyenne : $P_{moy} = \frac{W}{\Delta t} = \frac{40}{2} = 20$ W

### Exercice 2
1. Travail de la gravité : $W_g = mgh = 5 \times 9.8 \times 3 = 147$ J
   (Positif car la force de gravité et le déplacement sont dans la même direction verticale)
   
2. Pour trouver le travail du frottement, il faut d'abord trouver la distance parcourue le long du plan :
   $\sin\theta = \frac{h}{d}$ donc $d = \frac{h}{\sin\theta} = \frac{3}{\sin 30^\circ} = \frac{3}{0.5} = 6$ m
   
   La force normale : $N = mg\cos\theta = 5 \times 9.8 \times \cos 30^\circ = 49 \times 0.866 = 42.4$ N
   La force de frottement : $f_k = \mu_k N = 0.2 \times 42.4 = 8.48$ N
   Le travail du frottement : $W_f = -f_k d = -8.48 \times 6 = -50.9$ J
   (Négatif car la force de frottement s'oppose au déplacement)
   
3. Variation de l'énergie mécanique : $\Delta E_{mec} = W_g + W_f = 147 - 50.9 = 96.1$ J
   (Ceci est aussi égal au travail des forces non conservatrices, ici seulement le frottement)
   
4. Cette variation se traduit par une augmentation de l'énergie cinétique (en supposant qu'il parte du repos) :
   $\Delta E_{mec} = \frac{1}{2}mv^2 - 0$
   $96.1 = \frac{1}{2} \times 5 \times v^2$
   $v^2 = \frac{96.1 \times 2}{5} = 38.44$
   $v = \sqrt{38.44} = 6.20$ m/s

### Exercice 3
1. Puissance mécanique utile : $P_{mec} = Fv = mgv = 100 \times 9.8 \times 0.5 = 490$ W
   (La force nécessaire pour soulever la charge à vitesse constante est égale à son poids)
   
2. Efficacité : $\eta = \frac{P_{mec}}{P_{elec}} \times 100\% = \frac{490}{600} \times 100\% = 81.7\%$
   
3. Énergie électrique consommée en 1 heure : $E_{elec} = P_{elec} \times t = 600 \text{ W} \times 3600 \text{ s} = 2\,160\,000$ J = 2.16 MJ
   
4. En 1 heure, la distance parcourue est : $d = vt = 0.5 \times 3600 = 1800$ m
   Augmentation de l'énergie potentielle : $\Delta E_p = mgd = 100 \times 9.8 \times 1800 = 1\,764\,000$ J = 1.764 MJ
   (On peut aussi vérifier que : $E_{elec} \times \eta = 2.16 \times 0.817 = 1.764$ MJ, ce qui correspond)

---
*Source: Adapted from physics problem sets*