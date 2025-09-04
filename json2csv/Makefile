BIN_DIR=../bin
build: build-json2csv build-csv2json

build-json2csv:
	go build -o $(BIN_DIR)/json2csv json2csv.go

build-csv2json:
	go build -o $(BIN_DIR)/csv2json csv2json.go

generate-testdata: build-json2csv
	@echo "Generating testdata files..."
	@cat testdata/basic.json | json2csv > testdata/basic.csv
	@cat testdata/with-headers.json | json2csv -headers name,age > testdata/with-headers.csv
	@cat testdata/exclude.json | json2csv -exclude-headers secret > testdata/exclude.csv
	@cat testdata/merge-prefer-primary.json | json2csv -headers url -merge-headers url=path > testdata/merge-prefer-primary.csv
	@cat testdata/merge-fallback.json | json2csv -headers url -merge-headers url=path > testdata/merge-fallback.csv
	@cat testdata/append-all.json | json2csv -append-all-headers > testdata/append-all.csv
	@cat testdata/complex-merge.json | json2csv -headers time,requestID,method,url -merge-headers url=path > testdata/complex-merge.csv
	@cat testdata/all-flags.json | json2csv -headers time,requestID,method,url -append-all-headers -merge-headers url=path > testdata/all-flags.csv
	@cat testdata/merge-exclude.json | json2csv -merge-headers name=secret -exclude-headers secret > testdata/merge-exclude.csv
	@echo "✓ Generated all testdata CSV files"

test-json2csv: build-json2csv
	@echo "Running json2csv tests..."
	@passed=0; \
	failed=0; \
	\
	echo "Test 1: Basic CSV conversion"; \
	cat testdata/basic.json | json2csv > /tmp/test.csv 2>/dev/null; \
	if diff -q testdata/basic.csv /tmp/test.csv > /dev/null; then \
		echo "  ✓ Basic test passed"; \
		passed=$$((passed + 1)); \
	else \
		echo "  ✗ Basic test failed"; \
		echo "  Expected:"; cat testdata/basic.csv; \
		echo "  Got:"; cat /tmp/test.csv; \
		failed=$$((failed + 1)); \
	fi; \
	\
	echo "Test 2: Headers flag"; \
	cat testdata/with-headers.json | json2csv -headers name,age > /tmp/test.csv 2>/dev/null; \
	if diff -q testdata/with-headers.csv /tmp/test.csv > /dev/null; then \
		echo "  ✓ Headers test passed"; \
		passed=$$((passed + 1)); \
	else \
		echo "  ✗ Headers test failed"; \
		echo "  Expected:"; cat testdata/with-headers.csv; \
		echo "  Got:"; cat /tmp/test.csv; \
		failed=$$((failed + 1)); \
	fi; \
	\
	echo "Test 3: Exclude headers"; \
	cat testdata/exclude.json | json2csv -exclude-headers secret > /tmp/test.csv 2>/dev/null; \
	if diff -q testdata/exclude.csv /tmp/test.csv > /dev/null; then \
		echo "  ✓ Exclude headers test passed"; \
		passed=$$((passed + 1)); \
	else \
		echo "  ✗ Exclude headers test failed"; \
		echo "  Expected:"; cat testdata/exclude.csv; \
		echo "  Got:"; cat /tmp/test.csv; \
		failed=$$((failed + 1)); \
	fi; \
	\
	echo "Test 4: Merge headers (prefer primary)"; \
	cat testdata/merge-prefer-primary.json | json2csv -headers url -merge-headers url=path > /tmp/test.csv 2>/dev/null; \
	if diff -q testdata/merge-prefer-primary.csv /tmp/test.csv > /dev/null; then \
		echo "  ✓ Merge headers (prefer primary) test passed"; \
		passed=$$((passed + 1)); \
	else \
		echo "  ✗ Merge headers (prefer primary) test failed"; \
		echo "  Expected:"; cat testdata/merge-prefer-primary.csv; \
		echo "  Got:"; cat /tmp/test.csv; \
		failed=$$((failed + 1)); \
	fi; \
	\
	echo "Test 5: Merge headers (fallback)"; \
	cat testdata/merge-fallback.json | json2csv -headers url -merge-headers url=path > /tmp/test.csv 2>/dev/null; \
	if diff -q testdata/merge-fallback.csv /tmp/test.csv > /dev/null; then \
		echo "  ✓ Merge headers (fallback) test passed"; \
		passed=$$((passed + 1)); \
	else \
		echo "  ✗ Merge headers (fallback) test failed"; \
		echo "  Expected:"; cat testdata/merge-fallback.csv; \
		echo "  Got:"; cat /tmp/test.csv; \
		failed=$$((failed + 1)); \
	fi; \
	\
	echo "Test 6: Append all headers"; \
	cat testdata/append-all.json | json2csv -append-all-headers > /tmp/test.csv 2>/dev/null; \
	if diff -q testdata/append-all.csv /tmp/test.csv > /dev/null; then \
		echo "  ✓ Append all headers test passed"; \
		passed=$$((passed + 1)); \
	else \
		echo "  ✗ Append all headers test failed"; \
		echo "  Expected:"; cat testdata/append-all.csv; \
		echo "  Got:"; cat /tmp/test.csv; \
		failed=$$((failed + 1)); \
	fi; \
	\
	echo "Test 7: Complex merge scenario"; \
	cat testdata/complex-merge.json | json2csv -headers time,requestID,method,url -merge-headers url=path > /tmp/test.csv 2>/dev/null; \
	if diff -q testdata/complex-merge.csv /tmp/test.csv > /dev/null; then \
		echo "  ✓ Complex merge test passed"; \
		passed=$$((passed + 1)); \
	else \
		echo "  ✗ Complex merge test failed"; \
		echo "  Expected:"; cat testdata/complex-merge.csv; \
		echo "  Got:"; cat /tmp/test.csv; \
		failed=$$((failed + 1)); \
	fi; \
	\
	echo "Test 8: All flags combined (headers + append-all + merge)"; \
	cat testdata/all-flags.json | json2csv -headers time,requestID,method,url -append-all-headers -merge-headers url=path > /tmp/test.csv 2>/dev/null; \
	if diff -q testdata/all-flags.csv /tmp/test.csv > /dev/null; then \
		echo "  ✓ All flags test passed"; \
		passed=$$((passed + 1)); \
	else \
		echo "  ✗ All flags test failed"; \
		echo "  Expected:"; cat testdata/all-flags.csv; \
		echo "  Got:"; cat /tmp/test.csv; \
		failed=$$((failed + 1)); \
	fi; \
	\
	echo "Test 9: Merge with exclude (name=secret, exclude secret)"; \
	cat testdata/merge-exclude.json | json2csv -merge-headers name=secret -exclude-headers secret > /tmp/test.csv 2>/dev/null; \
	if diff -q testdata/merge-exclude.csv /tmp/test.csv > /dev/null; then \
		echo "  ✓ Merge with exclude test passed"; \
		passed=$$((passed + 1)); \
	else \
		echo "  ✗ Merge with exclude test failed"; \
		echo "  Expected:"; cat testdata/merge-exclude.csv; \
		echo "  Got:"; cat /tmp/test.csv; \
		failed=$$((failed + 1)); \
	fi; \
	\
	echo ""; \
	echo "Results: $$passed passed, $$failed failed"; \
	if [ $$failed -eq 0 ]; then \
		echo "All tests passed! ✅"; \
	else \
		echo "Some tests failed ❌"; \
		exit 1; \
	fi
