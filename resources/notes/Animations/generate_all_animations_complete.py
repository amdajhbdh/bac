from manim import *


class HormoneRegulation(Scene):
    def construct(self):
        title = Text("Hormone Regulation: FSH & LH", font_size=36)
        self.play(Write(title))
        self.wait()
        self.play(title.animate.to_edge(UP))

        axes = Axes(
            x_range=[0, 28, 7],
            y_range=[0, 10, 2],
            x_length=10,
            y_length=5,
            axis_config={"include_tip": True},
        )
        x_label = axes.get_x_axis_label("Days")
        y_label = axes.get_y_axis_label("Hormone Level")

        fsh_graph = axes.plot(lambda x: 5 + 3 * np.sin(x * np.pi / 14), color=BLUE)
        lh_graph = axes.plot(lambda x: 4 + 4 * np.sin((x - 7) * np.pi / 14), color=RED)

        fsh_label = Text("FSH", color=BLUE, font_size=24).next_to(fsh_graph, UP)
        lh_label = Text("LH", color=RED, font_size=24).next_to(lh_graph, DOWN)

        self.play(Create(axes), Write(x_label), Write(y_label))
        self.play(Create(fsh_graph), Write(fsh_label))
        self.play(Create(lh_graph), Write(lh_label))
        self.wait(2)


class NeuronSynapse(Scene):
    def construct(self):
        title = Text("Synaptic Transmission", font_size=36)
        self.play(Write(title))
        self.wait()
        self.play(title.animate.to_edge(UP))

        neuron1 = Rectangle(width=3, height=1, color=BLUE).shift(LEFT * 3)
        neuron2 = Rectangle(width=3, height=1, color=GREEN).shift(RIGHT * 3)
        synapse = Line(neuron1.get_right(), neuron2.get_left(), color=YELLOW)

        self.play(Create(neuron1), Create(neuron2), Create(synapse))

        for _ in range(3):
            vesicle = Dot(color=RED).move_to(neuron1.get_right())
            self.play(vesicle.animate.move_to(neuron2.get_left()), run_time=0.8)
            self.remove(vesicle)

        self.wait(2)


class GeneticCrossing(Scene):
    def construct(self):
        title = Text("Genetic Crossing: Drosophila", font_size=36)
        self.play(Write(title))
        self.wait()
        self.play(title.animate.to_edge(UP))

        p1 = Text("P1: AA × aa", font_size=28).shift(UP * 2)
        f1 = Text("F1: Aa (100%)", font_size=28)
        f2 = Text("F2: AA : Aa : aa = 1:2:1", font_size=28).shift(DOWN * 2)

        self.play(Write(p1))
        self.wait()
        self.play(Write(f1))
        self.wait()
        self.play(Write(f2))
        self.wait(2)


class ImmuneResponse(Scene):
    def construct(self):
        title = Text("Immune Response: B Lymphocyte", font_size=36)
        self.play(Write(title))
        self.wait()
        self.play(title.animate.to_edge(UP))

        b_cell = Circle(radius=0.5, color=BLUE).shift(LEFT * 3)
        plasma_cell = Circle(radius=0.7, color=RED).shift(RIGHT * 3)

        b_label = Text("B Cell", font_size=20).next_to(b_cell, DOWN)
        p_label = Text("Plasma Cell", font_size=20).next_to(plasma_cell, DOWN)

        self.play(Create(b_cell), Write(b_label))
        self.wait()
        self.play(Transform(b_cell.copy(), plasma_cell), Write(p_label))
        self.wait(2)


class GeneticBrassage(Scene):
    def construct(self):
        title = Text("Brassage de l'information génétique", font_size=32)
        self.play(Write(title))
        self.wait()
        self.play(title.animate.to_edge(UP))

        chromosomes = VGroup(
            *[Line(UP * 0.5, DOWN * 0.5, color=random_bright_color()) for _ in range(4)]
        ).arrange(RIGHT, buff=0.5)

        self.play(Create(chromosomes))
        self.play(chromosomes.animate.shift(DOWN))

        # Crossing over
        self.play(Swap(chromosomes[0], chromosomes[2]))
        self.play(Swap(chromosomes[1], chromosomes[3]))
        self.wait(2)


# ============================================
# CHEMISTRY ANIMATIONS (3 exercises)
# ============================================


class AlcoholReactions(Scene):
    def construct(self):
        title = Text("Alcohol Reactions", font_size=36)
        self.play(Write(title))
        self.wait()
        self.play(title.animate.to_edge(UP))

        alcohol = MathTex(r"R-OH", font_size=48)
        arrow = Arrow(LEFT, RIGHT, color=YELLOW).next_to(alcohol, RIGHT)
        product = MathTex(r"R-O-R'", font_size=48).next_to(arrow, RIGHT)

        self.play(Write(alcohol))
        self.wait()
        self.play(Create(arrow))
        self.play(Write(product))
        self.wait(2)


class OrganicChemistry(Scene):
    def construct(self):
        title = Text("Organic Chemistry: Functional Groups", font_size=32)
        self.play(Write(title))
        self.wait()
        self.play(title.animate.to_edge(UP))

        groups = VGroup(
            MathTex(r"-OH", color=BLUE),
            MathTex(r"-COOH", color=RED),
            MathTex(r"-NH_2", color=GREEN),
            MathTex(r"-CHO", color=YELLOW),
        ).arrange(DOWN, buff=0.5)

        labels = (
            VGroup(
                Text("Alcohol", font_size=20),
                Text("Carboxylic Acid", font_size=20),
                Text("Amine", font_size=20),
                Text("Aldehyde", font_size=20),
            )
            .arrange(DOWN, buff=0.5)
            .next_to(groups, RIGHT, buff=1)
        )

        self.play(Write(groups), Write(labels))
        self.wait(2)


class ChemicalKinetics(Scene):
    def construct(self):
        title = Text("Cinétique Chimique", font_size=36)
        self.play(Write(title))
        self.wait()
        self.play(title.animate.to_edge(UP))

        axes = Axes(x_range=[0, 10, 2], y_range=[0, 100, 20], x_length=8, y_length=5)
        x_label = axes.get_x_axis_label("Time (s)")
        y_label = axes.get_y_axis_label("Concentration")

        curve = axes.plot(lambda x: 100 * np.exp(-0.3 * x), color=BLUE)

        self.play(Create(axes), Write(x_label), Write(y_label))
        self.play(Create(curve))
        self.wait(2)


# ============================================
# MATHEMATICS ANIMATIONS (7 exercises)
# ============================================


class ComplexNumbersPlane(Scene):
    def construct(self):
        title = Text("Complex Numbers", font_size=36)
        self.play(Write(title))
        self.wait()
        self.play(title.animate.to_edge(UP))

        plane = ComplexPlane(x_range=[-3, 3], y_range=[-3, 3])
        self.play(Create(plane))

        z1 = Dot(plane.n2p(2 + 1j), color=RED)
        z2 = Dot(plane.n2p(-1 + 2j), color=BLUE)

        label1 = MathTex("z_1 = 2+i", font_size=24).next_to(z1, UR)
        label2 = MathTex("z_2 = -1+2i", font_size=24).next_to(z2, UL)

        self.play(Create(z1), Write(label1))
        self.play(Create(z2), Write(label2))
        self.wait(2)


class MatrixOperations(Scene):
    def construct(self):
        title = Text("Matrix Multiplication", font_size=36)
        self.play(Write(title))
        self.wait()
        self.play(title.animate.to_edge(UP))

        m1 = Matrix([[1, 2], [3, 4]], left_bracket="(", right_bracket=")")
        m2 = Matrix([[5, 6], [7, 8]], left_bracket="(", right_bracket=")")
        result = Matrix([[19, 22], [43, 50]], left_bracket="(", right_bracket=")")

        eq = VGroup(m1, MathTex(r"\times"), m2, MathTex("="), result).arrange(
            RIGHT, buff=0.3
        )

        self.play(Write(m1))
        self.play(Write(eq[1]), Write(m2))
        self.wait()
        self.play(Write(eq[3]), Write(result))
        self.wait(2)


class IntegralCalculation(Scene):
    def construct(self):
        title = Text("Definite Integral", font_size=36)
        self.play(Write(title))
        self.wait()
        self.play(title.animate.to_edge(UP))

        axes = Axes(x_range=[0, 5, 1], y_range=[0, 10, 2], x_length=8, y_length=5)
        curve = axes.plot(lambda x: x**2, color=BLUE)
        area = axes.get_riemann_rectangles(
            curve, x_range=[1, 3], dx=0.2, color=YELLOW, fill_opacity=0.5
        )

        integral = MathTex(r"\int_1^3 x^2 \, dx = \frac{26}{3}", font_size=32).to_edge(
            DOWN
        )

        self.play(Create(axes), Create(curve))
        self.play(Create(area))
        self.play(Write(integral))
        self.wait(2)


class ArithmeticSequence(Scene):
    def construct(self):
        title = Text("Arithmétique", font_size=36)
        self.play(Write(title))
        self.wait()
        self.play(title.animate.to_edge(UP))

        sequence = MathTex(r"u_n = u_0 + n \cdot r", font_size=40)
        example = MathTex(r"2, 5, 8, 11, 14, ...", font_size=36).shift(DOWN)

        self.play(Write(sequence))
        self.wait()
        self.play(Write(example))
        self.wait(2)


class ArithmeticIntegral(Scene):
    def construct(self):
        title = Text("Arithmétique & Intégrale", font_size=36)
        self.play(Write(title))
        self.wait()
        self.play(title.animate.to_edge(UP))

        formula1 = MathTex(r"S_n = \frac{n(u_0 + u_n)}{2}", font_size=36).shift(UP)
        formula2 = MathTex(r"\int_a^b f(x) dx", font_size=36).shift(DOWN)

        self.play(Write(formula1))
        self.wait()
        self.play(Write(formula2))
        self.wait(2)


class ComplexNumberAdvanced(Scene):
    def construct(self):
        title = Text("Nombre Complexe: Forme Exponentielle", font_size=32)
        self.play(Write(title))
        self.wait()
        self.play(title.animate.to_edge(UP))

        forms = VGroup(
            MathTex(r"z = a + bi", font_size=32),
            MathTex(r"z = r e^{i\theta}", font_size=32),
            MathTex(r"z = r(\cos\theta + i\sin\theta)", font_size=32),
        ).arrange(DOWN, buff=0.7)

        self.play(Write(forms[0]))
        self.wait()
        self.play(Write(forms[1]))
        self.wait()
        self.play(Write(forms[2]))
        self.wait(2)


class MatrixDeterminant(Scene):
    def construct(self):
        title = Text("Matrix Determinant", font_size=36)
        self.play(Write(title))
        self.wait()
        self.play(title.animate.to_edge(UP))

        matrix = Matrix([[1, 2], [3, 4]], left_bracket="|", right_bracket="|")
        formula = MathTex(r"= 1 \cdot 4 - 2 \cdot 3 = -2", font_size=32)

        eq = VGroup(matrix, formula).arrange(RIGHT, buff=0.5)

        self.play(Write(matrix))
        self.wait()
        self.play(Write(formula))
        self.wait(2)


# ============================================
# PHYSICS ANIMATIONS (15 exercises)
# ============================================


class ProjectileMotion(Scene):
    def construct(self):
        title = Text("Projectile Motion", font_size=36)
        self.play(Write(title))
        self.wait()
        self.play(title.animate.to_edge(UP))

        axes = Axes(x_range=[0, 10, 2], y_range=[0, 5, 1], x_length=10, y_length=5)
        trajectory = axes.plot(lambda x: x - 0.1 * x**2, color=BLUE)

        ball = Dot(color=RED).move_to(axes.c2p(0, 0))

        self.play(Create(axes))
        self.play(Create(trajectory))
        self.play(MoveAlongPath(ball, trajectory), run_time=3)
        self.wait(2)


class SpringOscillation(Scene):
    def construct(self):
        title = Text("Spring Oscillation", font_size=36)
        self.play(Write(title))
        self.wait()
        self.play(title.animate.to_edge(UP))

        spring = Line(UP * 2, DOWN * 2, color=BLUE)
        mass = Square(side_length=0.5, color=RED).next_to(spring, DOWN, buff=0)

        self.play(Create(spring), Create(mass))

        for _ in range(3):
            self.play(mass.animate.shift(DOWN * 0.5), run_time=0.5)
            self.play(mass.animate.shift(UP * 0.5), run_time=0.5)

        self.wait(2)


class ElectricField(Scene):
    def construct(self):
        title = Text("Electric Field", font_size=36)
        self.play(Write(title))
        self.wait()
        self.play(title.animate.to_edge(UP))

        charge1 = Dot(LEFT * 2, color=RED, radius=0.2)
        charge2 = Dot(RIGHT * 2, color=BLUE, radius=0.2)

        label1 = MathTex("+Q", font_size=24).next_to(charge1, UP)
        label2 = MathTex("-Q", font_size=24).next_to(charge2, UP)

        self.play(Create(charge1), Create(charge2), Write(label1), Write(label2))

        # Field lines
        for angle in range(0, 360, 30):
            line = Arrow(
                charge1.get_center(),
                charge1.get_center()
                + np.array([np.cos(np.radians(angle)), np.sin(np.radians(angle)), 0])
                * 1.5,
                color=YELLOW,
                buff=0.2,
            )
            self.play(Create(line), run_time=0.2)

        self.wait(2)


class DynamicsForces(Scene):
    def construct(self):
        title = Text("Dynamics: F = ma", font_size=36)
        self.play(Write(title))
        self.wait()
        self.play(title.animate.to_edge(UP))

        box = Square(side_length=1, color=BLUE)
        force = Arrow(box.get_right(), box.get_right() + RIGHT * 2, color=RED, buff=0)
        force_label = MathTex("F", font_size=32, color=RED).next_to(force, UP)

        self.play(Create(box))
        self.play(Create(force), Write(force_label))
        self.play(box.animate.shift(RIGHT * 3), run_time=2)
        self.wait(2)


class MagneticForce(Scene):
    def construct(self):
        title = Text("Magnetic Force (Laplace)", font_size=36)
        self.play(Write(title))
        self.wait()
        self.play(title.animate.to_edge(UP))

        wire = Line(LEFT * 3, RIGHT * 3, color=YELLOW)
        current = Arrow(LEFT * 2, RIGHT * 2, color=RED).shift(UP * 0.2)

        b_field = VGroup(
            *[
                Arrow(DOWN * 0.5, UP * 0.5, color=BLUE).shift(RIGHT * i)
                for i in range(-2, 3)
            ]
        ).shift(DOWN * 2)

        formula = MathTex(r"F = I L B \sin\theta", font_size=32).to_edge(DOWN)

        self.play(Create(wire), Create(current))
        self.play(Create(b_field))
        self.play(Write(formula))
        self.wait(2)


class MechanicsEnergy(Scene):
    def construct(self):
        title = Text("Mechanical Energy", font_size=36)
        self.play(Write(title))
        self.wait()
        self.play(title.animate.to_edge(UP))

        formulas = VGroup(
            MathTex(r"E_k = \frac{1}{2}mv^2", font_size=32),
            MathTex(r"E_p = mgh", font_size=32),
            MathTex(r"E_m = E_k + E_p", font_size=32),
        ).arrange(DOWN, buff=0.7)

        self.play(Write(formulas[0]))
        self.wait()
        self.play(Write(formulas[1]))
        self.wait()
        self.play(Write(formulas[2]))
        self.wait(2)


class ProjectileDetailed(Scene):
    def construct(self):
        title = Text("Projectile: Velocity Components", font_size=32)
        self.play(Write(title))
        self.wait()
        self.play(title.animate.to_edge(UP))

        axes = Axes(x_range=[0, 10, 2], y_range=[0, 5, 1], x_length=9, y_length=5)

        v_x = MathTex(r"v_x = v_0 \cos\theta", font_size=28).to_corner(UL).shift(DOWN)
        v_y = MathTex(r"v_y = v_0 \sin\theta - gt", font_size=28).next_to(v_x, DOWN)

        self.play(Create(axes))
        self.play(Write(v_x), Write(v_y))

        trajectory = axes.plot(lambda x: 2 * x - 0.2 * x**2, color=BLUE)
        self.play(Create(trajectory))
        self.wait(2)


class SpringEnergy(Scene):
    def construct(self):
        title = Text("Spring Energy", font_size=36)
        self.play(Write(title))
        self.wait()
        self.play(title.animate.to_edge(UP))

        formulas = VGroup(
            MathTex(r"E_p = \frac{1}{2}kx^2", font_size=36),
            MathTex(r"E_k = \frac{1}{2}mv^2", font_size=36),
            MathTex(r"E_m = \text{constant}", font_size=36),
        ).arrange(DOWN, buff=0.7)

        self.play(Write(formulas[0]))
        self.wait()
        self.play(Write(formulas[1]))
        self.wait()
        self.play(Write(formulas[2]))
        self.wait(2)


class InductionMagnetique(Scene):
    def construct(self):
        title = Text("Induction Magnétique", font_size=36)
        self.play(Write(title))
        self.wait()
        self.play(title.animate.to_edge(UP))

        coil = Circle(radius=1, color=BLUE)
        magnet = Rectangle(width=0.5, height=1.5, color=RED).shift(LEFT * 3)

        formula = MathTex(r"\varepsilon = -\frac{d\Phi}{dt}", font_size=36).to_edge(
            DOWN
        )

        self.play(Create(coil), Create(magnet))
        self.play(magnet.animate.move_to(coil.get_center()), run_time=2)
        self.play(Write(formula))
        self.wait(2)


class RajaSpring2024(Scene):
    def construct(self):
        title = Text("Raja Spring Problem 2024", font_size=32)
        self.play(Write(title))
        self.wait()
        self.play(title.animate.to_edge(UP))

        # Spring system
        spring = Line(UP * 2, ORIGIN, color=BLUE)
        mass = Square(side_length=0.6, color=RED).next_to(spring, DOWN, buff=0)

        data = (
            VGroup(
                MathTex(r"k = 100 \, \text{N/m}", font_size=24),
                MathTex(r"m = 0.5 \, \text{kg}", font_size=24),
                MathTex(r"T = 2\pi\sqrt{\frac{m}{k}}", font_size=24),
            )
            .arrange(DOWN, buff=0.3)
            .to_corner(UR)
        )

        self.play(Create(spring), Create(mass))
        self.play(Write(data))

        # Oscillation
        for _ in range(2):
            self.play(mass.animate.shift(DOWN * 0.7), run_time=0.6)
            self.play(mass.animate.shift(UP * 0.7), run_time=0.6)

        self.wait(2)


class ProjectileResume(Scene):
    def construct(self):
        title = Text("Projectile: Résumé", font_size=36)
        self.play(Write(title))
        self.wait()
        self.play(title.animate.to_edge(UP))

        equations = VGroup(
            MathTex(r"x(t) = v_0 \cos\theta \cdot t", font_size=28),
            MathTex(r"y(t) = v_0 \sin\theta \cdot t - \frac{1}{2}gt^2", font_size=28),
            MathTex(r"v_x = v_0 \cos\theta", font_size=28),
            MathTex(r"v_y = v_0 \sin\theta - gt", font_size=28),
        ).arrange(DOWN, buff=0.5)

        for eq in equations:
            self.play(Write(eq))
            self.wait(0.5)

        self.wait(2)


class EnergyRecall(Scene):
    def construct(self):
        title = Text("Rappel sur les Énergies", font_size=32)
        self.play(Write(title))
        self.wait()
        self.play(title.animate.to_edge(UP))

        energy_types = VGroup(
            Text("Énergie Cinétique", font_size=24, color=BLUE),
            Text("Énergie Potentielle", font_size=24, color=RED),
            Text("Énergie Mécanique", font_size=24, color=GREEN),
        ).arrange(DOWN, buff=0.8)

        formulas = (
            VGroup(
                MathTex(r"E_c = \frac{1}{2}mv^2", font_size=28),
                MathTex(r"E_p = mgh", font_size=28),
                MathTex(r"E_m = E_c + E_p", font_size=28),
            )
            .arrange(DOWN, buff=0.8)
            .next_to(energy_types, RIGHT, buff=1)
        )

        for i in range(3):
            self.play(Write(energy_types[i]), Write(formulas[i]))
            self.wait(0.5)

        self.wait(2)


class DynamicsDetailed(Scene):
    def construct(self):
        title = Text("Dynamique: Forces et Mouvement", font_size=32)
        self.play(Write(title))
        self.wait()
        self.play(title.animate.to_edge(UP))

        # Object
        obj = Circle(radius=0.5, color=BLUE)

        # Forces
        weight = Arrow(
            obj.get_bottom(), obj.get_bottom() + DOWN * 1.5, color=RED, buff=0
        )
        normal = Arrow(obj.get_top(), obj.get_top() + UP * 1.5, color=GREEN, buff=0)
        friction = Arrow(
            obj.get_left(), obj.get_left() + LEFT * 1.5, color=YELLOW, buff=0
        )

        labels = VGroup(
            MathTex(r"\vec{P}", color=RED).next_to(weight, RIGHT),
            MathTex(r"\vec{N}", color=GREEN).next_to(normal, RIGHT),
            MathTex(r"\vec{f}", color=YELLOW).next_to(friction, DOWN),
        )

        self.play(Create(obj))
        self.play(Create(weight), Create(normal), Create(friction))
        self.play(Write(labels))
        self.wait(2)


class MechanicsProblems(Scene):
    def construct(self):
        title = Text("Mécanique: Problèmes Types", font_size=32)
        self.play(Write(title))
        self.wait()
        self.play(title.animate.to_edge(UP))

        problems = (
            VGroup(
                Text("1. Chute libre", font_size=24),
                Text("2. Plan incliné", font_size=24),
                Text("3. Mouvement circulaire", font_size=24),
                Text("4. Oscillations", font_size=24),
            )
            .arrange(DOWN, buff=0.6, aligned_edge=LEFT)
            .shift(LEFT * 2)
        )

        for prob in problems:
            self.play(Write(prob))
            self.wait(0.3)

        self.wait(2)


class ElectricFieldDetailed(Scene):
    def construct(self):
        title = Text("Champ Électrique: Détails", font_size=32)
        self.play(Write(title))
        self.wait()
        self.play(title.animate.to_edge(UP))

        formula = MathTex(r"\vec{E} = \frac{F}{q} = \frac{kQ}{r^2}", font_size=36)

        charge = Dot(ORIGIN, color=RED, radius=0.3)
        charge_label = MathTex("+Q", font_size=28).next_to(charge, UP)

        self.play(Write(formula))
        self.wait()
        self.play(formula.animate.to_edge(DOWN))
        self.play(Create(charge), Write(charge_label))

        # Field vectors
        for angle in [0, 45, 90, 135, 180, 225, 270, 315]:
            vec = Arrow(
                ORIGIN,
                np.array([np.cos(np.radians(angle)), np.sin(np.radians(angle)), 0]) * 2,
                color=YELLOW,
                buff=0.3,
            )
            self.play(Create(vec), run_time=0.3)

        self.wait(2)
