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

import React from 'react';
import { Button, Card, Input } from '@douyinfe/semi-ui';
import {
  ArrowRightLeft,
  BarChart2,
  ClipboardList,
  Copy,
  Gift,
  Globe2,
  Link,
  ShieldCheck,
  Sparkles,
  TrendingUp,
  Users,
  Wallet,
  Zap,
} from 'lucide-react';

const TrustItem = ({ icon, title, description }) => (
  <div
    className='rounded-[10px] p-3 backdrop-blur'
    style={{
      background: 'rgba(255, 255, 255, 0.1)',
      border: '1px solid rgba(255, 255, 255, 0.15)',
      color: '#ffffff',
    }}
  >
    <div className='mb-1.5' style={{ color: '#ffffff' }}>
      {icon}
    </div>
    <div className='mb-1 text-[13px] font-medium' style={{ color: '#ffffff' }}>
      {title}
    </div>
    <div
      className='text-xs leading-5'
      style={{ color: 'rgba(255, 255, 255, 0.82)' }}
    >
      {description}
    </div>
  </div>
);

const StatCard = ({ label, value, tone = 'default', icon }) => {
  const toneClass = {
    default: 'text-gray-900 dark:text-gray-100',
    green: 'text-[#00b42a]',
    orange: 'text-[#ff7d00]',
    blue: 'text-[#0064fa]',
  }[tone];

  return (
    <div className='rounded-xl border border-[#f0f1f5] bg-white p-4 dark:border-zinc-700 dark:bg-zinc-900'>
      <div className='mb-1.5 flex items-center gap-1.5 text-xs text-[#86909c]'>
        {icon}
        <span>{label}</span>
      </div>
      <div className={`text-[22px] font-semibold leading-tight ${toneClass}`}>
        {value}
      </div>
    </div>
  );
};

const RuleLine = ({ title, children }) => (
  <div className='mb-2 text-[13px] leading-7 text-[#4e5969] dark:text-zinc-300'>
    <strong className='text-[#1c1f23] dark:text-zinc-100'>{title}</strong>
    <span className='mx-1'>·</span>
    {children}
  </div>
);

const formatCny = (cents) => {
  const value = Number(cents || 0) / 100;
  return `¥${value.toFixed(2)}`;
};

const AffiliatePromotionBlock = ({
  t,
  dashboard,
  affLink,
  handleAffLinkClick,
  handleAffCodeClick,
}) => {
  const affCode = affLink ? affLink.split('aff=').pop() : '';
  const availableCommission = dashboard?.commission_available || 0;
  const pendingCommission = dashboard?.commission_pending || 0;
  const totalCommission = dashboard?.commission_total || 0;
  const inviteCount = dashboard?.invite_count || 0;

  return (
    <div className='mt-6'>
      <div className='mb-3'>
        <div className='text-xl font-semibold text-[#1c1f23] dark:text-zinc-100'>
          {t('合作推广')}
        </div>
        <div className='mt-1 text-sm text-[#86909c]'>
          {t('邀请好友注册，好友充值后您可获得相应奖励')}
        </div>
      </div>

      <div
        className='rounded-2xl p-6'
        style={{
          background: 'linear-gradient(135deg, #1d4ed8 0%, #3730a3 100%)',
          color: '#ffffff',
        }}
      >
        <div
          className='inline-flex items-center gap-1.5 rounded-full px-3 py-1 text-xs font-medium'
          style={{
            background: 'rgba(255, 255, 255, 0.1)',
            border: '1px solid rgba(255, 255, 255, 0.3)',
            color: '#ffffff',
          }}
        >
          <Sparkles size={13} />
          {t('推广政策')}
        </div>
        <div
          className='mt-3 text-[32px] font-semibold leading-tight'
          style={{ color: '#ffffff' }}
        >
          {t('全部让利给推广人')}
        </div>
        <div
          className='mt-2 text-[15px] leading-7'
          style={{ color: 'rgba(255, 255, 255, 0.92)' }}
        >
          {t(
            '你邀请的每一笔充值，10% 直接进你账户 · 终身有效 · 21 天观察期后即可提现',
          )}
        </div>

        <div className='mt-4 grid grid-cols-1 gap-3 sm:grid-cols-2 lg:grid-cols-4'>
          <TrustItem
            icon={<Globe2 size={18} />}
            title={t('源头模型筛选')}
            description={t('全球化 AI 公司，模型自用同款')}
          />
          <TrustItem
            icon={<TrendingUp size={18} />}
            title={t('价格透明')}
            description={t('在源头模型基础上加价 10–12%')}
          />
          <TrustItem
            icon={<Zap size={18} />}
            title={t('纯消耗模式')}
            description={t('不搞会员制，充多少用多少')}
          />
          <TrustItem
            icon={<ShieldCheck size={18} />}
            title={t('服务稳定')}
            description={t('多通道冗余，24/7 监控')}
          />
        </div>
      </div>

      <div className='mt-4 grid grid-cols-1 gap-4 lg:grid-cols-2'>
        <Card className='!rounded-2xl border-0 shadow-sm'>
          <div className='mb-4 flex items-center'>
            <div className='mr-3 flex h-8 w-8 items-center justify-center rounded-lg bg-[#0064fa] text-white'>
              <Link size={16} />
            </div>
            <div>
              <div className='text-base font-medium text-[#1c1f23] dark:text-zinc-100'>
                {t('我的邀请')}
              </div>
              <div className='text-xs text-[#86909c]'>
                {t('复制链接或邀请码分享给朋友')}
              </div>
            </div>
          </div>

          <div className='space-y-3'>
            <Input
              value={affLink}
              readOnly
              className='!rounded-[10px]'
              prefix={
                <span className='text-xs text-[#86909c]'>{t('邀请链接')}</span>
              }
              suffix={
                <Button
                  type='primary'
                  theme='solid'
                  className='!rounded-md'
                  icon={<Copy size={14} />}
                  onClick={handleAffLinkClick}
                >
                  {t('复制')}
                </Button>
              }
            />
            <Input
              value={affCode}
              readOnly
              className='!rounded-[10px]'
              prefix={
                <span className='text-xs text-[#86909c]'>{t('邀请码')}</span>
              }
              suffix={
                <Button
                  type='primary'
                  theme='solid'
                  className='!rounded-md'
                  icon={<Copy size={14} />}
                  disabled={!affCode}
                  onClick={handleAffCodeClick}
                >
                  {t('复制')}
                </Button>
              }
            />
          </div>

          <div className='mt-4 flex items-start gap-3 rounded-[10px] border border-dashed border-[#ffd591] bg-[#fff7e6] p-3 text-xs leading-5 text-[#d4691f]'>
            <Gift size={18} className='mt-0.5 shrink-0' />
            <div>
              <strong>{t('新人首充福利')}</strong>
              {t(' · 朋友通过邀请链接注册后，后续每笔在线充值都会计入您的分佣。')}
            </div>
          </div>

          <div className='mt-4 grid grid-cols-1 gap-3 sm:grid-cols-2'>
            <Button
              type='primary'
              theme='solid'
              disabled
              className='!h-11 !rounded-[10px]'
              icon={<ArrowRightLeft size={16} />}
            >
              {t('待结算后可提现')}
            </Button>
            <Button
              theme='light'
              type='tertiary'
              disabled
              className='!h-11 !rounded-[10px]'
              icon={<Wallet size={16} />}
            >
              {t('申请提现')}
            </Button>
          </div>
        </Card>

        <div>
          <div className='mb-4 grid grid-cols-2 gap-3'>
            <StatCard
              icon={<Users size={14} />}
              label={t('邀请人数')}
              value={inviteCount}
              tone='blue'
            />
            <StatCard
              icon={<BarChart2 size={14} />}
              label={t('累计收益')}
              value={formatCny(totalCommission)}
            />
            <StatCard
              icon={<ClipboardList size={14} />}
              label={t('待结算')}
              value={formatCny(pendingCommission)}
              tone='orange'
            />
            <StatCard
              icon={<Wallet size={14} />}
              label={t('可提现')}
              value={formatCny(availableCommission)}
              tone='green'
            />
          </div>

          <Card className='!rounded-2xl border-0 shadow-sm'>
            <div className='mb-3 flex items-center'>
              <div className='mr-3 flex h-8 w-8 items-center justify-center rounded-lg bg-[#ff7d00] text-white'>
                <ClipboardList size={16} />
              </div>
              <div>
                <div className='text-base font-medium text-[#1c1f23] dark:text-zinc-100'>
                {t('推广规则说明')}
              </div>
              <div className='text-xs text-[#86909c]'>
                  {t('当前展示已接入的充值分佣数据')}
              </div>
            </div>
          </div>

          <div className='rounded-[10px] border-l-[3px] border-[#ff7d00] bg-gradient-to-r from-[#fff7e6] to-white p-3 text-[13px] font-medium text-[#d4691f] dark:to-zinc-900'>
              {t(
                '核心政策：邀请用户每笔在线充值按人民币支付金额的 10% 生成分佣，当前默认进入待结算状态。',
              )}
          </div>

            <div className='mt-4'>
              <RuleLine title={t('1. 邀请关系')}>
                {t('朋友通过邀请链接或邀请码注册后，会计入您的邀请人数。')}
              </RuleLine>
              <RuleLine title={t('2. 分佣比例')}>
                {t('邀请用户后续每笔在线充值，按实际支付金额的 10% 计算收益。')}
              </RuleLine>
              <RuleLine title={t('3. 结算状态')}>
                {t(
                  '分佣先进入待结算状态，21 天观察期后可继续接入提现或划转流程。',
                )}
              </RuleLine>
              <RuleLine title={t('4. 风控')}>
                {t('严禁自邀请、刷单、虚假宣传，违规账号可冻结未结算收益。')}
              </RuleLine>
            </div>
          </Card>
        </div>
      </div>
    </div>
  );
};

export default AffiliatePromotionBlock;
