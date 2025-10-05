import { Request, Response, NextFunction } from 'express';
import jwt from 'jsonwebtoken';

interface CustomRequest extends Request {
  userId?: number;
}

export const verifyToken = (req: CustomRequest, res: Response, next: NextFunction) => {
  const token = req.headers['authorization']?.split(' ')[1];

  if (!token) {
    return res.status(403).send({ message: '没有提供令牌' });
  }

  jwt.verify(token, 'your-secret-key', (err, decoded) => {
    if (err) {
      return res.status(401).send({ message: '未授权' });
    }
    if (decoded && typeof decoded !== 'string') {
      req.userId = decoded.id;
    }
    next();
  });
};