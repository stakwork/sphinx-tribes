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
