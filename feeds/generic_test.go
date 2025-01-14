package feeds

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAddedValue(t *testing.T) {
	tests := []struct {
		name             string
		value            *Value
		tribeOwnerPubkey string
		expected         *Value
	}{
		{
			name: "Empty tribe owner pubkey",
			value: &Value{
				Model: Model{
					Type:      "lightning",
					Suggested: "0.0001",
				},
				Destinations: []Destination{
					{
						Address: "existing_address",
						Split:   json.Number("1"),
						Type:    "node",
					},
				},
			},
			tribeOwnerPubkey: "",
			expected: &Value{
				Model: Model{
					Type:      "lightning",
					Suggested: "0.0001",
				},
				Destinations: []Destination{
					{
						Address: "existing_address",
						Split:   json.Number("1"),
						Type:    "node",
					},
				},
			},
		},
		{
			name:             "Nil value with tribe owner",
			value:            nil,
			tribeOwnerPubkey: "tribe_owner_key",
			expected: &Value{
				Model: Model{
					Type:      "lightning",
					Suggested: "0.00000015000",
				},
				Destinations: []Destination{
					{
						Address: "tribe_owner_key",
						Type:    "node",
						Split:   100,
					},
				},
			},
		},
		{
			name: "Single destination with split 1 and different address",
			value: &Value{
				Model: Model{
					Type:      "lightning",
					Suggested: "0.0001",
				},
				Destinations: []Destination{
					{
						Address: "different_address",
						Split:   json.Number("1"),
						Type:    "node",
					},
				},
			},
			tribeOwnerPubkey: "tribe_owner_key",
			expected: &Value{
				Model: Model{
					Type:      "lightning",
					Suggested: "0.0001",
				},
				Destinations: []Destination{
					{
						Address: "different_address",
						Split:   json.Number("1"),
						Type:    "node",
					},
					{
						Address: "tribe_owner_key",
						Split:   99,
						Type:    "node",
					},
				},
			},
		},
		{
			name: "Single destination with split 1 and same address",
			value: &Value{
				Model: Model{
					Type:      "lightning",
					Suggested: "0.0001",
				},
				Destinations: []Destination{
					{
						Address: "tribe_owner_key",
						Split:   json.Number("1"),
						Type:    "node",
					},
				},
			},
			tribeOwnerPubkey: "tribe_owner_key",
			expected: &Value{
				Model: Model{
					Type:      "lightning",
					Suggested: "0.0001",
				},
				Destinations: []Destination{
					{
						Address: "tribe_owner_key",
						Split:   json.Number("1"),
						Type:    "node",
					},
				},
			},
		},
		{
			name: "Multiple destinations",
			value: &Value{
				Model: Model{
					Type:      "lightning",
					Suggested: "0.0001",
				},
				Destinations: []Destination{
					{
						Address: "address1",
						Split:   json.Number("50"),
						Type:    "node",
					},
					{
						Address: "address2",
						Split:   json.Number("50"),
						Type:    "node",
					},
				},
			},
			tribeOwnerPubkey: "tribe_owner_key",
			expected: &Value{
				Model: Model{
					Type:      "lightning",
					Suggested: "0.0001",
				},
				Destinations: []Destination{
					{
						Address: "address1",
						Split:   json.Number("50"),
						Type:    "node",
					},
					{
						Address: "address2",
						Split:   json.Number("50"),
						Type:    "node",
					},
				},
			},
		},
		{
			name: "Empty destinations",
			value: &Value{
				Model: Model{
					Type:      "lightning",
					Suggested: "0.0001",
				},
				Destinations: []Destination{},
			},
			tribeOwnerPubkey: "tribe_owner_key",
			expected: &Value{
				Model: Model{
					Type:      "lightning",
					Suggested: "0.0001",
				},
				Destinations: []Destination{},
			},
		},
		{
			name: "Nil destinations",
			value: &Value{
				Model: Model{
					Type:      "lightning",
					Suggested: "0.0001",
				},
				Destinations: nil,
			},
			tribeOwnerPubkey: "tribe_owner_key",
			expected: &Value{
				Model: Model{
					Type:      "lightning",
					Suggested: "0.0001",
				},
				Destinations: nil,
			},
		},
		{
			name: "Split as Zero",
			value: &Value{
				Model: Model{
					Type:      "lightning",
					Suggested: "0.0001",
				},
				Destinations: []Destination{
					{
						Address: "address1",
						Split:   json.Number("0"),
						Type:    "node",
					},
				},
			},
			tribeOwnerPubkey: "tribe_owner_key",
			expected: &Value{
				Model: Model{
					Type:      "lightning",
					Suggested: "0.0001",
				},
				Destinations: []Destination{
					{
						Address: "address1",
						Split:   json.Number("0"),
						Type:    "node",
					},
				},
			},
		},
		{
			name: "Split as Negative Number",
			value: &Value{
				Model: Model{
					Type:      "lightning",
					Suggested: "0.0001",
				},
				Destinations: []Destination{
					{
						Address: "address1",
						Split:   json.Number("-1"),
						Type:    "node",
					},
				},
			},
			tribeOwnerPubkey: "tribe_owner_key",
			expected: &Value{
				Model: Model{
					Type:      "lightning",
					Suggested: "0.0001",
				},
				Destinations: []Destination{
					{
						Address: "address1",
						Split:   json.Number("-1"),
						Type:    "node",
					},
				},
			},
		},
		{
			name: "tribeOwnerPubkey as Special Characters",
			value: &Value{
				Model: Model{
					Type:      "lightning",
					Suggested: "0.0001",
				},
				Destinations: []Destination{
					{
						Address: "different_address",
						Split:   json.Number("1"),
						Type:    "node",
					},
				},
			},
			tribeOwnerPubkey: "!@#$%^&*()_+",
			expected: &Value{
				Model: Model{
					Type:      "lightning",
					Suggested: "0.0001",
				},
				Destinations: []Destination{
					{
						Address: "different_address",
						Split:   json.Number("1"),
						Type:    "node",
					},
					{
						Address: "!@#$%^&*()_+",
						Split:   99,
						Type:    "node",
					},
				},
			},
		},
		{
			name: "Split as a Large Number",
			value: &Value{
				Model: Model{
					Type:      "lightning",
					Suggested: "0.0001",
				},
				Destinations: []Destination{
					{
						Address: "address1",
						Split:   json.Number("999999999999"),
						Type:    "node",
					},
				},
			},
			tribeOwnerPubkey: "tribe_owner_key",
			expected: &Value{
				Model: Model{
					Type:      "lightning",
					Suggested: "0.0001",
				},
				Destinations: []Destination{
					{
						Address: "address1",
						Split:   json.Number("999999999999"),
						Type:    "node",
					},
				},
			},
		},
		{
			name: "tribeOwnerPubkey as a Very Long String",
			value: &Value{
				Model: Model{
					Type:      "lightning",
					Suggested: "0.0001",
				},
				Destinations: []Destination{
					{
						Address: "different_address",
						Split:   json.Number("1"),
						Type:    "node",
					},
				},
			},
			tribeOwnerPubkey: strings.Repeat("a", 1000),
			expected: &Value{
				Model: Model{
					Type:      "lightning",
					Suggested: "0.0001",
				},
				Destinations: []Destination{
					{
						Address: "different_address",
						Split:   json.Number("1"),
						Type:    "node",
					},
					{
						Address: strings.Repeat("a", 1000),
						Split:   99,
						Type:    "node",
					},
				},
			},
		},
		{
			name: "Large Number of Destinations",
			value: &Value{
				Model: Model{
					Type:      "lightning",
					Suggested: "0.0001",
				},
				Destinations: func() []Destination {
					dests := make([]Destination, 100)
					for i := 0; i < 100; i++ {
						dests[i] = Destination{
							Address: fmt.Sprintf("address_%d", i),
							Split:   json.Number(fmt.Sprintf("%d", i+1)),
							Type:    "node",
						}
					}
					return dests
				}(),
			},
			tribeOwnerPubkey: "tribe_owner_key",
			expected: &Value{
				Model: Model{
					Type:      "lightning",
					Suggested: "0.0001",
				},
				Destinations: func() []Destination {
					dests := make([]Destination, 100)
					for i := 0; i < 100; i++ {
						dests[i] = Destination{
							Address: fmt.Sprintf("address_%d", i),
							Split:   json.Number(fmt.Sprintf("%d", i+1)),
							Type:    "node",
						}
					}
					return dests
				}(),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := AddedValue(tt.value, tt.tribeOwnerPubkey)

			if tt.expected == nil {
				assert.Nil(t, result)
				return
			}

			assert.NotNil(t, result)

			assert.Equal(t, tt.expected.Model.Type, result.Model.Type)
			assert.Equal(t, tt.expected.Model.Suggested, result.Model.Suggested)

			assert.Equal(t, len(tt.expected.Destinations), len(result.Destinations))
			for i, expectedDest := range tt.expected.Destinations {
				assert.Equal(t, expectedDest.Address, result.Destinations[i].Address)
				assert.Equal(t, expectedDest.Type, result.Destinations[i].Type)
				assert.Equal(t, expectedDest.Split, result.Destinations[i].Split)
				assert.Equal(t, expectedDest.CustomKey, result.Destinations[i].CustomKey)
				assert.Equal(t, expectedDest.CustomValue, result.Destinations[i].CustomValue)
			}
		})
	}
}
