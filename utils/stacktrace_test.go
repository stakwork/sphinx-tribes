package utils

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func generateTestStackTrace() string {
	return `goroutine 1 [running]:
testing.tRunner()
	/usr/local/go/src/testing/testing.go:1689 +0x1b0
runtime/debug.Stack()
	/usr/local/go/src/runtime/debug/stack.go:24 +0x65
main.testFunction()
	/path/to/main.go:10 +0x26
another.function()
	/path/to/another.go:15 +0x45`
}

func TestFormatStacktraceToEdgeList(t *testing.T) {
	testCases := []struct {
		name           string
		stackTrace     string
		err            interface{}
		expectedChecks func(t *testing.T, edgeList EdgeList)
	}{
		{
			name:       "Basic Error Scenario",
			stackTrace: generateTestStackTrace(),
			err:        fmt.Errorf("test error"),
			expectedChecks: func(t *testing.T, edgeList EdgeList) {
				assert.NotNil(t, edgeList)
				assert.Greater(t, len(edgeList.EdgeList), 0)

				// Check GENERATED_BY edge
				generatedByEdge := findEdgeByType(edgeList, "GENERATED_BY")
				assert.NotNil(t, generatedByEdge)
				assert.Equal(t, "Report", generatedByEdge.Source.NodeType)
				assert.Equal(t, "Application", generatedByEdge.Targets[0].NodeType)

				// Check HAS edges
				hasEdges := findEdgesByType(edgeList, "HAS")
				assert.Greater(t, len(hasEdges), 0)

				// Check CONTAINS edges
				containsEdges := findEdgesByType(edgeList, "CONTAINS")
				assert.Greater(t, len(containsEdges), 0)
				for _, edge := range containsEdges {
					assert.Equal(t, "Stacktrace", edge.Source.NodeType)
					assert.Equal(t, "Trace", edge.Targets[0].NodeType)
				}

				// Check NEXT edges
				nextEdges := findEdgesByType(edgeList, "NEXT")
				assert.GreaterOrEqual(t, len(nextEdges), 0)
			},
		},
		{
			name:       "Nil Error Scenario",
			stackTrace: generateTestStackTrace(),
			err:        nil,
			expectedChecks: func(t *testing.T, edgeList EdgeList) {
				assert.NotNil(t, edgeList)
				assert.Greater(t, len(edgeList.EdgeList), 0)
			},
		},
		{
			name:       "Complex Error Scenario",
			stackTrace: generateTestStackTrace(),
			err:        struct{ message string }{"complex error"},
			expectedChecks: func(t *testing.T, edgeList EdgeList) {
				assert.NotNil(t, edgeList)
				assert.Greater(t, len(edgeList.EdgeList), 0)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			edgeList := FormatStacktraceToEdgeList(tc.stackTrace, tc.err)
			tc.expectedChecks(t, edgeList)
		})
	}
}

func TestParseStackTrace(t *testing.T) {
	testCases := []struct {
		name           string
		stackTrace     string
		expectedFrames int
	}{
		{
			name:           "Normal Stack Trace",
			stackTrace:     generateTestStackTrace(),
			expectedFrames: 4,
		},
		{
			name:           "Empty Stack Trace",
			stackTrace:     "",
			expectedFrames: 0,
		},
		{
			name: "Stack Trace Without Go Files",
			stackTrace: `goroutine 1 [running]:
some random text
another random line`,
			expectedFrames: 0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			frames := parseStackTrace(tc.stackTrace)
			assert.Equal(t, tc.expectedFrames, len(frames))
		})
	}
}

func TestGenerateReportID(t *testing.T) {

	reportID1 := generateReportID()
	reportID2 := generateReportID()

	assert.NotEqual(t, reportID1, reportID2)
	assert.Contains(t, reportID1, "stacktrace")
	assert.Contains(t, reportID2, "stacktrace")
}

func findEdgeByType(edgeList EdgeList, edgeType string) *Edge {
	for _, edge := range edgeList.EdgeList {
		if edge.Edge.EdgeType == edgeType {
			return &edge
		}
	}
	return nil
}

func findEdgesByType(edgeList EdgeList, edgeType string) []Edge {
	var matchedEdges []Edge
	for _, edge := range edgeList.EdgeList {
		if edge.Edge.EdgeType == edgeType {
			matchedEdges = append(matchedEdges, edge)
		}
	}
	return matchedEdges
}

func BenchmarkFormatStacktraceToEdgeList(b *testing.B) {
	stackTrace := generateTestStackTrace()
	err := fmt.Errorf("benchmark error")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		FormatStacktraceToEdgeList(stackTrace, err)
	}
}

func TestLongStackTraceHandling(t *testing.T) {

	var longStackTraceBuilder strings.Builder
	for i := 0; i < 1000; i++ {
		longStackTraceBuilder.WriteString(fmt.Sprintf("goroutine %d [running]:\n", i))
		longStackTraceBuilder.WriteString(fmt.Sprintf("/path/to/long/stack/trace%d.go:%d +0x26\n", i, i))
	}
	longStackTrace := longStackTraceBuilder.String()

	edgeList := FormatStacktraceToEdgeList(longStackTrace, "Stress Test Error")
	assert.NotNil(t, edgeList)
	assert.Greater(t, len(edgeList.EdgeList), 0)
}

func TestConcurrentStackTraceFormatting(t *testing.T) {
	stackTrace := generateTestStackTrace()
	err := fmt.Errorf("concurrent error")

	concurrentRuns := 100

	results := make(chan EdgeList, concurrentRuns)

	for i := 0; i < concurrentRuns; i++ {
		go func() {
			results <- FormatStacktraceToEdgeList(stackTrace, err)
		}()
	}

	for i := 0; i < concurrentRuns; i++ {
		edgeList := <-results
		assert.NotNil(t, edgeList)
		assert.Greater(t, len(edgeList.EdgeList), 0)
	}
}
