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

export const resolveQuotaCurrencyConfig = ({
  type = 'USD',
  priceRate = 1,
  usdExchangeRate = 1,
  customExchangeRate = 1,
  customCurrencySymbol = '¤',
} = {}) => {
  const rechargeRate =
    Number.isFinite(Number(priceRate)) && Number(priceRate) > 0
      ? Number(priceRate)
      : 1;
  const usdRate =
    Number.isFinite(Number(usdExchangeRate)) && Number(usdExchangeRate) > 0
      ? Number(usdExchangeRate)
      : 1;
  const customRate =
    Number.isFinite(Number(customExchangeRate)) &&
    Number(customExchangeRate) > 0
      ? Number(customExchangeRate)
      : 1;

  if (type === 'CNY') {
    return { symbol: '¥', rate: rechargeRate, type };
  }
  if (type === 'CUSTOM') {
    return {
      symbol: customCurrencySymbol,
      rate: (rechargeRate / usdRate) * customRate,
      type,
    };
  }
  if (type === 'TOKENS') {
    return { symbol: '', rate: 1, type };
  }
  return { symbol: '$', rate: rechargeRate / usdRate, type: 'USD' };
};

export const quotaBaseToDisplayAmount = (amount, currencyConfig) =>
  Number(amount) * (currencyConfig?.rate || 1);

export const displayAmountToQuotaBase = (amount, currencyConfig) =>
  Number(amount) / (currencyConfig?.rate || 1);
