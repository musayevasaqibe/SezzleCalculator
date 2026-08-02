import React from 'react';
import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import Calculator from '../Calculator';
import { vi } from 'vitest';

describe('Calculator UI', () => {
  beforeEach(() => {
    // mock fetch for backend calls
    vi.stubGlobal('fetch', vi.fn(async (url, opts) => {
      // basic mock: if op add and a/b provided return sum
      const body = JSON.parse((opts as any).body);
      const op = body.op;
      if (op === 'add') {
        return {
          json: async () => ({ result: body.a + (body.b ?? 0), error: null })
        };
      }
      return { json: async () => ({ result: 0, error: null }) };
    }));
  });

  afterEach(() => {
    vi.restoreAllMocks();
  });

  test('digits appear on display when clicked and form numbers', async () => {
    render(<Calculator />);
    const user = userEvent.setup();

    // initial display shows 0 (per component)
    const display = screen.getByTestId('display');
    expect(display.textContent).toBe('0' || '');

    // click 1 -> should show '1'
    await user.click(screen.getByRole('button', { name: '1' }));
    expect(display.textContent).toContain('1');

    // click 1 again -> should show '11'
    await user.click(screen.getByRole('button', { name: '1' }));
    expect(display.textContent).toContain('11');

    // click 2 -> should append -> '112'
    await user.click(screen.getByRole('button', { name: '2' }));
    expect(display.textContent).toContain('112');
  });

  test('performs addition via backend when equals pressed', async () => {
    render(<Calculator />);
    const user = userEvent.setup();

    // enter 4
    await user.click(screen.getByRole('button', { name: '4' }));
    // choose +
    await user.click(screen.getByRole('button', { name: '+' }));
    // enter 5
    await user.click(screen.getByRole('button', { name: '5' }));
    // press =
    await user.click(screen.getByRole('button', { name: '=' }));

    // wait a bit for async fetch to update display
    const display = await screen.findByTestId('display');
    expect(display.textContent).toContain('9');
  });
});
