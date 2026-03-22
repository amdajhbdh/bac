# API Reference

## FastRAGPipeline

Main class for the RAG pipeline.

### Methods

- `initialize(force_rebuild: bool = False)` - Load documents and build index
- `query(query: str, use_cache: bool = True, k: int = 5)` - Retrieve relevant documents

## HybridRetriever

Handles document retrieval using BM25 and vector search.

### Parameters

- `embedding_model`: Name of sentence-transformers model
- `persist_directory`: Path to store Chroma database
- `chunk_size`: Size of text chunks
- `chunk_overlap`: Overlap between chunks
