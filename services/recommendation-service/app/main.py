from fastapi import FastAPI
from fastapi.middleware.cors import CORSMiddleware
import uvicorn

app = FastAPI(
    title="Recommendation Service",
    description="ML-powered book recommendations service",
    version="1.0.0"
)

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
        "service": "recommendation-service"
    }

@app.get("/")
async def root():
    return {"message": "Recommendation Service API"}

# TODO: Add database connection
# TODO: Implement collaborative filtering
# TODO: Implement content-based filtering
# TODO: Add recommendation routes
# TODO: Add gRPC server
# TODO: Add RabbitMQ consumers
# TODO: Add model training pipeline

if __name__ == "__main__":
    uvicorn.run(app, host="0.0.0.0", port=8089)
