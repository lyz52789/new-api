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
import { Card, Avatar, Typography, Table, Tag } from '@douyinfe/semi-ui';
import { IconCoinMoneyStroked } from '@douyinfe/semi-icons';
import {
  calculateModelPrice,
  getModelPriceItems,
} from '../../../../../helpers';

const { Text } = Typography;

const ModelPricingTable = ({
  modelData,
  selectedGroup = 'all',
  groupRatio,
  currency,
  siteDisplayType,
  tokenUnit,
  displayPrice,
  showRatio,
  usableGroup,
  autoGroups = [],
  t,
}) => {
  const modelEnableGroups = Array.isArray(modelData?.enable_groups)
    ? modelData.enable_groups
    : [];
  const autoChain = autoGroups.filter((g) => modelEnableGroups.includes(g));
  const getAvailableGroups = () =>
    Object.keys(usableGroup || {})
      .filter((g) => g !== '')
      .filter((g) => g !== 'auto')
      .filter((g) => modelEnableGroups.includes(g));

  const resolveDisplayGroup = () => {
    const availableGroups = getAvailableGroups();
    if (
      selectedGroup !== 'all' &&
      availableGroups.includes(selectedGroup) &&
      groupRatio?.[selectedGroup] !== undefined
    ) {
      return {
        group: selectedGroup,
        ratio: groupRatio[selectedGroup],
      };
    }

    let minGroup = availableGroups[0] || '';
    let minRatio = groupRatio?.[minGroup] ?? 1;
    availableGroups.forEach((group) => {
      const ratio = groupRatio?.[group] ?? 1;
      if (ratio < minRatio) {
        minGroup = group;
        minRatio = ratio;
      }
    });
    return {
      group: minGroup,
      ratio: minRatio || 1,
    };
  };

  const formatRMB = (value, multiplier = 1) => {
    if (value === undefined || value === null || value === '') return '-';
    const n = Number(value) * Number(multiplier || 1);
    if (!Number.isFinite(n)) return '-';
    return `¥${n.toFixed(4).replace(/\.?0+$/, '')}`;
  };

  const renderVideoSalePrice = (row, ratio) => {
    const parts = [];
    if (row.sale_rmb_per_m_tokens) {
      parts.push(`${formatRMB(row.sale_rmb_per_m_tokens, ratio)} / 1M tokens`);
    }
    if (row.sale_rmb_per_video) {
      parts.push(`${formatRMB(row.sale_rmb_per_video, ratio)} / 条`);
    }
    if (row.sale_rmb_per_video_min || row.sale_rmb_per_video_max) {
      parts.push(
        `${formatRMB(row.sale_rmb_per_video_min, ratio)} - ${formatRMB(
          row.sale_rmb_per_video_max,
          ratio,
        )} / 条`,
      );
    }
    if (row.sale_rmb_per_second) {
      parts.push(`${formatRMB(row.sale_rmb_per_second, ratio)} / 秒`);
    }
    return parts.length > 0 ? parts.join('；') : '-';
  };

  const renderOfficialPrice = (row) => {
    const parts = [];
    if (row.official_usd_per_m_tokens) {
      parts.push(`$${row.official_usd_per_m_tokens} / 1M tokens`);
    }
    if (row.official_usd_per_video) {
      parts.push(`$${row.official_usd_per_video} / 条`);
    }
    if (row.official_usd_per_video_min || row.official_usd_per_video_max) {
      parts.push(
        `$${row.official_usd_per_video_min} - $${row.official_usd_per_video_max} / 条`,
      );
    }
    if (row.official_usd_per_second) {
      parts.push(`$${row.official_usd_per_second} / 秒`);
    }
    return parts.length > 0 ? parts.join('；') : '-';
  };

  const renderVideoPricingTable = () => {
    const rows = modelData?.video_pricing?.rows;
    if (!Array.isArray(rows) || rows.length === 0) return null;

    const displayGroup = resolveDisplayGroup();
    const tableData = rows.map((row, index) => ({
      key: `${row.resolution}-${row.scenario}-${index}`,
      resolution: row.resolution,
      scenario: row.scenario_label || row.scenario,
      official: renderOfficialPrice(row),
      sale: renderVideoSalePrice(row, displayGroup.ratio),
    }));

    const columns = [
      {
        title: t('分辨率'),
        dataIndex: 'resolution',
        render: (text) => (
          <Tag color='white' size='small' shape='circle'>
            {String(text).toUpperCase()}
          </Tag>
        ),
      },
      {
        title: t('场景'),
        dataIndex: 'scenario',
      },
      {
        title: t('官方价格'),
        dataIndex: 'official',
      },
      {
        title: `${t('售卖价')}（${displayGroup.group || t('默认')} ${displayGroup.ratio}x）`,
        dataIndex: 'sale',
        render: (text) => (
          <div className='font-semibold text-orange-600'>{text}</div>
        ),
      },
    ];

    return (
      <div className='mb-6'>
        <div className='mb-3'>
          <Text className='text-base font-medium'>{t('视频分辨率价格')}</Text>
          <div className='text-xs text-gray-600 mt-1'>
            {modelData.video_pricing.formula}
          </div>
        </div>
        <Table
          dataSource={tableData}
          columns={columns}
          pagination={false}
          size='small'
          bordered={false}
          className='!rounded-lg'
        />
      </div>
    );
  };

  const renderGroupPriceTable = () => {
    // 仅展示模型可用的分组：模型 enable_groups 与用户可用分组的交集

    const availableGroups = getAvailableGroups();

    // 准备表格数据
    const tableData = availableGroups.map((group) => {
      const priceData = modelData
        ? calculateModelPrice({
            record: modelData,
            selectedGroup: group,
            groupRatio,
            tokenUnit,
            displayPrice,
            currency,
            quotaDisplayType: siteDisplayType,
          })
        : { inputPrice: '-', outputPrice: '-', price: '-' };

      // 获取分组倍率
      const groupRatioValue =
        groupRatio && groupRatio[group] ? groupRatio[group] : 1;

      return {
        key: group,
        group: group,
        ratio: groupRatioValue,
        billingType:
          modelData?.quota_type === 0
            ? t('按量计费')
            : modelData?.quota_type === 1
              ? t('按次计费')
              : '-',
        priceItems: getModelPriceItems(priceData, t, siteDisplayType),
      };
    });

    // 定义表格列
    const columns = [
      {
        title: t('分组'),
        dataIndex: 'group',
        render: (text) => (
          <Tag color='white' size='small' shape='circle'>
            {text}
            {t('分组')}
          </Tag>
        ),
      },
    ];

    // 如果显示倍率，添加倍率列
    if (showRatio) {
      columns.push({
        title: t('倍率'),
        dataIndex: 'ratio',
        render: (text) => (
          <Tag color='white' size='small' shape='circle'>
            {text}x
          </Tag>
        ),
      });
    }

    // 添加计费类型列
    columns.push({
      title: t('计费类型'),
      dataIndex: 'billingType',
      render: (text) => {
        let color = 'white';
        if (text === t('按量计费')) color = 'violet';
        else if (text === t('按次计费')) color = 'teal';
        return (
          <Tag color={color} size='small' shape='circle'>
            {text || '-'}
          </Tag>
        );
      },
    });

    columns.push({
      title: siteDisplayType === 'TOKENS' ? t('计费摘要') : t('价格摘要'),
      dataIndex: 'priceItems',
      render: (items) => (
        <div className='space-y-1'>
          {items.map((item) => (
            <div key={item.key}>
              <div className='font-semibold text-orange-600'>
                {item.label} {item.value}
              </div>
              <div className='text-xs text-gray-500'>{item.suffix}</div>
            </div>
          ))}
        </div>
      ),
    });

    return (
      <Table
        dataSource={tableData}
        columns={columns}
        pagination={false}
        size='small'
        bordered={false}
        className='!rounded-lg'
      />
    );
  };

  return (
    <Card className='!rounded-2xl shadow-sm border-0'>
      <div className='flex items-center mb-4'>
        <Avatar size='small' color='orange' className='mr-2 shadow-md'>
          <IconCoinMoneyStroked size={16} />
        </Avatar>
        <div>
          <Text className='text-lg font-medium'>{t('分组价格')}</Text>
          <div className='text-xs text-gray-600'>
            {t('不同用户分组的价格信息')}
          </div>
        </div>
      </div>
      {autoChain.length > 0 && (
        <div className='flex flex-wrap items-center gap-1 mb-4'>
          <span className='text-sm text-gray-600'>{t('auto分组调用链路')}</span>
          <span className='text-sm'>→</span>
          {autoChain.map((g, idx) => (
            <React.Fragment key={g}>
              <Tag color='white' size='small' shape='circle'>
                {g}
                {t('分组')}
              </Tag>
              {idx < autoChain.length - 1 && <span className='text-sm'>→</span>}
            </React.Fragment>
          ))}
        </div>
      )}
      {renderVideoPricingTable()}
      {renderGroupPriceTable()}
    </Card>
  );
};

export default ModelPricingTable;
