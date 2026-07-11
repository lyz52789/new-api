/*
Copyright (C) 2025 QuantumNous

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU Affero General Public License as
published by the Free Software Foundation, either version 3 of the
License, or (at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
GNU Affero General Public License for more details.

You should have received a copy of the GNU Affero General Public License
along with this program. If not, see <https://www.gnu.org/licenses/>.

For commercial licensing, please contact support@quantumnous.com
*/

import { describe, expect, test } from 'bun:test';
import {
  displayAmountToQuotaBase,
  quotaBaseToDisplayAmount,
  resolveQuotaCurrencyConfig,
} from './quotaCurrency';

const token123Rates = {
  priceRate: 1,
  usdExchangeRate: 7,
  customExchangeRate: 0.91,
  customCurrencySymbol: '€',
};

describe('resolveQuotaCurrencyConfig', () => {
  test('keeps the recharge-priced quota amount unchanged in CNY', () => {
    expect(
      resolveQuotaCurrencyConfig({ type: 'CNY', ...token123Rates }),
    ).toEqual({ symbol: '¥', rate: 1, type: 'CNY' });
  });

  test('converts the recharge-priced quota amount from CNY to USD', () => {
    expect(
      resolveQuotaCurrencyConfig({ type: 'USD', ...token123Rates }),
    ).toEqual({ symbol: '$', rate: 1 / 7, type: 'USD' });
  });

  test('converts through USD for a custom currency', () => {
    expect(
      resolveQuotaCurrencyConfig({ type: 'CUSTOM', ...token123Rates }),
    ).toEqual({ symbol: '€', rate: 0.91 / 7, type: 'CUSTOM' });
  });

  test('leaves token quota untouched', () => {
    expect(
      resolveQuotaCurrencyConfig({ type: 'TOKENS', ...token123Rates }),
    ).toEqual({ symbol: '', rate: 1, type: 'TOKENS' });
  });

  test('converts the same wallet amount consistently in CNY and USD', () => {
    const cny = resolveQuotaCurrencyConfig({ type: 'CNY', ...token123Rates });
    const usd = resolveQuotaCurrencyConfig({ type: 'USD', ...token123Rates });

    expect(quotaBaseToDisplayAmount(7, cny)).toBe(7);
    expect(quotaBaseToDisplayAmount(7, usd)).toBe(1);
    expect(displayAmountToQuotaBase(1, usd)).toBe(7);
  });
});
