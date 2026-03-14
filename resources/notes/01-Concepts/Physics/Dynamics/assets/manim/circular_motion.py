from manim import *

class CircularMotion(Scene):
    def construct(self):
        title = Text("Uniform Circular Motion", font_size=24).to_edge(UP)
        self.add(title)
        
        circle = Circle(radius=2, color=WHITE)
        center = Dot(point=ORIGIN, color=WHITE)
        
        radius_line = Line(ORIGIN, RIGHT * 2, color=BLUE)
        
        position_dot = Dot(color=RED).move_to(RIGHT * 2)
        
        force_arrow = Arrow(
            start=RIGHT * 2,
            end=ORIGIN,
            color=YELLOW,
            buff=0
        )
        
        self.add(circle, center, radius_line, position_dot, title)
        
        eq1 = MathTex(r"F_c = \frac{mv^2}{r}", font_size=32).to_corner(DL)
        eq2 = MathTex(r"a_c = \frac{v^2}{r}", font_size=32).to_corner(DR)
        
        self.play(GrowArrow(force_arrow))
        self.play(Write(eq1))
        self.play(Write(eq2))
        
        self.wait(3)
