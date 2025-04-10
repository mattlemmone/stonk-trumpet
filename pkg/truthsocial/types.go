package truthsocial

import "time"

// Status represents a single status (tweet) from the Truth Social API.
type Status struct {
	ID                 string            `json:"id"`
	CreatedAt          time.Time         `json:"created_at"`
	InReplyToID        *string           `json:"in_reply_to_id"` // Use pointer for nullable fields
	QuoteID            *string           `json:"quote_id"`
	InReplyToAccountID *string           `json:"in_reply_to_account_id"`
	Sensitive          bool              `json:"sensitive"`
	SpoilerText        string            `json:"spoiler_text"`
	Visibility         string            `json:"visibility"`
	Language           *string           `json:"language"`
	URI                string            `json:"uri"`
	URL                string            `json:"url"`
	Content            string            `json:"content"` // Contains HTML
	Account            Account           `json:"account"`
	MediaAttachments   []MediaAttachment `json:"media_attachments"`
	Mentions           []Mention         `json:"mentions"`    // Assuming Mention struct needed if not empty
	Tags               []Tag             `json:"tags"`        // Assuming Tag struct needed if not empty
	Card               *Card             `json:"card"`        // Assuming Card struct needed if not null
	Group              *Group            `json:"group"`       // Assuming Group struct needed if not null
	Quote              *Status           `json:"quote"`       // Can be a nested Status
	InReplyTo          *Status           `json:"in_reply_to"` // Can be a nested Status
	Reblog             *Status           `json:"reblog"`      // Can be a nested Status
	Sponsored          bool              `json:"sponsored"`
	RepliesCount       int               `json:"replies_count"`
	ReblogsCount       int               `json:"reblogs_count"`
	FavouritesCount    int               `json:"favourites_count"`
	Favourited         bool              `json:"favourited"`
	Reblogged          bool              `json:"reblogged"`
	Muted              bool              `json:"muted"`
	Pinned             bool              `json:"pinned"`
	Bookmarked         bool              `json:"bookmarked"`
	Poll               *Poll             `json:"poll"`   // Assuming Poll struct needed if not null
	Emojis             []Emoji           `json:"emojis"` // Assuming Emoji struct needed if not empty
}

// Account represents the user account associated with a status.
type Account struct {
	ID                         string    `json:"id"`
	Username                   string    `json:"username"`
	Acct                       string    `json:"acct"`
	DisplayName                string    `json:"display_name"`
	Locked                     bool      `json:"locked"`
	Bot                        bool      `json:"bot"`
	Discoverable               bool      `json:"discoverable"`
	Group                      bool      `json:"group"`
	CreatedAt                  time.Time `json:"created_at"`
	Note                       string    `json:"note"` // Contains HTML
	URL                        string    `json:"url"`
	Avatar                     string    `json:"avatar"`
	AvatarStatic               string    `json:"avatar_static"`
	Header                     string    `json:"header"`
	HeaderStatic               string    `json:"header_static"`
	FollowersCount             int       `json:"followers_count"`
	FollowingCount             int       `json:"following_count"`
	StatusesCount              int       `json:"statuses_count"`
	LastStatusAt               string    `json:"last_status_at"` // Consider parsing as time.Time if needed
	Verified                   bool      `json:"verified"`
	Location                   string    `json:"location"`
	Website                    string    `json:"website"`
	UnauthVisibility           bool      `json:"unauth_visibility"`
	ChatsOnboarded             bool      `json:"chats_onboarded"`
	FeedsOnboarded             bool      `json:"feeds_onboarded"`
	AcceptingMessages          bool      `json:"accepting_messages"`
	ShowNonmemberGroupStatuses *bool     `json:"show_nonmember_group_statuses"`
	Emojis                     []Emoji   `json:"emojis"` // Assuming Emoji struct needed if not empty
	Fields                     []Field   `json:"fields"` // Assuming Field struct needed if not empty
	TVOnboarded                bool      `json:"tv_onboarded"`
	TVAccount                  bool      `json:"tv_account"`
}

// MediaAttachment represents media attached to a status.
type MediaAttachment struct {
	ID               string    `json:"id"`
	Type             string    `json:"type"` // e.g., "video", "image"
	URL              string    `json:"url"`
	PreviewURL       string    `json:"preview_url"`
	ExternalVideoID  *string   `json:"external_video_id"`
	RemoteURL        *string   `json:"remote_url"`
	PreviewRemoteURL *string   `json:"preview_remote_url"`
	TextURL          *string   `json:"text_url"`
	Meta             MediaMeta `json:"meta"`
	Description      *string   `json:"description"`
	Blurhash         *string   `json:"blurhash"`
	Processing       *string   `json:"processing"` // e.g., "complete"
}

// MediaMeta contains metadata about the media attachment.
type MediaMeta struct {
	Colors   *ColorsMeta `json:"colors,omitempty"`
	Original *MediaSize  `json:"original,omitempty"`
	Small    *MediaSize  `json:"small,omitempty"`
	// Add other potential sizes if needed
}

// ColorsMeta contains color information for media previews.
type ColorsMeta struct {
	Background string `json:"background"`
	Foreground string `json:"foreground"`
	Accent     string `json:"accent"`
}

// MediaSize contains dimensions and other info for a specific media size.
type MediaSize struct {
	Width     int     `json:"width"`
	Height    int     `json:"height"`
	FrameRate *string `json:"frame_rate,omitempty"` // e.g., "60/1"
	Duration  float64 `json:"duration,omitempty"`
	Bitrate   int     `json:"bitrate,omitempty"`
	Size      *string `json:"size,omitempty"` // e.g., "270x480"
	Aspect    float64 `json:"aspect,omitempty"`
}

// Mention represents a user mention within a status.
// Placeholder struct - Define fields if needed based on actual API response.
type Mention struct {
	// e.g., ID, Username, URL
}

// Tag represents a hashtag within a status.
// Placeholder struct - Define fields if needed based on actual API response.
type Tag struct {
	// e.g., Name, URL
}

// Card represents a preview card embedded in a status.
// Placeholder struct - Define fields if needed based on actual API response.
type Card struct {
	// Define fields based on card structure
}

// Group represents a group associated with a status.
// Placeholder struct - Define fields if needed based on actual API response.
type Group struct {
	// Define fields based on group structure
}

// Poll represents a poll attached to a status.
// Placeholder struct - Define fields if needed based on actual API response.
type Poll struct {
	// Define fields based on poll structure
}

// Emoji represents a custom emoji used.
// Placeholder struct - Define fields if needed based on actual API response.
type Emoji struct {
	// e.g., Shortcode, URL, StaticURL
}

// Field represents custom profile metadata fields.
// Placeholder struct - Define fields if needed based on actual API response.
type Field struct {
	// e.g., Name, Value
}
