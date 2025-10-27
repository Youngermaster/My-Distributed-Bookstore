import { Card, CardContent } from "@/components/ui/card";
import { type Category } from "@/types/book";
import { ChevronRight } from "lucide-react";

interface GenreCardProps {
  category: Category;
  bookCount?: number;
  onClick?: () => void;
}

const genreIcons: Record<string, string> = {
  programming: "💻",
  "distributed-systems": "🌐",
  "software-architecture": "🏗️",
  databases: "🗄️",
  "cloud-computing": "☁️",
  fiction: "📖",
  "science-fiction": "🚀",
  mystery: "🔍",
  romance: "❤️",
  fantasy: "🐉",
  thriller: "😱",
  biography: "👤",
  history: "📜",
  science: "🔬",
  default: "📚",
};

export default function GenreCard({
  category,
  bookCount,
  onClick,
}: GenreCardProps) {
  const icon = genreIcons[category.slug] || genreIcons.default;

  return (
    <Card
      className="overflow-hidden hover:shadow-lg transition-all cursor-pointer group hover:scale-105"
      onClick={onClick}
    >
      <CardContent className="p-6">
        <div className="flex items-center justify-between">
          <div className="flex-1">
            <div className="text-4xl mb-3">{icon}</div>
            <h3 className="font-semibold text-xl mb-1 group-hover:text-blue-600 transition-colors">
              {category.name}
            </h3>
            {bookCount !== undefined && (
              <p className="text-sm text-gray-500">
                {bookCount} {bookCount === 1 ? "book" : "books"}
              </p>
            )}
          </div>
          <ChevronRight className="h-6 w-6 text-gray-400 group-hover:text-blue-600 group-hover:translate-x-1 transition-all" />
        </div>
      </CardContent>
    </Card>
  );
}
