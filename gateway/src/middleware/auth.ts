import jwt from 'jsonwebtoken';
import { Request } from 'express';

export interface AuthContext {
  userId: string | null;
}

export function getAuthContext(req: Request, jwtSecret: string): AuthContext {
  const authHeader = req.headers.authorization;

  if (!authHeader || !authHeader.startsWith('Bearer ')) {
    return { userId: null };
  }

  const token = authHeader.substring(7);

  try {
    const decoded = jwt.verify(token, jwtSecret) as any;
    return { userId: decoded.sub || null };
  } catch (error) {
    return { userId: null };
  }
}

export function requireAuth(context: AuthContext): string {
  if (!context.userId) {
    throw new Error('Authentication required');
  }
  return context.userId;
}
