export const typeDefs = `#graphql
  scalar DateTime

  # ─── User Types ──────────────────────────────────
  type User {
    id: ID!
    email: String!
    name: String!
    avatarUrl: String
    currency: String!
    createdAt: DateTime
    updatedAt: DateTime
  }

  type AuthPayload {
    user: User!
    accessToken: String!
    refreshToken: String!
  }

  type TokenPayload {
    accessToken: String!
    refreshToken: String!
  }

  # ─── Group Types ─────────────────────────────────
  type GroupMember {
    userId: String!
    role: String!
    joinedAt: DateTime
    user: User
  }

  type Group {
    id: ID!
    name: String!
    description: String
    currency: String!
    inviteCode: String
    ownerId: String!
    createdAt: DateTime
    updatedAt: DateTime
    members: [GroupMember!]!
  }

  # ─── Expense Types ──────────────────────────────
  enum SplitType {
    EQUAL
    CUSTOM
    PERCENTAGE
  }

  type ExpenseSplit {
    id: ID!
    expenseId: String!
    userId: String!
    amount: Float!
    percentage: Float
    user: User
  }

  type Expense {
    id: ID!
    groupId: String!
    payerId: String!
    amount: Float!
    description: String
    splitType: SplitType!
    splits: [ExpenseSplit!]!
    payer: User
    createdAt: DateTime
    updatedAt: DateTime
  }

  type ExpenseList {
    expenses: [Expense!]!
    totalCount: Int!
  }

  type UserBalance {
    userId: String!
    balance: Float!
    user: User
  }

  type BalanceDetail {
    fromUserId: String!
    toUserId: String!
    amount: Float!
    fromUser: User
    toUser: User
  }

  type GroupBalances {
    balances: [UserBalance!]!
    details: [BalanceDetail!]!
  }

  # ─── Settlement Types ───────────────────────────
  enum SettlementStatus {
    PENDING
    COMPLETED
  }

  type Transaction {
    fromUserId: String!
    toUserId: String!
    amount: Float!
    fromUser: User
    toUser: User
  }

  type Settlement {
    id: ID!
    groupId: String!
    fromUserId: String!
    toUserId: String!
    amount: Float!
    status: SettlementStatus!
    fromUser: User
    toUser: User
    createdAt: DateTime
    updatedAt: DateTime
  }

  type SettlementPlan {
    transactions: [Transaction!]!
  }

  # ─── Inputs ─────────────────────────────────────
  input ExpenseSplitInput {
    userId: String!
    amount: Float
    percentage: Float
  }

  # ─── Queries ────────────────────────────────────
  type Query {
    # User
    me: User
    user(id: ID!): User

    # Groups
    myGroups: [Group!]!
    group(id: ID!): Group

    # Expenses
    expense(id: ID!): Expense
    groupExpenses(groupId: ID!, limit: Int, offset: Int): ExpenseList!
    groupBalances(groupId: ID!): GroupBalances!

    # Settlements
    settlementPlan(groupId: ID!): SettlementPlan!
    settlements(groupId: ID!): [Settlement!]!
  }

  # ─── Mutations ──────────────────────────────────
  type Mutation {
    # Auth
    register(email: String!, name: String!, password: String!, currency: String): AuthPayload!
    login(email: String!, password: String!): AuthPayload!
    refreshToken(refreshToken: String!): TokenPayload!

    # User
    updateProfile(name: String, avatarUrl: String, currency: String): User!

    # Groups
    createGroup(name: String!, description: String, currency: String): Group!
    updateGroup(id: ID!, name: String, description: String, currency: String): Group!
    deleteGroup(id: ID!): Boolean!
    addMember(groupId: ID!, userId: ID!, role: String): Group!
    removeMember(groupId: ID!, userId: ID!): Group!
    generateInviteCode(groupId: ID!): String!
    joinByInviteCode(inviteCode: String!): Group!

    # Expenses
    createExpense(
      groupId: ID!
      amount: Float!
      description: String!
      splitType: SplitType!
      splits: [ExpenseSplitInput!]!
    ): Expense!
    updateExpense(
      id: ID!
      amount: Float!
      description: String!
      splitType: SplitType!
      splits: [ExpenseSplitInput!]!
    ): Expense!
    deleteExpense(id: ID!): Boolean!

    # Settlements
    recordSettlement(groupId: ID!, toUserId: ID!, amount: Float!): Settlement!
    markSettlementCompleted(id: ID!): Settlement!
  }
`;
