import express from 'express';
import cors from 'cors';
import { ApolloServer } from '@apollo/server';
import { expressMiddleware } from '@apollo/server/express4';
import { typeDefs } from './schema/typeDefs';
import { createResolvers } from './resolvers';
import { UserClient, GroupClient, ExpenseClient, SettlementClient } from './grpc-clients';
import { getAuthContext } from './middleware/auth';

const PORT = process.env.PORT || 4000;
const JWT_SECRET = process.env.JWT_SECRET || 'split-jwt-secret-key-change-in-production';

// gRPC Service URLs
const USER_SERVICE_URL = process.env.USER_SERVICE_URL || 'localhost:50051';
const GROUP_SERVICE_URL = process.env.GROUP_SERVICE_URL || 'localhost:50052';
const EXPENSE_SERVICE_URL = process.env.EXPENSE_SERVICE_URL || 'localhost:50053';
const SETTLEMENT_SERVICE_URL = process.env.SETTLEMENT_SERVICE_URL || 'localhost:50054';

async function main() {
  // Initialize gRPC clients
  const userClient = new UserClient(USER_SERVICE_URL);
  const groupClient = new GroupClient(GROUP_SERVICE_URL);
  const expenseClient = new ExpenseClient(EXPENSE_SERVICE_URL);
  const settlementClient = new SettlementClient(SETTLEMENT_SERVICE_URL);

  // Create resolvers
  const resolvers = createResolvers(userClient, groupClient, expenseClient, settlementClient);

  // Create Apollo Server
  const server = new ApolloServer({
    typeDefs,
    resolvers,
  });

  await server.start();

  // Create Express app
  const app = express();

  app.use(cors({
    origin: process.env.CORS_ORIGIN || 'http://localhost:3000',
    credentials: true,
  }));

  app.use(express.json());

  // Health check endpoint
  app.get('/health', (_req, res) => {
    res.json({ status: 'ok', service: 'split-gateway' });
  });

  // GraphQL endpoint
  app.use('/graphql', expressMiddleware(server, {
    context: async ({ req }) => {
      const auth = getAuthContext(req, JWT_SECRET);
      return {
        auth,
        userCache: new Map<string, any>(), // Per-request user cache
      };
    },
  }));

  app.listen(PORT, () => {
    console.log(`🚀 Split Gateway running at http://localhost:${PORT}/graphql`);
  });
}

main().catch(console.error);
