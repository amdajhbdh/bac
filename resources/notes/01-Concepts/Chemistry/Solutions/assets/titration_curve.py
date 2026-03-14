from manim import *

class TitrationCurve(Scene):
    def construct(self):
        title = Text("Titration Curve: Strong Acid - Strong Base", font_size=20).to_edge(UP)
        self.add(title)
        
        axes = Axes(
            x_range=[0, 50, 10],
            y_range=[0, 14, 2],
            x_length=8,
            y_length=4,
            axis_config={"include_tip": True}
        )
        
        x_label = axes.get_x_axis_label("Volume NaOH (mL)")
        y_label = axes.get_y_axis_label("pH")
        
        # Titration curve - strong acid/strong base
        def ph(v):
            if v < 40:
                return 1 + np.log10(v/(50-v+0.001)) if v < 49 else 14
            return 14
        
        # Simplified S-curve
        curve = axes.plot(lambda v: 1 + 13/(1+np.exp(-0.3*(v-25))), x_range=[0.1, 49.9], color=BLUE)
        
        # Equivalence point
        eq_point = Dot(axes.coords_to_point(25, 7), color=RED, radius=0.1)
        eq_label = Text("Equivalence point", font_size=14).next_to(eq_point, UP)
        
        self.add(axes, x_label, y_label, title)
        self.play(Create(curve), run_time=2)
        self.play(Create(eq_point), Write(eq_label))
        self.wait(3)
