# System Dependencies for BAC Animation System

## Ubuntu 22.04 (Jammy)

### Required System Packages

```bash
# Update package list
sudo apt-get update

# Core build tools
sudo apt-get install -y \
    build-essential \
    git \
    curl \
    wget

# FFmpeg for video processing
sudo apt-get install -y \
    ffmpeg \
    libavdevice-dev \
    libavformat-dev \
    libavcodec-dev \
    libswscale-dev \
    libavutil-dev

# LaTeX for mathematical equations
sudo apt-get install -y \
    texlive-latex-base \
    texlive-pictures \
    texlive-science \
    texlive-latex-extra \
    latexmk

# Fonts
sudo apt-get install -y \
    fonts-dejavu-core \
    fonts-liberation \
    fonts-stix \
    fontconfig

# OpenGL and graphics
sudo apt-get install -y \
    xvfb \
    libgl1-mesa-glx \
    libglew-dev \
    libosmesa6-dev \
    mesa-common-dev

# Pango for text rendering
sudo apt-get install -y \
    pango1.0-tools \
    libpango1.0-dev \
    libcairo2-dev

# Python and development tools
sudo apt-get install -y \
    python3 \
    python3-dev \
    python3-pip \
    python3-venv \
    pkg-config

# Additional dependencies for manim
sudo apt-get install -y \
    libmagick++-dev \
    libpq-dev
```

## Docker Installation

### Using the Provided Dockerfile

```dockerfile
FROM ubuntu:22.04

# Set non-interactive mode
ENV DEBIAN_FRONTEND=noninteractive

# Install system dependencies
RUN apt-get update && apt-get install -y \
    build-essential \
    git \
    curl \
    wget \
    ffmpeg \
    libavdevice-dev \
    libavformat-dev \
    libavcodec-dev \
    libswscale-dev \
    libavutil-dev \
    texlive-latex-base \
    texlive-pictures \
    texlive-science \
    texlive-latex-extra \
    latexmk \
    fonts-dejavu-core \
    fonts-liberation \
    xvfb \
    libgl1-mesa-glx \
    libglew-dev \
    libosmesa6-dev \
    pango1.0-tools \
    libpango1.0-dev \
    libcairo2-dev \
    python3 \
    python3-dev \
    python3-pip \
    python3-venv \
    pkg-config \
    libmagick++-dev \
    && rm -rf /var/lib/apt/lists/*

# Create Python virtual environment
RUN python3 -m venv /opt/bac-animations
ENV PATH="/opt/bac-animations/bin:$PATH"

# Copy requirements and install Python packages
COPY requirements-animations.txt /tmp/
RUN pip install --no-cache-dir -r /tmp/requirements-animations.txt

# Set working directory
WORKDIR /app

# Default command
CMD ["/bin/bash"]
```

### Build and Run

```bash
# Build the image
docker build -t bac-animations:latest -f Dockerfile.anime .

# Run with volume mount
docker run -it \
    -v /path/to/output:/tmp/bac-animations \
    -v /path/to/scenes:/app/scenes \
    bac-animations:latest \
    python /app/bin/manim_renderer.py --scene-code "$(cat /app/scenes/scene.py)"
```

## Verification Commands

After installation, verify with:

```bash
# Check FFmpeg
ffmpeg -version

# Check LaTeX
latex --version

# Check Python
python3 --version

# Check Manim (after pip install)
python3 -c "import manim; print(manim.__version__)"

# Test xvfb
xvfb-run --auto-servernum manim --version

# Test rendering
xvfb-run -a manim -ql -c media/manim.cfg example_scenes.py ThreeDScene
```

## Output Directory

Ensure the output directory has proper permissions:

```bash
# Create and set permissions
sudo mkdir -p /tmp/bac-animations
sudo chown $USER:$USER /tmp/bac-animations
```

## Environment Variables

```bash
# Optional: Custom output directory
export ANIMATION_OUTPUT_DIR=/data/animations

# Optional: Custom quality
export ANIMATION_DEFAULT_QUALITY=medium

# Optional: Render timeout (seconds)
export ANIMATION_TIMEOUT=300
```
