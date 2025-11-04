import { useState } from "react";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import ReviewCard from "./ReviewCard";
import type { Review } from "@/mocks/reviews";

interface ReviewListProps {
  reviews: Review[];
  onHelpful?: (reviewId: string) => void;
}

type SortOption = "recent" | "helpful" | "highest" | "lowest";

export default function ReviewList({ reviews, onHelpful }: ReviewListProps) {
  const [sortBy, setSortBy] = useState<SortOption>("recent");

  const sortedReviews = [...reviews].sort((a, b) => {
    switch (sortBy) {
      case "recent":
        return (
          new Date(b.created_at).getTime() - new Date(a.created_at).getTime()
        );
      case "helpful":
        return b.helpful_votes - a.helpful_votes;
      case "highest":
        return b.rating - a.rating;
      case "lowest":
        return a.rating - b.rating;
      default:
        return 0;
    }
  });

  if (reviews.length === 0) {
    return (
      <div className="text-center py-12">
        <p className="text-muted-foreground">
          No reviews yet. Be the first to review this book!
        </p>
      </div>
    );
  }

  return (
    <div className="space-y-6">
      {/* Sort Controls */}
      <div className="flex items-center justify-between">
        <h3 className="text-xl font-semibold text-foreground">
          {reviews.length} {reviews.length === 1 ? "Review" : "Reviews"}
        </h3>
        <Select
          value={sortBy}
          onValueChange={(value: string) => setSortBy(value as SortOption)}
        >
          <SelectTrigger className="w-[180px]">
            <SelectValue placeholder="Sort by" />
          </SelectTrigger>
          <SelectContent>
            <SelectItem value="recent">Most Recent</SelectItem>
            <SelectItem value="helpful">Most Helpful</SelectItem>
            <SelectItem value="highest">Highest Rating</SelectItem>
            <SelectItem value="lowest">Lowest Rating</SelectItem>
          </SelectContent>
        </Select>
      </div>

      {/* Reviews */}
      <div className="space-y-4">
        {sortedReviews.map((review) => (
          <ReviewCard key={review.id} review={review} onHelpful={onHelpful} />
        ))}
      </div>
    </div>
  );
}
