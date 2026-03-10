from manim import *

class Scene1a_OvarianHormones(Scene):
    """Show zero ovarian hormones for M & N vs normal cycle"""
    def construct(self):
        title = Text("Exercise 1 - Part 1a: Ovarian Hormones", font_size=32)
        self.play(Write(title))
        self.wait()
        self.play(title.animate.to_edge(UP))
        
        # Problem statement
        problem = VGroup(
            Text("Two women M and N: absent menstruation", font_size=24),
            Text("Document 1: Ovarian hormone levels (30 days)", font_size=24)
        ).arrange(DOWN, aligned_edge=LEFT).next_to(title, DOWN, buff=0.5)
        self.play(Write(problem))
        self.wait(2)
        self.play(FadeOut(problem))
        
        # Create axes
        axes = Axes(
            x_range=[0, 30, 5],
            y_range=[0, 10, 2],
            x_length=10,
            y_length=5,
            axis_config={"include_tip": True},
            x_axis_config={"numbers_to_include": [0, 5, 10, 15, 20, 25, 30]},
            y_axis_config={"numbers_to_include": [0, 2, 4, 6, 8, 10]}
        ).shift(DOWN*0.5)
        
        x_label = axes.get_x_axis_label("Days", direction=RIGHT)
        y_label = axes.get_y_axis_label("Hormone Level", direction=UP)
        
        self.play(Create(axes), Write(x_label), Write(y_label))
        
        # Normal cycle (estrogen + progesterone)
        normal_estrogen = axes.plot(lambda x: 3 + 2*np.sin(2*np.pi*x/28), color=BLUE, x_range=[0, 30])
        normal_prog = axes.plot(lambda x: 2 + 3*np.sin(2*np.pi*(x-7)/28), color=GREEN, x_range=[0, 30])
        
        # M and N (zero hormones)
        m_line = axes.plot(lambda x: 0, color=RED, x_range=[0, 30])
        n_line = axes.plot(lambda x: 0, color=ORANGE, x_range=[0, 30])
        
        # Labels
        legend = VGroup(
            VGroup(Line(LEFT*0.3, RIGHT*0.3, color=BLUE), Text("Normal (Estrogen)", font_size=18)).arrange(RIGHT, buff=0.2),
            VGroup(Line(LEFT*0.3, RIGHT*0.3, color=GREEN), Text("Normal (Progesterone)", font_size=18)).arrange(RIGHT, buff=0.2),
            VGroup(Line(LEFT*0.3, RIGHT*0.3, color=RED), Text("Woman M = 0", font_size=18)).arrange(RIGHT, buff=0.2),
            VGroup(Line(LEFT*0.3, RIGHT*0.3, color=ORANGE), Text("Woman N = 0", font_size=18)).arrange(RIGHT, buff=0.2)
        ).arrange(DOWN, aligned_edge=LEFT, buff=0.2).to_corner(UR)
        
        self.play(Create(normal_estrogen), Create(normal_prog))
        self.play(Write(legend[0]), Write(legend[1]))
        self.wait()
        self.play(Create(m_line), Create(n_line))
        self.play(Write(legend[2]), Write(legend[3]))
        self.wait(2)

class Scene1b_NoMenstruation(Scene):
    """Explain why no menstruation"""
    def construct(self):
        title = Text("Exercise 1 - Part 1b: Why No Menstruation?", font_size=32)
        self.play(Write(title))
        self.wait()
        self.play(title.animate.to_edge(UP))
        
        # Reasoning chain
        step1 = Text("Ovarian hormones = 0 for 30 days", font_size=28, color=BLUE)
        arrow1 = Arrow(UP, DOWN, color=YELLOW)
        step2 = Text("No hormone → No endometrium development", font_size=28, color=RED)
        arrow2 = Arrow(UP, DOWN, color=YELLOW)
        step3 = Text("No menstruation (cycle stopped)", font_size=28, color=GREEN)
        
        chain = VGroup(step1, arrow1, step2, arrow2, step3).arrange(DOWN, buff=0.5)
        
        self.play(Write(step1))
        self.wait()
        self.play(Create(arrow1))
        self.play(Write(step2))
        self.wait()
        self.play(Create(arrow2))
        self.play(Write(step3))
        self.wait(2)
        
        # Conclusion box
        conclusion = VGroup(
            Text("Conclusion:", font_size=24, color=YELLOW),
            Text("Zero hormones justify absent menstruation", font_size=22)
        ).arrange(DOWN, buff=0.3)
        box = SurroundingRectangle(conclusion, color=YELLOW, buff=0.3)
        conclusion_group = VGroup(box, conclusion).to_edge(DOWN)
        
        self.play(FadeOut(chain))
        self.play(Create(box), Write(conclusion))
        self.wait(2)

class Scene2a_FSH_LH_Comparison(Scene):
    """Bar chart comparing FSH/LH levels"""
    def construct(self):
        title = Text("Exercise 1 - Part 2a: FSH & LH Comparison", font_size=32)
        self.play(Write(title))
        self.wait()
        self.play(title.animate.to_edge(UP))
        
        # Data table
        table_title = Text("Document 2: Hormone Levels (UI/L)", font_size=24)
        table = Table(
            [["Woman M", "92", "60"],
             ["Woman N", "4", "3"],
             ["Normal", "32", "30"]],
            col_labels=[Text(""), Text("FSH"), Text("LH")],
            include_outer_lines=True
        ).scale(0.6)
        
        table_group = VGroup(table_title, table).arrange(DOWN, buff=0.3).shift(UP*0.5)
        
        self.play(Write(table_title))
        self.play(Create(table))
        self.wait(2)
        self.play(table_group.animate.scale(0.7).to_corner(UL))
        
        # Bar chart
        axes = Axes(
            x_range=[0, 4, 1],
            y_range=[0, 100, 20],
            x_length=8,
            y_length=5,
            axis_config={"include_tip": False},
            y_axis_config={"numbers_to_include": [0, 20, 40, 60, 80, 100]}
        ).shift(DOWN*0.5)
        
        y_label = axes.get_y_axis_label("Hormone Level (UI/L)", direction=UP)
        self.play(Create(axes), Write(y_label))
        
        # FSH bars
        fsh_m = Rectangle(width=0.5, height=4.6, color=RED, fill_opacity=0.7).move_to(axes.c2p(0.7, 46))
        fsh_n = Rectangle(width=0.5, height=0.2, color=ORANGE, fill_opacity=0.7).move_to(axes.c2p(1.5, 2))
        fsh_normal = Rectangle(width=0.5, height=1.6, color=BLUE, fill_opacity=0.7).move_to(axes.c2p(2.3, 16))
        
        # LH bars
        lh_m = Rectangle(width=0.5, height=3.0, color=RED, fill_opacity=0.5).move_to(axes.c2p(1.0, 30))
        lh_n = Rectangle(width=0.5, height=0.15, color=ORANGE, fill_opacity=0.5).move_to(axes.c2p(1.8, 1.5))
        lh_normal = Rectangle(width=0.5, height=1.5, color=BLUE, fill_opacity=0.5).move_to(axes.c2p(2.6, 15))
        
        # Labels
        fsh_label = Text("FSH", font_size=20).next_to(axes.c2p(1.5, 0), DOWN)
        lh_label = Text("LH", font_size=20).next_to(axes.c2p(2.0, 0), DOWN)
        
        legend = VGroup(
            VGroup(Square(side_length=0.2, color=RED, fill_opacity=0.7), Text("Woman M", font_size=18)).arrange(RIGHT, buff=0.2),
            VGroup(Square(side_length=0.2, color=ORANGE, fill_opacity=0.7), Text("Woman N", font_size=18)).arrange(RIGHT, buff=0.2),
            VGroup(Square(side_length=0.2, color=BLUE, fill_opacity=0.7), Text("Normal", font_size=18)).arrange(RIGHT, buff=0.2)
        ).arrange(DOWN, aligned_edge=LEFT, buff=0.2).to_corner(DR)
        
        self.play(Create(fsh_m), Create(lh_m))
        self.play(Create(fsh_n), Create(lh_n))
        self.play(Create(fsh_normal), Create(lh_normal))
        self.play(Write(fsh_label), Write(lh_label), Write(legend))
        self.wait(2)
        
        # Analysis
        analysis = VGroup(
            Text("Woman M: FSH & LH very HIGH", font_size=22, color=RED),
            Text("Woman N: FSH & LH very LOW", font_size=22, color=ORANGE),
            Text("Normal: FSH=32, LH=30 (baseline)", font_size=22, color=BLUE)
        ).arrange(DOWN, aligned_edge=LEFT, buff=0.3).to_edge(DOWN)
        
        self.play(Write(analysis))
        self.wait(3)

class Scene2b_Hypotheses(Scene):
    """Two hypotheses for each woman"""
    def construct(self):
        title = Text("Exercise 1 - Part 2b: Hypotheses", font_size=32)
        self.play(Write(title))
        self.wait()
        self.play(title.animate.to_edge(UP))
        
        # Woman M (high FSH/LH)
        m_title = Text("Woman M: FSH=92, LH=60 (HIGH)", font_size=26, color=RED).shift(UP*2.5)
        m_hyp = VGroup(
            Text("Hypothesis 1: Ovarian failure", font_size=22),
            Text("  → Ovaries don't respond to FSH/LH", font_size=20),
            Text("  → Pituitary increases FSH/LH (no feedback)", font_size=20),
            Text("Hypothesis 2: Pituitary hyperactivity", font_size=22),
            Text("  → Excessive FSH/LH secretion", font_size=20)
        ).arrange(DOWN, aligned_edge=LEFT, buff=0.2).next_to(m_title, DOWN, buff=0.5)
        
        self.play(Write(m_title))
        for line in m_hyp:
            self.play(Write(line), run_time=0.8)
        self.wait(2)
        
        self.play(FadeOut(m_title), FadeOut(m_hyp))
        
        # Woman N (low FSH/LH)
        n_title = Text("Woman N: FSH=4, LH=3 (LOW)", font_size=26, color=ORANGE).shift(UP*2.5)
        n_hyp = VGroup(
            Text("Hypothesis 1: Pituitary insufficiency", font_size=22),
            Text("  → Pituitary doesn't secrete enough FSH/LH", font_size=20),
            Text("  → Ovaries not stimulated", font_size=20),
            Text("Hypothesis 2: Hypothalamus problem", font_size=22),
            Text("  → Low GnRH → Low FSH/LH", font_size=20)
        ).arrange(DOWN, aligned_edge=LEFT, buff=0.2).next_to(n_title, DOWN, buff=0.5)
        
        self.play(Write(n_title))
        for line in n_hyp:
            self.play(Write(line), run_time=0.8)
        self.wait(3)

class Scene3_OvarianStructure(Scene):
    """Ovarian ultrasound comparison"""
    def construct(self):
        title = Text("Exercise 1 - Part 3: Ovarian Ultrasound", font_size=32)
        self.play(Write(title))
        self.wait()
        self.play(title.animate.to_edge(UP))
        
        # Three ovaries
        normal_ovary = Circle(radius=1, color=BLUE).shift(LEFT*3.5)
        m_ovary = Circle(radius=1, color=RED).shift(ORIGIN)
        n_ovary = Circle(radius=1, color=ORANGE).shift(RIGHT*3.5)
        
        # Follicles in normal ovary
        follicles_normal = VGroup(*[Dot(radius=0.1, color=YELLOW).move_to(normal_ovary.point_at_angle(i*PI/4)) 
                                    for i in range(8)])
        
        # No follicles in M (ovarian failure)
        # Follicles in N (ovaries intact)
        follicles_n = VGroup(*[Dot(radius=0.1, color=YELLOW).move_to(n_ovary.point_at_angle(i*PI/4)) 
                               for i in range(8)])
        
        # Labels
        normal_label = Text("Normal", font_size=20, color=BLUE).next_to(normal_ovary, DOWN)
        m_label = Text("Woman M", font_size=20, color=RED).next_to(m_ovary, DOWN)
        n_label = Text("Woman N", font_size=20, color=ORANGE).next_to(n_ovary, DOWN)
        
        self.play(Create(normal_ovary), Create(m_ovary), Create(n_ovary))
        self.play(Write(normal_label), Write(m_label), Write(n_label))
        self.wait()
        
        self.play(Create(follicles_normal), Create(follicles_n))
        self.wait(2)
        
        # Analysis
        analysis = VGroup(
            Text("Woman M: NO follicles → Ovarian failure confirmed!", font_size=24, color=RED),
            Text("Woman N: Follicles present → Ovaries intact", font_size=24, color=ORANGE),
            Text("→ Problem must be pituitary or hypothalamus", font_size=22, color=ORANGE)
        ).arrange(DOWN, aligned_edge=LEFT, buff=0.3).to_edge(DOWN)
        
        self.play(Write(analysis))
        self.wait(3)

class Scene4_GnRH_Test(Scene):
    """GnRH injection test for Woman N"""
    def construct(self):
        title = Text("Exercise 1 - Part 4: GnRH Test (Woman N)", font_size=32)
        self.play(Write(title))
        self.wait()
        self.play(title.animate.to_edge(UP))
        
        # Test description
        test_desc = VGroup(
            Text("Inject 100 µg GnRH", font_size=24),
            Text("Measure FSH/LH response", font_size=24)
        ).arrange(DOWN, buff=0.3).shift(UP*2)
        
        self.play(Write(test_desc))
        self.wait()
        self.play(test_desc.animate.scale(0.7).to_corner(UL))
        
        # Response graph
        axes = Axes(
            x_range=[0, 120, 30],
            y_range=[0, 40, 10],
            x_length=9,
            y_length=5,
            axis_config={"include_tip": True},
            x_axis_config={"numbers_to_include": [0, 30, 60, 90, 120]}
        )
        
        x_label = axes.get_x_axis_label("Time (min)", direction=RIGHT)
        y_label = axes.get_y_axis_label("FSH/LH (UI/L)", direction=UP)
        
        self.play(Create(axes), Write(x_label), Write(y_label))
        
        # Normal response (pituitary works)
        normal_response = axes.plot(lambda x: 5 + 25*np.exp(-0.03*x) if x > 10 else 5, color=BLUE)
        
        # Woman N response (pituitary responds!)
        n_response = axes.plot(lambda x: 3 + 20*np.exp(-0.03*x) if x > 10 else 3, color=ORANGE)
        
        # Injection marker
        injection = Arrow(axes.c2p(10, 35), axes.c2p(10, 30), color=RED, buff=0)
        inj_label = Text("GnRH injection", font_size=18, color=RED).next_to(injection, UP)
        
        legend = VGroup(
            VGroup(Line(LEFT*0.3, RIGHT*0.3, color=BLUE), Text("Normal response", font_size=18)).arrange(RIGHT, buff=0.2),
            VGroup(Line(LEFT*0.3, RIGHT*0.3, color=ORANGE), Text("Woman N response", font_size=18)).arrange(RIGHT, buff=0.2)
        ).arrange(DOWN, aligned_edge=LEFT, buff=0.2).to_corner(UR)
        
        self.play(Create(injection), Write(inj_label))
        self.wait()
        self.play(Create(normal_response), Create(n_response))
        self.play(Write(legend))
        self.wait(2)
        
        # Conclusion
        conclusion = VGroup(
            Text("Woman N's pituitary RESPONDS to GnRH", font_size=24, color=GREEN),
            Text("→ Pituitary is functional", font_size=22),
            Text("→ Problem is HYPOTHALAMUS (low GnRH)", font_size=22, color=ORANGE)
        ).arrange(DOWN, aligned_edge=LEFT, buff=0.3)
        
        box = SurroundingRectangle(conclusion, color=GREEN, buff=0.3)
        conclusion_group = VGroup(box, conclusion).to_edge(DOWN)
        
        self.play(Create(box), Write(conclusion))
        self.wait(3)

class FinalSummary(Scene):
    """Complete exercise summary"""
    def construct(self):
        title = Text("Exercise 1 - Final Diagnosis", font_size=36, color=YELLOW)
        self.play(Write(title))
        self.wait()
        self.play(title.animate.to_edge(UP))
        
        summary = VGroup(
            Text("Woman M:", font_size=28, color=RED, weight=BOLD),
            Text("  • FSH=92, LH=60 (very high)", font_size=24),
            Text("  • No follicles in ovaries", font_size=24),
            Text("  • Diagnosis: OVARIAN FAILURE", font_size=26, color=RED),
            Text("", font_size=10),
            Text("Woman N:", font_size=28, color=ORANGE, weight=BOLD),
            Text("  • FSH=4, LH=3 (very low)", font_size=24),
            Text("  • Follicles present in ovaries", font_size=24),
            Text("  • Responds to GnRH injection", font_size=24),
            Text("  • Diagnosis: HYPOTHALAMUS INSUFFICIENCY", font_size=26, color=ORANGE)
        ).arrange(DOWN, aligned_edge=LEFT, buff=0.2)
        
        for line in summary:
            self.play(Write(line), run_time=0.6)
        
        self.wait(3)

# Print success message
print("✓ Bio_Ex1 created: 7 scenes, ~3 minutes total")
print("Scenes: 1a, 1b, 2a, 2b, 3, 4, Summary")
