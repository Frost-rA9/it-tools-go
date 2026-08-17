package metatag

// Option 表示 select 元素的选项，可为分组（Type="group"）。
type Option struct {
	Label    string   `json:"label"`
	Value    string   `json:"value,omitempty"`
	Type     string   `json:"type,omitempty"`
	Key      string   `json:"key,omitempty"`
	Children []Option `json:"children,omitempty"`
}

// Element 表示表单中的一个字段。
type Element struct {
	Key         string   `json:"key"`
	Label       string   `json:"label"`
	Placeholder string   `json:"placeholder"`
	Type        string   `json:"type"` // input | select | input-multiple
	Options     []Option `json:"options,omitempty"`
}

// Schema 表示一个表单区块。
type Schema struct {
	Name     string    `json:"name"`
	Elements []Element `json:"elements"`
}

// schemasAction 的返回结构。
type schemasResult struct {
	Base  []Schema          `json:"base"`
	Types map[string]Schema `json:"types"`
}

// baseSchemas 为始终展示的区块（General / Image / Twitter）。
var baseSchemas = []Schema{
	{
		Name: "General information",
		Elements: []Element{
			{
				Key: "type", Label: "Page type", Type: "select",
				Placeholder: "Select the type of your website...",
				Options: []Option{
					{Label: "Website", Value: "website"},
					{Label: "Article", Value: "article"},
					{Label: "Book", Value: "book"},
					{Label: "Profile", Value: "profile"},
					{
						Type: "group", Label: "Music", Key: "Music",
						Children: []Option{
							{Label: "Song", Value: "music.song"},
							{Label: "Music album", Value: "music.album"},
							{Label: "Playlist", Value: "music.playlist"},
							{Label: "Radio station", Value: "music.radio_station"},
						},
					},
					{
						Type: "group", Label: "Video", Key: "Video",
						Children: []Option{
							{Label: "Movie", Value: "video.movie"},
							{Label: "Episode", Value: "video.episode"},
							{Label: "TV show", Value: "video.tv_show"},
							{Label: "Other video", Value: "video.other"},
						},
					},
				},
			},
			{Key: "title", Label: "Title", Type: "input", Placeholder: "Enter the title of your website..."},
			{Key: "description", Label: "Description", Type: "input", Placeholder: "Enter the description of your website..."},
			{Key: "url", Label: "Page URL", Type: "input", Placeholder: "Enter the url of your website..."},
		},
	},
	{
		Name: "Image",
		Elements: []Element{
			{Key: "image", Label: "Image url", Type: "input", Placeholder: "The url of your website social image..."},
			{Key: "image:alt", Label: "Image alt", Type: "input", Placeholder: "The alternative text of your website social image..."},
			{Key: "image:width", Label: "Width", Type: "input", Placeholder: "Width in px of your website social image..."},
			{Key: "image:height", Label: "Height", Type: "input", Placeholder: "Height in px of your website social image..."},
		},
	},
	{
		Name: "Twitter",
		Elements: []Element{
			{
				Key: "twitter:card", Label: "Card type", Type: "select",
				Placeholder: "The Twitter card type...",
				Options: []Option{
					{Label: "Summary", Value: "summary"},
					{Label: "Summary with large image", Value: "summary_large_image"},
					{Label: "Application", Value: "app"},
					{Label: "Player", Value: "player"},
				},
			},
			{Key: "twitter:site", Label: "Site account", Type: "input", Placeholder: "The name of the Twitter account of the site (ex: @ittoolsdottech)..."},
			{Key: "twitter:creator", Label: "Creator acc.", Type: "input", Placeholder: "The name of the Twitter account of the creator (ex: @cthmsst)..."},
		},
	},
}

// videoMovieSchema 为视频相关区块的公共元素。
var videoMovieSchema = Schema{
	Name: "Movie details",
	Elements: []Element{
		{Key: "video:actor", Label: "Actor", Type: "input-multiple", Placeholder: "Name of the actress/actor..."},
		{Key: "video:director", Label: "Director", Type: "input-multiple", Placeholder: "Name of the director..."},
		{Key: "video:writer", Label: "Writer", Type: "input-multiple", Placeholder: "Writers of the movie..."},
		{Key: "video:duration", Label: "Duration", Type: "input", Placeholder: "The movie's length in seconds..."},
		{Key: "video:release_date", Label: "Release date", Type: "input", Placeholder: "The date the movie was released..."},
		{Key: "video:tag", Label: "Tag", Type: "input", Placeholder: "Tag words associated with this movie..."},
	},
}

// typeSchemas 为按页面类型追加的区块。
var typeSchemas = map[string]Schema{
	"article": {
		Name: "Article",
		Elements: []Element{
			{Key: "article:published_time", Label: "Publishing date", Type: "input", Placeholder: "When the article was first published..."},
			{Key: "article:modified_time", Label: "Modification date", Type: "input", Placeholder: "When the article was last changed..."},
			{Key: "article:expiration_time", Label: "Expiration date", Type: "input", Placeholder: "When the article is out of date after..."},
			{Key: "article:author", Label: "Author", Type: "input", Placeholder: "Writers of the article..."},
			{Key: "article:section", Label: "Section", Type: "input", Placeholder: "A high-level section name. E.g. Technology.."},
			{Key: "article:tag", Label: "Tag", Type: "input", Placeholder: "Tag words associated with this article..."},
		},
	},
	"book": {
		Name: "Book",
		Elements: []Element{
			{Key: "book:author", Label: "Author", Type: "input", Placeholder: "Who wrote this book..."},
			{Key: "book:isbn", Label: "ISBN", Type: "input", Placeholder: "The International Standard Book Number..."},
			{Key: "book:release_date", Label: "Release date", Type: "input", Placeholder: "The date the book was released..."},
			{Key: "book:tag", Label: "Tag", Type: "input", Placeholder: "Tag words associated with this book..."},
		},
	},
	"profile": {
		Name: "Profile",
		Elements: []Element{
			{Key: "profile:first_name", Label: "First name", Type: "input", Placeholder: "Enter the first name of the person..."},
			{Key: "profile:last_name", Label: "Last name", Type: "input", Placeholder: "Enter the last name of the person..."},
			{Key: "profile:username", Label: "Username", Type: "input", Placeholder: "Enter the username of the person..."},
			{Key: "profile:gender", Label: "Gender", Type: "input", Placeholder: "Enter the gender of the person..."},
		},
	},
	"music.song": {
		Name: "Song details",
		Elements: []Element{
			{Key: "music:duration", Label: "Duration", Type: "input", Placeholder: "The duration of the song..."},
			{Key: "music:album", Label: "Album", Type: "input", Placeholder: "The album this song is from..."},
			{Key: "music:album:disk", Label: "Disc", Type: "input", Placeholder: "Which disc of the album this song is on..."},
			{Key: "music:album:track", Label: "Track", Type: "input", Placeholder: "Which track this song is..."},
			{Key: "music:musician", Label: "Musician", Type: "input-multiple", Placeholder: "The musician that made this song..."},
		},
	},
	"music.album": {
		Name: "Album details",
		Elements: []Element{
			{Key: "music:song", Label: "Song", Type: "input", Placeholder: "The song on this album..."},
			{Key: "music:song:disc", Label: "Disc", Type: "input", Placeholder: "The same as music:album:disc but in reverse..."},
			{Key: "music:song:track", Label: "Track", Type: "input", Placeholder: "The same as music:album:track but in reverse..."},
			{Key: "music:musician", Label: "Musician", Type: "input", Placeholder: "The musician that made this song..."},
			{Key: "music:release_date", Label: "Release date", Type: "input", Placeholder: "The date the album was released..."},
		},
	},
	"music.playlist": {
		Name: "Playlist details",
		Elements: []Element{
			{Key: "music:song", Label: "Song", Type: "input", Placeholder: "The song on this album..."},
			{Key: "music:song:disc", Label: "Disc", Type: "input", Placeholder: "The same as music:album:disc but in reverse..."},
			{Key: "music:song:track", Label: "Track", Type: "input", Placeholder: "The same as music:album:track but in reverse..."},
			{Key: "music:creator", Label: "Creator", Type: "input", Placeholder: "The creator of this playlist..."},
		},
	},
	"music.radio_station": {
		Name: "Radio station details",
		Elements: []Element{
			{Key: "music:creator", Label: "Creator", Type: "input", Placeholder: "The creator of this radio station..."},
		},
	},
	"video.movie":   videoMovieSchema,
	"video.episode": {
		Name: "Video episode details",
		Elements: append(
			append([]Element{}, videoMovieSchema.Elements...),
			Element{Key: "video:series", Label: "Series", Type: "input", Placeholder: "Which series this episode belongs to..."},
		),
	},
	"video.tv_show": {
		Name:     "TV show details",
		Elements: append([]Element{}, videoMovieSchema.Elements...),
	},
	"video.other": {
		Name:     "Other video details",
		Elements: append([]Element{}, videoMovieSchema.Elements...),
	},
}

func schemas() schemasResult {
	return schemasResult{
		Base:  baseSchemas,
		Types: typeSchemas,
	}
}