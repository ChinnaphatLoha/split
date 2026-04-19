import * as grpc from '@grpc/grpc-js';
import * as protoLoader from '@grpc/proto-loader';
import path from 'path';

function loadProtoClient(protoPath: string, packageName: string, serviceName: string, serviceUrl: string): any {
  const absolutePath = path.resolve(__dirname, '../../', protoPath);
  const packageDefinition = protoLoader.loadSync(absolutePath, {
    keepCase: false,
    longs: String,
    enums: String,
    defaults: true,
    oneofs: true,
  });
  const proto = grpc.loadPackageDefinition(packageDefinition) as any;
  const ServiceClient = proto[packageName][serviceName];
  return new ServiceClient(serviceUrl, grpc.credentials.createInsecure());
}

export class UserClient {
  private client: any;

  constructor(serviceUrl: string) {
    this.client = loadProtoClient('proto/user/user.proto', 'user', 'UserService', serviceUrl);
  }

  register(email: string, name: string, password: string, currency: string): Promise<any> {
    return new Promise((resolve, reject) => {
      this.client.register({ email, name, password, currency }, (err: any, res: any) => {
        if (err) reject(err);
        else resolve(res);
      });
    });
  }

  login(email: string, password: string): Promise<any> {
    return new Promise((resolve, reject) => {
      this.client.login({ email, password }, (err: any, res: any) => {
        if (err) reject(err);
        else resolve(res);
      });
    });
  }

  getUser(id: string): Promise<any> {
    return new Promise((resolve, reject) => {
      this.client.getUser({ id }, (err: any, res: any) => {
        if (err) reject(err);
        else resolve(res);
      });
    });
  }

  updateUser(id: string, name: string, avatarUrl: string, currency: string): Promise<any> {
    return new Promise((resolve, reject) => {
      this.client.updateUser({ id, name, avatarUrl, currency }, (err: any, res: any) => {
        if (err) reject(err);
        else resolve(res);
      });
    });
  }

  getUsersByIds(ids: string[]): Promise<any> {
    return new Promise((resolve, reject) => {
      this.client.getUsersByIds({ ids }, (err: any, res: any) => {
        if (err) reject(err);
        else resolve(res);
      });
    });
  }

  refreshToken(refreshToken: string): Promise<any> {
    return new Promise((resolve, reject) => {
      this.client.refreshToken({ refreshToken }, (err: any, res: any) => {
        if (err) reject(err);
        else resolve(res);
      });
    });
  }
}

export class GroupClient {
  private client: any;

  constructor(serviceUrl: string) {
    this.client = loadProtoClient('proto/group/group.proto', 'group', 'GroupService', serviceUrl);
  }

  createGroup(name: string, description: string, currency: string, ownerId: string): Promise<any> {
    return new Promise((resolve, reject) => {
      this.client.createGroup({ name, description, currency, ownerId }, (err: any, res: any) => {
        if (err) reject(err);
        else resolve(res);
      });
    });
  }

  getGroup(id: string): Promise<any> {
    return new Promise((resolve, reject) => {
      this.client.getGroup({ id }, (err: any, res: any) => {
        if (err) reject(err);
        else resolve(res);
      });
    });
  }

  listUserGroups(userId: string): Promise<any> {
    return new Promise((resolve, reject) => {
      this.client.listUserGroups({ userId }, (err: any, res: any) => {
        if (err) reject(err);
        else resolve(res);
      });
    });
  }

  updateGroup(id: string, name: string, description: string, currency: string): Promise<any> {
    return new Promise((resolve, reject) => {
      this.client.updateGroup({ id, name, description, currency }, (err: any, res: any) => {
        if (err) reject(err);
        else resolve(res);
      });
    });
  }

  deleteGroup(id: string, userId: string): Promise<any> {
    return new Promise((resolve, reject) => {
      this.client.deleteGroup({ id, userId }, (err: any, res: any) => {
        if (err) reject(err);
        else resolve(res);
      });
    });
  }

  addMember(groupId: string, userId: string, role: string): Promise<any> {
    return new Promise((resolve, reject) => {
      this.client.addMember({ groupId, userId, role }, (err: any, res: any) => {
        if (err) reject(err);
        else resolve(res);
      });
    });
  }

  removeMember(groupId: string, userId: string): Promise<any> {
    return new Promise((resolve, reject) => {
      this.client.removeMember({ groupId, userId }, (err: any, res: any) => {
        if (err) reject(err);
        else resolve(res);
      });
    });
  }

  generateInviteCode(groupId: string, userId: string): Promise<any> {
    return new Promise((resolve, reject) => {
      this.client.generateInviteCode({ groupId, userId }, (err: any, res: any) => {
        if (err) reject(err);
        else resolve(res);
      });
    });
  }

  joinByInviteCode(inviteCode: string, userId: string): Promise<any> {
    return new Promise((resolve, reject) => {
      this.client.joinByInviteCode({ inviteCode, userId }, (err: any, res: any) => {
        if (err) reject(err);
        else resolve(res);
      });
    });
  }
}

export class ExpenseClient {
  private client: any;

  constructor(serviceUrl: string) {
    this.client = loadProtoClient('proto/expense/expense.proto', 'expense', 'ExpenseService', serviceUrl);
  }

  createExpense(groupId: string, payerId: string, amount: number, description: string, splitType: number, splits: any[]): Promise<any> {
    return new Promise((resolve, reject) => {
      this.client.createExpense({ groupId, payerId, amount, description, splitType, splits }, (err: any, res: any) => {
        if (err) reject(err);
        else resolve(res);
      });
    });
  }

  getExpense(id: string): Promise<any> {
    return new Promise((resolve, reject) => {
      this.client.getExpense({ id }, (err: any, res: any) => {
        if (err) reject(err);
        else resolve(res);
      });
    });
  }

  listGroupExpenses(groupId: string, limit: number, offset: number): Promise<any> {
    return new Promise((resolve, reject) => {
      this.client.listGroupExpenses({ groupId, limit, offset }, (err: any, res: any) => {
        if (err) reject(err);
        else resolve(res);
      });
    });
  }

  updateExpense(id: string, payerId: string, amount: number, description: string, splitType: number, splits: any[]): Promise<any> {
    return new Promise((resolve, reject) => {
      this.client.updateExpense({ id, payerId, amount, description, splitType, splits }, (err: any, res: any) => {
        if (err) reject(err);
        else resolve(res);
      });
    });
  }

  deleteExpense(id: string): Promise<any> {
    return new Promise((resolve, reject) => {
      this.client.deleteExpense({ id }, (err: any, res: any) => {
        if (err) reject(err);
        else resolve(res);
      });
    });
  }

  getGroupBalances(groupId: string): Promise<any> {
    return new Promise((resolve, reject) => {
      this.client.getGroupBalances({ groupId }, (err: any, res: any) => {
        if (err) reject(err);
        else resolve(res);
      });
    });
  }
}

export class SettlementClient {
  private client: any;

  constructor(serviceUrl: string) {
    this.client = loadProtoClient('proto/settlement/settlement.proto', 'settlement', 'SettlementService', serviceUrl);
  }

  computeSettlements(groupId: string, balances: any[]): Promise<any> {
    return new Promise((resolve, reject) => {
      this.client.computeSettlements({ groupId, balances }, (err: any, res: any) => {
        if (err) reject(err);
        else resolve(res);
      });
    });
  }

  recordSettlement(groupId: string, fromUserId: string, toUserId: string, amount: number): Promise<any> {
    return new Promise((resolve, reject) => {
      this.client.recordSettlement({ groupId, fromUserId, toUserId, amount }, (err: any, res: any) => {
        if (err) reject(err);
        else resolve(res);
      });
    });
  }

  markSettlementCompleted(id: string): Promise<any> {
    return new Promise((resolve, reject) => {
      this.client.markSettlementCompleted({ id }, (err: any, res: any) => {
        if (err) reject(err);
        else resolve(res);
      });
    });
  }

  listSettlements(groupId: string): Promise<any> {
    return new Promise((resolve, reject) => {
      this.client.listSettlements({ groupId }, (err: any, res: any) => {
        if (err) reject(err);
        else resolve(res);
      });
    });
  }
}
