from manim import *

class DNAHelix(Scene):
    def construct(self):
        title = Text("DNA Double Helix", font_size=20).to_edge(UP)
        
        # Create two helix strands
        strand1 = []
        strand2 = []
        for i in range(12):
            t = i * 0.5
            y = i * 0.4 - 2
            x1 = np.cos(t) * 0.8
            x2 = np.cos(t + np.pi) * 0.8
            z = np.sin(t) * 0.8
            z2 = np.sin(t + np.pi) * 0.8
            
            dot1 = Dot([x1, y, 0], color=BLUE, radius=0.08)
            dot2 = Dot([x2, y, 0], color=RED, radius=0.08)
            strand1.append(dot1)
            strand2.append(dot2)
        
        # Base pairs (lines between strands)
        pairs = VGroup()
        for i in range(0, 12, 2):
            line = Line(strand1[i].get_center(), strand2[i].get_center(), color=GREEN)
            pairs.add(line)
        
        self.add(title, *strand1, *strand2, pairs)
        self.wait(3)
