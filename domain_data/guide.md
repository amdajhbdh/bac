# User Guide

## Getting Started

1. Install the required packages
2. Place your documents in the domain_data folder
3. Initialize the pipeline with pipeline.initialize()
4. Query your documents using pipeline.query()

## Document Types Supported

- Markdown (.md)
- PDF (.pdf)
- Plain text (.txt)

## Configuration

You can customize:
- Embedding model via embedding_model parameter
- Chunk size and overlap
- Cache directory location
- Vector store persistence path
