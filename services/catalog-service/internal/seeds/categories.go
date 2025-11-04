package seeds

import "github.com/google/uuid"

// CategorySeed describes a catalog category and ensures a stable UUID across environments.
type CategorySeed struct {
	ID          uuid.UUID
	Name        string
	Slug        string
	Description string
}

var categorySeeds = []CategorySeed{
	{
		ID:          UUIDFromString("category:science-fiction"),
		Name:        "Science Fiction",
		Slug:        "science-fiction",
		Description: "Futuristic stories featuring advanced science, technology, space exploration, and speculative worlds.",
	},
	{
		ID:          UUIDFromString("category:fantasy"),
		Name:        "Fantasy",
		Slug:        "fantasy",
		Description: "Epic quests, mythical creatures, and magical worlds that push the boundaries of imagination.",
	},
	{
		ID:          UUIDFromString("category:mythology"),
		Name:        "Mythology",
		Slug:        "mythology",
		Description: "Retellings and studies of classical myths, legends, and folkloric traditions.",
	},
	{
		ID:          UUIDFromString("category:horror"),
		Name:        "Horror",
		Slug:        "horror",
		Description: "Stories designed to thrill, unsettle, and keep readers on the edge.",
	},
	{
		ID:          UUIDFromString("category:mystery"),
		Name:        "Mystery",
		Slug:        "mystery",
		Description: "Whodunits, investigations, and suspenseful tales loaded with twists.",
	},
	{
		ID:          UUIDFromString("category:romance"),
		Name:        "Romance",
		Slug:        "romance",
		Description: "Love stories ranging from slow-burn relationships to dramatic second chances.",
	},
	{
		ID:          UUIDFromString("category:thriller"),
		Name:        "Thriller",
		Slug:        "thriller",
		Description: "High-stakes, fast-paced narratives packed with tension and danger.",
	},
	{
		ID:          UUIDFromString("category:financial"),
		Name:        "Financial",
		Slug:        "financial",
		Description: "Books on investing, markets, and personal finance strategies.",
	},
	{
		ID:          UUIDFromString("category:business"),
		Name:        "Business",
		Slug:        "business",
		Description: "Leadership, entrepreneurship, and company-building insights.",
	},
	{
		ID:          UUIDFromString("category:self-help"),
		Name:        "Self-Help",
		Slug:        "self-help",
		Description: "Guides for mindset, productivity, and personal growth.",
	},
	{
		ID:          UUIDFromString("category:biography"),
		Name:        "Biography",
		Slug:        "biography",
		Description: "Life stories and memoirs from inspiring figures across disciplines.",
	},
	{
		ID:          UUIDFromString("category:history"),
		Name:        "History",
		Slug:        "history",
		Description: "Accounts of past events, civilizations, and pivotal world moments.",
	},
	{
		ID:          UUIDFromString("category:technology"),
		Name:        "Technology",
		Slug:        "technology",
		Description: "Innovation, software, and the impact of tech on society.",
	},
	{
		ID:          UUIDFromString("category:young-adult"),
		Name:        "Young Adult",
		Slug:        "young-adult",
		Description: "Coming-of-age stories and adventures for younger audiences and crossover readers.",
	},
	{
		ID:          UUIDFromString("category:philosophy"),
		Name:        "Philosophy",
		Slug:        "philosophy",
		Description: "Explorations of ideas, ethics, and the nature of existence.",
	},
	{
		ID:          UUIDFromString("category:science-nature"),
		Name:        "Science & Nature",
		Slug:        "science-nature",
		Description: "Popular science, natural history, and explorations of the world around us.",
	},
	{
		ID:          UUIDFromString("category:health-wellness"),
		Name:        "Health & Wellness",
		Slug:        "health-wellness",
		Description: "Nutrition, mental health, fitness, and holistic wellbeing.",
	},
	{
		ID:          UUIDFromString("category:cooking"),
		Name:        "Cooking",
		Slug:        "cooking",
		Description: "Cookbooks, culinary journeys, and kitchen inspiration.",
	},
	{
		ID:          UUIDFromString("category:travel-adventure"),
		Name:        "Travel & Adventure",
		Slug:        "travel-adventure",
		Description: "Travelogues, adventure tales, and destination guides.",
	},
}

// GetCategories returns a copy of the configured category seeds.
func GetCategories() []CategorySeed {
	result := make([]CategorySeed, len(categorySeeds))
	copy(result, categorySeeds)
	return result
}
