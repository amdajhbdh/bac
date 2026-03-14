from manim import *

class AreaUnderCurve(Scene):
    def construct(self):
        title = Text("Area Under Curve (Definite Integral)", font_size=20).to_edge(UP)
        self.add(title)
        
        axes = Axes(
            x_range=[0, 4, 1],
            y_range=[0, 8, 2],
            x_length=7,
            y_length=4,
            axis_config={"include_tip": True}
        )
        
        x_label = axes.get_x_axis_label("x")
        y_label = axes.get_y_axis_label("f(x)")
        
        # f(x) = x^2
        curve = axes.plot(lambda x: x**2, x_range=[0, 3], color=BLUE)
        
        # Area under curve from x=1 to x=2
        area = axes.get_area(curve, [1, 2], color=GREEN, opacity=0.5)
        
        eq = MathTex(r"\int_1^2 x^2 dx = \left[\frac{x^3}{3}\right]_1^2 = \frac{7}{3}", font_size=24).to_corner(DR)
        
        self.add(axes, x_label, y_label, title)
        self.play(Create(curve), run_time=2)
        self.play(Create(area), run_time=1)
        self.play(Write(eq))
        self.wait(3)
