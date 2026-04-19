import { UserClient, GroupClient, ExpenseClient, SettlementClient } from '../grpc-clients';
import { AuthContext, requireAuth } from '../middleware/auth';

const SPLIT_TYPE_MAP: Record<string, number> = {
  'EQUAL': 1,
  'CUSTOM': 2,
  'PERCENTAGE': 3,
};

const SPLIT_TYPE_REVERSE: Record<number, string> = {
  0: 'EQUAL',
  1: 'EQUAL',
  2: 'CUSTOM',
  3: 'PERCENTAGE',
};

const SETTLEMENT_STATUS_REVERSE: Record<number, string> = {
  0: 'PENDING',
  1: 'PENDING',
  2: 'COMPLETED',
};

function formatTimestamp(ts: any): string | null {
  if (!ts) return null;
  if (ts.seconds) {
    return new Date(parseInt(ts.seconds) * 1000).toISOString();
  }
  return ts;
}

function formatUser(u: any) {
  if (!u) return null;
  return {
    id: u.id,
    email: u.email,
    name: u.name,
    avatarUrl: u.avatarUrl || '',
    currency: u.currency || 'USD',
    createdAt: formatTimestamp(u.createdAt),
    updatedAt: formatTimestamp(u.updatedAt),
  };
}

function formatGroup(g: any) {
  if (!g) return null;
  return {
    id: g.id,
    name: g.name,
    description: g.description || '',
    currency: g.currency || 'USD',
    inviteCode: g.inviteCode || '',
    ownerId: g.ownerId,
    createdAt: formatTimestamp(g.createdAt),
    updatedAt: formatTimestamp(g.updatedAt),
    members: (g.members || []).map((m: any) => ({
      userId: m.userId,
      role: m.role,
      joinedAt: formatTimestamp(m.joinedAt),
    })),
  };
}

function formatExpense(e: any) {
  if (!e) return null;
  return {
    id: e.id,
    groupId: e.groupId,
    payerId: e.payerId,
    amount: parseFloat(e.amount) || 0,
    description: e.description || '',
    splitType: SPLIT_TYPE_REVERSE[e.splitType] || 'EQUAL',
    splits: (e.splits || []).map((s: any) => ({
      id: s.id,
      expenseId: s.expenseId,
      userId: s.userId,
      amount: parseFloat(s.amount) || 0,
      percentage: parseFloat(s.percentage) || 0,
    })),
    createdAt: formatTimestamp(e.createdAt),
    updatedAt: formatTimestamp(e.updatedAt),
  };
}

function formatSettlement(s: any) {
  if (!s) return null;
  return {
    id: s.id,
    groupId: s.groupId,
    fromUserId: s.fromUserId,
    toUserId: s.toUserId,
    amount: parseFloat(s.amount) || 0,
    status: SETTLEMENT_STATUS_REVERSE[s.status] || s.status || 'PENDING',
    createdAt: formatTimestamp(s.createdAt),
    updatedAt: formatTimestamp(s.updatedAt),
  };
}

export function createResolvers(
  userClient: UserClient,
  groupClient: GroupClient,
  expenseClient: ExpenseClient,
  settlementClient: SettlementClient
) {
  // Helper to resolve user by ID (with caching per request)
  const getUserById = async (userId: string, userCache: Map<string, any>): Promise<any> => {
    if (userCache.has(userId)) return userCache.get(userId);
    try {
      const result = await userClient.getUser(userId);
      const user = formatUser(result.user);
      userCache.set(userId, user);
      return user;
    } catch {
      return null;
    }
  };

  return {
    DateTime: {
      __serialize: (value: any) => value,
    },

    // ─── Field Resolvers ──────────────────────────
    GroupMember: {
      user: async (parent: any, _: any, context: any) => {
        return getUserById(parent.userId, context.userCache);
      },
    },

    Expense: {
      payer: async (parent: any, _: any, context: any) => {
        return getUserById(parent.payerId, context.userCache);
      },
    },

    ExpenseSplit: {
      user: async (parent: any, _: any, context: any) => {
        return getUserById(parent.userId, context.userCache);
      },
    },

    UserBalance: {
      user: async (parent: any, _: any, context: any) => {
        return getUserById(parent.userId, context.userCache);
      },
    },

    BalanceDetail: {
      fromUser: async (parent: any, _: any, context: any) => {
        return getUserById(parent.fromUserId, context.userCache);
      },
      toUser: async (parent: any, _: any, context: any) => {
        return getUserById(parent.toUserId, context.userCache);
      },
    },

    Transaction: {
      fromUser: async (parent: any, _: any, context: any) => {
        return getUserById(parent.fromUserId, context.userCache);
      },
      toUser: async (parent: any, _: any, context: any) => {
        return getUserById(parent.toUserId, context.userCache);
      },
    },

    Settlement: {
      fromUser: async (parent: any, _: any, context: any) => {
        return getUserById(parent.fromUserId, context.userCache);
      },
      toUser: async (parent: any, _: any, context: any) => {
        return getUserById(parent.toUserId, context.userCache);
      },
    },

    // ─── Queries ──────────────────────────────────
    Query: {
      me: async (_: any, __: any, context: { auth: AuthContext; userCache: Map<string, any> }) => {
        const userId = requireAuth(context.auth);
        const result = await userClient.getUser(userId);
        return formatUser(result.user);
      },

      user: async (_: any, { id }: { id: string }) => {
        const result = await userClient.getUser(id);
        return formatUser(result.user);
      },

      myGroups: async (_: any, __: any, context: { auth: AuthContext }) => {
        const userId = requireAuth(context.auth);
        const result = await groupClient.listUserGroups(userId);
        return (result.groups || []).map(formatGroup);
      },

      group: async (_: any, { id }: { id: string }) => {
        const result = await groupClient.getGroup(id);
        return formatGroup(result.group);
      },

      expense: async (_: any, { id }: { id: string }) => {
        const result = await expenseClient.getExpense(id);
        return formatExpense(result.expense);
      },

      groupExpenses: async (_: any, { groupId, limit, offset }: any) => {
        const result = await expenseClient.listGroupExpenses(groupId, limit || 50, offset || 0);
        return {
          expenses: (result.expenses || []).map(formatExpense),
          totalCount: result.totalCount || 0,
        };
      },

      groupBalances: async (_: any, { groupId }: { groupId: string }) => {
        const result = await expenseClient.getGroupBalances(groupId);
        return {
          balances: (result.balances || []).map((b: any) => ({
            userId: b.userId,
            balance: parseFloat(b.balance) || 0,
          })),
          details: (result.details || []).map((d: any) => ({
            fromUserId: d.fromUserId,
            toUserId: d.toUserId,
            amount: parseFloat(d.amount) || 0,
          })),
        };
      },

      settlementPlan: async (_: any, { groupId }: { groupId: string }) => {
        // First get balances, then compute optimal settlements
        const balanceResult = await expenseClient.getGroupBalances(groupId);
        const balances = (balanceResult.balances || []).map((b: any) => ({
          userId: b.userId,
          amount: parseFloat(b.balance) || 0,
        }));

        const result = await settlementClient.computeSettlements(groupId, balances);
        return {
          transactions: (result.transactions || []).map((t: any) => ({
            fromUserId: t.fromUserId,
            toUserId: t.toUserId,
            amount: parseFloat(t.amount) || 0,
          })),
        };
      },

      settlements: async (_: any, { groupId }: { groupId: string }) => {
        const result = await settlementClient.listSettlements(groupId);
        return (result.settlements || []).map(formatSettlement);
      },
    },

    // ─── Mutations ────────────────────────────────
    Mutation: {
      register: async (_: any, { email, name, password, currency }: any) => {
        const result = await userClient.register(email, name, password, currency || 'USD');
        return {
          user: formatUser(result.user),
          accessToken: result.accessToken,
          refreshToken: result.refreshToken,
        };
      },

      login: async (_: any, { email, password }: any) => {
        const result = await userClient.login(email, password);
        return {
          user: formatUser(result.user),
          accessToken: result.accessToken,
          refreshToken: result.refreshToken,
        };
      },

      refreshToken: async (_: any, { refreshToken }: { refreshToken: string }) => {
        const result = await userClient.refreshToken(refreshToken);
        return {
          accessToken: result.accessToken,
          refreshToken: result.refreshToken,
        };
      },

      updateProfile: async (_: any, { name, avatarUrl, currency }: any, context: { auth: AuthContext }) => {
        const userId = requireAuth(context.auth);
        const result = await userClient.updateUser(userId, name || '', avatarUrl || '', currency || '');
        return formatUser(result.user);
      },

      createGroup: async (_: any, { name, description, currency }: any, context: { auth: AuthContext }) => {
        const userId = requireAuth(context.auth);
        const result = await groupClient.createGroup(name, description || '', currency || 'USD', userId);
        return formatGroup(result.group);
      },

      updateGroup: async (_: any, { id, name, description, currency }: any) => {
        const result = await groupClient.updateGroup(id, name || '', description || '', currency || '');
        return formatGroup(result.group);
      },

      deleteGroup: async (_: any, { id }: { id: string }, context: { auth: AuthContext }) => {
        const userId = requireAuth(context.auth);
        const result = await groupClient.deleteGroup(id, userId);
        return result.success;
      },

      addMember: async (_: any, { groupId, userId, role }: any) => {
        const result = await groupClient.addMember(groupId, userId, role || 'member');
        return formatGroup(result.group);
      },

      removeMember: async (_: any, { groupId, userId }: any) => {
        const result = await groupClient.removeMember(groupId, userId);
        return formatGroup(result.group);
      },

      generateInviteCode: async (_: any, { groupId }: { groupId: string }, context: { auth: AuthContext }) => {
        const userId = requireAuth(context.auth);
        const result = await groupClient.generateInviteCode(groupId, userId);
        return result.inviteCode;
      },

      joinByInviteCode: async (_: any, { inviteCode }: { inviteCode: string }, context: { auth: AuthContext }) => {
        const userId = requireAuth(context.auth);
        const result = await groupClient.joinByInviteCode(inviteCode, userId);
        return formatGroup(result.group);
      },

      createExpense: async (_: any, { groupId, amount, description, splitType, splits }: any, context: { auth: AuthContext }) => {
        const userId = requireAuth(context.auth);
        const grpcSplitType = SPLIT_TYPE_MAP[splitType] || 1;
        const grpcSplits = splits.map((s: any) => ({
          userId: s.userId,
          amount: s.amount || 0,
          percentage: s.percentage || 0,
        }));
        const result = await expenseClient.createExpense(groupId, userId, amount, description, grpcSplitType, grpcSplits);
        return formatExpense(result.expense);
      },

      updateExpense: async (_: any, { id, amount, description, splitType, splits }: any, context: { auth: AuthContext }) => {
        const userId = requireAuth(context.auth);
        const grpcSplitType = SPLIT_TYPE_MAP[splitType] || 1;
        const grpcSplits = splits.map((s: any) => ({
          userId: s.userId,
          amount: s.amount || 0,
          percentage: s.percentage || 0,
        }));
        const result = await expenseClient.updateExpense(id, userId, amount, description, grpcSplitType, grpcSplits);
        return formatExpense(result.expense);
      },

      deleteExpense: async (_: any, { id }: { id: string }) => {
        const result = await expenseClient.deleteExpense(id);
        return result.success;
      },

      recordSettlement: async (_: any, { groupId, toUserId, amount }: any, context: { auth: AuthContext }) => {
        const userId = requireAuth(context.auth);
        const result = await settlementClient.recordSettlement(groupId, userId, toUserId, amount);
        return formatSettlement(result.settlement);
      },

      markSettlementCompleted: async (_: any, { id }: { id: string }) => {
        const result = await settlementClient.markSettlementCompleted(id);
        return formatSettlement(result.settlement);
      },
    },
  };
}
