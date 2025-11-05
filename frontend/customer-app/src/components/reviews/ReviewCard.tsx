import { Card, CardContent } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Star, ThumbsUp, BadgeCheck } from "lucide-react";
import type { Review } from "@/mocks/reviews";
import { formatReviewDate } from "@/mocks/reviews";

interface ReviewCardProps {
  review: Review;
  onHelpful?: (reviewId: string) => void;
}

export default function ReviewCard({ review, onHelpful }: ReviewCardProps) {
  const {
    id,
    user_name,
    rating,
    title,
    content,
    helpful_votes,
    verified_purchase,
    created_at,
  } = review;

  return (
    <Card>
      <CardContent className="pt-6 space-y-4">
        {/* Header */}
        <div className="flex items-start justify-between">
          <div>
            <div className="flex items-center gap-2 mb-1">
              <span className="font-semibold text-foreground">{user_name}</span>
              {verified_purchase && (
                <div className="flex items-center gap-1 text-sm text-green-600">
                  <BadgeCheck className="h-4 w-4" />
                  <span>Verified Purchase</span>
                </div>
              )}
            </div>
            <div className="text-sm text-muted-foreground">
              {formatReviewDate(created_at)}
            </div>
          </div>

          {/* Rating */}
          <div className="flex items-center gap-1">
            {Array.from({ length: 5 }).map((_, i) => (
              <Star
                key={i}
                className={`h-4 w-4 ${
                  i < rating
                    ? "fill-yellow-400 text-yellow-400"
                    : "text-gray-300"
                }`}
              />
            ))}
          </div>
        </div>

        {/* Review Title */}
        <div>
          <h3 className="font-semibold text-lg text-foreground">{title}</h3>
        </div>

        {/* Review Content */}
        <div>
          <p className="text-muted-foreground leading-relaxed">{content}</p>
        </div>

        {/* Helpful Button */}
        <div className="flex items-center gap-4 pt-2">
          <Button
            variant="ghost"
            size="sm"
            onClick={() => onHelpful?.(id)}
            className="text-muted-foreground hover:text-foreground"
          >
            <ThumbsUp className="h-4 w-4 mr-2" />
            Helpful {helpful_votes > 0 && `(${helpful_votes})`}
          </Button>
        </div>
      </CardContent>
    </Card>
  );
}
