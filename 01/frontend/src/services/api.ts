import axios, { isAxiosError } from 'axios';
import type { CreateUserInput, User } from '../types/user';

const baseURL = import.meta.env.VITE_API_BASE ?? '/api';

export const api = axios.create({
  baseURL,
  headers: { 'Content-Type': 'application/json' },
  timeout: 15000,
});

export async function listUsers(): Promise<User[]> {
  const { data } = await api.get<User[] | null>('/users');
  // Backend may send JSON null for empty lists if a nil Go slice is serialized; normalize to [].
  return Array.isArray(data) ? data : [];
}

export async function createUser(body: CreateUserInput): Promise<User> {
  const { data } = await api.post<User>('/users', body);
  return data;
}

export async function updateUser(id: string, body: Partial<CreateUserInput>): Promise<User> {
  const { data } = await api.put<User>(`/users/${id}`, body);
  return data;
}

export async function deleteUser(id: string): Promise<void> {
  await api.delete(`/users/${id}`);
}

export { isAxiosError };
