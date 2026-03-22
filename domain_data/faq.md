# Frequently Asked Questions

## What is this project?
This is a domain-specific RAG pipeline using hybrid BM25 + vector retrieval.

## How do I install dependencies?
Run: pip install -r requirements.txt

## How do I use the pipeline?
```python
from rag_pipeline import FastRAGPipeline
pipeline = FastRAGPipeline()
pipeline.initialize()
results = pipeline.query("your question")
```

## What embedding models are supported?
The default is all-MiniLM-L6-v2, but any sentence-transformers model works.
