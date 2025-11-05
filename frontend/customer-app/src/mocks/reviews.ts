// Mock review data for UI development
// This will be replaced with actual API calls when review service is integrated

export interface Review {
  id: string;
  book_id: string;
  user_id: string;
  user_name: string;
  rating: number; // 1-5
  title: string;
  content: string;
  helpful_votes: number;
  verified_purchase: boolean;
  sentiment_score?: number; // -1 to 1
  sentiment_label?: "positive" | "negative" | "neutral";
  created_at: string;
  updated_at: string;
}

export interface ReviewStats {
  average_rating: number;
  total_reviews: number;
  rating_distribution: {
    5: number;
    4: number;
    3: number;
    2: number;
    1: number;
  };
}

// Mock reviews for different books
export const mockReviews: Record<string, Review[]> = {
  // This will be populated with book IDs as keys
  default: [
    {
      id: "review-1",
      book_id: "book-1",
      user_id: "user-1",
      user_name: "Sarah Johnson",
      rating: 5,
      title: "Absolutely brilliant!",
      content:
        "This book exceeded all my expectations. The writing is engaging, the concepts are explained clearly, and I learned so much. Highly recommended for anyone interested in the subject.",
      helpful_votes: 42,
      verified_purchase: true,
      sentiment_score: 0.95,
      sentiment_label: "positive",
      created_at: "2024-10-15T10:30:00Z",
      updated_at: "2024-10-15T10:30:00Z",
    },
    {
      id: "review-2",
      book_id: "book-1",
      user_id: "user-2",
      user_name: "Michael Chen",
      rating: 4,
      title: "Great read, with minor issues",
      content:
        "Overall, this is an excellent book. The author does a fantastic job explaining complex topics. However, some sections could have been more concise. Still, I'd recommend it to others.",
      helpful_votes: 28,
      verified_purchase: true,
      sentiment_score: 0.7,
      sentiment_label: "positive",
      created_at: "2024-10-12T14:20:00Z",
      updated_at: "2024-10-12T14:20:00Z",
    },
    {
      id: "review-3",
      book_id: "book-1",
      user_id: "user-3",
      user_name: "Emily Rodriguez",
      rating: 5,
      title: "A must-read masterpiece",
      content:
        "I've read many books on this topic, and this one stands out. The practical examples are incredibly helpful, and the writing style makes even complex concepts accessible. Worth every penny!",
      helpful_votes: 35,
      verified_purchase: true,
      sentiment_score: 0.92,
      sentiment_label: "positive",
      created_at: "2024-10-08T09:15:00Z",
      updated_at: "2024-10-08T09:15:00Z",
    },
    {
      id: "review-4",
      book_id: "book-1",
      user_id: "user-4",
      user_name: "David Thompson",
      rating: 3,
      title: "Good but not great",
      content:
        "The book covers the basics well, but I was hoping for more advanced topics. It's a solid introduction for beginners, but experienced readers might find it too elementary.",
      helpful_votes: 15,
      verified_purchase: false,
      sentiment_score: 0.1,
      sentiment_label: "neutral",
      created_at: "2024-10-05T16:45:00Z",
      updated_at: "2024-10-05T16:45:00Z",
    },
    {
      id: "review-5",
      book_id: "book-1",
      user_id: "user-5",
      user_name: "Jennifer Park",
      rating: 5,
      title: "Changed my perspective completely",
      content:
        "This book opened my eyes to new ways of thinking. The author presents ideas in such a compelling way that it's hard to put down. I found myself taking notes throughout. Absolutely transformative!",
      helpful_votes: 51,
      verified_purchase: true,
      sentiment_score: 0.98,
      sentiment_label: "positive",
      created_at: "2024-09-28T11:00:00Z",
      updated_at: "2024-09-28T11:00:00Z",
    },
    {
      id: "review-6",
      book_id: "book-1",
      user_id: "user-6",
      user_name: "Robert Martinez",
      rating: 4,
      title: "Solid and well-researched",
      content:
        "The research behind this book is impressive. The author clearly knows their stuff. A few chapters felt a bit dry, but overall it's a valuable resource that I'll refer back to often.",
      helpful_votes: 22,
      verified_purchase: true,
      sentiment_score: 0.65,
      sentiment_label: "positive",
      created_at: "2024-09-22T13:30:00Z",
      updated_at: "2024-09-22T13:30:00Z",
    },
    {
      id: "review-7",
      book_id: "book-1",
      user_id: "user-7",
      user_name: "Amanda Wilson",
      rating: 2,
      title: "Disappointed",
      content:
        "I had high hopes based on the reviews, but this book didn't live up to the hype for me. The writing style wasn't my cup of tea, and some examples felt dated. It's not terrible, but I expected more.",
      helpful_votes: 8,
      verified_purchase: true,
      sentiment_score: -0.45,
      sentiment_label: "negative",
      created_at: "2024-09-18T08:20:00Z",
      updated_at: "2024-09-18T08:20:00Z",
    },
  ],
};

// Mock review stats
export const mockReviewStats: Record<string, ReviewStats> = {
  default: {
    average_rating: 4.3,
    total_reviews: 7,
    rating_distribution: {
      5: 3,
      4: 2,
      3: 1,
      2: 1,
      1: 0,
    },
  },
};

// Helper function to get reviews for a book
export function getReviewsForBook(bookId: string): Review[] {
  return mockReviews[bookId] || mockReviews.default;
}

// Helper function to get review stats for a book
export function getReviewStatsForBook(bookId: string): ReviewStats {
  return mockReviewStats[bookId] || mockReviewStats.default;
}

// Helper function to format date
export function formatReviewDate(dateString: string): string {
  const date = new Date(dateString);
  const now = new Date();
  const diffInMs = now.getTime() - date.getTime();
  const diffInDays = Math.floor(diffInMs / (1000 * 60 * 60 * 24));

  if (diffInDays === 0) {
    return "Today";
  } else if (diffInDays === 1) {
    return "Yesterday";
  } else if (diffInDays < 7) {
    return `${diffInDays} days ago`;
  } else if (diffInDays < 30) {
    const weeks = Math.floor(diffInDays / 7);
    return `${weeks} ${weeks === 1 ? "week" : "weeks"} ago`;
  } else if (diffInDays < 365) {
    const months = Math.floor(diffInDays / 30);
    return `${months} ${months === 1 ? "month" : "months"} ago`;
  } else {
    return date.toLocaleDateString();
  }
}
