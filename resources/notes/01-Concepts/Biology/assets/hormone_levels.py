from manim import *

class HormoneLevels(Scene):
    def construct(self):
        title = Text("Menstrual Cycle: Hormone Levels", font_size=18).to_edge(UP)
        self.add(title)
        
        axes = Axes(
            x_range=[0, 28, 7],
            y_range=[0, 100, 20],
            x_length=8,
            y_length=4,
            axis_config={"include_tip": True}
        )
        
        x_label = axes.get_x_axis_label("Day")
        y_label = axes.get_y_axis_label("Level")
        
        # FSH curve
        fsh = axes.plot(lambda d: 20 + 60*np.exp(-(d-14)**2/20), x_range=[1, 28], color=BLUE)
        
        # LH curve
        lh = axes.plot(lambda d: 15 + 70*np.exp(-(d-14)**2/10), x_range=[1, 28], color=GREEN)
        
        # Estrogen curve
        est = axes.plot(lambda d: 20 + 60*np.exp(-(d-14)**2/15), x_range=[1, 28], color=RED)
        
        fsh_label = Text("FSH", color=BLUE, font_size=14).to_corner(UL).shift(RIGHT*2)
        lh_label = Text("LH", color=GREEN, font_size=14).next_to(fsh_label, DOWN)
        est_label = Text("Estrogen", color=RED, font_size=14).next_to(lh_label, DOWN)
        
        self.add(axes, x_label, y_label, title)
        self.play(Create(fsh), run_time=1)
        self.play(Create(lh), run_time=1)
        self.play(Create(est), run_time=1)
        self.add(fsh_label, lh_label, est_label)
        self.wait(3)
