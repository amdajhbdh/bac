#!/usr/bin/env python3
"""
Vault MCP Server - Model Context Protocol server for Study Vault
Provides tools for vault operations with semantic understanding
"""
import json
import sqlite3
from pathlib import Path
from typing import List, Dict, Any

VAULT = Path("/home/med/Documents/bac/resources/notes/Study-Vault")
DB = VAULT / ".vault-index.db"

class VaultMCPServer:
    def __init__(self):
        self.vault_path = VAULT
        self.db_path = DB
    
    def search(self, query: str, limit: int = 10, filters: Dict = None) -> List[Dict]:
        """Full-text search with optional filters"""
        conn = sqlite3.connect(self.db_path)
        c = conn.cursor()
        
        sql = "SELECT path, title FROM fts_notes WHERE fts_notes MATCH ? "
        params = [query]
        
        if filters:
            if filters.get('subject'):
                sql += "AND subject = ? "
                params.append(filters['subject'])
            if filters.get('type'):
                sql += "AND path IN (SELECT path FROM notes WHERE type = ?) "
                params.append(filters['type'])
        
        sql += "ORDER BY rank LIMIT ?"
        params.append(limit)
        
        c.execute(sql, params)
        results = [{"path": p, "title": t} for p, t in c.fetchall()]
        conn.close()
        return results
    
    def get_note(self, path: str) -> Dict:
        """Get note content and metadata"""
        note_path = self.vault_path / path
        if not note_path.exists():
            return {"error": "Note not found"}
        
        content = note_path.read_text(encoding='utf-8')
        
        conn = sqlite3.connect(self.db_path)
        c = conn.cursor()
        c.execute("SELECT * FROM notes WHERE path = ?", (path,))
        meta = c.fetchone()
        conn.close()
        
        return {
            "path": path,
            "content": content,
            "metadata": meta if meta else {}
        }
    
    def find_related(self, path: str, limit: int = 5) -> List[Dict]:
        """Find related notes via links and tags"""
        conn = sqlite3.connect(self.db_path)
        c = conn.cursor()
        
        # Get linked notes
        c.execute("""
            SELECT DISTINCT target_path as path, 'link' as relation
            FROM links WHERE source_path = ?
            UNION
            SELECT DISTINCT source_path as path, 'backlink' as relation
            FROM links WHERE target_path = ?
            LIMIT ?
        """, (path, path, limit))
        
        results = [{"path": p, "relation": r} for p, r in c.fetchall()]
        conn.close()
        return results
    
    def get_stats(self) -> Dict:
        """Get vault statistics"""
        conn = sqlite3.connect(self.db_path)
        c = conn.cursor()
        
        stats = {}
        c.execute("SELECT COUNT(*) FROM notes")
        stats['total_notes'] = c.fetchone()[0]
        
        c.execute("SELECT COUNT(*) FROM links")
        stats['total_links'] = c.fetchone()[0]
        
        c.execute("SELECT COUNT(*) FROM tags")
        stats['total_tags'] = c.fetchone()[0]
        
        c.execute("SELECT COUNT(*) FROM formulas")
        stats['total_formulas'] = c.fetchone()[0]
        
        c.execute("SELECT subject, COUNT(*) FROM notes WHERE subject IS NOT NULL GROUP BY subject")
        stats['by_subject'] = dict(c.fetchall())
        
        conn.close()
        return stats
    
    def find_orphans(self) -> List[str]:
        """Find notes with no links"""
        conn = sqlite3.connect(self.db_path)
        c = conn.cursor()
        c.execute("""
            SELECT path FROM notes 
            WHERE path NOT IN (SELECT DISTINCT source_path FROM links)
            AND path NOT IN (SELECT DISTINCT target_path FROM links)
            AND type != 'moc'
        """)
        orphans = [r[0] for r in c.fetchall()]
        conn.close()
        return orphans
    
    def search_formulas(self, query: str, limit: int = 10) -> List[Dict]:
        """Search mathematical formulas"""
        conn = sqlite3.connect(self.db_path)
        c = conn.cursor()
        c.execute("""
            SELECT DISTINCT path, formula FROM formulas 
            WHERE formula LIKE ? LIMIT ?
        """, (f'%{query}%', limit))
        results = [{"path": p, "formula": f} for p, f in c.fetchall()]
        conn.close()
        return results
    
    def rebuild_index(self) -> Dict:
        """Rebuild search index"""
        import subprocess
        result = subprocess.run(
            ['python3', str(VAULT / '../../../vault_indexer.py')],
            capture_output=True,
            text=True
        )
        return {
            "success": result.returncode == 0,
            "output": result.stdout
        }

# MCP Protocol Handler
def handle_mcp_request(request: Dict) -> Dict:
    """Handle MCP protocol requests"""
    server = VaultMCPServer()
    method = request.get('method')
    params = request.get('params', {})
    
    handlers = {
        'vault/search': lambda: server.search(params.get('query'), params.get('limit', 10)),
        'vault/get': lambda: server.get_note(params.get('path')),
        'vault/related': lambda: server.find_related(params.get('path'), params.get('limit', 5)),
        'vault/stats': lambda: server.get_stats(),
        'vault/orphans': lambda: server.find_orphans(),
        'vault/formulas': lambda: server.search_formulas(params.get('query'), params.get('limit', 10)),
        'vault/rebuild': lambda: server.rebuild_index()
    }
    
    if method in handlers:
        try:
            result = handlers[method]()
            return {"result": result}
        except Exception as e:
            return {"error": str(e)}
    else:
        return {"error": f"Unknown method: {method}"}

if __name__ == "__main__":
    import sys
    
    if len(sys.argv) < 2:
        print("Vault MCP Server")
        print("\nAvailable methods:")
        print("  vault/search    - Search notes")
        print("  vault/get       - Get note content")
        print("  vault/related   - Find related notes")
        print("  vault/stats     - Get statistics")
        print("  vault/orphans   - Find orphaned notes")
        print("  vault/formulas  - Search formulas")
        print("  vault/rebuild   - Rebuild index")
        sys.exit(0)
    
    # Test mode
    if sys.argv[1] == "test":
        test_requests = [
            {"method": "vault/stats", "params": {}},
            {"method": "vault/search", "params": {"query": "complex", "limit": 3}},
            {"method": "vault/orphans", "params": {}}
        ]
        
        for req in test_requests:
            print(f"\n{'='*60}")
            print(f"Request: {req['method']}")
            response = handle_mcp_request(req)
            print(json.dumps(response, indent=2))
