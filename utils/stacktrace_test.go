package utils

import (
	"fmt"
	"strings"
	"testing"
	"encoding/json"
	"reflect"

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
			expectedFrames: 2,
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

func BenchmarkStacktrace(b *testing.B) {
  type containsCheck struct {
    shouldCheck bool
    expectedFile string
    expectedLine int
  }

  // Test Case 6: Very Large Stacktrace (Stress Case)
  // Build a stack trace with 2 header lines (non-matching) and 100 matching lines.
  var sb strings.Builder
  // Two header lines that do not contain ".go:"
  sb.WriteString("Header line 1\nHeader line 2\n")
  for i := 1; i <= 100; i++ {
    // each matching line: "file{i}.go:120"
    sb.WriteString(fmt.Sprintf("file%d.go:120", i))
    if i < 100 {
      sb.WriteString("\n")
    }
  }
  largeStackTrace := sb.String()

	testCases := []struct {
		name           string
		stackTrace     string
		err            interface{}
    expectedEdgesCount int
    expectedEffectiveFrames int
    expectedReportError string
    containsEdgeCheck containsCheck
  }{
    {
      name:               "Standard Multi-Frame Stacktrace",
      stackTrace:         generateTestStackTrace(),
      err:                fmt.Errorf("benchmark error"),
      // Matching lines: 4, but since len >1, effective frames = frames[2:] = 2.
      // Base edges = 4, plus 2 CONTAINS and 1 NEXT = 7.
      expectedEdgesCount:      7,
      expectedEffectiveFrames: 2,
      expectedReportError:     "benchmark error",
      containsEdgeCheck: containsCheck{
        shouldCheck:  true,
        // The first effective frame is "/path/to/main.go:10 +0x26"
        expectedFile: "\t/path/to/main.go",
        expectedLine: 10,
      },
    },
    {
      name:               "Single-Frame Stacktrace",
      stackTrace:         "file.go:123",
      err:                fmt.Errorf("single frame error"),
      // One matching frame yields effective frames = 1. Total edges = 4 + 1 = 5.
      expectedEdgesCount:      5,
      expectedEffectiveFrames: 1,
      expectedReportError:     "single frame error",
      containsEdgeCheck: containsCheck{
        shouldCheck:  true,
        expectedFile: "file.go",
        expectedLine: 123,
      },
    },
    {
      name:               "Stacktrace with No Valid Frames",
      stackTrace:         "this is not a valid stack trace",
      err:                fmt.Errorf("no frames"),
      // No matching frames. Total edges = 4.
      expectedEdgesCount:      4,
      expectedEffectiveFrames: 0,
      expectedReportError:     "no frames",
      containsEdgeCheck: containsCheck{
        shouldCheck: false,
      },
    },
    {
      name:               "Nil Error Value",
      stackTrace:         "file.go:50",
      err:                nil,
      // One matching frame. Effective frames = 1. Total edges = 5.
      expectedEdgesCount:      5,
      expectedEffectiveFrames: 1,
      // fmt.Sprintf("%v", nil) yields "<nil>"
      expectedReportError: "<nil>",
      containsEdgeCheck: containsCheck{
        shouldCheck:  true,
        expectedFile: "file.go",
        expectedLine: 50,
      },
    },
    {
      name:               "Frame with Non-Numeric Line Number",
      stackTrace:         "file.go:abc",
      err:                fmt.Errorf("bad line number"),
      // One matching frame. Effective frames = 1. Total edges = 5.
      expectedEdgesCount:      5,
      expectedEffectiveFrames: 1,
      expectedReportError:     "bad line number",
      containsEdgeCheck: containsCheck{
        shouldCheck:  true,
        expectedFile: "file.go",
        // On parsing failure, line number remains 0.
        expectedLine: 0,
      },
    },
    {
      name:               "Very Large Stacktrace (Stress Case)",
      stackTrace:         largeStackTrace,
      err:                fmt.Errorf("stress test"),
      // 100 matching lines minus 2 header discard equals 98 effective frames.
      // Total edges = 4 (base) + 98 (CONTAINS) + 97 (NEXT) = 199.
      expectedEdgesCount:      199,
      expectedEffectiveFrames: 98,
      expectedReportError:     "stress test",
      containsEdgeCheck: containsCheck{
        shouldCheck:  true,
        // First effective frame is "file3.go:120"
        expectedFile: "file3.go",
        expectedLine: 120,
      },
    },
    {
      name:               "Non-error Type as err Parameter",
      stackTrace:         "file.go:75",
      err:                123,
      // One matching frame. Total edges = 5.
      expectedEdgesCount:      5,
      expectedEffectiveFrames: 1,
      // fmt.Sprintf("%v", 123) yields "123"
      expectedReportError: "123",
      containsEdgeCheck: containsCheck{
        shouldCheck:  true,
        expectedFile: "file.go",
        expectedLine: 75,
      },
    },
    {
      name:               "Exactly Two Matching Frames",
      stackTrace:         "file1.go:100\nfile2.go:200",
      err:                fmt.Errorf("two frames only"),
      // Two matching frames are discarded (frames[2:]) so effective frames = 0.
      // Total edges = 4.
      expectedEdgesCount:      4,
      expectedEffectiveFrames: 0,
      expectedReportError:     "two frames only",
      containsEdgeCheck: containsCheck{
        shouldCheck: false,
      },
    },
    {
      name:               "Empty Stacktrace",
      stackTrace:         "",
      err:                fmt.Errorf("empty stacktrace"),
      // No matching frames so effective frames = 0. Total edges = 4.
      expectedEdgesCount:      4,
      expectedEffectiveFrames: 0,
      expectedReportError:     "empty stacktrace",
      containsEdgeCheck: containsCheck{
        shouldCheck: false,
      },
    },
  }

  b.ResetTimer()

  for _, tc := range testCases {
    tc := tc
    for i := 0; i < b.N; i++ {
      b.Run(tc.name, func(b *testing.B) {
        edgeList := FormatStacktraceToEdgeList(tc.stackTrace, tc.err)
        totalEdges := len(edgeList.EdgeList)
        assert.Equal(b, tc.expectedEdgesCount, totalEdges, "Total edges count mismatch for test case: %s", tc.name)

        // Validate the report node error inside the first base edge.
        if totalEdges >= 1 {
          reportEdge := edgeList.EdgeList[0]
          reportData, ok := reportEdge.Source.NodeData.(ReportNodeData)
          assert.True(b, ok, "Report node data type assertion failed in test case: %s", tc.name)
          assert.Equal(b, tc.expectedReportError, reportData.Errors, "ReportNodeData.Errors mismatch in test case: %s", tc.name)
        }

        // Determine effective frames using parseStackTrace.
        actualFrames := parseStackTrace(tc.stackTrace)
        assert.Equal(b, tc.expectedEffectiveFrames, len(actualFrames), "Effective frames count mismatch in test case: %s", tc.name)

        // Validate CONTAINS edges for effective frames.
        if tc.containsEdgeCheck.shouldCheck && tc.expectedEffectiveFrames > 0 {
          containsEdgeIndex := 4
          edge := edgeList.EdgeList[containsEdgeIndex]
          assert.Equal(b, "CONTAINS", edge.Edge.EdgeType, "First CONTAINS edge type mismatch in test case: %s", tc.name)
          if len(edge.Targets) > 0 {
            traceData, ok := edge.Targets[0].NodeData.(TraceNodeData)
            assert.True(b, ok, "Trace node data type assertion failed in test case: %s", tc.name)
            assert.Equal(b, tc.containsEdgeCheck.expectedFile, traceData.File, "TraceNodeData.File mismatch in test case: %s", tc.name)
            assert.Equal(b, tc.containsEdgeCheck.expectedLine, traceData.LineNumber, "TraceNodeData.LineNumber mismatch in test case: %s", tc.name)
          } else {
            b.Errorf("CONTAINS edge missing Targets in test case: %s", tc.name)
          }
        }

        // Validate NEXT edges if more than one effective frame exists.
        if tc.expectedEffectiveFrames > 1 {
          nextEdgesStart := 4 + tc.expectedEffectiveFrames
          nextEdgesCount := tc.expectedEffectiveFrames - 1
          for i := 0; i < nextEdgesCount; i++ {
            index := nextEdgesStart + i
            assert.Less(b, index, totalEdges, "NEXT edge index out of range in test case: %s", tc.name)
            nextEdge := edgeList.EdgeList[index]
            assert.Equal(b, "NEXT", nextEdge.Edge.EdgeType, "NEXT edge type mismatch at index %d in test case: %s", index, tc.name)
          }
        }
      })
    }
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
func TestPrettyPrintEdgeList(t *testing.T) {
	tests := []struct {
		name      string
		input     EdgeList
		validator func(t *testing.T, output string)
	}{
		{
			name: "Basic Functionality with a Single Edge",
			input: EdgeList{
				EdgeList: []Edge{
					{
						Edge:EdgeInfo{
							EdgeType: "CALLS",
							Weight: 1,
						},
						Source: Node{
							NodeType: "source_type",
							NodeData: "node-A",
						},
						Targets: []Node{
							{
							NodeType: "target_type",
							NodeData: "node-B",
							},
						
						},
						Properties: map[string]interface{}{
							"call_start": 100,
							"call_end": 150,
							"weight": 1,
	
						},
					},
				},
			},
			validator: func(t *testing.T, output string) {
				var result map[string]interface{}
				err := json.Unmarshal([]byte(output), &result)
				if err != nil {
					t.Fatalf("Unexpected error unmarshaling JSON: %v", err)
				}
				expected := map[string]interface{}{
					"edge_list": []interface{}{
						map[string]interface{}{

							"edge": map[string]interface{}{
								"edge_type": "CALLS",
								"weight": float64(1),
							},
							"source": map[string]interface{}{
								"node_type": "source_type",
								"node_data": "node-A",

							},
							"targets":[]interface{}{
								map[string]interface{}{
									"node_type": "target_type",
									"node_data":"node-B",
								},
							},
							"properties": map[string]interface{}{
								"call_start": float64(100),
								"call_end": float64(150),
								"weight": float64(1),
							},
						},
					},
				}
				if !reflect.DeepEqual(result, expected) {
					t.Errorf("Output structure mismatch.\nGot: %v\nExpected: %v", result, expected)
				}
			},
		},
		{
			name: "Edge Case: Empty EdgeList",
			input: EdgeList{
				EdgeList: []Edge{},
			},
			validator: func(t *testing.T, output string) {
				var result map[string]interface{}
				err := json.Unmarshal([]byte(output), &result)
				if err != nil {
					t.Fatalf("Unexpected error unmarshaling JSON: %v", err)
				}
				edges, ok := result["edge_list"].([]interface{})
				if !ok {
					t.Errorf("Expected key 'edge_list' to be an array, got: %v", result["edge_list"])
				}
				if len(edges) != 0 {
					t.Errorf("Expected an empty edge list, but got %d items", len(edges))
				}
			},
		},
		{
			name: "Basic Functionality with Multiple Edges",
			input: EdgeList{
				EdgeList: []Edge{
					{
						Edge: EdgeInfo{
							EdgeType: "CALLS",
							Weight: 1,
						},
						Source: Node{
							NodeType: "source_type",
							NodeData: "node-1",
						},
						Targets: []Node{
						{
							NodeType: "target_type",
							NodeData: "node-2",
						},

					},
				},
				{
					Edge: EdgeInfo{
						EdgeType: "CONTAINS",
						Weight: 2,
					},
					Source: Node{
						NodeType: "source_type",
						NodeData: "node-2",
					},
					Targets: []Node{
						{
							NodeType: "target_type",
							NodeData: "node-3",
						},
					},
					Properties: map[string]interface{}{
						"call_start": 200,
						"call_end": 250,
						"weight": 1,

					},
				},
			},
		},
		validator: func(t *testing.T, output string) {
			var result map[string]interface{}
			err := json.Unmarshal([]byte(output), &result)
			if err != nil {
				t.Fatalf("Unexpected error unmarshaling JSON: %v", err)
			}
		
			edges, ok := result["edge_list"].([]interface{})
			if !ok {
				t.Fatalf("Expected 'edge_list' to be an array")
			}
			if len(edges) != 2 {
				t.Fatalf("Expected 2 edges, got %d", len(edges))
			}
		
			first, ok := edges[0].(map[string]interface{})
			if !ok {
				t.Fatalf("Expected first edge to be a map")
			}
			edgeData1, ok := first["edge"].(map[string]interface{}) 
			if !ok {
				t.Fatalf("Expected 'edge' to be a map, got: %T", first["edge"])
			}
			if edgeData1["edge_type"] != "CALLS" {
				t.Errorf("Expected first edge type 'CALLS', got %v", edgeData1["edge_type"])
			}
		
			second, ok := edges[1].(map[string]interface{})
			if !ok {
				t.Fatalf("Expected second edge to be a map")
			}
			edgeData2, ok := second["edge"].(map[string]interface{})
			if !ok {
				t.Fatalf("Expected 'edge' to be a map, got: %T", second["edge"])
			}
			if edgeData2["edge_type"] != "CONTAINS" {
				t.Errorf("Expected second edge type 'CONTAINS', got %v", edgeData2["edge_type"])
			}
		},
	},
		{
			name: "Error Condition: Unserializable Field",
			input: EdgeList{
				EdgeList: []Edge{
					{
						Edge: EdgeInfo{
							EdgeType: "CONTAINS",
							Weight: 2,
						},
						
					
					Source: Node{
						NodeType: "source_type",
						NodeData: "node-X",

					},
					Targets: []Node{
						{
						NodeType: "target_type",
						NodeData: "Node-Y",
						},
					},
					Properties: map[string]interface{}{
						"call_start": 300,
						"call_end": 350,
						"weight": 1,
						"invalid": func(){},

					},

				},
			},
		},
			validator: func(t *testing.T, output string) {
				expectedPrefix := "Error formatting edge list:"
				if !strings.HasPrefix(output, expectedPrefix) {
					t.Errorf("Expected output to start with %q, but got: %q", expectedPrefix, output)
				}
			},
		},
		{
			name: "Performance and Scale: Very Large Edge List",
			input: func() EdgeList {
				edges := make([]Edge, 5000)
				for i := 1; i <= 5000; i++ {
					edges[i-1] = Edge{
						Edge: EdgeInfo{
							EdgeType: "CALLS",
							Weight: 1,
						},
						Properties: map[string]interface{}{
							"call_start": i,
							"call_end": i+50,
							"weight": 1,
						},
						Source: Node{
							NodeType: "large_source",
							NodeData: "node-large",
						},
						Targets: []Node{
							{
								NodeType: "large_target",
								NodeData: "node-large-target",
							},
						},
			
						
					}
				}
				return EdgeList{EdgeList: edges}
			}(),
			validator: func(t *testing.T, output string) {
				if len(output) == 0 {
					t.Errorf("Expected non-empty JSON output for large edge list")
				}
				trimmed := strings.TrimSpace(output)
				if !strings.HasPrefix(trimmed, "{") || !strings.HasSuffix(trimmed, "}") {
					t.Errorf("Expected output to start with '{' and end with '}', got: %q", trimmed)
				}
				count := strings.Count(output, `"edge_type"`)
				if count != 5000 {
					t.Errorf("Expected 5000 occurrences of \"edge_type\", but found %d", count)
				}
			},
		},
		{
			name: "Special Case: Special Characters within Properties",
			input: EdgeList{
				EdgeList: []Edge{
					{
						Edge: EdgeInfo{
							EdgeType: "CALLS",
							Weight: 1,
						},
						Properties: map[string]interface{}{
							"message": "This is a \"special\" message.\nIt contains newlines, \t tabs, and Unicode — ✓",
	
						},

					Source: Node{
						NodeType: "special_type",
						NodeData: "node-special",

					},
					Targets: []Node{
						{
						NodeType: "destination_type",
						NodeData: "node-destination",
						},
					},
					
				},
			},
		},
			validator: func(t *testing.T, output string) {
				var result map[string]interface{}
				err := json.Unmarshal([]byte(output), &result)
				if err != nil {
					t.Fatalf("Failed to unmarshal output: %v", err)
				}
				edges, ok := result["edge_list"].([]interface{})
				if !ok || len(edges) != 1 {
					t.Fatalf("Expected one edge in output, got: %v", result["edge_list"])
				}
				edge, ok := edges[0].(map[string]interface{})
				if !ok {
					t.Fatalf("Expected edge to be a map")
				}
				props, ok := edge["properties"].(map[string]interface{})
				if !ok {
					t.Fatalf("Expected 'properties' to be a map")
				}
				expectedMsg := "This is a \"special\" message.\nIt contains newlines, \t tabs, and Unicode — ✓"
				if props["message"] != expectedMsg {
					t.Errorf("Mismatch in message.\nGot: %q\nExpected: %q", props["message"], expectedMsg)
				}
			},
		},
		{
			name: "Special Case: Nil Properties Field",
			input: EdgeList{
				EdgeList: []Edge{
					{
						Edge: EdgeInfo{
							EdgeType: "NULL_TEST",
							Weight: 1,
							RefID: "edge-NULL",
						},
						Properties: map[string]interface{}(nil),
						Source: Node{
							NodeType: "null_source",
							NodeData: "node-NULL",
						},
						Targets: []Node{
							{
								NodeType: "null_target",
								NodeData: "node-NULL-target",

							},
							
						},
						
					},
				},
			},
			validator: func(t *testing.T, output string) {
				var result map[string]interface{}
				err := json.Unmarshal([]byte(output), &result)
				if err != nil {
					t.Fatalf("Error unmarshaling output: %v", err)
				}
				edges, ok := result["edge_list"].([]interface{})
				if !ok || len(edges) != 1 {
					t.Fatalf("Expected one edge, got: %v", result["edge_list"])
				}
				edge, ok := edges[0].(map[string]interface{})
				if !ok {
					t.Fatalf("Expected edge to be a map")
				}
				if v, present := edge["properties"]; !present || v != nil {
					t.Errorf("Expected properties to be null, got: %v", v)
				}
			},
		},
		{
			name: "Special Case: Nested Properties",
			input: EdgeList{
				EdgeList: []Edge{
					{
						Edge: EdgeInfo{
							EdgeType: "NESTED",
							Weight: 1,
							RefID: "edge-nested",

						},
						Properties: map[string]interface{}{
							"metadata" : map[string]interface{}{
								"version": "1.0",
								"flags": []string{"alpha", "beta"},
							},
							"status":"active",
						},
						Source: Node{
							NodeType: "nested_source",
							NodeData: "node-nested",

						},
						Targets: []Node{
							{
								NodeType: "nested_target",
								NodeData: "node-target",
							},
						},
					},
				},
			},
			validator: func(t *testing.T, output string) {
				var result map[string]interface{}
				err := json.Unmarshal([]byte(output), &result)
				if err != nil {
					t.Fatalf("Error unmarshaling output: %v", err)
				}
				edges, ok := result["edge_list"].([]interface{})
				if !ok || len(edges) != 1 {
					t.Fatalf("Expected one edge, got: %v", result["edge_list"])
				}
				edge, ok := edges[0].(map[string]interface{})
				if !ok {
					t.Fatalf("Expected edge to be a map")
				}
				props, ok := edge["properties"].(map[string]interface{})
				if !ok {
					t.Fatalf("Expected properties to be a map")
				}
				metadata, ok := props["metadata"].(map[string]interface{})
				if !ok {
					t.Fatalf("Expected metadata to be a map")
				}
				if metadata["version"] != "1.0" {
					t.Errorf("Expected metadata.version to be '1.0', got: %v", metadata["version"])
				}
				flags, ok := metadata["flags"].([]interface{})
				if !ok || len(flags) != 2 {
					t.Fatalf("Expected metadata.flags to be an array of 2 elements")
				}
				if flags[0] != "alpha" || flags[1] != "beta" {
					t.Errorf("Unexpected flags values: %v", flags)
				}
				if props["status"] != "active" {
					t.Errorf("Expected status 'active', got: %v", props["status"])
				}
			},
		},
	}

	for _, tc := range tests {
		tc := tc // capture range variable
		t.Run(tc.name, func(t *testing.T) {
			output := PrettyPrintEdgeList(tc.input)
			tc.validator(t, output)
		})
	}
}
