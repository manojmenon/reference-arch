import { render, screen, waitFor } from '@testing-library/react';
import { vi } from 'vitest';
import App from './App';

vi.mock('./services/api', () => ({
  listUsers: vi.fn().mockResolvedValue([]),
  createUser: vi.fn(),
  updateUser: vi.fn(),
  deleteUser: vi.fn(),
  isAxiosError: vi.fn().mockReturnValue(false),
}));

describe('App', () => {
  it('renders directory heading', async () => {
    render(<App />);
    expect(screen.getByRole('heading', { name: /directory/i })).toBeInTheDocument();
    await waitFor(() => expect(screen.queryByText(/loading/i)).not.toBeInTheDocument());
  });
});
