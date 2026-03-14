---
tags: [practice, energy, hard]
---

# Energy Practice - Hard Difficulty

## Exercices avancés sur l'énergie, le travail et la puissance

### Exercice 1: Mouvement avec résistance de l'air
Un objet de masse $m = 0.5$ kg est lancé verticalement vers le haut avec une vitesse initiale $v_0 = 20$ m/s. On considère que la résistance de l'air exerce une force proportionnelle à la vitesse : $F_{air} = -kv$, où $k = 0.1$ N·s/m.

1. Écrire l'équation du mouvement de l'objet pendant sa montée.
2. Déterminer la vitesse maximale que l'objet atteindra (vitesse terminale pendant la montée).
3. Calculer la hauteur maximale atteinte par l'objet.
4. Comparer avec le cas sans résistance de l'air.

### Exercice 2: Ressort et travail variable
Un bloc de masse $m = 2$ kg est attaché à un ressort de raideur $k = 200$ N/m. Le bloc est initialement au repos à la position d'équilibre du ressort. On tire alors le bloc jusqu'à une elongation $x_0 = 0.1$ m et on le laisse partir.

1. Quelle est l'énergie potentielle élastique stockée dans le ressort lorsqu'il est allongé de $x_0$ ?
2. Quelle est la vitesse maximale du bloc pendant son oscillation ?
3. Quelle est la valeur de cette vitesse maximale ?
4. Si le coefficient de frottement entre le bloc et le sol est $\mu_k = 0.05$, quelle est la diminution de l'amplitude après chaque oscillation complète ?

### Exercice 3: Centrale hydroélectrique
Une centrale hydroélectrique utilise l'eau d'un réservoir situé à une hauteur $h = 100$ m au-dessus des turbines. Le débit d'eau traversé par les turbines est $Q = 10$ m³/s. La densité de l'eau est $\rho = 1000$ kg/m³. Le rendement global de la centrale (turbines + alternateur) est $\eta = 85\%$.

1. Quelle est la puissance mécanique disponible liée au chute d'eau ?
2. Quelle est la puissance électrique délivrée par la centrale ?
3. Quel est le rendement hydraulique si on sait que les pertes dans les turbines représentent 15% de l'énergie disponible ?
4. Quelle masse d'eau doit traverser les turbines chaque seconde pour produire une puissance électrique de 100 MW avec le même rendement ?

## Corrigés

### Exercice 1
1. Pendant la montée, les forces agissant sur l'objet sont :
   - Le poids : $mg$ vers le bas
   - La résistance de l'air : $kv$ vers le bas (puisque $v > 0$ vers le haut et $F_{air} = -kv$)
   
   Équation du mouvement (en prenant le vers le haut comme positif) :
   $mg + kv = m \frac{dv}{dt}$
   Donc : $\frac{dv}{dt} = -g - \frac{k}{m}v$
   
2. La vitesse terminale pendant la montée se produit quand $\frac{dv}{dt} = 0$ :
   $0 = -g - \frac{k}{m}v_{term}$
   Donc : $v_{term} = -\frac{mg}{k}$ (le signe négatif indique que c'est une vitesse vers le bas, ce qui n'a pas de sens pendant la montée)
   
   En fait, pendant la montée, la vitesse diminue continuellement jusqu'à devenir nulle au sommet, puis devient négative pendant la descente.
   Il n'y a pas de vitesse terminale pendant la montée au sens où l'objet arrêterait d'accélérer, car la vitesse diminue toujours jusqu'à zéro.
   
   Cependant, on peut trouver la vitesse maximale atteinte, qui est simplement la vitesse initiale $v_0 = 20$ m/s (au lancement).
   
   Pour être plus précis, la question demande probablement la vitesse terminale absolue que l'objet atteindrait finalement pendant sa chute, ce qui se produit quand la force de l'air équilibre le poids :
   $mg = kv_{term}$ donc $v_{term} = \frac{mg}{k} = \frac{0.5 \times 9.8}{0.1} = \frac{4.9}{0.1} = 49$ m/s
   
   Mais cette vitesse n'est atteinte que pendant la descente longue après de nombreuses oscillations si on considérait un mouvement vertical répété.
   
   Dans le cas d'un lancement simple vers le haut suivi d'une chute, l'objet n'atteindra jamais cette vitesse car il retombera avant d'avoir eu le temps d'accélérer suffisamment.
   
   Reformulons la question : quelle est la vitesse maximale que l'objet atteint pendant son mouvement ?
   Réponse : C'est la vitesse initiale $v_0 = 20$ m/s (au point de lancement).
   
   En fait, pendant la montée, la vitesse diminue de $v_0$ à 0, donc la vitesse maximale est bien $v_0$.
   Pendant la descente, la vitesse augmente de 0 à une valeur maximale inférieure à $v_0$ à cause de la dissipation d'énergie par la résistance de l'air.
   
   Donc la vitesse maximale absolue atteinte est $v_0 = 20$ m/s.
   
3. Pour trouver la hauteur maximale, on peut utiliser l'approximation énergétique en considérant le travail de la résistance de l'air, ou résoudre l'équation différentielle.
   
   Approximation : Travail de la résistance de l'air ≈ force moyenne × distance
   Mais une approche plus précise consiste à utiliser la conservation de l'énergie avec travail de la force non conservative :
   
   $E_{méc,i} + W_{nc} = E_{méc,f}$
   Au point de lancement : $E_{méc,i} = E_{c,i} + E_{p,i} = \frac{1}{2}mv_0^2 + 0$
   Au sommet : $E_{méc,f} = E_{c,f} + E_{p,f} = 0 + mgh_{max}$
   $W_{nc} = W_{air}$ = travail de la résistance de l'air (négatif)
   
   Donc : $\frac{1}{2}mv_0^2 + W_{air} = mgh_{max}$
   
   Pour estimer $W_{air}$, on peut utiliser que la force de l'air varie avec la vitesse.
   Une méthode consiste à résoudre l'équation différentielle :
   
   $\frac{dv}{dt} = -g - \frac{k}{m}v$
   
   Cette équation est séparable :
   $\frac{dv}{g + \frac{k}{m}v} = -dt$
   
   En intégrant de $v_0$ à $v$ et de 0 à $t$ :
   $\int_{v_0}^{v} \frac{dv'}{g + \frac{k}{m}v'} = -\int_0^t dt' = -t$
   
   $\frac{m}{k} \left[ \ln\left(g + \frac{k}{m}v'\right) \right]_{v_0}^{v} = -t$
   $\frac{m}{k} \ln\left(\frac{g + \frac{k}{m}v}{g + \frac{k}{m}v_0}\right) = -t$
   
   Pour trouver la hauteur, on utilise $v = \frac{dh}{dt}$, donc $dt = \frac{dh}{v}$ :
   
   $\frac{dv}{dt} = v \frac{dv}{dh} = -g - \frac{k}{m}v$
   $v \frac{dv}{dh} = -g - \frac{k}{m}v$
   $\frac{dv}{dh} = -\frac{g}{v} - \frac{k}{m}$
   
   En séparant les variables :
   $\frac{dv}{-\frac{g}{v} - \frac{k}{m}} = dh$
   $-\frac{v dv}{\frac{mg}{k} + v} = dh$ (en multipliant numérateur et dénominateur par $\frac{m}{k}$)
   
   En intégrant de $h=0$ à $h=h_{max}$ et de $v=v_0$ à $v=0$ :
   $-\int_{v_0}^{0} \frac{v dv}{\frac{mg}{k} + v} = \int_0^{h_{max}} dh = h_{max}$
   
   $\int_{0}^{v_0} \frac{v dv}{\frac{mg}{k} + v} = h_{max}$
   
   En posant $u = \frac{mg}{k} + v$, donc $du = dv$ et $v = u - \frac{mg}{k}$ :
   $\int_{\frac{mg}{k}}^{\frac{mg}{k}+v_0} \frac{u - \frac{mg}{k}}{u} du = h_{max}$
   $\int_{\frac{mg}{k}}^{\frac{mg}{k}+v_0} \left(1 - \frac{mg}{ku}\right) du = h_{max}$
   $\left[ u - \frac{mg}{k} \ln(u) \right]_{\frac{mg}{k}}^{\frac{mg}{k}+v_0} = h_{max}$
   
   $\left(\frac{mg}{k}+v_0 - \frac{mg}{k} \ln\left(\frac{mg}{k}+v_0\right)\right) - \left(\frac{mg}{k} - \frac{mg}{k} \ln\left(\frac{mg}{k}\right)\right) = h_{max}$
   $v_0 - \frac{mg}{k} \left[ \ln\left(\frac{mg}{k}+v_0\right) - \ln\left(\frac{mg}{k}\right) \right] = h_{max}$
   $h_{max} = v_0 - \frac{mg}{k} \ln\left(1 + \frac{kv_0}{mg}\right)$
   
   Numériquement :
   $\frac{mg}{k} = \frac{0.5 \times 9.8}{0.1} = 49$ m
   $\frac{kv_0}{mg} = \frac{0.1 \times 20}{0.5 \times 9.8} = \frac{2}{4.9} = 0.408$
   $h_{max} = 20 - 49 \times \ln(1 + 0.408) = 20 - 49 \times \ln(1.408) = 20 - 49 \times 0.342 = 20 - 16.76 = 3.24$ m
   
4. Sans résistance de l'air : $h_{max,0} = \frac{v_0^2}{2g} = \frac{20^2}{2 \times 9.8} = \frac{400}{19.6} = 20.4$ m
   
   La résistance de l'air réduit donc considérablement la hauteur maximale (de 20.4 m à 3.24 m dans notre cas).

### Exercice 2
1. $E_{el} = \frac{1}{2}kx_0^2 = \frac{1}{2} \times 200 \times (0.1)^2 = 100 \times 0.01 = 1$ J
   
2. Par conservation de l'énergie mécanique (sans frottement) :
   $E_{el,max} = E_{c,max}$
   $\frac{1}{2}kx_0^2 = \frac{1}{2}mv_{max}^2$
   Donc : $v_{max} = \sqrt{\frac{k}{m}} x_0$
   
3. $v_{max} = \sqrt{\frac{200}{2}} \times 0.1 = \sqrt{100} \times 0.1 = 10 \times 0.1 = 1$ m/s
   
4. Une oscillation complète (aller-retour) implique un déplacement total de $4x_0$ (de $+x_0$ à $-x_0$ et retour à $+x_0$).
   
   Travail du frottement sur une distance $d$ : $W_f = -f_k d = -\mu_k mg d$
   
   Pour une oscillation complète : $d = 4x_0$
   $W_f = -\mu_k mg (4x_0) = -4\mu_k mg x_0$
   
   Cette travail diminue l'énergie mécanique totale du système.
   Au point de déplacement maximal, toute l'énergie mécanique est sous forme d'énergie potentielle élastique :
   $E_{el} = \frac{1}{2}kx^2$
   
   Si l'amplitude passe de $x_0$ à $x_1$ après une oscillation complète :
   $\frac{1}{2}kx_1^2 = \frac{1}{2}kx_0^2 + W_f$
   $\frac{1}{2}kx_1^2 = \frac{1}{2}kx_0^2 - 4\mu_k mg x_0$
   $x_1^2 = x_0^2 - \frac{8\mu_k mg x_0}{k}$
   $x_1 = \sqrt{x_0^2 - \frac{8\mu_k mg x_0}{k}}$
   
   La diminution d'amplitude est : $\Delta x = x_0 - x_1$
   
   Numériquement :
   $\frac{8\mu_k mg x_0}{k} = \frac{8 \times 0.05 \times 2 \times 9.8 \times 0.1}{200} = \frac{8 \times 0.05 \times 2 \times 9.8 \times 0.1}{200}$
   $= \frac{0.784}{200} = 0.00392$
   $x_1 = \sqrt{0.1^2 - 0.00392} = \sqrt{0.01 - 0.00392} = \sqrt{0.00608} = 0.078$ m
   $\Delta x = 0.1 - 0.078 = 0.022$ m = 2.2 cm
   
   Alternativement, pour des amortissements faibles, on peut utiliser l'approximation :
   $\frac{\Delta x}{x_0} \approx \frac{2\mu_k mg}{k}$
   $\Delta x \approx x_0 \frac{2\mu_k mg}{k} = 0.1 \times \frac{2 \times 0.05 \times 2 \times 9.8}{200} = 0.1 \times \frac{3.92}{200} = 0.1 \times 0.0196 = 0.00196$ m
   
   Cette approximation n'est bonne que pour de très petits amortissements. Notre premier calcul est plus précis.
   
   En fait, en revisitant le calcul :
   Le travail du frottement sur une oscillation complète (de $+x_0$ à $-x_0$) est en fait :
   De $+x_0$ à $0$ : $d = x_0$, travail = $-\mu_k mg x_0$
   De $0$ à $-x_0$ : $d = x_0$, travail = $-\mu_k mg x_0$
   Total demi-oscillation : $-2\mu_k mg x_0$
   
   Puis de $-x_0$ à $0$ : $d = x_0$, travail = $-\mu_k mg x_0$ (le frottement s'oppose toujours au mouvement)
   De $0$ à $+x_0$ : $d = x_0$, travail = $-\mu_k mg x_0$
   Total retour : $-2\mu_k mg x_0$
   
   Donc oscillation complète : $W_f = -4\mu_k mg x_0$ (comme je l'avais initialement)
   
   Recalculons avec plus de précision :
   $W_f = -4 \times 0.05 \times 2 \times 9.8 \times 0.1 = -4 \times 0.098 = -0.392$ J
   
   Énergie initiale : $E_i = \frac{1}{2} \times 200 \times 0.1^2 = 1$ J
   Énergie après une oscillation complète : $E_f = E_i + W_f = 1 - 0.392 = 0.608$ J
   
   Nouvelle amplitude : $\frac{1}{2}k x_1^2 = 0.608$
   $x_1^2 = \frac{2 \times 0.608}{200} = \frac{1.216}{200} = 0.00608$
   $x_1 = \sqrt{0.00608} = 0.078$ m
   
   Donc l'amplitude diminue de $0.1 - 0.078 = 0.022$ m = 2.2 cm par oscillation complète.

### Exercice 3
1. Puissance mécanique disponible : $P_{méc} = \rho g h Q$
   $P_{méc} = 1000 \times 9.8 \times 100 \times 10 = 9\,800\,000$ W = 9.8 MW
   
2. Puissance électrique délivrée : $P_{elec} = \eta P_{méc} = 0.85 \times 9.8$ MW = 8.33 MW
   
3. Rendement hydraulique : Si les pertes dans les turbines représentent 15% de l'énergie disponible,
   alors l'énergie récupérée par les turbines est 85% de l'énergie disponible.
   Donc le rendement hydraulique est $\eta_{hyd} = 85\% = 0.85$
   
   (On peut aussi voir que : rendement global = rendement hydraulique × rendement électrique
   Donc : rendement électrique = $\frac{\eta_{global}}{\eta_{hyd}} = \frac{0.85}{0.85} = 1 = 100\%$
   Ce qui signifie que l'alternateur est considéré comme parfait dans ce problème.)
   
4. Pour produire $P_{elec} = 100$ MW avec le même rendement global $\eta = 85\%$ :
   Puissance mécanique nécessaire : $P_{méc}' = \frac{P_{elec}}{\eta} = \frac{100}{0.85}$ MW = 117.65 MW
   
   Comme $P_{méc} = \rho g h Q'$, on a :
   $Q' = \frac{P_{méc}'}{\rho g h} = \frac{117.65 \times 10^6}{1000 \times 9.8 \times 100} = \frac{117.65 \times 10^6}{980 \times 10^3} = \frac{117.65 \times 10^3}{980} = 120.05$ m³/s
   
   Masse d'eau par seconde : $m' = \rho Q' = 1000 \times 120.05 = 120\,050$ kg/s ≈ 120 tonnes/s

---
*Source: Adapted from advanced physics problem sets*