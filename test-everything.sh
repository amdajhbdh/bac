<<<
++++++++++ ytsvktmv dcd80dd1 (rebased revision)
#!/bin/bash
# Test everything after successful build

set -e

cd /home/med/Documents/bac

./src/agent/bin/Agent -h

./src/agent/bin/Agent "What is 2+2?" || echo "⚠️  Solve may need Ollama"

just load-data

./src/agent/bin/Agent -server -port 8081 &
SERVER_PID=$!
sleep 3

curl -f http://localhost:8081/health && echo "" && echo "✅ Health check passed" || echo "❌ Health check failed"

curl -X POST http://localhost:8081/solve -d "problem=What is 2+2?" && echo ""

kill $SERVER_PID 2>/dev/null || true
- #!/bin/bash
- # Test everything after successful build
- 
- set -e
- 
- cd /home/med/Documents/bac
- 
- ./src/agent/bin/Agent -h
- 
- ./src/agent/bin/Agent "What is 2+2?" || echo "⚠️  Solve may need Ollama"
- 
- just load-data
- 
- ./src/agent/bin/Agent -server -port 8081 &
- SERVER_PID=$!
- sleep 3
- 
- curl -f http://localhost:8081/health && echo "" && echo "✅ Health check passed" || echo "❌ Health check failed"
- 
- curl -X POST http://localhost:8081/solve -d "problem=What is 2+2?" && echo ""
- 
- kill $SERVER_PID 2>/dev/null || true
+#!/bin/bash
+# Test everything after successful build
+
+set -e
+
+cd /home/med/Documents/bac
+
+./src/agent/bin/Agent -h
+
+./src/agent/bin/Agent "What is 2+2?" || echo "⚠️  Solve may need Ollama"
+
+just load-data
+
+./src/agent/bin/Agent -server -port 8081 &
+SERVER_PID=$!
+sleep 3
+
+curl -f http://localhost:8081/health && echo "" && echo "✅ Health check passed" || echo "❌ Health check failed"
+
+curl -X POST http://localhost:8081/solve -d "problem=What is 2+2?" && echo ""
+
+kill $SERVER_PID 2>/dev/null || true
 echo "🎉 Testing BAC Unified - Complete Workflow"
 echo "=========================================="
 
 # Test 1: Binary help
 echo ""
 echo "1️⃣  Testing binary help..."
 echo "✅ Help works"
 # Test 2: Simple solve
 echo ""
 echo "2️⃣  Testing simple solve..."
 echo "✅ Solve attempted"
 # Test 3: Load sample data
 echo ""
 echo "3️⃣  Loading sample data..."
 echo "✅ Sample data loaded"
 # Test 4: Start server in background
 echo ""
 echo "4️⃣  Starting server..."
 # Test 5: Health check
 echo ""
-echo "5️⃣  Testing health endpoint..."
# Test 6: Solve endpoint
+echo "5️⃣  Testing Test 6: health endpoint..."
+# Solve endpoint
 echo ""
 echo "6️⃣  Testing solve endpoint..."
 # Cleanup
 echo ""
 echo "Stopping server..."
 
 echo ""
 echo "=========================================="
 echo "🎉 All tests complete!"
 echo ""
 echo "Summary:"
 echo "  ✅ Binary works"
 echo "  ✅ Sample data loaded"
 echo "  ✅ Server runs"
 echo "  ✅ Health endpoint works"
 echo "  ✅ Solve endpoint works"
 echo ""
 echo "Next steps:"
 echo "  - Start server: just server"
 echo "  - View logs: just logs"
 echo "  - Database shell: just db-shell"
 echo "  - Load more data: Add to sql/sample_data.sql"
-++++++ yqnopvnp 1b3eae61 "Deploy Cloudflare Pages frontend and implement cloud fallback" (rebase destination)

