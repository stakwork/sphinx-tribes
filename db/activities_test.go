package db

import (
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestValidateActivity(t *testing.T) {
	teardownSuite := SetupSuite(t)
	defer teardownSuite(t)

	tests := []struct {
		name     string
		input    *Activity
		expected error
	}{
		{"Empty Content", &Activity{Content: "", AuthorRef: "valid_author", ContentType: FeatureCreation, Author: HumansAuthor, Workspace: "valid_workspace"}, ErrInvalidContent},
		{"Content Too Long", &Activity{Content: string(make([]byte, 10001)), AuthorRef: "valid_author", ContentType: FeatureCreation, Author: HumansAuthor, Workspace: "valid_workspace"}, ErrInvalidContent},
		{"Empty AuthorRef", &Activity{Content: "Valid", AuthorRef: "", ContentType: FeatureCreation, Author: HumansAuthor, Workspace: "valid_workspace"}, ErrInvalidAuthorRef},
		{"Empty Workspace", &Activity{Content: "Valid", AuthorRef: "valid_author", ContentType: FeatureCreation, Author: HumansAuthor, Workspace: ""}, ErrInvalidWorkspace},
		{"Invalid ContentType", &Activity{Content: "Valid", AuthorRef: "valid_author", ContentType: "invalid_type", Author: HumansAuthor, Workspace: "valid_workspace"}, ErrInvalidContentType},
		{"Invalid AuthorType", &Activity{Content: "Valid", AuthorRef: "valid_author", ContentType: FeatureCreation, Author: "invalid_author", Workspace: "valid_workspace"}, ErrInvalidAuthorType},
		{"Invalid Human Author Public Key", &Activity{Content: "Valid", AuthorRef: "short_key", ContentType: FeatureCreation, Author: HumansAuthor, Workspace: "valid_workspace"}, errors.New("invalid public key format for human author")},
		{"Invalid Hive Author UUID", &Activity{Content: "Valid", AuthorRef: "not-a-uuid", ContentType: FeatureCreation, Author: HiveAuthor, Workspace: "valid_workspace"}, errors.New("invalid UUID format for hive author")},
		{"Valid Human Author", &Activity{Content: "Valid", AuthorRef: "abcdefghijklmnopqrstuvwxyz123456", ContentType: FeatureCreation, Author: HumansAuthor, Workspace: "valid_workspace"}, nil},
		{"Valid Hive Author", &Activity{Content: "Valid", AuthorRef: uuid.NewString(), ContentType: FeatureCreation, Author: HiveAuthor, Workspace: "valid_workspace"}, nil},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := validateActivity(test.input)
			if test.expected == nil {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, test.expected.Error())
			}
		})
	}
}

func TestCreateActivity(t *testing.T) {
	teardownSuite := SetupSuite(t)
	defer teardownSuite(t)

	tests := []struct {
		name        string
		setup       func() *Activity
		expectError bool
		errorMsg    string
	}{
		{
			name: "Successfully create activity",
			setup: func() *Activity {
				return &Activity{
					ID:          uuid.Nil,
					ThreadID:    uuid.New(),
					Sequence:    1,
					ContentType: FeatureCreation,
					Content:     "Valid content",
					Workspace:   "workspace-1",
					Author:      HumansAuthor,
					AuthorRef:   "12345678901234567890123456789012",
				}
			},
			expectError: false,
		},
		{
			name: "Invalid activity: empty content",
			setup: func() *Activity {
				return &Activity{
					Content: "",
				}
			},
			expectError: true,
			errorMsg:    ErrInvalidContent.Error(),
		},
		{
			name: "Invalid activity: unknown content type",
			setup: func() *Activity {
				return &Activity{
					ContentType: "unknown_type",
					Content:     "Valid content",
					Author:      HumansAuthor,
					AuthorRef:   "12345678901234567890123456789012",
					Workspace:   "workspace-1",
				}
			},
			expectError: true,
			errorMsg:    ErrInvalidContentType.Error(),
		},
		{
			name: "Invalid activity: invalid human author ref",
			setup: func() *Activity {
				return &Activity{
					Content:     "Valid content",
					ContentType: FeatureCreation,
					Author:      HumansAuthor,
					AuthorRef:   "short",
					Workspace:   "workspace-1",
				}
			},
			expectError: true,
			errorMsg:    "invalid public key format for human author",
		},
		{
			name: "Invalid activity: invalid UUID for hive author",
			setup: func() *Activity {
				return &Activity{
					Content:     "Valid content",
					ContentType: FeatureCreation,
					Author:      HiveAuthor,
					AuthorRef:   "invalid-uuid",
					Workspace:   "workspace-1",
				}
			},
			expectError: true,
			errorMsg:    "invalid UUID format for hive author",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			activity := tt.setup()
			_, err := TestDB.CreateActivity(activity)

			if tt.expectError {
				assert.Error(t, err)
				assert.Equal(t, tt.errorMsg, err.Error())
			} else {
				assert.NoError(t, err)
				assert.NotEqual(t, uuid.Nil, activity.ID)
				assert.NotZero(t, activity.TimeCreated)
				assert.NotZero(t, activity.TimeUpdated)
			}
		})
	}
}

func TestUpdateActivity(t *testing.T) {
	teardownSuite := SetupSuite(t)
	defer teardownSuite(t)

	existingActivity := &Activity{
		ID:          uuid.New(),
		ThreadID:    uuid.New(),
		Sequence:    1,
		ContentType: FeatureCreation,
		Content:     "Original content",
		Workspace:   "workspace-1",
		Author:      HumansAuthor,
		AuthorRef:   "12345678901234567890123456789012",
		TimeCreated: time.Now(),
		TimeUpdated: time.Now(),
	}
	TestDB.db.Create(existingActivity)

	tests := []struct {
		name        string
		setup       func() *Activity
		expectError bool
		errorMsg    string
	}{
		{
			name: "Successfully update activity",
			setup: func() *Activity {
				return &Activity{
					ID:          existingActivity.ID,
					ThreadID:    existingActivity.ThreadID,
					Sequence:    existingActivity.Sequence,
					ContentType: FeatureCreation,
					Content:     "Updated content",
					Workspace:   "workspace-1",
					Author:      HumansAuthor,
					AuthorRef:   "12345678901234567890123456789012",
				}
			},
			expectError: false,
		},
		{
			name: "Fail to update non-existent activity",
			setup: func() *Activity {
				return &Activity{
					ID:          uuid.New(),
					ThreadID:    uuid.New(),
					Sequence:    1,
					ContentType: FeatureCreation,
					Content:     "New content",
					Workspace:   "workspace-1",
					Author:      HumansAuthor,
					AuthorRef:   "12345678901234567890123456789012",
				}
			},
			expectError: true,
			errorMsg:    "activity not found",
		},
		{
			name: "Fail to update thread_id",
			setup: func() *Activity {
				return &Activity{
					ID:          existingActivity.ID,
					ThreadID:    uuid.New(),
					Sequence:    existingActivity.Sequence,
					ContentType: FeatureCreation,
					Content:     "Updated content",
					Workspace:   "workspace-1",
					Author:      HumansAuthor,
					AuthorRef:   "12345678901234567890123456789012",
				}
			},
			expectError: true,
			errorMsg:    "thread_id cannot be modified",
		},
		{
			name: "Fail to update sequence",
			setup: func() *Activity {
				return &Activity{
					ID:          existingActivity.ID,
					ThreadID:    existingActivity.ThreadID,
					Sequence:    999,
					ContentType: FeatureCreation,
					Content:     "Updated content",
					Workspace:   "workspace-1",
					Author:      HumansAuthor,
					AuthorRef:   "12345678901234567890123456789012",
				}
			},
			expectError: true,
			errorMsg:    "sequence cannot be modified",
		},
		{
			name: "Invalid activity: empty content",
			setup: func() *Activity {
				return &Activity{
					ID:       existingActivity.ID,
					ThreadID: existingActivity.ThreadID,
					Sequence: existingActivity.Sequence,
					Content:  "",
				}
			},
			expectError: true,
			errorMsg:    ErrInvalidContent.Error(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			activity := tt.setup()
			_, err := TestDB.UpdateActivity(activity)

			if tt.expectError {
				assert.Error(t, err)
				assert.Equal(t, tt.errorMsg, err.Error())
			} else {
				assert.NoError(t, err)
				assert.NotEqual(t, existingActivity.TimeUpdated, activity.TimeUpdated)
			}
		})
	}
}

func TestGetActivity(t *testing.T) {
	teardownSuite := SetupSuite(t)
	defer teardownSuite(t)

	id := uuid.New()

	tests := []struct {
		name        string
		id          string
		setup       func()
		expectError bool
		errorMsg    string
	}{
		{
			name: "Successfully retrieve activity",
			id:   id.String(),
			setup: func() {
				activity := &Activity{
					ID:          id,
					ThreadID:    uuid.New(),
					Sequence:    1,
					ContentType: FeatureCreation,
					Content:     "Valid content",
					Workspace:   "workspace-1",
					Author:      HumansAuthor,
					AuthorRef:   "12345678901234567890123456789012",
				}
				TestDB.CreateActivity(activity)
			},
			expectError: false,
		},
		{
			name:        "Activity not found",
			id:          uuid.NewString(),
			setup:       func() {},
			expectError: true,
			errorMsg:    "record not found",
		},
		{
			name:        "Invalid UUID format",
			id:          "invalid-uuid",
			setup:       func() {},
			expectError: true,
			errorMsg:    "invalid activity ID format",
		},
		{
			name:        "Empty ID string",
			id:          "",
			setup:       func() {},
			expectError: true,
			errorMsg:    "invalid activity ID format",
		},
		{
			name:        "UUID with special characters",
			id:          uuid.NewString() + "@!#",
			setup:       func() {},
			expectError: true,
			errorMsg:    "invalid activity ID format",
		},
		{
			name:        "UUID with whitespace",
			id:          "  " + uuid.NewString() + "  ",
			setup:       func() {},
			expectError: true,
			errorMsg:    "invalid activity ID format",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			_, err := TestDB.GetActivity(tt.id)

			if tt.expectError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorMsg)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestGetActivitiesByThread(t *testing.T) {
	teardownSuite := SetupSuite(t)
	defer teardownSuite(t)

	validThreadID := uuid.New()
	invalidThreadID := "invalid-uuid"

	activities := []Activity{
		{ID: uuid.New(), ThreadID: validThreadID, Sequence: 1, Content: "First Activity"},
		{ID: uuid.New(), ThreadID: validThreadID, Sequence: 2, Content: "Second Activity"},
	}
	TestDB.db.Create(&activities)

	tests := []struct {
		name          string
		threadID      string
		expectError   bool
		expectedCount int
		errorMsg      string
	}{
		{
			name:          "Valid thread ID with activities",
			threadID:      validThreadID.String(),
			expectError:   false,
			expectedCount: 2,
		},
		{
			name:          "Valid thread ID with no activities",
			threadID:      uuid.New().String(),
			expectError:   false,
			expectedCount: 0,
		},
		{
			name:        "Invalid thread ID format",
			threadID:    invalidThreadID,
			expectError: true,
			errorMsg:    "invalid thread ID format",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := TestDB.GetActivitiesByThread(tt.threadID)

			if tt.expectError {
				assert.Error(t, err)
				assert.Equal(t, tt.errorMsg, err.Error())
			} else {
				assert.NoError(t, err)
				assert.Len(t, result, tt.expectedCount)
			}
		})
	}
}

func TestGetActivitiesByFeature(t *testing.T) {
	teardownSuite := SetupSuite(t)
	defer teardownSuite(t)

	tests := []struct {
		name        string
		setup       func() string
		expectError bool
		errorMsg    string
		expectCount int
	}{
		{
			name: "Valid feature UUID with activities",
			setup: func() string {
				featureUUID := uuid.New().String()
				activities := []Activity{
					{ID: uuid.New(), FeatureUUID: featureUUID, Content: "Activity 1", TimeCreated: time.Now().Add(-time.Hour)},
					{ID: uuid.New(), FeatureUUID: featureUUID, Content: "Activity 2", TimeCreated: time.Now()}}
				TestDB.db.Create(&activities)
				return featureUUID

			},

			expectError: false,
			expectCount: 2,
		},
		{
			name: "Valid feature UUID with no activities",
			setup: func() string {
				return uuid.New().String()
			},
			expectError: false,
			expectCount: 0,
		},
		{
			name: "Empty feature UUID",
			setup: func() string {
				return ""
			},
			expectError: true,
			errorMsg:    "feature UUID is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			featureUUID := tt.setup()
			activities, err := TestDB.GetActivitiesByFeature(featureUUID)

			if tt.expectError {
				assert.Error(t, err)
				assert.Equal(t, tt.errorMsg, err.Error())
			} else {
				assert.NoError(t, err)
				assert.Len(t, activities, tt.expectCount)
			}
		})
	}
}

func TestGetActivitiesByPhase(t *testing.T) {
	teardownSuite := SetupSuite(t)
	defer teardownSuite(t)

	tests := []struct {
		name        string
		setup       func() string
		expectError bool
		errorMsg    string
		expectCount int
	}{
		{
			name: "Valid phase UUID with activities",
			setup: func() string {
				phaseUUID := uuid.New().String()
				activities := []Activity{
					{
						ID:          uuid.New(),
						PhaseUUID:   phaseUUID,
						Content:     "Activity 1",
						TimeCreated: time.Now().Add(-time.Hour),
						ContentType: FeatureCreation,
						Author:      HumansAuthor,
						AuthorRef:   "12345678901234567890123456789012",
						Workspace:   "workspace-1",
					},
					{
						ID:          uuid.New(),
						PhaseUUID:   phaseUUID,
						Content:     "Activity 2",
						TimeCreated: time.Now(),
						ContentType: FeatureCreation,
						Author:      HumansAuthor,
						AuthorRef:   "12345678901234567890123456789012",
						Workspace:   "workspace-1",
					},
				}
				TestDB.db.Create(&activities)
				return phaseUUID
			},
			expectError: false,
			expectCount: 2,
		},
		{
			name: "Valid phase UUID with no activities",
			setup: func() string {
				return uuid.New().String()
			},
			expectError: false,
			expectCount: 0,
		},
		{
			name: "Empty phase UUID",
			setup: func() string {
				return ""
			},
			expectError: true,
			errorMsg:    "phase UUID is required",
		},
		{
			name: "Whitespace phase UUID",
			setup: func() string {
				return "   "
			},
			expectError: true,
			errorMsg:    "phase UUID is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			phaseUUID := tt.setup()
			activities, err := TestDB.GetActivitiesByPhase(phaseUUID)

			if tt.expectError {
				assert.Error(t, err)
				assert.Equal(t, tt.errorMsg, err.Error())
			} else {
				assert.NoError(t, err)
				assert.Len(t, activities, tt.expectCount)

				if tt.expectCount > 0 {
					for i := 0; i < len(activities)-1; i++ {
						assert.True(t, activities[i].TimeCreated.After(activities[i+1].TimeCreated) ||
							activities[i].TimeCreated.Equal(activities[i+1].TimeCreated))
					}
				}
			}
		})
	}
}

func TestGetActivitiesByWorkspace(t *testing.T) {
	teardownSuite := SetupSuite(t)
	defer teardownSuite(t)

	tests := []struct {
		name        string
		setup       func() string
		expectError bool
		errorMsg    string
		expectCount int
	}{
		{
			name: "Valid workspace with activities",
			setup: func() string {
				workspace := "test-workspace-1"
				activities := []Activity{
					{
						ID:          uuid.New(),
						Workspace:   workspace,
						Content:     "Activity 1",
						TimeCreated: time.Now().Add(-time.Hour),
						ContentType: FeatureCreation,
						Author:      HumansAuthor,
						AuthorRef:   "12345678901234567890123456789012",
					},
					{
						ID:          uuid.New(),
						Workspace:   workspace,
						Content:     "Activity 2",
						TimeCreated: time.Now(),
						ContentType: FeatureCreation,
						Author:      HumansAuthor,
						AuthorRef:   "12345678901234567890123456789012",
					},
				}
				TestDB.db.Create(&activities)
				return workspace
			},
			expectError: false,
			expectCount: 2,
		},
		{
			name: "Valid workspace with no activities",
			setup: func() string {
				return "empty-workspace"
			},
			expectError: false,
			expectCount: 0,
		},
		{
			name: "Empty workspace",
			setup: func() string {
				return ""
			},
			expectError: true,
			errorMsg:    "workspace is required",
		},
		{
			name: "Whitespace workspace",
			setup: func() string {
				return "   "
			},
			expectError: true,
			errorMsg:    "workspace is required",
		},
		{
			name: "Multiple workspaces",
			setup: func() string {
				workspace1 := "test-workspace-2"
				workspace2 := "test-workspace-3"
				activities := []Activity{
					{
						ID:          uuid.New(),
						Workspace:   workspace1,
						Content:     "Activity 1",
						TimeCreated: time.Now(),
						ContentType: FeatureCreation,
						Author:      HumansAuthor,
						AuthorRef:   "12345678901234567890123456789012",
					},
					{
						ID:          uuid.New(),
						Workspace:   workspace2,
						Content:     "Activity 2",
						TimeCreated: time.Now(),
						ContentType: FeatureCreation,
						Author:      HumansAuthor,
						AuthorRef:   "12345678901234567890123456789012",
					},
				}
				TestDB.db.Create(&activities)
				return workspace1
			},
			expectError: false,
			expectCount: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			DeleteAllActivities()
			workspace := tt.setup()
			activities, err := TestDB.GetActivitiesByWorkspace(workspace)

			if tt.expectError {
				assert.Error(t, err)
				assert.Equal(t, tt.errorMsg, err.Error())
			} else {
				assert.NoError(t, err)
				assert.Len(t, activities, tt.expectCount)

				if tt.expectCount > 0 {
					for i := 0; i < len(activities)-1; i++ {
						assert.True(t, activities[i].TimeCreated.After(activities[i+1].TimeCreated) ||
							activities[i].TimeCreated.Equal(activities[i+1].TimeCreated))
					}

					for _, activity := range activities {
						assert.Equal(t, workspace, activity.Workspace)
					}
				}
			}
		})
	}
}

func TestGetLatestActivityByThread(t *testing.T) {
	teardownSuite := SetupSuite(t)
	defer teardownSuite(t)

	tests := []struct {
		name        string
		setup       func() string
		expectError bool
		errorMsg    string
		validate    func(*testing.T, *Activity)
	}{
		{
			name: "Valid thread ID with multiple activities",
			setup: func() string {
				threadID := uuid.New()
				activities := []Activity{
					{
						ID:          uuid.New(),
						ThreadID:    threadID,
						Sequence:    1,
						Content:     "First Activity",
						TimeCreated: time.Now().Add(-2 * time.Hour),
						ContentType: FeatureCreation,
						Author:      HumansAuthor,
						AuthorRef:   "12345678901234567890123456789012",
						Workspace:   "test-workspace",
					},
					{
						ID:          uuid.New(),
						ThreadID:    threadID,
						Sequence:    2,
						Content:     "Second Activity",
						TimeCreated: time.Now().Add(-time.Hour),
						ContentType: FeatureCreation,
						Author:      HumansAuthor,
						AuthorRef:   "12345678901234567890123456789012",
						Workspace:   "test-workspace",
					},
					{
						ID:          uuid.New(),
						ThreadID:    threadID,
						Sequence:    3,
						Content:     "Latest Activity",
						TimeCreated: time.Now(),
						ContentType: FeatureCreation,
						Author:      HumansAuthor,
						AuthorRef:   "12345678901234567890123456789012",
						Workspace:   "test-workspace",
					},
				}
				result := TestDB.db.Create(&activities)
				assert.NoError(t, result.Error, "Failed to create test activities")
				return threadID.String()
			},
			expectError: false,
			validate: func(t *testing.T, activity *Activity) {
				assert.NotNil(t, activity)
				assert.Equal(t, 3, activity.Sequence, "Should return activity with highest sequence")
				assert.Equal(t, "Latest Activity", activity.Content)
			},
		},
		{
			name: "Valid thread ID with single activity",
			setup: func() string {
				threadID := uuid.New()
				activity := Activity{
					ID:          uuid.New(),
					ThreadID:    threadID,
					Sequence:    1,
					Content:     "Single Activity",
					TimeCreated: time.Now(),
					ContentType: FeatureCreation,
					Author:      HumansAuthor,
					AuthorRef:   "12345678901234567890123456789012",
					Workspace:   "test-workspace",
				}
				result := TestDB.db.Create(&activity)
				assert.NoError(t, result.Error, "Failed to create test activity")
				return threadID.String()
			},
			expectError: false,
			validate: func(t *testing.T, activity *Activity) {
				assert.NotNil(t, activity)
				assert.Equal(t, 1, activity.Sequence)
				assert.Equal(t, "Single Activity", activity.Content)
			},
		},
		{
			name: "Non-existent thread ID",
			setup: func() string {
				return uuid.New().String()
			},
			expectError: true,
			errorMsg:    "record not found",
			validate:    func(t *testing.T, activity *Activity) {},
		},
		{
			name: "Invalid thread ID format",
			setup: func() string {
				return "invalid-uuid"
			},
			expectError: true,
			errorMsg:    "invalid thread ID format",
			validate:    func(t *testing.T, activity *Activity) {},
		},
		{
			name: "Empty thread ID",
			setup: func() string {
				return ""
			},
			expectError: true,
			errorMsg:    "invalid thread ID format",
			validate:    func(t *testing.T, activity *Activity) {},
		},
		{
			name: "Thread ID with special characters",
			setup: func() string {
				return uuid.New().String() + "@!#"
			},
			expectError: true,
			errorMsg:    "invalid thread ID format",
			validate:    func(t *testing.T, activity *Activity) {},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			threadID := tt.setup()
			activity, err := TestDB.GetLatestActivityByThread(threadID)

			if tt.expectError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorMsg)
				assert.Nil(t, activity)
			} else {
				assert.NoError(t, err)
				tt.validate(t, activity)
			}
		})
	}
}

func TestCreateActivityThread(t *testing.T) {
	teardownSuite := SetupSuite(t)
	defer teardownSuite(t)

	existingThreadID := uuid.New().String()

	tests := []struct {
		name        string
		sourceID    string
		activity    *Activity
		setup       func()
		expectError bool
		errorMsg    string
		validate    func(*testing.T, *Activity)
	}{
		{
			name:     "Successfully create first activity in thread",
			sourceID: uuid.New().String(),
			activity: &Activity{
				Content:     "First thread activity",
				ContentType: FeatureCreation,
				Author:      HumansAuthor,
				AuthorRef:   "12345678901234567890123456789012",
				Workspace:   "test-workspace",
			},
			expectError: false,
			validate: func(t *testing.T, activity *Activity) {
				assert.NotNil(t, activity)
				assert.Equal(t, 1, activity.Sequence, "First activity should have sequence 1")
				assert.NotEqual(t, uuid.Nil, activity.ThreadID)
				assert.Equal(t, "First thread activity", activity.Content)
			},
		},
		{
			name:     "Successfully add to existing thread",
			sourceID: existingThreadID,
			activity: &Activity{
				Content:     "Second thread activity",
				ContentType: FeatureCreation,
				Author:      HumansAuthor,
				AuthorRef:   "12345678901234567890123456789012",
				Workspace:   "test-workspace",
			},
			setup: func() {
				existingActivity := Activity{
					ID:          uuid.New(),
					ThreadID:    uuid.MustParse(existingThreadID),
					Sequence:    1,
					Content:     "First Activity",
					ContentType: FeatureCreation,
					Author:      HumansAuthor,
					AuthorRef:   "12345678901234567890123456789012",
					Workspace:   "test-workspace",
				}
				result := TestDB.db.Create(&existingActivity)
				assert.NoError(t, result.Error, "Failed to create existing activity")
			},
			expectError: false,
			validate: func(t *testing.T, activity *Activity) {
				assert.NotNil(t, activity)
				assert.Equal(t, 2, activity.Sequence, "Second activity should have sequence 2")
				assert.NotEqual(t, uuid.Nil, activity.ThreadID)
				assert.Equal(t, "Second thread activity", activity.Content)
			},
		},
		{
			name:     "Invalid source ID format",
			sourceID: "invalid-uuid",
			activity: &Activity{
				Content:     "Test content",
				ContentType: FeatureCreation,
				Author:      HumansAuthor,
				AuthorRef:   "12345678901234567890123456789012",
				Workspace:   "test-workspace",
			},
			expectError: true,
			errorMsg:    "invalid source ID format",
		},
		{
			name:     "Empty source ID",
			sourceID: "",
			activity: &Activity{
				Content:     "Test content",
				ContentType: FeatureCreation,
				Author:      HumansAuthor,
				AuthorRef:   "12345678901234567890123456789012",
				Workspace:   "test-workspace",
			},
			expectError: true,
			errorMsg:    "invalid source ID format",
		},
		{
			name:     "Invalid activity content",
			sourceID: uuid.New().String(),
			activity: &Activity{
				Content:     "",
				ContentType: FeatureCreation,
				Author:      HumansAuthor,
				AuthorRef:   "12345678901234567890123456789012",
				Workspace:   "test-workspace",
			},
			expectError: true,
			errorMsg:    ErrInvalidContent.Error(),
		},
		{
			name:     "Invalid content type",
			sourceID: uuid.New().String(),
			activity: &Activity{
				Content:     "Valid content",
				ContentType: "invalid_type",
				Author:      HumansAuthor,
				AuthorRef:   "12345678901234567890123456789012",
				Workspace:   "test-workspace",
			},
			expectError: true,
			errorMsg:    ErrInvalidContentType.Error(),
		},
		{
			name:     "Invalid author type",
			sourceID: uuid.New().String(),
			activity: &Activity{
				Content:     "Valid content",
				ContentType: FeatureCreation,
				Author:      "invalid_author",
				AuthorRef:   "12345678901234567890123456789012",
				Workspace:   "test-workspace",
			},
			expectError: true,
			errorMsg:    ErrInvalidAuthorType.Error(),
		},
		{
			name:     "Invalid author reference",
			sourceID: uuid.New().String(),
			activity: &Activity{
				Content:     "Valid content",
				ContentType: FeatureCreation,
				Author:      HumansAuthor,
				AuthorRef:   "short",
				Workspace:   "test-workspace",
			},
			expectError: true,
			errorMsg:    "invalid public key format for human author",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if tt.setup != nil {
				tt.setup()
			}

			activity, err := TestDB.CreateActivityThread(tt.sourceID, tt.activity)

			if tt.expectError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorMsg)
				assert.Nil(t, activity)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, activity)
				assert.NotEqual(t, uuid.Nil, activity.ID)
				assert.NotZero(t, activity.TimeCreated)
				assert.NotZero(t, activity.TimeUpdated)

				if tt.validate != nil {
					tt.validate(t, activity)
				}

				savedActivity, err := TestDB.GetActivity(activity.ID.String())
				assert.NoError(t, err)
				assert.Equal(t, activity.ID, savedActivity.ID)
				assert.Equal(t, activity.Sequence, savedActivity.Sequence)
				assert.Equal(t, activity.ThreadID, savedActivity.ThreadID)
			}
		})
	}
}

func TestValidateActivityTitle(t *testing.T) {
    teardownSuite := SetupSuite(t)
    defer teardownSuite(t)

    tests := []struct {
        name     string
        input    *Activity
        expected error
    }{
        {
            name: "Valid title",
            input: &Activity{
                Title:       "Test Title",
                Content:     "Valid content",
                AuthorRef:   "12345678901234567890123456789012",
                ContentType: FeatureCreation,
                Author:      HumansAuthor,
                Workspace:   "valid_workspace",
            },
            expected: nil,
        },
        {
            name: "Empty title",
            input: &Activity{
                Title:       "",
                Content:     "Valid content",
                AuthorRef:   "12345678901234567890123456789012",
                ContentType: FeatureCreation,
                Author:      HumansAuthor,
                Workspace:   "valid_workspace",
            },
            expected: nil,
        },
        {
            name: "Title at max length",
            input: &Activity{
                Title:       strings.Repeat("a", 200),
                Content:     "Valid content",
                AuthorRef:   "12345678901234567890123456789012",
                ContentType: FeatureCreation,
                Author:      HumansAuthor,
                Workspace:   "valid_workspace",
            },
            expected: nil,
        },
        {
            name: "Title exceeds max length",
            input: &Activity{
                Title:       strings.Repeat("a", 201),
                Content:     "Valid content",
                AuthorRef:   "12345678901234567890123456789012",
                ContentType: FeatureCreation,
                Author:      HumansAuthor,
                Workspace:   "valid_workspace",
            },
            expected: ErrInvalidTitle,
        },
    }

    for _, test := range tests {
        t.Run(test.name, func(t *testing.T) {
            err := validateActivity(test.input)
            if test.expected == nil {
                assert.NoError(t, err)
            } else {
                assert.EqualError(t, err, test.expected.Error())
            }
        })
    }
}
