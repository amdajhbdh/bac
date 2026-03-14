from manim import *

class HookesLaw(Scene):
    def construct(self):
        axes = Axes(
            x_range=[0, 0.5, 0.1],
            y_range=[0, 100, 20],
            x_length=6,
            y_length=4,
            axis_config={"include_tip": True}
        )
        
        x_label = axes.get_x_axis_label("x (m)")
        y_label = axes.get_y_axis_label("F (N)")
        
        k = 200
        def func(x):
            return k * x
        
        graph = axes.plot(func, x_range=[0, 0.4], color=BLUE)
        
        title = Text("Hooke's Law: F = kx", font_size=24).to_edge(UP)
        
        self.add(axes, x_label, y_label, title)
        self.play(Create(graph), run_time=2)
        self.wait(1)
        
        label = MathTex("F = -kx", font_size=36).next_to(graph, RIGHT)
        self.play(Write(label))
        self.wait(2)
