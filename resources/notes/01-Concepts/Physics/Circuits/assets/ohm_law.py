from manim import *

class OhmLaw(Scene):
    def construct(self):
        title = Text("Ohm's Law: V = IR", font_size=24).to_edge(UP)
        self.add(title)
        
        axes = Axes(
            x_range=[0, 5, 1],
            y_range=[0, 10, 2],
            x_length=7,
            y_length=4,
            axis_config={"include_tip": True}
        )
        
        x_label = axes.get_x_axis_label("I (A)")
        y_label = axes.get_y_axis_label("V (V)")
        
        # V = IR with R = 2
        def voltage(i):
            return 2 * i
        
        graph = axes.plot(voltage, x_range=[0, 4], color=BLUE)
        
        self.add(axes, x_label, y_label, title)
        self.play(Create(graph), run_time=2)
        
        eq = MathTex("V = IR", font_size=36).to_corner(DR)
        self.play(Write(eq))
        self.wait(3)
