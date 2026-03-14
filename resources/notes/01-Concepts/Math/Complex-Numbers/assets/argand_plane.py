from manim import *

class ArgandPlane(Scene):
    def construct(self):
        title = Text("Argand Plane (Complex Numbers)", font_size=20).to_edge(UP)
        self.add(title)
        
        axes = Axes(
            x_range=[-3, 3, 1],
            y_range=[-3, 3, 1],
            x_length=6,
            y_length=6,
            axis_config={"include_tip": True}
        )
        
        x_label = axes.get_x_axis_label("Re")
        y_label = axes.get_y_axis_label("Im")
        
        # Complex number z = 2 + 1.5i
        point = axes.coords_to_point(2, 1.5)
        dot = Dot(point, color=RED, radius=0.1)
        
        # Line from origin to point
        line = Line(ORIGIN, point, color=YELLOW)
        
        # Angle arc
        arc = Arc(radius=0.8, start_angle=0, angle=0.64, color=GREEN)
        
        # Labels
        z_label = MathTex("z = 2 + 1.5i", font_size=24).next_to(dot, RIGHT + UP)
        r_label = MathTex("r = |z|", font_size=18).next_to(line, LEFT)
        theta_label = MathTex("theta = arg(z)", font_size=18).next_to(arc, RIGHT)
        
        self.add(axes, x_label, y_label, title)
        self.play(Create(line), Create(dot), run_time=2)
        self.play(Write(z_label), Write(r_label))
        self.wait(3)
