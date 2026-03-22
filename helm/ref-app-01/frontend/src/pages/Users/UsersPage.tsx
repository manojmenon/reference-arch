import { FormEvent, useCallback, useEffect, useState } from 'react';
import type { User } from '../../types/user';
import * as usersApi from '../../services/api';
import { isAxiosError } from '../../services/api';
import styles from './UsersPage.module.css';

type FormState = { name: string; email: string };
const emptyForm: FormState = { name: '', email: '' };

function validateEmail(email: string): string | null {
  if (!email.trim()) return 'Email is required';
  if (!/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(email)) return 'Enter a valid email';
  return null;
}

export function UsersPage() {
  const [users, setUsers] = useState<User[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [form, setForm] = useState<FormState>(emptyForm);
  const [formErrors, setFormErrors] = useState<Partial<FormState>>({});
  const [editingId, setEditingId] = useState<string | null>(null);
  const [busy, setBusy] = useState(false);

  const load = useCallback(async () => {
    setLoading(true);
    setError(null);
    try {
      const data = await usersApi.listUsers();
      setUsers(data);
    } catch (e) {
      setError(isAxiosError(e) ? e.message : 'Failed to load users');
    } finally {
      setLoading(false);
    }
  }, []);

  useEffect(() => {
    void load();
  }, [load]);

  function validateForm(): boolean {
    const next: Partial<FormState> = {};
    if (!form.name.trim()) next.name = 'Name is required';
    const em = validateEmail(form.email);
    if (em) next.email = em;
    setFormErrors(next);
    return Object.keys(next).length === 0;
  }

  async function onSubmit(e: FormEvent) {
    e.preventDefault();
    if (!validateForm()) return;
    setBusy(true);
    setError(null);
    try {
      if (editingId) {
        const updated = await usersApi.updateUser(editingId, {
          name: form.name.trim(),
          email: form.email.trim(),
        });
        setUsers((prev) => prev.map((u) => (u.id === updated.id ? updated : u)));
        setEditingId(null);
      } else {
        const created = await usersApi.createUser({
          name: form.name.trim(),
          email: form.email.trim(),
        });
        setUsers((prev) => [created, ...prev]);
      }
      setForm(emptyForm);
      setFormErrors({});
    } catch (e) {
      const msg = isAxiosError(e) ? (e.response?.data as { error?: string })?.error ?? e.message : 'Request failed';
      setError(String(msg));
    } finally {
      setBusy(false);
    }
  }

  function startEdit(u: User) {
    setEditingId(u.id);
    setForm({ name: u.name, email: u.email });
    setFormErrors({});
  }

  async function onDelete(id: string) {
    if (!confirm('Delete this user?')) return;
    setBusy(true);
    setError(null);
    try {
      await usersApi.deleteUser(id);
      setUsers((prev) => prev.filter((u) => u.id !== id));
      if (editingId === id) {
        setEditingId(null);
        setForm(emptyForm);
      }
    } catch (e) {
      setError(isAxiosError(e) ? e.message : 'Delete failed');
    } finally {
      setBusy(false);
    }
  }

  return (
    <div className={styles.layout}>
      <header className={styles.header}>
        <h1>Directory</h1>
        <p className={styles.lede}>Manage users with validation and clear feedback.</p>
      </header>

      <section className={styles.card}>
        <h2 className={styles.h2}>{editingId ? 'Edit user' : 'Create user'}</h2>
        <form className={styles.form} onSubmit={onSubmit}>
          <label className={styles.field}>
            <span>Name</span>
            <input
              value={form.name}
              onChange={(e) => setForm((f) => ({ ...f, name: e.target.value }))}
              disabled={busy}
              autoComplete="name"
            />
            {formErrors.name && <span className={styles.fieldError}>{formErrors.name}</span>}
          </label>
          <label className={styles.field}>
            <span>Email</span>
            <input
              type="email"
              value={form.email}
              onChange={(e) => setForm((f) => ({ ...f, email: e.target.value }))}
              disabled={busy}
              autoComplete="email"
            />
            {formErrors.email && <span className={styles.fieldError}>{formErrors.email}</span>}
          </label>
          <div className={styles.actions}>
            <button type="submit" className={styles.primary} disabled={busy}>
              {editingId ? 'Save changes' : 'Add user'}
            </button>
            {editingId && (
              <button
                type="button"
                className={styles.ghost}
                disabled={busy}
                onClick={() => {
                  setEditingId(null);
                  setForm(emptyForm);
                  setFormErrors({});
                }}
              >
                Cancel
              </button>
            )}
          </div>
        </form>
        {error && <p className={styles.banner}>{error}</p>}
      </section>

      <section className={styles.card}>
        <div className={styles.toolbar}>
          <h2 className={styles.h2}>People</h2>
          <button type="button" className={styles.ghost} onClick={() => void load()} disabled={loading}>
            Refresh
          </button>
        </div>
        {loading ? (
          <p className={styles.muted}>Loading…</p>
        ) : users.length === 0 ? (
          <p className={styles.muted}>No users yet.</p>
        ) : (
          <ul className={styles.list}>
            {users.map((u) => (
              <li key={u.id} className={styles.row}>
                <div>
                  <div className={styles.name}>{u.name}</div>
                  <div className={styles.email}>{u.email}</div>
                </div>
                <div className={styles.rowActions}>
                  <button type="button" className={styles.link} onClick={() => startEdit(u)} disabled={busy}>
                    Edit
                  </button>
                  <button
                    type="button"
                    className={styles.danger}
                    onClick={() => void onDelete(u.id)}
                    disabled={busy}
                  >
                    Delete
                  </button>
                </div>
              </li>
            ))}
          </ul>
        )}
      </section>
    </div>
  );
}
