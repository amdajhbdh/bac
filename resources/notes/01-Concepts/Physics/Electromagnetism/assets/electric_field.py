from manim import *

class ElectricField(Scene):
    def construct(self):
        title = Text("Electric Field: E = F/q", font_size=24).to_edge(UP)
        self.add(title)
        
        # Point charge
        charge = Circle(radius=0.3, color=RED, fill_opacity=0.8)
        plus = Text("+", font_size=24, color=WHITE).move_to(charge)
        
        # Field lines
        lines = VGroup()
        for angle in np.arange(0, 2*np.pi, np.pi/4):
            line = Line(ORIGIN, RIGHT * 2.5, color=YELLOW)
            line.rotate(angle)
            line.shift(RIGHT * 0.5)
            lines.add(line)
        
        # E field equation
        eq = MathTex(r"\vec{E} = k\frac{q}{r^2}", font_size=36).to_corner(DL)
        
        self.add(charge, plus, title)
        self.play(Create(lines), run_time=2)
        self.play(Write(eq))
        self.wait(3)
