import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Star } from "lucide-react";
import type { ReviewStatsResponse } from "@/types/review";

interface ReviewStatsProps {
  stats: ReviewStatsResponse;
}

export default function ReviewStats({ stats }: ReviewStatsProps) {
  const { average_rating, total_reviews, rating_distribution } = stats;

  const distribution = {
    5: rating_distribution[5] ?? 0,
    4: rating_distribution[4] ?? 0,
    3: rating_distribution[3] ?? 0,
    2: rating_distribution[2] ?? 0,
    1: rating_distribution[1] ?? 0,
  };

  const ratingPercentages = {
    5: total_reviews ? (distribution[5] / total_reviews) * 100 : 0,
    4: total_reviews ? (distribution[4] / total_reviews) * 100 : 0,
    3: total_reviews ? (distribution[3] / total_reviews) * 100 : 0,
    2: total_reviews ? (distribution[2] / total_reviews) * 100 : 0,
    1: total_reviews ? (distribution[1] / total_reviews) * 100 : 0,
  };

  return (
    <Card>
      <CardHeader>
        <CardTitle>Customer Reviews</CardTitle>
      </CardHeader>
      <CardContent className="space-y-6">
        {/* Overall Rating */}
        <div className="flex items-center gap-8">
          <div className="text-center">
            <div className="text-5xl font-bold text-foreground mb-1">
              {total_reviews ? average_rating.toFixed(1) : "0.0"}
            </div>
            <div className="flex items-center gap-1 justify-center mb-1">
              {Array.from({ length: 5 }).map((_, i) => (
                <Star
                  key={i}
                  className={`h-5 w-5 ${
                    i < Math.floor(average_rating)
                      ? "fill-yellow-400 text-yellow-400"
                      : "text-gray-300"
                  }`}
                />
              ))}
            </div>
            <div className="text-sm text-muted-foreground">
              {total_reviews} {total_reviews === 1 ? "review" : "reviews"}
            </div>
          </div>

          {/* Rating Distribution */}
          <div className="flex-1 space-y-2">
            {([5, 4, 3, 2, 1] as const).map((rating) => (
              <div key={rating} className="flex items-center gap-2">
                <div className="flex items-center gap-1 w-16">
                  <span className="text-sm font-medium text-foreground">
                    {rating}
                  </span>
                  <Star className="h-3 w-3 fill-yellow-400 text-yellow-400" />
                </div>
                <div className="flex-1 bg-gray-200 rounded-full h-2">
                  <div
                    className="bg-yellow-400 h-2 rounded-full transition-all duration-300"
                    style={{ width: `${ratingPercentages[rating]}%` }}
                  />
                </div>
                <div className="text-sm text-muted-foreground w-16 text-right">
                  {distribution[rating]}{" "}
                  {distribution[rating] === 1 ? "review" : "reviews"}
                </div>
              </div>
            ))}
          </div>
        </div>
      </CardContent>
    </Card>
  );
}
