import { Request, Response } from 'express';
import jwt from 'jsonwebtoken';
import bcrypt from 'bcryptjs';

// Mock user data
const users = [
  {
    id: 1,
    username: 'admin',
    passwordHash: bcrypt.hashSync('password', 8),
    name: 'Admin User',
    email: 'admin@example.com',
    role: 1,
    status: 1,
  },
];

export const login = async (req: Request, res: Response) => {
  const { username, password } = req.body;

  if (!username || !password) {
    return res.status(400).send({ message: '用户名和密码是必填项' });
  }

  const user = users.find((u) => u.username === username);

  if (!user) {
    return res.status(404).send({ message: '用户不存在' });
  }

  const passwordIsValid = bcrypt.compareSync(password, user.passwordHash);

  if (!passwordIsValid) {
    return res.status(401).send({
      accessToken: null,
      message: '无效的密码',
    });
  }

  const token = jwt.sign({ id: user.id }, 'your-secret-key', {
    expiresIn: 86400, // 24 hours
  });

  res.status(200).send({
    token: token,
    user: {
      id: user.id,
      username: user.username,
      name: user.name,
      email: user.email,
      role: user.role,
      status: user.status,
    },
  });
};