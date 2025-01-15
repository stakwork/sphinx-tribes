package feeds

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestParseMediumFeed(t *testing.T) {
	zeroTime := time.Time{}.Unix()

	tests := []struct {
		name     string
		url      string
		bod      []byte
		expected *Feed
		wantErr  bool
	}{
		{
			name: "Valid Complete Medium Feed",
			url:  "https://medium.com/feed",
			bod: []byte(`<?xml version="1.0" encoding="UTF-8"?>
				<rss>
					<channel>
						<title>Test Blog</title>
						<link>https://medium.com/blog</link>
						<description>Test Description</description>
						<image>
							<url>https://test.com/image.jpg</url>
						</image>
						<generator>Medium</generator>
						<lastBuildDate>Mon, 02 Jan 2006 15:04:05 GMT</lastBuildDate>
						<creator>Test Author</creator>
						<item>
							<title>Test Post</title>
							<description>Test Post Description</description>
							<link>https://medium.com/post1</link>
							<guid>post-1</guid>
							<pubDate>Mon, 02 Jan 2006 15:04:05 GMT</pubDate>
							<updated>Mon, 02 Jan 2006 15:04:05 GMT</updated>
							<creator>Post Author</creator>
						</item>
					</channel>
				</rss>`),
			expected: &Feed{
				ID:          "https://medium.com/feed",
				FeedType:    FeedTypeBlog,
				Title:       "Test Blog",
				Url:         "https://medium.com/feed",
				Link:        "https://medium.com/blog",
				Description: "Test Description",
				ImageUrl:    "https://test.com/image.jpg",
				Generator:   "Medium",
				Author:      "Test Author",
				DateUpdated: 1136214245,
				Items: []Item{
					{
						Id:            "post-1",
						Title:         "Test Post",
						Link:          "https://medium.com/post1",
						EnclosureURL:  "https://medium.com/post1",
						Description:   "Test Post Description",
						Author:        "Post Author",
						DatePublished: 1136214245,
						DateUpdated:   1136214245,
					},
				},
			},
			wantErr: false,
		},
		{
			name: "Empty the Feed",
			url:  "https://medium.com/feed",
			bod: []byte(`<?xml version="1.0" encoding="UTF-8"?>
				<rss>
					<channel>
					</channel>
				</rss>`),
			expected: &Feed{
				ID:          "https://medium.com/feed",
				FeedType:    FeedTypeBlog,
				Url:         "https://medium.com/feed",
				DateUpdated: zeroTime,
				Items:       []Item{},
			},
			wantErr: false,
		},
		{
			name:     "Invalid XML",
			url:      "https://medium.com/feed",
			bod:      []byte(`invalid xml`),
			expected: nil,
			wantErr:  true,
		},
		{
			name: "Feed with Special Characters",
			url:  "https://medium.com/feed",
			bod: []byte(`<?xml version="1.0" encoding="UTF-8"?>
				<rss>
					<channel>
						<title>Test &amp; Blog</title>
						<description>Test &lt;Description&gt;</description>
						<item>
							<title>Test &quot;Post&quot;</title>
						</item>
					</channel>
				</rss>`),
			expected: &Feed{
				ID:          "https://medium.com/feed",
				FeedType:    FeedTypeBlog,
				Title:       "Test & Blog",
				Description: "Test <Description>",
				Url:         "https://medium.com/feed",
				DateUpdated: zeroTime,
				Items: []Item{
					{
						Title:         "Test \"Post\"",
						DatePublished: zeroTime,
						DateUpdated:   zeroTime,
					},
				},
			},
			wantErr: false,
		},
		{
			name: "Feed with Invalid Dates",
			url:  "https://medium.com/feed",
			bod: []byte(`<?xml version="1.0" encoding="UTF-8"?>
				<rss>
					<channel>
						<lastBuildDate>invalid date</lastBuildDate>
						<item>
							<pubDate>invalid date</pubDate>
							<updated>invalid date</updated>
						</item>
					</channel>
				</rss>`),
			expected: &Feed{
				ID:          "https://medium.com/feed",
				FeedType:    FeedTypeBlog,
				Url:         "https://medium.com/feed",
				DateUpdated: zeroTime,
				Items: []Item{
					{
						DatePublished: zeroTime,
						DateUpdated:   zeroTime,
					},
				},
			},
			wantErr: false,
		},
		{
			name:     "Empty Body",
			url:      "https://medium.com/feed",
			bod:      []byte{},
			expected: nil,
			wantErr:  true,
		},
		{
			name: "Feed with Multiple Items",
			url:  "https://medium.com/feed",
			bod: []byte(`<?xml version="1.0" encoding="UTF-8"?>
				<rss>
					<channel>
						<item>
							<title>Post 1</title>
							<guid>1</guid>
						</item>
						<item>
							<title>Post 2</title>
							<guid>2</guid>
						</item>
					</channel>
				</rss>`),
			expected: &Feed{
				ID:          "https://medium.com/feed",
				FeedType:    FeedTypeBlog,
				Url:         "https://medium.com/feed",
				DateUpdated: zeroTime,
				Items: []Item{
					{
						Id:            "1",
						Title:         "Post 1",
						DatePublished: zeroTime,
						DateUpdated:   zeroTime,
					},
					{
						Id:            "2",
						Title:         "Post 2",
						DatePublished: zeroTime,
						DateUpdated:   zeroTime,
					},
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseMediumFeed(tt.url, tt.bod)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, got)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tt.expected, got)
		})
	}
}
