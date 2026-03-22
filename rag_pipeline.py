import os
import hashlib
import json
from pathlib import Path
from typing import List, Optional, Dict, Any
from dataclasses import dataclass, field

from langchain_text_splitters import RecursiveCharacterTextSplitter
from langchain_community.document_loaders import (
    TextLoader,
    DirectoryLoader,
    PyPDFLoader,
)
from langchain_community.embeddings import SentenceTransformerEmbeddings
from langchain_chroma import Chroma
from langchain_core.documents import Document
from langchain_core.vectorstores import VectorStore

try:
    from rank_bm25 import BM25Okapi
except ImportError:
    import subprocess

    subprocess.run(
        ["pip", "install", "rank-bm25", "--break-system-packages"], check=True
    )
    from rank_bm25 import BM25Okapi


@dataclass
class RetrievalResult:
    document: Document
    score: float
    source: str


@dataclass
class CacheEntry:
    query: str
    results: List[RetrievalResult]
    timestamp: float


class QueryCache:
    def __init__(self, cache_dir: str = "./cache", max_entries: int = 1000):
        self.cache_dir = Path(cache_dir)
        self.cache_dir.mkdir(exist_ok=True)
        self.max_entries = max_entries
        self._memory_cache: Dict[str, CacheEntry] = {}

    def _get_cache_key(self, query: str) -> str:
        return hashlib.md5(query.lower().encode()).hexdigest()

    def get(self, query: str) -> Optional[List[RetrievalResult]]:
        key = self._get_cache_key(query)

        if key in self._memory_cache:
            return self._memory_cache[key].results

        cache_file = self.cache_dir / f"{key}.json"
        if cache_file.exists():
            try:
                data = json.loads(cache_file.read_text())
                results = [
                    RetrievalResult(
                        document=Document(
                            page_content=r["content"], metadata=r["metadata"]
                        ),
                        score=r["score"],
                        source=r["source"],
                    )
                    for r in data["results"]
                ]
                self._memory_cache[key] = CacheEntry(
                    query=data["query"], results=results, timestamp=data["timestamp"]
                )
                return results
            except Exception:
                return None
        return None

    def set(self, query: str, results: List[RetrievalResult]):
        import time

        key = self._get_cache_key(query)

        entry = CacheEntry(query=query, results=results, timestamp=time.time())
        self._memory_cache[key] = entry

        cache_file = self.cache_dir / f"{key}.json"
        data = {
            "query": query,
            "timestamp": entry.timestamp,
            "results": [
                {
                    "content": r.document.page_content,
                    "metadata": r.document.metadata,
                    "score": r.score,
                    "source": r.source,
                }
                for r in results
            ],
        }
        cache_file.write_text(json.dumps(data))


class DocumentFilter:
    PRIORITY_TYPES = {
        "faq": 1.0,
        "documentation": 0.9,
        "guide": 0.8,
        "reference": 0.7,
        "default": 0.5,
    }

    def __init__(self, priority_types: Optional[List[str]] = None):
        self.priority_types = priority_types or list(self.PRIORITY_TYPES.keys())

    def get_priority(self, doc_type: str) -> float:
        return self.PRIORITY_TYPES.get(doc_type, self.PRIORITY_TYPES["default"])

    def filter_and_boost(
        self, results: List[RetrievalResult], boost_types: Optional[List[str]] = None
    ) -> List[RetrievalResult]:
        boost_types = boost_types or self.priority_types[:3]

        for result in results:
            doc_type = result.document.metadata.get("type", "default")
            if doc_type in boost_types:
                result.score *= self.get_priority(doc_type)

        return sorted(results, key=lambda x: x.score, reverse=True)


class HybridRetriever:
    def __init__(
        self,
        embedding_model: str = "all-MiniLM-L6-v2",
        persist_directory: str = "./chroma_db",
        chunk_size: int = 1000,
        chunk_overlap: int = 200,
    ):
        self.persist_directory = persist_directory
        self.chunk_size = chunk_size
        self.chunk_overlap = chunk_overlap

        self.embeddings = SentenceTransformerEmbeddings(model_name=embedding_model)
        self.text_splitter = RecursiveCharacterTextSplitter(
            chunk_size=chunk_size, chunk_overlap=chunk_overlap
        )

        self.vector_store: Optional[Chroma] = None
        self.bm25: Optional[BM25Okapi] = None
        self.documents: List[Document] = []

    def load_documents(self, path: str, file_type: str = "auto"):
        if os.path.isfile(path):
            if path.endswith(".pdf"):
                loader = PyPDFLoader(path)
            else:
                loader = TextLoader(path)
            docs = loader.load()
        elif os.path.isdir(path):
            loader = DirectoryLoader(path, glob="**/*.md", loader_cls=TextLoader)
            docs = loader.load()

            pdf_loader = DirectoryLoader(path, glob="**/*.pdf", loader_cls=PyPDFLoader)
            pdf_docs = pdf_loader.load()
            docs.extend(pdf_docs)
        else:
            raise ValueError(f"Invalid path: {path}")

        self.documents = self.text_splitter.split_documents(docs)
        return self

    def build_index(self):
        if not self.documents:
            raise ValueError("No documents loaded. Call load_documents first.")

        texts = [doc.page_content for doc in self.documents]
        self.bm25 = BM25Okapi(texts)

        self.vector_store = Chroma.from_documents(
            documents=self.documents,
            embedding=self.embeddings,
            persist_directory=self.persist_directory,
        )

        return self

    def retrieve(
        self,
        query: str,
        k: int = 5,
        alpha: float = 0.5,
        doc_filter: Optional[DocumentFilter] = None,
    ) -> List[RetrievalResult]:
        results: Dict[str, RetrievalResult] = {}

        if self.vector_store:
            vector_results = self.vector_store.similarity_search_with_score(
                query, k=k * 2
            )
            for doc, score in vector_results:
                key = doc.page_content[:100]
                results[key] = RetrievalResult(
                    document=doc, score=score * (1 - alpha), source="vector"
                )

        if self.bm25:
            bm25_scores = self.bm25.get_scores(query.split()).tolist()
            top_bm25_indices = sorted(
                range(len(bm25_scores)), key=lambda i: bm25_scores[i], reverse=True
            )[: k * 2]

            for idx in top_bm25_indices:
                if idx < len(self.documents):
                    doc = self.documents[idx]
                    key = doc.page_content[:100]
                    if key in results:
                        results[key].score += bm25_scores[idx] * alpha
                        results[key].source = "hybrid"
                    else:
                        results[key] = RetrievalResult(
                            document=doc, score=bm25_scores[idx] * alpha, source="bm25"
                        )

        combined = list(results.values())

        if doc_filter:
            combined = doc_filter.filter_and_boost(combined)
        else:
            combined.sort(key=lambda x: x.score, reverse=True)

        return combined[:k]


class FastRAGPipeline:
    def __init__(
        self,
        data_path: str = "./domain_data",
        embedding_model: str = "all-MiniLM-L6-v2",
        persist_directory: str = "./chroma_db",
        cache_dir: str = "./cache",
    ):
        self.retriever = HybridRetriever(
            embedding_model=embedding_model, persist_directory=persist_directory
        )
        self.cache = QueryCache(cache_dir=cache_dir)
        self.filter = DocumentFilter()
        self.data_path = data_path

    def initialize(self, force_rebuild: bool = False):
        if force_rebuild or not os.path.exists(self.retriever.persist_directory):
            print(f"Loading documents from {self.data_path}...")
            self.retriever.load_documents(self.data_path)
            print("Building index...")
            self.retriever.build_index()
            print("Index ready!")
        else:
            print("Loading existing index...")
            self.retriever.vector_store = Chroma(
                persist_directory=self.retriever.persist_directory,
                embedding_function=self.retriever.embeddings,
            )

    def query(
        self, query: str, use_cache: bool = True, k: int = 5
    ) -> List[RetrievalResult]:
        if use_cache:
            cached = self.cache.get(query)
            if cached:
                print("Returning cached results")
                return cached

        results = self.retriever.retrieve(query, k=k, doc_filter=self.filter)

        if use_cache and results:
            self.cache.set(query, results)

        return results


if __name__ == "__main__":
    pipeline = FastRAGPipeline()

    domain_data = Path("./domain_data")
    if domain_data.exists():
        pipeline.initialize()

        print("\n--- Testing RAG Pipeline ---")
        test_queries = ["What is this about?", "How do I use this?"]

        for q in test_queries:
            print(f"\nQuery: {q}")
            results = pipeline.query(q)
            for i, r in enumerate(results[:3], 1):
                print(
                    f"  {i}. {r.document.page_content[:100]}... (score: {r.score:.3f})"
                )
