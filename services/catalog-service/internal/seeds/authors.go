package seeds

import "github.com/google/uuid"

// AuthorSeed captures metadata for an author including a stable identifier
// and lightweight profile information used by the sample catalog.
type AuthorSeed struct {
	ID       uuid.UUID
	Code     string
	Name     string
	Bio      string
	Country  string
	ImageURL string
}

var authorSeeds = []AuthorSeed{
	{
		ID:      UUIDFromString("author:nk-jemisin"),
		Code:    "nk-jemisin",
		Name:    "N. K. Jemisin",
		Bio:     "Award-winning speculative fiction author known for rich world-building and sociopolitical commentary.",
		Country: "USA",
	},
	{
		ID:      UUIDFromString("author:brandon-sanderson"),
		Code:    "brandon-sanderson",
		Name:    "Brandon Sanderson",
		Bio:     "Prolific fantasy author celebrated for intricate magic systems and expansive universes.",
		Country: "USA",
	},
	{
		ID:      UUIDFromString("author:ursula-k-le-guin"),
		Code:    "ursula-k-le-guin",
		Name:    "Ursula K. Le Guin",
		Bio:     "Iconic voice in science fiction and fantasy whose work explores anthropology, gender, and society.",
		Country: "USA",
	},
	{
		ID:      UUIDFromString("author:james-sa-corey"),
		Code:    "james-sa-corey",
		Name:    "James S. A. Corey",
		Bio:     "Collaborative pen name of Daniel Abraham and Ty Franck, authors of The Expanse series.",
		Country: "USA",
	},
	{
		ID:      UUIDFromString("author:tana-french"),
		Code:    "tana-french",
		Name:    "Tana French",
		Bio:     "Critically acclaimed mystery writer blending police procedural elements with literary depth.",
		Country: "Ireland",
	},
	{
		ID:      UUIDFromString("author:gillian-flynn"),
		Code:    "gillian-flynn",
		Name:    "Gillian Flynn",
		Bio:     "Bestselling thriller author renowned for psychological twists and unreliable narrators.",
		Country: "USA",
	},
	{
		ID:      UUIDFromString("author:stephen-king"),
		Code:    "stephen-king",
		Name:    "Stephen King",
		Bio:     "Master of modern horror and suspense with decades of genre-defining novels.",
		Country: "USA",
	},
	{
		ID:      UUIDFromString("author:silvia-moreno-garcia"),
		Code:    "silvia-moreno-garcia",
		Name:    "Silvia Moreno-Garcia",
		Bio:     "Mexican-Canadian author blending horror, fantasy, and noir with lush storytelling.",
		Country: "Mexico",
	},
	{
		ID:      UUIDFromString("author:neil-gaiman"),
		Code:    "neil-gaiman",
		Name:    "Neil Gaiman",
		Bio:     "Cross-genre storyteller whose work spans mythic fantasy, comics, and children's literature.",
		Country: "UK",
	},
	{
		ID:      UUIDFromString("author:margaret-atwood"),
		Code:    "margaret-atwood",
		Name:    "Margaret Atwood",
		Bio:     "Renowned Canadian author exploring power, society, and speculative futures.",
		Country: "Canada",
	},
	{
		ID:      UUIDFromString("author:agatha-christie"),
		Code:    "agatha-christie",
		Name:    "Agatha Christie",
		Bio:     "Legendary mystery writer and creator of beloved detectives Hercule Poirot and Miss Marple.",
		Country: "UK",
	},
	{
		ID:      UUIDFromString("author:colson-whitehead"),
		Code:    "colson-whitehead",
		Name:    "Colson Whitehead",
		Bio:     "Pulitzer-winning novelist whose work often examines American history and culture.",
		Country: "USA",
	},
	{
		ID:      UUIDFromString("author:chimamanda-ngozi-adichie"),
		Code:    "chimamanda-ngozi-adichie",
		Name:    "Chimamanda Ngozi Adichie",
		Bio:     "Esteemed Nigerian writer spotlighting identity, feminism, and diaspora experiences.",
		Country: "Nigeria",
	},
	{
		ID:      UUIDFromString("author:michelle-obama"),
		Code:    "michelle-obama",
		Name:    "Michelle Obama",
		Bio:     "Former First Lady of the United States and bestselling memoirist.",
		Country: "USA",
	},
	{
		ID:      UUIDFromString("author:walter-isaacson"),
		Code:    "walter-isaacson",
		Name:    "Walter Isaacson",
		Bio:     "Biographer of innovators and leaders, known for deeply researched narratives.",
		Country: "USA",
	},
	{
		ID:      UUIDFromString("author:yuval-noah-harari"),
		Code:    "yuval-noah-harari",
		Name:    "Yuval Noah Harari",
		Bio:     "Historian and philosopher dissecting humanity's past and future in accessible prose.",
		Country: "Israel",
	},
	{
		ID:      UUIDFromString("author:cal-newport"),
		Code:    "cal-newport",
		Name:    "Cal Newport",
		Bio:     "Computer science professor and productivity expert advocating deep work and focus.",
		Country: "USA",
	},
	{
		ID:      UUIDFromString("author:brene-brown"),
		Code:    "brene-brown",
		Name:    "Brené Brown",
		Bio:     "Research professor studying vulnerability, courage, and wholehearted leadership.",
		Country: "USA",
	},
	{
		ID:      UUIDFromString("author:morgan-housel"),
		Code:    "morgan-housel",
		Name:    "Morgan Housel",
		Bio:     "Financial writer translating behavioral economics into practical investing insights.",
		Country: "USA",
	},
	{
		ID:      UUIDFromString("author:robert-greene"),
		Code:    "robert-greene",
		Name:    "Robert Greene",
		Bio:     "Author of strategic nonfiction exploring power dynamics and mastery.",
		Country: "USA",
	},
	{
		ID:      UUIDFromString("author:david-baldacci"),
		Code:    "david-baldacci",
		Name:    "David Baldacci",
		Bio:     "Thriller author delivering high-stakes political suspense and action.",
		Country: "USA",
	},
	{
		ID:      UUIDFromString("author:madeline-miller"),
		Code:    "madeline-miller",
		Name:    "Madeline Miller",
		Bio:     "Classically trained author reimagining ancient myths with modern sensitivity.",
		Country: "USA",
	},
	{
		ID:      UUIDFromString("author:haruki-murakami"),
		Code:    "haruki-murakami",
		Name:    "Haruki Murakami",
		Bio:     "Internationally acclaimed novelist blending magical realism with introspective themes.",
		Country: "Japan",
	},
	{
		ID:      UUIDFromString("author:andy-weir"),
		Code:    "andy-weir",
		Name:    "Andy Weir",
		Bio:     "Science-driven storyteller known for grounded, high-stakes space adventures.",
		Country: "USA",
	},
	{
		ID:      UUIDFromString("author:sally-rooney"),
		Code:    "sally-rooney",
		Name:    "Sally Rooney",
		Bio:     "Contemporary novelist capturing millennial relationships with sharp dialogue.",
		Country: "Ireland",
	},
	{
		ID:      UUIDFromString("author:rebecca-skloot"),
		Code:    "rebecca-skloot",
		Name:    "Rebecca Skloot",
		Bio:     "Science writer bringing investigative rigor to medical history narratives.",
		Country: "USA",
	},
	{
		ID:      UUIDFromString("author:anthony-bourdain"),
		Code:    "anthony-bourdain",
		Name:    "Anthony Bourdain",
		Bio:     "Chef, traveler, and storyteller who celebrated global food culture.",
		Country: "USA",
	},
	{
		ID:      UUIDFromString("author:bill-bryson"),
		Code:    "bill-bryson",
		Name:    "Bill Bryson",
		Bio:     "Humorous travel and science author uncovering curiosities of the world.",
		Country: "USA",
	},
	{
		ID:      UUIDFromString("author:mary-roach"),
		Code:    "mary-roach",
		Name:    "Mary Roach",
		Bio:     "Popular science writer known for witty dives into overlooked topics.",
		Country: "USA",
	},
	{
		ID:      UUIDFromString("author:neil-degrasse-tyson"),
		Code:    "neil-degrasse-tyson",
		Name:    "Neil deGrasse Tyson",
		Bio:     "Astrophysicist and science communicator making the cosmos accessible to all.",
		Country: "USA",
	},
}

// GetAuthors returns a copy of the configured author seeds.
func GetAuthors() []AuthorSeed {
	result := make([]AuthorSeed, len(authorSeeds))
	copy(result, authorSeeds)
	return result
}
