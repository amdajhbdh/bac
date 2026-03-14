from manim import *

class ReactionRate(Scene):
    def construct(self):
        title = Text("Reaction Rate vs Concentration", font_size=20).to_edge(UP)
        self.add(title)
        
        axes = Axes(
            x_range=[0, 2, 0.5],
            y_range=[0, 2, 0.5],
            x_length=7,
            y_length=4,
            axis_config={"include_tip": True}
        )
        
        x_label = axes.get_x_axis_label("[A]")
        y_label = axes.get_y_axis_label("Rate")
        
        first_order = axes.plot(lambda c: c, x_range=[0, 2], color=BLUE)
        second_order = axes.plot(lambda c: c**2, x_range=[0, 2], color=RED)
        
        self.add(axes, x_label, y_label, title)
        self.play(Create(first_order), run_time=1)
        self.play(Create(second_order), run_time=1)
        self.wait(3)
