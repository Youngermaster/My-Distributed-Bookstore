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
  isSubmitting?: boolean;
  onSubmit?: (data: ReviewFormData) => Promise<void> | void;
}

export interface ReviewFormData {
  rating: number;
  title: string;
  content: string;
}

export default function ReviewForm({
  bookId: _bookId,
  bookTitle,
  isAuthenticated,
  isSubmitting = false,
  onSubmit,
}: ReviewFormProps) {
  const [rating, setRating] = useState(0);
  const [hoverRating, setHoverRating] = useState(0);
  const [title, setTitle] = useState("");
  const [content, setContent] = useState("");
  const [formError, setFormError] = useState<string | null>(null);
  const [formSuccess, setFormSuccess] = useState<string | null>(null);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setFormError(null);
    setFormSuccess(null);

    if (rating <= 0 || !title.trim() || !content.trim()) {
      setFormError("Please provide a rating, title, and review content.");
      return;
    }

    try {
      await onSubmit?.({
        rating,
        title: title.trim(),
        content: content.trim(),
      });

      setRating(0);
      setTitle("");
      setContent("");
      setFormSuccess("Review submitted! Thank you for your feedback.");
    } catch (err: any) {
      const message =
        err?.response?.data?.message ||
        "Failed to submit review. Please try again.";
      setFormError(message);
    }
  };

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
        ) : (
          <form onSubmit={handleSubmit} className="space-y-4">
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
                disabled={isSubmitting}
              />
            </div>

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
                disabled={isSubmitting}
              />
              <p className="text-xs text-muted-foreground mt-1">
                {content.length}/2000 characters
              </p>
            </div>

            {formError && (
              <Alert variant="destructive">
                <AlertCircle className="h-4 w-4" />
                <AlertDescription>{formError}</AlertDescription>
              </Alert>
            )}

            {formSuccess && (
              <Alert>
                <AlertCircle className="h-4 w-4" />
                <AlertDescription>{formSuccess}</AlertDescription>
              </Alert>
            )}

            <Button
              type="submit"
              disabled={
                rating === 0 || !title.trim() || !content.trim() || isSubmitting
              }
              className="w-full"
            >
              {isSubmitting ? "Submitting..." : "Submit Review"}
            </Button>
          </form>
        )}
      </CardContent>
    </Card>
  );
}
