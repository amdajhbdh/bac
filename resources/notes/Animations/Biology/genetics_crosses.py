from manim import *
import numpy as np

config.background_color = "#1a1a2e"


class TitleScene(Scene):
    def construct(self):
        title = Text("Génétique - Croisements", font_size=48, color=WHITE)
        subtitle = Text("Animations Manim", font_size=32, color=BLUE)
        subtitle.next_to(title, DOWN)

        self.play(Write(title))
        self.play(FadeIn(subtitle, shift=UP))
        self.wait(2)


class MonohybridCompleteDominance(Scene):
    def construct(self):
        title = Text("Dominance Complète", font_size=40, color=YELLOW)
        title.to_edge(UP)
        self.add(title)

        parent_label = Text("Parents: Aa × Aa", font_size=28, color=WHITE)
        parent_label.next_to(title, DOWN, buff=0.5)
        self.play(Write(parent_label))
        self.wait()

        gamete_labels = VGroup(
            Text("A", font_size=24, color=GREEN), Text("a", font_size=24, color=RED)
        ).arrange(RIGHT, buff=2)

        parent1 = VGroup(
            Circle(radius=0.4, color=BLUE, fill_opacity=0.3),
            Text("Aa", font_size=20).move_to([-1.5, 1, 0]),
        )

        parent2 = VGroup(
            Circle(radius=0.4, color=BLUE, fill_opacity=0.3),
            Text("Aa", font_size=20).move_to([1.5, 1, 0]),
        )

        self.play(Create(parent1), Create(parent2))

        gametes_row1 = VGroup(
            Circle(radius=0.25, color=GREEN, fill_opacity=0.5),
            Text("A", font_size=16).move_to([-1.5, -0.5, 0]),
        )
        gametes_row2 = VGroup(
            Circle(radius=0.25, color=RED, fill_opacity=0.5),
            Text("a", font_size=16).move_to([-0.5, -0.5, 0]),
        )
        gametes_row3 = VGroup(
            Circle(radius=0.25, color=GREEN, fill_opacity=0.5),
            Text("A", font_size=16).move_to([0.5, -0.5, 0]),
        )
        gametes_row4 = VGroup(
            Circle(radius=0.25, color=RED, fill_opacity=0.5),
            Text("a", font_size=16).move_to([1.5, -0.5, 0]),
        )

        self.play(
            gametes_row1.animate.move_to([-1.2, 0, 0]),
            gametes_row2.animate.move_to([-0.4, 0, 0]),
            gametes_row3.animate.move_to([0.4, 0, 0]),
            gametes_row4.animate.move_to([1.2, 0, 0]),
        )

        punnett_square = VGroup()

        colors = [GREEN, YELLOW, YELLOW, RED]
        genotypes = ["AA", "Aa", "Aa", "aa"]
        positions = [(-1.2, -1, 0), (-0.4, -1, 0), (0.4, -1, 0), (1.2, -1, 0)]

        for i, (color, genotype, pos) in enumerate(zip(colors, genotypes, positions)):
            square = Square(side_length=0.8, color=WHITE, fill_opacity=0.3)
            square.set_fill(color, opacity=0.4)
            square.move_to(pos)
            label = Text(genotype, font_size=20, color=WHITE).move_to(pos)
            punnett_square.add(square, label)

        self.play(Create(punnett_square))
        self.wait()

        ratio_text = Text("Ratio: 3:1 (Jaune:Vert)", font_size=28, color=YELLOW)
        ratio_text.to_edge(DOWN)
        self.play(Write(ratio_text))
        self.wait(3)


class MonohybridIncompleteDominance(Scene):
    def construct(self):
        title = Text("Dominance Incomplète", font_size=40, color=YELLOW)
        title.to_edge(UP)
        self.add(title)

        parent_label = Text("Parents: RR × BB", font_size=28, color=WHITE)
        parent_label.next_to(title, DOWN, buff=0.5)
        self.play(Write(parent_label))
        self.wait()

        colors_list = [RED, PINK, PINK, WHITE]
        genotypes_list = ["RR", "RB", "RB", "BB"]

        for i, (color, genotype) in enumerate(zip(colors_list, genotypes_list)):
            square = Square(side_length=1.2, color=WHITE, fill_opacity=0.5)
            square.set_fill(color, opacity=0.7)
            square.shift(RIGHT * (i * 1.5 - 2.25))

            label = Text(
                genotype, font_size=24, color=BLACK if color != WHITE else WHITE
            )
            label.move_to(square.get_center())

            phenotype = Text(
                "Rouge" if i == 0 else ("Rose" if i in [1, 2] else "Blanc"),
                font_size=18,
                color=BLACK,
            ).next_to(square, DOWN)

            self.play(Create(square), Write(label), Write(phenotype))
            self.wait(0.3)

        ratio_text = Text("Ratio: 1:2:1 (Rouge:Rose:Blanc)", font_size=28, color=YELLOW)
        ratio_text.to_edge(DOWN)
        self.play(Write(ratio_text))
        self.wait(3)


class MonohybridCodominance(Scene):
    def construct(self):
        title = Text("Codominance - Groupes Sanguins", font_size=40, color=YELLOW)
        title.to_edge(UP)
        self.add(title)

        blood_groups = [
            ("I^A I^A", "A", "#ff6b6b"),
            ("I^A I^B", "AB", "#9b59b6"),
            ("I^B I^B", "B", "#3498db"),
            ("ii", "O", "#95a5a6"),
        ]

        for i, (genotype, phenotype, color) in enumerate(blood_groups):
            square = Square(side_length=1.3, color=WHITE, fill_opacity=0.5)
            square.set_fill(color, opacity=0.7)
            square.shift(RIGHT * (i * 1.7 - 2.55))

            geno_label = (
                Text(genotype, font_size=16, color=WHITE)
                .move_to(square.get_center())
                .shift(UP * 0.2)
            )
            pheno_label = (
                Text(phenotype, font_size=28, color=WHITE)
                .move_to(square.get_center())
                .shift(DOWN * 0.2)
            )

            self.play(Create(square), Write(geno_label), Write(pheno_label))
            self.wait(0.3)

        explanation = Text(
            "AB: Les deux antigènes exprimés!", font_size=22, color=YELLOW
        )
        explanation.to_edge(DOWN)
        self.play(Write(explanation))
        self.wait(3)


class DihybridIndependentAssortment(Scene):
    def construct(self):
        title = Text("Dihybridisme - Indépendance", font_size=40, color=YELLOW)
        title.to_edge(UP)
        self.add(title)

        subtitle = Text("AaBb × AaBb", font_size=28, color=WHITE)
        subtitle.next_to(title, DOWN)
        self.play(Write(subtitle))

        gametes = ["AB", "Ab", "aB", "ab"]

        grid_group = VGroup()

        for i, g1 in enumerate(gametes):
            for j, g2 in enumerate(gametes):
                square = Rectangle(width=1.5, height=0.8, color=WHITE, fill_opacity=0.3)
                x = (i - 1.5) * 1.6
                y = -j * 1.0 - 1.5
                square.move_to([x, y, 0])

                genotype = g1 + "/" + g2

                if g1[0] == "a" and g2[0] == "a":
                    color = GREEN if "b" in [g1[1], g2[1]] else YELLOW
                elif g1[1] == "b" and g2[1] == "b":
                    color = GREEN if "a" in [g1[0], g2[0]] else YELLOW
                else:
                    color = YELLOW

                if "a" in [g1[0], g2[0]] and "b" in [g1[1], g2[1]]:
                    color = RED
                if "a" in [g1[0], g2[0]] and "a" not in [g1[0], g2[0]]:
                    if "b" not in [g1[1], g2[1]]:
                        color = GREEN

                label = Text(
                    genotype, font_size=12, color=color if color != YELLOW else WHITE
                )
                label.move_to(square.get_center())

                grid_group.add(square, label)

        self.play(Create(grid_group), run_time=3)

        ratio_text = Text("Ratio: 9:3:3:1", font_size=36, color=YELLOW, weight=BOLD)
        ratio_text.to_edge(DOWN)
        self.play(Write(ratio_text))

        categories = [
            ("9: Jaune-Lisse", YELLOW),
            ("3: Jaune-Ridé", ORANGE),
            ("3: Vert-Lisse", GREEN),
            ("1: Vert-Ridé", RED),
        ]

        legend = VGroup()
        for i, (text, color) in enumerate(categories):
            t = Text(text, font_size=16, color=color)
            t.shift(DOWN * 2 + RIGHT * (i * 2.2 - 3.3))
            legend.add(t)

        self.play(Write(legend))
        self.wait(3)


class LinkedGenes(Scene):
    def construct(self):
        title = Text("Gènes Liés - Configuration Cis/Trans", font_size=36, color=YELLOW)
        title.to_edge(UP)
        self.add(title)

        cis_label = Text("Configuration CIS", font_size=28, color=GREEN)
        cis_label.shift(UP * 1.5 + LEFT * 3)

        cis_diagram = VGroup()

        chrom1 = Line(LEFT * 1.5, RIGHT * 1.5, color=BLUE, stroke_width=8)
        chrom2 = Line(LEFT * 1.5, RIGHT * 1.5, color=RED, stroke_width=8)
        chrom1.shift(UP * 0.3)
        chrom2.shift(DOWN * 0.3)

        cis_diagram.add(chrom1, chrom2)

        gene_a = Text("A", font_size=14).move_to(chrom1.get_center() + LEFT * 0.8)
        gene_b = Text("b", font_size=14).move_to(chrom1.get_center() + RIGHT * 0.8)
        gene_a2 = Text("a", font_size=14).move_to(chrom2.get_center() + LEFT * 0.8)
        gene_b2 = Text("B", font_size=14).move_to(chrom2.get_center() + RIGHT * 0.8)

        cis_diagram.add(gene_a, gene_b, gene_a2, gene_b2)
        cis_label.shift(LEFT * 3)

        trans_label = Text("Configuration TRANS", font_size=28, color=ORANGE)
        trans_label.shift(UP * 1.5 + RIGHT * 3)

        trans_diagram = VGroup()

        chrom3 = Line(LEFT * 1.5, RIGHT * 1.5, color=BLUE, stroke_width=8)
        chrom4 = Line(LEFT * 1.5, RIGHT * 1.5, color=RED, stroke_width=8)
        chrom3.shift(UP * 0.3)
        chrom4.shift(DOWN * 0.3)

        gene_a3 = Text("A", font_size=14).move_to(chrom3.get_center() + LEFT * 0.8)
        gene_b3 = Text("B", font_size=14).move_to(chrom3.get_center() + RIGHT * 0.8)
        gene_a4 = Text("a", font_size=14).move_to(chrom4.get_center() + LEFT * 0.8)
        gene_b4 = Text("b", font_size=14).move_to(chrom4.get_center() + RIGHT * 0.8)

        trans_diagram.add(chrom3, chrom4, gene_a3, gene_b3, gene_a4, gene_b4)
        trans_diagram.shift(RIGHT * 6 - LEFT * 6 + LEFT * 3)

        self.play(Write(cis_label), Write(trans_label))
        self.play(Create(cis_diagram), Create(trans_diagram))

        explanation = Text("CIS: AB/ab  |  TRANS: Ab/aB", font_size=24, color=WHITE)
        explanation.to_edge(DOWN)
        self.play(Write(explanation))
        self.wait(3)


class CrossingOver(Scene):
    def construct(self):
        title = Text("Crossing-over", font_size=40, color=YELLOW)
        title.to_edge(UP)
        self.add(title)

        before_label = Text("Avant crossing-over", font_size=24, color=GREEN)
        before_label.shift(UP * 1.5)

        chrom_left = Line(UP * 1, DOWN * 1, color=BLUE, stroke_width=6)
        chrom_right = Line(UP * 1, DOWN * 1, color=RED, stroke_width=6)

        gene_a1 = Text("A", font_size=12, color=BLUE).move_to(
            chrom_left.get_center() + UP * 0.5
        )
        gene_b1 = Text("B", font_size=12, color=RED).move_to(
            chrom_left.get_center() + DOWN * 0.5
        )
        gene_a2 = Text("a", font_size=12, color=RED).move_to(
            chrom_right.get_center() + UP * 0.5
        )
        gene_b2 = Text("b", font_size=12, color=BLUE).move_to(
            chrom_right.get_center() + DOWN * 0.5
        )

        self.play(
            Write(before_label),
            Create(chrom_left),
            Create(chrom_right),
            Write(gene_a1),
            Write(gene_b1),
            Write(gene_a2),
            Write(gene_b2),
        )

        self.wait()

        after_label = Text("Après crossing-over", font_size=24, color=ORANGE)
        after_label.shift(UP * 1.5 + RIGHT * 3)

        cross_chrom1 = Line(
            LEFT * 0.3 + UP * 1, RIGHT * 0.3 + DOWN * 1, color=BLUE, stroke_width=6
        )
        cross_chrom2 = Line(
            LEFT * 0.3 + DOWN * 1, RIGHT * 0.3 + UP * 1, color=RED, stroke_width=6
        )

        new_gene1 = Text("A-b", font_size=10, color=WHITE).move_to(
            cross_chrom1.get_center()
        )
        new_gene2 = Text("a-B", font_size=10, color=WHITE).move_to(
            cross_chrom2.get_center()
        )

        self.play(Write(after_label))
        self.play(
            Transform(chrom_left.copy(), cross_chrom1),
            Transform(chrom_right.copy(), cross_chrom2),
            Write(new_gene1),
            Write(new_gene2),
        )

        formula = Text(
            "Fréquence recombinaison = (recombinants / total) × 100%",
            font_size=22,
            color=YELLOW,
        )
        formula.to_edge(DOWN)
        self.play(Write(formula))
        self.wait(3)


class EpistasisTypes(Scene):
    def construct(self):
        title = Text("Types d'Épistasie", font_size=40, color=YELLOW)
        title.to_edge(UP)
        self.add(title)

        epistasis_types = [
            ("9:3:4", "Épistasie récéssive", RED),
            ("12:3:1", "Épistasie dominante", BLUE),
            ("9:7", "Épistasie complémentaire", GREEN),
            ("15:1", "Épistasie duplicate dominante", ORANGE),
            ("9:3:3:1", "Indépendance (pas d'épistasie)", WHITE),
        ]

        for i, (ratio, name, color) in enumerate(epistasis_types):
            ratio_text = Text(ratio, font_size=32, color=color, weight=BOLD)
            ratio_text.shift(LEFT * 3 + DOWN * (i * 0.8 - 1.5))

            name_text = Text(name, font_size=20, color=WHITE)
            name_text.shift(RIGHT * 2 + DOWN * (i * 0.8 - 1.5))

            self.play(Write(ratio_text), Write(name_text))
            self.wait(0.3)

        self.wait(3)


class ProbabilityCalculations(Scene):
    def construct(self):
        title = Text("Calculs de Probabilités", font_size=40, color=YELLOW)
        title.to_edge(UP)
        self.add(title)

        formula1 = MathTex(r"P(Aa \times Aa) = \frac{1}{2}", font_size=36, color=WHITE)
        formula1.shift(UP * 1.5)

        formula2 = MathTex(
            r"P(A-B-) = \frac{3}{4} \times \frac{3}{4} = \frac{9}{16}",
            font_size=36,
            color=WHITE,
        )

        formula3 = MathTex(
            r"P(aabb) = \frac{1}{4} \times \frac{1}{4} = \frac{1}{16}",
            font_size=36,
            color=WHITE,
        )
        formula3.shift(DOWN * 1.5)

        self.play(Write(formula1))
        self.play(Write(formula2))
        self.play(Write(formula3))

        rule = Text(
            "Règle: événements indépendants → multiplication",
            font_size=22,
            color=YELLOW,
        )
        rule.to_edge(DOWN)
        self.play(Write(rule))

        self.wait(3)


class PunnettSquareAnimation(Scene):
    def construct(self):
        title = Text("Construction Grille de Punnett", font_size=36, color=YELLOW)
        title.to_edge(UP)
        self.add(title)

        gametes_top = ["A", "a"]
        gametes_left = ["A", "a"]

        for i, gamete in enumerate(gametes_top):
            label = Text(gamete, font_size=24, color=GREEN)
            label.move_to([i * 1.5 - 0.75, 1.5, 0])
            self.play(Write(label))

        for i, gamete in enumerate(gametes_left):
            label = Text(gamete, font_size=24, color=RED)
            label.move_to([-2, -i * 1.5 + 0.5, 0])
            self.play(Write(label))

        for i, g1 in enumerate(gametes_left):
            for j, g2 in enumerate(gametes_top):
                square = Square(side_length=1.2, color=WHITE, fill_opacity=0.3)
                square.move_to([j * 1.5 - 0.75, -i * 1.5 + 0.5, 0])

                genotype = g1 + g2
                if genotype in ["AA", "aa"]:
                    color = GREEN if g1 == g2 == "A" else RED
                else:
                    color = YELLOW
                square.set_fill(color, opacity=0.5)

                label = Text(genotype, font_size=20, color=WHITE)
                label.move_to(square.get_center())

                self.play(Create(square), Write(label))
                self.wait(0.2)

        self.wait(3)


class RecapScene(Scene):
    def construct(self):
        title = Text("Récapitulatif - Tous les Cas", font_size=40, color=YELLOW)
        title.to_edge(UP)
        self.add(title)

        cases = [
            (
                "1. Monohybridisme",
                [
                    "• Dominance complète (3:1)",
                    "• Dominance incomplète (1:2:1)",
                    "• Codominance (1:2:1)",
                    "• Allèles multiples",
                ],
            ),
            (
                "2. Dihybridisme",
                [
                    "• Indépendance (9:3:3:1)",
                    "• Gènes liés",
                    "• Crossing-over",
                    "• Épistasie",
                ],
            ),
        ]

        for i, (category, items) in enumerate(cases):
            cat_text = Text(category, font_size=24, color=GREEN if i == 0 else BLUE)
            cat_text.shift(LEFT * 3 + UP * (1 - i * 2))

            for j, item in enumerate(items):
                item_text = Text(item, font_size=16, color=WHITE)
                item_text.shift(LEFT * 3 + UP * (1 - i * 2 - 0.5) + DOWN * j * 0.4)
                self.play(Write(cat_text), Write(item_text), run_time=0.5)

        self.wait(3)


class MultipleAlleles(Scene):
    def construct(self):
        title = Text("Allèles Multiples - Système ABO", font_size=38, color=YELLOW)
        title.to_edge(UP)
        self.add(title)

        alleles = [
            (r"I^A", "Groupe A", "#ff6b6b"),
            (r"I^B", "Groupe B", "#3498db"),
            ("I^I^B", "Groupe AB", "#9b59b6"),
            ("i", "Groupe O", "#95a5a6"),
        ]

        for i, (allele, group, color) in enumerate(alleles):
            circle = Circle(radius=0.7, color=color, fill_opacity=0.6)
            circle.shift(RIGHT * (i * 2 - 3))

            allele_text = Text(allele, font_size=16, color=WHITE)
            allele_text.move_to(circle.get_center() + UP * 0.2)

            group_text = Text(group, font_size=18, color=WHITE)
            group_text.move_to(circle.get_center() + DOWN * 0.3)

            self.play(Create(circle), Write(allele_text), Write(group_text))
            self.wait(0.3)

        note = Text(
            "3 alleles dans la population: I^A, I^B, i", font_size=20, color=YELLOW
        )
        note.to_edge(DOWN)
        self.play(Write(note))

        self.wait(3)


class GeneticDistance(Scene):
    def construct(self):
        title = Text("Distance Génétique", font_size=40, color=YELLOW)
        title.to_edge(UP)
        self.add(title)

        formula = MathTex(
            r"d = \frac{N_{recombinants}}{N_{total}} \times 100\%",
            font_size=42,
            color=WHITE,
        )
        formula.shift(UP * 0.5)

        self.play(Write(formula))

        examples = [
            ("< 10 cM", "Liaison forte", GREEN),
            ("10-25 cM", "Liaison modérée", YELLOW),
            ("> 25 cM", "Indépendance apparente", RED),
        ]

        for i, (distance, desc, color) in enumerate(examples):
            dist_text = Text(distance, font_size=24, color=color)
            dist_text.shift(DOWN * 1 + LEFT * 2 + DOWN * i * 0.8)

            desc_text = Text(desc, font_size=18, color=WHITE)
            desc_text.shift(DOWN * 1 + RIGHT * 2 + DOWN * i * 0.8)

            self.play(Write(dist_text), Write(desc_text))

        self.wait(3)
