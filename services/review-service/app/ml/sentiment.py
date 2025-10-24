"""
Sentiment analysis using TextBlob and NLTK.
"""
from typing import Tuple
import logging

logger = logging.getLogger(__name__)


def analyze_sentiment(text: str) -> Tuple[str, float]:
    """
    Analyze sentiment of review text using TextBlob.

    Args:
        text: Review content to analyze

    Returns:
        Tuple of (sentiment_label, sentiment_score)
        - sentiment_label: "positive", "negative", or "neutral"
        - sentiment_score: Float between -1.0 (negative) and 1.0 (positive)
    """
    try:
        from textblob import TextBlob

        # Create TextBlob object
        blob = TextBlob(text)

        # Get polarity score (-1 to 1)
        polarity = blob.sentiment.polarity

        # Classify sentiment based on polarity
        if polarity > 0.3:
            label = "positive"
        elif polarity < -0.3:
            label = "negative"
        else:
            label = "neutral"

        logger.debug(f"Sentiment analysis: {label} (score: {polarity:.3f})")

        return label, polarity

    except Exception as e:
        logger.error(f"Error in sentiment analysis: {e}")
        # Return neutral sentiment on error
        return "neutral", 0.0


def download_nltk_data():
    """
    Download required NLTK data for TextBlob.
    Should be called on application startup.
    """
    try:
        import nltk
        nltk.download('brown', quiet=True)
        nltk.download('punkt', quiet=True)
        nltk.download('punkt_tab', quiet=True)
        logger.info("NLTK data downloaded successfully")
    except Exception as e:
        logger.warning(f"Could not download NLTK data: {e}")
