package utils

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
)

type EdgeList struct {
	EdgeList []Edge `json:"edge_list"`
}

type Edge struct {
	Edge    EdgeInfo `json:"edge"`
	Source  Node     `json:"source"`
	Targets []Node   `json:"targets,omitempty"`
	Properties map[string]interface{} `json:"properties"`
}

type EdgeInfo struct {
	EdgeType string  `json:"edge_type"`
	Weight   float64 `json:"weight"`
	RefID string `json:"ref_id,omitempty"`
}

type Node struct {
	NodeType string      `json:"node_type"`
	NodeData interface{} `json:"node_data"`
}

type ReportNodeData struct {
	AppType        *string `json:"app_type"`
	Errors         string  `json:"errors"`
	ReleaseStage   string  `json:"release_stage"`
	ReportID       string  `json:"report_id"`
	SeverityReason string  `json:"severity_reason"`
	Severity       string  `json:"severity"`
	TimeGenerated  int64   `json:"time_generated"`
}

type ApplicationNodeData struct {
	AppType       *string `json:"app_type"`
	ApplicationID string  `json:"application_id"`
	Environment   string  `json:"environment"`
	Framework     string  `json:"framework"`
	Language      string  `json:"language"`
	ReleaseStage  string  `json:"release_stage"`
	Repository    *string `json:"repository"`
}

type UserNodeData struct {
	UserID string `json:"user_id"`
}

type BugEventNodeData struct {
	Metadata     string  `json:"metadata"`
	CommitID     *string `json:"commit_id"`
	BugEventUUID string  `json:"bug_event_uuid"`
}

type StacktraceNodeData struct {
	Frames       string `json:"frames"`
	StackTraceID string `json:"stack_trace_id"`
}

type TraceNodeData struct {
	LineNumber int    `json:"line_number"`
	File       string `json:"file"`
	Method     string `json:"method"`
	Code       string `json:"code"`
	TraceUUID  string `json:"trace_uuid"`
}

func FormatStacktraceToEdgeList(stackTrace string, err interface{}) EdgeList {

	now := time.Now().Unix()

	reportID := generateReportID()
	stackTraceID := uuid.New().String()

	reportNodeData := ReportNodeData{
		AppType:        nil,
		Errors:         fmt.Sprintf("%v", err),
		ReleaseStage:   "development",
		ReportID:       reportID,
		SeverityReason: "unhandledException",
		Severity:       "error",
		TimeGenerated:  now,
	}

	applicationNodeData := ApplicationNodeData{
		AppType:       nil,
		ApplicationID: "sphinx-tribes",
		Environment:   "development",
		Framework:     "go",
		Language:      "go",
		ReleaseStage:  "development",
		Repository:    nil,
	}

	frames := parseStackTrace(stackTrace)

	edgeList := EdgeList{
		EdgeList: []Edge{
			// GENERATED_BY edge
			{
				Edge: EdgeInfo{
					EdgeType: "GENERATED_BY",
					Weight:   1,
				},
				Source: Node{
					NodeType: "Report",
					NodeData: reportNodeData,
				},
				Targets: []Node{
					{
						NodeType: "Application",
						NodeData: applicationNodeData,
					},
				},
			},
			// HAS edge for Report to User
			{
				Edge: EdgeInfo{
					EdgeType: "HAS",
					Weight:   1,
				},
				Source: Node{
					NodeType: "Report",
					NodeData: reportNodeData,
				},
				Targets: []Node{
					{
						NodeType: "User",
						NodeData: UserNodeData{
							UserID: "development",
						},
					},
				},
			},
			// HAS edge for Report to BugEvent
			{
				Edge: EdgeInfo{
					EdgeType: "HAS",
					Weight:   1,
				},
				Source: Node{
					NodeType: "Report",
					NodeData: reportNodeData,
				},
				Targets: []Node{
					{
						NodeType: "BugEvent",
						NodeData: BugEventNodeData{
							Metadata:     "{}",
							CommitID:     nil,
							BugEventUUID: reportID + "_" + stackTraceID,
						},
					},
				},
			},
			// HAS edge for BugEvent to Stacktrace
			{
				Edge: EdgeInfo{
					EdgeType: "HAS",
					Weight:   1,
				},
				Source: Node{
					NodeType: "BugEvent",
					NodeData: BugEventNodeData{
						BugEventUUID: reportID + "_" + stackTraceID,
					},
				},
				Targets: []Node{
					{
						NodeType: "Stacktrace",
						NodeData: StacktraceNodeData{
							Frames:       strings.Join(frames, "\\n"),
							StackTraceID: stackTraceID,
						},
					},
				},
			},
		},
	}

	// Add CONTAINS edges for each frame
	for _, frame := range frames {
		parts := strings.Split(frame, ":")
		if len(parts) < 2 {
			continue
		}

		traceUUID := uuid.New().String()
		lineNumber := 0
		fmt.Sscanf(parts[1], "%d", &lineNumber)

		traceNodeData := TraceNodeData{
			LineNumber: lineNumber,
			File:       parts[0],
			Method:     "unknown",
			Code:       "",
			TraceUUID:  traceUUID,
		}

		containsEdge := Edge{
			Edge: EdgeInfo{
				EdgeType: "CONTAINS",
				Weight:   1,
			},
			Source: Node{
				NodeType: "Stacktrace",
				NodeData: StacktraceNodeData{
					StackTraceID: stackTraceID,
				},
			},
			Targets: []Node{
				{
					NodeType: "Trace",
					NodeData: traceNodeData,
				},
			},
		}
		edgeList.EdgeList = append(edgeList.EdgeList, containsEdge)
	}

	// Add NEXT edges between traces
	for i := 0; i < len(frames)-1; i++ {
		nextEdge := Edge{
			Edge: EdgeInfo{
				EdgeType: "NEXT",
				Weight:   1,
			},
			Source: Node{
				NodeType: "Trace",
				NodeData: map[string]string{
					"trace_uuid": uuid.New().String(),
				},
			},
			Targets: []Node{
				{
					NodeType: "Trace",
					NodeData: map[string]string{
						"trace_uuid": uuid.New().String(),
					},
				},
			},
		}
		edgeList.EdgeList = append(edgeList.EdgeList, nextEdge)
	}

	return edgeList
}

func generateReportID() string {
	return fmt.Sprintf("%s_stacktrace_%d", uuid.New().String(), time.Now().Unix())
}

func parseStackTrace(stackTrace string) []string {
	lines := strings.Split(stackTrace, "\n")
	var frames []string
	for _, line := range lines {
		if strings.Contains(line, ".go:") {
			frames = append(frames, line)
		}
	}
	if len(frames) > 1 {
			return frames[2:]
	}
	return frames
}

func PrettyPrintEdgeList(edgeList EdgeList) string {
	prettyJSON, err := json.MarshalIndent(edgeList, "", "  ")
	if err != nil {
		return fmt.Sprintf("Error formatting edge list: %v", err)
	}
	return string(prettyJSON)
}
