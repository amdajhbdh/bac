# BAC Unified - Knowledge Base Optimization
# Enhances retrieval performance for Vector DB

# Optimization: Pre-warm HNSW index and tune search speed (ef_search)
# Higher ef_search = more accurate but slower
# Default is 40. We'll set it to 100 for better accuracy on first run.

def optimize_db [] {
    let db_url = "postgresql://neondb_owner:npg_ubkCLmerS03Z@ep-fragrant-violet-ai2ew4vx-pooler.c-4.us-east-1.aws.neon.tech/neondb"
    
    print "Optimizing Knowledge Base Index..."
    
    # Set ef_search for the current session to improve retrieval accuracy
    let sql = "SET hnsw.ef_search = 100;"
    try {
        ^psql $db_url -c $sql
        print "✓ HNSW ef_search optimized"
    } catch {
        print "✗ Failed to set HNSW parameters"
    }
}

# Implementation of Query Optimization Layer (Caching)
# We can use a simple JSON file to cache frequent search results
const SEARCH_CACHE = "data/search_cache.json"

def search_with_cache [query_hash: string] {
    if not (SEARCH_CACHE | path exists) {
        "{}" | save SEARCH_CACHE
        return null
    }
    
    let cache = (open SEARCH_CACHE)
    if ($cache | has $query_hash) {
        return ($cache | get $query_hash)
    }
    
    return null
}
