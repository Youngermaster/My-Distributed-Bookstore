import express, { Express, Request, Response } from 'express';
import cors from 'cors';
import helmet from 'helmet';
import dotenv from 'dotenv';

// Load environment variables
dotenv.config();

const app: Express = express();
const PORT = process.env.PORT || 8085;
const GRPC_PORT = process.env.GRPC_PORT || 50055;

// Middleware
app.use(helmet());
app.use(cors());
app.use(express.json());
app.use(express.urlencoded({ extended: true }));

// Health check endpoints
app.get('/health', (req: Request, res: Response) => {
  res.json({
    status: 'ok',
    service: 'payment-service',
    timestamp: new Date().toISOString(),
  });
});

app.get('/ready', (req: Request, res: Response) => {
  // TODO: Check database connection
  // TODO: Check Stripe connection
  // TODO: Check RabbitMQ connection
  res.json({
    status: 'ready',
    checks: {
      database: 'ok',
      stripe: 'ok',
      rabbitmq: 'ok',
    },
  });
});

// TODO: Initialize Prisma client
// TODO: Initialize Stripe client
// TODO: Setup gRPC server
// TODO: Setup RabbitMQ consumers
// TODO: Setup routes
// TODO: Setup error handling middleware
// TODO: Setup Jaeger tracing

// Stripe webhook endpoint
app.post('/webhook/stripe', express.raw({ type: 'application/json' }), (req: Request, res: Response) => {
  // TODO: Implement Stripe webhook handler
  res.json({ received: true });
});

// API routes placeholder
app.use('/api/v1', (req: Request, res: Response) => {
  res.status(404).json({ error: 'Route not implemented yet' });
});

// Error handling middleware
app.use((err: Error, req: Request, res: Response, next: any) => {
  console.error(err.stack);
  res.status(500).json({
    error: 'Internal server error',
    message: process.env.NODE_ENV === 'development' ? err.message : undefined,
  });
});

// Start HTTP server
app.listen(PORT, () => {
  console.log(`✅ Payment Service HTTP server running on port ${PORT}`);
  console.log(`📊 Environment: ${process.env.NODE_ENV}`);
  console.log(`🔌 gRPC server should be running on port ${GRPC_PORT}`);
});

// Graceful shutdown
process.on('SIGTERM', () => {
  console.log('SIGTERM signal received: closing HTTP server');
  process.exit(0);
});

process.on('SIGINT', () => {
  console.log('SIGINT signal received: closing HTTP server');
  process.exit(0);
});

export default app;
