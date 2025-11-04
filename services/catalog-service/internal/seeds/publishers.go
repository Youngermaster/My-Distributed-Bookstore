package seeds

import "github.com/google/uuid"

// PublisherSeed provides metadata for publishers referenced by multiple books.
type PublisherSeed struct {
	ID          uuid.UUID
	Code        string
	Name        string
	Country     string
	Website     string
	Description string
}

var publisherSeeds = []PublisherSeed{
	{
		ID:          UUIDFromString("publisher:oreilly-media"),
		Code:        "oreilly-media",
		Name:        "O'Reilly Media",
		Country:     "USA",
		Website:     "https://www.oreilly.com",
		Description: "Technology and business learning platform publishing hands-on guides for professionals.",
	},
	{
		ID:          UUIDFromString("publisher:tor-books"),
		Code:        "tor-books",
		Name:        "Tor Books",
		Country:     "USA",
		Website:     "https://www.tor.com",
		Description: "Leading publisher of science fiction and fantasy across adult and young adult markets.",
	},
	{
		ID:          UUIDFromString("publisher:orbit-books"),
		Code:        "orbit-books",
		Name:        "Orbit Books",
		Country:     "USA",
		Website:     "https://www.orbitbooks.net",
		Description: "Imprint focused on imaginative and epic speculative fiction.",
	},
	{
		ID:          UUIDFromString("publisher:penguin-random-house"),
		Code:        "penguin-random-house",
		Name:        "Penguin Random House",
		Country:     "USA",
		Website:     "https://www.penguinrandomhouse.com",
		Description: "Global trade publisher with diverse storytelling across every genre.",
	},
	{
		ID:          UUIDFromString("publisher:crown"),
		Code:        "crown",
		Name:        "Crown Publishing Group",
		Country:     "USA",
		Website:     "https://www.penguinrandomhouse.com",
		Description: "Imprint of Penguin Random House focusing on narrative nonfiction and quality fiction.",
	},
	{
		ID:          UUIDFromString("publisher:harpercollins"),
		Code:        "harpercollins",
		Name:        "HarperCollins",
		Country:     "USA",
		Website:     "https://www.harpercollins.com",
		Description: "One of the world's largest publishing companies with imprints spanning fiction and nonfiction.",
	},
	{
		ID:          UUIDFromString("publisher:hachette-book-group"),
		Code:        "hachette-book-group",
		Name:        "Hachette Book Group",
		Country:     "France",
		Website:     "https://www.hachettebookgroup.com",
		Description: "Major publisher offering bestselling fiction, nonfiction, and children's books.",
	},
	{
		ID:          UUIDFromString("publisher:little-brown"),
		Code:        "little-brown",
		Name:        "Little, Brown and Company",
		Country:     "USA",
		Website:     "https://www.littlebrown.com",
		Description: "Literary imprint of Hachette producing acclaimed fiction and nonfiction titles.",
	},
	{
		ID:          UUIDFromString("publisher:simon-schuster"),
		Code:        "simon-schuster",
		Name:        "Simon & Schuster",
		Country:     "USA",
		Website:     "https://www.simonandschuster.com",
		Description: "Publisher recognized for commercial and literary successes across genres.",
	},
	{
		ID:          UUIDFromString("publisher:bloomsbury"),
		Code:        "bloomsbury",
		Name:        "Bloomsbury Publishing",
		Country:     "UK",
		Website:     "https://www.bloomsbury.com",
		Description: "Independent worldwide publisher known for literary and academic lists.",
	},
	{
		ID:          UUIDFromString("publisher:scholastic-press"),
		Code:        "scholastic-press",
		Name:        "Scholastic Press",
		Country:     "USA",
		Website:     "https://www.scholastic.com",
		Description: "Children's publisher providing engaging fiction, nonfiction, and educational materials.",
	},
	{
		ID:          UUIDFromString("publisher:macmillan"),
		Code:        "macmillan",
		Name:        "Macmillan Publishers",
		Country:     "USA",
		Website:     "https://us.macmillan.com",
		Description: "Global publisher with strong portfolios in fiction, science, and education.",
	},
	{
		ID:          UUIDFromString("publisher:wiley-business"),
		Code:        "wiley-business",
		Name:        "Wiley Business",
		Country:     "USA",
		Website:     "https://www.wiley.com",
		Description: "Business and professional imprint publishing finance, leadership, and technical titles.",
	},
	{
		ID:          UUIDFromString("publisher:ww-norton"),
		Code:        "ww-norton",
		Name:        "W. W. Norton & Company",
		Country:     "USA",
		Website:     "https://wwnorton.com",
		Description: "Employee-owned publisher known for authoritative nonfiction and literary works.",
	},
	{
		ID:          UUIDFromString("publisher:national-geographic"),
		Code:        "national-geographic",
		Name:        "National Geographic",
		Country:     "USA",
		Website:     "https://www.nationalgeographic.com",
		Description: "Publisher of visually rich books on science, exploration, and culture.",
	},
}

// GetPublishers returns a copy of the configured publisher seeds.
func GetPublishers() []PublisherSeed {
	result := make([]PublisherSeed, len(publisherSeeds))
	copy(result, publisherSeeds)
	return result
}
