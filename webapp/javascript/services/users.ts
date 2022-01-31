import { Result } from '@utils/fp';
import { Users, parse, type User } from '@models/users';
import type { ZodError } from 'zod';
import { request } from './base';
import type { RequestError } from './base';

export interface FetchUsersError {
  message?: string;
}

export async function fetchUsers(): Promise<
  Result<Users, RequestError | ZodError>
> {
  const response = await request('/api/users');

  if (response.isOk) {
    return parse(response.value);
  }

  return Result.err<Users, RequestError>(response.error);
}

export async function disableUser(
  user: User
): Promise<Result<boolean, RequestError | ZodError>> {
  const response = await request(`/api/users/${user.id}/disable`);

  if (response.isOk) {
    return Result.ok(true);
  }

  return Result.err<false, RequestError>(response.error);
}

export async function enableUser(
  user: User
): Promise<Result<boolean, RequestError | ZodError>> {
  const response = await request(`/api/users/${user.id}/enable`);

  if (response.isOk) {
    return Result.ok(true);
  }

  return Result.err<false, RequestError>(response.error);
}

export async function createUser(
  user: User
): Promise<Result<boolean, RequestError | ZodError>> {
  const response = await request(`/api/users`, {
    method: 'POST',
    body: JSON.stringify(user),
  });

  if (response.isOk) {
    return Result.ok(true);
  }

  return Result.err<false, RequestError>(response.error);
}