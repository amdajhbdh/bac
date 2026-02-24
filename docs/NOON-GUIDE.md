# Noon Animation Library - AI Reference Guide

This document provides complete reference for generating Noon (Rust) animation code. AI agents should use this to generate visualizations for any BAC problem.

---

## Table of Contents

1. [Core Concepts](#core-concepts)
2. [Scene Setup](#scene-setup)
3. [Creating Objects](#creating-objects)
4. [Animations](#animations)
5. [Colors](#colors)
6. [Text](#text)
7. [Math Visualizations](#math-visualizations)
8. [Physics Visualizations](#physics-visualizations)
9. [Chemistry Visualizations](#chemistry-visualizations)
10. [Complete Examples](#complete-examples)

---

## Core Concepts

Noon is a Rust animation library built on nannou. Key concepts:

- **Scene**: The main container for all animations
- **Objects**: Shapes, text, lines that can be animated
- **Animations**: Actions applied to objects (show, move, fade, etc.)
- **Builders**: Fluent API for configuring objects

---

## Scene Setup

### Basic Template

```rust
use noon::prelude::*;

fn scene(win_rect: Rect) -> Scene {
    let mut scene = Scene::new(win_rect);
    
    // Your objects and animations here
    
    scene
}

fn main() {
    noon::app(model).update(update).view(view).run();
}

fn model<'a>(app: &App) -> Scene {
    app.new_window()
        .size(1920, 1080)
        .title("Your Title")
        .view(view)
        .build()
        .unwrap();
    scene(app.window_rect())
}

fn update(app: &App, scene: &mut Scene, _update: Update) {
    scene.update(app.time, app.window_rect());
}

fn view(app: &App, scene: &Scene, frame: Frame) {
    let draw = app.draw();
    draw.background().color(BLACK);  // Use BLACK, WHITE, etc.
    scene.draw(draw.clone());
    draw.to_frame(app, &frame).unwrap();
}
```

### Window Size
- Full HD: `.size(1920, 1080)`
- Square: `.size(1080, 1080)`
- Custom: `.size(width, height)`

---

## Creating Objects

### 1. Circle

```rust
// Basic circle
let circle = scene.circle()
    .with_radius(0.5)                    // Radius (f32)
    .with_position(x, y)                 // Position (f32, f32)
    .with_color(Color::CYAN)              // Color
    .make();

// With all options
let circle = scene.circle()
    .with_radius(0.5)
    .with_position(2.0, 3.0)
    .with_color(Color::rgb(0.5, 0.8, 1.0))
    .with_angle(45.0)                     // Rotation in degrees
    .make();
```

### 2. Rectangle

```rust
// Basic rectangle
let rect = scene.rectangle()
    .with_size(width, height)             // (f32, f32)
    .with_position(x, y)                  // Center position
    .with_color(Color::BLUE)
    .make();

// Square
let square = scene.rectangle()
    .with_size(1.0, 1.0)
    .make();
```

### 3. Line

```rust
// Line from point to point
let line = scene.line()
    .from(x1, y1)                        // Start (f32, f32)
    .to(x2, y2)                          // End (f32, f32)
    .with_color(Color::WHITE)
    .with_thickness(2.0)                  // Line thickness
    .make();

// From another object
let line = scene.line()
    .from_point(other_object)             // Use another object's position
    .to(x, y)
    .make();
```

### 4. Dot

```rust
// Small dot/point
let dot = scene.dot()
    .with_position(x, y)
    .with_color(Color::RED)
    .with_radius(0.1)                     // Smaller than circles
    .make();
```

### 5. Arrow

```rust
let arrow = scene.arrow()
    .from(x1, y1)
    .to(x2, y2)
    .with_color(Color::YELLOW)
    .with_thickness(2.0)
    .make();
```

### 6. Triangle

```rust
let triangle = scene.triangle()
    .with_position(x, y)
    .with_size(base, height)
    .with_color(Color::GREEN)
    .with_angle(30.0)                     // Rotation
    .make();
```

---

## Animations

### Available Animation Methods

```rust
// Creation animations
scene.play(object.show_creation())     // Draw the object
scene.play(object.fade_in())          // Fade from invisible
scene.play(object.draw_stroke())      // Draw outline only

// Movement
scene.play(object.move_to(x, y))       // Move to position
scene.play(object.move_by(dx, dy))    // Move relative
scene.play(object.shift(x, y))        // Shift position

// Transformation
scene.play(object.scale(2.0))          // Scale by factor
scene.play(object.rotate(45.0))      // Rotate by degrees
scene.play(object.morph(other_object)) // Morph into another shape

// Visibility
scene.play(object.fade_out())         // Fade to invisible
scene.play(object.hide())            // Hide immediately

// Color
scene.play(object.set_color(Color::RED))
scene.play(object.set_thickness(3.0))
```

### Timing Control

```rust
// Duration
scene.play(object.show_creation()).run_time(1.5);  // 1.5 seconds

// Delay between animations (for sequential)
scene.play(object1.show_creation()).run_time(1.0);
scene.wait(0.5);  // Wait 0.5 seconds
scene.play(object2.show_creation()).run_time(1.0);

// Lag (delay between multiple animations)
scene.play(vec![obj1.show_creation(), obj2.show_creation()])
    .lag(0.3);  // 0.3 second delay between each

// Easing (smooth motion)
use noon::ease::EaseType;
scene.play(object.move_to(x, y))
    .rate_func(EaseType::Quad)     // Quadratic easing
    .rate_func(EaseType::Bounce)   // Bounce effect
    .rate_func(EaseType::Elastic)  // Elastic effect
    .rate_func(EaseType::Linear);  // Linear (constant speed)
```

### Waiting

```rust
scene.wait(1.0);              // Wait 1 second (no animation)
scene.wait_for(0.5);         // Wait for animation to complete + 0.5s
```

---

## Colors

### Predefined Colors

```rust
// Basic colors (use without ::)
BLACK, WHITE, RED, GREEN, BLUE, YELLOW, CYAN, MAGENTA
ORANGE, PURPLE, PINK, BROWN, GRAY

// Extended
DARK_RED, DARK_GREEN, DARK_BLUE
LIGHT_RED, LIGHT_GREEN, LIGHT_BLUE
TRANSPARENT
```

### Custom Colors

```rust
// RGB (0.0 - 1.0)
.with_color(Color::rgb(0.5, 0.8, 0.2))

// RGBA (with alpha/transparency)
.with_color(Color::rgba(1.0, 0.0, 0.0, 0.5))

// Using alpha on existing color
.with_color(Color::RED.with_alpha(0.5))  // 50% transparent
```

---

## Text

```rust
// Basic text
let text = scene.text("Hello")
    .with_position(x, y)
    .with_color(Color::WHITE)
    .with_font_size(24.0)
    .make();

// Math expressions (as plain text)
let formula = scene.text("x^2 + y^2 = r^2")
    .with_position(-5.0, 5.0)
    .with_color(Color::YELLOW)
    .with_font_size(32.0)
    .make();

// Special characters
// Note: Noon uses basic ASCII. For math, use approximations:
// ^ for superscript: x^2
// _ for subscript: x_1
// Greek letters as words: alpha, beta, theta, pi, etc. in lowercase
```

---

## Math Visualizations

### Coordinate System

```rust
// Draw axes
let x_axis = scene.line()
    .from(-8.0, 0.0)
    .to(8.0, 0.0)
    .with_color(Color::WHITE)
    .make();

let y_axis = scene.line()
    .from(0.0, -5.0)
    .to(0.0, 6.0)
    .with_color(Color::WHITE)
    .make();

// Labels
let x_label = scene.text("x")
    .with_position(7.5, -0.5)
    .with_color(Color::WHITE)
    .with_font_size(18.0)
    .make();

let y_label = scene.text("y")
    .with_position(0.3, 5.5)
    .with_color(Color::WHITE)
    .with_font_size(18.0)
    .make();
```

### Plotting Functions

```rust
// Plot y = f(x) using dots and lines
let mut prev_x = -5.0;
let mut prev_y = f(prev_x);
let mut prev_point = None;

for i in 1..100 {
    let x = -5.0 + (i as f32 / 100.0 * 10.0);  // -5 to 5
    let y = x * x;  // y = x^2 (example)
    
    let point = scene.circle()
        .with_radius(0.05)
        .with_position(x, y)
        .with_color(Color::CYAN)
        .make();
    
    // Connect to previous point
    if let Some(prev) = prev_point {
        let line = scene.line()
            .from(prev.0, prev.1)
            .to(x, y)
            .with_color(Color::CYAN)
            .make();
        scene.play(line.show_creation()).run_time(0.01);
    }
    
    prev_point = Some((x, y));
}
```

### Common Functions

```rust
// y = x^2
let y = x * x;

// y = x^3  
let y = x * x * x;

// y = sin(x)
let y = x.sin();

// y = cos(x)
let y = x.cos();

// y = e^x
let y = x.exp();

// y = ln(x)
let y = x.ln();

// y = sqrt(x)
let y = x.sqrt();

// y = |x|
let y = x.abs();

// y = 1/x
let y = 1.0 / x;

// y = a*x + b (linear)
let y = a * x + b;
```

### Shading Area Under Curve

```rust
// Shade area between curve and x-axis
for i in 0..=50 {
    let x = x_start + (i as f32 / 50.0) * (x_end - x_start);
    let y = f(x);
    
    let area_dot = scene.circle()
        .with_radius(0.1)
        .with_position(x, 0.0)  // Start from x-axis
        .with_color(Color::CYAN.with_alpha(0.3))
        .make();
    
    // Or create vertical line to curve
    let line = scene.line()
        .from(x, 0.0)
        .to(x, y)
        .with_color(Color::CYAN.with_alpha(0.2))
        .make();
}
```

### Geometry Shapes

```rust
// Triangle with vertices
let tri = scene.triangle()
    .with_position(center_x, center_y)
    .with_size(base, height)
    .with_angle(rotation)
    .make();

// Circle
let circ = scene.circle()
    .with_radius(radius)
    .with_position(cx, cy)
    .make();

// Arc
// (Draw as small line segments in a curve)
```

---

## Physics Visualizations

### Projectile Motion

```rust
// Parameters
let v0 = 30.0;      // Initial velocity (m/s)
let angle = 45.0;    // Launch angle (degrees)
let g = 9.8;        // Gravity

// Convert angle to radians
let rad = angle * 3.14159 / 180.0;

// Animate
let mut ball = scene.circle()
    .with_radius(0.3)
    .with_position(start_x, start_y)
    .with_color(Color::RED)
    .make();

let mut t = 0.0;
let dt = 0.05;

while t < total_time {
    let x = start_x + v0 * rad.cos() * t;
    let y = start_y + v0 * rad.sin() * t - 0.5 * g * t * t;
    
    if y < ground_y { break; }  // Hit ground
    
    scene.play(ball.move_to(x, y)).run_time(dt);
    t += dt;
}
```

### Vectors

```rust
// Vector as arrow
let vector = scene.arrow()
    .from(x1, y1)
    .to(x2, y2)
    .with_color(Color::YELLOW)
    .with_thickness(3.0)
    .make();

// Arrow head (small triangle at end)
// (No built-in - create small triangle)
```

### Circuit Diagrams

```rust
// Battery (rectangle with + and -)
let battery = scene.rectangle()
    .with_size(0.5, 1.5)
    .with_position(x, y)
    .with_color(Color::GRAY)
    .make();

// Resistor (zigzag - use multiple lines)
// Or rectangle
let resistor = scene.rectangle()
    .with_size(1.5, 0.4)
    .with_position(x, y)
    .with_color(Color::BROWN)
    .make();

// Capacitor (two parallel lines)
let cap1 = scene.line()
    .from(x1, y1)
    .to(x1, y2)
    .with_color(Color::BLUE)
    .make();

let cap2 = scene.line()
    .from(x2, y1)
    .to(x2, y2)
    .with_color(Color::BLUE)
    .make();

// Wire connections
let wire = scene.line()
    .from(x1, y1)
    .to(x2, y2)
    .with_color(Color::YELLOW)
    .with_thickness(2.0)
    .make();
```

---

## Chemistry Visualizations

### Atoms

```rust
// Carbon (gray, larger)
let carbon = scene.circle()
    .with_radius(0.35)
    .with_position(x, y)
    .with_color(Color::GRAY)
    .make();

// Oxygen (red)
let oxygen = scene.circle()
    .with_radius(0.3)
    .with_position(x, y)
    .with_color(Color::RED)
    .make();

// Hydrogen (white/small)
let hydrogen = scene.circle()
    .with_radius(0.2)
    .with_position(x, y)
    .with_color(Color::WHITE)
    .make();

// Nitrogen (blue)
let nitrogen = scene.circle()
    .with_radius(0.3)
    .with_position(x, y)
    .with_color(Color::BLUE)
    .make();
```

### Bonds

```rust
// Single bond (single line)
let bond = scene.line()
    .from(x1, y1)
    .to(x2, y2)
    .with_color(Color::WHITE)
    .with_thickness(3.0)
    .make();

// Double bond (two parallel lines)
// Draw two lines slightly offset

// Bond angle
// Use trigonometry to calculate positions
```

### Molecules

```rust
// Water (H2O)
// Oxygen at center, two H at angles ~104.5 degrees

// Ethanol (C2H5OH)
// Two C atoms bonded
// Each C has H atoms
// O-H group at end
```

### Reactions

```rust
// Show reactants on left, products on right
// Use arrows to show direction
let arrow = scene.arrow()
    .from(react_x, y)
    .to(product_x, y)
    .with_color(Color::WHITE)
    .make();

// Plus sign
let plus = scene.text("+")
    .with_position(x, y)
    .make();
```

---

## Complete Examples

### Example 1: Simple Function Plot

```rust
use noon::prelude::*;

fn scene(win_rect: Rect) -> Scene {
    let mut scene = Scene::new(win_rect);
    
    // Axes
    let x_axis = scene.line().from(-5.0, 0.0).to(5.0, 0.0).with_color(BLACK).make();
    let y_axis = scene.line().from(0.0, -5.0).to(0.0, 5.0).with_color(BLACK).make();
    scene.play(x_axis.show_creation()).run_time(1.0);
    scene.play(y_axis.show_creation()).run_time(1.0);
    
    // Plot y = x^2
    for i in 0..100 {
        let x = -5.0 + (i as f32 / 100.0 * 10.0);
        let y = x * x;
        
        if y < 5.0 {
            let dot = scene.dot().with_position(x, y).with_color(BLUE).make();
            scene.play(dot.fade_in()).run_time(0.01);
        }
    }
    
    scene
}

fn main() { noon::app(|_| scene(app.window_rect())).run(); }
```

### Example 2: Multiple Objects with Labels

```rust
use noon::prelude::*;

fn scene(win_rect: Rect) -> Scene {
    let mut scene = Scene::new(win_rect);
    
    // Title
    let title = scene.text("Graph of f(x) = x^2")
        .with_position(-4.0, 4.0)
        .with_color(BLUE)
        .with_font_size(32.0)
        .make();
    scene.play(title.show_creation()).run_time(1.5);
    
    // Point A at (2, 4)
    let point_a = scene.circle()
        .with_radius(0.2)
        .with_position(2.0, 4.0)
        .with_color(RED)
        .make();
    
    let label_a = scene.text("A(2, 4)")
        .with_position(2.3, 4.3)
        .with_color(RED)
        .with_font_size(18.0)
        .make();
    
    scene.play(point_a.show_creation()).run_time(0.5);
    scene.play(label_a.show_creation()).run_time(0.5);
    
    scene
}

fn main() { noon::app(|_| scene(app.window_rect())).run(); }
```

### Example 3: Physics Projectile

```rust
use noon::prelude::*;

fn scene(win_rect: Rect) -> Scene {
    let mut scene = Scene::new(win_rect);
    
    // Ground
    let ground = scene.rectangle()
        .with_size(20.0, 0.5)
        .with_position(0.0, -5.0)
        .with_color(GREEN)
        .make();
    scene.play(ground.fade_in()).run_time(0.5);
    
    // Ball
    let mut ball = scene.circle()
        .with_radius(0.3)
        .with_position(-7.0, -4.0)
        .with_color(RED)
        .make();
    scene.play(ball.show_creation()).run_time(0.5);
    
    // Animate projectile
    let v0 = 25.0;
    let angle = 45.0;
    let g = 9.8;
    let rad = angle * 3.14159 / 180.0;
    
    let mut t = 0.0;
    while t < 3.0 {
        let x = -7.0 + v0 * rad.cos() * t;
        let y = -4.0 + v0 * rad.sin() * t - 0.5 * g * t * t;
        
        if y < -4.75 { break; }
        
        scene.play(ball.move_to(x, y)).run_time(0.05);
        t += 0.05;
    }
    
    scene
}

fn main() { noon::app(|_| scene(app.window_rect())).run(); }
```

---

## Error Prevention

### Common Issues and Solutions

1. **Color not found**: Use predefined colors (BLACK, RED, etc.) or `Color::rgb(r, g, b)`

2. **Object not visible**: Check position is within window bounds (-8 to +8 horizontal, -6 to +6 vertical)

3. **Animation too fast/slow**: Adjust `.run_time()` - try 0.5-2.0 seconds

4. **Objects overlap**: Use `with_position()` to space them out

5. **Text too large/small**: Adjust `with_font_size()` - typical range 12-64

6. **Build errors**: Make sure to use `.make()` at the end of builder chains

---

## Coordinate System

- Center of screen: (0, 0)
- X increases to the right
- Y increases upward
- Typical range: X from -8 to +8, Y from -6 to +6 (for 1920x1080)
- Adjust based on your window size

---

## AI Generation Guidelines

When generating Noon code:

1. **Always use the template** with proper scene, model, update, view functions
2. **Use predefined colors** when possible (RED, BLUE, etc.)
3. **Keep animations simple** - use show_creation(), fade_in(), move_to()
4. **Set appropriate run_times** - 0.5-2.0 seconds typical
5. **Use loops for plotting** functions
6. **Add labels** for key points
7. **Include axes** for math problems
8. **Use comments** to explain sections

---

*Last Updated: 2026-02-23*
*For BAC Unified Project*
