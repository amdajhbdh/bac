from manim import *

class FaradayLaw(Scene):
    def construct(self):
        title = Text("Faraday's Law: emf = -dPhi/dt", font_size=22).to_edge(UP)
        self.add(title)
        
        # Magnetic flux graph
        axes = Axes(
            x_range=[0, 4, 1],
            y_range=[0, 2, 0.5],
            x_length=7,
            y_length=3,
            axis_config={"include_tip": True}
        )
        
        x_label = axes.get_x_axis_label("t")
        y_label = axes.get_y_axis_label("Phi")
        
        # Magnetic flux increasing
        flux = axes.plot(lambda t: 0.5*t, x_range=[0, 4], color=BLUE)
        
        # Induced emf (negative derivative = constant negative)
        axes2 = Axes(
            x_range=[0, 4, 1],
            y_range=[-1, 0.5, 0.5],
            x_length=7,
            y_length=2.5,
            axis_config={"include_tip": True}
        ).shift(DOWN * 1.5)
        
        emf = axes2.plot(lambda t: -0.5, x_range=[0, 4], color=RED)
        
        self.add(axes, x_label, y_label, flux, title)
        self.add(axes2)
        self.play(Create(flux), run_time=2)
        self.play(Create(emf), run_time=1)
        
        eq = MathTex(r"\mathcal{E} = -\frac{d\Phi_B}{dt}", font_size=32).to_corner(DR)
        self.play(Write(eq))
        self.wait(3)
