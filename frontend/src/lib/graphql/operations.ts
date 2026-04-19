import { gql } from '@apollo/client';

// ─── Auth ──────────────────────────────────────────
export const REGISTER = gql`
  mutation Register($email: String!, $name: String!, $password: String!, $currency: String) {
    register(email: $email, name: $name, password: $password, currency: $currency) {
      user { id email name avatarUrl currency }
      accessToken
      refreshToken
    }
  }
`;

export const LOGIN = gql`
  mutation Login($email: String!, $password: String!) {
    login(email: $email, password: $password) {
      user { id email name avatarUrl currency }
      accessToken
      refreshToken
    }
  }
`;

export const REFRESH_TOKEN = gql`
  mutation RefreshToken($refreshToken: String!) {
    refreshToken(refreshToken: $refreshToken) {
      accessToken
      refreshToken
    }
  }
`;

// ─── User ──────────────────────────────────────────
export const GET_ME = gql`
  query GetMe {
    me { id email name avatarUrl currency createdAt }
  }
`;

export const UPDATE_PROFILE = gql`
  mutation UpdateProfile($name: String, $avatarUrl: String, $currency: String) {
    updateProfile(name: $name, avatarUrl: $avatarUrl, currency: $currency) {
      id email name avatarUrl currency
    }
  }
`;

// ─── Groups ────────────────────────────────────────
export const GET_MY_GROUPS = gql`
  query GetMyGroups {
    myGroups {
      id name description currency inviteCode ownerId createdAt
      members { userId role joinedAt user { id name email avatarUrl } }
    }
  }
`;

export const GET_GROUP = gql`
  query GetGroup($id: ID!) {
    group(id: $id) {
      id name description currency inviteCode ownerId createdAt updatedAt
      members { userId role joinedAt user { id name email avatarUrl } }
    }
  }
`;

export const CREATE_GROUP = gql`
  mutation CreateGroup($name: String!, $description: String, $currency: String) {
    createGroup(name: $name, description: $description, currency: $currency) {
      id name description currency inviteCode ownerId
      members { userId role }
    }
  }
`;

export const UPDATE_GROUP = gql`
  mutation UpdateGroup($id: ID!, $name: String, $description: String, $currency: String) {
    updateGroup(id: $id, name: $name, description: $description, currency: $currency) {
      id name description currency
    }
  }
`;

export const DELETE_GROUP = gql`
  mutation DeleteGroup($id: ID!) {
    deleteGroup(id: $id)
  }
`;

export const JOIN_BY_INVITE_CODE = gql`
  mutation JoinByInviteCode($inviteCode: String!) {
    joinByInviteCode(inviteCode: $inviteCode) {
      id name description currency
    }
  }
`;

export const GENERATE_INVITE_CODE = gql`
  mutation GenerateInviteCode($groupId: ID!) {
    generateInviteCode(groupId: $groupId)
  }
`;

export const REMOVE_MEMBER = gql`
  mutation RemoveMember($groupId: ID!, $userId: ID!) {
    removeMember(groupId: $groupId, userId: $userId) {
      id members { userId role }
    }
  }
`;

// ─── Expenses ──────────────────────────────────────
export const GET_GROUP_EXPENSES = gql`
  query GetGroupExpenses($groupId: ID!, $limit: Int, $offset: Int) {
    groupExpenses(groupId: $groupId, limit: $limit, offset: $offset) {
      expenses {
        id groupId payerId amount description splitType createdAt
        payer { id name avatarUrl }
        splits { id userId amount percentage user { id name avatarUrl } }
      }
      totalCount
    }
  }
`;

export const CREATE_EXPENSE = gql`
  mutation CreateExpense(
    $groupId: ID!
    $amount: Float!
    $description: String!
    $splitType: SplitType!
    $splits: [ExpenseSplitInput!]!
  ) {
    createExpense(
      groupId: $groupId
      amount: $amount
      description: $description
      splitType: $splitType
      splits: $splits
    ) {
      id groupId payerId amount description splitType
      splits { id userId amount percentage }
    }
  }
`;

export const UPDATE_EXPENSE = gql`
  mutation UpdateExpense(
    $id: ID!
    $amount: Float!
    $description: String!
    $splitType: SplitType!
    $splits: [ExpenseSplitInput!]!
  ) {
    updateExpense(
      id: $id
      amount: $amount
      description: $description
      splitType: $splitType
      splits: $splits
    ) {
      id amount description splitType
      splits { id userId amount percentage }
    }
  }
`;

export const DELETE_EXPENSE = gql`
  mutation DeleteExpense($id: ID!) {
    deleteExpense(id: $id)
  }
`;

// ─── Balances & Settlements ────────────────────────
export const GET_GROUP_BALANCES = gql`
  query GetGroupBalances($groupId: ID!) {
    groupBalances(groupId: $groupId) {
      balances { userId balance user { id name avatarUrl } }
      details { fromUserId toUserId amount fromUser { id name } toUser { id name } }
    }
  }
`;

export const GET_SETTLEMENT_PLAN = gql`
  query GetSettlementPlan($groupId: ID!) {
    settlementPlan(groupId: $groupId) {
      transactions {
        fromUserId toUserId amount
        fromUser { id name avatarUrl }
        toUser { id name avatarUrl }
      }
    }
  }
`;

export const GET_SETTLEMENTS = gql`
  query GetSettlements($groupId: ID!) {
    settlements(groupId: $groupId) {
      id groupId fromUserId toUserId amount status createdAt
      fromUser { id name avatarUrl }
      toUser { id name avatarUrl }
    }
  }
`;

export const RECORD_SETTLEMENT = gql`
  mutation RecordSettlement($groupId: ID!, $toUserId: ID!, $amount: Float!) {
    recordSettlement(groupId: $groupId, toUserId: $toUserId, amount: $amount) {
      id fromUserId toUserId amount status
    }
  }
`;

export const MARK_SETTLEMENT_COMPLETED = gql`
  mutation MarkSettlementCompleted($id: ID!) {
    markSettlementCompleted(id: $id) {
      id status
    }
  }
`;
