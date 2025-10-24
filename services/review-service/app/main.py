from fastapi import FastAPI
from fastapi.middleware.cors import CORSMiddleware
import uvicorn

app = FastAPI(
    title="Review Service",
    description="Book reviews and ratings service",
    version="1.0.0"
)

# CORS middleware
app.add_middleware(
    CORSMiddleware,
    allow_origins=["*"],
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)

@app.get("/health")
async def health_check():
    return {
        "status": "ok",
        "service": "review-service"
    }

@app.get("/")
async def root():
    return {"message": "Review Service API"}

# TODO: Add database connection
# TODO: Add review routes
# TODO: Add gRPC server
# TODO: Add sentiment analysis
# TODO: Add RabbitMQ consumers

if __name__ == "__main__":
    uvicorn.run(app, host="0.0.0.0", port=8088)
