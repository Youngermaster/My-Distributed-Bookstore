import { useState } from "react";
import {
  Card,
  CardContent,
  CardHeader,
  CardTitle,
  CardDescription,
} from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Textarea } from "@/components/ui/textarea";
import { Star, AlertCircle } from "lucide-react";
import { Alert, AlertDescription } from "@/components/ui/alert";

interface ReviewFormProps {
  bookId: string;
  bookTitle: string;
  isAuthenticated: boolean;
  onSubmit?: (data: ReviewFormData) => void;
}

export interface ReviewFormData {
  rating: number;
  title: string;
  content: string;
}

export default function ReviewForm({
  bookId: _bookId, // Prefix with _ to indicate intentionally unused (for future integration)
  bookTitle,
  isAuthenticated,
  onSubmit,
}: ReviewFormProps) {
  const [rating, setRating] = useState(0);
  const [hoverRating, setHoverRating] = useState(0);
  const [title, setTitle] = useState("");
  const [content, setContent] = useState("");

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    if (rating > 0 && title.trim() && content.trim()) {
      onSubmit?.({
        rating,
        title: title.trim(),
        content: content.trim(),
      });
      // Reset form
      setRating(0);
      setTitle("");
      setContent("");
    }
  };

  // Review service not yet deployed - show disabled form
  const isDisabled = true;

  return (
    <Card>
      <CardHeader>
        <CardTitle>Write a Review</CardTitle>
        <CardDescription>Share your thoughts about {bookTitle}</CardDescription>
      </CardHeader>
      <CardContent>
        {!isAuthenticated ? (
          <Alert>
            <AlertCircle className="h-4 w-4" />
            <AlertDescription>
              You must be logged in to write a review.
            </AlertDescription>
          </Alert>
        ) : isDisabled ? (
          <Alert>
            <AlertCircle className="h-4 w-4" />
            <AlertDescription>
              Review submission is coming soon! The review service is currently
              being integrated.
            </AlertDescription>
          </Alert>
        ) : (
          <form onSubmit={handleSubmit} className="space-y-4">
            {/* Rating Selection */}
            <div>
              <label className="block text-sm font-medium text-foreground mb-2">
                Your Rating *
              </label>
              <div className="flex items-center gap-2">
                {Array.from({ length: 5 }).map((_, i) => {
                  const starValue = i + 1;
                  return (
                    <button
                      key={i}
                      type="button"
                      onClick={() => setRating(starValue)}
                      onMouseEnter={() => setHoverRating(starValue)}
                      onMouseLeave={() => setHoverRating(0)}
                      className="transition-transform hover:scale-110"
                    >
                      <Star
                        className={`h-8 w-8 ${
                          starValue <= (hoverRating || rating)
                            ? "fill-yellow-400 text-yellow-400"
                            : "text-gray-300"
                        }`}
                      />
                    </button>
                  );
                })}
                {rating > 0 && (
                  <span className="ml-2 text-sm text-muted-foreground">
                    {rating} {rating === 1 ? "star" : "stars"}
                  </span>
                )}
              </div>
            </div>

            {/* Review Title */}
            <div>
              <label className="block text-sm font-medium text-foreground mb-2">
                Review Title *
              </label>
              <Input
                type="text"
                value={title}
                onChange={(e) => setTitle(e.target.value)}
                placeholder="Sum up your review in one line"
                maxLength={100}
                required
              />
            </div>

            {/* Review Content */}
            <div>
              <label className="block text-sm font-medium text-foreground mb-2">
                Your Review *
              </label>
              <Textarea
                value={content}
                onChange={(e: React.ChangeEvent<HTMLTextAreaElement>) =>
                  setContent(e.target.value)
                }
                placeholder="Share your thoughts about this book..."
                rows={6}
                maxLength={2000}
                required
              />
              <p className="text-xs text-muted-foreground mt-1">
                {content.length}/2000 characters
              </p>
            </div>

            {/* Submit Button */}
            <Button
              type="submit"
              disabled={rating === 0 || !title.trim() || !content.trim()}
              className="w-full"
            >
              Submit Review
            </Button>
          </form>
        )}
      </CardContent>
    </Card>
  );
}
