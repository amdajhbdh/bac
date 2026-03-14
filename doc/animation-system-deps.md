# System Dependencies for BAC Animation System

## Ubuntu/Debian Installation

```bash
# Update package list
sudo apt-get update

# Install FFmpeg and video processing tools
sudo apt-get install -y \
    ffmpeg \
    libavdevice-dev \
    libavformat-dev \
    libavcodec-dev \
    libswscale-dev \
    libavutil-dev

# Install LaTeX for mathematical notation
sudo apt-get install -y \
    texlive-latex-base \
    texlive-pictures \
    texlive-science \
    texlive-latex-extra \
    fonts-dejavu-core \
    fonts-liberation

# Install OpenGL and rendering dependencies
sudo apt-get install -y \
    xvfb \
    libgl1-mesa-glx \
    libglew-dev \
    libosmesa6-dev \
    libopengl0

# Install Pango and Cairo for text rendering
sudo apt-get install -y \
    pango1.0-tools \
    libpango1.0-dev \
    libcairo2-dev

# Install Python and pip
sudo apt-get install -y \
    python3 \
    python3-pip \
    python3-venv
```

## Docker Installation (Recommended)

```dockerfile
FROM docker.io/manimcommunity/manim:latest

# Install additional dependencies
RUN apt-get update && apt-get install -y \
    ffmpeg \
    texlive-latex-base \
    texlive-pictures \
    texlive-science \
    texlive-latex-extra \
    fonts-dejavu-core \
    fonts-liberation \
    && rm -rf /var/lib/apt/lists/*

# Set working directory
WORKDIR /workspace

# Default command
CMD ["manim"]
```

## Verification Commands

After installation, verify with:

```bash
# Check FFmpeg
ffmpeg -version

# Check LaTeX
latex --version

# Check xvfb
xvfb-run --auto-servernum --server-args="-screen 0 1024x768x24" manim --version

# Check Python and pip
python3 --version
pip3 --version

# Test Manim import (requires venv)
python3 -c "import manim; print(manim.__version__)"
```

## Python Virtual Environment Setup

```bash
# Create virtual environment
python3 -m venv bac-animations

# Activate
source bac-animations/bin/activate

# Install dependencies
pip install -r requirements-animations.txt

# Verify installation
python -c "import manim; print(manim.__version__)"
python -c "import manim_physics; print('physics OK')"
python -c "import manim_chemistry; print('chemistry OK')"
```

## Environment Variables

Configure the animation system with these environment variables:

```bash
# Animation settings
export ANIMATION_OUTPUT_DIR="/tmp/bac-animations"
export ANIMATION_TIMEOUT="300"  # seconds

# Manim container (optional, uses Docker/Podman)
export MANIM_CONTAINER_IMAGE="docker.io/manimcommunity/manim:latest"
export MANIM_CONTAINER_NAME="bac-animator"

# Redis cache (optional)
export REDIS_URL="redis://localhost:6379"
```
