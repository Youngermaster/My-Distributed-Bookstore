import express from 'express';
import dotenv from 'dotenv';

dotenv.config();

const app = express();
const PORT = process.env.PORT || 8087;

app.use(express.json());

app.get('/health', (req, res) => {
  res.json({ status: 'ok', service: 'notification-service' });
});

// TODO: Initialize SendGrid
// TODO: Setup RabbitMQ consumers
// TODO: Load email templates
// TODO: Implement notification handlers

app.listen(PORT, () => {
  console.log(`✅ Notification Service running on port ${PORT}`);
});

export default app;
