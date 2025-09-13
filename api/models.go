package api

// Member represents a band member
type Member struct {
	ID          int      `json:"id"`
	Name        string   `json:"name"`
	StageName   string   `json:"stage_name"`
	Position    string   `json:"position"`
	Instrument  string   `json:"instrument"`
	Birthdate   string   `json:"birthdate"`
	Nationality string   `json:"nationality"`
	Height      string   `json:"height"`
	BloodType   string   `json:"blood_type"`
	MBTI        string   `json:"mbti"`
	Hobbies     []string `json:"hobbies"`
	Facts       []string `json:"facts"`
}

// Song represents a band's song
type Song struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Album       string `json:"album"`
	ReleaseDate string `json:"release_date"`
	Duration    string `json:"duration"`
	Genre       string `json:"genre"`
	Language    string `json:"language"`
	YouTubeURL  string `json:"youtube_url"`
	SpotifyURL  string `json:"spotify_url"`
}

// Album represents a band's album
type Album struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	ReleaseDate string `json:"release_date"`
	Type        string `json:"type"` // Single, EP, Album
	TrackCount  int    `json:"track_count"`
	CoverImage  string `json:"cover_image"`
	Songs       []Song `json:"songs"`
}

// Award represents an award or prize
type Award struct {
	Year      int    `json:"year"`
	Event     string `json:"event"`
	Category  string `json:"category"`
	Recipient string `json:"recipient"`
}

// Band represents the band information
type Band struct {
	ID           int               `json:"id"`
	Name         string            `json:"name"`
	KoreanName   string            `json:"korean_name"`
	DebutDate    string            `json:"debut_date"`
	Company      string            `json:"company"`
	Genre        []string          `json:"genre"`
	MemberCount  int               `json:"member_count"`
	Members      []Member          `json:"members"`
	Description  string            `json:"description"`
	SocialMedia  map[string]string `json:"social_media"`
	OfficialSite string            `json:"official_site"`
	Discography  []Album           `json:"discography"`
	Awards       []Award           `json:"awards"`
}

// APIResponse represents standard API response format
type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}
